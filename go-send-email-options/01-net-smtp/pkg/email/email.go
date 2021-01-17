package email

import (
	"fmt"
	"strings"
)

// Message represents a simple email message.
type Message struct {
	From     string
	Tos      []string
	Subject  string
	Body     string
	BodyHtml bool
}

// NewMessage creates a message.
func NewMessage(from, to, subject, body string, bodyHtml bool) *Message {
	e := Message{
		From:     from,
		Tos:      []string{to},
		Subject:  subject,
		Body:     body,
		BodyHtml: bodyHtml,
	}
	return &e
}

// NewMessageToMany creates a message that has multiple recipients
func NewMessageToMany(from string, tos []string, subject, body string, bodyHtml bool) *Message {
	e := Message{
		From:     from,
		Tos:      tos,
		Subject:  subject,
		Body:     body,
		BodyHtml: bodyHtml,
	}
	return &e
}

// Build creates the email data that is sent as-is to the server.
func (m *Message) Build() string {
	var sb strings.Builder
	contentType := "plain"
	if m.BodyHtml {
		contentType = "html"
	}
	fmt.Fprintf(&sb, "Content-Type: text/%s; charset=utf-8\r\n", contentType)
	fmt.Fprintf(&sb, "From: %s\r\n", m.From)
	if len(m.Tos) > 0 {
		fmt.Fprintf(&sb, "To: %s\r\n", strings.Join(m.Tos, ";"))
	}
	fmt.Fprintf(&sb, "Subject: %s\r\n", m.Subject)
	fmt.Fprintf(&sb, "\r\n%s", m.Body)

	return sb.String()
}
