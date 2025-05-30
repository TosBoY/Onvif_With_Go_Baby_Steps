package main

import (
	"fmt"
	"log"
	"time"

	"github.com/videonext/onvif/profiles/media2"
	"github.com/videonext/onvif/soap"
)

func applyResolution2(ip, username, password string, width, height, fps, bitrate, gop int) error {
	client := soap.NewClient(soap.WithTimeout(5 * time.Second))
	client.AddHeader(soap.NewWSSSecurityHeader(username, password, time.Now()))

	mediaService := media2.NewMedia2(client, fmt.Sprintf("http://%s/onvif/media_service", ip)) // Adjusted for media package

	// Get profiles
	profiles, err := mediaService.GetProfiles(&media2.GetProfiles{})
	if err != nil || len(profiles.Profiles) == 0 {
		fmt.Println("Profiles:", profiles) // Print profiles
		return fmt.Errorf("failed to get profiles: %v", err)
	}
	fmt.Println("Profiles:", profiles) // Print profiles

	// Get video source configurations (equivalent to getting configs in the example)
	configs, err := mediaService.GetVideoEncoderConfigurations(&media2.GetConfiguration{})
	if err != nil || len(configs.Configurations) == 0 {
		fmt.Println("Configs:", configs) // Print configs
		return fmt.Errorf("failed to get configurations: %v", err)
	}
	fmt.Println("Configs:", configs) // Print configs

	// Modify the first available config
	cfg := configs.Configurations[0]      // Assuming first config; adjust as needed
	cfg.Resolution.Width = int32(width)   // Set new width
	cfg.Resolution.Height = int32(height) // Set new height
	// Note: FPS, bitrate, and GOP might not directly apply here; adjust if needed for media package

	setRequest := &media2.SetVideoEncoderConfiguration{
		Configuration: cfg,
	}
	_, err = mediaService.SetVideoEncoderConfiguration(setRequest)
	if err != nil {
		return fmt.Errorf("failed to apply configuration: %v", err)
	}

	// Add verification step: Get configurations again and check
	updatedConfigs, err := mediaService.GetVideoSourceConfigurations(&media2.GetVideoSourceConfigurations{})
	if err != nil {
		return fmt.Errorf("failed to retrieve updated configurations: %v", err)
	}
	if len(updatedConfigs.Configurations) > 0 {
		updatedCfg := updatedConfigs.Configurations[0]
		if updatedCfg.Bounds.Width == int32(width) && updatedCfg.Bounds.Height == int32(height) {
			fmt.Println("✅ Resolution successfully applied.")
		} else {
			return fmt.Errorf("verification failed: resolution was not updated as expected")
		}
	} else {
		return fmt.Errorf("verification failed: no configurations found after update")
	}

	return nil
}

func main() {
	// Example usage; replace with actual values
	ip := "192.168.1.30" // Your camera IP
	username := "admin"
	password := "Admin123"
	width := 1280
	height := 720
	fps := 30       // Example
	bitrate := 5000 // Example
	gop := 30       // Example

	err := applyResolution2(ip, username, password, width, height, fps, bitrate, gop)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}
