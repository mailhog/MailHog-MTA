package local

import (
	"code.google.com/p/go.crypto/bcrypt"
	"github.com/ian-kent/Go-MailHog/MailHog-MTA/backend"
	"github.com/ian-kent/Go-MailHog/MailHog-MTA/config"
)

// Backend implements local disk storage for all backend services
type Backend struct {
	authMap    map[string][]byte
	resolveMap map[string]map[string]backend.ResolvedState
	config     *config.Config
	server     *config.Server
}

// Configure implements Service.Configure
func (l *Backend) Configure(config *config.Config, server *config.Server) error {
	c, _ := bcrypt.GenerateFromPassword([]byte("test"), 11)
	l.authMap = map[string][]byte{
		"test@mailhog.example": c,
	}
	l.resolveMap = map[string]map[string]backend.ResolvedState{
		"mailhog.example": map[string]backend.ResolvedState{
			"test": backend.ResolvedPrimaryLocal,
		},
		"mailhog.internal": map[string]backend.ResolvedState{
			"test": backend.ResolvedSecondaryLocal,
		},
		"mailhog.remote": map[string]backend.ResolvedState{
			"test": backend.ResolvedRemote,
		},
	}
	return nil
}
