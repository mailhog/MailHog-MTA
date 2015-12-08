package auth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/mailhog/MailHog-MTA/backend"
	"github.com/mailhog/MailHog-MTA/config"
	"github.com/mailhog/smtp"
)

// Service represents an authentication service implementation
type Service interface {
	backend.Service
	Authenticate(mechanism string, args ...string) (identity *backend.Identity, errorReply *smtp.Reply, ok bool)
	Mechanisms() []string
}

// Load loads an auth backend
func Load(cfg *config.Config, server *config.Server) Service {
	if a := server.Backends.Auth; a != nil {
		if len(a.Ref) > 0 {
			if a2, ok := cfg.Backends[a.Ref]; ok {
				a = &a2
			} else {
				fmt.Printf("Backend not found: %s\n", a.Ref)
				os.Exit(1)
			}
		}

		var authMap map[string]*LocalUser

		if c, ok := a.Data["config"]; ok {
			if s, ok := c.(string); ok && len(s) > 0 {
				if !strings.HasPrefix(s, "/") {
					s = filepath.Join(cfg.RelPath(), s)
				}

				b, err := ioutil.ReadFile(s)
				if err != nil {
					log.Fatal(err)
				}
				err = json.Unmarshal(b, &authMap)
				if err != nil {
					log.Fatal(err)
				}
			}
		}

		localAuth := NewLocalAuth(authMap)
		localAuth.Configure(cfg, server)
		return localAuth
	}

	return nil
}
