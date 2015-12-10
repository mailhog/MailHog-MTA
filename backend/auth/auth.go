package auth

import (
	"fmt"
	"os"
	"strings"

	"github.com/mailhog/MailHog-MTA/config"
	"github.com/mailhog/smtp"
)

/*
  FIXME

  Consider whether mechanisms are defined per-backend or not.

  - Are all (available) mechanisms automatically supported by all backends, e.g. EXTERNAL?
  - Or, are the mechanisms supported specific to a particular backend?
  - Parsing is done in mailhog/smtp, does that make mechanism support a policy decision?
*/

// Service represents an authentication service implementation
type Service interface {
	Authenticate(mechanism string, args ...string) (identity Identity, errorReply *smtp.Reply, ok bool)
	Mechanisms() []string
}

// Identity represents an identity
type Identity interface {
	String() string
	IsValidSender(string) bool
	PolicySet() config.IdentityPolicySet
}

// Load loads an auth backend
func Load(cfg *config.Config, server *config.Server) Service {
	// FIXME auth backend could be loaded multiple times, should cache this
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
			return NewLocalAuth(*a, *server, *cfg)
		default:
			fmt.Printf("Backend type not recognised\n")
			os.Exit(1)
		}
	}

	return nil
}
