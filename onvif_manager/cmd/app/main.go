package main

import (
	"fmt"
	"log"
	"os"

	"onvif_manager/internal/backend/api"
	"onvif_manager/internal/backend/camera"
	"onvif_manager/internal/backend/config"
	"onvif_manager/internal/cli"
)

func main() {
	// Check if CLI arguments are provided
	if len(os.Args) > 1 {
		// Check for special server command
		if os.Args[1] == "server" {
			// API server mode
			fmt.Println("🚀 Starting ONVIF Manager API Server")
			fmt.Println("📊 API endpoints will be available at http://localhost:8090/cameras")
			fmt.Println("")

			// Initialize cameras for the API server
			cams, err := config.LoadCameraList()
			if err != nil {
				log.Fatalf("Failed to load camera config: %v", err)
			}

			err = camera.InitializeAllCameras(cams)
			if err != nil {
				log.Printf("Warning: Failed to connect to some cameras: %v", err)
				log.Println("🔄 API server will continue, but some cameras may not be accessible")
			}

			// Start the API server
			api.StartServer(":8090")
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
	fmt.Println("🎥 ONVIF Manager - CLI and API Server")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  onvif-manager [command]")
	fmt.Println("")
	fmt.Println("Available Commands:")
	fmt.Println("  server            Start API server on port 8090")
	fmt.Println("  list              List all cameras")
	fmt.Println("  config            Configuration management commands")
	fmt.Println("  help              Help about any command")
	fmt.Println("")
	fmt.Println("Use 'onvif-manager help [command]' for more information about a command.")
}
