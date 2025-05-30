package main

import (
	"fmt"
	"log"
	"time"

	"github.com/videonext/onvif/profiles/media"
	"github.com/videonext/onvif/soap"
)

func applyResolution(ip, username, password string, width, height, quality, fps int) error {
	client := soap.NewClient(soap.WithTimeout(5 * time.Second))
	client.AddHeader(soap.NewWSSSecurityHeader(username, password, time.Now()))

	mediaService := media.NewMedia(client, fmt.Sprintf("http://%s/onvif/media_service", ip)) // Adjusted for media package

	// Get profiles
	profiles, err := mediaService.GetProfiles(&media.GetProfiles{})
	if err != nil || len(profiles.Profiles) == 0 {
		fmt.Println("Profiles:", profiles) // Print profiles
		return fmt.Errorf("failed to get profiles: %v", err)
	}
	fmt.Println("Profiles:", profiles) // Print profiles

	// Get video source configurations (equivalent to getting configs in the example)
	configs, err := mediaService.GetVideoEncoderConfigurations(&media.GetVideoEncoderConfigurations{})
	if err != nil || len(configs.Configurations) == 0 {
		fmt.Println("Configs:", configs) // Print configs
		return fmt.Errorf("failed to get configurations: %v", err)
	}
	fmt.Println("Configs:", configs) // Print configs

	// Modify the first available config
	cfg := configs.Configurations[0]            // Assuming first config; adjust as needed
	cfg.Resolution.Width = int32(width)         // Set new width
	cfg.Resolution.Height = int32(height)       // Set new height
	cfg.RateControl.FrameRateLimit = int32(fps) // Set new frame rate
	// Note: FPS, bitrate, and GOP might not directly apply here; adjust if needed for media package

	setRequest := &media.SetVideoEncoderConfiguration{
		Configuration:    cfg,
		ForcePersistence: true,
	}
	_, err = mediaService.SetVideoEncoderConfiguration(setRequest)
	if err != nil {
		return fmt.Errorf("failed to apply configuration: %v", err)
	}

	// Add verification step: Get configurations again and check
	updatedConfigs, err := mediaService.GetVideoSourceConfigurations(&media.GetVideoSourceConfigurations{})
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
	width := 2048
	height := 1536
	quality := 6
	fps := 10

	err := applyResolution(ip, username, password, width, height, quality, fps)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}
