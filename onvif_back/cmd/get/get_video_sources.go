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

	fmt.Println("📹 ONVIF Video Sources Information 📹")
	fmt.Printf("Connecting to camera at %s:%d...\n", *ipPtr, *portPtr)

	// Create and connect the camera using lib functions
	camera := lib.NewCamera(*ipPtr, *portPtr, *userPtr, *passPtr)
	err := camera.Connect()
	if err != nil {
		log.Fatalf("❌ Could not connect to the camera: %v", err)
	}
	fmt.Println("✅ Connected to camera successfully")

	// Get all video sources
	fmt.Println("\n🔍 Getting all video sources...")
	sources, err := lib.GetAllVideoSources(camera)
	if err != nil {
		log.Fatalf("❌ Could not retrieve video sources: %v", err)
	}

	if len(sources) == 0 {
		log.Fatalf("❌ No video sources found")
	}

	// If debug mode is enabled, print raw XML response for examination
	if *debugPtr {
		// Get raw XML response for debugging
		rawResp, err := lib.GetRawVideoSourcesXML(camera)
		if err == nil {
			fmt.Println("\n===== DEBUG: Raw XML Response =====")
			fmt.Println(rawResp)
		}

		// Print the sources as JSON for easier debugging
		fmt.Println("\n===== DEBUG: Sources as JSON =====")
		sourcesJSON, _ := json.MarshalIndent(sources, "", "  ")
		fmt.Println(string(sourcesJSON))
	}

	// Print all video sources with detailed information
	fmt.Println("\n===== Available Video Sources =====")
	for i, source := range sources {
		fmt.Printf("%d. Video Source: %s (Token: %s)\n", i+1, source.Name, source.Token)
		fmt.Printf("   Framerate: %.2f fps\n", source.Framerate)
		fmt.Printf("   Resolution: %dx%d\n", source.Resolution.Width, source.Resolution.Height)

		// Track if we found any imaging settings
		hasImagingSettings := false

		// Display Imaging Settings if available
		if source.Imaging != nil {
			// Check if there's any actual data in the imaging settings
			if source.Imaging.Brightness != 0 || source.Imaging.ColorSaturation != 0 ||
				source.Imaging.Contrast != 0 || source.Imaging.Sharpness != 0 ||
				source.Imaging.BacklightCompensation.Mode != "" ||
				source.Imaging.IrCutFilter != "" ||
				source.Imaging.Exposure.Mode != "" ||
				source.Imaging.Focus.AutoFocusMode != "" ||
				source.Imaging.WideDynamicRange.Mode != "" ||
				source.Imaging.WhiteBalance.Mode != "" {

				hasImagingSettings = true
				fmt.Println("\n   💡 Imaging Settings:")

				// Basic settings
				if source.Imaging.Brightness != 0 {
					fmt.Printf("     • Brightness: %.2f\n", source.Imaging.Brightness)
				}
				if source.Imaging.ColorSaturation != 0 {
					fmt.Printf("     • Color Saturation: %.2f\n", source.Imaging.ColorSaturation)
				}
				if source.Imaging.Contrast != 0 {
					fmt.Printf("     • Contrast: %.2f\n", source.Imaging.Contrast)
				}
				if source.Imaging.Sharpness != 0 {
					fmt.Printf("     • Sharpness: %.2f\n", source.Imaging.Sharpness)
				}

				// Backlight Compensation
				if source.Imaging.BacklightCompensation.Mode != "" {
					fmt.Printf("     • Backlight Compensation: %s (Level: %.2f)\n",
						source.Imaging.BacklightCompensation.Mode,
						source.Imaging.BacklightCompensation.Level)
				}

				// IR Cut Filter
				if source.Imaging.IrCutFilter != "" {
					fmt.Printf("     • IR Cut Filter: %s\n", source.Imaging.IrCutFilter)
				}

				// Exposure
				if source.Imaging.Exposure.Mode != "" {
					fmt.Printf("     • Exposure Mode: %s\n", source.Imaging.Exposure.Mode)
					if source.Imaging.Exposure.Priority != "" {
						fmt.Printf("       - Priority: %s\n", source.Imaging.Exposure.Priority)
					}
					if source.Imaging.Exposure.MinExposureTime > 0 || source.Imaging.Exposure.MaxExposureTime > 0 {
						fmt.Printf("       - Exposure Time Range: %.2f - %.2f μs\n",
							source.Imaging.Exposure.MinExposureTime, source.Imaging.Exposure.MaxExposureTime)
					}
					if source.Imaging.Exposure.MinGain > 0 || source.Imaging.Exposure.MaxGain > 0 {
						fmt.Printf("       - Gain Range: %.2f - %.2f dB\n",
							source.Imaging.Exposure.MinGain, source.Imaging.Exposure.MaxGain)
					}
					if source.Imaging.Exposure.MinIris > 0 || source.Imaging.Exposure.MaxIris > 0 {
						fmt.Printf("       - Iris Range: %.2f - %.2f dB\n",
							source.Imaging.Exposure.MinIris, source.Imaging.Exposure.MaxIris)
					}
					if source.Imaging.Exposure.ExposureTime > 0 {
						fmt.Printf("       - Current Exposure Time: %.2f μs\n", source.Imaging.Exposure.ExposureTime)
					}
					if source.Imaging.Exposure.Gain > 0 {
						fmt.Printf("       - Current Gain: %.2f dB\n", source.Imaging.Exposure.Gain)
					}
					if source.Imaging.Exposure.Iris > 0 {
						fmt.Printf("       - Current Iris: %.2f dB\n", source.Imaging.Exposure.Iris)
					}
				}

				// Focus
				if source.Imaging.Focus.AutoFocusMode != "" {
					fmt.Printf("     • Focus Mode: %s\n", source.Imaging.Focus.AutoFocusMode)
					if source.Imaging.Focus.DefaultSpeed > 0 {
						fmt.Printf("       - Default Speed: %.2f\n", source.Imaging.Focus.DefaultSpeed)
					}
					if source.Imaging.Focus.NearLimit > 0 {
						fmt.Printf("       - Near Limit: %.2f m\n", source.Imaging.Focus.NearLimit)
					}
					if source.Imaging.Focus.FarLimit > 0 {
						fmt.Printf("       - Far Limit: %.2f m\n", source.Imaging.Focus.FarLimit)
					}
				}

				// Wide Dynamic Range
				if source.Imaging.WideDynamicRange.Mode != "" {
					fmt.Printf("     • Wide Dynamic Range: %s (Level: %.2f)\n",
						source.Imaging.WideDynamicRange.Mode, source.Imaging.WideDynamicRange.Level)
				}

				// White Balance
				if source.Imaging.WhiteBalance.Mode != "" {
					fmt.Printf("     • White Balance Mode: %s\n", source.Imaging.WhiteBalance.Mode)
					if source.Imaging.WhiteBalance.CrGain > 0 || source.Imaging.WhiteBalance.CbGain > 0 {
						fmt.Printf("       - Cr Gain: %.2f, Cb Gain: %.2f\n",
							source.Imaging.WhiteBalance.CrGain, source.Imaging.WhiteBalance.CbGain)
					}
				}
			}
		}

		// Display Extended Imaging Settings (v2.0) if available
		if source.Extension != nil && source.Extension.Imaging != nil {
			img := source.Extension.Imaging

			// Check if there's any actual data in the extended imaging settings
			if img.Brightness != 0 || img.ColorSaturation != 0 ||
				img.Contrast != 0 || img.Sharpness != 0 ||
				img.BacklightCompensation.Mode != "" ||
				img.IrCutFilter != "" ||
				img.Exposure.Mode != "" ||
				img.Focus.AutoFocusMode != "" || img.Focus.AFMode != "" ||
				img.WideDynamicRange.Mode != "" ||
				img.WhiteBalance.Mode != "" ||
				img.ImageStabilization.Mode != "" ||
				len(img.IrCutFilterAutoAdjustment) > 0 ||
				img.ToneCompensation.Mode != "" ||
				img.Defogging.Mode != "" ||
				img.NoiseReduction.Level > 0 {

				hasImagingSettings = true
				fmt.Println("\n   🔍 Extended Imaging Settings (v2.0):")

				// Basic settings
				if img.Brightness != 0 {
					fmt.Printf("     • Brightness: %.2f\n", img.Brightness)
				}
				if img.ColorSaturation != 0 {
					fmt.Printf("     • Color Saturation: %.2f\n", img.ColorSaturation)
				}
				if img.Contrast != 0 {
					fmt.Printf("     • Contrast: %.2f\n", img.Contrast)
				}
				if img.Sharpness != 0 {
					fmt.Printf("     • Sharpness: %.2f\n", img.Sharpness)
				}

				// IR Cut Filter
				if img.IrCutFilter != "" {
					fmt.Printf("     • IR Cut Filter: %s\n", img.IrCutFilter)
				}

				// Backlight Compensation
				if img.BacklightCompensation.Mode != "" {
					fmt.Printf("     • Backlight Compensation: %s (Level: %.2f)\n",
						img.BacklightCompensation.Mode, img.BacklightCompensation.Level)
				}

				// Exposure
				if img.Exposure.Mode != "" {
					fmt.Printf("     • Exposure Mode: %s\n", img.Exposure.Mode)
					if img.Exposure.Priority != "" {
						fmt.Printf("       - Priority: %s\n", img.Exposure.Priority)
					}
					if img.Exposure.MinExposureTime > 0 || img.Exposure.MaxExposureTime > 0 {
						fmt.Printf("       - Exposure Time Range: %.2f - %.2f μs\n",
							img.Exposure.MinExposureTime, img.Exposure.MaxExposureTime)
					}
					if img.Exposure.MinGain > 0 || img.Exposure.MaxGain > 0 {
						fmt.Printf("       - Gain Range: %.2f - %.2f dB\n",
							img.Exposure.MinGain, img.Exposure.MaxGain)
					}
					if img.Exposure.MinIris > 0 || img.Exposure.MaxIris > 0 {
						fmt.Printf("       - Iris Range: %.2f - %.2f dB\n",
							img.Exposure.MinIris, img.Exposure.MaxIris)
					}
					if img.Exposure.ExposureTime > 0 {
						fmt.Printf("       - Current Exposure Time: %.2f μs\n", img.Exposure.ExposureTime)
					}
					if img.Exposure.Gain > 0 {
						fmt.Printf("       - Current Gain: %.2f dB\n", img.Exposure.Gain)
					}
					if img.Exposure.Iris > 0 {
						fmt.Printf("       - Current Iris: %.2f dB\n", img.Exposure.Iris)
					}
				}

				// Focus
				if img.Focus.AutoFocusMode != "" || img.Focus.AFMode != "" {
					fmt.Printf("     • Focus Mode: %s\n", img.Focus.AutoFocusMode)
					if img.Focus.AFMode != "" {
						fmt.Printf("       - AF Mode: %s\n", img.Focus.AFMode)
					}
					if img.Focus.DefaultSpeed > 0 {
						fmt.Printf("       - Default Speed: %.2f\n", img.Focus.DefaultSpeed)
					}
					if img.Focus.NearLimit > 0 {
						fmt.Printf("       - Near Limit: %.2f m\n", img.Focus.NearLimit)
					}
					if img.Focus.FarLimit > 0 {
						fmt.Printf("       - Far Limit: %.2f m\n", img.Focus.FarLimit)
					}
				}

				// Wide Dynamic Range
				if img.WideDynamicRange.Mode != "" {
					fmt.Printf("     • Wide Dynamic Range: %s (Level: %.2f)\n",
						img.WideDynamicRange.Mode, img.WideDynamicRange.Level)
				}

				// White Balance
				if img.WhiteBalance.Mode != "" {
					fmt.Printf("     • White Balance Mode: %s\n", img.WhiteBalance.Mode)
					if img.WhiteBalance.CrGain > 0 || img.WhiteBalance.CbGain > 0 {
						fmt.Printf("       - Cr Gain: %.2f, Cb Gain: %.2f\n",
							img.WhiteBalance.CrGain, img.WhiteBalance.CbGain)
					}
				}

				// Image Stabilization
				if img.ImageStabilization.Mode != "" {
					fmt.Printf("     • Image Stabilization: %s (Level: %.2f)\n",
						img.ImageStabilization.Mode, img.ImageStabilization.Level)
				}

				// IR Cut Filter Auto Adjustment
				if len(img.IrCutFilterAutoAdjustment) > 0 {
					fmt.Println("     • IR Cut Filter Auto Adjustments:")
					for i, adj := range img.IrCutFilterAutoAdjustment {
						fmt.Printf("       %d) Type: %s, Offset: %.2f, Response Time: %s\n",
							i+1, adj.BoundaryType, adj.BoundaryOffset, adj.ResponseTime)
					}
				}

				// Tone Compensation
				if img.ToneCompensation.Mode != "" {
					fmt.Printf("     • Tone Compensation: %s (Level: %.2f)\n",
						img.ToneCompensation.Mode, img.ToneCompensation.Level)
				}

				// Defogging
				if img.Defogging.Mode != "" {
					fmt.Printf("     • Defogging: %s (Level: %.2f)\n",
						img.Defogging.Mode, img.Defogging.Level)
				}

				// Noise Reduction
				if img.NoiseReduction.Level > 0 {
					fmt.Printf("     • Noise Reduction Level: %.2f\n", img.NoiseReduction.Level)
				}
			}
		}

		if !hasImagingSettings {
			fmt.Println("\n   ⚠️  No imaging settings provided by this camera model")
			fmt.Println("   → Try using the imaging service to get more details (see onvif_back/cmd/get/get_imaging_settings.go)")
		}

		fmt.Println()
	}

	// Get all video source configurations
	fmt.Println("🔍 Getting all video source configurations...")
	configs, err := lib.GetAllVideoSourceConfigurations(camera)
	if err != nil {
		log.Fatalf("❌ Could not retrieve video source configurations: %v", err)
	}

	// Print all video source configurations
	fmt.Println("\n===== Video Source Configurations =====")
	for i, config := range configs {
		fmt.Printf("%d. %s (Token: %s)\n", i+1, config.Name, config.Token)
		fmt.Printf("   Source Token: %s\n", config.SourceToken)
		fmt.Printf("   UseCount: %d\n", config.UseCount)
		if config.ViewMode != "" {
			fmt.Printf("   View Mode: %s\n", config.ViewMode)
		}
		fmt.Printf("   Bounds: X=%d, Y=%d, Width=%d, Height=%d\n",
			config.Bounds.X, config.Bounds.Y, config.Bounds.Width, config.Bounds.Height)

		// Find and display the corresponding source information
		for _, source := range sources {
			if source.Token == config.SourceToken {
				fmt.Printf("   → Associated Source: %s (Native Resolution: %dx%d, Framerate: %.2f fps)\n",
					source.Name, source.Resolution.Width, source.Resolution.Height, source.Framerate)
				break
			}
		}
		fmt.Println()
	}

	// Summarize the relationship between sources and configurations
	fmt.Println("\n===== Video Source to Configuration Mapping =====")
	for _, source := range sources {
		fmt.Printf("Source '%s' (Token: %s) is used by:\n", source.Name, source.Token)

		usedByAny := false
		for _, config := range configs {
			if config.SourceToken == source.Token {
				fmt.Printf("  • Configuration '%s' (Token: %s)\n", config.Name, config.Token)
				usedByAny = true
			}
		}

		if !usedByAny {
			fmt.Printf("  • Not currently used by any configuration\n")
		}
		fmt.Println()
	}

	fmt.Println("\n===== Summary =====")
	fmt.Printf("Total video sources: %d\n", len(sources))
	fmt.Printf("Total video source configurations: %d\n", len(configs))

	// Explanation of the difference
	fmt.Println("\nNOTE: A video source represents the physical input (camera sensor),")
	fmt.Println("while a video source configuration defines how that input is processed")
	fmt.Println("(cropping, scaling, etc.) before being encoded.")

	// Note about imaging settings
	fmt.Println("\nNOTE: Some cameras do not provide detailed imaging settings through the")
	fmt.Println("GetVideoSources method. To get more detailed imaging settings, you may need")
	fmt.Println("to use the specific Imaging service methods.")

	// Run with debug flag suggestion
	fmt.Println("\nTIP: Run this program with -debug flag to see the raw response from the camera.")
}
