package routers

import (
	"encoding/json"
	"net/http"
	"time"
	"urlshort-backend/modules"
	"urlshort-backend/smtp"
)

type InviteRequest struct {
	Email string `json:"email"`
}

// HandleCreateInvite handles the creation of an invitation link
func HandleCreateInvite(w http.ResponseWriter, r *http.Request) {
	var req InviteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Email == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	invite, err := modules.CreateInvite(req.Email)
	if err != nil {
		http.Error(w, "Failed to create invite", http.StatusInternalServerError)
		return
	}

	cfg, err := modules.Load()
	if err != nil {
		http.Error(w, "Failed to load config", http.StatusInternalServerError)
		return
	}

	link := cfg.BaseURL + "/invite/" + invite.Token

	// Try to build SMTP config and send email
	if smtpSettings, err := smtp.BuildFromConfig(*cfg); err == nil {
		htmlBody, err := smtp.RenderTemplate("invite.html", map[string]any{
			"ServiceName": cfg.BrandingName,
			"SignUpLink":  link,
		})
		if err != nil {
			return
		}
		smtp.Send(smtpSettings, smtp.Message{
			To:       req.Email,
			Subject:  "You're invited to GoShort.ly",
			HTML:     true,
			HTMLBody: htmlBody,
		})
	}

	// Respond either way
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "ok",
		"invite":  link,
		"email":   invite.Email,
		"token":   invite.Token,
		"expires": invite.ExpiresAt.Format(time.RFC3339),
	})
}
