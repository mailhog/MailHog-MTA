package smtp

import (
	"io"
	"log"
	"net"

	"github.com/ian-kent/Go-MailHog/MailHog-MTA/config"
)

func Listen(cfg *config.Config, server *config.Server, exitCh chan int) *net.TCPListener {
	log.Printf("[SMTP] Binding to address: %s\n", server.BindAddr)
	ln, err := net.Listen("tcp", server.BindAddr)
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
			server.Hostname,
			cfg,
			server,
		)
	}
}
