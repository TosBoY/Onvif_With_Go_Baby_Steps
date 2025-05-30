package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	lib "onvif_back/lib"
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
	debugPtr := flag.Bool("debug", false, "Enable debug mode to print raw response")
	flag.Parse()

	fmt.Println("📹 ONVIF Video Source Configurations 📹")
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

	// If debug mode is enabled, print configurations as JSON for easier inspection
	if *debugPtr {
		fmt.Println("\n===== DEBUG: Configurations as JSON =====")
		configsJSON, _ := json.MarshalIndent(configs, "", "  ")
		fmt.Println(string(configsJSON))
	}

	// Print all video source configurations with detailed information
	fmt.Println("\n===== Video Source Configurations =====")
	for i, config := range configs {
		fmt.Printf("%d. Configuration: %s (Token: %s)\n", i+1, config.Name, config.Token)
		fmt.Printf("   • Source Token: %s\n", config.SourceToken)
		fmt.Printf("   • UseCount: %d\n", config.UseCount)

		if config.ViewMode != "" {
			fmt.Printf("   • View Mode: %s\n", config.ViewMode)
		}

		fmt.Printf("   • Bounds: X=%d, Y=%d, Width=%d, Height=%d\n",
			config.Bounds.X, config.Bounds.Y, config.Bounds.Width, config.Bounds.Height)
	}

	// Get all video sources to link them with configurations
	fmt.Println("\n🔍 Getting all video sources to link with configurations...")
	sources, err := lib.GetAllVideoSources(camera)
	if err != nil {
		fmt.Printf("⚠️ Could not retrieve video sources: %v\n", err)
		fmt.Println("Skipping source-to-configuration mapping...")
	} else {
		// Print the mapping between sources and configurations
		fmt.Println("\n===== Configuration to Source Mapping =====")
		for _, config := range configs {
			fmt.Printf("Configuration '%s' (Token: %s) uses:\n", config.Name, config.Token)

			foundSource := false
			for _, source := range sources {
				if source.Token == config.SourceToken {
					fmt.Printf("   • Source: %s (Token: %s)\n", source.Name, source.Token)
					fmt.Printf("     Native Resolution: %dx%d, Framerate: %.2f fps\n",
						source.Resolution.Width, source.Resolution.Height, source.Framerate)

					// Calculate if the configuration is using full resolution or cropped
					isFullRes := (config.Bounds.Width == source.Resolution.Width &&
						config.Bounds.Height == source.Resolution.Height &&
						config.Bounds.X == 0 && config.Bounds.Y == 0)

					if isFullRes {
						fmt.Println("     ✅ Using full source resolution")
					} else {
						fmt.Println("     ⚠️ Using cropped resolution")
						cropPercentW := float64(config.Bounds.Width) / float64(source.Resolution.Width) * 100
						cropPercentH := float64(config.Bounds.Height) / float64(source.Resolution.Height) * 100
						fmt.Printf("       - Crop: %.1f%% of width, %.1f%% of height\n",
							cropPercentW, cropPercentH)
					}

					foundSource = true
					break
				}
			}

			if !foundSource {
				fmt.Printf("   • Unknown source (Token: %s) - source may have been removed\n", config.SourceToken)
			}
			fmt.Println()
		}

		fmt.Println("\n===== Source to Configuration Mapping =====")
		for _, source := range sources {
			fmt.Printf("Source '%s' (Token: %s) is used by:\n", source.Name, source.Token)

			usedByAny := false
			for _, config := range configs {
				if config.SourceToken == source.Token {
					fmt.Printf("   • Configuration '%s' (Token: %s)\n", config.Name, config.Token)
					fmt.Printf("     Bounds: X=%d, Y=%d, Width=%d, Height=%d\n",
						config.Bounds.X, config.Bounds.Y, config.Bounds.Width, config.Bounds.Height)
					usedByAny = true
				}
			}

			if !usedByAny {
				fmt.Printf("   • Not currently used by any configuration\n")
			}
			fmt.Println()
		}
	}

	fmt.Println("\n===== Summary =====")
	fmt.Printf("Total video source configurations: %d\n", len(configs))
	if len(sources) > 0 {
		fmt.Printf("Total video sources: %d\n", len(sources))
	}

	// Explanation of video source configurations
	fmt.Println("\nℹ️  INFORMATION:")
	fmt.Println("• A video source represents the physical input (camera sensor)")
	fmt.Println("• A video source configuration defines how that input is processed")
	fmt.Println("  (cropping, scaling, etc.) before being encoded")
	fmt.Println("• The 'Bounds' parameter defines the cropping region from the source")
	fmt.Println("• Multiple configurations can refer to the same source with different")
	fmt.Println("  bounds for different streams/profiles")

	// Usage tips
	fmt.Println("\n💡 TIPS:")
	fmt.Println("• Run with -debug flag to see configuration details in JSON format")
	fmt.Println("• Use cmd/change/interactive_source_change.go to modify source configurations")
}
