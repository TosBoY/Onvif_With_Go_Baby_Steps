package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"main_back/internal/camera"
	"main_back/internal/config"
	"main_back/internal/ffprobe"
	"main_back/internal/vlc"
	"main_back/pkg/models"

	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/cameras", HandleGetCameras).Methods("GET")
	r.HandleFunc("/cameras", HandleAddCamera).Methods("POST")
	r.HandleFunc("/cameras/{id}", HandleDeleteCamera).Methods("DELETE")
	r.HandleFunc("/apply-config", HandleApplyConfig).Methods("POST")
	r.HandleFunc("/vlc", HandleVLC).Methods("POST")
}

func HandleGetCameras(w http.ResponseWriter, r *http.Request) {
	cameras, err := config.LoadCameraList()
	if err != nil {
		log.Printf("Error loading camera list: %v", err)
		http.Error(w, "Failed to load camera list", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(cameras)
}

func HandleAddCamera(w http.ResponseWriter, r *http.Request) {
	log.Println("Received /cameras POST request to add a new camera")

	var input struct {
		IP       string `json:"ip"`
		Username string `json:"username"`
		Password string `json:"password"`
		IsFake   bool   `json:"isFake"`
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

	log.Printf("Adding new camera with IP: %s, Username: %s, IsFake: %v", input.IP, input.Username, input.IsFake)
	newID, err := camera.AddNewCamera(input.IP, input.Username, input.Password, input.IsFake)
	if err != nil {
		log.Printf("Error adding new camera: %v", err)
		http.Error(w, fmt.Sprintf("Failed to add new camera: %v", err), http.StatusInternalServerError)
		return
	}

	// Return the new camera ID and details
	newCamera := models.Camera{
		ID:       newID,
		IP:       input.IP,
		Username: input.Username,
		Password: input.Password,
		IsFake:   input.IsFake,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newCamera)
}

// HandleDeleteCamera handles DELETE requests for removing cameras from the system.
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

func HandleApplyConfig(w http.ResponseWriter, r *http.Request) {
	log.Println("Received /apply-config request")

	var input struct {
		CameraID  string   `json:"cameraId"`  // For backward compatibility
		CameraIDs []string `json:"cameraIds"` // New field for multiple cameras
		Width     int      `json:"width"`
		Height    int      `json:"height"`
		FPS       int      `json:"fps"`
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
		log.Printf("Applying config for multiple cameras (%d): Width: %d, Height: %d, FPS: %d",
			len(cameraIDs), input.Width, input.Height, input.FPS)
	} else if input.CameraID != "" {
		cameraIDs = []string{input.CameraID}
		log.Printf("Applying config for single camera %s: Width: %d, Height: %d, FPS: %d",
			input.CameraID, input.Width, input.Height, input.FPS)
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
		IsFake             bool
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

		// Check if camera is fake and handle it differently
		if client.Camera.IsFake {
			log.Printf("Camera %s is a simulated device, skipping real configuration", cameraID)

			// For fake cameras, we simulate successful config application
			result.Success = true
			result.IsFake = true
			result.AppliedConfig = map[string]interface{}{
				"resolution": map[string]int{
					"width":  input.Width,
					"height": input.Height,
				},
				"fps": input.FPS,
			}
			results[cameraID] = result
			continue
		}

		// For real cameras, proceed with config application
		log.Printf("\n Getting profiles and configs for camera %s", cameraID)
		profileTokens, configTokens, err := camera.GetProfilesAndConfigs(client)
		if err != nil {
			log.Printf("Failed to get camera profiles and configs for %s: %v", cameraID, err)
			result.Error = fmt.Errorf("failed to get camera profiles and configs: %w", err)
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

		// Prepare the new configuration
		newConfig := models.EncoderConfig{
			Resolution: closestResolution,
			Quality:    currentConfig.Quality, // Keep the current quality
			FPS:        input.FPS,
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

		// The ONVIF GetStreamUri typically doesn't include credentials. Embed them for ffprobe.
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
			"fps": input.FPS,
		}
		result.ResolutionAdjusted = input.Width != closestResolution.Width || input.Height != closestResolution.Height
		result.StreamURL = fullStreamURL

		results[cameraID] = result
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

		if !result.Success || result.IsFake {
			// Skip validation for failed configs or fake cameras
			if result.IsFake {
				position := findCameraPosition(cameraID)
				log.Printf("Adding simulated validation results for fake camera %s (position %d in original order)", cameraID, position)
				// Add simulated validation results for fake cameras
				validationResults[cameraID] = map[string]interface{}{
					"isValid":        true,
					"expectedWidth":  input.Width,
					"expectedHeight": input.Height,
					"expectedFPS":    input.FPS,
					"actualWidth":    input.Width, // For fake cameras, actual matches expected
					"actualHeight":   input.Height,
					"actualFPS":      float64(input.FPS),
					"streamInfo": map[string]interface{}{
						"width":    input.Width,
						"height":   input.Height,
						"fps":      input.FPS,
						"codec":    "h264 (simulated)",
						"bitrate":  "variable (simulated)",
						"duration": "live stream",
					},
					"message": "Simulated camera configuration applied successfully",
				}
			}
			continue
		}

		// Validate the stream using ffprobe
		position := findCameraPosition(cameraID)
		log.Printf("Starting ffprobe validation for camera %s (position %d in original order)", cameraID, position)
		validationResult, validationErr := ffprobe.ValidateStream(result.StreamURL, input.Width, input.Height, input.FPS)

		if validationErr != nil {
			log.Printf("FFprobe validation failed for camera %s: %v", cameraID, validationErr)
			validationResults[cameraID] = map[string]interface{}{
				"isValid":        false,
				"error":          validationErr.Error(),
				"expectedWidth":  input.Width,
				"expectedHeight": input.Height,
				"expectedFPS":    input.FPS,
			}
		} else {
			// Create a map from the validation result
			validationMap := map[string]interface{}{
				"isValid":        validationResult.IsValid,
				"actualWidth":    validationResult.ActualWidth,
				"actualHeight":   validationResult.ActualHeight,
				"actualFPS":      validationResult.ActualFPS,
				"expectedWidth":  validationResult.ExpectedWidth,
				"expectedHeight": validationResult.ExpectedHeight,
				"expectedFPS":    validationResult.ExpectedFPS,
			}

			// Add error if present
			if validationResult.Error != "" {
				validationMap["error"] = validationResult.Error
			}

			log.Printf("FFprobe validation completed for camera %s: valid=%v", cameraID, validationResult.IsValid)
			validationResults[cameraID] = validationMap
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
			"fps": input.FPS,
		},
		"results": make(map[string]interface{}),
	}

	// Add individual camera results
	for cameraID, result := range results {
		cameraResult := map[string]interface{}{
			"success": result.Success,
		}

		if result.Success {
			cameraResult["appliedConfig"] = result.AppliedConfig
			cameraResult["resolutionAdjusted"] = result.ResolutionAdjusted
			cameraResult["isFake"] = result.IsFake

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

	// Get the camera client
	client, err := camera.GetCameraClient(input.CameraID)
	if err != nil {
		log.Printf("Error getting camera client for %s: %v", input.CameraID, err)
		http.Error(w, fmt.Sprintf("Camera not found: %v", err), http.StatusNotFound)
		return
	}

	// Check if camera is fake and handle it differently
	if client.Camera.IsFake {
		log.Printf("Camera %s is a simulated device, providing simulated stream URL", input.CameraID)

		// For fake cameras, we'll return a simulated response
		simulatedStreamUrl := fmt.Sprintf("rtsp://fake-stream-simulation/%s", input.CameraID)

		response := map[string]interface{}{
			"message":   "Simulated stream URL generated (no actual VLC launched for simulated camera)",
			"streamUrl": simulatedStreamUrl,
			"isFake":    true,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
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
