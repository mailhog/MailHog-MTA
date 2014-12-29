package local

import (
	"code.google.com/p/go.crypto/bcrypt"
	"github.com/mailhog/MailHog-MTA/backend"
	"github.com/mailhog/MailHog-MTA/config"
)

// Backend implements local disk storage for all backend services
type Backend struct {
	authMap    map[string]*BackendUser
	resolveMap map[string]map[string]backend.ResolvedState
	config     *config.Config
	server     *config.Server
}

type BackendUser struct {
	Username     string
	Password     []byte
	ValidSenders []string
}

// Configure implements Service.Configure
func (l *Backend) Configure(config *config.Config, server *config.Server) error {
	l.server = server
	l.config = config

	c, _ := bcrypt.GenerateFromPassword([]byte("test"), 11)
	l.authMap = map[string]*BackendUser{
		"test@mailhog.example": &BackendUser{
			Username:     "test@mailhog.example",
			Password:     c,
			ValidSenders: []string{"test@mailhog.example", "alias@mailhog.example"},
		},
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
