package smtp

import (
	"io"
	"log"
	"net"

	"github.com/ian-kent/Go-MailHog/MailHog-MTA/config"
)

func ListenSMTP(cfg *config.Config, exitCh chan int) *net.TCPListener {
	log.Printf("[SMTP] Binding to address: %s\n", cfg.SMTPBindAddr)
	ln, err := net.Listen("tcp", cfg.SMTPBindAddr)
	if err != nil {
		log.Fatalf("[SMTP] Error listening on socket: %s\n", err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("[SMTP] Error accepting connection: %s\n", err)
			continue
		}
		defer conn.Close()

		go Accept(
			conn.(*net.TCPConn).RemoteAddr().String(),
			io.ReadWriteCloser(conn),
			cfg.Hostname,
			false,
		)
	}
}

func ListenSubmission(cfg *config.Config, exitCh chan int) *net.TCPListener {
	log.Printf("[SUBMISSION] Binding to address: %s\n", cfg.SubmissionBindAddr)
	ln, err := net.Listen("tcp", cfg.SubmissionBindAddr)
	if err != nil {
		log.Fatalf("[SUBMISSION] Error listening on socket: %s\n", err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("[SUBMISSION] Error accepting connection: %s\n", err)
			continue
		}
		defer conn.Close()

		go Accept(
			conn.(*net.TCPConn).RemoteAddr().String(),
			io.ReadWriteCloser(conn),
			cfg.Hostname,
			true,
		)
	}
}
