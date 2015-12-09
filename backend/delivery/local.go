package delivery

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/mailhog/MailHog-MTA/backend/auth"
	"github.com/mailhog/MailHog-MTA/config"
	"github.com/mailhog/data"
)

// LocalDelivery implements delivery.Service
type LocalDelivery struct {
	spoolPath string
	server    config.Server
	app       config.Config
}

// NewLocalDelivery creates a new LocalDelivery backend
func NewLocalDelivery(cfg config.BackendConfig, srvCfg config.Server, appCfg config.Config) Service {
	spoolPath := os.TempDir()

	if c, ok := cfg.Data["spool_path"]; ok {
		if s, ok := c.(string); ok && len(s) > 0 {
			spoolPath = s
		}
	}

	return &LocalDelivery{
		server:    srvCfg,
		app:       appCfg,
		spoolPath: spoolPath,
	}
}

// Deliver implements DeliveryService.Deliver
func (l *LocalDelivery) Deliver(msg *data.Message) (id string, err error) {
	var mid data.MessageID

	// FIXME should use server hostname
	// FIXME also, this is for storage, so isn't strictly the "Message-ID"
	// as defined by the message header, or what the data.NewMessageID function
	// was intended for.
	mid, err = data.NewMessageID("mailhog.example")
	if err != nil {
		return
	}
	id = string(mid)

	dp := l.spoolPath
	if !strings.HasPrefix(dp, "/") {
		dp = filepath.Join(l.app.RelPath(), dp)
	}

	os.MkdirAll(dp, 0660)

	dp = filepath.Join(dp, id)

	b, err := ioutil.ReadAll(msg.Raw.Bytes())
	if err != nil {
		return "", err
	}

	err = ioutil.WriteFile(dp, b, 0660)

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
