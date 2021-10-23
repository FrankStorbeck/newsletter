// sendmail implements a mail sender

// Package sendmail can be used to send an e-mail to a number of addressees.
// It uses TLS for delivering e-mails.
package sendmail

import (
	"bytes"
	"fmt"
	"html/template"
	"path/filepath"
	"strconv"
	"time"

	"github.com/go-mail/mail/v2"
)

// Config contains data for building a message
type Config struct {
	Sender      string         // Tells the SMTP server the sender of the e-mail
	From        NamedAddress   // Reply address
	To          []NamedAddress // Recipients
	Cc          []NamedAddress // Carbon copy recipients
	Bcc         []NamedAddress // Blind carbon copy recipients
	Subject     string         // Subject
	PlainText   string         // plain version of the mail body
	HTMLText    string         // HTML version of the mail body
	Embedd      []string       // Embedded files, only used when HTMLText defined
	Attachments []string       // Attachments
}

// NamedAddress holds the e-mail address and the name of a recipient or sender.
type NamedAddress struct {
	EMail string // email address
	Name  string // real name
}

// Body creates an e-mail body
func Body(bodyType string, tmplt *template.Template, data interface{}) (string, error) {
	body := new(bytes.Buffer)
	var err error
	if len(bodyType) > 0 {
		err = tmplt.ExecuteTemplate(body, bodyType, data)
	} else {
		err = tmplt.Execute(body, data)
	}
	if err != nil &&
		err.Error() == fmt.Sprintf("html/template: %q is undefined", bodyType) {
		body = new(bytes.Buffer)
		err = nil
	}

	return body.String(), err
}

// BuildMessage builds a complete mail.Message
func (c *Config) BuildMessage() *mail.Message {
	m := mail.NewMessage()

	m.SetHeader("Sender", c.Sender)

	SetAddresses(m, "From", []NamedAddress{c.From}...)
	if len(c.To) > 0 {
		SetAddresses(m, "To", c.To...)
	}

	if len(c.Cc) > 0 {
		SetAddresses(m, "Cc", c.Cc...)
	}
	if len(c.Bcc) > 0 {
		SetAddresses(m, "Bcc", c.Bcc...)
	}

	m.SetHeader("Subject", c.Subject)

	if len(c.PlainText) > 0 {
		m.SetBody("text/plain", c.PlainText)
	}

	if len(c.HTMLText) > 0 {
		setBody := m.AddAlternative
		if len(c.PlainText) == 0 {
			setBody = m.SetBody
		}
		setBody("text/html", c.HTMLText)
		for _, e := range c.Embedd {
			_, baseName := filepath.Split(e)
			m.Embed(e, mail.Rename(baseName))
		}
	}

	for _, att := range c.Attachments {
		_, baseName := filepath.Split(att)
		m.Attach(att, mail.Rename(baseName))
	}

	m.SetDateHeader("X-Date", time.Now())

	return m
}

// NewDialer returns a dialer to access the SMTP server using the auth data.
func NewDialer(a *Auth) (*mail.Dialer, error) {
	prt, err := strconv.Atoi(a.Value("port"))
	if err != nil {
		return nil, err
	}

	dialer := mail.NewDialer(a.Value("hostname"), prt, a.Value("username"),
		a.Value("password"))
	dialer.Timeout = 5 * time.Second
	dialer.StartTLSPolicy = mail.MandatoryStartTLS

	return dialer, nil
}

// SetAddresses sets recipients or sender. hdr can be  "From", "Sender", "To",
// "Cc" or "Bcc". When addresses is empty, nothing happens. When hdr is "From"
// or "Sender" only the first NamedAddress will be used.
func SetAddresses(m *mail.Message, hdr string, addresses ...NamedAddress) {
	if len(addresses) > 0 {
		switch hdr {
		case "From", "Sender":
			m.SetHeader(hdr,
				m.FormatAddress(addresses[0].EMail, addresses[0].Name))
		case "To", "Cc", "Bcc":
			addrs := make([]string, len(addresses))
			for i, a := range addresses {
				addrs[i] = m.FormatAddress(a.EMail, a.Name)
			}
			m.SetHeader(hdr, addrs...) // takes care of duplicates
		default:
		}
	}
}
