package camera

import (
	"encoding/json"
	"main_back/pkg/models"
	"os"
	"path/filepath"
	"testing"
)

func TestAddNewCamera(t *testing.T) {
	// Create a temporary file for testing
	tempDir, err := os.MkdirTemp("", "camera-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create config directory in temp directory
	configDir := filepath.Join(tempDir, "config")
	if err := os.Mkdir(configDir, 0755); err != nil {
		t.Fatalf("Failed to create config directory: %v", err)
	}

	// Create a test cameras.json file
	testCameras := struct {
		Cameras []models.Camera `json:"cameras"`
	}{
		Cameras: []models.Camera{
			{
				ID:       "1",
				IP:       "192.168.1.10",
				Username: "admin",
				Password: "password1",
				IsFake:   false,
			},
			{
				ID:       "5", // Non-sequential ID to test finding highest
				IP:       "192.168.1.20",
				Username: "admin",
				Password: "password2",
				IsFake:   false,
			},
		},
	}

	testData, err := json.MarshalIndent(testCameras, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal test data: %v", err)
	}

	testFilePath := filepath.Join(configDir, "cameras.json")
	if err := os.WriteFile(testFilePath, testData, 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Store original working directory
	origWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}
	defer os.Chdir(origWd) // Return to original working directory when done

	// Change to temp directory for test
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change working directory: %v", err)
	}
	// Test adding a new camera
	testIP := "192.168.1.30"
	testUsername := "admin"
	testPassword := "password3"
	testIsFake := false

	newID, err := AddNewCamera(testIP, testUsername, testPassword, testIsFake)
	if err != nil {
		t.Fatalf("AddNewCamera returned error: %v", err)
	}

	// Check that the new ID is "6" (one higher than the highest existing ID)
	expectedID := "6"
	if newID != expectedID {
		t.Errorf("Expected new ID to be %s, got %s", expectedID, newID)
	}

	// Read the file back and verify the new camera was added correctly
	updatedData, err := os.ReadFile(testFilePath)
	if err != nil {
		t.Fatalf("Failed to read updated config file: %v", err)
	}

	var updatedConfig struct {
		Cameras []models.Camera `json:"cameras"`
	}

	if err := json.Unmarshal(updatedData, &updatedConfig); err != nil {
		t.Fatalf("Failed to unmarshal updated config: %v", err)
	}

	// Check that we have 3 cameras now
	if len(updatedConfig.Cameras) != 3 {
		t.Errorf("Expected 3 cameras, got %d", len(updatedConfig.Cameras))
	}

	// Find the newly added camera
	var foundNewCamera bool
	for _, cam := range updatedConfig.Cameras {
		if cam.ID == newID {
			foundNewCamera = true
			// Verify all properties
			if cam.IP != testIP {
				t.Errorf("Expected IP %s, got %s", testIP, cam.IP)
			}
			if cam.Username != testUsername {
				t.Errorf("Expected Username %s, got %s", testUsername, cam.Username)
			}
			if cam.Password != testPassword {
				t.Errorf("Expected Password %s, got %s", testPassword, cam.Password)
			}
			if cam.IsFake != false {
				t.Errorf("Expected IsFake to be false, got %v", cam.IsFake)
			}
			break
		}
	}
	if !foundNewCamera {
		t.Errorf("New camera with ID %s not found in updated config", newID)
	}
}

func TestAddNewFakeCamera(t *testing.T) {
	// Create a temporary file for testing
	tempDir, err := os.MkdirTemp("", "fake-camera-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create config directory in temp directory
	configDir := filepath.Join(tempDir, "config")
	if err := os.Mkdir(configDir, 0755); err != nil {
		t.Fatalf("Failed to create config directory: %v", err)
	}

	// Create a test cameras.json file
	testCameras := struct {
		Cameras []models.Camera `json:"cameras"`
	}{
		Cameras: []models.Camera{
			{
				ID:       "1",
				IP:       "192.168.1.10",
				Username: "admin",
				Password: "password1",
				IsFake:   false,
			},
		},
	}

	testData, err := json.MarshalIndent(testCameras, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal test data: %v", err)
	}

	testFilePath := filepath.Join(configDir, "cameras.json")
	if err := os.WriteFile(testFilePath, testData, 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Store original working directory
	origWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}
	defer os.Chdir(origWd) // Return to original working directory when done

	// Change to temp directory for test
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change working directory: %v", err)
	}

	// Test adding a fake camera
	testIP := "192.168.1.200"
	testUsername := "admin"
	testPassword := "password"
	testIsFake := true

	newID, err := AddNewCamera(testIP, testUsername, testPassword, testIsFake)
	if err != nil {
		t.Fatalf("AddNewCamera returned error: %v", err)
	}

	// Check that the new ID is "2" (one higher than the highest existing ID)
	expectedID := "2"
	if newID != expectedID {
		t.Errorf("Expected new ID to be %s, got %s", expectedID, newID)
	}

	// Read the file back and verify the new camera was added correctly
	updatedData, err := os.ReadFile(testFilePath)
	if err != nil {
		t.Fatalf("Failed to read updated config file: %v", err)
	}

	var updatedConfig struct {
		Cameras []models.Camera `json:"cameras"`
	}

	if err := json.Unmarshal(updatedData, &updatedConfig); err != nil {
		t.Fatalf("Failed to unmarshal updated config: %v", err)
	}

	// Check that we have 2 cameras now
	if len(updatedConfig.Cameras) != 2 {
		t.Errorf("Expected 2 cameras, got %d", len(updatedConfig.Cameras))
	}

	// Find the newly added camera
	var foundNewCamera bool
	for _, cam := range updatedConfig.Cameras {
		if cam.ID == newID {
			foundNewCamera = true
			// Verify all properties
			if cam.IP != testIP {
				t.Errorf("Expected IP %s, got %s", testIP, cam.IP)
			}
			if cam.Username != testUsername {
				t.Errorf("Expected Username %s, got %s", testUsername, cam.Username)
			}
			if cam.Password != testPassword {
				t.Errorf("Expected Password %s, got %s", testPassword, cam.Password)
			}
			if cam.IsFake != testIsFake {
				t.Errorf("Expected IsFake to be %v, got %v", testIsFake, cam.IsFake)
			}
			break
		}
	}

	if !foundNewCamera {
		t.Errorf("New fake camera with ID %s not found in updated config", newID)
	}
}
