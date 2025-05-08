package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/use-go/onvif/device"
)

// Configuration for the application
type Config struct {
	ServerPort    int    `json:"server_port"`
	CameraIP      string `json:"camera_ip"`
	CameraPort    int    `json:"camera_port"`
	CameraUser    string `json:"camera_user"`
	CameraPass    string `json:"camera_pass"`
	AllowedOrigin string `json:"allowed_origin"`
}

// Global variables
var (
	globalCamera *Camera
	cameraLock   sync.Mutex
	appConfig    Config
)

func main() {
	// Load configuration
	loadConfig()

	// Initialize logging
	logFile, err := os.OpenFile("onvif_pi.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Failed to open log file: %v. Logs will only be printed to console.", err)
	} else {
		defer logFile.Close()
		log.SetOutput(io.MultiWriter(os.Stdout, logFile))
	}

	log.Println("Starting ONVIF Pi proxy server...")
	log.Printf("Server configuration: port=%d, camera=%s:%d",
		appConfig.ServerPort, appConfig.CameraIP, appConfig.CameraPort)

	// Initialize camera connection
	initCamera()

	// Set up router
	r := mux.NewRouter()

	// API Routes
	apiRouter := r.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/camera/info", getCameraInfo).Methods("GET")
	apiRouter.HandleFunc("/camera/resolutions", getResolutions).Methods("GET")
	apiRouter.HandleFunc("/camera/change-resolution", changeResolution).Methods("POST")
	apiRouter.HandleFunc("/camera/stream-url", getStreamURL).Methods("GET")
	apiRouter.HandleFunc("/camera/config", getConfigDetails).Methods("GET")
	apiRouter.HandleFunc("/camera/device-info", getDeviceInfo).Methods("GET")
	apiRouter.HandleFunc("/system/status", getSystemStatus).Methods("GET")

	// Add CORS support
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{appConfig.AllowedOrigin}),
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type"}),
	)(r)

	// Setup graceful shutdown
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", appConfig.ServerPort),
		Handler: corsHandler,
	}

	// Channel to listen for interrupt signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		log.Printf("Server listening on port %d", appConfig.ServerPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-stop
	log.Println("Shutting down server...")

	// Cleanup before exit
	cleanup()
}

// Load configuration from config.json or use defaults
func loadConfig() {
	// Default configuration
	appConfig = Config{
		ServerPort:    8090,
		CameraIP:      "192.168.1.12",
		CameraPort:    80,
		CameraUser:    "admin",
		CameraPass:    "admin123",
		AllowedOrigin: "http://localhost:3000",
	}

	// Try to load configuration from file
	configFile, err := os.Open("config.json")
	if err != nil {
		log.Println("Could not open config file, using defaults")
		// Create default config file for user to modify
		saveConfig()
		return
	}
	defer configFile.Close()

	// Decode JSON
	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&appConfig)
	if err != nil {
		log.Printf("Error parsing config file: %v. Using defaults.", err)
	}
}

// Save current configuration to config.json
func saveConfig() {
	configFile, err := os.Create("config.json")
	if err != nil {
		log.Printf("Could not create config file: %v", err)
		return
	}
	defer configFile.Close()

	encoder := json.NewEncoder(configFile)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(appConfig); err != nil {
		log.Printf("Error writing config file: %v", err)
	}
}

// Cleanup resources before exit
func cleanup() {
	log.Println("Cleaning up resources...")

	// Reset camera connection
	cameraLock.Lock()
	globalCamera = nil
	cameraLock.Unlock()

	log.Println("Cleanup complete")
}

// Get or initialize camera connection
func getCamera() (*Camera, error) {
	cameraLock.Lock()
	defer cameraLock.Unlock()

	// If camera isn't connected yet, create a new connection
	if globalCamera == nil {
		log.Println("Creating new camera connection...")
		globalCamera = NewCamera(
			appConfig.CameraIP,
			appConfig.CameraPort,
			appConfig.CameraUser,
			appConfig.CameraPass,
		)

		if err := globalCamera.Connect(); err != nil {
			globalCamera = nil
			return nil, fmt.Errorf("failed to connect to camera: %v", err)
		}
		log.Println("Successfully connected to camera")
	}

	return globalCamera, nil
}

// Initialize camera at startup
func initCamera() {
	log.Println("Initializing camera connection...")
	_, err := getCamera()
	if err != nil {
		log.Printf("Warning: Failed to initialize camera connection: %v", err)
		log.Println("Will try to reconnect on first request")
	} else {
		log.Println("Camera connection initialized successfully")
	}
}

// Reset camera connection
func resetCamera() {
	cameraLock.Lock()
	defer cameraLock.Unlock()
	globalCamera = nil
	log.Println("Camera connection reset")
}

// API handlers
func getCameraInfo(w http.ResponseWriter, r *http.Request) {
	log.Println("API request: Get camera info")

	camera, err := getCamera()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to connect to camera: %v", err), http.StatusInternalServerError)
		return
	}

	configs, err := GetAllVideoEncoderConfigurations(camera)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get configurations: %v", err), http.StatusInternalServerError)
		return
	}

	profiles, err := GetAllProfiles(camera)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get profiles: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"configs":  configs,
		"profiles": profiles,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getResolutions(w http.ResponseWriter, r *http.Request) {
	log.Println("API request: Get resolutions")

	camera, err := getCamera()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to connect to camera: %v", err), http.StatusInternalServerError)
		return
	}

	configToken := r.URL.Query().Get("configToken")
	profileToken := r.URL.Query().Get("profileToken")

	log.Printf("Getting resolutions for configToken=%s, profileToken=%s", configToken, profileToken)

	options, err := GetVideoEncoderOptions(camera, configToken, profileToken)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get video encoder options: %v", err), http.StatusInternalServerError)
		return
	}

	h264Options := ParseH264Options(options)

	// Debug log the parsed options
	resolutionCount := len(h264Options.ResolutionOptions)
	log.Printf("Parsed %d resolutions from camera", resolutionCount)

	if resolutionCount > 0 {
		log.Printf("First resolution: %dx%d",
			h264Options.ResolutionOptions[0].Width,
			h264Options.ResolutionOptions[0].Height)
	} else {
		log.Printf("WARNING: No resolutions were parsed from camera response!")

		// If no resolutions were found, add default resolutions
		log.Printf("Adding default resolutions as fallback")
		h264Options.ResolutionOptions = []VideoResolution{
			{Width: 1920, Height: 1080},
			{Width: 1280, Height: 720},
			{Width: 640, Height: 480},
			{Width: 320, Height: 240},
		}
	}

	// Log the response that will be sent
	jsonBytes, _ := json.MarshalIndent(h264Options, "", "  ")
	log.Printf("Sending JSON response: %s", string(jsonBytes))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(h264Options)
}

func changeResolution(w http.ResponseWriter, r *http.Request) {
	log.Println("API request: Change resolution")

	var payload struct {
		ConfigToken  string `json:"configToken"`
		ProfileToken string `json:"profileToken"`
		Width        int    `json:"width"`
		Height       int    `json:"height"`
		FrameRate    int    `json:"frameRate"`
		BitRate      int    `json:"bitRate"`
		GopLength    int    `json:"gopLength"`
		H264Profile  string `json:"h264Profile"`
		ConfigName   string `json:"configName"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	log.Printf("Changing to %dx%d, frameRate=%d, bitRate=%d, gopLength=%d, profile=%s",
		payload.Width, payload.Height, payload.FrameRate, payload.BitRate, payload.GopLength, payload.H264Profile)

	camera, err := getCamera()
	if err != nil {
		http.Error(w, "Failed to connect to camera", http.StatusInternalServerError)
		return
	}

	// If configName is not provided, get the original name
	configName := payload.ConfigName
	if configName == "" {
		// Get the current configuration to preserve its name
		currentConfig, err := GetVideoEncoderConfiguration(camera, payload.ConfigToken)
		if err != nil {
			http.Error(w, "Failed to get current configuration", http.StatusInternalServerError)
			return
		}
		configName = currentConfig.Name
	}

	err = SetVideoEncoderConfiguration(
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
		http.Error(w, "Failed to change resolution", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "Resolution updated successfully"}`))
}

func getStreamURL(w http.ResponseWriter, r *http.Request) {
	log.Println("API request: Get stream URL")

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

	streamURI, err := GetStreamURI(camera, profileToken)
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

func getConfigDetails(w http.ResponseWriter, r *http.Request) {
	log.Println("API request: Get config details")

	camera, err := getCamera()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to connect to camera: %v", err), http.StatusInternalServerError)
		return
	}

	configToken := r.URL.Query().Get("configToken")
	if configToken == "" {
		http.Error(w, "Missing configToken parameter", http.StatusBadRequest)
		return
	}

	configDetails, err := GetVideoEncoderConfiguration(camera, configToken)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get configuration details: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(configDetails)
}

func getDeviceInfo(w http.ResponseWriter, r *http.Request) {
	log.Println("API request: Get device info")

	camera, err := getCamera()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to connect to camera: %v", err), http.StatusInternalServerError)
		return
	}

	// Use the device library to get device information
	deviceInfoRequest := device.GetDeviceInformation{}
	deviceInfoResponse, err := camera.Device.CallMethod(deviceInfoRequest)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get device information: %v", err), http.StatusInternalServerError)
		return
	}

	// Read the response body
	rawDeviceInfoXML, err := io.ReadAll(deviceInfoResponse.Body)
	if err != nil {
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getSystemStatus(w http.ResponseWriter, r *http.Request) {
	log.Println("API request: Get system status")

	// Check camera connection
	cameraStatus := "connected"
	camera, err := getCamera()
	if err != nil || camera == nil {
		cameraStatus = "disconnected"
	}

	// Get basic system information
	hostname, _ := os.Hostname()

	systemInfo := map[string]interface{}{
		"hostname":     hostname,
		"uptime":       getUptime(),
		"cameraStatus": cameraStatus,
		"serverTime":   time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(systemInfo)
}

// Get system uptime (simplified for cross-platform compatibility)
func getUptime() string {
	// This is a simplified version; on a real Pi you might use syscall to get actual uptime
	return "Unknown" // Placeholder
}
