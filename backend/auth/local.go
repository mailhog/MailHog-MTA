package auth

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"code.google.com/p/go.crypto/bcrypt"

	"github.com/mailhog/MailHog-MTA/config"
	"github.com/mailhog/smtp"
)

var localMechanisms = []string{"PLAIN"}

// LocalAuth implements auth.Service
type LocalAuth struct {
	authMap map[string]*LocalUser
}

// NewLocalAuth returns a new LocalAuth using the provided map
func NewLocalAuth(cfg config.BackendConfig, srvCfg config.Server, appCfg config.Config) *LocalAuth {
	var authMap map[string]*LocalUser

	if c, ok := cfg.Data["config"]; ok {
		if s, ok := c.(string); ok && len(s) > 0 {
			log.Printf("loading auth data from: %s", s)
			if !strings.HasPrefix(s, "/") {
				s = filepath.Join(appCfg.RelPath(), s)
			}

			b, err := ioutil.ReadFile(s)
			if err != nil {
				log.Fatal(err)
			}
			err = json.Unmarshal(b, &authMap)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	return &LocalAuth{
		authMap: authMap,
	}
}

// LocalUser represents a virtual user
type LocalUser struct {
	Username     string
	Password     []byte
	ValidSenders []string
}

func (l LocalUser) String() string {
	return l.Username
}

// IsValidSender implements Identity.IsValidSender
func (l LocalUser) IsValidSender(sender string) bool {
	for _, s := range l.ValidSenders {
		if s == sender {
			return true
		}
	}
	return false
}

// TODO abstract away password mechanism and identity retrieval

// Authenticate implements AuthService.Authenticate
func (l *LocalAuth) Authenticate(mechanism string, args ...string) (identity Identity, errorReply *smtp.Reply, ok bool) {
	log.Println(mechanism)
	log.Println(args)

	if len(args) < 2 {
		errorReply = smtp.ReplyError(errors.New("Missing arguments"))
		ok = false
		return
	}

	user, pass := args[0], args[1]

	if u, k := l.authMap[user]; k {
		err := bcrypt.CompareHashAndPassword(u.Password, []byte(pass))

		if err != nil {
			// FIXME
			errorReply = smtp.ReplyError(errors.New("Invalid password"))
			ok = false
			return
		}
		// FIXME
		id := Identity(LocalUser{user, []byte{}, []string{user}})
		identity = id
		ok = true
		return
	}

	// FIXME
	errorReply = smtp.ReplyError(errors.New("User not found"))
	ok = false
	return
}

// Mechanisms implements AuthService.Mechanisms
func (l *LocalAuth) Mechanisms() []string {
	return localMechanisms
}
