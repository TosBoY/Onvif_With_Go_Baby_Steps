package camera

import (
	"fmt"
	"onvif_manager/pkg/models"
	"strconv"
	"strings"
)

var connectedCameras map[string]*CameraClient
var inMemoryCameras []models.Camera

func init() {
	connectedCameras = make(map[string]*CameraClient)
	inMemoryCameras = []models.Camera{}
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

// GetAllCameras returns the list of all cameras from in-memory storage
func GetAllCameras() []models.Camera {
	return inMemoryCameras
}

// AddNewCamera adds a new camera to the in-memory storage and assigns it an ID
// that is one greater than the largest existing ID.
// Returns the new camera ID and any error encountered.
func AddNewCamera(ip string, port int, url string, username string, password string, isFake bool) (string, error) {
	// Check if a camera with the same IP already exists
	for _, cam := range inMemoryCameras {
		if cam.IP == ip {
			return "", fmt.Errorf("camera with IP address %s already exists (ID: %s)", ip, cam.ID)
		}
	}

	// Find the highest ID
	highestID := 0
	for _, cam := range inMemoryCameras {
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

	// Add the new camera to the in-memory list
	inMemoryCameras = append(inMemoryCameras, newCamera)

	// Initialize the camera client and add it to connectedCameras
	if isFake {
		// Create a fake camera client
		client := NewFakeCameraClient(newCamera)
		connectedCameras[newID] = client
	} else {
		// Attempt to create a real camera client
		client, err := NewCameraClient(newCamera)
		if err != nil {
			// Log the error but don't fail - we still want to add the camera to the list
			fmt.Printf("Warning: Failed to initialize camera client for %s: %v\n", newID, err)
		} else {
			// Add to connected cameras
			connectedCameras[newID] = client
		}
	}

	// Return the new camera ID
	return newID, nil
}

// RemoveCamera removes a camera from the in-memory storage by its ID.
// It also removes the camera client from the connected cameras map if it exists.
// Returns any error encountered.
func RemoveCamera(id string) error {
	// Find the camera and remove it
	var found bool
	var updatedCameras []models.Camera

	for _, cam := range inMemoryCameras {
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

	// Update the in-memory list
	inMemoryCameras = updatedCameras
	return nil
}
