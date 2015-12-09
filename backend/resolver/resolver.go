package resolver

import (
	"fmt"
	"os"
	"strings"

	"github.com/mailhog/MailHog-MTA/config"
)

// Service represents an address resolver implementation
// FIXME what this all actually means is "will you accept messages for this address"
// FIXME and the only responses are: yes, no
// FIXME if yes, the reasons can be: i own the mailbox, i can deliver to the mailbox, i'll relay mail for you
// FIXME it might be clearer to use that terminology?
type Service interface {
	Resolve(address string) (ResolvedState, DeliveryState)
}

// ResolvedState represents the resolved state of an address
type ResolvedState uint8

// DeliveryState represents the deliverability of the address
type DeliveryState uint8

const (
	// ResolvedNotFound is returned for non-existant mailboxes at local domains
	ResolvedNotFound = ResolvedState(iota)
	// ResolvedPrimaryLocal is returned for mailboxes at local primary domains
	ResolvedPrimaryLocal
	// ResolvedSecondaryLocal is returned for mailboxes at local secondary domains
	ResolvedSecondaryLocal
	// ResolvedRemote is returned for mailboxes at unrecognised domains
	ResolvedRemote

	// DeliveryRejected is returned if delivery to the mailbox is not possible
	DeliveryRejected = DeliveryState(iota)
	// DeliveryDirect is returned if direct delivery to the mailbox is possible
	DeliveryDirect
	// DeliveryRelay is returned if relay delivery to the mailbox is possible
	DeliveryRelay
)

// Load loads a resolver backend
func Load(cfg *config.Config, server *config.Server) Service {
	if a := server.Backends.Resolver; a != nil {
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
			return NewLocalResolver(*server.Backends.Resolver, *server, *cfg)
		default:
			fmt.Printf("Backend type not recognised\n")
			os.Exit(1)
		}
	}

	return nil
}
