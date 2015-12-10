package main

import (
	"flag"
	"log"
	"sync"

	"github.com/mailhog/MailHog-MTA/config"
	"github.com/mailhog/MailHog-MTA/smtp"
	"github.com/mailhog/backends/auth"
	sconfig "github.com/mailhog/backends/config"
	"github.com/mailhog/backends/delivery"
	"github.com/mailhog/backends/resolver"
)

var conf *config.Config
var wg sync.WaitGroup

func configure() {
	config.RegisterFlags()
	flag.Parse()
	conf = config.Configure()
}

func main() {
	configure()

	for _, s := range conf.Servers {
		wg.Add(1)
		go func(s *config.Server) {
			defer wg.Done()
			err := newServer(conf, s)
			if err != nil {
				log.Fatal(err)
			}
		}(s)
	}

	wg.Wait()
}

func newServer(cfg *config.Config, server *config.Server) error {
	var a, d, r sconfig.BackendConfig

	if server.Backends.Auth != nil {
		a = *server.Backends.Auth
		if len(server.Backends.Auth.Ref) > 0 {
			if _, ok := cfg.Backends[server.Backends.Auth.Ref]; ok {
				a = cfg.Backends[server.Backends.Auth.Ref]
			}
		}
	}

	if server.Backends.Delivery != nil {
		d = *server.Backends.Delivery
		if len(server.Backends.Delivery.Ref) > 0 {
			if _, ok := cfg.Backends[server.Backends.Delivery.Ref]; ok {
				d = cfg.Backends[server.Backends.Delivery.Ref]
			}
		}
	}

	if server.Backends.Resolver != nil {
		r = *server.Backends.Resolver
		if len(server.Backends.Resolver.Ref) > 0 {
			if _, ok := cfg.Backends[server.Backends.Resolver.Ref]; ok {
				r = cfg.Backends[server.Backends.Resolver.Ref]
			}
		}
	}

	s := &smtp.Server{
		BindAddr:        server.BindAddr,
		Hostname:        server.Hostname,
		PolicySet:       server.PolicySet,
		AuthBackend:     auth.Load(a, *cfg),
		DeliveryBackend: delivery.Load(d, *cfg),
		ResolverBackend: resolver.Load(r, *cfg),
		Config:          cfg,
		Server:          server,
	}

	return s.Listen()
}
