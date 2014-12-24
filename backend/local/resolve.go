package local

import (
	"github.com/mailhog/MailHog-MTA/backend"
	"github.com/mailhog/data"
)

// Resolve implements ResolverService.Resolve
func (l *Backend) Resolve(address string) (backend.ResolvedState, error) {
	path := data.PathFromString(address)

	if m, ok := l.resolveMap[path.Domain]; ok {
		if s, ok := m[path.Mailbox]; ok {
			return s, nil
		}
		return backend.ResolvedNotFound, nil
	}

	return backend.ResolvedRemote, nil
}
