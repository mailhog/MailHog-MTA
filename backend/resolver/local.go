package resolver

import (
	"github.com/mailhog/MailHog-MTA/backend"
	"github.com/mailhog/data"
)

// LocalResolver implements resolver.Service
type LocalResolver struct {
	backend.DefaultBackend
	resolveMap map[string]map[string]ResolvedState
}

// NewLocalResolver returns a new local resolver using the provided map
func NewLocalResolver(m map[string]map[string]ResolvedState) *LocalResolver {
	return &LocalResolver{
		resolveMap: m,
	}
}

// Resolve implements ResolverService.Resolve
func (l *LocalResolver) Resolve(address string) (ResolvedState, DeliveryState) {
	path := data.PathFromString(address)

	if m, ok := l.resolveMap[path.Domain]; ok {
		if s, ok := m[path.Mailbox]; ok {
			if s == ResolvedPrimaryLocal {
				return s, DeliveryDirect
			}
			if s == ResolvedSecondaryLocal || s == ResolvedRemote {
				return s, DeliveryRelay
			}
		}
		return ResolvedNotFound, DeliveryRejected
	}

	return ResolvedRemote, DeliveryRejected
}
