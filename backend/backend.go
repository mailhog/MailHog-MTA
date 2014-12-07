package backend

import (
	"github.com/ian-kent/Go-MailHog/data"
	"github.com/ian-kent/Go-MailHog/smtp/protocol"
)

// Identity represents an identity
type Identity interface {
	String() string
}

// UserIdentity represents a users identity
type UserIdentity struct {
	Username string
}

// String implements Identity.String
func (i UserIdentity) String() string {
	return i.Username
}

// Service represents a service implementation
type Service interface {
	Configure(map[string]string) error
}

// AuthService represents an authentication service implementation
type AuthService interface {
	Service
	Authenticate(mechanism string, args ...string) (identity *Identity, errorReply *protocol.Reply, ok bool)
	Mechanisms() []string
}

// DeliveryService represents a delivery service implementation
type DeliveryService interface {
	Service
	Deliver(msg *data.Message) (id string, err error)
	WillDeliver(from, to string, as *Identity) bool
}
