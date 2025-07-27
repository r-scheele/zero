package services

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"net/smtp"
	"strings"

	"github.com/r-scheele/zero/config"
	"github.com/r-scheele/zero/pkg/log"
	"maragu.dev/gomponents"

	"github.com/labstack/echo/v4"
)

type (
	// MailClient provides a client for sending email
	// This is purposely not completed because there are many different methods and services
	// for sending email, many of which are very different. Choose what works best for you
	// and populate the methods below. For now, emails will just be logged.
	MailClient struct {
		// config stores application configuration.
		config *config.Config
	}

	// mail represents an email to be sent.
	mail struct {
		client    *MailClient
		from      string
		to        string
		subject   string
		body      string
		component gomponents.Node
	}
)

// NewMailClient creates a new MailClient.
func NewMailClient(cfg *config.Config) (*MailClient, error) {
	return &MailClient{
		config: cfg,
	}, nil
}

// Compose creates a new email.
func (m *MailClient) Compose() *mail {
	return &mail{
		client: m,
		from:   m.config.Mail.FromAddress,
	}
}

// skipSend determines if mail sending should be skipped.
func (m *MailClient) skipSend() bool {
	return m.config.App.Environment != config.EnvProduction
}

// send attempts to send the email.
func (m *MailClient) send(email *mail, ctx echo.Context) error {
	switch {
	case email.to == "":
		return errors.New("email cannot be sent without a to address")
	case email.body == "" && email.component == nil:
		return errors.New("email cannot be sent without a body or component to render")
	}

	// Check if a component was supplied.
	if email.component != nil {
		// Render the component and use as the body.
		// TODO pool the buffers?
		buf := bytes.NewBuffer(nil)
		if err := email.component.Render(buf); err != nil {
			return err
		}

		email.body = buf.String()
	}

	// Check if mail sending should be skipped.
	if m.skipSend() {
		log.Ctx(ctx).Debug("skipping email delivery in non-production environment",
			"to", email.to,
			"subject", email.subject,
		)
		return nil
	}

	// Send the actual email via SMTP
	return m.sendSMTP(email, ctx)
}

// sendSMTP sends the email using SMTP
func (m *MailClient) sendSMTP(email *mail, ctx echo.Context) error {
	// SMTP server configuration
	host := m.config.Mail.Hostname
	port := fmt.Sprintf("%d", m.config.Mail.Port)
	addr := fmt.Sprintf("%s:%s", host, port)

	// Authentication
	var auth smtp.Auth
	if m.config.Mail.User != "" && m.config.Mail.Password != "" {
		auth = smtp.PlainAuth("", m.config.Mail.User, m.config.Mail.Password, host)
	}

	// Create the email message
	subject := email.subject
	body := email.body

	// Determine if body is HTML or plain text
	contentType := "text/plain; charset=UTF-8"
	if strings.Contains(body, "<html>") || strings.Contains(body, "<HTML>") {
		contentType = "text/html; charset=UTF-8"
	}

	// Build the email message
	msg := fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: %s\r\n"+
		"\r\n"+
		"%s\r\n",
		email.from,
		email.to,
		subject,
		contentType,
		body)

	// Recipients
	to := []string{email.to}

	// Attempt to send the email
	var err error
	if port == "587" || port == "25" {
		// Use STARTTLS for ports 587 and 25
		err = m.sendWithStartTLS(addr, auth, email.from, to, []byte(msg))
	} else {
		// Use direct SMTP for other ports
		err = smtp.SendMail(addr, auth, email.from, to, []byte(msg))
	}

	if err != nil {
		log.Ctx(ctx).Error("failed to send email",
			"to", email.to,
			"subject", email.subject,
			"error", err,
		)
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Ctx(ctx).Info("email sent successfully",
		"to", email.to,
		"subject", email.subject,
	)
	return nil
}

// sendWithStartTLS sends email using STARTTLS (for ports like 587)
func (m *MailClient) sendWithStartTLS(addr string, auth smtp.Auth, from string, to []string, msg []byte) error {
	// Create connection
	c, err := smtp.Dial(addr)
	if err != nil {
		return err
	}
	defer c.Close()

	// Use STARTTLS if supported
	if ok, _ := c.Extension("STARTTLS"); ok {
		config := &tls.Config{ServerName: m.config.Mail.Hostname}
		if err = c.StartTLS(config); err != nil {
			return err
		}
	}

	// Authenticate if auth is provided
	if auth != nil {
		if err = c.Auth(auth); err != nil {
			return err
		}
	}

	// Set sender
	if err = c.Mail(from); err != nil {
		return err
	}

	// Set recipients
	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}

	// Send the email body
	w, err := c.Data()
	if err != nil {
		return err
	}

	_, err = w.Write(msg)
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	return c.Quit()
}

// From sets the email from address.
func (m *mail) From(from string) *mail {
	m.from = from
	return m
}

// To sets the email address this email will be sent to.
func (m *mail) To(to string) *mail {
	m.to = to
	return m
}

// Subject sets the subject line of the email.
func (m *mail) Subject(subject string) *mail {
	m.subject = subject
	return m
}

// Body sets the body of the email.
// This is not required and will be ignored if a component is set via Component().
func (m *mail) Body(body string) *mail {
	m.body = body
	return m
}

// Component sets a renderable component to use as the body of the email.
func (m *mail) Component(component gomponents.Node) *mail {
	m.component = component
	return m
}

// Send attempts to send the email.
func (m *mail) Send(ctx echo.Context) error {
	return m.client.send(m, ctx)
}
