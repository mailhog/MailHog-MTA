package main

import (
	"flag"
	"sync"

	"github.com/ian-kent/Go-MailHog/MailHog-MTA/config"
	"github.com/ian-kent/Go-MailHog/MailHog-MTA/smtp"
)

var conf *config.Config
var exitCh chan int

func configure() {
	config.RegisterFlags()
	flag.Parse()
	conf = config.Configure()
}

func main() {
	configure()

	exitCh = make(chan int)

	var wg sync.WaitGroup
	for _, s := range conf.Servers {
		wg.Add(1)
		go func(s *config.Server) {
			defer wg.Done()
			smtp.Listen(conf, s, exitCh)
		}(s)
	}
	wg.Wait()
}
