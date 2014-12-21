package local

import (
	"code.google.com/p/go.crypto/bcrypt"
	"github.com/ian-kent/Go-MailHog/MailHog-MTA/config"
)

// Backend implements local disk storage for all backend services
type Backend struct {
	authMap map[string][]byte
	config  *config.Config
	server  *config.Server
}

// Configure implements Service.Configure
func (l *Backend) Configure(config *config.Config, server *config.Server) error {
	c, _ := bcrypt.GenerateFromPassword([]byte("test"), 11)
	l.authMap = map[string][]byte{
		"test@mailhog.example": c,
	}
	return nil
}
