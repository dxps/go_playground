package smtp

import (
	"crypto/tls"
	"log"
	"net/smtp"

	"github.com/devisions/go-playground/go-send-email-options/01-net-smtp/pkg/email"
	"github.com/pkg/errors"
)

type SMTPClient struct {
	host string
	port string
	auth smtp.Auth
	conn *tls.Conn
}

func NewSMTPClient(host, port string) *SMTPClient {
	s := SMTPClient{
		host: host,
		port: port,
	}
	return &s
}

// NewSMTPClientGMailTLS ...
func NewSMTPClientGMailTLS(username, password string) *SMTPClient {
	s := SMTPClient{
		host: "smtp.gmail.com",
		port: "465", // This is the standard port. 587 is a cleartext SMTP port with STARTTLS support.
	}
	s.auth = smtp.PlainAuth("", username, password, s.host)
	return &s
}

func (c *SMTPClient) connect() error {
	tlscfg := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         c.host,
	}
	conn, err := tls.Dial("tcp", c.serverEndpoint(), tlscfg)
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

func (c *SMTPClient) Send(msg email.Message) error {

	if c.conn == nil {
		if err := c.connect(); err != nil {
			return errors.Wrap(err, "initing the connection")
		}
	}

	client, err := smtp.NewClient(c.conn, c.host)
	if err != nil {
		return errors.Wrap(err, "initing the client")
	}
	defer func() { _ = client.Quit() }()

	if err = client.Auth(c.auth); err != nil {
		return errors.Wrap(err, "setting the client auth")
	}
	// Setting the sender ('from').
	if err = client.Mail(msg.From); err != nil {
		return errors.Wrap(err, "setting the 'from'")
	}
	// Setting the recipients ('rcpt's).
	for _, r := range msg.Tos {
		if err = client.Rcpt(r); err != nil {
			return errors.Wrap(err, "setting recipients")
		}
	}
	// Setting the data ('data').
	w, err := client.Data()
	if err != nil {
		return errors.Wrap(err, "setting the data")
	}
	_, err = w.Write([]byte(msg.Build()))
	if err != nil {
		return errors.Wrap(err, "sending the message data")
	}
	if err := w.Close(); err != nil {
		log.Println("[warn] Failed closing the sending writer. Reason:", err)
	}

	return nil
}

func (c *SMTPClient) SendAndDone(msg email.Message) error {
	if err := c.Send(msg); err != nil {
		return err
	}
	if err := c.conn.Close(); err != nil {
		log.Println("[wrn] Failed to close the connection. Reason:", err)
	}
	return nil
}

func (s *SMTPClient) serverEndpoint() string {
	return s.host + ":" + s.port
}
