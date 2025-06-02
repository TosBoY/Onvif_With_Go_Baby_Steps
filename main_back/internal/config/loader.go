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

	workingDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get working directory: %w", err)
	}

	configPath := filepath.Join(workingDir, "..", "..", "config", "cameras.json")

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
