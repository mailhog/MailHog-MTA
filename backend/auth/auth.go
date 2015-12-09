package auth

import (
	"fmt"
	"os"
	"strings"

	"github.com/mailhog/MailHog-MTA/config"
	"github.com/mailhog/smtp"
)

// Service represents an authentication service implementation
type Service interface {
	Authenticate(mechanism string, args ...string) (identity Identity, errorReply *smtp.Reply, ok bool)
	Mechanisms() []string
}

// Identity represents an identity
type Identity interface {
	String() string
	IsValidSender(string) bool
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

		switch strings.ToLower(a.Type) {
		case "local":
			return NewLocalAuth(*server.Backends.Auth, *server, *cfg)
		default:
			fmt.Printf("Backend type not recognised\n")
			os.Exit(1)
		}
	}

	return nil
}
