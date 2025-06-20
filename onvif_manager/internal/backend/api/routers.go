package api

import (
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func StartServer(addr string) {
	r := mux.NewRouter()
	RegisterRoutes(r)

	// Configure CORS
	corsOptions := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}), // Allow all origins for development
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	log.Println("Listening on", addr)
	// Wrap the router with the CORS handler
	err := http.ListenAndServe(addr, corsOptions(r))
	if err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
