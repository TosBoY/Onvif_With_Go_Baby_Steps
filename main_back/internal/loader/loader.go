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

// LoadCameraList reads a CSV file containing camera configurations with cam_id and rtsp columns.
// It determines the path to the config file relative to the executable's location.
func LoadCameraList() ([]models.Camera, error) {
	workingDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get working directory: %w", err)
	}

	configPath := filepath.Join(workingDir, "..", "..", "internal", "loader", "cameras.csv")

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
		Port:     0, // Use ONVIF port instead of RTSP port
		URL:      "",
		Username: username,
		Password: password,
		IsFake:   isFake,
	}

	return camera, nil
}
