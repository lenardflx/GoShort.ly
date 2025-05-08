package routers

import (
	"encoding/json"
	"net/http"

	"urlshort-backend/smtp"
)

type smtpPayload struct {
	smtp.Settings
	To string `json:"smtp_test_to"`
}

func HandleSendTestEmail(w http.ResponseWriter, r *http.Request) {
	var payload smtpPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	html, err := smtp.RenderTemplate("test-smtp.html", map[string]string{
		"DashboardURL": "https://your-url-shortener.local/dashboard",
	})
	if err != nil {
		http.Error(w, "Failed to load template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = smtp.Send(payload.Settings, smtp.Message{
		To:       payload.To,
		Subject:  "📬 SMTP Test from GoShort.ly",
		HTMLBody: html,
		HTML:     true,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"ok"}`))
}
