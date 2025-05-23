package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	lib "onvif_test2/lib"
	"onvif_test2/lib/validator"

	"github.com/use-go/onvif/device"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type CameraConfig struct {
	IP       string `json:"ip"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Camera struct {
	ID       string `json:"id"`
	IP       string `json:"ip"`
	Username string `json:"username"`
	Password string `json:"password"`
	IsFake   bool   `json:"isFake"`
}

type CameraList struct {
	Cameras []Camera `json:"cameras"`
}

// Fixed camera config that's known to work
var cameraConfig = CameraConfig{
	IP:       "192.168.1.12",
	Username: "admin",
	Password: "admin123",
}

// Global camera instance and cameras list
var (
	globalCamera *lib.Camera
	cameraLock   sync.Mutex
	cameras      []Camera
)

// Initialize camera connection
func initCamera() error {
	cameraLock.Lock()
	defer cameraLock.Unlock()

	if globalCamera != nil {
		// Already initialized
		return nil
	}

	// Create new camera instance
	globalCamera = lib.NewCamera(cameraConfig.IP, 80, cameraConfig.Username, cameraConfig.Password)

	// Connect to the camera
	err := globalCamera.Connect()
	if err != nil {
		globalCamera = nil
		return fmt.Errorf("failed to connect to camera: %v", err)
	}

	log.Println("Camera connection initialized successfully")
	return nil
}

// Get the camera instance, reconnecting if necessary
func getCamera() (*lib.Camera, error) {
	cameraLock.Lock()
	defer cameraLock.Unlock()

	// If we don't have a camera yet, initialize it
	if globalCamera == nil {
		log.Println("Creating new camera connection...")
		globalCamera = lib.NewCamera(cameraConfig.IP, 80, cameraConfig.Username, cameraConfig.Password)

		if err := globalCamera.Connect(); err != nil {
			globalCamera = nil
			return nil, fmt.Errorf("failed to connect to camera: %v", err)
		}
		log.Println("Successfully connected to camera")
	}

	return globalCamera, nil
}

func loadCamerasFromConfig() ([]Camera, error) {
	file, err := os.ReadFile("./config/cameras.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read cameras.json: %v", err)
	}

	var cameraList CameraList
	if err := json.Unmarshal(file, &cameraList); err != nil {
		return nil, fmt.Errorf("failed to parse cameras.json: %v", err)
	}

	return cameraList.Cameras, nil
}

// Initialize default cameras if needed
func initDefaultCameras() {
	if len(cameras) == 0 {
		loadedCameras, err := loadCamerasFromConfig()
		if err != nil {
			log.Printf("Warning: Failed to load cameras from config: %v. Using default camera.", err)
			cameras = []Camera{
				{
					ID:       "1",
					IP:       "192.168.1.12",
					Username: "admin",
					Password: "admin123",
					IsFake:   false,
				},
			}
		} else {
			cameras = loadedCameras
			log.Printf("Loaded %d cameras from config file", len(cameras))
		}
	}
}

func getCamerasHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cameras)
}

func applyConfigHandler(w http.ResponseWriter, r *http.Request) {
	var selectedCameras []string
	if err := json.NewDecoder(r.Body).Decode(&selectedCameras); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	for _, cameraID := range selectedCameras {
		for _, camera := range cameras {
			if camera.ID == cameraID {
				log.Printf("Applying configuration to camera: %s (IP: %s, IsFake: %v)", camera.ID, camera.IP, camera.IsFake)
				// Add logic to apply configuration to the camera here
			}
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Configuration applied successfully"))
}

// getNextCameraID returns the next available camera ID
func getNextCameraID(cameras []Camera) string {
	if len(cameras) == 0 {
		return "1"
	}

	// Convert all IDs to integers and find the highest one
	maxID := 0
	for _, cam := range cameras {
		if id, err := strconv.Atoi(cam.ID); err == nil {
			if id > maxID {
				maxID = id
			}
		}
	}

	// Return the next ID as a string
	return strconv.Itoa(maxID + 1)
}

func addCameraHandler(w http.ResponseWriter, r *http.Request) {
	var newCamera Camera
	if err := json.NewDecoder(r.Body).Decode(&newCamera); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// If it's not a fake camera, try to connect to verify credentials
	if !newCamera.IsFake {
		testCamera := lib.NewCamera(newCamera.IP, 80, newCamera.Username, newCamera.Password)
		if err := testCamera.Connect(); err != nil {
			http.Error(w, fmt.Sprintf("Failed to connect to camera: %v", err), http.StatusBadRequest)
			return
		}
	}

	// Read existing cameras
	content, err := os.ReadFile("config/cameras.json")
	if err != nil {
		http.Error(w, "Failed to read camera configuration", http.StatusInternalServerError)
		return
	}

	var cameraList CameraList
	if err := json.Unmarshal(content, &cameraList); err != nil {
		http.Error(w, "Failed to parse camera configuration", http.StatusInternalServerError)
		return
	}

	// Generate the next available ID
	newCamera.ID = getNextCameraID(cameraList.Cameras)

	// Add the new camera to the list
	cameraList.Cameras = append(cameraList.Cameras, newCamera)

	// Create a file with proper formatting
	configFile, err := os.Create("config/cameras.json")
	if err != nil {
		http.Error(w, "Failed to create camera configuration file", http.StatusInternalServerError)
		return
	}
	defer configFile.Close()

	encoder := json.NewEncoder(configFile)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(cameraList); err != nil {
		http.Error(w, "Failed to save camera configuration", http.StatusInternalServerError)
		return
	}

	// Add to current cameras list
	cameras = append(cameras, newCamera)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newCamera)
}

func deleteCameraHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cameraID := vars["id"]

	// Read existing cameras
	content, err := os.ReadFile("config/cameras.json")
	if err != nil {
		http.Error(w, "Failed to read camera configuration", http.StatusInternalServerError)
		return
	}

	var cameraList CameraList
	if err := json.Unmarshal(content, &cameraList); err != nil {
		http.Error(w, "Failed to parse camera configuration", http.StatusInternalServerError)
		return
	}

	// Find and remove the camera
	var found bool
	var updatedCameras []Camera
	for _, cam := range cameraList.Cameras {
		if cam.ID != cameraID {
			updatedCameras = append(updatedCameras, cam)
		} else {
			found = true
		}
	}

	if !found {
		http.Error(w, "Camera not found", http.StatusNotFound)
		return
	}
	// Save updated camera list with proper formatting
	cameraList.Cameras = updatedCameras

	// Create a file with proper formatting
	configFile, err := os.Create("config/cameras.json")
	if err != nil {
		http.Error(w, "Failed to create camera configuration file", http.StatusInternalServerError)
		return
	}
	defer configFile.Close()

	// Use encoder with indentation for pretty JSON
	encoder := json.NewEncoder(configFile)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(cameraList); err != nil {
		http.Error(w, "Failed to save camera configuration", http.StatusInternalServerError)
		return
	}

	// Update the in-memory cameras list
	cameras = updatedCameras

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Camera deleted successfully"})
}

func main() {
	// Initialize with default camera
	initDefaultCameras()

	// Initialize camera connection at startup
	log.Println("Initializing camera connection...")
	err := initCamera()
	if err != nil {
		log.Printf("Warning: Failed to initialize camera connection: %v", err)
		log.Println("Will try to reconnect on first request")
	}

	r := mux.NewRouter()
	// Define API routes
	r.HandleFunc("/api/camera/info", getCameraInfo).Methods("GET")
	r.HandleFunc("/api/camera/resolutions", getResolutions).Methods("GET")
	r.HandleFunc("/api/camera/change-resolution", changeResolution).Methods("POST")
	r.HandleFunc("/api/camera/change-resolution-simple", changeResolutionSimple).Methods("POST")
	r.HandleFunc("/api/camera/stream-url", getStreamURL).Methods("GET")
	r.HandleFunc("/api/camera/launch-vlc", launchVLCWithStream).Methods("POST")
	r.HandleFunc("/api/camera/config", getConfigDetails).Methods("GET")
	r.HandleFunc("/api/camera/device-info", getDeviceInfo).Methods("GET")
	r.HandleFunc("/api/cameras", getCamerasHandler).Methods("GET")
	r.HandleFunc("/api/cameras", addCameraHandler).Methods("POST")
	r.HandleFunc("/api/cameras/{id}", deleteCameraHandler).Methods("DELETE")
	r.HandleFunc("/api/apply-config", applyConfigHandler).Methods("POST")
	r.HandleFunc("/api/camera/resolutions-simple", getResolutionsSimple).Methods("GET")
	r.HandleFunc("/api/camera/details-simple", getCameraDetailsSimple).Methods("GET")

	// Add CORS support
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:3000", "http://localhost:5173"}),
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS", "DELETE"}),
		handlers.AllowedHeaders([]string{"Content-Type"}),
	)(r)

	// Start the server
	log.Println("Starting server on :8090...")
	log.Fatal(http.ListenAndServe(":8090", corsHandler))
}

// Handler to get camera info
func getCameraInfo(w http.ResponseWriter, r *http.Request) {
	// Add debugging to understand exactly what's happening
	fmt.Println("getCameraInfo: Starting camera info request")

	camera, err := getCamera()
	if err != nil {
		fmt.Printf("getCameraInfo: Failed to get camera connection: %v\n", err)
		http.Error(w, fmt.Sprintf("Failed to connect to camera: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Println("getCameraInfo: Successfully got camera connection, retrieving configurations...")
	configs, err := lib.GetAllVideoEncoderConfigurations(camera)
	if err != nil {
		fmt.Printf("getCameraInfo: Failed to get configurations: %v\n", err)
		http.Error(w, fmt.Sprintf("Failed to get configurations: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Println("getCameraInfo: Successfully retrieved configurations, getting profiles...")
	profiles, err := lib.GetAllProfiles(camera)
	if err != nil {
		fmt.Printf("getCameraInfo: Failed to get profiles: %v\n", err)
		http.Error(w, fmt.Sprintf("Failed to get profiles: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Println("getCameraInfo: Successfully retrieved profiles, sending response")
	response := map[string]interface{}{
		"configs":  configs,
		"profiles": profiles,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Handler to get available resolutions
func getResolutions(w http.ResponseWriter, r *http.Request) {
	fmt.Println("getResolutions: Starting request for resolution options")

	camera, err := getCamera()
	if err != nil {
		fmt.Printf("getResolutions: Failed to get camera connection: %v\n", err)
		http.Error(w, fmt.Sprintf("Failed to connect to camera: %v", err), http.StatusInternalServerError)
		return
	}

	configToken := r.URL.Query().Get("configToken")
	profileToken := r.URL.Query().Get("profileToken")
	fmt.Printf("getResolutions: Using configToken=%s, profileToken=%s\n", configToken, profileToken)

	options, err := lib.GetVideoEncoderOptions(camera, configToken, profileToken)
	if err != nil {
		fmt.Printf("getResolutions: Failed to get video encoder options: %v\n", err)
		http.Error(w, fmt.Sprintf("Failed to get resolutions: %v", err), http.StatusInternalServerError)
		return
	}

	// Log the response structure before parsing
	h264 := options.Body.GetVideoEncoderConfigurationOptionsResponse.Options.H264
	fmt.Printf("getResolutions: Raw H264 options before parsing: %+v\n", h264)
	fmt.Printf("getResolutions: Found %d resolutions in raw response\n", len(h264.ResolutionsAvailable))
	for _, res := range h264.ResolutionsAvailable {
		fmt.Printf("getResolutions: Resolution from raw response: %dx%d\n", res.Width, res.Height)
	}

	h264Options := lib.ParseH264Options(options)
	fmt.Printf("getResolutions: Parsed options result: %+v\n", h264Options)

	// Set content type header and write the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(h264Options)
}

// Handler to get available resolutions without needing profile/config selection
func getResolutionsSimple(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\ngetResolutionsSimple: Starting request for resolution options")

	camera, err := getCamera()
	if err != nil {
		fmt.Printf("getResolutionsSimple: Failed to connect to camera: %v\n", err)
		http.Error(w, fmt.Sprintf("Failed to connect to camera: %v", err), http.StatusInternalServerError)
		return
	}

	// Get active profile and config
	profileToken, configToken, err := lib.GetActiveProfile(camera)
	if err != nil {
		fmt.Printf("getResolutionsSimple: Failed to get active profile/config: %v\n", err)
		http.Error(w, fmt.Sprintf("Failed to get active profile/config: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Printf("getResolutionsSimple: Using profileToken=%s, configToken=%s\n", profileToken, configToken)

	options, err := lib.GetVideoEncoderOptions(camera, configToken, profileToken)
	if err != nil {
		fmt.Printf("getResolutionsSimple: Failed to get video encoder options: %v\n", err)
		http.Error(w, fmt.Sprintf("Failed to get resolutions: %v", err), http.StatusInternalServerError)
		return
	}

	// Log available resolutions
	h264 := options.Body.GetVideoEncoderConfigurationOptionsResponse.Options.H264
	fmt.Printf("getResolutionsSimple: Found %d supported resolutions:\n", len(h264.ResolutionsAvailable))
	for _, res := range h264.ResolutionsAvailable {
		fmt.Printf("- %dx%d\n", res.Width, res.Height)
	}

	h264Options := lib.ParseH264Options(options)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(h264Options)
}

// Handler to get stream URL for a profile
func getStreamURL(w http.ResponseWriter, r *http.Request) {
	profileToken := r.URL.Query().Get("profileToken")
	if profileToken == "" {
		http.Error(w, "Missing profileToken parameter", http.StatusBadRequest)
		return
	}

	camera, err := getCamera()
	if err != nil {
		http.Error(w, "Failed to connect to camera", http.StatusInternalServerError)
		return
	}

	streamURI, err := lib.GetStreamURI(camera, profileToken)
	if err != nil {
		http.Error(w, "Failed to get stream URI", http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"streamUrl": streamURI,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Handler to launch VLC with a stream
func launchVLCWithStream(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		CameraId string `json:"cameraId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request payload: %v", err), http.StatusBadRequest)
		return
	}

	// Find the requested camera
	var targetCamera *Camera
	for _, cam := range cameras {
		if cam.ID == payload.CameraId {
			targetCamera = &cam
			break
		}
	}

	if targetCamera == nil {
		http.Error(w, "Camera not found", http.StatusNotFound)
		return
	}

	if targetCamera.IsFake {
		http.Error(w, "Cannot stream from simulated camera", http.StatusBadRequest)
		return
	}

	camera := lib.NewCamera(targetCamera.IP, 80, targetCamera.Username, targetCamera.Password)
	if err := camera.Connect(); err != nil {
		http.Error(w, fmt.Sprintf("Failed to connect to camera: %v", err), http.StatusInternalServerError)
		return
	}

	// Get active profile
	profileToken, _, err := lib.GetActiveProfile(camera)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get active profile: %v", err), http.StatusInternalServerError)
		return
	}

	streamURI, err := lib.GetStreamURI(camera, profileToken)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get stream URI: %v", err), http.StatusInternalServerError)
		return
	}

	// Add authentication credentials to the stream URL
	parsedURL, err := url.Parse(streamURI)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse stream URI: %v", err), http.StatusInternalServerError)
		return
	}

	// Add username and password to the URL
	parsedURL.User = url.UserPassword(targetCamera.Username, targetCamera.Password)
	authenticatedStreamURI := parsedURL.String()

	fmt.Printf("Got stream URI with auth: %s\n", authenticatedStreamURI)

	// Check if VLC is already running
	vlcRunning := isVLCRunning()
	fmt.Printf("VLC running status: %v\n", vlcRunning)

	// Get VLC path to verify it's installed
	vlcPath, err := getVLCPath()
	if err != nil {
		http.Error(w, fmt.Sprintf("VLC is not installed or not found: %v", err), http.StatusInternalServerError)
		return
	}
	fmt.Printf("Found VLC at: %s\n", vlcPath)

	// Launch or inject the stream into VLC
	fmt.Printf("Attempting to %s VLC with stream...\n",
		map[bool]string{true: "inject stream into", false: "launch"}[vlcRunning])

	err = launchOrInjectVLC(authenticatedStreamURI)
	if err != nil {
		// Log more details about the failure
		fmt.Printf("Failed to handle VLC: %v\n", err)
		if vlcRunning {
			fmt.Println("Attempting to close existing VLC instances and retry...")
			closeVLCInstances()
			time.Sleep(1 * time.Second)
			err = launchNewVLCInstance(authenticatedStreamURI)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to launch new VLC instance after closing existing: %v", err), http.StatusInternalServerError)
				return
			}
			vlcRunning = false
		} else {
			http.Error(w, fmt.Sprintf("Failed to launch VLC: %v", err), http.StatusInternalServerError)
			return
		}
	}

	// Create appropriate success message based on what we did
	actionMsg := "VLC launched successfully"
	if vlcRunning {
		actionMsg = "Stream added to running VLC instance"
	}

	response := map[string]string{
		"message":   actionMsg,
		"streamUrl": authenticatedStreamURI,
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
		fmt.Printf("isVLCHttpInterfaceActive: Failed to create request: %v\n", err)
		return false
	}

	// Add Basic auth - empty username with password "123"
	req.SetBasicAuth("", "123")

	// Try to access a simple VLC HTTP endpoint
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("isVLCHttpInterfaceActive: VLC HTTP interface not responding: %v\n", err)
		return false
	}
	defer resp.Body.Close()

	// Check for successful response
	fmt.Printf("isVLCHttpInterfaceActive: VLC HTTP interface responded with status: %d\n", resp.StatusCode)
	return resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusUnauthorized
}

// Helper function to launch VLC or inject stream into running instance
func launchOrInjectVLC(rtspURI string) error {
	// First check if VLC is running
	if isVLCRunning() {
		fmt.Println("launchOrInjectVLC: VLC is already running")

		// Check if HTTP interface is accessible
		if isVLCHttpInterfaceActive() {
			fmt.Println("launchOrInjectVLC: VLC HTTP interface is active, attempting to inject stream")

			// Try to add the stream to the running instance
			err := addStreamToRunningVLC(rtspURI)
			if err == nil {
				fmt.Println("launchOrInjectVLC: Successfully added stream to running VLC")
				return nil
			}
			fmt.Printf("launchOrInjectVLC: Failed to add stream to running VLC: %v\n", err)
		} else {
			fmt.Println("launchOrInjectVLC: VLC HTTP interface is not active in the running instance")
		}

		// At this point either:
		// 1. HTTP interface is not enabled in the running VLC
		// 2. Or we failed to add the stream
		// Ask user if they want to close current VLC and launch a new one with HTTP enabled
		fmt.Println("launchOrInjectVLC: Will close current VLC and launch new one with HTTP interface enabled")

		// Close existing VLC instances
		closeVLCInstances()

		// Small delay to ensure VLC has fully closed
		time.Sleep(1 * time.Second)
	}

	// Launch a new VLC instance with HTTP interface enabled
	fmt.Println("launchOrInjectVLC: Launching new VLC instance with HTTP interface enabled")
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
	fmt.Println("addStreamToRunningVLC: Attempting to add stream to running VLC instance")
	fmt.Println("addStreamToRunningVLC: Stream URL:", rtspURI)

	// VLC HTTP interface default port is 8080
	// URL encode the RTSP URI to ensure it's properly passed to VLC
	encodedURI := url.QueryEscape(rtspURI)

	// First check the current playlist to see if we need to clear it
	client := &http.Client{
		Timeout: 2 * time.Second,
	}

	// Step 1: Create a request to get the current playlist status
	playlistReq, err := http.NewRequest("GET", "http://localhost:8080/requests/playlist.xml", nil)
	if err != nil {
		fmt.Printf("addStreamToRunningVLC: Failed to create playlist request: %v\n", err)
		return err
	}
	playlistReq.SetBasicAuth("", "123")

	playlistResp, err := client.Do(playlistReq)
	if err != nil {
		fmt.Printf("addStreamToRunningVLC: Error getting playlist: %v\n", err)
		return fmt.Errorf("failed to get VLC playlist: %v", err)
	}
	defer playlistResp.Body.Close()

	// Step 2: Clear the current playlist
	clearReq, err := http.NewRequest("GET", "http://localhost:8080/requests/status.xml?command=pl_empty", nil)
	if err != nil {
		fmt.Printf("addStreamToRunningVLC: Failed to create clear request: %v\n", err)
		return err
	}
	clearReq.SetBasicAuth("", "123")

	clearResp, err := client.Do(clearReq)
	if err != nil {
		fmt.Printf("addStreamToRunningVLC: Error clearing playlist: %v\n", err)
	} else {
		defer clearResp.Body.Close()
		fmt.Println("addStreamToRunningVLC: Cleared existing playlist")
	}

	// Step 3: Add the stream to the playlist using in_play instead of in_enqueue
	// This will both add and play the stream in one command
	playURL := fmt.Sprintf("http://localhost:8080/requests/status.xml?command=in_play&input=%s", encodedURI)
	fmt.Println("addStreamToRunningVLC: Sending play request to VLC HTTP interface:", playURL)

	playReq, err := http.NewRequest("GET", playURL, nil)
	if err != nil {
		fmt.Printf("addStreamToRunningVLC: Failed to create play request: %v\n", err)
		return err
	}
	playReq.SetBasicAuth("", "123")

	// Try to play the stream via HTTP API with authentication
	playResp, err := client.Do(playReq)
	if err != nil {
		fmt.Printf("addStreamToRunningVLC: Error connecting to VLC HTTP interface: %v\n", err)
		return fmt.Errorf("failed to connect to VLC HTTP interface: %v", err)
	}
	defer playResp.Body.Close()

	// Read response body for debugging
	body, _ := io.ReadAll(playResp.Body)
	fmt.Printf("addStreamToRunningVLC: VLC HTTP interface response code: %d\n", playResp.StatusCode)
	fmt.Printf("addStreamToRunningVLC: VLC HTTP interface response body: %s\n", string(body))

	if playResp.StatusCode != http.StatusOK {
		return fmt.Errorf("VLC HTTP interface returned error: %d", playResp.StatusCode)
	}

	fmt.Println("addStreamToRunningVLC: Successfully added and started playing stream in running VLC instance")
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

// Handler to change resolution
func changeResolution(w http.ResponseWriter, r *http.Request) {
	fmt.Println("changeResolution: Starting resolution change request")

	var payload struct {
		ConfigToken  string   `json:"configToken"`
		ProfileToken string   `json:"profileToken"`
		CameraIds    []string `json:"cameraIds"`
		Width        int      `json:"width"`
		Height       int      `json:"height"`
		FrameRate    int      `json:"frameRate"`
		BitRate      int      `json:"bitRate"`
		GopLength    int      `json:"gopLength"`
		H264Profile  string   `json:"h264Profile"`
		ConfigName   string   `json:"configName"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		fmt.Printf("changeResolution: Failed to decode request payload: %v\n", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	type UpdateResult struct {
		CameraId         string                      `json:"cameraId"`
		Success          bool                        `json:"success"`
		Error            string                      `json:"error,omitempty"`
		ValidationResult *validator.ValidationResult `json:"validationResult,omitempty"`
	}
	results := make([]UpdateResult, 0)

	expectedConfig := validator.VideoConfig{
		Width:       payload.Width,
		Height:      payload.Height,
		FrameRate:   payload.FrameRate,
		BitRate:     payload.BitRate,
		GopLength:   payload.GopLength,
		H264Profile: payload.H264Profile,
	}

	for _, cameraId := range payload.CameraIds {
		var targetCamera *Camera
		for _, cam := range cameras {
			if cam.ID == cameraId {
				targetCamera = &cam
				break
			}
		}
		if targetCamera == nil {
			results = append(results, UpdateResult{
				CameraId: cameraId,
				Success:  false,
				Error:    "Camera not found",
			})
			continue
		}

		if targetCamera.IsFake {
			results = append(results, UpdateResult{
				CameraId: cameraId,
				Success:  true,
				Error:    "Skipped: This is a simulated camera",
			})
			continue
		}

		camera := lib.NewCamera(targetCamera.IP, 80, targetCamera.Username, targetCamera.Password)
		if err := camera.Connect(); err != nil {
			results = append(results, UpdateResult{
				CameraId: cameraId,
				Success:  false,
				Error:    fmt.Sprintf("Failed to connect: %v", err),
			})
			continue
		}

		configName := payload.ConfigName
		if configName == "" {
			currentConfig, err := lib.GetVideoEncoderConfiguration(camera, payload.ConfigToken)
			if err != nil {
				results = append(results, UpdateResult{
					CameraId: cameraId,
					Success:  false,
					Error:    fmt.Sprintf("Failed to get config: %v", err),
				})
				continue
			}
			configName = currentConfig.Name
		}

		// Apply the configuration changes		fmt.Printf("Applying config to camera %s (IP: %s)...\n", cameraId, targetCamera.IP)
		err := lib.SetVideoEncoderConfiguration(
			camera,
			payload.ConfigToken,
			configName,
			payload.Width,
			payload.Height,
			payload.FrameRate,
			payload.BitRate,
			payload.GopLength,
			payload.H264Profile,
		)

		if err != nil {
			fmt.Printf("Failed to set config for camera %s: %v\n", cameraId, err)
			results = append(results, UpdateResult{
				CameraId:         cameraId,
				Success:          false,
				Error:            fmt.Sprintf("Failed to set config: %v", err),
				ValidationResult: nil,
			})
			continue
		}

		fmt.Printf("Successfully applied config to camera %s\n", cameraId)
		// Get the stream URL for validation
		streamURI, err := lib.GetStreamURI(camera, payload.ProfileToken)
		if err != nil {
			results = append(results, UpdateResult{
				CameraId:         cameraId,
				Success:          true, // Config was applied, but validation couldn't be performed
				Error:            fmt.Sprintf("Config applied but couldn't validate: %v", err),
				ValidationResult: nil,
			})
			continue
		}
		// Add authentication credentials to the stream URL for camera 7
		if cameraId == "7" {
			parsedURL, err := url.Parse(streamURI)
			if err == nil {
				// Add username and password to the URL
				parsedURL.User = url.UserPassword(targetCamera.Username, targetCamera.Password)
				streamURI = parsedURL.String()
				fmt.Printf("Added authentication to stream URL for camera 7: %s\n", streamURI)
			}
		}
		// Validate the configuration
		var errStr string
		validationResult, err := validator.ValidateVideoConfig(streamURI, expectedConfig)
		if err != nil {
			errStr = err.Error()
		}
		results = append(results, UpdateResult{
			CameraId:         cameraId,
			Success:          err == nil && validationResult.IsValid,
			ValidationResult: validationResult,
			Error:            errStr,
		})
	}

	successCount := 0
	validatedCount := 0
	for _, result := range results {
		if result.Success {
			successCount++
			if result.ValidationResult != nil && result.ValidationResult.IsValid {
				validatedCount++
			}
		}
	}

	response := struct {
		Message string         `json:"message"`
		Results []UpdateResult `json:"results"`
	}{
		Message: fmt.Sprintf(
			"Updated %d of %d cameras (%d validated successfully)",
			successCount,
			len(payload.CameraIds),
			validatedCount,
		),
		Results: results,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Handler to change resolution with automatic profile/config detection
func changeResolutionSimple(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\nchangeResolutionSimple: Starting resolution change request")

	var payload struct {
		CameraIds   []string `json:"cameraIds"`
		Width       int      `json:"width"`
		Height      int      `json:"height"`
		FrameRate   int      `json:"frameRate"`
		BitRate     int      `json:"bitRate"`
		GopLength   int      `json:"gopLength"`
		H264Profile string   `json:"h264Profile"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		fmt.Printf("changeResolutionSimple: Failed to decode request payload: %v\n", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	fmt.Printf("changeResolutionSimple: Request to set resolution %dx%d for %d cameras\n",
		payload.Width, payload.Height, len(payload.CameraIds))

	type UpdateResult struct {
		CameraId         string                      `json:"cameraId"`
		Success          bool                        `json:"success"`
		Error            string                      `json:"error,omitempty"`
		ValidationResult *validator.ValidationResult `json:"validationResult,omitempty"`
	}
	results := make([]UpdateResult, 0)

	expectedConfig := validator.VideoConfig{
		Width:       payload.Width,
		Height:      payload.Height,
		FrameRate:   payload.FrameRate,
		BitRate:     payload.BitRate,
		GopLength:   payload.GopLength,
		H264Profile: payload.H264Profile,
	}

	for _, cameraId := range payload.CameraIds {
		var targetCamera *Camera
		for _, cam := range cameras {
			if cam.ID == cameraId {
				targetCamera = &cam
				break
			}
		}
		if targetCamera == nil {
			results = append(results, UpdateResult{
				CameraId: cameraId,
				Success:  false,
				Error:    "Camera not found",
			})
			continue
		}

		if targetCamera.IsFake {
			results = append(results, UpdateResult{
				CameraId: cameraId,
				Success:  true,
				Error:    "Skipped: This is a simulated camera",
			})
			continue
		}

		camera := lib.NewCamera(targetCamera.IP, 80, targetCamera.Username, targetCamera.Password)
		if err := camera.Connect(); err != nil {
			results = append(results, UpdateResult{
				CameraId: cameraId,
				Success:  false,
				Error:    fmt.Sprintf("Failed to connect: %v", err),
			})
			continue
		}
		// Get active profile and config
		profileToken, configToken, err := lib.GetActiveProfile(camera)
		if err != nil {
			fmt.Printf("Failed to get active profile/config for camera %s: %v\n", cameraId, err)
			results = append(results, UpdateResult{
				CameraId: cameraId,
				Success:  false,
				Error:    fmt.Sprintf("Failed to get active profile/config: %v", err),
			})
			continue
		}
		fmt.Printf("Got active profile=%s, config=%s for camera %s\n", profileToken, configToken, cameraId)

		// Get available video encoder options
		options, err := lib.GetVideoEncoderOptions(camera, configToken, profileToken)
		if err != nil {
			fmt.Printf("Failed to get video encoder options for camera %s: %v\n", cameraId, err)
		} else {
			h264 := options.Body.GetVideoEncoderConfigurationOptionsResponse.Options.H264
			fmt.Printf("Available resolutions for camera %s:\n", cameraId)
			for _, res := range h264.ResolutionsAvailable {
				fmt.Printf("- %dx%d\n", res.Width, res.Height)
			}
		}

		// Get the current config to preserve the name
		currentConfig, err := lib.GetVideoEncoderConfiguration(camera, configToken)
		if err != nil {
			fmt.Printf("Failed to get current config for camera %s: %v\n", cameraId, err)
			results = append(results, UpdateResult{
				CameraId: cameraId,
				Success:  false,
				Error:    fmt.Sprintf("Failed to get current config: %v", err),
			})
			continue
		}
		fmt.Printf("Current config for camera %s: %dx%d @ %dfps\n",
			cameraId, currentConfig.Width, currentConfig.Height, currentConfig.FrameRate)

		// Apply the configuration changes
		var applyErr error
		if cameraId == "7" {
			// Special handling for camera 7 - don't set GOP length
			fmt.Printf("Special handling for camera 7 - not setting GOP length\n")

			// Attempt #1: Try using the existing GOP length
			applyErr = lib.SetVideoEncoderConfiguration(
				camera,
				configToken,
				currentConfig.Name,
				payload.Width,
				payload.Height,
				payload.FrameRate,
				payload.BitRate,
				currentConfig.GovLength, // Use existing GOP length
				payload.H264Profile,
			)

			if applyErr != nil {
				fmt.Printf("Attempt #1 failed for camera 7: %v\n", applyErr)

				// Attempt #2: Try setting only resolution and framerate
				fmt.Printf("Trying second approach with minimal parameters...\n")
				applyErr = lib.SetVideoEncoderConfiguration(
					camera,
					configToken,
					currentConfig.Name,
					payload.Width,
					payload.Height,
					payload.FrameRate,
					currentConfig.BitRate,     // Use existing bitrate
					currentConfig.GovLength,   // Use existing GOP length
					currentConfig.H264Profile, // Use existing profile
				)

				if applyErr != nil {
					fmt.Printf("Attempt #2 failed for camera 7: %v\n", applyErr)

					// Attempt #3: Try getting "Profile_2" instead
					profileDetails, profileErr := lib.GetProfileDetails(camera)
					if profileErr == nil {
						fmt.Printf("Got profile details with %d profiles\n",
							len(profileDetails.Body.GetProfilesResponse.Profiles))

						// Try to find a profile containing "Profile_2"
						for _, profile := range profileDetails.Body.GetProfilesResponse.Profiles {
							if strings.Contains(profile.Token, "Profile_2") {
								fmt.Printf("Found Profile_2, trying with token: %s\n", profile.Token)

								// Get the active profile's config token
								alternateProfile, alternateConfig, err := lib.GetActiveProfile(camera)
								if err != nil {
									fmt.Printf("Failed to get config token: %v\n", err)
									continue
								}

								fmt.Printf("Using profile token %s with config token %s\n",
									alternateProfile, alternateConfig)

								// Try setting configuration on the alternate profile
								applyErr = lib.SetVideoEncoderConfiguration(
									camera,
									alternateConfig,
									currentConfig.Name,
									payload.Width,
									payload.Height,
									payload.FrameRate,
									currentConfig.BitRate,     // Use existing bitrate
									currentConfig.GovLength,   // Use existing GOP length
									currentConfig.H264Profile, // Use existing profile
								)

								if applyErr == nil {
									fmt.Printf("Successfully set configuration on alternate profile!\n")
									configToken = alternateConfig   // Update token for validation
									profileToken = alternateProfile // Update for stream URL
								} else {
									fmt.Printf("Alternate profile attempt failed: %v\n", applyErr)
								}
								break
							}
						}
					}
				}
			}
		} else {
			// Regular handling for other cameras
			applyErr = lib.SetVideoEncoderConfiguration(
				camera,
				configToken,
				currentConfig.Name,
				payload.Width,
				payload.Height,
				payload.FrameRate,
				payload.BitRate,
				payload.GopLength,
				payload.H264Profile,
			)
		}

		if applyErr != nil {
			fmt.Printf("Failed to set config for camera %s: %v\n", cameraId, applyErr)
			results = append(results, UpdateResult{
				CameraId:         cameraId,
				Success:          false,
				Error:            fmt.Sprintf("Failed to set config: %v", applyErr),
				ValidationResult: nil,
			})
			continue
		}

		fmt.Printf("Successfully applied config to camera %s\n", cameraId)

		// Verify that the configuration was actually applied by reading it back
		if cameraId == "7" {
			time.Sleep(5 * time.Second)
			verifyConfig, verifyErr := lib.GetVideoEncoderConfiguration(camera, configToken)
			if verifyErr != nil {
				fmt.Printf("Warning: Could not verify camera 7 configuration: %v\n", verifyErr)
			} else {
				fmt.Printf("Camera 7 post-change config: %dx%d @ %dfps (GOV length: %d)\n",
					verifyConfig.Width, verifyConfig.Height, verifyConfig.FrameRate, verifyConfig.GovLength)
				if verifyConfig.Width != payload.Width || verifyConfig.Height != payload.Height {
					fmt.Printf("Warning: Camera 7 config does not match requested resolution!\n")
				}
			}
		}

		// Add special handling for camera 7 - longer delay to ensure config is applied
		if cameraId == "7" {
			fmt.Printf("Adding extended delay for camera 7 to ensure configuration is applied...\n")
			time.Sleep(10 * time.Second)
		}

		// Get the stream URL for validation
		streamURI, err := lib.GetStreamURI(camera, profileToken)
		if err != nil {
			results = append(results, UpdateResult{
				CameraId:         cameraId,
				Success:          true,
				Error:            fmt.Sprintf("Config applied but couldn't validate: %v", err),
				ValidationResult: nil,
			})
			continue
		}

		// Add authentication credentials to the stream URL for camera 7
		if cameraId == "7" {
			parsedURL, err := url.Parse(streamURI)
			if err == nil {
				// Add username and password to the URL
				parsedURL.User = url.UserPassword(targetCamera.Username, targetCamera.Password)
				streamURI = parsedURL.String()
				fmt.Printf("Added authentication to stream URL for camera 7: %s\n", streamURI)
			}
		}

		// Validate the configuration
		var errStr string
		validationResult, err := validator.ValidateVideoConfig(streamURI, expectedConfig)
		if err != nil {
			errStr = err.Error()
		}
		results = append(results, UpdateResult{
			CameraId:         cameraId,
			Success:          err == nil && validationResult.IsValid,
			ValidationResult: validationResult,
			Error:            errStr,
		})
	}

	successCount := 0
	validatedCount := 0
	for _, result := range results {
		if result.Success {
			successCount++
			if result.ValidationResult != nil && result.ValidationResult.IsValid {
				validatedCount++
			}
		}
	}

	response := struct {
		Message string         `json:"message"`
		Results []UpdateResult `json:"results"`
	}{
		Message: fmt.Sprintf(
			"Updated %d of %d cameras (%d validated successfully)",
			successCount,
			len(payload.CameraIds),
			validatedCount,
		),
		Results: results,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Handler to get single configuration details
func getConfigDetails(w http.ResponseWriter, r *http.Request) {
	fmt.Println("getConfigDetails: Starting request for configuration details")

	camera, err := getCamera()
	if err != nil {
		fmt.Printf("getConfigDetails: Failed to get camera connection: %v\n", err)
		http.Error(w, fmt.Sprintf("Failed to connect to camera: %v", err), http.StatusInternalServerError)
		return
	}

	configToken := r.URL.Query().Get("configToken")
	if configToken == "" {
		http.Error(w, "Missing configToken parameter", http.StatusBadRequest)
		return
	}

	fmt.Printf("getConfigDetails: Using configToken=%s\n", configToken)

	configDetails, err := lib.GetVideoEncoderConfiguration(camera, configToken)
	if err != nil {
		fmt.Printf("getConfigDetails: Failed to get configuration details: %v\n", err)
		http.Error(w, fmt.Sprintf("Failed to get configuration details: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Println("getConfigDetails: Successfully retrieved configuration details")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(configDetails)
}

// Handler to get device information
func getDeviceInfo(w http.ResponseWriter, r *http.Request) {
	fmt.Println("getDeviceInfo: Starting request for device information")

	camera, err := getCamera()
	if err != nil {
		fmt.Printf("getDeviceInfo: Failed to get camera connection: %v\n", err)
		http.Error(w, fmt.Sprintf("Failed to connect to camera: %v", err), http.StatusInternalServerError)
		return
	}

	// Use the device library to get device information
	deviceInfoRequest := device.GetDeviceInformation{}
	deviceInfoResponse, err := camera.Device.CallMethod(deviceInfoRequest)
	if err != nil {
		fmt.Printf("getDeviceInfo: Failed to get device information: %v\n", err)
		http.Error(w, fmt.Sprintf("Failed to get device information: %v", err), http.StatusInternalServerError)
		return
	}

	// Read the response body
	rawDeviceInfoXML, err := io.ReadAll(deviceInfoResponse.Body)
	if err != nil {
		fmt.Printf("getDeviceInfo: Failed to read device information response: %v\n", err)
		http.Error(w, fmt.Sprintf("Failed to read device information response: %v", err), http.StatusInternalServerError)
		return
	}

	// Parse the XML response
	var deviceInfo struct {
		XMLName xml.Name `xml:"Envelope"`
		Body    struct {
			GetDeviceInformationResponse struct {
				Manufacturer    string `xml:"Manufacturer"`
				Model           string `xml:"Model"`
				FirmwareVersion string `xml:"FirmwareVersion"`
				SerialNumber    string `xml:"SerialNumber"`
				HardwareId      string `xml:"HardwareId"`
			} `xml:"GetDeviceInformationResponse"`
		} `xml:"Body"`
	}

	if err := xml.Unmarshal(rawDeviceInfoXML, &deviceInfo); err != nil {
		fmt.Printf("getDeviceInfo: Failed to parse device information XML: %v\n", err)
		http.Error(w, fmt.Sprintf("Failed to parse device information: %v", err), http.StatusInternalServerError)
		return
	}

	// Create a response object with the device information
	response := map[string]interface{}{
		"manufacturer":    deviceInfo.Body.GetDeviceInformationResponse.Manufacturer,
		"model":           deviceInfo.Body.GetDeviceInformationResponse.Model,
		"firmwareVersion": deviceInfo.Body.GetDeviceInformationResponse.FirmwareVersion,
		"serialNumber":    deviceInfo.Body.GetDeviceInformationResponse.SerialNumber,
		"hardwareId":      deviceInfo.Body.GetDeviceInformationResponse.HardwareId,
	}

	fmt.Println("getDeviceInfo: Successfully retrieved device information")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Handler to get camera details using active profile and config
func getCameraDetailsSimple(w http.ResponseWriter, r *http.Request) {
	cameraId := r.URL.Query().Get("cameraId")
	if cameraId == "" {
		http.Error(w, "Missing cameraId parameter", http.StatusBadRequest)
		return
	}

	// Find the requested camera
	var targetCamera *Camera
	for _, cam := range cameras {
		if cam.ID == cameraId {
			targetCamera = &cam
			break
		}
	}

	if targetCamera == nil {
		http.Error(w, "Camera not found", http.StatusNotFound)
		return
	}

	if targetCamera.IsFake {
		http.Error(w, "Cannot get details for simulated camera", http.StatusBadRequest)
		return
	}

	camera := lib.NewCamera(targetCamera.IP, 80, targetCamera.Username, targetCamera.Password)
	if err := camera.Connect(); err != nil {
		http.Error(w, fmt.Sprintf("Failed to connect to camera: %v", err), http.StatusInternalServerError)
		return
	}

	// Get active profile and config
	profileToken, configToken, err := lib.GetActiveProfile(camera)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get active profile/config: %v", err), http.StatusInternalServerError)
		return
	}

	// Get stream URL
	streamURI, err := lib.GetStreamURI(camera, profileToken)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get stream URI: %v", err), http.StatusInternalServerError)
		return
	}

	// Get config details
	configDetails, err := lib.GetVideoEncoderConfiguration(camera, configToken)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get config details: %v", err), http.StatusInternalServerError)
		return
	}

	// Get device info
	deviceInfoRequest := device.GetDeviceInformation{}
	deviceInfoResponse, err := camera.Device.CallMethod(deviceInfoRequest)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get device info: %v", err), http.StatusInternalServerError)
		return
	}

	// Read and parse device info response
	rawDeviceInfoXML, err := io.ReadAll(deviceInfoResponse.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read device info: %v", err), http.StatusInternalServerError)
		return
	}

	var deviceInfoXML struct {
		Body struct {
			GetDeviceInformationResponse struct {
				Manufacturer    string `xml:"Manufacturer"`
				Model           string `xml:"Model"`
				FirmwareVersion string `xml:"FirmwareVersion"`
				SerialNumber    string `xml:"SerialNumber"`
				HardwareId      string `xml:"HardwareId"`
			} `xml:"GetDeviceInformationResponse"`
		} `xml:"Body"`
	}

	if err := xml.Unmarshal(rawDeviceInfoXML, &deviceInfoXML); err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse device info: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"streamUrl":  streamURI,
		"config":     configDetails,
		"deviceInfo": deviceInfoXML.Body.GetDeviceInformationResponse,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
		return
	}
}
