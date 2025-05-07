package main

import (
	"log"
	"net/http"

	"github.com/rs/cors"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/send-test-email", handleSendTestEmail)

	handler := cors.Default().Handler(mux)

	log.Println("✅ Backend running on http://localhost:8080")
	http.ListenAndServe(":8080", handler)
}
