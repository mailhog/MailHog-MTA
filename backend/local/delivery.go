package local

import (
	"github.com/mailhog/MailHog-MTA/backend"
	"github.com/mailhog/data"
)

// Deliver implements DeliveryService.Deliver
func (l *Backend) Deliver(msg *data.Message) (id string, err error) {
	var mid data.MessageID

	mid, err = data.NewMessageID("mailhog.example")
	if err != nil {
		return
	}
	id = string(mid)

	l.deliveryQueue = append(l.deliveryQueue, msg)

	return
}

// WillDeliver implements DeliveryService.WillDeliver
func (l *Backend) WillDeliver(from, to string, as *backend.Identity) bool {
	return true
}

// MaxRecipients implements DeliveryService.MaxRecipients
func (l *Backend) MaxRecipients(as *backend.Identity) int {
	return l.server.PolicySet.MaximumRecipients
}
