package server

import (
	"log"
	"net/http"

	"onvif_manager/internal/backend/api"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// StartAPIServer starts the API server only (no frontend)
func StartAPIServer(addr string) {
	r := mux.NewRouter()

	// Register API routes
	api.RegisterRoutes(r)

	// Configure CORS
	corsOptions := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}), // Allow all origins for development
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	log.Printf("📊 API server starting on http://localhost%s", addr)
	log.Printf("� API endpoints available at http://localhost%s/cameras", addr)

	// Wrap the router with the CORS handler
	err := http.ListenAndServe(addr, corsOptions(r))
	if err != nil {
		log.Fatal("API server failed to start:", err)
	}
}
