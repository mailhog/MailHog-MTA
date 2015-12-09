package resolver

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/mailhog/MailHog-MTA/config"
	"github.com/mailhog/data"
)

// LocalResolver implements resolver.Service
type LocalResolver struct {
	resolveMap map[string]domain
}

type domain struct {
	Name      string
	State     DomainState
	Mailboxes map[string]mailbox `json:",omitempty"`
}

type mailbox struct {
	Name  string
	State MailboxState
}

// NewLocalResolver returns a new local resolver using the provided map
func NewLocalResolver(cfg config.BackendConfig, srvCfg config.Server, appCfg config.Config) *LocalResolver {
	var resolveMap map[string]domain

	if c, ok := cfg.Data["config"]; ok {
		if s, ok := c.(string); ok && len(s) > 0 {
			if !strings.HasPrefix(s, "/") {
				s = filepath.Join(appCfg.RelPath(), s)
			}
			b, err := ioutil.ReadFile(s)
			if err != nil {
				log.Fatal(err)
			}
			err = json.Unmarshal(b, &resolveMap)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	return &LocalResolver{
		resolveMap: resolveMap,
	}
}

// Resolve implements ResolverService.Resolve
func (l *LocalResolver) Resolve(address string) (r Result) {
	path := data.PathFromString(address)

	log.Printf("resolving: %s", path)

	if m, ok := l.resolveMap[path.Domain]; ok {
		log.Printf("found domain: %s", path.Domain)
		r.Domain = m.State

		if s, ok := m.Mailboxes[path.Mailbox]; ok {
			log.Printf("found mailbox: %s [%d]", path.Mailbox, s.State)
			r.Mailbox = s.State
			return
		}

		log.Printf("mailbox doesn't exist at local domain: %s", path.Mailbox)
		return
	}

	log.Printf("not a local address")
	return
}
