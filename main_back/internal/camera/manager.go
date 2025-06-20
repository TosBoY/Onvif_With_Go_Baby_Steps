package camera

import (
	"encoding/json"
	"fmt"
	"main_back/pkg/models"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var connectedCameras map[string]*CameraClient

func init() {
	connectedCameras = make(map[string]*CameraClient)
}

// InitializeAllCameras connects to all provided cameras and stores their clients.
// It attempts to initialize all cameras and returns a single error containing details
// of all failures, or nil if all cameras were successfully initialized.
func InitializeAllCameras(cameras []models.Camera) error {
	var initializationErrors []string

	for _, cam := range cameras {
		client, err := NewCameraClient(cam)
		if err != nil {
			initializationErrors = append(initializationErrors, fmt.Sprintf("failed to initialize camera %s: %v", cam.ID, err))
		} else {
			connectedCameras[cam.ID] = client
		}
	}

	if len(initializationErrors) > 0 {
		return fmt.Errorf("failed to initialize some cameras:\n%s", strings.Join(initializationErrors, "\n"))
	}

	return nil
}

// GetCameraClient returns a connected camera client by its ID.
func GetCameraClient(id string) (*CameraClient, error) {
	client, ok := connectedCameras[id]
	if !ok {
		return nil, fmt.Errorf("camera with ID %s not found", id)
	}
	return client, nil
}

// AddNewCamera adds a new camera to the cameras.json file and assigns it an ID
// that is one greater than the largest existing ID.
// Returns the new camera ID and any error encountered.
func AddNewCamera(ip string, port int, url string, username string, password string, isFake bool) (string, error) {
	// Get the path to cameras.json
	workingDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get working directory: %w", err)
	}

	configPath := filepath.Join(workingDir, "..", "..", "config", "cameras.json")

	// Read the existing cameras
	data, err := os.ReadFile(configPath)
	if err != nil {
		return "", fmt.Errorf("failed to read camera config file %s: %w", configPath, err)
	}

	var config struct {
		Cameras []models.Camera `json:"cameras"`
	}
	// Unmarshal the JSON data
	if err := json.Unmarshal(data, &config); err != nil {
		return "", fmt.Errorf("failed to parse camera config: %w", err)
	}

	// Check if a camera with the same IP already exists
	for _, cam := range config.Cameras {
		if cam.IP == ip {
			return "", fmt.Errorf("camera with IP address %s already exists (ID: %s)", ip, cam.ID)
		}
	}

	// Find the highest ID
	highestID := 0
	for _, cam := range config.Cameras {
		// Try to convert the ID to an integer
		if camID, err := strconv.Atoi(cam.ID); err == nil {
			if camID > highestID {
				highestID = camID
			}
		}
	}

	// Create a new ID by incrementing the highest ID
	newID := strconv.Itoa(highestID + 1) // Create a new camera
	newCamera := models.Camera{
		ID:       newID,
		IP:       ip,
		Port:     port,
		URL:      url,
		Username: username,
		Password: password,
		IsFake:   isFake,
	}

	// Add the new camera to the list
	config.Cameras = append(config.Cameras, newCamera)

	// Marshal the updated config with indentation for better readability
	updatedData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal updated camera config: %w", err)
	}

	// Write the updated config back to the file
	if err := os.WriteFile(configPath, updatedData, 0644); err != nil {
		return "", fmt.Errorf("failed to write updated camera config: %w", err)
	}
	// Return the new camera ID
	return newID, nil
}

// RemoveCamera removes a camera from the cameras.json file by its ID.
// It also removes the camera client from the connected cameras map if it exists.
// Returns any error encountered.
func RemoveCamera(id string) error {
	// Get the path to cameras.json
	workingDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	configPath := filepath.Join(workingDir, "..", "..", "config", "cameras.json")

	// Read the existing cameras
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read camera config file %s: %w", configPath, err)
	}

	var config struct {
		Cameras []models.Camera `json:"cameras"`
	}

	// Unmarshal the JSON data
	if err := json.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse camera config: %w", err)
	}

	// Find the camera and remove it
	found := false
	updatedCameras := make([]models.Camera, 0, len(config.Cameras))
	for _, cam := range config.Cameras {
		if cam.ID != id {
			updatedCameras = append(updatedCameras, cam)
		} else {
			found = true
			// Remove the camera from the connected cameras map if it exists
			delete(connectedCameras, id)
		}
	}

	if !found {
		return fmt.Errorf("camera with ID %s not found", id)
	}

	// Update the cameras list
	config.Cameras = updatedCameras

	// Marshal the updated config with indentation for better readability
	updatedData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal updated camera config: %w", err)
	}

	// Write the updated config back to the file
	if err := os.WriteFile(configPath, updatedData, 0644); err != nil {
		return fmt.Errorf("failed to write updated camera config: %w", err)
	}

	return nil
}
