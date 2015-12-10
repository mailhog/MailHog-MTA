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
	var err error

	if server.Backends.Auth != nil {
		a, err = server.Backends.Auth.Resolve(cfg.Backends)
		if err != nil {
			return err
		}
	}
	if server.Backends.Delivery != nil {
		d, err = server.Backends.Delivery.Resolve(cfg.Backends)
		if err != nil {
			return err
		}
	}
	if server.Backends.Resolver != nil {
		r, err = server.Backends.Resolver.Resolve(cfg.Backends)
		if err != nil {
			return err
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
