package webserver

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"

	"onvif_manager/internal/backend/api"
	"onvif_manager/internal/cli"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// We need web assets to be embedded where this file can see them
//
//go:embed web/*
var webFiles embed.FS

// StartApp handles CLI arguments and starts the appropriate mode
func StartApp() {
	if len(os.Args) > 1 {
		// Check for web server command
		if os.Args[1] == "web" {
			// Combined web server mode (API + Frontend)
			fmt.Println("🌐 Starting ONVIF Manager Web Application")
			fmt.Println("📱 Frontend will be available at http://localhost:8090")
			fmt.Println("🔌 API endpoints will be available at http://localhost:8090/api")
			fmt.Println("")

			// Camera initialization is handled by in-memory storage
			StartWebServer(":8090")
			return
		}

		// Check for API server only command
		if os.Args[1] == "server" {
			// API server mode only
			fmt.Println("🚀 Starting ONVIF Manager API Server")
			fmt.Println("📊 API endpoints will be available at http://localhost:8090")
			fmt.Println("")

			// Camera initialization is handled by in-memory storage
			StartAPIServer(":8090")
			return
		}

		// CLI mode - run the CLI interface
		if err := cli.Execute(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return
	}
	// No arguments provided - show usage information
	fmt.Println("🎥 ONVIF Manager - CLI, API Server, and Web Application")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  onvif-manager [command]")
	fmt.Println("")
	fmt.Println("Available Commands:")
	fmt.Println("  web               Start combined web application (frontend + API) on port 8090")
	fmt.Println("  server            Start API server only on port 8090")
	fmt.Println("  config            Configuration management commands")
	fmt.Println("  help              Help about any command")
	fmt.Println("")
	fmt.Println("Use 'onvif-manager help [command]' for more information about a command.")
}

// StartWebServer starts the combined web server with both API and frontend
func StartWebServer(addr string) {
	r := mux.NewRouter()

	// Register API routes under /api prefix
	apiRouter := r.PathPrefix("/api").Subrouter()
	api.RegisterRoutes(apiRouter)

	// Serve embedded static files
	webFS, err := fs.Sub(webFiles, "web")
	if err != nil {
		log.Fatal("Failed to get web files:", err)
	}

	// Create a custom file server that handles SPA routing
	fileServer := http.FileServer(http.FS(webFS))

	// Handle static assets
	r.PathPrefix("/assets/").Handler(fileServer)
	r.Handle("/vite.svg", fileServer)
	r.Handle("/debug.html", fileServer)

	// Handle SPA routing - serve index.html for all other routes
	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		// If it's a file request with extension, serve it directly
		if strings.Contains(path, ".") && !strings.HasSuffix(path, "/") {
			fileServer.ServeHTTP(w, r)
			return
		}

		// For all other routes (SPA routes), serve index.html
		r.URL.Path = "/"
		fileServer.ServeHTTP(w, r)
	})

	// Configure CORS
	corsOptions := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	log.Printf("🌐 Starting web server on %s", addr)
	log.Printf("📱 Frontend available at: http://localhost%s", addr)
	log.Printf("🔌 API available at: http://localhost%s/api", addr)

	err = http.ListenAndServe(addr, corsOptions(r))
	if err != nil {
		log.Fatal("Web server failed to start:", err)
	}
}

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
	log.Printf("🔌 API endpoints available at http://localhost%s/cameras", addr)

	// Wrap the router with the CORS handler
	err := http.ListenAndServe(addr, corsOptions(r))
	if err != nil {
		log.Fatal("API server failed to start:", err)
	}
}
