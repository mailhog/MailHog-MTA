package auth

import (
	"errors"
	"log"

	"code.google.com/p/go.crypto/bcrypt"
	"github.com/mailhog/MailHog-MTA/backend"
	"github.com/mailhog/smtp"
)

var mechanisms = []string{"PLAIN"}

// LocalAuth implements auth.Service
type LocalAuth struct {
	backend.DefaultBackend
	authMap map[string]*LocalUser
}

// NewLocalAuth returns a new LocalAuth using the provided map
func NewLocalAuth(m map[string]*LocalUser) *LocalAuth {
	return &LocalAuth{
		authMap: m,
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
func (l *LocalAuth) Authenticate(mechanism string, args ...string) (identity *backend.Identity, errorReply *smtp.Reply, ok bool) {
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
		id := backend.Identity(LocalUser{user, []byte{}, []string{user}})
		identity = &id
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
	return mechanisms
}
