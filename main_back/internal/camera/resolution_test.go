package camera

import (
    "testing"
    "main_back/internal/models"
)

func TestFindClosestResolution_Camera2(t *testing.T) {
    // Define the target resolution
    targetResolution := models.Resolution{
        Width:  1280,
        Height: 720,
    }

    // Define the available resolutions for camera 2
    availableResolutions := []models.Resolution{
        {Width: 640, Height: 480},
        {Width: 1280, Height: 720},
        {Width: 1920, Height: 1080},
    }

    // Call the function
    closestResolution := FindClosestResolution(targetResolution, availableResolutions)

    // Assert the result
    if closestResolution.Width != 1280 || closestResolution.Height != 720 {
        t.Errorf("Expected resolution 1280x720, but got %dx%d", closestResolution.Width, closestResolution.Height)
    }
}