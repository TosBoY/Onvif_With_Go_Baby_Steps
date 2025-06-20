package vlc

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

// Constants for VLC HTTP interface
const (
	DefaultHTTPPort     = "8080"
	DefaultHTTPPassword = "123"
)

// LaunchVLCWithStream launches VLC with the provided stream URI
// Returns error if VLC cannot be launched or the stream cannot be played
func LaunchVLCWithStream(streamURI string) (string, error) {
	// Check if VLC is already running
	vlcRunning := IsVLCRunning()
	fmt.Printf("VLC running status: %v\n", vlcRunning)

	// Get VLC path to verify it's installed
	vlcPath, err := GetVLCPath()
	if err != nil {
		return "", fmt.Errorf("VLC is not installed or not found: %v", err)
	}
	fmt.Printf("Found VLC at: %s\n", vlcPath)

	// Launch or inject the stream into VLC
	fmt.Printf("Attempting to %s VLC with stream...\n",
		map[bool]string{true: "inject stream into", false: "launch"}[vlcRunning])

	err = LaunchOrInjectVLC(streamURI)
	if err != nil {
		// Log more details about the failure
		fmt.Printf("Failed to handle VLC: %v\n", err)
		if vlcRunning {
			fmt.Println("Attempting to close existing VLC instances and retry...")
			CloseVLCInstances()
			time.Sleep(1 * time.Second)
			err = LaunchNewVLCInstance(streamURI)
			if err != nil {
				return "", fmt.Errorf("failed to launch new VLC instance after closing existing: %v", err)
			}
			vlcRunning = false
		} else {
			return "", fmt.Errorf("failed to launch VLC: %v", err)
		}
	}

	// Create appropriate success message based on what we did
	actionMsg := "VLC launched successfully"
	if vlcRunning {
		actionMsg = "Stream added to running VLC instance"
	}

	return actionMsg, nil
}

// IsVLCHttpInterfaceActive checks if VLC HTTP interface is accessible
func IsVLCHttpInterfaceActive() bool {
	// Create client with short timeout to quickly check if HTTP interface responds
	client := &http.Client{
		Timeout: 500 * time.Millisecond,
	}

	// Create request with authentication
	req, err := http.NewRequest("GET", "http://localhost:"+DefaultHTTPPort+"/requests/status.xml", nil)
	if err != nil {
		fmt.Printf("IsVLCHttpInterfaceActive: Failed to create request: %v\n", err)
		return false
	}

	// Add Basic auth - empty username with password "123"
	req.SetBasicAuth("", DefaultHTTPPassword)

	// Try to access a simple VLC HTTP endpoint
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("IsVLCHttpInterfaceActive: VLC HTTP interface not responding: %v\n", err)
		return false
	}
	defer resp.Body.Close()

	// Check for successful response
	fmt.Printf("IsVLCHttpInterfaceActive: VLC HTTP interface responded with status: %d\n", resp.StatusCode)
	return resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusUnauthorized
}

// LaunchOrInjectVLC launches VLC or injects stream into running instance
func LaunchOrInjectVLC(rtspURI string) error {
	// First check if VLC is running
	if IsVLCRunning() {
		fmt.Println("LaunchOrInjectVLC: VLC is already running")

		// Check if HTTP interface is accessible
		if IsVLCHttpInterfaceActive() {
			fmt.Println("LaunchOrInjectVLC: VLC HTTP interface is active, attempting to inject stream")

			// Try to add the stream to the running instance
			err := AddStreamToRunningVLC(rtspURI)
			if err == nil {
				fmt.Println("LaunchOrInjectVLC: Successfully added stream to running VLC")
				return nil
			}
			fmt.Printf("LaunchOrInjectVLC: Failed to add stream to running VLC: %v\n", err)
		} else {
			fmt.Println("LaunchOrInjectVLC: VLC HTTP interface is not active in the running instance")
		}

		// At this point either:
		// 1. HTTP interface is not enabled in the running VLC
		// 2. Or we failed to add the stream
		fmt.Println("LaunchOrInjectVLC: Will close current VLC and launch new one with HTTP interface enabled")

		// Close existing VLC instances
		CloseVLCInstances()

		// Small delay to ensure VLC has fully closed
		time.Sleep(1 * time.Second)
	}

	// Launch a new VLC instance with HTTP interface enabled
	fmt.Println("LaunchOrInjectVLC: Launching new VLC instance with HTTP interface enabled")
	return LaunchNewVLCInstance(rtspURI)
}

// LaunchNewVLCInstance launches a new VLC instance with the provided stream URI
func LaunchNewVLCInstance(rtspURI string) error {
	vlcPath, err := GetVLCPath()
	if err != nil {
		return err
	}

	// Launch VLC with the stream URL and enable HTTP interface
	cmd := exec.Command(
		vlcPath,
		"--no-video-title-show",
		"--rtsp-tcp",
		"--extraintf", "http",
		"--http-host", "localhost",
		"--http-port", DefaultHTTPPort,
		"--http-password", DefaultHTTPPassword,
		rtspURI,
	)

	return cmd.Start()
}

// AddStreamToRunningVLC adds a stream to a running VLC instance using VLC's HTTP interface
func AddStreamToRunningVLC(rtspURI string) error {
	// Add detailed logging
	fmt.Println("AddStreamToRunningVLC: Attempting to add stream to running VLC instance")
	fmt.Println("AddStreamToRunningVLC: Stream URL:", rtspURI)

	// VLC HTTP interface default port is 8080
	// URL encode the RTSP URI to ensure it's properly passed to VLC
	encodedURI := url.QueryEscape(rtspURI)

	// First check the current playlist to see if we need to clear it
	client := &http.Client{
		Timeout: 2 * time.Second,
	}

	// Step 1: Create a request to get the current playlist status
	playlistReq, err := http.NewRequest("GET", "http://localhost:"+DefaultHTTPPort+"/requests/playlist.xml", nil)
	if err != nil {
		fmt.Printf("AddStreamToRunningVLC: Failed to create playlist request: %v\n", err)
		return err
	}
	playlistReq.SetBasicAuth("", DefaultHTTPPassword)

	playlistResp, err := client.Do(playlistReq)
	if err != nil {
		fmt.Printf("AddStreamToRunningVLC: Error getting playlist: %v\n", err)
		return fmt.Errorf("failed to get VLC playlist: %v", err)
	}
	defer playlistResp.Body.Close()

	// Step 2: Clear the current playlist
	clearReq, err := http.NewRequest("GET", "http://localhost:"+DefaultHTTPPort+"/requests/status.xml?command=pl_empty", nil)
	if err != nil {
		fmt.Printf("AddStreamToRunningVLC: Failed to create clear request: %v\n", err)
		return err
	}
	clearReq.SetBasicAuth("", DefaultHTTPPassword)

	clearResp, err := client.Do(clearReq)
	if err != nil {
		fmt.Printf("AddStreamToRunningVLC: Error clearing playlist: %v\n", err)
	} else {
		defer clearResp.Body.Close()
		fmt.Println("AddStreamToRunningVLC: Cleared existing playlist")
	}

	// Step 3: Add the stream to the playlist using in_play instead of in_enqueue
	// This will both add and play the stream in one command
	playURL := fmt.Sprintf("http://localhost:%s/requests/status.xml?command=in_play&input=%s", DefaultHTTPPort, encodedURI)
	fmt.Println("AddStreamToRunningVLC: Sending play request to VLC HTTP interface:", playURL)

	playReq, err := http.NewRequest("GET", playURL, nil)
	if err != nil {
		fmt.Printf("AddStreamToRunningVLC: Failed to create play request: %v\n", err)
		return err
	}
	playReq.SetBasicAuth("", DefaultHTTPPassword)

	// Try to play the stream via HTTP API with authentication
	playResp, err := client.Do(playReq)
	if err != nil {
		fmt.Printf("AddStreamToRunningVLC: Error connecting to VLC HTTP interface: %v\n", err)
		return fmt.Errorf("failed to connect to VLC HTTP interface: %v", err)
	}
	defer playResp.Body.Close()

	// Read response body for debugging
	body, _ := io.ReadAll(playResp.Body)
	fmt.Printf("AddStreamToRunningVLC: VLC HTTP interface response code: %d\n", playResp.StatusCode)
	fmt.Printf("AddStreamToRunningVLC: VLC HTTP interface response body: %s\n", string(body))

	if playResp.StatusCode != http.StatusOK {
		return fmt.Errorf("VLC HTTP interface returned error: %d", playResp.StatusCode)
	}

	fmt.Println("AddStreamToRunningVLC: Successfully added and started playing stream in running VLC instance")
	return nil
}

// IsVLCRunning checks if VLC is already running
func IsVLCRunning() bool {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("tasklist", "/FI", "IMAGENAME eq vlc.exe", "/NH")
	case "darwin": // macOS
		cmd = exec.Command("pgrep", "VLC")
	case "linux":
		cmd = exec.Command("pgrep", "vlc")
	}

	output, err := cmd.Output()
	if err != nil {
		// Command failed or no VLC process found
		return false
	}

	// Check if the output contains evidence of VLC running
	return strings.Contains(string(output), "vlc")
}

// GetVLCPath returns the path to VLC executable based on OS
func GetVLCPath() (string, error) {
	var vlcPath string
	switch runtime.GOOS {
	case "windows":
		// Check common installation paths
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
		// Try to find using which command
		out, err := exec.Command("which", "vlc").Output()
		if err == nil {
			vlcPath = strings.TrimSpace(string(out))
		}
	}

	if vlcPath == "" {
		return "", fmt.Errorf("VLC not found on system")
	}
	return vlcPath, nil
}

// CloseVLCInstances closes any existing VLC instances
func CloseVLCInstances() {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("taskkill", "/F", "/IM", "vlc.exe")
	case "darwin": // macOS
		cmd = exec.Command("killall", "VLC")
	case "linux":
		cmd = exec.Command("killall", "vlc")
	}
	_ = cmd.Run()               // Ignore errors since it's fine if no instances are running
	time.Sleep(1 * time.Second) // Increased delay to ensure VLC has fully released resources
}

// Configuration for VLC HTTP interface
type VLCConfig struct {
	Port     string
	Password string
	Host     string
}

// NewDefaultVLCConfig creates a new VLC configuration with default values
func NewDefaultVLCConfig() *VLCConfig {
	return &VLCConfig{
		Port:     DefaultHTTPPort,
		Password: DefaultHTTPPassword,
		Host:     "localhost",
	}
}
