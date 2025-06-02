package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"main_back/pkg/models"
)

// CameraConfig is a helper struct to unmarshal the cameras JSON file.
type CameraConfig struct {
	Cameras []models.Camera `json:"cameras"`
}

// LoadCameraList reads a JSON file containing camera configurations.
// It determines the path to the config file relative to the executable's location.
func LoadCameraList() ([]models.Camera, error) {
	// Determine the executable directory
	execPath, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("failed to get executable path: %w", err)
	}
	execDir := filepath.Dir(execPath)

	// Construct the path to cameras.json relative to the executable directory
	// Corrected path: go up two directories from the executable's location to reach the project root relative to the config folder
	configPath := filepath.Join(execDir, "..", "..", "config", "cameras.json")

	// Read the JSON file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read camera config file %s: %w", configPath, err)
	}

	var config CameraConfig
	// Unmarshal the JSON data into the struct
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal camera config file %s: %w", configPath, err)
	}

	return config.Cameras, nil
}
