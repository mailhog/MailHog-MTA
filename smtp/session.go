package smtp

// http://www.rfc-editor.org/rfc/rfc5321.txt

import (
	"io"
	"log"
	"strings"

	"github.com/ian-kent/Go-MailHog/MailHog-MTA/backend"
	"github.com/ian-kent/Go-MailHog/data"
	"github.com/ian-kent/Go-MailHog/smtp/protocol"
)

// Session represents a SMTP session using net.TCPConn
type Session struct {
	server *Server

	conn          io.ReadWriteCloser
	proto         *protocol.Protocol
	remoteAddress string
	isTLS         bool
	line          string
	identity      *backend.Identity
}

// Accept starts a new SMTP session using io.ReadWriteCloser
func (s *Server) Accept(remoteAddress string, conn io.ReadWriteCloser) {
	proto := protocol.NewProtocol()
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

	session.logf("Starting session")
	session.Write(proto.Start())
	for session.Read() == true {
	}
	session.logf("Session ended")
}

func (c *Session) validateAuthentication(mechanism string, args ...string) (errorReply *protocol.Reply, ok bool) {
	i, e, ok := c.server.AuthBackend.Authenticate(mechanism, args...)
	if e != nil || !ok {
		return protocol.ReplyInvalidAuth(), false
	}
	c.identity = i
	return nil, true
}

func (c *Session) validateRecipient(to string) bool {
	return c.server.DeliveryBackend.WillDeliver(to, c.proto.Message.From, c.identity)
}

func (c *Session) validateSender(from string) bool {
	return true
}

func (c *Session) verbFilter(verb string, args ...string) (errorReply *protocol.Reply) {
	if c.server.PolicySet.RequireAuthentication && c.proto.State == protocol.MAIL && c.identity == nil {
		verb = strings.ToUpper(verb)
		if verb == "RSET" || verb == "QUIT" || verb == "NOOP" ||
			verb == "EHLO" || verb == "HELO" || verb == "AUTH" {
			return nil
		}
		// FIXME more appropriate error
		return protocol.ReplyUnrecognisedCommand()
	}
	return nil
}

func (c *Session) acceptMessage(msg *data.Message) (id string, err error) {
	c.logf("Storing message %s", msg.ID)
	//id, err = c.storage.Store(msg)
	//c.messageChan <- msg
	return
}

func (c *Session) logf(message string, args ...interface{}) {
	message = strings.Join([]string{"[SMTP %s]", message}, " ")
	args = append([]interface{}{c.remoteAddress}, args...)
	log.Printf(message, args...)
}

// Read reads from the underlying net.TCPConn
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

	for strings.Contains(c.line, "\n") {
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

// Write writes a reply to the underlying net.TCPConn
func (c *Session) Write(reply *protocol.Reply) {
	lines := reply.Lines()
	for _, l := range lines {
		logText := strings.Replace(l, "\n", "\\n", -1)
		logText = strings.Replace(logText, "\r", "\\r", -1)
		c.logf("Sent %d bytes: '%s'", len(l), logText)
		io.Writer(c.conn).Write([]byte(l))
	}
}
