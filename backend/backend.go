package backend

import (
	"github.com/mailhog/MailHog-MTA/config"
	"github.com/mailhog/data"
	"github.com/mailhog/smtp"
)

// Identity represents an identity
type Identity interface {
	String() string
	IsValidSender(string) bool
}

// Service represents a service implementation
//
// Combined service implementations should not assume that
// all individual service components will be used.
//
// Configure will only be called once for each distinct Go value.
// E.g., using a combined service implementation which provides
// multiple services will only have its Configure function called once.
type Service interface {
	Configure(*config.Config, *config.Server) error
}

// AuthService represents an authentication service implementation
type AuthService interface {
	Service
	Authenticate(mechanism string, args ...string) (identity *Identity, errorReply *smtp.Reply, ok bool)
	Mechanisms() []string
}

// DeliveryService represents a delivery service implementation
type DeliveryService interface {
	Service
	Deliver(msg *data.Message) (id string, err error)
	WillDeliver(from, to string, as *Identity) bool
	MaxRecipients(as *Identity) int
}

// ResolverService represents an address resolver implementation
// FIXME what this all actually means is "will you accept messages for this address"
// FIXME and the only responses are: yes, no
// FIXME if yes, the reasons can be: i own the mailbox, i can deliver to the mailbox, i'll relay mail for you
// FIXME it might be clearer to use that terminology?
type ResolverService interface {
	Service
	Resolve(address string) ResolvedState
}

// ResolvedState represents the resolved state of an address
type ResolvedState uint8

const (
	// ResolvedNotFound is returned for non-existant mailboxes at local domains
	ResolvedNotFound = ResolvedState(iota)
	// ResolvedPrimaryLocal is returned for mailboxes at local primary domains
	ResolvedPrimaryLocal
	// ResolvedSecondaryLocal is returned for mailboxes at local secondary domains
	ResolvedSecondaryLocal
	// ResolvedRemote is returned for mailboxes at unrecognised domains
	ResolvedRemote
)
