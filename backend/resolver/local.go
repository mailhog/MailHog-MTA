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
	resolveMap map[string]map[string]ResolvedState
}

// NewLocalResolver returns a new local resolver using the provided map
func NewLocalResolver(cfg config.BackendConfig, srvCfg config.Server, appCfg config.Config) *LocalResolver {
	var resolveMap map[string]map[string]ResolvedState

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
