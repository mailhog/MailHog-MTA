package main

import (
	"flag"
	"log"
	"sync"

	"github.com/ian-kent/Go-MailHog/MailHog-MTA/backend/local"
	"github.com/ian-kent/Go-MailHog/MailHog-MTA/config"
	"github.com/ian-kent/Go-MailHog/MailHog-MTA/smtp"
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
	// FIXME make configurable
	localBackend := &local.Backend{}
	localBackend.Configure(cfg, server)

	s := &smtp.Server{
		BindAddr:        server.BindAddr,
		Hostname:        server.Hostname,
		PolicySet:       server.PolicySet,
		AuthBackend:     localBackend,
		DeliveryBackend: localBackend,
	}

	return s.Listen()
}
