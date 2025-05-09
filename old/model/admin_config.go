package model

import "gorm.io/gorm"

type AdminConfig struct {
	gorm.Model
	// Admin Settings
	AdminPassHash string

	// General Settings
	BaseURL      string `json:"base_url"`
	MaxURLLength int    `json:"max_url_length"`
	BrandingName string `json:"branding_name"`

	// User Settings
	DefaultUserLimit       int  `json:"default_user_limit"`
	DefaultExpirationDays  int  `json:"default_expiration"`
	UserOverrideExpiration bool `json:"user_override_expiration"`
	AllowAnonymous         bool `json:"allow_anonymous"`
	FailedLoginLimit       int  `json:"failed_login_limit"`

	// SMTP Settings
	SMTPHost       string `json:"smtp_host"`
	SMTPPort       int    `json:"smtp_port,string"`
	SMTPSecurity   string `json:"smtp_security"`
	SMTPFrom       string `json:"smtp_from"`
	SMTPName       string `json:"smtp_name"`
	SMTPUser       string `json:"smtp_user"`
	SMTPPass       string `json:"smtp_pass"`
	SMTPAuthMethod string `json:"smtp_auth_method"`
}
