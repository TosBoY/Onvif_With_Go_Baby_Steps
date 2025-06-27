package api

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"main_back/internal/camera"
	"main_back/internal/ffmpeg"
	"main_back/internal/loader"
	"main_back/internal/vlc"
	"main_back/pkg/models"

	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/cameras", HandleGetCameras).Methods("GET")
	r.HandleFunc("/cameras", HandleAddCamera).Methods("POST")
	r.HandleFunc("/cameras/{id}", HandleDeleteCamera).Methods("DELETE")
	r.HandleFunc("/load-cam-list", HandleLoadCamList).Methods("GET")
	r.HandleFunc("/check-single-cam/{id}", HandleCheckSingleCam).Methods("GET")
	r.HandleFunc("/config-single-cam/{id}", HandleConfigSingleCam).Methods("POST")
	r.HandleFunc("/validate-cam/{id}", HandleValidateCam).Methods("GET")
	r.HandleFunc("/cameras/import-csv", HandleImportCamerasCSV).Methods("POST")
	r.HandleFunc("/import-config-csv", HandleImportConfigCSV).Methods("POST")
	r.HandleFunc("/choose-cam-from-csv", HandleChooseCamFromCSV).Methods("POST")
	r.HandleFunc("/apply-config", HandleApplyConfig).Methods("POST")
	r.HandleFunc("/export-validation-csv", HandleExportValidationCSV).Methods("POST")
	r.HandleFunc("/vlc", HandleVLC).Methods("POST")

	// Debug: catch-all route to log unmatched requests
	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Unmatched request: %s %s", r.Method, r.URL.Path)
		http.NotFound(w, r)
	})

	// Debug: log all registered routes
	log.Println("Registered routes:")
	r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		if err == nil {
			methods, err := route.GetMethods()
			if err == nil {
				log.Printf("Route: %s %v", pathTemplate, methods)
			} else {
				log.Printf("Route: %s (no methods)", pathTemplate)
			}
		}
		return nil
	})
}

func HandleGetCameras(w http.ResponseWriter, r *http.Request) {
	// Get cameras from in-memory storage instead of CSV file
	cameras := camera.GetAllCameras()
	json.NewEncoder(w).Encode(cameras)
}

func HandleAddCamera(w http.ResponseWriter, r *http.Request) {
	log.Println("Received /cameras POST request to add a new camera")

	var input struct {
		IP       string `json:"ip"`
		Port     int    `json:"port"`
		URL      string `json:"url"`
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Printf("Error decoding add-camera request body: %v", err)
		http.Error(w, fmt.Sprintf("Failed to decode request body: %v", err), http.StatusBadRequest)
		return
	}

	// Validate the input
	if input.IP == "" {
		log.Println("Error: Missing IP address in add-camera request")
		http.Error(w, "IP address is required", http.StatusBadRequest)
		return
	}

	if input.Username == "" {
		log.Println("Error: Missing username in add-camera request")
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}
	// Password can be empty for some cameras, so we don't check for it
	log.Printf("Adding new camera with IP: %s, Port: %d, URL: %s, Username: %s",
		input.IP, input.Port, input.URL, input.Username)
	newID, err := camera.AddNewCamera(input.IP, input.Port, input.URL, input.Username, input.Password)
	if err != nil {
		log.Printf("Error adding new camera: %v", err)
		http.Error(w, fmt.Sprintf("Failed to add new camera: %v", err), http.StatusInternalServerError)
		return
	}
	// Return the new camera ID and details
	newCamera := models.Camera{
		ID:       newID,
		IP:       input.IP,
		Port:     input.Port,
		URL:      input.URL,
		Username: input.Username,
		Password: input.Password,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newCamera)
}

func HandleDeleteCamera(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cameraID := vars["id"]

	log.Printf("Received request to delete camera with ID: %s", cameraID)

	if cameraID == "" {
		log.Println("Error: Missing camera ID in delete request")
		http.Error(w, "Camera ID is required", http.StatusBadRequest)
		return
	}

	err := camera.RemoveCamera(cameraID)
	if err != nil {
		log.Printf("Error deleting camera: %v", err)
		http.Error(w, fmt.Sprintf("Failed to delete camera: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully deleted camera with ID: %s", cameraID)

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{
		"message": fmt.Sprintf("Camera %s successfully deleted", cameraID),
	}
	json.NewEncoder(w).Encode(response)
}

func HandleLoadCamList(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to load camera list from CSV")

	cameras, err := loader.LoadCameraList()
	if err != nil {
		log.Printf("Error loading cameras from CSV: %v", err)
		http.Error(w, fmt.Sprintf("Failed to load cameras from CSV: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("Loaded %d cameras from CSV", len(cameras))

	// Convert cameras to response format
	var results []map[string]interface{}
	for _, cam := range cameras {
		result := map[string]interface{}{
			"cameraId": cam.ID,
			"ip":       cam.IP,
			"port":     cam.Port,
			"username": cam.Username,
			"url":      cam.URL,
		}
		results = append(results, result)
	}

	response := map[string]interface{}{
		"cameras": results,
		"total":   len(cameras),
		"message": fmt.Sprintf("Successfully loaded %d cameras from CSV", len(cameras)),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func HandleCheckSingleCam(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cameraID := vars["id"]

	log.Printf("Received request to check camera ID: %s", cameraID)

	if cameraID == "" {
		log.Println("Error: Missing camera ID in request")
		http.Error(w, "Camera ID is required", http.StatusBadRequest)
		return
	}

	// Load cameras from CSV to find the requested camera
	cameras, err := loader.LoadCameraList()
	if err != nil {
		log.Printf("Error loading cameras from CSV: %v", err)
		http.Error(w, fmt.Sprintf("Failed to load cameras from CSV: %v", err), http.StatusInternalServerError)
		return
	}

	// Find the camera with the specified ID
	var targetCamera *models.Camera
	for _, cam := range cameras {
		if cam.ID == cameraID {
			targetCamera = &cam
			break
		}
	}

	if targetCamera == nil {
		log.Printf("Camera with ID %s not found in CSV", cameraID)
		http.Error(w, fmt.Sprintf("Camera with ID %s not found", cameraID), http.StatusNotFound)
		return
	}

	log.Printf("Checking camera %s (IP: %s:%d)", targetCamera.ID, targetCamera.IP, targetCamera.Port)

	// Prepare response structure
	result := map[string]interface{}{
		"cameraId": targetCamera.ID,
		"ip":       targetCamera.IP,
		"port":     targetCamera.Port,
		"username": targetCamera.Username,
		"status":   "unknown",
		"error":    "",
	}

	// Try to initialize the camera client
	client, err := camera.NewCameraClient(*targetCamera)
	if err != nil {
		log.Printf("Failed to initialize camera %s: %v", targetCamera.ID, err)

		// Ping the camera to determine if it's a network issue (offline) or ONVIF issue (error)
		log.Printf("Pinging camera %s at IP %s to determine connectivity", targetCamera.ID, targetCamera.IP)
		pingSuccess, pingOutput, pingErr := pingHost(targetCamera.IP, 5*time.Second)

		if pingSuccess {
			// Camera is reachable but ONVIF failed - this is an error
			result["status"] = "error"
			result["error"] = fmt.Sprintf("Camera is reachable but ONVIF initialization failed: %v", err)
			log.Printf("Camera %s is reachable via ping but ONVIF failed", targetCamera.ID)
		} else {
			// Camera is not reachable - this is offline
			result["status"] = "offline"
			if pingErr != nil {
				result["error"] = fmt.Sprintf("Camera not reachable: %v", pingErr)
			} else {
				result["error"] = fmt.Sprintf("Camera not reachable: %s", pingOutput)
			}
			log.Printf("Camera %s is not reachable via ping", targetCamera.ID)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
		return
	}

	// Try to get configuration
	log.Printf("Getting profiles and configs for camera %s (IP: %s:%d)", targetCamera.ID, client.Camera.IP, client.Camera.Port)
	profileTokens, configTokens, err := camera.GetProfilesAndConfigs(client)
	if err != nil {
		log.Printf("Failed to get camera profiles and configs for %s (IP: %s:%d): %v", targetCamera.ID, client.Camera.IP, client.Camera.Port, err)

		// Ping the camera to determine if it's a network issue (offline) or ONVIF issue (error)
		log.Printf("Pinging camera %s at IP %s to determine connectivity after profiles failure", targetCamera.ID, targetCamera.IP)
		pingSuccess, pingOutput, pingErr := pingHost(targetCamera.IP, 5*time.Second)

		if pingSuccess {
			// Camera is reachable but ONVIF profiles failed - this is an error
			result["status"] = "error"
			errorMsg := err.Error()
			if strings.Contains(errorMsg, "i/o timeout") || strings.Contains(errorMsg, "dial tcp") {
				result["error"] = "Camera reachable but ONVIF timeout: check ONVIF service"
			} else if strings.Contains(errorMsg, "connection refused") {
				result["error"] = "Connection refused: check ONVIF port and service"
			} else if strings.Contains(errorMsg, "no route to host") {
				result["error"] = "No route to host: check network connectivity"
			} else {
				result["error"] = fmt.Sprintf("Camera reachable but ONVIF profiles failed: %v", err)
			}
			log.Printf("Camera %s is reachable via ping but ONVIF profiles failed", targetCamera.ID)
		} else {
			// Camera is not reachable - this is offline
			result["status"] = "offline"
			if pingErr != nil {
				result["error"] = fmt.Sprintf("Camera not reachable: %v", pingErr)
			} else {
				result["error"] = fmt.Sprintf("Camera not reachable: %s", pingOutput)
			}
			log.Printf("Camera %s is not reachable via ping during profiles check", targetCamera.ID)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
		return
	}

	if len(profileTokens) == 0 {
		log.Printf("No profiles found for camera %s", targetCamera.ID)
		result["status"] = "error"
		result["error"] = "No profiles found"
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
		return
	}

	if len(configTokens) == 0 {
		log.Printf("No video encoder configuration found for camera %s", targetCamera.ID)
		result["status"] = "error"
		result["error"] = "No video encoder configuration found"
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
		return
	}

	// Use the first token found
	profileToken := profileTokens[0]
	configToken := configTokens[0]
	result["profileToken"] = profileToken
	result["configToken"] = configToken

	log.Printf("Using profile token %s and config token %s for camera %s", profileToken, configToken, targetCamera.ID)

	// Try to get current encoder config
	log.Printf("Getting current encoder config for camera %s", targetCamera.ID)
	currentConfig, err := camera.GetCurrentConfig(client, configToken)
	if err != nil {
		log.Printf("Failed to get current encoder config for %s: %v", targetCamera.ID, err)
		result["status"] = "partial"
		result["error"] = fmt.Sprintf("Failed to get current config: %v", err)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
		return
	}

	// Try to get available encoder options
	log.Printf("Getting available encoder options for camera %s", targetCamera.ID)
	encoderOptions, err := camera.GetCurrentEncoderOptions(client, profileToken, configToken)
	if err != nil {
		log.Printf("Failed to get encoder options for %s: %v", targetCamera.ID, err)
		// Still mark as online since we got the current config
		result["status"] = "partial"
		result["error"] = fmt.Sprintf("Failed to get encoder options: %v", err)
		result["currentConfig"] = map[string]interface{}{
			"resolution": map[string]int{
				"width":  currentConfig.Resolution.Width,
				"height": currentConfig.Resolution.Height,
			},
			"fps":      currentConfig.FPS,
			"bitrate":  currentConfig.Bitrate,
			"encoding": currentConfig.Encoding,
			"quality":  currentConfig.Quality,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
		return
	}

	// Success - camera is fully online and configured
	result["status"] = "online"
	result["currentConfig"] = map[string]interface{}{
		"resolution": map[string]int{
			"width":  currentConfig.Resolution.Width,
			"height": currentConfig.Resolution.Height,
		},
		"fps":      currentConfig.FPS,
		"bitrate":  currentConfig.Bitrate,
		"encoding": currentConfig.Encoding,
		"quality":  currentConfig.Quality,
	}

	// Prepare available resolutions
	availableResolutions := make([]map[string]int, len(encoderOptions.Resolutions))
	for i, resolution := range encoderOptions.Resolutions {
		availableResolutions[i] = map[string]int{
			"width":  resolution.Width,
			"height": resolution.Height,
		}
	}
	result["availableResolutions"] = availableResolutions

	// Add encoder options information
	result["encoderOptions"] = map[string]interface{}{
		"availableFPS":     encoderOptions.FPSOptions,
		"availableBitrate": encoderOptions.Bitrate,
		"availableQuality": encoderOptions.Quality,
	}

	log.Printf("Successfully checked camera %s - status: online", targetCamera.ID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func HandleApplyConfig(w http.ResponseWriter, r *http.Request) {
	log.Println("Received /apply-config request")
	var input struct {
		CameraID  string   `json:"cameraId"`  // For backward compatibility
		CameraIDs []string `json:"cameraIds"` // New field for multiple cameras
		Width     int      `json:"width"`
		Height    int      `json:"height"`
		FPS       int      `json:"fps"`
		Bitrate   int      `json:"bitrate"`
		Encoding  string   `json:"encoding"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Printf("Error decoding apply-config request body: %v", err)
		http.Error(w, fmt.Sprintf("Failed to decode request body: %v", err), http.StatusBadRequest)
		return
	}
	// Handle both legacy (single camera) and new (multiple cameras) format
	var cameraIDs []string
	if len(input.CameraIDs) > 0 {
		cameraIDs = input.CameraIDs
		log.Printf("Applying config for multiple cameras (%d): Width: %d, Height: %d, FPS: %d, Bitrate: %d, Encoding: %s",
			len(cameraIDs), input.Width, input.Height, input.FPS, input.Bitrate, input.Encoding)
	} else if input.CameraID != "" {
		cameraIDs = []string{input.CameraID}
		log.Printf("Applying config for single camera %s: Width: %d, Height: %d, FPS: %d, Bitrate: %d, Encoding: %s",
			input.CameraID, input.Width, input.Height, input.FPS, input.Bitrate, input.Encoding)
	} else {
		log.Println("Error: No camera IDs provided in request")
		http.Error(w, "No camera IDs provided", http.StatusBadRequest)
		return
	}

	// Create a target resolution object once
	targetResolution := models.Resolution{Width: input.Width, Height: input.Height}

	// Create result structures for tracking progress
	type CameraConfigResult struct {
		CameraID           string
		Success            bool
		Error              error
		AppliedConfig      map[string]interface{}
		ResolutionAdjusted bool
		ProfileToken       string
		StreamURL          string
	}

	results := make(map[string]CameraConfigResult)
	validationResults := make(map[string]interface{})

	log.Printf("===== PHASE 1: Applying configuration to all cameras =====")
	// PHASE 1: Apply configuration to all cameras first
	for _, cameraID := range cameraIDs {
		// Initialize result for this camera
		result := CameraConfigResult{
			CameraID: cameraID,
			Success:  false,
		}

		// Get the camera client
		client, err := camera.GetCameraClient(cameraID)
		if err != nil {
			log.Printf("Error getting camera client for %s: %v", cameraID, err)
			result.Error = fmt.Errorf("camera not found: %w", err)
			results[cameraID] = result
			continue
		}

		// Proceed with config application
		log.Printf("\n Getting profiles and configs for camera %s (IP: %s:%d)", cameraID, client.Camera.IP, client.Camera.Port)
		profileTokens, configTokens, err := camera.GetProfilesAndConfigs(client)
		if err != nil {
			log.Printf("Failed to get camera profiles and configs for %s (IP: %s:%d): %v", cameraID, client.Camera.IP, client.Camera.Port, err)
			// Add more specific error information for network issues
			errorMsg := err.Error()
			if strings.Contains(errorMsg, "i/o timeout") || strings.Contains(errorMsg, "dial tcp") {
				result.Error = fmt.Errorf("network timeout: camera at %s:%d is not responding. Please check: 1) Camera is powered on and connected to network, 2) IP address %s is correct, 3) Port %d is the correct ONVIF port, 4) Camera supports ONVIF protocol", client.Camera.IP, client.Camera.Port, client.Camera.IP, client.Camera.Port)
			} else if strings.Contains(errorMsg, "connection refused") {
				result.Error = fmt.Errorf("connection refused: camera at %s:%d refused connection. Please check: 1) Correct ONVIF port (common ports: 80, 8080, 554), 2) ONVIF service is enabled on camera, 3) Firewall settings", client.Camera.IP, client.Camera.Port)
			} else if strings.Contains(errorMsg, "no route to host") {
				result.Error = fmt.Errorf("no route to host: cannot reach camera at %s:%d. Please check: 1) Camera and server are on same network, 2) IP address is correct, 3) Network routing", client.Camera.IP, client.Camera.Port)
			} else {
				result.Error = fmt.Errorf("failed to get camera profiles and configs: %w", err)
			}
			results[cameraID] = result
			continue
		}

		if len(profileTokens) == 0 {
			log.Printf("No profiles found for camera %s", cameraID)
			result.Error = fmt.Errorf("no profiles found")
			results[cameraID] = result
			continue
		}

		if len(configTokens) == 0 {
			log.Printf("No video encoder configuration found for camera %s", cameraID)
			result.Error = fmt.Errorf("no video encoder configuration found")
			results[cameraID] = result
			continue
		}

		// Use the first token found
		profileToken := profileTokens[0]
		configToken := configTokens[0]
		result.ProfileToken = profileToken

		log.Printf("Using profile token %s and config token %s for camera %s", profileToken, configToken, cameraID)

		// Get current encoder config
		log.Printf("Getting current encoder config for camera %s", cameraID)
		currentConfig, err := camera.GetCurrentConfig(client, configToken)
		if err != nil {
			log.Printf("Failed to get current encoder config for %s: %v", cameraID, err)
			result.Error = fmt.Errorf("failed to get current encoder config: %w", err)
			results[cameraID] = result
			continue
		}
		// log.Printf("Current encoding for camera %s: %+v", cameraID, currentConfig.Encoding)

		// Get available encoder options
		log.Printf("Getting available encoder options for camera %s", cameraID)
		encoderOptions, err := camera.GetCurrentEncoderOptions(client, profileToken, configToken)
		if err != nil {
			log.Printf("Failed to get encoder options for %s: %v", cameraID, err)
			result.Error = fmt.Errorf("failed to get encoder options: %w", err)
			results[cameraID] = result
			continue
		}

		// Find closest matching resolution
		log.Printf("Finding closest matching resolution for camera %s", cameraID)
		closestResolution := camera.FindClosestResolution(targetResolution, encoderOptions.Resolutions)
		log.Printf("Closest resolution found for camera %s: %dx%d", cameraID, closestResolution.Width, closestResolution.Height)

		newConfig := models.EncoderConfig{
			Resolution: closestResolution,
			Quality:    currentConfig.Quality, // Keep the current quality
			FPS:        input.FPS,
			Bitrate:    input.Bitrate,
			Encoding:   input.Encoding,
		}
		log.Printf("Prepared new config for camera %s: %+v", cameraID, newConfig)

		// Set the new encoder config
		log.Printf("Setting new encoder config for camera %s", cameraID)
		if err := camera.SetEncoderConfig(client, configToken, currentConfig, newConfig); err != nil {
			log.Printf("Failed to set encoder config for %s: %v", cameraID, err)
			result.Error = fmt.Errorf("failed to set encoder config: %w", err)
			results[cameraID] = result
			continue
		}
		log.Printf("Successfully applied config for camera %s", cameraID)

		// Get stream URI for later validation
		streamURI, err := client.GetStreamURI(profileToken)
		if err != nil {
			log.Printf("Failed to get stream URI for %s: %v", cameraID, err)
			result.Error = fmt.Errorf("failed to get stream URI: %w", err)
			results[cameraID] = result
			continue
		}

		// The ONVIF GetStreamUri typically doesn't include credentials. Embed them for FFmpeg validation.
		parsedURI, err := url.Parse(streamURI)
		if err != nil {
			log.Printf("Failed to parse stream URI for %s: %v", cameraID, err)
			result.Error = fmt.Errorf("failed to parse stream URI: %w", err)
			results[cameraID] = result
			continue
		}

		// Construct the URL with embedded credentials
		fullStreamURL := fmt.Sprintf("%s://%s:%s@%s%s", parsedURI.Scheme, client.Camera.Username, client.Camera.Password, parsedURI.Host, parsedURI.RequestURI())
		// Mark this camera as successfully configured
		result.Success = true
		result.AppliedConfig = map[string]interface{}{
			"resolution": map[string]int{
				"width":  closestResolution.Width,
				"height": closestResolution.Height,
			},
			"fps":      input.FPS,
			"bitrate":  input.Bitrate,
			"encoding": input.Encoding,
		}
		result.ResolutionAdjusted = input.Width != closestResolution.Width || input.Height != closestResolution.Height
		result.StreamURL = fullStreamURL

		results[cameraID] = result
		continue
	}

	// Helper function to find position of a camera ID in the original order
	findCameraPosition := func(cameraID string) int {
		for i, id := range cameraIDs {
			if id == cameraID {
				return i
			}
		}
		return -1
	}

	// PHASE 2: Validate all successfully configured cameras
	// Wait for camera configurations to stabilize
	time.Sleep(1 * time.Second)

	// Log that we're starting validation phase
	log.Printf("===== PHASE 2: Validating all cameras in original order =====")
	log.Printf("Original camera order: %v", cameraIDs)

	// Iterate through cameras in the SAME order as the configuration phase
	for _, cameraID := range cameraIDs {
		result, exists := results[cameraID]
		if !exists {
			log.Printf("No configuration result found for camera %s, skipping validation", cameraID)
			continue
		}
		if !result.Success {
			// Skip validation for failed configs - add validation results for failed configuration cameras
			position := findCameraPosition(cameraID)
			log.Printf("Adding failed validation results for camera %s (position %d in original order)", cameraID, position)
			errorMessage := "Configuration failed - camera not reachable"
			if result.Error != nil {
				errorMessage = result.Error.Error()
			}
			validationResults[cameraID] = map[string]interface{}{
				"isValid":          false,
				"expectedWidth":    input.Width,
				"expectedHeight":   input.Height,
				"expectedFPS":      input.FPS,
				"expectedBitrate":  input.Bitrate,
				"expectedEncoding": input.Encoding,
				"actualWidth":      0,
				"actualHeight":     0,
				"actualFPS":        0.0,
				"actualBitrate":    0,
				"error":            errorMessage,
			}
			continue
		}
		// Validate the stream using FFmpeg CGO
		position := findCameraPosition(cameraID)
		log.Printf("Starting FFmpeg validation for camera %s (position %d in original order)", cameraID, position)
		validationResult, validationErr := ffmpeg.ValidateStream(result.StreamURL, input.Width, input.Height, input.FPS, input.Bitrate, input.Encoding)
		if validationErr != nil {
			log.Printf("FFmpeg validation failed for camera %s: %v", cameraID, validationErr)
			validationResults[cameraID] = map[string]interface{}{
				"isValid":          false,
				"error":            validationErr.Error(),
				"expectedWidth":    input.Width,
				"expectedHeight":   input.Height,
				"expectedFPS":      input.FPS,
				"expectedBitrate":  input.Bitrate,
				"expectedEncoding": input.Encoding,
			}
		} else {
			// Determine validation status based on business rules
			// Resolution mismatch = failure, FPS/bitrate mismatch = warning
			resolutionMatches := validationResult.ActualWidth > 0 && validationResult.ActualHeight > 0 &&
				validationResult.ActualWidth == validationResult.ExpectedWidth &&
				validationResult.ActualHeight == validationResult.ExpectedHeight

			fpsMatches := validationResult.ActualFPS > 0 &&
				int(validationResult.ActualFPS+0.5) == validationResult.ExpectedFPS

			bitrateMatches := true // Default to true if no expected bitrate
			if validationResult.ExpectedBitrate > 0 && validationResult.ActualBitrate > 0 {
				tolerance := float64(validationResult.ExpectedBitrate) * 0.1
				diff := float64(validationResult.ActualBitrate - validationResult.ExpectedBitrate)
				if diff < 0 {
					diff = -diff
				}
				bitrateMatches = diff <= tolerance
			}

			// Override validation result: resolution mismatch = failure, others = warning
			overrideIsValid := resolutionMatches // Only consider valid if resolution matches

			// Create a map from the validation result
			validationMap := map[string]interface{}{
				"isValid":          overrideIsValid, // Use our override logic
				"actualWidth":      validationResult.ActualWidth,
				"actualHeight":     validationResult.ActualHeight,
				"actualFPS":        validationResult.ActualFPS,
				"actualBitrate":    validationResult.ActualBitrate,
				"actualEncoding":   validationResult.ActualEncoding,
				"expectedWidth":    validationResult.ExpectedWidth,
				"expectedHeight":   validationResult.ExpectedHeight,
				"expectedFPS":      validationResult.ExpectedFPS,
				"expectedBitrate":  validationResult.ExpectedBitrate,
				"expectedEncoding": input.Encoding,
			}

			// Build warning/error messages
			var messages []string
			if !resolutionMatches {
				if validationResult.ActualWidth > 0 && validationResult.ActualHeight > 0 {
					messages = append(messages, fmt.Sprintf("RESOLUTION MISMATCH: got %dx%d, expected %dx%d",
						validationResult.ActualWidth, validationResult.ActualHeight,
						validationResult.ExpectedWidth, validationResult.ExpectedHeight))
				} else {
					messages = append(messages, "RESOLUTION VALIDATION FAILED: unable to detect actual resolution")
				}
			}

			if !fpsMatches && validationResult.ActualFPS > 0 {
				messages = append(messages, fmt.Sprintf("FPS DIFFERENCE (warning): got %.2f fps, expected %d fps",
					validationResult.ActualFPS, validationResult.ExpectedFPS))
			}

			if !bitrateMatches && validationResult.ExpectedBitrate > 0 && validationResult.ActualBitrate > 0 {
				messages = append(messages, fmt.Sprintf("BITRATE DIFFERENCE (warning): got %d kbps, expected %d kbps",
					validationResult.ActualBitrate, validationResult.ExpectedBitrate))
			}

			// Check for encoding mismatches
			encodingMatches := true
			if input.Encoding != "" && validationResult.ActualEncoding != "" {
				encodingMatches = strings.EqualFold(validationResult.ActualEncoding, input.Encoding)
				if !encodingMatches {
					messages = append(messages, fmt.Sprintf("ENCODING DIFFERENCE (warning): got %s, expected %s",
						validationResult.ActualEncoding, input.Encoding))
				}
			}

			// Set error/warning message
			if len(messages) > 0 {
				validationMap["error"] = strings.Join(messages, "; ")
			} else if validationResult.Error != "" {
				validationMap["error"] = validationResult.Error
			}

			log.Printf("FFmpeg validation completed for camera %s: valid=%v", cameraID, validationResult.IsValid)
			validationResults[cameraID] = validationMap
		}
	}
	log.Printf("===== PHASE 2 COMPLETED =====")
	// Log summary of configuration results
	successCount := 0
	failureCount := 0

	for _, result := range results {
		if result.Success {
			successCount++
		} else {
			failureCount++
		}
	}

	log.Printf("Configuration Summary: %d successful, %d failed", successCount, failureCount)
	if failureCount > 0 {
		log.Printf("Failed cameras:")
		for cameraID, result := range results {
			if !result.Success {
				log.Printf("  - Camera %s: %v", cameraID, result.Error)
			}
		}
	}
	log.Printf("\n")

	// Collect configuration errors separately
	configurationErrors := make([]map[string]interface{}, 0)
	for _, cameraID := range cameraIDs {
		if result, exists := results[cameraID]; exists && !result.Success {
			errorMessage := "Configuration failed - camera not reachable"
			if result.Error != nil {
				errorMessage = result.Error.Error()
			}
			configurationErrors = append(configurationErrors, map[string]interface{}{
				"cameraId": cameraID,
				"error":    errorMessage,
			})
		}
	}

	// Prepare the final response
	finalResponse := map[string]interface{}{
		"status": "configuration applied",
		"originalRequest": map[string]interface{}{
			"resolution": map[string]int{
				"width":  input.Width,
				"height": input.Height,
			},
			"fps":      input.FPS,
			"bitrate":  input.Bitrate,
			"encoding": input.Encoding,
		},
		"results":             make(map[string]interface{}),
		"configurationErrors": configurationErrors,
	}

	// Add individual camera results
	for cameraID, result := range results {
		cameraResult := map[string]interface{}{
			"success": result.Success,
		}

		if result.Success {
			cameraResult["appliedConfig"] = result.AppliedConfig
			cameraResult["resolutionAdjusted"] = result.ResolutionAdjusted

			if validation, ok := validationResults[cameraID]; ok {
				cameraResult["validation"] = validation
			}
		} else if result.Error != nil {
			cameraResult["error"] = result.Error.Error()
		}

		finalResponse["results"].(map[string]interface{})[cameraID] = cameraResult
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(finalResponse)
}

func HandleVLC(w http.ResponseWriter, r *http.Request) {
	log.Println("Received /vlc request")

	var input struct {
		CameraID string `json:"cameraId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Printf("Error decoding VLC request body: %v", err)
		http.Error(w, fmt.Sprintf("Failed to decode request body: %v", err), http.StatusBadRequest)
		return
	}

	log.Printf("Launching VLC for camera ID: %s", input.CameraID)

	// Load cameras from CSV to find the requested camera
	cameras, err := loader.LoadCameraList()
	if err != nil {
		log.Printf("Error loading cameras from CSV: %v", err)
		http.Error(w, fmt.Sprintf("Failed to load cameras from CSV: %v", err), http.StatusInternalServerError)
		return
	}

	// Find the camera with the specified ID
	var targetCamera *models.Camera
	for _, cam := range cameras {
		if cam.ID == input.CameraID {
			targetCamera = &cam
			break
		}
	}

	if targetCamera == nil {
		log.Printf("Camera with ID %s not found in CSV", input.CameraID)
		http.Error(w, fmt.Sprintf("Camera with ID %s not found", input.CameraID), http.StatusNotFound)
		return
	}

	log.Printf("Found camera %s (IP: %s:%d)", targetCamera.ID, targetCamera.IP, targetCamera.Port)

	// Create a camera client for this specific camera
	client, err := camera.NewCameraClient(*targetCamera)
	if err != nil {
		log.Printf("Error creating camera client for %s: %v", input.CameraID, err)
		http.Error(w, fmt.Sprintf("Failed to create camera client: %v", err), http.StatusInternalServerError)
		return
	}

	// Get active profile
	log.Printf("Getting active profile for camera %s", input.CameraID)
	profileTokens, _, err := camera.GetProfilesAndConfigs(client)
	if err != nil {
		log.Printf("Failed to get camera profiles for %s: %v", input.CameraID, err)
		http.Error(w, fmt.Sprintf("Failed to get camera profiles: %v", err), http.StatusInternalServerError)
		return
	}

	if len(profileTokens) == 0 {
		log.Printf("No profiles found for camera %s", input.CameraID)
		http.Error(w, "No profiles found", http.StatusInternalServerError)
		return
	}

	// Use the first profile token
	profileToken := profileTokens[0]
	log.Printf("Using profile token %s for camera %s", profileToken, input.CameraID)

	// Get stream URI for the profile
	streamURI, err := client.GetStreamURI(profileToken)
	if err != nil {
		log.Printf("Failed to get stream URI for camera %s: %v", input.CameraID, err)
		http.Error(w, fmt.Sprintf("Failed to get stream URI: %v", err), http.StatusInternalServerError)
		return
	}

	// Add credentials to the stream URL
	parsedURI, err := url.Parse(streamURI)
	if err != nil {
		log.Printf("Failed to parse stream URI for camera %s: %v", input.CameraID, err)
		http.Error(w, fmt.Sprintf("Failed to parse stream URI: %v", err), http.StatusInternalServerError)
		return
	}

	// Add username and password to the URL
	parsedURI.User = url.UserPassword(client.Camera.Username, client.Camera.Password)
	authenticatedStreamURI := parsedURI.String()

	log.Printf("Stream URI with auth: %s", authenticatedStreamURI)

	// Launch VLC with the stream
	message, err := vlc.LaunchVLCWithStream(authenticatedStreamURI)
	if err != nil {
		log.Printf("Failed to launch VLC for camera %s: %v", input.CameraID, err)
		http.Error(w, fmt.Sprintf("Failed to launch VLC: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("VLC launched successfully for camera %s: %s", input.CameraID, message)

	response := map[string]interface{}{
		"message":   message,
		"streamUrl": authenticatedStreamURI,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func HandleExportValidationCSV(w http.ResponseWriter, r *http.Request) {
	log.Println("Received /export-validation-csv request")

	var input struct {
		Validation          interface{}   `json:"validation"`
		ConfigurationErrors []interface{} `json:"configurationErrors"`
		CameraOrder         []string      `json:"cameraOrder"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Printf("Error decoding export-validation-csv request body: %v", err)
		http.Error(w, fmt.Sprintf("Failed to decode request body: %v", err), http.StatusBadRequest)
		return
	}

	if input.Validation == nil && len(input.ConfigurationErrors) == 0 {
		log.Println("Error: No data provided for export")
		http.Error(w, "Validation data or configuration errors are required", http.StatusBadRequest)
		return
	}

	// Convert validation data to map format
	validationMap := make(map[string]interface{})
	if input.Validation != nil {
		var err error
		validationMap, err = convertValidationToMap(input.Validation)
		if err != nil {
			log.Printf("Error converting validation data: %v", err)
			http.Error(w, fmt.Sprintf("Invalid validation data format: %v", err), http.StatusBadRequest)
			return
		}
	}
	// Get camera information from in-memory storage
	cameras := camera.GetAllCameras()
	// Process configuration errors into a map for easy lookup
	configErrorsMap := make(map[string]string)
	for _, errItem := range input.ConfigurationErrors {
		if errMap, ok := errItem.(map[string]interface{}); ok {
			if cameraID, hasID := errMap["cameraId"]; hasID {
				if errorMsg, hasError := errMap["error"]; hasError {
					if id, ok := cameraID.(string); ok {
						if msg, ok := errorMsg.(string); ok {
							configErrorsMap[id] = msg
						}
					}
				}
			}
		}
	}

	// Generate CSV content
	csvContent, err := generateValidationCSV(validationMap, configErrorsMap, input.CameraOrder, cameras)
	if err != nil {
		log.Printf("Error generating CSV: %v", err)
		http.Error(w, fmt.Sprintf("Failed to generate CSV: %v", err), http.StatusInternalServerError)
		return
	}

	// Set headers for CSV download
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=\"validation_results.csv\"")

	// Write CSV content
	w.Write([]byte(csvContent))
	log.Println("CSV export completed successfully")
}

func generateValidationCSV(validation map[string]interface{}, configErrors map[string]string, cameraOrder []string, cameras []models.Camera) (string, error) {
	var csvBuilder strings.Builder
	writer := csv.NewWriter(&csvBuilder)

	// Create a map of camera ID to camera info for quick lookup
	cameraMap := make(map[string]models.Camera)
	for _, camera := range cameras {
		cameraMap[camera.ID] = camera
	}
	// Write CSV header with IP column and notes
	header := []string{"cam_id", "cam_ip", "result", "reso_expected", "reso_actual", "fps_expected", "fps_actual", "encoding_expected", "encoding_actual", "notes"}
	if err := writer.Write(header); err != nil {
		return "", fmt.Errorf("failed to write CSV header: %v", err)
	}
	// Create a list of camera IDs to process in the right order
	cameraIDs := make([]string, 0)

	// Collect all IDs from validation data and config errors
	idMap := make(map[string]bool)

	// Add IDs from validation results
	for cameraID := range validation {
		idMap[cameraID] = true
	}

	// Add IDs from configuration errors
	for cameraID := range configErrors {
		idMap[cameraID] = true
	}

	// Convert to slice
	for cameraID := range idMap {
		cameraIDs = append(cameraIDs, cameraID)
	}

	// Always sort by cam_id numerically, regardless of whether cameraOrder was provided
	sort.Slice(cameraIDs, func(i, j int) bool {
		// Extract numeric part using regular expression
		idI := extractNumericID(cameraIDs[i])
		idJ := extractNumericID(cameraIDs[j])

		// If both have numeric IDs, compare them numerically
		if idI != -1 && idJ != -1 {
			return idI < idJ
		}

		// Fall back to string comparison if not both numeric
		return cameraIDs[i] < cameraIDs[j]
	})

	// Process each camera in order
	for _, cameraID := range cameraIDs {
		// Check if this is a configuration error
		if errorMsg, hasError := configErrors[cameraID]; hasError {
			// Get camera IP from camera map
			cameraIP := "Unknown"
			if camera, exists := cameraMap[cameraID]; exists {
				cameraIP = camera.IP
			}

			// Write CSV row for configuration error
			row := []string{cameraID, cameraIP, "CONFIG_ERROR", "", "", "", "", "", "", fmt.Sprintf("Configuration Error: %s", errorMsg)}
			if err := writer.Write(row); err != nil {
				return "", fmt.Errorf("failed to write CSV row for camera %s: %v", cameraID, err)
			}
			continue
		}

		// Process validation data if available
		validationData, hasValidation := validation[cameraID]
		if !hasValidation {
			continue
		}
		// Convert interface{} to map[string]interface{}
		validationMap, ok := validationData.(map[string]interface{})
		if !ok {
			log.Printf("Warning: Invalid validation data format for camera %s", cameraID)
			continue
		}

		// Get camera IP from camera map
		cameraIP := "Unknown"
		if camera, exists := cameraMap[cameraID]; exists {
			cameraIP = camera.IP
		} // Extract data with safe type conversion and determine result status
		result := "FAIL" // Default to fail
		var notes strings.Builder

		if isValid, exists := validationMap["isValid"]; exists {
			if valid, ok := isValid.(bool); ok && valid {
				// Camera is valid, but check if there are warnings (FPS/bitrate mismatches)

				// Check resolution match
				resolutionMatches := true
				if expectedWidth, hasExpWidth := validationMap["expectedWidth"]; hasExpWidth {
					if expectedHeight, hasExpHeight := validationMap["expectedHeight"]; hasExpHeight {
						if actualWidth, hasActWidth := validationMap["actualWidth"]; hasActWidth {
							if actualHeight, hasActHeight := validationMap["actualHeight"]; hasActHeight {
								if expW, ok1 := expectedWidth.(float64); ok1 {
									if expH, ok2 := expectedHeight.(float64); ok2 {
										if actW, ok3 := actualWidth.(float64); ok3 {
											if actH, ok4 := actualHeight.(float64); ok4 {
												resolutionMatches = (actW > 0 && actH > 0 &&
													int(actW) == int(expW) && int(actH) == int(expH))
											}
										}
									}
								}
							}
						}
					}
				}

				// Check FPS match
				fpsMatches := true
				if expectedFPS, hasExpFPS := validationMap["expectedFPS"]; hasExpFPS {
					if actualFPS, hasActFPS := validationMap["actualFPS"]; hasActFPS {
						if expFPS, ok1 := expectedFPS.(float64); ok1 {
							if actFPS, ok2 := actualFPS.(float64); ok2 {
								if actFPS > 0 {
									fpsMatches = (int(actFPS+0.5) == int(expFPS))
								}
							}
						}
					}
				}

				// Check bitrate match
				bitrateMatches := true
				if expectedBitrate, hasExpBitrate := validationMap["expectedBitrate"]; hasExpBitrate {
					if actualBitrate, hasActBitrate := validationMap["actualBitrate"]; hasActBitrate {
						if expBitrate, ok1 := expectedBitrate.(float64); ok1 {
							if actBitrate, ok2 := actualBitrate.(float64); ok2 {
								if expBitrate > 0 && actBitrate > 0 {
									tolerance := expBitrate * 0.1
									diff := actBitrate - expBitrate
									if diff < 0 {
										diff = -diff
									}
									bitrateMatches = (diff <= tolerance)
								}
							}
						}
					}
				} // Determine final result and notes

				// Check encoding match
				encodingMatches := true
				if expectedEncoding, hasExpEncoding := validationMap["expectedEncoding"]; hasExpEncoding {
					if actualEncoding, hasActEncoding := validationMap["actualEncoding"]; hasActEncoding {
						if expEnc, ok1 := expectedEncoding.(string); ok1 && expEnc != "" {
							if actEnc, ok2 := actualEncoding.(string); ok2 && actEnc != "" {
								encodingMatches = strings.EqualFold(actEnc, expEnc)
							}
						}
					}
				}

				if resolutionMatches && fpsMatches && bitrateMatches && encodingMatches {
					result = "PASS"
					notes.WriteString("All parameters match expected values")
				} else if resolutionMatches {
					// Resolution matches but FPS/bitrate/encoding doesn't = warning
					result = "WARNING"
					if !fpsMatches {
						notes.WriteString("FPS mismatch")
					}
					if !bitrateMatches {
						if notes.Len() > 0 {
							notes.WriteString("; ")
						}
						notes.WriteString("Bitrate mismatch")
					}
					if !encodingMatches {
						if notes.Len() > 0 {
							notes.WriteString("; ")
						}
						notes.WriteString("Encoding mismatch")
					}
				} else {
					// Resolution doesn't match = fail (this shouldn't happen if isValid=true, but just in case)
					result = "FAIL"
					notes.WriteString("Resolution mismatch")
				}
			} else {
				notes.WriteString("Validation failed")
			}
		} else {
			// Get error message if available
			if errorMsg, hasError := validationMap["error"]; hasError {
				if msg, ok := errorMsg.(string); ok && msg != "" {
					notes.WriteString(msg)
				} else {
					notes.WriteString("Unknown error")
				}
			} else {
				notes.WriteString("Validation failed without error details")
			}
		}

		// Format resolution expected
		resoExpected := ""
		if expectedWidth, hasWidth := validationMap["expectedWidth"]; hasWidth {
			if expectedHeight, hasHeight := validationMap["expectedHeight"]; hasHeight {
				if w, ok1 := expectedWidth.(float64); ok1 {
					if h, ok2 := expectedHeight.(float64); ok2 {
						resoExpected = fmt.Sprintf("%dx%d", int(w), int(h))
					}
				}
			}
		}

		// Format resolution actual
		resoActual := ""
		if actualWidth, hasWidth := validationMap["actualWidth"]; hasWidth {
			if actualHeight, hasHeight := validationMap["actualHeight"]; hasHeight {
				if w, ok1 := actualWidth.(float64); ok1 {
					if h, ok2 := actualHeight.(float64); ok2 {
						if w > 0 && h > 0 {
							resoActual = fmt.Sprintf("%dx%d", int(w), int(h))
						}
					}
				}
			}
		}

		// Format FPS expected
		fpsExpected := ""
		if expectedFPS, exists := validationMap["expectedFPS"]; exists {
			if fps, ok := expectedFPS.(float64); ok {
				fpsExpected = strconv.Itoa(int(fps))
			}
		}

		// Format FPS actual
		fpsActual := ""
		if actualFPS, exists := validationMap["actualFPS"]; exists {
			if fps, ok := actualFPS.(float64); ok && fps > 0 {
				fpsActual = fmt.Sprintf("%.2f", fps)
			}
		}

		// Format encoding expected
		encodingExpected := ""
		if expectedEncoding, exists := validationMap["expectedEncoding"]; exists {
			if enc, ok := expectedEncoding.(string); ok {
				encodingExpected = enc
			}
		}

		// Format encoding actual
		encodingActual := ""
		if actualEncoding, exists := validationMap["actualEncoding"]; exists {
			if enc, ok := actualEncoding.(string); ok {
				encodingActual = enc
			}
		}

		// Write CSV row with IP column and notes
		row := []string{cameraID, cameraIP, result, resoExpected, resoActual, fpsExpected, fpsActual, encodingExpected, encodingActual, notes.String()}
		if err := writer.Write(row); err != nil {
			return "", fmt.Errorf("failed to write CSV row for camera %s: %v", cameraID, err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return "", fmt.Errorf("failed to flush CSV writer: %v", err)
	}

	return csvBuilder.String(), nil
}

// extractNumericID extracts the numeric part from a camera ID
// Returns -1 if no numeric part is found
func extractNumericID(id string) int {
	// First, try to parse the entire ID as a number (for cases like "1", "2", etc.)
	if num, err := strconv.Atoi(id); err == nil {
		return num
	}

	// If that fails, extract numeric part using regex
	re := regexp.MustCompile(`^[^\d]*(\d+).*$`)
	matches := re.FindStringSubmatch(id)

	if len(matches) > 1 {
		if num, err := strconv.Atoi(matches[1]); err == nil {
			return num
		}
	}
	return -1
}

func convertValidationToMap(validation interface{}) (map[string]interface{}, error) {
	// If it's already a map, return it directly
	if validationMap, ok := validation.(map[string]interface{}); ok {
		return validationMap, nil
	}

	// If it's an array, convert it to a map using indices as keys
	if validationArray, ok := validation.([]interface{}); ok {
		result := make(map[string]interface{})
		for i, item := range validationArray {
			if itemMap, ok := item.(map[string]interface{}); ok {
				// Try to get camera ID from the item
				cameraID := fmt.Sprintf("camera%d", i+1) // Default camera ID
				if id, exists := itemMap["cameraId"]; exists {
					if idStr, ok := id.(string); ok && idStr != "" {
						cameraID = idStr
					}
				}
				result[cameraID] = itemMap
			} else {
				return nil, fmt.Errorf("invalid validation item at index %d: expected object", i)
			}
		}
		return result, nil
	}

	return nil, fmt.Errorf("validation data must be either an object or an array of objects")
}

func HandleImportCamerasCSV(w http.ResponseWriter, r *http.Request) {
	log.Println("Received /cameras/import-csv request")

	// Parse multipart form with a 10MB size limit
	err := r.ParseMultipartForm(10 << 20) // 10MB
	if err != nil {
		log.Printf("Error parsing multipart form: %v", err)
		http.Error(w, fmt.Sprintf("Failed to parse form: %v", err), http.StatusBadRequest)
		return
	}

	// Get the CSV file from the form
	file, header, err := r.FormFile("csvFile")
	if err != nil {
		log.Printf("Error getting CSV file from form: %v", err)
		http.Error(w, "CSV file is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	log.Printf("Processing CSV file: %s (size: %d bytes)", header.Filename, header.Size)

	// Read and parse CSV content
	csvReader := csv.NewReader(file)
	csvReader.FieldsPerRecord = -1 // Allow variable number of fields

	// Read all records
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Printf("Error reading CSV file: %v", err)
		http.Error(w, fmt.Sprintf("Failed to read CSV file: %v", err), http.StatusBadRequest)
		return
	}

	if len(records) == 0 {
		log.Println("Error: CSV file is empty")
		http.Error(w, "CSV file is empty", http.StatusBadRequest)
		return
	}

	// Parse header to determine column indices
	header_row := records[0]
	columnIndices := make(map[string]int)

	for i, column := range header_row {
		columnName := strings.ToLower(strings.TrimSpace(column))
		columnIndices[columnName] = i
	}

	// Required columns
	requiredColumns := []string{"ip", "username"}
	for _, reqCol := range requiredColumns {
		if _, exists := columnIndices[reqCol]; !exists {
			log.Printf("Error: Required column '%s' not found in CSV", reqCol)
			http.Error(w, fmt.Sprintf("Required column '%s' not found in CSV header", reqCol), http.StatusBadRequest)
			return
		}
	}

	log.Printf("CSV header parsed successfully. Found columns: %v", columnIndices)

	// Process each data row
	var results []map[string]interface{}
	var successCount, errorCount int

	for rowIndex, record := range records[1:] { // Skip header row
		rowNum := rowIndex + 2 // +2 because we start from row 1 (skipping header) and want 1-based numbering

		log.Printf("Processing row %d: %v", rowNum, record)

		// Extract camera data with defaults
		cameraData := struct {
			IP       string
			Port     int
			URL      string
			Username string
			Password string
		}{
			Port: 80, // Default ONVIF port
			URL:  "", // Default empty
		}

		// Extract IP (required)
		if ipIndex, exists := columnIndices["ip"]; exists && ipIndex < len(record) {
			cameraData.IP = strings.TrimSpace(record[ipIndex])
		}
		if cameraData.IP == "" {
			log.Printf("Row %d: Missing IP address", rowNum)
			results = append(results, map[string]interface{}{
				"row":     rowNum,
				"success": false,
				"error":   "Missing IP address",
				"data":    record,
			})
			errorCount++
			continue
		}

		// Extract Username (required)
		if usernameIndex, exists := columnIndices["username"]; exists && usernameIndex < len(record) {
			cameraData.Username = strings.TrimSpace(record[usernameIndex])
		}
		if cameraData.Username == "" {
			log.Printf("Row %d: Missing username", rowNum)
			results = append(results, map[string]interface{}{
				"row":     rowNum,
				"success": false,
				"error":   "Missing username",
				"data":    record,
			})
			errorCount++
			continue
		}

		// Extract optional fields
		if portIndex, exists := columnIndices["port"]; exists && portIndex < len(record) {
			if portStr := strings.TrimSpace(record[portIndex]); portStr != "" {
				if port, err := strconv.Atoi(portStr); err == nil {
					cameraData.Port = port
				} else {
					log.Printf("Row %d: Invalid port value '%s', using default 80", rowNum, portStr)
				}
			}
		}

		if urlIndex, exists := columnIndices["url"]; exists && urlIndex < len(record) {
			cameraData.URL = strings.TrimSpace(record[urlIndex])
		}

		if passwordIndex, exists := columnIndices["password"]; exists && passwordIndex < len(record) {
			cameraData.Password = strings.TrimSpace(record[passwordIndex])
		}

		// Attempt to add the camera
		log.Printf("Adding camera from row %d: IP=%s, Port=%d, Username=%s",
			rowNum, cameraData.IP, cameraData.Port, cameraData.Username)

		newID, err := camera.AddNewCamera(cameraData.IP, cameraData.Port, cameraData.URL, cameraData.Username, cameraData.Password)
		if err != nil {
			log.Printf("Row %d: Failed to add camera: %v", rowNum, err)
			results = append(results, map[string]interface{}{
				"row":     rowNum,
				"success": false,
				"error":   err.Error(),
				"data":    record,
				"camera":  cameraData,
			})
			errorCount++
		} else {
			log.Printf("Row %d: Successfully added camera with ID: %s", rowNum, newID)
			results = append(results, map[string]interface{}{
				"row":      rowNum,
				"success":  true,
				"cameraId": newID,
				"camera": models.Camera{
					ID:       newID,
					IP:       cameraData.IP,
					Port:     cameraData.Port,
					URL:      cameraData.URL,
					Username: cameraData.Username,
					Password: cameraData.Password,
				},
			})
			successCount++
		}
	}

	log.Printf("CSV import completed: %d successful, %d errors", successCount, errorCount)

	// Prepare response
	response := map[string]interface{}{
		"message":      fmt.Sprintf("CSV import completed: %d cameras added successfully, %d errors", successCount, errorCount),
		"totalRows":    len(records) - 1, // Exclude header
		"successCount": successCount,
		"errorCount":   errorCount,
		"results":      results,
	}

	w.Header().Set("Content-Type", "application/json")
	if errorCount > 0 && successCount == 0 {
		w.WriteHeader(http.StatusBadRequest)
	} else if errorCount > 0 {
		w.WriteHeader(http.StatusPartialContent) // 206 - Some succeeded, some failed
	} else {
		w.WriteHeader(http.StatusCreated) // 201 - All succeeded
	}
	json.NewEncoder(w).Encode(response)
}

func HandleImportConfigCSV(w http.ResponseWriter, r *http.Request) {
	log.Println("Received /import-config-csv request")

	// Parse multipart form with a 10MB size limit
	err := r.ParseMultipartForm(10 << 20) // 10MB
	if err != nil {
		log.Printf("Error parsing multipart form: %v", err)
		http.Error(w, fmt.Sprintf("Failed to parse form: %v", err), http.StatusBadRequest)
		return
	}

	// Get the CSV file from the form
	file, header, err := r.FormFile("csvFile")
	if err != nil {
		log.Printf("Error getting CSV file from form: %v", err)
		http.Error(w, "CSV file is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	log.Printf("Processing config CSV file: %s (size: %d bytes)", header.Filename, header.Size)

	// Read and parse CSV content
	csvReader := csv.NewReader(file)
	csvReader.FieldsPerRecord = -1 // Allow variable number of fields

	// Read all records
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Printf("Error reading CSV file: %v", err)
		http.Error(w, fmt.Sprintf("Failed to read CSV file: %v", err), http.StatusBadRequest)
		return
	}

	if len(records) == 0 {
		log.Println("Error: CSV file is empty")
		http.Error(w, "CSV file is empty", http.StatusBadRequest)
		return
	}

	if len(records) < 2 {
		log.Println("Error: CSV file must contain at least header and one data row")
		http.Error(w, "CSV file must contain header and configuration data", http.StatusBadRequest)
		return
	}

	// Parse header to determine column indices
	headerRow := records[0]
	columnIndices := make(map[string]int)

	for i, column := range headerRow {
		columnName := strings.ToLower(strings.TrimSpace(column))
		columnIndices[columnName] = i
	}

	// Required columns (bitrate is optional)
	requiredColumns := []string{"width", "height", "fps"}
	for _, reqCol := range requiredColumns {
		if _, exists := columnIndices[reqCol]; !exists {
			log.Printf("Error: Required column '%s' not found in CSV", reqCol)
			http.Error(w, fmt.Sprintf("Required column '%s' not found in CSV header", reqCol), http.StatusBadRequest)
			return
		}
	}

	log.Printf("Config CSV header parsed successfully. Found columns: %v", columnIndices)

	// Process the first data row (should only be 1 row)
	dataRow := records[1]
	log.Printf("Processing config data row: %v", dataRow)

	// Parse configuration values
	configData := struct {
		Width   int
		Height  int
		FPS     int
		Bitrate int // Optional, defaults to 0 if not provided
	}{
		Bitrate: 0, // Default value for optional bitrate
	}

	// Extract Width (required)
	if widthIndex, exists := columnIndices["width"]; exists && widthIndex < len(dataRow) {
		widthStr := strings.TrimSpace(dataRow[widthIndex])
		if widthStr == "" {
			log.Println("Error: Width value is empty")
			http.Error(w, "Width value is required", http.StatusBadRequest)
			return
		}
		width, err := strconv.Atoi(widthStr)
		if err != nil || width <= 0 {
			log.Printf("Error: Invalid width value '%s'", widthStr)
			http.Error(w, fmt.Sprintf("Invalid width value: %s", widthStr), http.StatusBadRequest)
			return
		}
		configData.Width = width
	} else {
		log.Println("Error: Width column not found or empty")
		http.Error(w, "Width value is required", http.StatusBadRequest)
		return
	}

	// Extract Height (required)
	if heightIndex, exists := columnIndices["height"]; exists && heightIndex < len(dataRow) {
		heightStr := strings.TrimSpace(dataRow[heightIndex])
		if heightStr == "" {
			log.Println("Error: Height value is empty")
			http.Error(w, "Height value is required", http.StatusBadRequest)
			return
		}
		height, err := strconv.Atoi(heightStr)
		if err != nil || height <= 0 {
			log.Printf("Error: Invalid height value '%s'", heightStr)
			http.Error(w, fmt.Sprintf("Invalid height value: %s", heightStr), http.StatusBadRequest)
			return
		}
		configData.Height = height
	} else {
		log.Println("Error: Height column not found or empty")
		http.Error(w, "Height value is required", http.StatusBadRequest)
		return
	}

	// Extract FPS (required)
	if fpsIndex, exists := columnIndices["fps"]; exists && fpsIndex < len(dataRow) {
		fpsStr := strings.TrimSpace(dataRow[fpsIndex])
		if fpsStr == "" {
			log.Println("Error: FPS value is empty")
			http.Error(w, "FPS value is required", http.StatusBadRequest)
			return
		}
		fps, err := strconv.Atoi(fpsStr)
		if err != nil || fps <= 0 {
			log.Printf("Error: Invalid FPS value '%s'", fpsStr)
			http.Error(w, fmt.Sprintf("Invalid FPS value: %s", fpsStr), http.StatusBadRequest)
			return
		}
		configData.FPS = fps
	} else {
		log.Println("Error: FPS column not found or empty")
		http.Error(w, "FPS value is required", http.StatusBadRequest)
		return
	}

	// Extract Bitrate (optional)
	if bitrateIndex, exists := columnIndices["bitrate"]; exists && bitrateIndex < len(dataRow) {
		bitrateStr := strings.TrimSpace(dataRow[bitrateIndex])
		if bitrateStr != "" {
			bitrate, err := strconv.Atoi(bitrateStr)
			if err != nil || bitrate < 0 {
				log.Printf("Warning: Invalid bitrate value '%s', using default 0", bitrateStr)
			} else {
				configData.Bitrate = bitrate
			}
		}
	}

	log.Printf("Parsed config from CSV: Width=%d, Height=%d, FPS=%d, Bitrate=%d",
		configData.Width, configData.Height, configData.FPS, configData.Bitrate)

	// Prepare response with parsed configuration
	response := map[string]interface{}{
		"message": "Configuration CSV imported successfully",
		"config": map[string]interface{}{
			"width":   configData.Width,
			"height":  configData.Height,
			"fps":     configData.FPS,
			"bitrate": configData.Bitrate,
		},
		"status": "ready_to_apply",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

	log.Printf("Config CSV import completed successfully: %+v", configData)
}

func HandleChooseCamFromCSV(w http.ResponseWriter, r *http.Request) {
	log.Println("Received /choose-cam-from-csv request")

	// Parse multipart form with a 10MB size limit
	err := r.ParseMultipartForm(10 << 20) // 10MB
	if err != nil {
		log.Printf("Error parsing multipart form: %v", err)
		http.Error(w, fmt.Sprintf("Failed to parse form: %v", err), http.StatusBadRequest)
		return
	}

	// Get the CSV file from the form
	file, header, err := r.FormFile("csvFile")
	if err != nil {
		log.Printf("Error getting CSV file from form: %v", err)
		http.Error(w, "CSV file is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	log.Printf("Processing camera selection CSV file: %s (size: %d bytes)", header.Filename, header.Size)

	// Read and parse CSV content
	csvReader := csv.NewReader(file)
	csvReader.FieldsPerRecord = -1 // Allow variable number of fields

	// Read all records
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Printf("Error reading CSV file: %v", err)
		http.Error(w, fmt.Sprintf("Failed to read CSV file: %v", err), http.StatusBadRequest)
		return
	}

	if len(records) == 0 {
		log.Println("Error: CSV file is empty")
		http.Error(w, "CSV file is empty", http.StatusBadRequest)
		return
	}

	// Parse header to find the IP column
	headerRow := records[0]
	ipColumnIndex := -1

	for i, column := range headerRow {
		columnName := strings.ToLower(strings.TrimSpace(column))
		if columnName == "ip" {
			ipColumnIndex = i
			break
		}
	}

	if ipColumnIndex == -1 {
		log.Println("Error: 'ip' column not found in CSV header")
		http.Error(w, "Required column 'ip' not found in CSV header", http.StatusBadRequest)
		return
	}

	log.Printf("Found IP column at index %d", ipColumnIndex)
	// Get cameras from in-memory storage to match IPs with camera IDs
	cameras := camera.GetAllCameras()

	// Create a map of IP to camera ID for quick lookup
	ipToCameraMap := make(map[string]string)
	for _, camera := range cameras {
		ipToCameraMap[camera.IP] = camera.ID
	}

	// Process each data row to extract IPs and find matching cameras
	var selectedCameraIDs []string
	var matchedCameras []models.Camera
	var unmatchedIPs []string
	var invalidRows []map[string]interface{}

	for rowIndex, record := range records[1:] { // Skip header row
		rowNum := rowIndex + 2 // +2 because we start from row 1 (skipping header) and want 1-based numbering

		// Check if row has enough columns
		if ipColumnIndex >= len(record) {
			log.Printf("Row %d: Insufficient columns, expected at least %d", rowNum, ipColumnIndex+1)
			invalidRows = append(invalidRows, map[string]interface{}{
				"row":   rowNum,
				"error": "Insufficient columns",
				"data":  record,
			})
			continue
		}

		ip := strings.TrimSpace(record[ipColumnIndex])
		if ip == "" {
			log.Printf("Row %d: Empty IP address", rowNum)
			invalidRows = append(invalidRows, map[string]interface{}{
				"row":   rowNum,
				"error": "Empty IP address",
				"data":  record,
			})
			continue
		}

		log.Printf("Processing row %d: IP = %s", rowNum, ip)

		// Check if this IP exists in our camera list
		if cameraID, exists := ipToCameraMap[ip]; exists {
			// Find the full camera object
			for _, camera := range cameras {
				if camera.ID == cameraID {
					selectedCameraIDs = append(selectedCameraIDs, cameraID)
					matchedCameras = append(matchedCameras, camera)
					log.Printf("Row %d: Found camera %s for IP %s", rowNum, cameraID, ip)
					break
				}
			}
		} else {
			log.Printf("Row %d: No camera found for IP %s", rowNum, ip)
			unmatchedIPs = append(unmatchedIPs, ip)
		}
	}

	log.Printf("Camera selection completed: %d cameras selected, %d IPs not found, %d invalid rows",
		len(selectedCameraIDs), len(unmatchedIPs), len(invalidRows))

	// Prepare response
	response := map[string]interface{}{
		"message":           fmt.Sprintf("Camera selection completed: %d cameras selected", len(selectedCameraIDs)),
		"totalRows":         len(records) - 1, // Exclude header
		"selectedCameraIds": selectedCameraIDs,
		"selectedCameras":   matchedCameras,
		"matchedCount":      len(selectedCameraIDs),
		"unmatchedIPs":      unmatchedIPs,
		"unmatchedCount":    len(unmatchedIPs),
		"invalidRows":       invalidRows,
		"invalidRowCount":   len(invalidRows),
		"status":            "cameras_selected_for_config",
	}

	w.Header().Set("Content-Type", "application/json")

	// Set appropriate status code based on results
	if len(selectedCameraIDs) == 0 {
		w.WriteHeader(http.StatusBadRequest) // No cameras found
	} else if len(unmatchedIPs) > 0 || len(invalidRows) > 0 {
		w.WriteHeader(http.StatusPartialContent) // Some cameras found, some not
	} else {
		w.WriteHeader(http.StatusOK) // All cameras found
	}

	json.NewEncoder(w).Encode(response)

	log.Printf("Camera selection response sent: %d selected cameras ready for configuration", len(selectedCameraIDs))
}

func HandleConfigSingleCam(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cameraID := vars["id"]

	log.Printf("Received request to configure camera ID: %s", cameraID)

	if cameraID == "" {
		log.Println("Error: Missing camera ID in config request")
		http.Error(w, "Camera ID is required", http.StatusBadRequest)
		return
	}

	var input struct {
		Width    int    `json:"width"`
		Height   int    `json:"height"`
		FPS      int    `json:"fps"`
		Bitrate  int    `json:"bitrate"`
		Encoding string `json:"encoding"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Printf("Error decoding config request body: %v", err)
		http.Error(w, fmt.Sprintf("Failed to decode request body: %v", err), http.StatusBadRequest)
		return
	}

	// Validate input parameters
	if input.Width < 320 || input.Width > 3840 {
		log.Printf("Invalid width %d for camera %s (must be between 320 and 3840)", input.Width, cameraID)
		http.Error(w, "Width must be between 320 and 3840 pixels", http.StatusBadRequest)
		return
	}

	if input.Height < 240 || input.Height > 2160 {
		log.Printf("Invalid height %d for camera %s (must be between 240 and 2160)", input.Height, cameraID)
		http.Error(w, "Height must be between 240 and 2160 pixels", http.StatusBadRequest)
		return
	}

	if input.FPS < 1 || input.FPS > 60 {
		log.Printf("Invalid FPS %d for camera %s (must be between 1 and 60)", input.FPS, cameraID)
		http.Error(w, "FPS must be between 1 and 60", http.StatusBadRequest)
		return
	}

	if input.Bitrate < 100 || input.Bitrate > 50000 {
		log.Printf("Invalid bitrate %d for camera %s (must be between 100 and 50000)", input.Bitrate, cameraID)
		http.Error(w, "Bitrate must be between 100 and 50000 kbps", http.StatusBadRequest)
		return
	}

	validEncodings := map[string]bool{
		"h264":  true,
		"h265":  true,
		"mjpeg": true,
		"H264":  true,
		"H265":  true,
		"MJPEG": true,
	}

	if !validEncodings[input.Encoding] {
		log.Printf("Invalid encoding '%s' for camera %s (must be h264, h265, or mjpeg)", input.Encoding, cameraID)
		http.Error(w, "Encoding must be h264, h265, or mjpeg", http.StatusBadRequest)
		return
	}

	log.Printf("Configuring camera %s: Width=%d, Height=%d, FPS=%d, Bitrate=%d, Encoding=%s",
		cameraID, input.Width, input.Height, input.FPS, input.Bitrate, input.Encoding)

	// Load cameras from CSV to find the requested camera
	cameras, err := loader.LoadCameraList()
	if err != nil {
		log.Printf("Error loading cameras from CSV: %v", err)
		http.Error(w, fmt.Sprintf("Failed to load cameras from CSV: %v", err), http.StatusInternalServerError)
		return
	}

	// Find the camera with the specified ID
	var targetCamera *models.Camera
	for _, cam := range cameras {
		if cam.ID == cameraID {
			targetCamera = &cam
			break
		}
	}

	if targetCamera == nil {
		log.Printf("Camera with ID %s not found in CSV", cameraID)
		http.Error(w, fmt.Sprintf("Camera with ID %s not found", cameraID), http.StatusNotFound)
		return
	}

	// Initialize camera client
	client, err := camera.NewCameraClient(*targetCamera)
	if err != nil {
		log.Printf("Failed to initialize camera %s: %v", targetCamera.ID, err)
		http.Error(w, fmt.Sprintf("Failed to initialize camera: %v", err), http.StatusInternalServerError)
		return
	}

	// Proceed with actual configuration
	log.Printf("Getting profiles and configs for camera %s", targetCamera.ID)
	profileTokens, configTokens, err := camera.GetProfilesAndConfigs(client)
	if err != nil {
		log.Printf("Failed to get camera profiles and configs for %s: %v", targetCamera.ID, err)
		http.Error(w, fmt.Sprintf("Failed to get camera profiles and configs: %v", err), http.StatusInternalServerError)
		return
	}

	if len(profileTokens) == 0 {
		log.Printf("No profiles found for camera %s", targetCamera.ID)
		http.Error(w, "No profiles found", http.StatusInternalServerError)
		return
	}

	if len(configTokens) == 0 {
		log.Printf("No video encoder configuration found for camera %s", targetCamera.ID)
		http.Error(w, "No video encoder configuration found", http.StatusInternalServerError)
		return
	}

	// Use the first token found
	profileToken := profileTokens[0]
	configToken := configTokens[0]

	log.Printf("Using profile token %s and config token %s for camera %s", profileToken, configToken, targetCamera.ID)

	// Get current encoder config for quality value
	currentConfig, err := camera.GetCurrentConfig(client, configToken)
	if err != nil {
		log.Printf("Failed to get current encoder config for %s: %v", targetCamera.ID, err)
		http.Error(w, fmt.Sprintf("Failed to get current encoder config: %v", err), http.StatusInternalServerError)
		return
	}

	// Get available encoder options to find closest resolution
	encoderOptions, err := camera.GetCurrentEncoderOptions(client, profileToken, configToken)
	if err != nil {
		log.Printf("Failed to get encoder options for %s: %v", targetCamera.ID, err)
		http.Error(w, fmt.Sprintf("Failed to get encoder options: %v", err), http.StatusInternalServerError)
		return
	}

	// Find closest matching resolution
	targetResolution := models.Resolution{Width: input.Width, Height: input.Height}
	closestResolution := camera.FindClosestResolution(targetResolution, encoderOptions.Resolutions)
	log.Printf("Closest resolution found for camera %s: %dx%d", targetCamera.ID, closestResolution.Width, closestResolution.Height)

	// Create new configuration
	newConfig := models.EncoderConfig{
		Resolution: closestResolution,
		Quality:    currentConfig.Quality, // Keep the current quality
		FPS:        input.FPS,
		Bitrate:    input.Bitrate,
		Encoding:   input.Encoding,
	}

	// Apply the configuration
	log.Printf("Applying new configuration to camera %s", targetCamera.ID)
	err = camera.SetEncoderConfig(client, configToken, currentConfig, newConfig)
	if err != nil {
		log.Printf("Failed to apply configuration to camera %s: %v", targetCamera.ID, err)
		http.Error(w, fmt.Sprintf("Failed to apply configuration: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully applied configuration to camera %s", targetCamera.ID)

	// Prepare response
	response := map[string]interface{}{
		"cameraId": targetCamera.ID,
		"status":   "success",
		"message":  "Configuration applied successfully",
		"appliedConfig": map[string]interface{}{
			"resolution": map[string]int{
				"width":  closestResolution.Width,
				"height": closestResolution.Height,
			},
			"fps":                input.FPS,
			"bitrate":            input.Bitrate,
			"encoding":           input.Encoding,
			"resolutionAdjusted": input.Width != closestResolution.Width || input.Height != closestResolution.Height,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func HandleValidateCam(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cameraID := vars["id"]

	log.Printf("Received request to validate camera ID: %s", cameraID)

	if cameraID == "" {
		log.Println("Error: Missing camera ID in validation request")
		http.Error(w, "Camera ID is required", http.StatusBadRequest)
		return
	}

	// Load cameras from CSV to find the requested camera
	cameras, err := loader.LoadCameraList()
	if err != nil {
		log.Printf("Error loading cameras from CSV: %v", err)
		http.Error(w, fmt.Sprintf("Failed to load cameras from CSV: %v", err), http.StatusInternalServerError)
		return
	}

	// Find the camera with the specified ID
	var targetCamera *models.Camera
	for _, cam := range cameras {
		if cam.ID == cameraID {
			targetCamera = &cam
			break
		}
	}

	if targetCamera == nil {
		log.Printf("Camera with ID %s not found in CSV", cameraID)
		http.Error(w, fmt.Sprintf("Camera with ID %s not found", cameraID), http.StatusNotFound)
		return
	}

	log.Printf("Found camera %s (IP: %s:%d)", targetCamera.ID, targetCamera.IP, targetCamera.Port)

	// Create a camera client for this specific camera
	client, err := camera.NewCameraClient(*targetCamera)
	if err != nil {
		log.Printf("Error creating camera client for %s: %v", cameraID, err)
		http.Error(w, fmt.Sprintf("Failed to create camera client: %v", err), http.StatusInternalServerError)
		return
	}

	// Get current camera configuration to use as expected values
	log.Printf("Getting current configuration for camera %s", cameraID)
	profileTokens, configTokens, err := camera.GetProfilesAndConfigs(client)
	if err != nil {
		log.Printf("Failed to get camera profiles and configs for %s: %v", cameraID, err)
		http.Error(w, fmt.Sprintf("Failed to get camera profiles and configs: %v", err), http.StatusInternalServerError)
		return
	}

	if len(profileTokens) == 0 || len(configTokens) == 0 {
		log.Printf("No profiles or configs found for camera %s", cameraID)
		http.Error(w, "No profiles or configs found", http.StatusInternalServerError)
		return
	}

	// Get current encoder config
	currentConfig, err := camera.GetCurrentConfig(client, configTokens[0])
	if err != nil {
		log.Printf("Failed to get current encoder config for %s: %v", cameraID, err)
		http.Error(w, fmt.Sprintf("Failed to get current encoder config: %v", err), http.StatusInternalServerError)
		return
	}

	// Get stream URI
	streamURI, err := client.GetStreamURI(profileTokens[0])
	if err != nil {
		log.Printf("Failed to get stream URI for camera %s: %v", cameraID, err)
		http.Error(w, fmt.Sprintf("Failed to get stream URI: %v", err), http.StatusInternalServerError)
		return
	}

	// Add credentials to the stream URL
	parsedURI, err := url.Parse(streamURI)
	if err != nil {
		log.Printf("Failed to parse stream URI for camera %s: %v", cameraID, err)
		http.Error(w, fmt.Sprintf("Failed to parse stream URI: %v", err), http.StatusInternalServerError)
		return
	}

	// Add username and password to the URL
	parsedURI.User = url.UserPassword(client.Camera.Username, client.Camera.Password)
	authenticatedStreamURI := parsedURI.String()

	log.Printf("Validating stream URI: %s", authenticatedStreamURI)

	// Use RTSP analyzer to validate the stream
	validationResult, err := ffmpeg.ValidateStream(
		authenticatedStreamURI,
		currentConfig.Resolution.Width,
		currentConfig.Resolution.Height,
		currentConfig.FPS,
		currentConfig.Bitrate,
		currentConfig.Encoding,
	)

	if err != nil {
		log.Printf("Failed to validate stream for camera %s: %v", cameraID, err)
		http.Error(w, fmt.Sprintf("Failed to validate stream: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("Validation completed for camera %s: valid=%t", cameraID, validationResult.IsValid)

	// Prepare response
	response := map[string]interface{}{
		"cameraId":         targetCamera.ID,
		"isValid":          validationResult.IsValid,
		"message":          validationResult.Error,
		"validationResult": validationResult,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// pingHost checks if a host is reachable via ping
func pingHost(host string, timeout time.Duration) (bool, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		// Windows: ping -n 1 -w 3000 hostname
		cmd = exec.CommandContext(ctx, "ping", "-n", "1", "-w", "3000", host)
	} else {
		// Linux/macOS: ping -c 1 -W 3 hostname
		cmd = exec.CommandContext(ctx, "ping", "-c", "1", "-W", "3", host)
	}

	output, err := cmd.CombinedOutput()
	outputStr := string(output)

	if ctx.Err() == context.DeadlineExceeded {
		return false, "Ping timeout", fmt.Errorf("ping timeout after %v", timeout)
	}

	// Check for unreachable host patterns in the output (cross-platform)
	outputLower := strings.ToLower(outputStr)
	unreachablePatterns := []string{
		"destination host unreachable",
		"destination net unreachable",
		"host unreachable",
		"no route to host",
		"request timed out",
		"could not find host",
		"network is unreachable",          // Linux specific
		"connect: network is unreachable", // Linux specific
		"ping: cannot resolve",            // Linux DNS failure
		"name or service not known",       // Linux DNS failure
	}

	for _, pattern := range unreachablePatterns {
		if strings.Contains(outputLower, pattern) {
			return false, outputStr, fmt.Errorf("host unreachable: %s", pattern)
		}
	}

	// Platform-specific success/failure detection
	if runtime.GOOS == "windows" {
		// Windows: look for successful ping patterns
		if strings.Contains(outputLower, "reply from") && !strings.Contains(outputLower, "unreachable") {
			return true, outputStr, nil
		}
		// If no "reply from" or contains unreachable, it's failed
		return false, outputStr, fmt.Errorf("ping failed: no valid reply received")
	} else {
		// Linux/macOS: check for successful ping patterns
		successPatterns := []string{
			"1 received",
			"1 packets transmitted, 1 received",
			"1 packets transmitted, 1 packets received", // Some Linux variants
			"rtt min/avg/max",                           // Success indicator
		}

		hasSuccess := false
		for _, pattern := range successPatterns {
			if strings.Contains(outputLower, pattern) {
				hasSuccess = true
				break
			}
		}

		if hasSuccess {
			return true, outputStr, nil
		}

		// Check exit code as fallback for Linux/macOS
		if err != nil {
			return false, outputStr, err
		}

		if cmd.ProcessState.ExitCode() == 0 {
			return true, outputStr, nil
		}

		return false, outputStr, fmt.Errorf("ping failed: exit code %d", cmd.ProcessState.ExitCode())
	}
}
