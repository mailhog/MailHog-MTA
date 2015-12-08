package delivery

import (
	"fmt"
	"os"

	"github.com/mailhog/MailHog-MTA/backend"
	"github.com/mailhog/MailHog-MTA/config"
	"github.com/mailhog/data"
)

// Service represents a delivery service implementation
type Service interface {
	backend.Service
	Deliver(msg *data.Message) (id string, err error)
	WillDeliver(from, to string, as *backend.Identity) bool
	MaxRecipients(as *backend.Identity) int
}

// Load loads a delivery backend
func Load(cfg *config.Config, server *config.Server) Service {
	if a := server.Backends.Delivery; a != nil {
		if len(a.Ref) > 0 {
			if a2, ok := cfg.Backends[a.Ref]; ok {
				a = &a2
			} else {
				fmt.Printf("Backend not found: %s\n", a.Ref)
				os.Exit(1)
			}
		}

		localDelivery := &LocalDelivery{}
		localDelivery.Configure(cfg, server)
		return localDelivery
	}

	return nil
}
