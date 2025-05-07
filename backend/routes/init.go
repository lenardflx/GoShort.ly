package routes

import (
	"net/http"

	"github.com/rs/cors"
)

// SetupRoutes binds all handlers and returns the full HTTP handler with middleware
func SetupRoutes() http.Handler {
	mux := http.NewServeMux()

	// Admin config API
	mux.HandleFunc("/api/admin-config", HandleGetAdminConfig)       // GET
	mux.HandleFunc("/api/admin-config/save", HandleSaveAdminConfig) // POST
	mux.HandleFunc("/api/admin-config/hash", HandleGetAdminHash)    // GET

	// SMTP test API
	mux.HandleFunc("/api/send-test-email", HandleSendTestEmail) // POST

	// Invite API
	mux.HandleFunc("/api/invite", HandleCreateInvite) // POST

	handler := cors.Default().Handler(mux)
	return handler
}
