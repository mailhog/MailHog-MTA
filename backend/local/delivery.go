package local

import (
	"github.com/ian-kent/Go-MailHog/MailHog-MTA/backend"
	"github.com/ian-kent/Go-MailHog/data"
)

// Deliver implements DeliveryService.Deliver
func (l *Backend) Deliver(msg *data.Message) (id string, err error) {
	return
}

// WillDeliver implements DeliveryService.WillDeliver
func (l *Backend) WillDeliver(from, to string, as *backend.Identity) bool {
	return true
}
