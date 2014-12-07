package local

import (
	"errors"
	"log"

	"code.google.com/p/go.crypto/bcrypt"
	"github.com/ian-kent/Go-MailHog/MailHog-MTA/backend"
	"github.com/ian-kent/Go-MailHog/smtp/protocol"
)

var mechanisms = []string{"PLAIN"}

// TODO abstract away password mechanism and identity retrieval

// Authenticate implements AuthService.Authenticate
func (l *Backend) Authenticate(mechanism string, args ...string) (identity *backend.Identity, errorReply *protocol.Reply, ok bool) {
	log.Println(mechanism)
	log.Println(args)

	if len(args) < 2 {
		errorReply = protocol.ReplyError(errors.New("Missing arguments"))
		ok = false
		return
	}

	user, pass := args[0], args[1]

	if pw, k := l.authMap[user]; k {
		err := bcrypt.CompareHashAndPassword(pw, []byte(pass))

		if err != nil {
			// FIXME
			errorReply = protocol.ReplyError(errors.New("Invalid password"))
			ok = false
			return
		}
		id := backend.Identity(backend.UserIdentity{user})
		identity = &id
		ok = true
		return
	}

	// FIXME
	errorReply = protocol.ReplyError(errors.New("User not found"))
	ok = false
	return
}

// Mechanisms implements AuthService.Mechanisms
func (l *Backend) Mechanisms() []string {
	return mechanisms
}
