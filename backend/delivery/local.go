package delivery

import (
	"github.com/mailhog/MailHog-MTA/backend"
	"github.com/mailhog/MailHog-MTA/config"
	"github.com/mailhog/data"
)

// LocalDelivery implements delivery.Service
type LocalDelivery struct {
	backend.DefaultBackend
	deliveryQueue []*data.Message
}

// Configure implements backend.Service.Configure
func (l *LocalDelivery) Configure(config *config.Config, server *config.Server) error {
	return l.DefaultBackend.Configure(config, server)
}

// Deliver implements DeliveryService.Deliver
func (l *LocalDelivery) Deliver(msg *data.Message) (id string, err error) {
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
func (l *LocalDelivery) WillDeliver(from, to string, as *backend.Identity) bool {
	return true
}

// MaxRecipients implements DeliveryService.MaxRecipients
func (l *LocalDelivery) MaxRecipients(as *backend.Identity) int {
	return l.DefaultBackend.Server.PolicySet.MaximumRecipients
}
