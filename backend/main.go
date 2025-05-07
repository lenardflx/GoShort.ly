package main

import (
	"log"
	"net/http"
	"urlshort-backend/db"
	"urlshort-backend/modules"

	"urlshort-backend/routes"
)

func main() {
	db.InitDB()
	modules.AutoMigrate()
	modules.InitDefaults()

	log.Println("🚀 Server running at http://localhost:8080")
	err := http.ListenAndServe(":8080", routes.SetupRoutes())
	if err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	} else {
		log.Println("✅ Server started successfully")
	}
}
