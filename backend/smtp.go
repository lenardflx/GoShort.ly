package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-mail/mail"
)

type SMTPSettings struct {
	Host       string `json:"smtp_host"`
	Port       int    `json:"smtp_port"`
	From       string `json:"smtp_from"`
	FromName   string `json:"smtp_name"`
	Username   string `json:"smtp_user"`
	Password   string `json:"smtp_pass"`
	Security   string `json:"smtp_security"`    // "none", "starttls", "ssl"
	AuthMethod string `json:"smtp_auth_method"` // not used, but included
}

func handleSendTestEmail(w http.ResponseWriter, r *http.Request) {
	var settings SMTPSettings
	if err := json.NewDecoder(r.Body).Decode(&settings); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	m := mail.NewMessage()
	m.SetHeader("From", m.FormatAddress(settings.From, settings.FromName))
	m.SetHeader("To", settings.From)
	m.SetHeader("Subject", "SMTP Test from GoShort.ly")
	m.SetBody("text/plain", "✅ If you're reading this, SMTP works!")

	d := mail.NewDialer(settings.Host, settings.Port, settings.Username, settings.Password)

	switch settings.Security {
	case "ssl":
		d.SSL = true
	case "starttls":
		d.StartTLSPolicy = mail.MandatoryStartTLS
	case "none":
		d.SSL = false
		d.StartTLSPolicy = mail.NoStartTLS
	}

	if err := d.DialAndSend(m); err != nil {
		http.Error(w, "Send failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"ok"}`))
}
