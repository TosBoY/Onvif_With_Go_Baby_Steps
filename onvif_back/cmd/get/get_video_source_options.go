package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	lib "onvif_test2/lib"
)

// Default camera connection details
const (
	defaultCameraIP = "192.168.1.12"
	defaultUsername = "admin"
	defaultPassword = "admin123"
)

func main() {
	// Define flags for camera connection parameters
	ipPtr := flag.String("ip", defaultCameraIP, "Camera IP address")
	portPtr := flag.Int("port", 80, "Camera port")
	userPtr := flag.String("user", defaultUsername, "Username")
	passPtr := flag.String("pass", defaultPassword, "Password")
	configTokenPtr := flag.String("config", "", "Optional video source configuration token")
	profileTokenPtr := flag.String("profile", "", "Optional profile token")
	flag.Parse()

	fmt.Println("📹 ONVIF Video Source Configuration Options 📹")
	fmt.Printf("Connecting to camera at %s:%d...\n", *ipPtr, *portPtr)

	// Create and connect the camera using lib functions
	camera := lib.NewCamera(*ipPtr, *portPtr, *userPtr, *passPtr)
	err := camera.Connect()
	if err != nil {
		log.Fatalf("❌ Could not connect to the camera: %v", err)
	}
	fmt.Println("✅ Connected to camera successfully")

	// Get all video source configurations
	fmt.Println("\n🔍 Getting all video source configurations...")
	configs, err := lib.GetAllVideoSourceConfigurations(camera)
	if err != nil {
		log.Fatalf("❌ Could not retrieve video source configurations: %v", err)
	}

	if len(configs) == 0 {
		log.Fatalf("❌ No video source configurations found")
	}

	// Print all video source configurations
	fmt.Println("\n===== Available Video Source Configurations =====")
	for i, config := range configs {
		fmt.Printf("%d. %s (Token: %s)\n", i+1, config.Name, config.Token)
		fmt.Printf("   Source Token: %s\n", config.SourceToken)
		fmt.Printf("   Use Count: %d\n", config.UseCount)
		if config.ViewMode != "" {
			fmt.Printf("   View Mode: %s\n", config.ViewMode)
		}
		fmt.Printf("   Bounds: X=%d, Y=%d, Width=%d, Height=%d\n",
			config.Bounds.X, config.Bounds.Y, config.Bounds.Width, config.Bounds.Height)
		fmt.Println()
	}

	// Get all profiles
	fmt.Println("🔍 Getting all profiles...")
	profiles, err := lib.GetAllProfiles(camera)
	if err != nil {
		log.Fatalf("❌ Could not retrieve profiles: %v", err)
	}

	if len(profiles) == 0 {
		log.Fatalf("❌ No profiles found")
	}

	// Print all profiles
	fmt.Println("\n===== Available Profiles =====")
	for i, profile := range profiles {
		fmt.Printf("%d. %s (Token: %s)\n", i+1, profile.Name, profile.Token)
	}
	fmt.Println()

	// Determine which configuration and profile tokens to use
	configToken := *configTokenPtr
	profileToken := *profileTokenPtr

	// If no config token is provided, use the first one
	if configToken == "" && len(configs) > 0 {
		configToken = configs[0].Token
		fmt.Printf("No configuration token specified, using first one: %s (%s)\n",
			configs[0].Name, configToken)
	}

	// If no profile token is provided, use the first one
	if profileToken == "" && len(profiles) > 0 {
		profileToken = profiles[0].Token
		fmt.Printf("No profile token specified, using first one: %s (%s)\n",
			profiles[0].Name, profileToken)
	}

	// Exit if we don't have the tokens we need
	if configToken == "" {
		fmt.Println("❌ No configuration token available")
		os.Exit(1)
	}

	if profileToken == "" {
		fmt.Println("❌ No profile token available")
		os.Exit(1)
	}

	// Get the video source configuration options
	fmt.Printf("\n🔍 Getting video source configuration options for config '%s' and profile '%s'...\n",
		configToken, profileToken)

	optionsResp, err := lib.GetVideoSourceConfigurationOptions(camera, configToken, profileToken)
	if err != nil {
		log.Fatalf("❌ Could not get video source configuration options: %v", err)
	}

	// Parse the options into more manageable structure
	options := lib.ParseVideoSourceConfigOptions(optionsResp)

	// Print options
	fmt.Println("\n===== Video Source Configuration Options =====")
	fmt.Printf("Maximum Number of Profiles: %d\n", options.MaximumNumberOfProfiles)

	fmt.Println("\nBounds Range:")
	fmt.Printf("  X Range: %d to %d\n",
		options.BoundsRange.XRange.Min, options.BoundsRange.XRange.Max)
	fmt.Printf("  Y Range: %d to %d\n",
		options.BoundsRange.YRange.Min, options.BoundsRange.YRange.Max)
	fmt.Printf("  Width Range: %d to %d\n",
		options.BoundsRange.WidthRange.Min, options.BoundsRange.WidthRange.Max)
	fmt.Printf("  Height Range: %d to %d\n",
		options.BoundsRange.HeightRange.Min, options.BoundsRange.HeightRange.Max)

	if len(options.VideoSourceTokens) > 0 {
		fmt.Println("\nAvailable Video Source Tokens:")
		for i, token := range options.VideoSourceTokens {
			fmt.Printf("  %d. %s\n", i+1, token)
		}
	} else {
		fmt.Println("\nNo video source tokens available")
	}

	// Get the current configuration for comparison
	fmt.Println("\n🔍 Getting current video source configuration details...")
	config, err := lib.GetVideoSourceConfiguration(camera, configToken)
	if err != nil {
		log.Fatalf("❌ Could not get current configuration: %v", err)
	}

	// Print the current configuration details
	fmt.Println("\n===== Current Configuration =====")
	fmt.Printf("Name: %s\n", config.Name)
	fmt.Printf("Token: %s\n", config.Token)
	fmt.Printf("Source Token: %s\n", config.SourceToken)
	fmt.Printf("Use Count: %d\n", config.UseCount)
	if config.ViewMode != "" {
		fmt.Printf("View Mode: %s\n", config.ViewMode)
	}
	fmt.Printf("Bounds: X=%d, Y=%d, Width=%d, Height=%d\n",
		config.Bounds.X, config.Bounds.Y, config.Bounds.Width, config.Bounds.Height)

	// Print usage guidance based on current configuration and available options
	fmt.Println("\n===== Usage Guidance =====")
	fmt.Println("Based on the available options, you can modify this video source configuration with the following boundaries:")

	// Check if current bounds are within allowed ranges
	withinBounds := true
	if config.Bounds.X < options.BoundsRange.XRange.Min || config.Bounds.X > options.BoundsRange.XRange.Max {
		withinBounds = false
		fmt.Printf("⚠️  Current X value (%d) is outside the allowed range (%d-%d)\n",
			config.Bounds.X, options.BoundsRange.XRange.Min, options.BoundsRange.XRange.Max)
	}

	if config.Bounds.Y < options.BoundsRange.YRange.Min || config.Bounds.Y > options.BoundsRange.YRange.Max {
		withinBounds = false
		fmt.Printf("⚠️  Current Y value (%d) is outside the allowed range (%d-%d)\n",
			config.Bounds.Y, options.BoundsRange.YRange.Min, options.BoundsRange.YRange.Max)
	}

	if config.Bounds.Width < options.BoundsRange.WidthRange.Min || config.Bounds.Width > options.BoundsRange.WidthRange.Max {
		withinBounds = false
		fmt.Printf("⚠️  Current Width value (%d) is outside the allowed range (%d-%d)\n",
			config.Bounds.Width, options.BoundsRange.WidthRange.Min, options.BoundsRange.WidthRange.Max)
	}

	if config.Bounds.Height < options.BoundsRange.HeightRange.Min || config.Bounds.Height > options.BoundsRange.HeightRange.Max {
		withinBounds = false
		fmt.Printf("⚠️  Current Height value (%d) is outside the allowed range (%d-%d)\n",
			config.Bounds.Height, options.BoundsRange.HeightRange.Min, options.BoundsRange.HeightRange.Max)
	}

	if withinBounds {
		fmt.Println("✅ Current configuration bounds are within the allowed ranges.")
	}

	// Example command to modify the video source configuration
	fmt.Println("\n===== Example Command =====")
	fmt.Println("To create an interactive tool for changing video source configuration, refer to:")
	fmt.Println("onvif_back/cmd/change/interactive_source_change.go")
}
