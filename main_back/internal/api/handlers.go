package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"main_back/internal/camera"
	"main_back/internal/config"
	"main_back/internal/ffprobe"
	"main_back/pkg/models"

	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/cameras", HandleGetCameras).Methods("GET")
	r.HandleFunc("/apply-config", HandleApplyConfig).Methods("POST")
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

func HandleApplyConfig(w http.ResponseWriter, r *http.Request) {
	log.Println("Received /apply-config request")

	var input struct {
		CameraID string `json:"cameraId"`
		Width    int    `json:"width"`
		Height   int    `json:"height"`
		FPS      int    `json:"fps"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Printf("Error decoding apply-config request body: %v", err)
		http.Error(w, fmt.Sprintf("Failed to decode request body: %v", err), http.StatusBadRequest)
		return
	}

	log.Printf("Applying config for CameraID: %s, Width: %d, Height: %d, FPS: %d", input.CameraID, input.Width, input.Height, input.FPS)

	client, err := camera.GetCameraClient(input.CameraID)
	if err != nil {
		log.Printf("Error getting camera client for %s: %v", input.CameraID, err)
		http.Error(w, fmt.Sprintf("Camera not found: %v", err), http.StatusNotFound)
		return
	}

	// Get profiles and configs
	log.Printf("Getting profiles and configs for camera %s", input.CameraID)
	profileTokens, configTokens, err := camera.GetProfilesAndConfigs(client)
	if err != nil {
		log.Printf("Failed to get camera profiles and configs for %s: %v", input.CameraID, err)
		http.Error(w, fmt.Sprintf("Failed to get camera profiles and configs: %v", err), http.StatusInternalServerError)
		return
	}

	if len(profileTokens) == 0 {
		log.Printf("No profiles found for camera %s", input.CameraID)
		http.Error(w, "No profiles found", http.StatusInternalServerError)
		return
	}

	if len(configTokens) == 0 {
		log.Printf("No video encoder configuration found for camera %s", input.CameraID)
		http.Error(w, "No video encoder configuration found", http.StatusInternalServerError)
		return
	}

	// Use the first token found
	profileToken := profileTokens[0]
	configToken := configTokens[0]
	log.Printf("Using profile token %s and config token %s for camera %s", profileToken, configToken, input.CameraID)

	// Get current encoder config
	log.Printf("Getting current encoder config for camera %s", input.CameraID)
	currentConfig, err := camera.GetCurrentConfig(client, configToken)
	if err != nil {
		log.Printf("Failed to get current encoder config for %s: %v", input.CameraID, err)
		http.Error(w, fmt.Sprintf("Failed to get current encoder config: %v", err), http.StatusInternalServerError)
		return
	}

	// Get available encoder options
	log.Printf("Getting available encoder options for camera %s", input.CameraID)
	encoderOptions, err := camera.GetCurrentEncoderOptions(client, profileToken, configToken)
	if err != nil {
		log.Printf("Failed to get encoder options for %s: %v", input.CameraID, err)
		http.Error(w, fmt.Sprintf("Failed to get encoder options: %v", err), http.StatusInternalServerError)
		return
	}

	// Find closest matching resolution
	log.Printf("Finding closest matching resolution for camera %s", input.CameraID)
	targetResolution := models.Resolution{Width: input.Width, Height: input.Height}
	closestResolution := camera.FindClosestResolution(targetResolution, encoderOptions.Resolutions)
	log.Printf("Closest resolution found for camera %s: %dx%d", input.CameraID, closestResolution.Width, closestResolution.Height)

	// Prepare the new configuration
	newConfig := models.EncoderConfig{
		Resolution: closestResolution,
		Quality:    currentConfig.Quality, // Keep the current quality
		FPS:        input.FPS,
	}
	log.Printf("Prepared new config for camera %s: %+v", input.CameraID, newConfig)

	// Set the new encoder config
	log.Printf("Setting new encoder config for camera %s", input.CameraID)
	if err := camera.SetEncoderConfig(client, configToken, currentConfig, newConfig); err != nil {
		log.Printf("Failed to set encoder config for %s: %v", input.CameraID, err)
		http.Error(w, fmt.Sprintf("Failed to set encoder config: %v", err), http.StatusInternalServerError)
		return
	}

	streamURI, err := client.GetStreamURI(profileToken)
	if err != nil {
		log.Printf("Failed to get stream URI for %s: %v", input.CameraID, err)
		http.Error(w, fmt.Sprintf("Failed to get stream URI: %v", err), http.StatusInternalServerError)
		return
	}

	// The ONVIF GetStreamUri typically doesn't include credentials. Embed them for ffprobe.
	parsedURI, err := url.Parse(streamURI)
	if err != nil {
		log.Printf("Failed to parse stream URI for %s: %v", input.CameraID, err)
		http.Error(w, fmt.Sprintf("Failed to parse stream URI: %v", err), http.StatusInternalServerError)
		return
	}

	// Construct the URL with embedded credentials
	fullStreamURL := fmt.Sprintf("%s://%s:%s@%s%s", parsedURI.Scheme, client.Camera.Username, client.Camera.Password, parsedURI.Host, parsedURI.RequestURI())

	log.Printf("Using stream URL for validation: %s", fullStreamURL)

	// Validate the stream using ffprobe
	log.Printf("Starting ffprobe validation for camera %s", input.CameraID)
	isValid, err := ffprobe.ValidateStream(fullStreamURL, newConfig)
	if err != nil {
		log.Printf("FFprobe validation failed for camera %s: %v", input.CameraID, err)
	} else if !isValid {
		log.Printf("FFprobe validation failed for camera %s: stream mismatch", input.CameraID)
	} else {
		log.Printf("FFprobe validation successful for camera %s", input.CameraID)
	}

	log.Printf("Successfully applied config for camera %s", input.CameraID)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "configuration applied"})
}
