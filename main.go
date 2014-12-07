package main

import (
	"flag"
	"os"

	"github.com/ian-kent/Go-MailHog/MailHog-MTA/config"
	"github.com/ian-kent/Go-MailHog/MailHog-MTA/smtp"
	"github.com/ian-kent/go-log/log"
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

	go smtp.Listen(conf, exitCh)

	for {
		select {
		case <-exitCh:
			log.Printf("Received exit signal")
			os.Exit(0)
		}
	}
}
