package local

import "code.google.com/p/go.crypto/bcrypt"

// Backend implements local disk storage for all backend services
type Backend struct {
	authMap map[string][]byte
}

// Configure implements Service.Configure
func (l *Backend) Configure(args map[string]string) error {
	c, _ := bcrypt.GenerateFromPassword([]byte("test"), 11)
	l.authMap = map[string][]byte{
		"test@mailhog.example": c,
	}
	return nil
}
