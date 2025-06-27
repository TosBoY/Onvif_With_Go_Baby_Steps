package loader

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"main_back/pkg/models"
)

// findCameraCSVPath dynamically searches for the main_back folder first, then uses the fixed path to cameras.csv
func findCameraCSVPath() (string, error) {
	mainBackDir, err := findMainFolder()
	if err != nil {
		return "", fmt.Errorf("failed to locate main_back folder: %w", err)
	}

	// Fixed path within main_back folder
	csvPath := filepath.Join(mainBackDir, "internal", "loader", "cameras.csv")

	// Verify the CSV file exists at the expected location
	if _, err := os.Stat(csvPath); err != nil {
		return "", fmt.Errorf("cameras.csv not found at expected path %s: %w", csvPath, err)
	}

	return csvPath, nil
}

// findMainFolder searches for the main_back folder starting from common locations
func findMainFolder() (string, error) {
	// Get the current executable path
	execPath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %w", err)
	}
	execDir := filepath.Dir(execPath)

	// Get current working directory
	workingDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get working directory: %w", err)
	}

	// List of starting points to search from
	searchDirs := []string{workingDir, execDir}

	// Search from each starting point
	for _, startDir := range searchDirs {
		if onvifDir := searchUpwardsForMain(startDir); onvifDir != "" {
			return onvifDir, nil
		}
	}

	return "", fmt.Errorf("main_back folder not found from executable path (%s) or working directory (%s)", execDir, workingDir)
}

// searchUpwardsForOnvifManager searches upwards from a starting directory for the onvif_manager folder
func searchUpwardsForMain(startDir string) string {
	currentDir := startDir

	// Search upwards through directory tree (max 10 levels to prevent infinite loops)
	for i := 0; i < 10; i++ {
		// Check if current directory is main_back
		if filepath.Base(currentDir) == "main_back" {
			return currentDir
		}

		// Check if main_back exists as a subdirectory
		mainBackPath := filepath.Join(currentDir, "main_back")
		if info, err := os.Stat(mainBackPath); err == nil && info.IsDir() {
			return mainBackPath
		}

		// Move up one directory level
		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			// Reached root directory
			break
		}
		currentDir = parentDir
	}

	return ""
}

// LoadCameraList reads a CSV file containing camera configurations with cam_id and rtsp columns.
// It dynamically finds the onvif_manager folder and locates the CSV file within it.
func LoadCameraList() ([]models.Camera, error) {
	configPath, err := findCameraCSVPath()
	if err != nil {
		return nil, fmt.Errorf("failed to locate cameras.csv file: %w", err)
	}

	// Open the CSV file
	file, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open camera config file %s: %w", configPath, err)
	}
	defer file.Close()

	// Create a new CSV reader
	csvReader := csv.NewReader(file)

	// Read the header row first
	header, err := csvReader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV header: %w", err)
	}

	// Create a map of column indices for each field
	columnMap := make(map[string]int)
	for i, column := range header {
		columnMap[strings.TrimSpace(strings.ToLower(column))] = i
	}

	// Verify all required columns exist
	requiredColumns := []string{"cam_id", "rtsp"}
	for _, col := range requiredColumns {
		if _, exists := columnMap[col]; !exists {
			return nil, fmt.Errorf("required column '%s' not found in CSV header", col)
		}
	}

	// Read all records and convert to cameras
	var cameras []models.Camera
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading CSV record: %w", err)
		}

		camID := strings.TrimSpace(record[columnMap["cam_id"]])
		rtspURL := strings.TrimSpace(record[columnMap["rtsp"]])

		// Parse RTSP URL to extract camera details
		camera, err := parseRTSPURL(camID, rtspURL)
		if err != nil {
			return nil, fmt.Errorf("failed to parse RTSP URL for camera ID %s: %w", camID, err)
		}

		cameras = append(cameras, camera)
	}

	return cameras, nil
}

// parseRTSPURL parses an RTSP URL and extracts camera information
func parseRTSPURL(camID, rtspURL string) (models.Camera, error) {
	// Parse the URL
	parsedURL, err := url.Parse(rtspURL)
	if err != nil {
		return models.Camera{}, fmt.Errorf("invalid RTSP URL format: %w", err)
	}

	// Extract IP and port from host
	host := parsedURL.Hostname()
	if host == "" {
		return models.Camera{}, fmt.Errorf("no hostname found in RTSP URL")
	}

	// Extract username and password from userinfo
	var username, password string
	if parsedURL.User != nil {
		username = parsedURL.User.Username()
		password, _ = parsedURL.User.Password()
	}

	// Determine if camera is fake based on common patterns or default to false
	isFake := false
	// You can add logic here to determine if a camera is fake based on IP patterns
	// For example: if strings.HasPrefix(host, "192.168.1.15") { isFake = true }

	camera := models.Camera{
		ID:       camID,
		IP:       host,
		Port:     0,
		URL:      "",
		Username: username,
		Password: password,
		IsFake:   isFake,
	}

	return camera, nil
}
