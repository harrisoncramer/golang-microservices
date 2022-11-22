package main

import (
	"bytes"
	"html/template"

	"github.com/vanng822/go-premailer/premailer"
)

/* For the mail server we'll be communicating with */
type Mail struct {
	Domain     string
	Host       string
	Port       int
	Username   string
	Password   string
	Encryption string
	From       string
	FromName   string
}

/* The content of an email message */
type Message struct {
	From        string
	FromName    string
	To          string
	Subject     string
	Attachments []string /* Paths to attachments */
	Data        any
	DataMap     map[string]any /* Things to pass to the email templates */
}

func (m *Mail) SendSMTPMessage(msg Message) error {
	if msg.From == "" {
		msg.From = m.From
	}

	if msg.FromName == "" {
		msg.FromName = m.FromName
	}

	/* Build DataMap on message from it's own text */
	data := map[string]any{
		"message": msg.Data,
	}

	msg.DataMap = data

	formattedMessage, err := m.buildHTMLMessage(msg)

	if err != nil {
		return err
	}

	plainMessage, err := m.buildPlainTextMessage(msg)

	if err != nil {
		return err
	}

	// TODO: Create Mail Server to Send the Message (15:51)

	return nil

}

func (m *Mail) buildHTMLMessage(msg Message) (string, error) {

	templateToRender := "./templates/mail.html"
	t, err := template.New("email-html").ParseFiles(templateToRender)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer

	if err = t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		return "", err
	}

	formattedMessage := tpl.String()
	formattedMessage, err = m.inlineCSS(formattedMessage)

	if err != nil {
		return "", err
	}

	return formattedMessage, nil

}

func (m *Mail) buildPlainTextMessage(msg Message) (string, error) {
	templateToRender := "./templates/plain.html"
	t, err := template.New("email-plain").ParseFiles(templateToRender)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer

	if err = t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		return "", err
	}

	plainMessage := tpl.String()

	return plainMessage, nil

}

/* Adds styling to the email template */
func (m *Mail) inlineCSS(s string) (string, error) {
	options := premailer.Options{
		RemoveClasses:     false,
		CssToAttributes:   false,
		KeepBangImportant: true,
	}

	pre, err := premailer.NewPremailerFromString(s, &options)

	if err != nil {
		return "", err
	}

	html, err := pre.Transform()

	if err != nil {
		return "", err
	}

	return html, nil

}
