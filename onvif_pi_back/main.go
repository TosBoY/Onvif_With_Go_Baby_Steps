package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Config struct {
	PiProxy struct {
		Address string `json:"address"`
		Port    int    `json:"port"`
	} `json:"pi_proxy"`
	Server struct {
		Port int `json:"port"`
	} `json:"server"`
}

var config = Config{
	PiProxy: struct {
		Address string `json:"address"`
		Port    int    `json:"port"`
	}{
		Address: "raspberrypi.local", // Default Pi address, can be changed in config.json
		Port:    8090,                // Default Pi port
	},
	Server: struct {
		Port int `json:"port"`
	}{
		Port: 8091, // Default server port
	},
}

// Global Pi client instance
var globalPiClient *PiProxyClient
var piClientLock sync.Mutex

// Initialize Pi client connection
func initPiClient() error {
	piClientLock.Lock()
	defer piClientLock.Unlock()

	if globalPiClient != nil {
		// Already initialized
		return nil
	}

	// Create new Pi client instance
	globalPiClient = NewPiProxyClient(config.PiProxy.Address, config.PiProxy.Port)

	// Test connection
	if !globalPiClient.IsConnected() {
		globalPiClient = nil
		return fmt.Errorf("failed to connect to Pi proxy at %s:%d",
			config.PiProxy.Address, config.PiProxy.Port)
	}

	log.Printf("Pi proxy client connection initialized successfully at %s:%d",
		config.PiProxy.Address, config.PiProxy.Port)
	return nil
}

// Get the Pi client instance, reconnecting if necessary
func getPiClient() (*PiProxyClient, error) {
	piClientLock.Lock()
	defer piClientLock.Unlock()

	// If we don't have a client yet, initialize it
	if globalPiClient == nil {
		log.Println("Creating new Pi proxy client connection...")
		globalPiClient = NewPiProxyClient(config.PiProxy.Address, config.PiProxy.Port)

		// Test connection
		if !globalPiClient.IsConnected() {
			globalPiClient = nil
			return nil, fmt.Errorf("failed to connect to Pi proxy at %s:%d",
				config.PiProxy.Address, config.PiProxy.Port)
		}
		log.Println("Successfully connected to Pi proxy")
	}

	return globalPiClient, nil
}

// Reset the Pi client connection
func resetPiClient() {
	piClientLock.Lock()
	defer piClientLock.Unlock()
	globalPiClient = nil
	log.Println("Pi proxy client connection reset")
}

// Load configuration from config.json
func loadConfig() {
	file, err := os.Open("config.json")
	if err != nil {
		log.Println("Warning: Could not open config.json, using default settings")
		saveConfig() // Create default config
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		log.Printf("Warning: Error parsing config.json: %v. Using default settings", err)
	}
}

// Save configuration to config.json
func saveConfig() {
	file, err := os.Create("config.json")
	if err != nil {
		log.Printf("Error creating config file: %v", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(config)
	if err != nil {
		log.Printf("Error writing config file: %v", err)
	}
}

func main() {
	// Initialize logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Load configuration
	loadConfig()

	// Initialize Pi client connection at startup
	log.Println("Initializing Pi proxy client connection...")
	err := initPiClient()
	if err != nil {
		log.Printf("Warning: Failed to initialize Pi proxy client: %v", err)
		log.Println("Will try to reconnect on first request")
	}

	r := mux.NewRouter()

	// Define API routes
	r.HandleFunc("/api/camera/info", getCameraInfo).Methods("GET")
	r.HandleFunc("/api/camera/resolutions", getResolutions).Methods("GET")
	r.HandleFunc("/api/camera/change-resolution", changeResolution).Methods("POST")
	r.HandleFunc("/api/camera/stream-url", getStreamURL).Methods("GET")
	r.HandleFunc("/api/camera/launch-vlc", launchVLCWithStream).Methods("POST")
	r.HandleFunc("/api/camera/config", getConfigDetails).Methods("GET")
	r.HandleFunc("/api/camera/device-info", getDeviceInfo).Methods("GET")
	r.HandleFunc("/api/pi/status", getPiStatus).Methods("GET")

	// Add CORS support
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:3000", "http://localhost:5173"}),
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type"}),
	)(r)

	// Start the server
	log.Printf("Starting server on port %d...", config.Server.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Server.Port), corsHandler))
}

// Handler to get camera info via Pi proxy
func getCameraInfo(w http.ResponseWriter, r *http.Request) {
	log.Println("getCameraInfo: Starting camera info request via Pi proxy")

	piClient, err := getPiClient()
	if err != nil {
		log.Printf("getCameraInfo: Failed to get Pi proxy connection: %v", err)
		http.Error(w, fmt.Sprintf("Failed to connect to Pi proxy: %v", err), http.StatusInternalServerError)
		return
	}

	info, err := piClient.GetCameraInfo()
	if err != nil {
		log.Printf("getCameraInfo: Failed to get camera info from Pi proxy: %v", err)
		http.Error(w, fmt.Sprintf("Failed to get camera info: %v", err), http.StatusInternalServerError)
		return
	}

	log.Println("getCameraInfo: Successfully retrieved camera info via Pi proxy")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(info)
}

// Handler to get available resolutions via Pi proxy
func getResolutions(w http.ResponseWriter, r *http.Request) {
	log.Println("getResolutions: Starting request for resolution options via Pi proxy")

	configToken := r.URL.Query().Get("configToken")
	profileToken := r.URL.Query().Get("profileToken")
	log.Printf("getResolutions: Using configToken=%s, profileToken=%s", configToken, profileToken)

	piClient, err := getPiClient()
	if err != nil {
		log.Printf("getResolutions: Failed to get Pi proxy connection: %v", err)
		http.Error(w, fmt.Sprintf("Failed to connect to Pi proxy: %v", err), http.StatusInternalServerError)
		return
	}

	resolutions, err := piClient.GetResolutions(configToken, profileToken)
	if err != nil {
		log.Printf("getResolutions: Failed to get resolutions from Pi proxy: %v", err)
		http.Error(w, fmt.Sprintf("Failed to get resolutions: %v", err), http.StatusInternalServerError)
		return
	}

	log.Println("getResolutions: Successfully retrieved resolutions via Pi proxy")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resolutions)
}

// Handler to change resolution via Pi proxy
func changeResolution(w http.ResponseWriter, r *http.Request) {
	log.Println("changeResolution: Starting resolution change request via Pi proxy")

	var payload map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Printf("changeResolution: Failed to decode request payload: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Convert numeric fields to ensure they match Pi proxy expectations
	if width, ok := payload["width"].(float64); ok {
		payload["width"] = int(width)
	}
	if height, ok := payload["height"].(float64); ok {
		payload["height"] = int(height)
	}
	if frameRate, ok := payload["frameRate"].(float64); ok {
		payload["frameRate"] = int(frameRate)
	}
	if bitRate, ok := payload["bitRate"].(float64); ok {
		payload["bitRate"] = int(bitRate)
	}
	if govLength, ok := payload["govLength"].(float64); ok {
		payload["govLength"] = int(govLength)
	}

	log.Printf("changeResolution: Change request details: %+v", payload)

	piClient, err := getPiClient()
	if err != nil {
		log.Printf("changeResolution: Failed to get Pi proxy connection: %v", err)
		http.Error(w, fmt.Sprintf("Failed to connect to Pi proxy: %v", err), http.StatusInternalServerError)
		return
	}

	err = piClient.ChangeResolution(payload)
	if err != nil {
		log.Printf("changeResolution: Failed to change resolution via Pi proxy: %v", err)
		http.Error(w, fmt.Sprintf("Failed to change resolution: %v", err), http.StatusInternalServerError)
		return
	}

	log.Println("changeResolution: Resolution updated successfully via Pi proxy")
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "Resolution updated successfully"}`))
}

// Handler to get stream URL for a profile via Pi proxy
func getStreamURL(w http.ResponseWriter, r *http.Request) {
	log.Println("getStreamURL: Starting request for stream URL via Pi proxy")

	profileToken := r.URL.Query().Get("profileToken")
	if profileToken == "" {
		http.Error(w, "Missing profileToken parameter", http.StatusBadRequest)
		return
	}

	piClient, err := getPiClient()
	if err != nil {
		log.Printf("getStreamURL: Failed to get Pi proxy connection: %v", err)
		http.Error(w, fmt.Sprintf("Failed to connect to Pi proxy: %v", err), http.StatusInternalServerError)
		return
	}

	streamURI, err := piClient.GetStreamURL(profileToken)
	if err != nil {
		log.Printf("getStreamURL: Failed to get stream URL from Pi proxy: %v", err)
		http.Error(w, fmt.Sprintf("Failed to get stream URL: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"streamUrl": streamURI,
	}

	log.Println("getStreamURL: Successfully retrieved stream URL via Pi proxy")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Handler to get configuration details via Pi proxy
func getConfigDetails(w http.ResponseWriter, r *http.Request) {
	log.Println("getConfigDetails: Starting request for configuration details via Pi proxy")

	configToken := r.URL.Query().Get("configToken")
	if configToken == "" {
		http.Error(w, "Missing configToken parameter", http.StatusBadRequest)
		return
	}

	piClient, err := getPiClient()
	if err != nil {
		log.Printf("getConfigDetails: Failed to get Pi proxy connection: %v", err)
		http.Error(w, fmt.Sprintf("Failed to connect to Pi proxy: %v", err), http.StatusInternalServerError)
		return
	}

	configDetails, err := piClient.GetConfigDetails(configToken)
	if err != nil {
		log.Printf("getConfigDetails: Failed to get configuration details from Pi proxy: %v", err)
		http.Error(w, fmt.Sprintf("Failed to get configuration details: %v", err), http.StatusInternalServerError)
		return
	}

	log.Println("getConfigDetails: Successfully retrieved configuration details via Pi proxy")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(configDetails)
}

// Handler to get device information via Pi proxy
func getDeviceInfo(w http.ResponseWriter, r *http.Request) {
	log.Println("getDeviceInfo: Starting request for device information via Pi proxy")

	piClient, err := getPiClient()
	if err != nil {
		log.Printf("getDeviceInfo: Failed to get Pi proxy connection: %v", err)
		http.Error(w, fmt.Sprintf("Failed to connect to Pi proxy: %v", err), http.StatusInternalServerError)
		return
	}

	deviceInfo, err := piClient.GetDeviceInfo()
	if err != nil {
		log.Printf("getDeviceInfo: Failed to get device information from Pi proxy: %v", err)
		http.Error(w, fmt.Sprintf("Failed to get device information: %v", err), http.StatusInternalServerError)
		return
	}

	log.Println("getDeviceInfo: Successfully retrieved device information via Pi proxy")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deviceInfo)
}

// Handler to get Pi system status
func getPiStatus(w http.ResponseWriter, r *http.Request) {
	log.Println("getPiStatus: Starting request for Pi system status")

	piClient, err := getPiClient()
	if err != nil {
		log.Printf("getPiStatus: Failed to get Pi proxy connection: %v", err)
		http.Error(w, fmt.Sprintf("Failed to connect to Pi proxy: %v", err), http.StatusInternalServerError)
		return
	}

	status, err := piClient.GetSystemStatus()
	if err != nil {
		log.Printf("getPiStatus: Failed to get system status from Pi proxy: %v", err)
		http.Error(w, fmt.Sprintf("Failed to get Pi system status: %v", err), http.StatusInternalServerError)
		return
	}

	log.Println("getPiStatus: Successfully retrieved Pi system status")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

// Handler to launch VLC with a stream
func launchVLCWithStream(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		ProfileToken string `json:"profileToken"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Get the stream URL from the Pi proxy
	piClient, err := getPiClient()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to connect to Pi proxy: %v", err), http.StatusInternalServerError)
		return
	}

	// Get stream URL via the Pi proxy
	streamURI, err := piClient.GetStreamURL(payload.ProfileToken)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get stream URL: %v", err), http.StatusInternalServerError)
		return
	}

	// Check if VLC is already running
	vlcRunning := isVLCRunning()

	// Launch or inject the stream into VLC
	log.Printf("Attempting to %s VLC with stream...\n",
		map[bool]string{true: "inject stream into", false: "launch"}[vlcRunning])

	err = launchOrInjectVLC(streamURI)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to launch or inject stream into VLC: %v", err), http.StatusInternalServerError)
		return
	}

	// Create appropriate success message based on what we did
	actionMsg := "VLC launched successfully"
	if vlcRunning {
		actionMsg = "Stream added to running VLC instance"
	}

	response := map[string]string{
		"message":   actionMsg,
		"streamUrl": streamURI,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Helper function to check if VLC HTTP interface is accessible
func isVLCHttpInterfaceActive() bool {
	// Create client with short timeout to quickly check if HTTP interface responds
	client := &http.Client{
		Timeout: 500 * time.Millisecond,
	}

	// Create request with authentication
	req, err := http.NewRequest("GET", "http://localhost:8080/requests/status.xml", nil)
	if err != nil {
		log.Printf("isVLCHttpInterfaceActive: Failed to create request: %v\n", err)
		return false
	}

	// Add Basic auth - empty username with password "123"
	req.SetBasicAuth("", "123")

	// Try to access a simple VLC HTTP endpoint
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("isVLCHttpInterfaceActive: VLC HTTP interface not responding: %v\n", err)
		return false
	}
	defer resp.Body.Close()

	// Check for successful response
	log.Printf("isVLCHttpInterfaceActive: VLC HTTP interface responded with status: %d\n", resp.StatusCode)
	return resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusUnauthorized
}

// Helper function to launch VLC or inject stream into running instance
func launchOrInjectVLC(rtspURI string) error {
	// First check if VLC is running
	if isVLCRunning() {
		log.Println("launchOrInjectVLC: VLC is already running")

		// Check if HTTP interface is accessible
		if isVLCHttpInterfaceActive() {
			log.Println("launchOrInjectVLC: VLC HTTP interface is active, attempting to inject stream")

			// Try to add the stream to the running instance
			err := addStreamToRunningVLC(rtspURI)
			if err == nil {
				log.Println("launchOrInjectVLC: Successfully added stream to running VLC")
				return nil
			}
			log.Printf("launchOrInjectVLC: Failed to add stream to running VLC: %v\n", err)
		} else {
			log.Println("launchOrInjectVLC: VLC HTTP interface is not active in the running instance")
		}

		// At this point either:
		// 1. HTTP interface is not enabled in the running VLC
		// 2. Or we failed to add the stream
		// Close existing VLC and launch a new one with HTTP enabled
		log.Println("launchOrInjectVLC: Will close current VLC and launch new one with HTTP interface enabled")

		// Close existing VLC instances
		closeVLCInstances()

		// Small delay to ensure VLC has fully closed
		time.Sleep(1 * time.Second)
	}

	// Launch a new VLC instance with HTTP interface enabled
	log.Println("launchOrInjectVLC: Launching new VLC instance with HTTP interface enabled")
	return launchNewVLCInstance(rtspURI)
}

// Helper function to launch a new VLC instance
func launchNewVLCInstance(rtspURI string) error {
	vlcPath, err := getVLCPath()
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
		"--http-port", "8080",
		"--http-password", "123", // Set HTTP password
		rtspURI,
	)

	return cmd.Start()
}

// Add stream to running VLC instance using VLC's HTTP interface
func addStreamToRunningVLC(rtspURI string) error {
	// Add detailed logging
	log.Println("addStreamToRunningVLC: Attempting to add stream to running VLC instance")
	log.Println("addStreamToRunningVLC: Stream URL:", rtspURI)

	// URL encode the RTSP URI to ensure it's properly passed to VLC
	encodedURI := url.QueryEscape(rtspURI)

	// First check the current playlist to see if we need to clear it
	client := &http.Client{
		Timeout: 2 * time.Second,
	}

	// Step 1: Create a request to get the current playlist status
	playlistReq, err := http.NewRequest("GET", "http://localhost:8080/requests/playlist.xml", nil)
	if err != nil {
		log.Printf("addStreamToRunningVLC: Failed to create playlist request: %v\n", err)
		return err
	}
	playlistReq.SetBasicAuth("", "123")

	playlistResp, err := client.Do(playlistReq)
	if err != nil {
		log.Printf("addStreamToRunningVLC: Error getting playlist: %v\n", err)
		return fmt.Errorf("failed to get VLC playlist: %v", err)
	}
	defer playlistResp.Body.Close()

	// Step 2: Clear the current playlist
	clearReq, err := http.NewRequest("GET", "http://localhost:8080/requests/status.xml?command=pl_empty", nil)
	if err != nil {
		log.Printf("addStreamToRunningVLC: Failed to create clear request: %v\n", err)
		return err
	}
	clearReq.SetBasicAuth("", "123")

	clearResp, err := client.Do(clearReq)
	if err != nil {
		log.Printf("addStreamToRunningVLC: Error clearing playlist: %v\n", err)
	} else {
		defer clearResp.Body.Close()
		log.Println("addStreamToRunningVLC: Cleared existing playlist")
	}

	// Step 3: Add the stream to the playlist using in_play instead of in_enqueue
	// This will both add and play the stream in one command
	playURL := fmt.Sprintf("http://localhost:8080/requests/status.xml?command=in_play&input=%s", encodedURI)
	log.Println("addStreamToRunningVLC: Sending play request to VLC HTTP interface:", playURL)

	playReq, err := http.NewRequest("GET", playURL, nil)
	if err != nil {
		log.Printf("addStreamToRunningVLC: Failed to create play request: %v\n", err)
		return err
	}
	playReq.SetBasicAuth("", "123")

	// Try to play the stream via HTTP API with authentication
	playResp, err := client.Do(playReq)
	if err != nil {
		log.Printf("addStreamToRunningVLC: Error connecting to VLC HTTP interface: %v\n", err)
		return fmt.Errorf("failed to connect to VLC HTTP interface: %v", err)
	}
	defer playResp.Body.Close()

	if playResp.StatusCode != http.StatusOK {
		return fmt.Errorf("VLC HTTP interface returned error: %d", playResp.StatusCode)
	}

	log.Println("addStreamToRunningVLC: Successfully added and started playing stream in running VLC instance")
	return nil
}

// Check if VLC is running
func isVLCRunning() bool {
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

// Get the path to VLC executable based on OS
func getVLCPath() (string, error) {
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

// Helper function to close existing VLC instances
func closeVLCInstances() {
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
