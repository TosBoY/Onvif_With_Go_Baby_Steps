package main

import (
	"fmt"
	"log"
	"main_back/internal/camera"
	"main_back/pkg/models"
)

func main() {

	cam2 := models.Camera{
		ID:       "2",
		IP:       "192.168.1.30",
		Username: "admin",
		Password: "Admin123",
	}

	cam2_Client, err := camera.NewCameraClient(cam2)
	if err != nil {
		log.Fatalf("Error initializing camera client: %v", err)
	}

	// Mock profile and config tokens for Camera 2
	profileTokens, configTokens, err := camera.GetProfilesAndConfigs(cam2_Client)
	if err != nil {
		log.Fatalf("Error fetching profiles and configs: %v", err)
	}

	// Get current encoder options
	fmt.Println("Fetching current encoder options for Camera 2...")
	encoderOptions, err := camera.GetCurrentEncoderOptions(cam2_Client, profileTokens[0], configTokens[0])
	if err != nil {
		log.Fatalf("Error fetching encoder options: %v", err)
	}
	fmt.Printf("Available resolutions: %+v\n", encoderOptions.Resolutions)
	fmt.Printf("Available FPS options: %+v\n", encoderOptions.FPSOptions)

	// Define a target resolution
	targetResolution := models.Resolution{
		Width:  1440,
		Height: 3000,
	}

	// Find the closest resolution
	fmt.Println("Finding closest resolution to target...")
	closestResolution := camera.FindClosestResolution(targetResolution, encoderOptions.Resolutions)
	fmt.Printf("Closest resolution to %dx%d: %dx%d\n", targetResolution.Width, targetResolution.Height, closestResolution.Width, closestResolution.Height)
}
