package main

import (
	"flag"
	"fmt"
	lib "onvif_test2/lib"
	"os"
)

func main() {
	// Define flags for camera connection parameters
	ipPtr := flag.String("ip", "192.168.1.12", "Camera IP address")
	portPtr := flag.Int("port", 80, "Camera port")
	userPtr := flag.String("user", "admin", "Username")
	passPtr := flag.String("pass", "admin123", "Password")
	flag.Parse()

	// Connect to the camera
	camera := lib.NewCamera(*ipPtr, *portPtr, *userPtr, *passPtr)
	err := camera.Connect()
	if err != nil {
		fmt.Printf("Failed to connect to camera: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Successfully connected to camera")

	// Get all video encoder configurations
	configs, err := lib.GetAllVideoEncoderConfigurations(camera)
	if err != nil {
		fmt.Printf("Failed to get video encoder configurations: %v\n", err)
		os.Exit(1)
	}

	// Find the "UpdatedConfig" configuration
	var configToRename *lib.VideoEncoderConfig
	for _, config := range configs {
		if config.Name == "UpdatedConfig" {
			configToRename = &config
			break
		}
	}

	if configToRename == nil {
		fmt.Println("Could not find any configuration with name 'UpdatedConfig'")
		fmt.Println("Available configurations:")
		for _, config := range configs {
			fmt.Printf("  - Token: %s, Name: %s\n", config.Token, config.Name)
		}
		os.Exit(1)
	}

	// Display current configuration details
	fmt.Printf("Found configuration to rename:\n")
	fmt.Printf("  Token: %s\n", configToRename.Token)
	fmt.Printf("  Name: %s\n", configToRename.Name)
	fmt.Printf("  Resolution: %dx%d\n", configToRename.Width, configToRename.Height)
	fmt.Printf("  Frame Rate: %d\n", configToRename.FrameRate)
	fmt.Printf("  Bit Rate: %d\n", configToRename.BitRate)

	// Set the desired name
	const desiredName = "VideoEncoderConfig_Channel1_MainStream1"

	// Keep the same settings, just change the name
	err = lib.SetVideoEncoderConfiguration(
		camera,
		configToRename.Token,
		desiredName,
		configToRename.Width,
		configToRename.Height,
		configToRename.FrameRate,
		configToRename.BitRate,
		configToRename.GovLength,
		configToRename.H264Profile,
	)

	if err != nil {
		fmt.Printf("Failed to update configuration name: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully renamed configuration to: %s\n", desiredName)

	// Verify the change
	updatedConfig, err := lib.GetVideoEncoderConfiguration(camera, configToRename.Token)
	if err != nil {
		fmt.Printf("Failed to get updated configuration: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Updated configuration:\n")
	fmt.Printf("  Token: %s\n", updatedConfig.Token)
	fmt.Printf("  Name: %s\n", updatedConfig.Name)
}
