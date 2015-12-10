package smtp

// http://www.rfc-editor.org/rfc/rfc5321.txt

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"strings"

	"github.com/mailhog/backends/auth"
	"github.com/mailhog/backends/resolver"
	"github.com/mailhog/data"
	"github.com/mailhog/smtp"
)

// Session represents a SMTP session using net.TCPConn
type Session struct {
	server *Server

	conn          io.ReadWriteCloser
	proto         *smtp.Protocol
	remoteAddress string
	isTLS         bool
	line          string
	identity      auth.Identity

	maximumBufferLength int
}

// Accept starts a new SMTP session using io.ReadWriteCloser
func (s *Server) Accept(remoteAddress string, conn io.ReadWriteCloser) {
	proto := smtp.NewProtocol()
	proto.Hostname = s.Hostname

	session := &Session{
		server:              s,
		conn:                conn,
		proto:               proto,
		remoteAddress:       remoteAddress,
		isTLS:               false,
		line:                "",
		identity:            nil,
		maximumBufferLength: 2048000,
	}

	// FIXME this all feels nasty
	proto.LogHandler = session.logf
	proto.MessageReceivedHandler = session.acceptMessage
	proto.ValidateSenderHandler = session.validateSender
	proto.ValidateRecipientHandler = session.validateRecipient
	proto.ValidateAuthenticationHandler = session.validateAuthentication
	if session.server != nil && session.server.AuthBackend != nil {
		proto.GetAuthenticationMechanismsHandler = session.server.AuthBackend.Mechanisms
	}
	proto.SMTPVerbFilter = session.verbFilter
	proto.MaximumRecipients = session.server.PolicySet.MaximumRecipients
	proto.MaximumLineLength = session.server.PolicySet.MaximumLineLength

	if !session.server.PolicySet.DisableTLS {
		session.logf("Enabling TLS support")
		proto.TLSHandler = session.tlsHandler
		proto.RequireTLS = session.server.PolicySet.RequireTLS
	}

	session.logf("Starting session")
	session.Write(proto.Start())
	for session.Read() == true {
	}
	io.Closer(conn).Close()
	session.logf("Session ended")
}

func (c *Session) validateAuthentication(mechanism string, args ...string) (errorReply *smtp.Reply, ok bool) {
	if c.server.AuthBackend == nil {
		return smtp.ReplyInvalidAuth(), false
	}
	i, e, ok := c.server.AuthBackend.Authenticate(mechanism, args...)
	if e != nil || !ok {
		if e != nil {
			c.logf("error authenticating: %s", e)
		}
		return smtp.ReplyInvalidAuth(), false
	}
	c.identity = i
	return nil, true
}

func (c *Session) validateRecipient(to string) bool {
	if c.server.DeliveryBackend == nil {
		return false
	}

	maxRecipients := c.server.PolicySet.MaximumRecipients
	if maxRecipients > -1 && len(c.proto.Message.To) > maxRecipients {
		return false
	}

	if c.identity != nil &&
		c.identity.PolicySet().MaximumRecipients != nil &&
		*c.identity.PolicySet().MaximumRecipients <= len(c.proto.Message.To) {
		return false
	}

	r := c.server.ResolverBackend.Resolve(to)

	if c.server.PolicySet.RequireLocalDelivery {
		if r.Domain == resolver.DomainNotFound {
			return false
		}
	}

	if r.Domain == resolver.DomainPrimaryLocal &&
		r.Mailbox != resolver.MailboxFound &&
		(c.server.PolicySet.RejectInvalidRecipients ||
			(c.identity != nil &&
				c.identity.PolicySet().RejectInvalidRecipients != nil &&
				*c.identity.PolicySet().RejectInvalidRecipients)) {
		return false
	}

	return c.server.DeliveryBackend.WillDeliver(to, c.proto.Message.From, c.identity)
}

func (c *Session) validateSender(from string) bool {
	// we have a user (authenticated outbound SMTP)
	if c.identity != nil {
		return c.identity.IsValidSender(from)
	}

	// we don't, but we should (unauthenticated outbound SMTP)
	if c.server.PolicySet.RequireAuthentication {
		return false
	}

	// we don't, but we don't care (inbound SMTP)
	return true
}

func (c *Session) verbFilter(verb string, args ...string) (errorReply *smtp.Reply) {
	if c.server.PolicySet.RequireAuthentication && c.identity == nil {
		verb = strings.ToUpper(verb)
		if verb == "RSET" || verb == "QUIT" || verb == "NOOP" ||
			verb == "EHLO" || verb == "HELO" || verb == "AUTH" ||
			verb == "STARTTLS" {
			return nil
		}
		// FIXME more appropriate error
		c.logf("Use of verb not permitted in this state")
		return smtp.ReplyUnrecognisedCommand()
	}
	return nil
}

// tlsHandler handles the STARTTLS command
func (c *Session) tlsHandler(done func(ok bool)) (errorReply *smtp.Reply, callback func(), ok bool) {
	c.logf("Returning TLS handler")
	return nil, func() {
		c.logf("Upgrading session to TLS")
		// FIXME errors reading TLS config? should preload it
		tConn := tls.Server(c.conn.(net.Conn), c.server.getTLSConfig())
		err := tConn.Handshake()
		c.conn = tConn
		if err != nil {
			c.logf("handshake error in TLS connection: %s", err)
			done(false)
			return
		}
		c.isTLS = true
		c.logf("Session upgrade complete")
		done(true)
	}, true
}

func (c *Session) acceptMessage(msg *data.SMTPMessage) (id string, err error) {
	id, err = c.server.DeliveryBackend.Deliver(msg)
	c.logf("Storing message %s", id)
	return
}

func (c *Session) logf(message string, args ...interface{}) {
	message = strings.Join([]string{"[SMTP %s]", message}, " ")
	args = append([]interface{}{c.remoteAddress}, args...)
	log.Printf(message, args...)
}

// Read reads from the underlying io.Reader
func (c *Session) Read() bool {
	buf := make([]byte, 1024)
	n, err := io.Reader(c.conn).Read(buf)

	if n == 0 {
		c.logf("Connection closed by remote host\n")
		io.Closer(c.conn).Close() // not sure this is necessary?
		return false
	}

	if err != nil {
		c.logf("Error reading from socket: %s\n", err)
		return false
	}

	text := string(buf[0:n])
	logText := strings.Replace(text, "\n", "\\n", -1)
	logText = strings.Replace(logText, "\r", "\\r", -1)
	c.logf("Received %d bytes: '%s'\n", n, logText)

	if c.maximumBufferLength > -1 && len(c.line+text) > c.maximumBufferLength {
		// FIXME what is the "expected" behaviour for this?
		c.Write(smtp.ReplyError(fmt.Errorf("Maximum buffer length exceeded")))
		return false
	}

	c.line += text

	for strings.Contains(c.line, "\r\n") {
		line, reply := c.proto.Parse(c.line)
		c.line = line

		if reply != nil {
			c.Write(reply)
			if reply.Status == 221 {
				return false
			}
		}
	}

	return true
}

// Write writes a reply to the underlying io.Writer
func (c *Session) Write(reply *smtp.Reply) {
	lines := reply.Lines()
	for _, l := range lines {
		logText := strings.Replace(l, "\n", "\\n", -1)
		logText = strings.Replace(logText, "\r", "\\r", -1)
		c.logf("Sent %d bytes: '%s'", len(l), logText)
		io.Writer(c.conn).Write([]byte(l))
	}
	if reply.Done != nil {
		reply.Done()
	}
}
