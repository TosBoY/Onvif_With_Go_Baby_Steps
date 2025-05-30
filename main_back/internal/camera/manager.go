package camera

import (
	"fmt"
	"main_back/pkg/models"
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
