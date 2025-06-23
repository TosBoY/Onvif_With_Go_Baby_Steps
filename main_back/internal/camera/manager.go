package camera

import (
	"encoding/csv"
	"fmt"
	"main_back/internal/config"
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

// AddNewCamera adds a new camera to the cameras.csv file and assigns it an ID
// that is one greater than the largest existing ID.
// Returns the new camera ID and any error encountered.
func AddNewCamera(ip string, port int, url string, username string, password string, isFake bool) (string, error) {
	// Get the path to cameras.csv
	workingDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get working directory: %w", err)
	}

	configPath := filepath.Join(workingDir, "..", "..", "config", "cameras.csv")

	// Load the existing cameras using the config package
	cameras, err := config.LoadCameraList()
	if err != nil {
		return "", fmt.Errorf("failed to load existing cameras: %w", err)
	}

	// Check if a camera with the same IP already exists
	for _, cam := range cameras {
		if cam.IP == ip {
			return "", fmt.Errorf("camera with IP address %s already exists (ID: %s)", ip, cam.ID)
		}
	}

	// Find the highest ID
	highestID := 0
	for _, cam := range cameras {
		// Try to convert the ID to an integer
		if camID, err := strconv.Atoi(cam.ID); err == nil {
			if camID > highestID {
				highestID = camID
			}
		}
	}

	// Create a new ID by incrementing the highest ID
	newID := strconv.Itoa(highestID + 1)

	// Create a new camera
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
	cameras = append(cameras, newCamera)

	// Save the updated list
	err = saveCameraListToCSV(cameras, configPath)
	if err != nil {
		return "", fmt.Errorf("failed to save updated camera list: %w", err)
	}

	// Return the new camera ID
	return newID, nil
}

// saveCameraListToCSV saves the list of cameras to a CSV file
func saveCameraListToCSV(cameras []models.Camera, filePath string) error {
	// Create or open the file
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create or open camera CSV file: %w", err)
	}
	defer file.Close()

	// Create a CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write the header row
	header := []string{"id", "ip", "port", "url", "username", "password", "isFake"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write each camera record
	for _, cam := range cameras {
		record := []string{
			cam.ID,
			cam.IP,
			strconv.Itoa(cam.Port),
			cam.URL,
			cam.Username,
			cam.Password,
			strconv.FormatBool(cam.IsFake),
		}

		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write camera record: %w", err)
		}
	}

	return nil
}

// RemoveCamera removes a camera from the cameras.csv file by its ID.
// It also removes the camera client from the connected cameras map if it exists.
// Returns any error encountered.
func RemoveCamera(id string) error {
	// Get the path to cameras.csv
	workingDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	configPath := filepath.Join(workingDir, "..", "..", "config", "cameras.csv")

	// Load the existing cameras using the config package
	cameras, err := config.LoadCameraList()
	if err != nil {
		return fmt.Errorf("failed to load existing cameras: %w", err)
	}

	// Find the camera and remove it
	found := false
	updatedCameras := make([]models.Camera, 0, len(cameras))
	for _, cam := range cameras {
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

	// Save the updated list
	err = saveCameraListToCSV(updatedCameras, configPath)
	if err != nil {
		return fmt.Errorf("failed to save updated camera list: %w", err)
	}

	return nil
}
