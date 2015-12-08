package backend

import "github.com/mailhog/MailHog-MTA/config"

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

var _ Service = &DefaultBackend{}

// DefaultBackend is a default struct to hold the current config and server
type DefaultBackend struct {
	Config *config.Config
	Server *config.Server
}

// Configure implements Service.Configure
func (l *DefaultBackend) Configure(config *config.Config, server *config.Server) error {
	l.Server = server
	l.Config = config
	return nil
}
