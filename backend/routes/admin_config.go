package routes

import (
	"encoding/json"
	"net/http"
	"urlshort-backend/models"
	"urlshort-backend/modules"
)

// HandleGetAdminConfig handles GET requests for the admin config
func HandleGetAdminConfig(w http.ResponseWriter, r *http.Request) {
	cfg, err := modules.Load()
	if err != nil {
		http.Error(w, "Could not load admin config", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cfg)
}

// HandleSaveAdminConfig handles POST requests to save the admin config
func HandleSaveAdminConfig(w http.ResponseWriter, r *http.Request) {
	var input models.AdminConfig
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Save to DB
	if err := modules.Save(&input); err != nil {
		http.Error(w, "Failed to save config", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"ok"}`))
}

// HandleGetAdminHash handles GET requests to get the admin password hash
func HandleGetAdminHash(w http.ResponseWriter, r *http.Request) {
	cfg, err := modules.Load()
	if err != nil {
		http.Error(w, "Config not found", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"hash": cfg.AdminPassHash,
	})
}
