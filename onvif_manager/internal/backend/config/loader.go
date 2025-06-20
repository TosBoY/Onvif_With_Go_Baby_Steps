package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"onvif_manager/pkg/models"
)

// CameraConfig is a helper struct to unmarshal the cameras JSON file.
type CameraConfig struct {
	Cameras []models.Camera `json:"cameras"`
}

// FindConfigPath attempts to locate the config directory using multiple strategies
func FindConfigPath() (string, error) {
	// Strategy 1: Try relative to current working directory (for built executable)
	workingDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get working directory: %w", err)
	}

	candidatePaths := []string{
		// For built executable in onvif_manager directory
		filepath.Join(workingDir, "config", "cameras.json"),
		// For go run from cmd/app directory
		filepath.Join(workingDir, "..", "..", "config", "cameras.json"),
	}

	// Try each candidate path
	for _, path := range candidatePaths {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	// Strategy 2: Try to find config relative to executable location
	execPath, err := os.Executable()
	if err == nil {
		execDir := filepath.Dir(execPath)
		execConfigPath := filepath.Join(execDir, "config", "cameras.json")
		if _, err := os.Stat(execConfigPath); err == nil {
			return execConfigPath, nil
		}
	}

	// If all else fails, return the primary candidate with error info
	primaryPath := filepath.Join(workingDir, "config", "cameras.json")
	return primaryPath, fmt.Errorf("config file not found. Tried multiple locations, primary expected location: %s", primaryPath)
}

// LoadCameraList reads a JSON file containing camera configurations.
// It determines the path to the config file relative to the executable's location.
func LoadCameraList() ([]models.Camera, error) {
	configPath, err := FindConfigPath()
	if err != nil {
		return nil, fmt.Errorf("failed to find config path: %w", err)
	}

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
