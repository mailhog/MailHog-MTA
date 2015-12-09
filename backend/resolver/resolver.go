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
	Resolve(address string) Result
}

// Result represents an address resolution result
type Result struct {
	Domain  DomainState
	Mailbox MailboxState
}

// DomainState is the result of a domain lookup
type DomainState uint8

// MailboxState is the result of a mailbox lookup
type MailboxState uint8

const (
	// DomainNotFound is returned for unknown domains.
	// This includes outbound SMTP to domains not at this host.
	DomainNotFound = DomainState(iota)
	// DomainPrimaryLocal is returned for local primary domains, i.e.
	// - domains this host is responsible for
	DomainPrimaryLocal
	// DomainSecondaryLocal is returned for secondary local domains, i.e.
	// - domains this host is a backup MX for
	// - domains this host acts as a inter-network router for, including
	//   private/public mail relaying
	DomainSecondaryLocal

	// MailboxLookupSkipped is returned when no local mailbox lookup is performed
	// e.g. for secondary local domains
	MailboxLookupSkipped = MailboxState(iota)
	// MailboxNotFound is returned when a lookup fails to locate a mailbox
	MailboxNotFound
	// MailboxFound is returned when a lookup finds a mailbox
	MailboxFound
)

// Load loads a resolver backend
func Load(cfg *config.Config, server *config.Server) Service {
	// FIXME resolver backend could be loaded multiple times, should cache this
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
			return NewLocalResolver(*a, *server, *cfg)
		default:
			fmt.Printf("Backend type not recognised\n")
			os.Exit(1)
		}
	}

	return nil
}
