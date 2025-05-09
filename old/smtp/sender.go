package smtp

import (
	"bytes"
	"errors"
	"fmt"
	"goshortly/models"
	"html/template"
	"path/filepath"

	"github.com/go-mail/mail"
)

type Settings struct {
	Host     string `json:"smtp_host"`
	Port     int    `json:"smtp_port"`
	From     string `json:"smtp_from"`
	FromName string `json:"smtp_name"`
	Username string `json:"smtp_user"`
	Password string `json:"smtp_pass"`
	Security string `json:"smtp_security"`
}

type Message struct {
	To       string
	Subject  string
	Body     string
	HTMLBody string
	HTML     bool
}

func Send(settings Settings, msg Message) error {
	m := mail.NewMessage()
	m.SetHeader("From", m.FormatAddress(settings.From, settings.FromName))
	m.SetHeader("To", msg.To)
	m.SetHeader("Subject", msg.Subject)

	if msg.HTML {
		m.SetBody("text/html", msg.HTMLBody)
	} else {
		m.SetBody("text/plain", msg.Body)
	}

	d := mail.NewDialer(settings.Host, settings.Port, settings.Username, settings.Password)

	switch settings.Security {
	case "ssl":
		d.SSL = true
	case "starttls":
		d.StartTLSPolicy = mail.MandatoryStartTLS
	default:
		d.SSL = false
		d.StartTLSPolicy = mail.NoStartTLS
	}

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("smtp send failed: %w", err)
	}
	return nil
}

func RenderTemplate(templateName string, data any) (string, error) {
	path := filepath.Join("smtp", "templates", templateName)
	tpl, err := template.ParseFiles(path)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func BuildFromConfig(cfg models.AdminConfig) (Settings, error) {
	if cfg.SMTPHost == "" ||
		cfg.SMTPPort <= 0 ||
		cfg.SMTPFrom == "" ||
		cfg.SMTPUser == "" ||
		cfg.SMTPPass == "" {
		return Settings{}, errors.New("SMTP is not fully configured")
	}

	return Settings{
		Host:     cfg.SMTPHost,
		Port:     cfg.SMTPPort,
		From:     cfg.SMTPFrom,
		FromName: cfg.SMTPName,
		Username: cfg.SMTPUser,
		Password: cfg.SMTPPass,
		Security: cfg.SMTPSecurity,
	}, nil
}
