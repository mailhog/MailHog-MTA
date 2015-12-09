package delivery

import (
	"github.com/mailhog/MailHog-MTA/backend/auth"
	"github.com/mailhog/MailHog-MTA/config"
	"github.com/mailhog/data"
)

// LocalDelivery implements delivery.Service
type LocalDelivery struct {
	deliveryQueue []*data.Message
	server        config.Server
}

// NewLocalDelivery creates a new LocalDelivery backend
func NewLocalDelivery(cfg config.BackendConfig, srvCfg config.Server, appCfg config.Config) Service {
	return &LocalDelivery{
		server: srvCfg,
	}
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
func (l *LocalDelivery) WillDeliver(from, to string, as auth.Identity) bool {
	return true
}

// MaxRecipients implements DeliveryService.MaxRecipients
func (l *LocalDelivery) MaxRecipients(as auth.Identity) int {
	return l.server.PolicySet.MaximumRecipients
}
