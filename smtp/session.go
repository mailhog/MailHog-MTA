package smtp

// http://www.rfc-editor.org/rfc/rfc5321.txt

import (
	"crypto/tls"
	"io"
	"log"
	"net"
	"strings"

	"github.com/mailhog/MailHog-MTA/backend"
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
	identity      *backend.Identity
}

// Accept starts a new SMTP session using io.ReadWriteCloser
func (s *Server) Accept(remoteAddress string, conn io.ReadWriteCloser) {
	proto := smtp.NewProtocol()
	proto.Hostname = s.Hostname

	session := &Session{
		server:        s,
		conn:          conn,
		proto:         proto,
		remoteAddress: remoteAddress,
		isTLS:         false,
		line:          "",
		identity:      nil,
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

	if session.server.PolicySet.EnableTLS {
		proto.TLSHandler = session.tlsHandler
		proto.RequireTLS = session.server.PolicySet.RequireTLS
	}

	session.logf("Starting session")
	session.Write(proto.Start())
	for session.Read() == true {
	}
	session.logf("Session ended")
}

func (c *Session) validateAuthentication(mechanism string, args ...string) (errorReply *smtp.Reply, ok bool) {
	if c.server.AuthBackend == nil {
		return smtp.ReplyInvalidAuth(), false
	}
	i, e, ok := c.server.AuthBackend.Authenticate(mechanism, args...)
	if e != nil || !ok {
		return smtp.ReplyInvalidAuth(), false
	}
	c.identity = i
	return nil, true
}

func (c *Session) validateRecipient(to string) bool {
	if c.server.DeliveryBackend == nil {
		return false
	}
	maxRecipients := c.server.DeliveryBackend.MaxRecipients(c.identity)
	if maxRecipients > -1 && len(c.proto.Message.To) > maxRecipients {
		return false
	}
	if c.server.PolicySet.RequireLocalDelivery {
		r := c.server.ResolverBackend.Resolve(to)
		if r != backend.ResolvedPrimaryLocal && r != backend.ResolvedSecondaryLocal {
			return false
		}
	}
	return c.server.DeliveryBackend.WillDeliver(to, c.proto.Message.From, c.identity)
}

func (c *Session) validateSender(from string) bool {
	// FIXME better policy for this?
	if c.server.PolicySet.RequireAuthentication {
		if c.identity == nil {
			return false
		}
		return (*c.identity).IsValidSender(from)
	}
	return true
}

func (c *Session) verbFilter(verb string, args ...string) (errorReply *smtp.Reply) {
	// FIXME consider moving this to smtp proto? since STARTTLS is there anyway...
	if c.server.PolicySet.RequireAuthentication && c.proto.State == smtp.MAIL && c.identity == nil {
		verb = strings.ToUpper(verb)
		if verb == "RSET" || verb == "QUIT" || verb == "NOOP" ||
			verb == "EHLO" || verb == "HELO" || verb == "AUTH" {
			return nil
		}
		// FIXME more appropriate error
		return smtp.ReplyUnrecognisedCommand()
	}
	return nil
}

// tlsHandler handles the STARTTLS command
func (c *Session) tlsHandler(done func(ok bool)) (errorReply *smtp.Reply, callback func(), ok bool) {
	return nil, func() {
		c.logf("Upgrading session to TLS")
		c.conn = tls.Server(c.conn.(net.Conn), c.server.getTLSConfig())
		c.isTLS = true
		c.logf("Session upgrade complete")
		done(true)
	}, true
}

func (c *Session) acceptMessage(msg *data.Message) (id string, err error) {
	c.logf("Storing message %s", msg.ID)
	return c.server.DeliveryBackend.Deliver(msg)
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

	c.line += text

	for strings.Contains(c.line, "\r\n") {
		line, reply := c.proto.Parse(c.line)
		c.line = line

		if reply != nil {
			c.Write(reply)
			if reply.Status == 221 {
				io.Closer(c.conn).Close()
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
