package delivery

import (
	"fmt"
	"os"
	"strings"

	"github.com/mailhog/MailHog-MTA/backend/auth"
	"github.com/mailhog/MailHog-MTA/config"
	"github.com/mailhog/data"
)

// Service represents a delivery service implementation
type Service interface {
	Deliver(msg *data.Message) (id string, err error)
	WillDeliver(from, to string, as auth.Identity) bool
	MaxRecipients(as auth.Identity) int
}

// Load loads a delivery backend
func Load(cfg *config.Config, server *config.Server) Service {
	// FIXME delivery backend could be loaded multiple times, should cache this
	if a := server.Backends.Delivery; a != nil {
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
			return NewLocalDelivery(*a, *server, *cfg)
		default:
			fmt.Printf("Backend type not recognised\n")
			os.Exit(1)
		}
	}

	return nil
}
