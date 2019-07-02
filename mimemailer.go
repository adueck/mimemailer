package mimemailer

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"net/mail"
	"net/smtp"
	"os"
	"text/template"
	"time"

	"jaytaylor.com/html2text"
)

// Config holds the SMTP connection and sending configuration infos
type Config struct {
	Host          string // SMTP Mail Server Name
	Port          string // Port to connect to SMTP Server
	Username      string // Username for SMTP Server Login
	Password      string // Password for SMTP SERVER Login
	SenderName    string // Name that will appear in the From: of email
	SenderAddress string // Email that will appear in the From: of email
}

// Email carries the content and information needed to send a single email
type Email struct {
	ToAddress       string
	ToName          string
	Subject         string
	HTML            string
	Date            time.Time
	ListUnsubscribe string // Optional value for List-Unsubscribe header. ie. "<mailto:unsubscribe@example.com?subject=unsubscribe-request>"
}

// For use in the text template in making the email
type emailInfoForTemplate struct {
	From            string
	To              string
	Subject         string
	ListUnsubscribe string
	HTMLQP          string
	TextQP          string
	Date            string
}

// Mailer implements a connection for sending SMTP emails
type Mailer struct {
	Config    Config
	connected bool
	client    *smtp.Client
}

var (
	// tmpl is used to create the email with headers etc for sending via smtp
	tmpl *template.Template
	// Connected is a state holder to see if the client is connected or non
	connected = false
)

// NewMailer returns a new instance for connecting and sending mail.
// Note that currently only secure TLS connections are supported.
func NewMailer(c Config) (*Mailer, error) {
	newMailer := &Mailer{
		Config: c,
	}
	newMailer.connected = false
	return newMailer, nil
}

func init() {
	// Create the email template
	var err error
	tmpl, err = template.New("Email Template").Parse(emailTemplate)
	if err != nil {
		fmt.Println("Error parsing email template")
		os.Exit(1)
	}
}

// IsConnected checks if the mailer is connected to an SMTP server
func (m *Mailer) IsConnected() bool {
	if m.connected {
		return true
	}
	return false
}

// Connect connects to SMTP Mail server
func (m *Mailer) Connect() error {
	if m.connected {
		return errors.New("mimemailer: smtp client is already connected")
	}
	// NOTE: This only works if the app has one SMTP connection at a time!
	auth := smtp.PlainAuth(
		"",
		m.Config.Username,
		m.Config.Password,
		m.Config.Host,
	)

	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         m.Config.Host,
	}

	// connect to client
	addr := net.JoinHostPort(m.Config.Host, m.Config.Port)
	tlsConn, err := tls.Dial("tcp", addr, tlsconfig)
	if err != nil {
		return err
	}
	m.client, err = smtp.NewClient(tlsConn, m.Config.Host)
	if err != nil {
		return err
	}

	// do auth
	err = m.client.Auth(auth)
	if err != nil {
		return err
	}
	// connected successfully
	m.connected = true
	return nil
}

// Disconnect disconnects from SMTP server
func (m *Mailer) Disconnect() error {
	err := m.client.Quit()
	if err != nil {
		return err
	}
	m.connected = false
	return nil
}

// SendEmail sends a single email. Note that the client must first be connected to the SMTP server.
func (m *Mailer) SendEmail(email Email) error {
	// Get from Adderss from config
	fromAddress := mail.Address{Name: m.Config.SenderName, Address: m.Config.SenderAddress}
	// Create the headers and body
	message, err := email.make(fromAddress)
	if err != nil {
		return err
	}

	// do sending
	// NOTE: In error handling, need to do mailClient.Reset() to abort message sending halfway through
	// Otherwise you get the Error: nested MAIL command from the SMTP
	temporarilyConnected := false
	if m.IsConnected() == false {
		err = m.Connect()
		if err != nil {
			return err
		}
		temporarilyConnected = true
	}
	err = m.client.Mail(m.Config.SenderAddress)
	if err != nil {
		_ = m.client.Reset()
		return err
	}

	err = m.client.Rcpt(email.ToAddress)
	if err != nil {
		_ = m.client.Reset()
		return err
	}

	// send body to smtp server
	var w io.WriteCloser
	w, err = m.client.Data()
	if err != nil {
		_ = m.client.Reset()
		return err
	}

	_, err = w.Write(message)
	if err != nil {
		_ = m.client.Reset()
		return err
	}

	err = w.Close()
	if err != nil {
		_ = m.client.Reset()
		return err
	}

	if temporarilyConnected == true {
		err = m.Disconnect()
		if err != nil {
			return err
		}
	}

	// mail sent successfully
	return nil
}

// make() will return an RFC 2045 compliant version of email and headers ready for sending
func (email Email) make(fromAddress mail.Address) ([]byte, error) {
	// Make a text version of the HTML email body
	text, err := html2text.FromString(email.HTML, html2text.Options{PrettyTables: false})
	if err != nil {
		return nil, err
	}

	// Prepare Sending Address for header fields
	toAddress := mail.Address{Name: email.ToName, Address: email.ToAddress}

	// Put the content sections (plaintext and HTML) into Quoted Printable format re: RFC 2045
	var textQP string
	textQP, err = convertToQuotedPrintable(text)
	if err != nil {
		return nil, err
	}

	var HTMLQP string
	HTMLQP, err = convertToQuotedPrintable(email.HTML)
	if err != nil {
		return nil, err
	}
	// Prepare the struct to pass to the templatea
	e := emailInfoForTemplate{
		From:            fromAddress.String(),
		To:              toAddress.String(),
		ListUnsubscribe: email.ListUnsubscribe,
		Subject:         email.Subject,
		Date:            email.Date.Format(time.RFC1123Z),
		HTMLQP:          HTMLQP,
		TextQP:          textQP,
	}

	// Compile the template with current email information
	var outputBuffer bytes.Buffer
	err = tmpl.Execute(&outputBuffer, e)
	if err != nil {
		return nil, err
	}
	// Need to convert all the mixed line endings to CRLF for email
	// At this point it is miexed. The templated headers use LF, and the printed-quotable sections use CRLF
	mixed := outputBuffer.Bytes()
	// Convert to have all CRLF (\r\n) line endings re: RFC 2045
	allCRLF := makeAllCRLF(mixed)
	return allCRLF, nil
}
