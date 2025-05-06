package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	lib "onvif_test2/lib"
)

// Camera connection details
const (
	cameraIP = "192.168.1.12" // Replace with your camera's IP
	username = "admin"        // Replace with your camera's username
	password = "admin123"     // Replace with your camera's password
)

func main() {
	fmt.Println("🎥 ONVIF Interactive Resolution Changer 🎥")
	fmt.Printf("Connecting to camera at %s...\n", cameraIP)

	// Create and connect the camera using lib functions
	camera := lib.NewCamera(cameraIP, 80, username, password)
	err := camera.Connect()
	if err != nil {
		log.Fatalf("❌ Could not connect to the camera: %v", err)
	}
	fmt.Println("✅ Connected to camera successfully")

	// Get all video encoder configurations using lib function
	configs, err := lib.GetAllVideoEncoderConfigurations(camera)
	if err != nil {
		log.Fatalf("❌ Could not retrieve video encoder configurations: %v", err)
	}

	if len(configs) == 0 {
		log.Fatalf("❌ No video encoder configurations found")
	}

	// Get all profiles using lib function
	profiles, err := lib.GetAllProfiles(camera)
	if err != nil {
		log.Fatalf("❌ Could not retrieve profiles: %v", err)
	}

	if len(profiles) == 0 {
		log.Fatalf("❌ No profiles found")
	}

	fmt.Println("\n===== Available Video Encoder Configurations =====")
	for i, config := range configs {
		fmt.Printf("%d. %s (Token: %s, Profile: %s, Resolution: %dx%d)\n",
			i+1, config.Name, config.Token, config.H264Profile, config.Width, config.Height)
	}

	fmt.Println("\nSelect a configuration to modify (enter the number):")
	selectedConfig := readIntInput(1, len(configs)) - 1
	config := configs[selectedConfig]

	fmt.Println("\n===== Available Profiles =====")
	for i, profile := range profiles {
		fmt.Printf("%d. %s (Token: %s)\n", i+1, profile.Name, profile.Token)
	}

	fmt.Println("\nSelect a profile to associate with (enter the number):")
	selectedProfile := readIntInput(1, len(profiles)) - 1
	profile := profiles[selectedProfile]

	// Get the available options for this config+profile combination using lib function
	fmt.Printf("\nGetting available resolutions for configuration '%s' and profile '%s'...\n",
		config.Name, profile.Name)

	options, err := lib.GetVideoEncoderOptions(camera, config.Token, profile.Token)
	if err != nil {
		log.Fatalf("❌ Could not get video encoder options: %v", err)
	}

	// Extract H264 options using lib function
	h264Options := lib.ParseH264Options(options)
	if len(h264Options.ResolutionsAvailable) == 0 {
		log.Fatalf("❌ No available resolutions found for this configuration and profile")
	}

	// Display available resolutions
	fmt.Println("\n===== Available Resolutions =====")
	for i, res := range h264Options.ResolutionsAvailable {
		fmt.Printf("%d. %dx%d\n", i+1, res.Width, res.Height)
	}

	// Let user select a resolution
	fmt.Println("\nSelect a resolution to apply (enter the number):")
	selectedResolution := readIntInput(1, len(h264Options.ResolutionsAvailable)) - 1
	resolution := h264Options.ResolutionsAvailable[selectedResolution]

	// Display other H264 options
	fmt.Println("\n===== Available H264 Options =====")
	fmt.Printf("Frame Rate Range: %d-%d fps\n",
		h264Options.FrameRateRange.Min, h264Options.FrameRateRange.Max)
	fmt.Printf("GOP Length Range: %d-%d\n",
		h264Options.GovLengthRange.Min, h264Options.GovLengthRange.Max)
	fmt.Printf("Encoding Interval Range: %d-%d\n",
		h264Options.EncodingIntervalRange.Min, h264Options.EncodingIntervalRange.Max)

	// Let user choose frame rate
	fmt.Printf("\nEnter frame rate (%d-%d fps):\n",
		h264Options.FrameRateRange.Min, h264Options.FrameRateRange.Max)
	frameRate := readIntRangeInput(
		h264Options.FrameRateRange.Min,
		h264Options.FrameRateRange.Max)

	// Let user choose bitrate
	fmt.Println("\nEnter bitrate in kbps (e.g., 4096 for 4Mbps):")
	bitRate := readIntInput(256, 20000) // Reasonable range for most cameras

	// Let user choose GOP length
	fmt.Printf("\nEnter GOP length (%d-%d):\n",
		h264Options.GovLengthRange.Min, h264Options.GovLengthRange.Max)
	gopLength := readIntRangeInput(
		h264Options.GovLengthRange.Min,
		h264Options.GovLengthRange.Max)

	// Let user select H264 profile if options are available
	var h264Profile string
	if len(h264Options.H264ProfilesSupported) > 0 {
		fmt.Println("\n===== Available H264 Profiles =====")
		for i, profile := range h264Options.H264ProfilesSupported {
			fmt.Printf("%d. %s\n", i+1, profile)
		}

		fmt.Println("\nSelect an H264 profile (enter the number):")
		selectedH264Profile := readIntInput(1, len(h264Options.H264ProfilesSupported)) - 1
		h264Profile = h264Options.H264ProfilesSupported[selectedH264Profile]
	} else {
		// Use the current profile as fallback
		h264Profile = config.H264Profile
		fmt.Printf("\nNo H264 profiles reported. Using current profile: %s\n", h264Profile)
	}

	// Confirm the changes
	fmt.Println("\n===== Configuration Summary =====")
	fmt.Printf("Configuration: %s (Token: %s)\n", config.Name, config.Token)
	fmt.Printf("Profile: %s (Token: %s)\n", profile.Name, profile.Token)
	fmt.Printf("Resolution: %dx%d\n", resolution.Width, resolution.Height)
	fmt.Printf("Frame Rate: %d fps\n", frameRate)
	fmt.Printf("Bitrate: %d kbps\n", bitRate)
	fmt.Printf("GOP Length: %d\n", gopLength)
	fmt.Printf("H264 Profile: %s\n", h264Profile)

	fmt.Println("\nApply these settings? (y/n):")
	reader := bufio.NewReader(os.Stdin)
	confirm, _ := reader.ReadString('\n')
	confirm = strings.TrimSpace(confirm)
	if confirm != "y" && confirm != "Y" {
		fmt.Println("Operation canceled by user.")
		return
	}

	// Get the current stream URI (before change) using lib functions
	fmt.Println("\nGetting current stream URI (before change)...")
	currentUri, err := lib.GetStreamURI(camera, profile.Token)
	if err != nil {
		fmt.Printf("Warning: Could not get current stream URI: %v\n", err)
	} else {
		fmt.Println("Current Stream URI:", currentUri)
	}

	// Apply the configuration changes using lib function
	fmt.Println("\nApplying configuration changes...")
	err = lib.SetVideoEncoderConfiguration(
		camera,
		config.Token,
		config.Name,
		resolution.Width,
		resolution.Height,
		frameRate,
		bitRate,
		gopLength,
		h264Profile,
	)

	if err != nil {
		log.Fatalf("❌ Failed to apply configuration changes: %v", err)
	}

	fmt.Println("✅ Configuration changes applied successfully")

	// Associate the config with the profile using lib function
	err = lib.AddVideoEncoderConfiguration(camera, profile.Token, config.Token)
	if err != nil {
		log.Fatalf("❌ Failed to add video encoder configuration to profile: %v", err)
	}

	fmt.Println("✅ AddVideoEncoderConfiguration successful")

	// Verify the changes using lib function
	fmt.Println("\nVerifying configuration changes...")
	updatedConfig, err := lib.GetVideoEncoderConfiguration(camera, config.Token)
	if err != nil {
		fmt.Printf("Warning: Could not verify configuration changes: %v\n", err)
	} else {
		fmt.Println("\n===== Updated Configuration =====")
		fmt.Printf("Name: %s\n", updatedConfig.Name)
		fmt.Printf("Resolution: %dx%d\n", updatedConfig.Width, updatedConfig.Height)
		fmt.Printf("Frame Rate: %d fps\n", updatedConfig.FrameRate)
		fmt.Printf("Bitrate: %d kbps\n", updatedConfig.BitRate)
		fmt.Printf("GOP Length: %d\n", updatedConfig.GovLength)
		fmt.Printf("H264 Profile: %s\n", updatedConfig.H264Profile)

		// Check if desired resolution matches the actual resolution
		if updatedConfig.Width == resolution.Width && updatedConfig.Height == resolution.Height {
			fmt.Println("\n✅ Resolution successfully changed to the requested values")
		} else {
			fmt.Printf("\n❌ Resolution does not match the requested values. Current: %dx%d\n",
				updatedConfig.Width, updatedConfig.Height)
		}
	}

	// Get the new stream URI using lib function
	fmt.Println("\nGetting new stream URI...")
	streamUri, err := lib.GetStreamURI(camera, profile.Token)
	if err != nil {
		fmt.Printf("Warning: Could not get stream URI: %v\n", err)
	} else {
		fmt.Println("New Stream URI:", streamUri)

		// Check if the URI changed
		if currentUri != "" && streamUri == currentUri {
			fmt.Println("Note: Stream URI is the same as before the configuration change")
		}

		// Ask if the user wants to open the stream in VLC
		fmt.Println("\nDo you want to open the stream in VLC? (y/n):")
		openStream, _ := reader.ReadString('\n')
		openStream = strings.TrimSpace(openStream)
		if openStream == "y" || openStream == "Y" {
			err := refreshStream(streamUri)
			if err != nil {
				fmt.Printf("Warning: Could not open VLC: %v\n", err)
				fmt.Println("Please manually open the stream URI in VLC or another player")
			}
		}
	}

	fmt.Println("\n🎬 All operations completed successfully 🎬")
}

// refreshStream attempts to refresh the RTSP stream
func refreshStream(rtspUri string) error {
	// Method 1: Try to make a TCP connection to the RTSP server to wake it up
	parts := strings.Split(rtspUri, "//")
	if len(parts) != 2 {
		return fmt.Errorf("invalid RTSP URI format")
	}

	hostPort := strings.Split(parts[1], "/")[0]
	if !strings.Contains(hostPort, ":") {
		hostPort = hostPort + ":554" // Default RTSP port
	}

	fmt.Println("Connecting to RTSP server at", hostPort)
	conn, err := net.DialTimeout("tcp", hostPort, 5*time.Second)
	if err != nil {
		fmt.Println("Warning: Could not connect to RTSP server for refresh:", err)
	} else {
		conn.Close()
		fmt.Println("Successfully connected to RTSP server")
	}

	// Close any existing VLC instances before opening a new one
	fmt.Println("Closing any existing VLC instances...")
	closeVLCInstances()

	// Open in a new VLC instance
	fmt.Println("Opening stream in a new VLC instance...")
	if err := openStreamInVLC(rtspUri); err != nil {
		fmt.Println("Warning: Could not launch VLC:", err)
		fmt.Println("RTSP Stream URI:", rtspUri)
		fmt.Println("Please manually refresh your VLC player or use the URI above to open in VLC")
	}

	return nil
}

// closeVLCInstances attempts to close all running VLC instances
func closeVLCInstances() error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		// On Windows, use taskkill to forcefully terminate VLC processes
		cmd = exec.Command("taskkill", "/F", "/IM", "vlc.exe")
	case "darwin": // macOS
		cmd = exec.Command("pkill", "-x", "VLC")
	case "linux":
		cmd = exec.Command("pkill", "-x", "vlc")
	default:
		return fmt.Errorf("unsupported operating system")
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		// If the error is because no VLC instances are running, that's fine
		fmt.Printf("Note: No VLC instances were running or could not close VLC: %v\n", err)
		fmt.Println("Command output:", string(output))
		// Return nil because not finding VLC to kill isn't an error for our purpose
		return nil
	}

	fmt.Println("Successfully closed existing VLC instances")
	// Add a small delay to ensure VLC has time to fully close
	time.Sleep(1 * time.Second)
	return nil
}

// openStreamInVLC tries to open the stream in VLC
func openStreamInVLC(rtspUri string) error {
	// VLC path varies by OS
	var vlcPath string

	switch runtime.GOOS {
	case "windows":
		// Common VLC paths on Windows
		possiblePaths := []string{
			"C:\\Program Files\\VideoLAN\\VLC\\vlc.exe",
			"C:\\Program Files (x86)\\VideoLAN\\VLC\\vlc.exe",
		}

		for _, path := range possiblePaths {
			if _, err := os.Stat(path); err == nil {
				vlcPath = path
				break
			}
		}

	case "darwin": // macOS
		vlcPath = "/Applications/VLC.app/Contents/MacOS/VLC"

	case "linux":
		// Use which command to find VLC
		out, err := exec.Command("which", "vlc").Output()
		if err == nil {
			vlcPath = strings.TrimSpace(string(out))
		}
	}

	if vlcPath == "" {
		return fmt.Errorf("VLC not found on your system")
	}

	fmt.Println("Found VLC at:", vlcPath)
	fmt.Println("Attempting to open stream in VLC...")

	// Build command to launch VLC with minimal interface and the RTSP stream
	cmd := exec.Command(vlcPath, "--no-video-title-show", "--play-and-exit", rtspUri)

	// Run in background
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start VLC: %v", err)
	}

	fmt.Println("VLC launched successfully with the stream")
	return nil
}

// Helper function to read integer input from console
func readIntInput(min, max int) int {
	reader := bufio.NewReader(os.Stdin)
	for {
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		choice, err := strconv.Atoi(input)
		if err != nil || choice < min || choice > max {
			fmt.Printf("Please enter a number between %d and %d: ", min, max)
			continue
		}
		return choice
	}
}

// Helper function to read integer input from console with range validation
func readIntRangeInput(min, max int) int {
	fmt.Printf("Enter a value between %d and %d: ", min, max)
	return readIntInput(min, max)
}
