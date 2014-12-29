package smtp

import (
	"io"
	"log"
	"net"

	"github.com/mailhog/MailHog-MTA/backend"
	"github.com/mailhog/MailHog-MTA/config"
)

// Server represents an SMTP server instance
type Server struct {
	BindAddr  string
	Hostname  string
	PolicySet config.PolicySet

	AuthBackend     backend.AuthService
	DeliveryBackend backend.DeliveryService
	ResolverBackend backend.ResolverService
}

// Listen starts listening on the configured bind address
func (s *Server) Listen() error {
	log.Printf("[SMTP] Binding to address: %s\n", s.BindAddr)
	ln, err := net.Listen("tcp", s.BindAddr)
	if err != nil {
		log.Fatalf("[SMTP] Error listening on socket: %s\n", err)
		return err
	}

	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("[SMTP] Error accepting connection: %s\n", err)
			continue
		}
		defer conn.Close()

		go s.Accept(
			conn.(*net.TCPConn).RemoteAddr().String(),
			io.ReadWriteCloser(conn),
		)
	}
}
