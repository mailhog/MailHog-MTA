package smtp

import (
	"crypto/tls"
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

	tlsConfig *tls.Config
}

func (s *Server) getTLSConfig() *tls.Config {
	if s.tlsConfig != nil {
		return s.tlsConfig
	}
	cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	if err != nil {
		log.Fatal(err)
	}
	s.tlsConfig = &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	return s.tlsConfig
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

	sem := make(chan int, s.PolicySet.MaximumConnections)

	for {
		sem <- 1

		conn, err := ln.Accept()
		if err != nil {
			log.Printf("[SMTP] Error accepting connection: %s\n", err)
			continue
		}

		go func() {
			s.Accept(
				conn.(*net.TCPConn).RemoteAddr().String(),
				io.ReadWriteCloser(conn),
			)

			<-sem
		}()
	}
}
