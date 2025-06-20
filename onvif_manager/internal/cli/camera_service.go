package cli

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"onvif_manager/internal/backend/camera"
	"onvif_manager/internal/backend/config"
	"onvif_manager/internal/backend/ffmpeg"
	"onvif_manager/pkg/models"
)

// CameraService handles camera-related CLI operations
type CameraService struct{}

// NewCameraService creates a new camera service instance
func NewCameraService() *CameraService {
	return &CameraService{}
}

// EnsureCamerasInitialized ensures that all cameras are initialized and connected
func (cs *CameraService) EnsureCamerasInitialized() error {
	cameras, err := config.LoadCameraList()
	if err != nil {
		return fmt.Errorf("failed to load camera list: %w", err)
	}

	// Initialize all cameras (this will populate the connectedCameras map)
	if err := camera.InitializeAllCameras(cameras); err != nil {
		// Don't fail completely if some cameras can't be initialized
		// The individual operations will handle missing cameras gracefully
		log.Printf("Warning: Some cameras failed to initialize: %v", err)
	}

	return nil
}

// GetCameraList returns all cameras in the system
func (cs *CameraService) GetCameraList() ([]models.Camera, error) {
	return config.LoadCameraList()
}

// ImportCamerasFromCSV imports cameras from a CSV file
func (cs *CameraService) ImportCamerasFromCSV(csvFilePath string) (*ImportResult, error) {
	file, err := os.Open(csvFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	csvReader.FieldsPerRecord = -1

	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV file: %w", err)
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("CSV file is empty")
	}

	return cs.processCameraRecords(records)
}

// SelectCamerasFromCSV selects cameras based on IPs in CSV file
func (cs *CameraService) SelectCamerasFromCSV(csvFilePath string) (*SelectionResult, error) {
	file, err := os.Open(csvFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	csvReader.FieldsPerRecord = -1

	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV file: %w", err)
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("CSV file is empty")
	}

	return cs.processCameraSelection(records)
}

// ImportConfigFromCSV imports configuration from CSV file
func (cs *CameraService) ImportConfigFromCSV(csvFilePath string) (*ConfigData, error) {
	file, err := os.Open(csvFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	csvReader.FieldsPerRecord = -1

	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV file: %w", err)
	}

	if len(records) < 2 {
		return nil, fmt.Errorf("CSV file must contain header and configuration data")
	}

	return cs.processConfigData(records)
}

// ImportConfigFromCSVAndSave imports configuration from CSV file and saves it as saved config
func (cs *CameraService) ImportConfigFromCSVAndSave(csvFilePath string) (*ConfigData, error) {
	configData, err := cs.ImportConfigFromCSV(csvFilePath)
	if err != nil {
		return nil, err
	}

	// Save to saved config
	configService := NewConfigService()
	if err := configService.ImportFromConfigData(configData, "csv"); err != nil {
		log.Printf("Warning: Failed to save config to saved config file: %v", err)
		// Don't fail the entire operation if saving fails
	}

	return configData, nil
}

// ApplyConfigToCameras applies configuration to selected cameras
func (cs *CameraService) ApplyConfigToCameras(cameraIDs []string, config *ConfigData) (*ValidationResults, error) {
	log.Printf("Applying configuration to %d cameras", len(cameraIDs))

	results := &ValidationResults{
		CameraResults:     make(map[string]*CameraResult),
		ValidationResults: make(map[string]*ValidationResult),
	}

	// Phase 1: Apply configuration
	for _, cameraID := range cameraIDs {
		result := cs.applyCameraConfig(cameraID, config)
		results.CameraResults[cameraID] = result
	}

	// Wait for configurations to stabilize
	time.Sleep(1 * time.Second)

	// Phase 2: Validate configurations
	for _, cameraID := range cameraIDs {
		if result, exists := results.CameraResults[cameraID]; exists && result.Success && !result.IsFake {
			validation := cs.validateCameraStream(result.StreamURL, config)
			results.ValidationResults[cameraID] = validation
		} else if result.IsFake {
			// Add simulated validation for fake cameras
			validation := &ValidationResult{
				IsValid:         true,
				ExpectedWidth:   config.Width,
				ExpectedHeight:  config.Height,
				ExpectedFPS:     config.FPS,
				ExpectedBitrate: config.Bitrate,
				ActualWidth:     config.Width,
				ActualHeight:    config.Height,
				ActualFPS:       float64(config.FPS),
				ActualBitrate:   config.Bitrate,
				Message:         "Simulated camera configuration applied successfully",
			}
			results.ValidationResults[cameraID] = validation
		}
	}

	return results, nil
}

// ApplyConfigToCamerasFromSaved applies saved configuration to selected cameras
func (cs *CameraService) ApplyConfigToCamerasFromSaved(cameraIDs []string) (*ValidationResults, error) {
	configService := NewConfigService()
	savedConfig, err := configService.LoadSavedConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load saved config: %w", err)
	}

	configData := savedConfig.ToConfigData()
	return cs.ApplyConfigToCameras(cameraIDs, configData)
}

// ExportValidationToCSV exports validation results to CSV file
func (cs *CameraService) ExportValidationToCSV(validation *ValidationResults, outputPath string) error {
	cameras, err := cs.GetCameraList()
	if err != nil {
		return fmt.Errorf("failed to load camera list: %w", err)
	}

	// Create camera IP map
	cameraMap := make(map[string]models.Camera)
	for _, camera := range cameras {
		cameraMap[camera.ID] = camera
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{"cam_id", "cam_ip", "result", "reso_expected", "reso_actual", "fps_expected", "fps_actual"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write validation results
	for cameraID, validationResult := range validation.ValidationResults {
		cameraIP := "Unknown"
		if camera, exists := cameraMap[cameraID]; exists {
			cameraIP = camera.IP
		}

		result := "FAIL"
		if validationResult.IsValid {
			result = "PASS"
		}

		resoExpected := fmt.Sprintf("%dx%d", validationResult.ExpectedWidth, validationResult.ExpectedHeight)
		resoActual := ""
		if validationResult.ActualWidth > 0 && validationResult.ActualHeight > 0 {
			resoActual = fmt.Sprintf("%dx%d", validationResult.ActualWidth, validationResult.ActualHeight)
		}

		fpsExpected := strconv.Itoa(validationResult.ExpectedFPS)
		fpsActual := ""
		if validationResult.ActualFPS > 0 {
			fpsActual = fmt.Sprintf("%.2f", validationResult.ActualFPS)
		}

		row := []string{cameraID, cameraIP, result, resoExpected, resoActual, fpsExpected, fpsActual}
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write CSV row: %w", err)
		}
	}

	return nil
}

// Helper methods for processing data
func (cs *CameraService) processCameraRecords(records [][]string) (*ImportResult, error) {
	// Parse header to determine column indices
	headerRow := records[0]
	columnIndices := make(map[string]int)

	for i, column := range headerRow {
		columnName := strings.ToLower(strings.TrimSpace(column))
		columnIndices[columnName] = i
	}

	// Required columns
	requiredColumns := []string{"ip", "username"}
	for _, reqCol := range requiredColumns {
		if _, exists := columnIndices[reqCol]; !exists {
			return nil, fmt.Errorf("required column '%s' not found in CSV header", reqCol)
		}
	}

	// Process each data row
	var results []ImportRowResult
	var successCount, errorCount int

	for rowIndex, record := range records[1:] { // Skip header row
		rowNum := rowIndex + 2 // +2 because we start from row 1 (skipping header) and want 1-based numbering

		// Extract camera data with defaults
		cameraData := struct {
			IP       string
			Port     int
			URL      string
			Username string
			Password string
			IsFake   bool
		}{
			Port:   80,    // Default ONVIF port
			URL:    "",    // Default empty
			IsFake: false, // Default to real camera
		}

		// Extract IP (required)
		if ipIndex, exists := columnIndices["ip"]; exists && ipIndex < len(record) {
			cameraData.IP = strings.TrimSpace(record[ipIndex])
		}
		if cameraData.IP == "" {
			results = append(results, ImportRowResult{
				Row:     rowNum,
				Success: false,
				Error:   "Missing IP address",
				Data:    record,
			})
			errorCount++
			continue
		}

		// Extract Username (required)
		if usernameIndex, exists := columnIndices["username"]; exists && usernameIndex < len(record) {
			cameraData.Username = strings.TrimSpace(record[usernameIndex])
		}
		if cameraData.Username == "" {
			results = append(results, ImportRowResult{
				Row:     rowNum,
				Success: false,
				Error:   "Missing username",
				Data:    record,
			})
			errorCount++
			continue
		}

		// Extract optional fields
		if portIndex, exists := columnIndices["port"]; exists && portIndex < len(record) {
			if portStr := strings.TrimSpace(record[portIndex]); portStr != "" {
				if port, err := strconv.Atoi(portStr); err == nil {
					cameraData.Port = port
				}
			}
		}

		if urlIndex, exists := columnIndices["url"]; exists && urlIndex < len(record) {
			cameraData.URL = strings.TrimSpace(record[urlIndex])
		}

		if passwordIndex, exists := columnIndices["password"]; exists && passwordIndex < len(record) {
			cameraData.Password = strings.TrimSpace(record[passwordIndex])
		}

		if fakeIndex, exists := columnIndices["isfake"]; exists && fakeIndex < len(record) {
			if fakeStr := strings.ToLower(strings.TrimSpace(record[fakeIndex])); fakeStr != "" {
				cameraData.IsFake = (fakeStr == "true" || fakeStr == "1" || fakeStr == "yes")
			}
		}

		// Attempt to add the camera
		newID, err := camera.AddNewCamera(cameraData.IP, cameraData.Port, cameraData.URL, cameraData.Username, cameraData.Password, cameraData.IsFake)
		if err != nil {
			results = append(results, ImportRowResult{
				Row:     rowNum,
				Success: false,
				Error:   err.Error(),
				Data:    record,
			})
			errorCount++
		} else {
			newCamera := &models.Camera{
				ID:       newID,
				IP:       cameraData.IP,
				Port:     cameraData.Port,
				URL:      cameraData.URL,
				Username: cameraData.Username,
				Password: cameraData.Password,
				IsFake:   cameraData.IsFake,
			}
			results = append(results, ImportRowResult{
				Row:      rowNum,
				Success:  true,
				CameraID: newID,
				Camera:   newCamera,
			})
			successCount++
		}
	}

	return &ImportResult{
		Message:      fmt.Sprintf("CSV import completed: %d cameras added successfully, %d errors", successCount, errorCount),
		TotalRows:    len(records) - 1, // Exclude header
		SuccessCount: successCount,
		ErrorCount:   errorCount,
		Results:      results,
	}, nil
}

func (cs *CameraService) processCameraSelection(records [][]string) (*SelectionResult, error) {
	// Parse header to find the IP column
	headerRow := records[0]
	ipColumnIndex := -1

	for i, column := range headerRow {
		columnName := strings.ToLower(strings.TrimSpace(column))
		if columnName == "ip" {
			ipColumnIndex = i
			break
		}
	}

	if ipColumnIndex == -1 {
		return nil, fmt.Errorf("required column 'ip' not found in CSV header")
	}

	// Load existing cameras to match IPs with camera IDs
	cameras, err := config.LoadCameraList()
	if err != nil {
		return nil, fmt.Errorf("failed to load camera list: %w", err)
	}

	// Create a map of IP to camera ID for quick lookup
	ipToCameraMap := make(map[string]string)
	for _, camera := range cameras {
		ipToCameraMap[camera.IP] = camera.ID
	}

	// Process each data row to extract IPs and find matching cameras
	var selectedCameraIDs []string
	var matchedCameras []models.Camera
	var unmatchedIPs []string
	var invalidRows []InvalidRowInfo

	for rowIndex, record := range records[1:] { // Skip header row
		rowNum := rowIndex + 2 // +2 because we start from row 1 (skipping header) and want 1-based numbering

		// Check if row has enough columns
		if ipColumnIndex >= len(record) {
			invalidRows = append(invalidRows, InvalidRowInfo{
				Row:   rowNum,
				Error: "Insufficient columns",
				Data:  record,
			})
			continue
		}

		ip := strings.TrimSpace(record[ipColumnIndex])
		if ip == "" {
			invalidRows = append(invalidRows, InvalidRowInfo{
				Row:   rowNum,
				Error: "Empty IP address",
				Data:  record,
			})
			continue
		}

		// Check if this IP exists in our camera list
		if cameraID, exists := ipToCameraMap[ip]; exists {
			// Find the full camera object
			for _, camera := range cameras {
				if camera.ID == cameraID {
					selectedCameraIDs = append(selectedCameraIDs, cameraID)
					matchedCameras = append(matchedCameras, camera)
					break
				}
			}
		} else {
			unmatchedIPs = append(unmatchedIPs, ip)
		}
	}

	return &SelectionResult{
		Message:           fmt.Sprintf("Camera selection completed: %d cameras selected", len(selectedCameraIDs)),
		TotalRows:         len(records) - 1, // Exclude header
		SelectedCameraIDs: selectedCameraIDs,
		SelectedCameras:   matchedCameras,
		MatchedCount:      len(selectedCameraIDs),
		UnmatchedIPs:      unmatchedIPs,
		UnmatchedCount:    len(unmatchedIPs),
		InvalidRows:       invalidRows,
		InvalidRowCount:   len(invalidRows),
	}, nil
}

func (cs *CameraService) processConfigData(records [][]string) (*ConfigData, error) {
	// Parse header to determine column indices
	headerRow := records[0]
	columnIndices := make(map[string]int)

	for i, column := range headerRow {
		columnName := strings.ToLower(strings.TrimSpace(column))
		columnIndices[columnName] = i
	}

	// Required columns (bitrate is optional)
	requiredColumns := []string{"width", "height", "fps"}
	for _, reqCol := range requiredColumns {
		if _, exists := columnIndices[reqCol]; !exists {
			return nil, fmt.Errorf("required column '%s' not found in CSV header", reqCol)
		}
	}

	// Process the first data row (should only be 1 row)
	dataRow := records[1]

	// Parse configuration values
	configData := &ConfigData{
		Bitrate: 0, // Default value for optional bitrate
	}

	// Extract Width (required)
	if widthIndex, exists := columnIndices["width"]; exists && widthIndex < len(dataRow) {
		widthStr := strings.TrimSpace(dataRow[widthIndex])
		if widthStr == "" {
			return nil, fmt.Errorf("width value is required")
		}
		width, err := strconv.Atoi(widthStr)
		if err != nil || width <= 0 {
			return nil, fmt.Errorf("invalid width value: %s", widthStr)
		}
		configData.Width = width
	} else {
		return nil, fmt.Errorf("width value is required")
	}

	// Extract Height (required)
	if heightIndex, exists := columnIndices["height"]; exists && heightIndex < len(dataRow) {
		heightStr := strings.TrimSpace(dataRow[heightIndex])
		if heightStr == "" {
			return nil, fmt.Errorf("height value is required")
		}
		height, err := strconv.Atoi(heightStr)
		if err != nil || height <= 0 {
			return nil, fmt.Errorf("invalid height value: %s", heightStr)
		}
		configData.Height = height
	} else {
		return nil, fmt.Errorf("height value is required")
	}

	// Extract FPS (required)
	if fpsIndex, exists := columnIndices["fps"]; exists && fpsIndex < len(dataRow) {
		fpsStr := strings.TrimSpace(dataRow[fpsIndex])
		if fpsStr == "" {
			return nil, fmt.Errorf("FPS value is required")
		}
		fps, err := strconv.Atoi(fpsStr)
		if err != nil || fps <= 0 {
			return nil, fmt.Errorf("invalid FPS value: %s", fpsStr)
		}
		configData.FPS = fps
	} else {
		return nil, fmt.Errorf("FPS value is required")
	}

	// Extract Bitrate (optional)
	if bitrateIndex, exists := columnIndices["bitrate"]; exists && bitrateIndex < len(dataRow) {
		bitrateStr := strings.TrimSpace(dataRow[bitrateIndex])
		if bitrateStr != "" {
			bitrate, err := strconv.Atoi(bitrateStr)
			if err != nil || bitrate < 0 {
				log.Printf("Warning: Invalid bitrate value '%s', using default 0", bitrateStr)
			} else {
				configData.Bitrate = bitrate
			}
		}
	}

	return configData, nil
}

func (cs *CameraService) applyCameraConfig(cameraID string, config *ConfigData) *CameraResult {
	result := &CameraResult{
		CameraID: cameraID,
		Success:  false,
	}

	// Get the camera client
	client, err := camera.GetCameraClient(cameraID)
	if err != nil {
		result.Error = fmt.Errorf("camera not found: %w", err)
		return result
	}

	// Check if camera is fake and handle it differently
	if client.Camera.IsFake {
		log.Printf("Camera %s is a simulated device, skipping real configuration", cameraID)
		// For fake cameras, we simulate successful config application
		result.Success = true
		result.IsFake = true
		result.AppliedConfig = map[string]interface{}{
			"resolution": map[string]int{
				"width":  config.Width,
				"height": config.Height,
			},
			"fps":     config.FPS,
			"bitrate": config.Bitrate,
		}
		return result
	}

	// For real cameras, proceed with config application
	log.Printf("Getting profiles and configs for camera %s (IP: %s:%d)", cameraID, client.Camera.IP, client.Camera.Port)
	profileTokens, configTokens, err := camera.GetProfilesAndConfigs(client)
	if err != nil {
		log.Printf("Failed to get camera profiles and configs for %s (IP: %s:%d): %v", cameraID, client.Camera.IP, client.Camera.Port, err)
		// Add more specific error information for network issues
		errorMsg := err.Error()
		if strings.Contains(errorMsg, "i/o timeout") || strings.Contains(errorMsg, "dial tcp") {
			result.Error = fmt.Errorf("network timeout: camera at %s:%d is not responding. Please check: 1) Camera is powered on and connected to network, 2) IP address %s is correct, 3) Port %d is the correct ONVIF port, 4) Camera supports ONVIF protocol", client.Camera.IP, client.Camera.Port, client.Camera.IP, client.Camera.Port)
		} else if strings.Contains(errorMsg, "connection refused") {
			result.Error = fmt.Errorf("connection refused: camera at %s:%d refused connection. Please check: 1) Correct ONVIF port (common ports: 80, 8080, 554), 2) ONVIF service is enabled on camera, 3) Firewall settings", client.Camera.IP, client.Camera.Port)
		} else if strings.Contains(errorMsg, "no route to host") {
			result.Error = fmt.Errorf("no route to host: cannot reach camera at %s:%d. Please check: 1) Camera and server are on same network, 2) IP address is correct, 3) Network routing", client.Camera.IP, client.Camera.Port)
		} else {
			result.Error = fmt.Errorf("failed to get camera profiles and configs: %w", err)
		}
		return result
	}

	if len(profileTokens) == 0 {
		log.Printf("No profiles found for camera %s", cameraID)
		result.Error = fmt.Errorf("no profiles found")
		return result
	}

	if len(configTokens) == 0 {
		log.Printf("No video encoder configuration found for camera %s", cameraID)
		result.Error = fmt.Errorf("no video encoder configuration found")
		return result
	}

	// Use the first token found
	profileToken := profileTokens[0]
	configToken := configTokens[0]
	result.ProfileToken = profileToken

	log.Printf("Using profile token %s and config token %s for camera %s", profileToken, configToken, cameraID)

	// Get current encoder config
	log.Printf("Getting current encoder config for camera %s", cameraID)
	currentConfig, err := camera.GetCurrentConfig(client, configToken)
	if err != nil {
		log.Printf("Failed to get current encoder config for %s: %v", cameraID, err)
		result.Error = fmt.Errorf("failed to get current encoder config: %w", err)
		return result
	}

	// Get available encoder options
	log.Printf("Getting available encoder options for camera %s", cameraID)
	encoderOptions, err := camera.GetCurrentEncoderOptions(client, profileToken, configToken)
	if err != nil {
		log.Printf("Failed to get encoder options for %s: %v", cameraID, err)
		result.Error = fmt.Errorf("failed to get encoder options: %w", err)
		return result
	}

	// Create target resolution object
	targetResolution := models.Resolution{Width: config.Width, Height: config.Height}

	// Find closest matching resolution
	log.Printf("Finding closest matching resolution for camera %s", cameraID)
	closestResolution := camera.FindClosestResolution(targetResolution, encoderOptions.Resolutions)
	log.Printf("Closest resolution found for camera %s: %dx%d", cameraID, closestResolution.Width, closestResolution.Height)

	// Check if current configuration already matches the requested configuration
	currentMatches := currentConfig.Resolution.Width == closestResolution.Width &&
		currentConfig.Resolution.Height == closestResolution.Height &&
		currentConfig.FPS == config.FPS &&
		(config.Bitrate == 0 || currentConfig.Bitrate == config.Bitrate)

	if currentMatches {
		log.Printf("Camera %s already has the requested configuration (Resolution: %dx%d, FPS: %d, Bitrate: %d), skipping config change",
			cameraID, closestResolution.Width, closestResolution.Height, config.FPS, currentConfig.Bitrate)
		// Mark as successful but indicate no change was needed
		result.Success = true
		result.AppliedConfig = map[string]interface{}{
			"resolution": map[string]int{
				"width":  closestResolution.Width,
				"height": closestResolution.Height,
			},
			"fps":       config.FPS,
			"bitrate":   currentConfig.Bitrate,
			"unchanged": true, // Indicate no change was needed
		}
		result.ResolutionAdjusted = config.Width != closestResolution.Width || config.Height != closestResolution.Height

		// Still get stream URI for validation
		streamURI, err := client.GetStreamURI(profileToken)
		if err != nil {
			log.Printf("Failed to get stream URI for %s: %v", cameraID, err)
			result.Error = fmt.Errorf("failed to get stream URI: %w", err)
			return result
		}

		// Parse and construct the URL with embedded credentials
		parsedURI, err := url.Parse(streamURI)
		if err != nil {
			log.Printf("Failed to parse stream URI for %s: %v", cameraID, err)
			result.Error = fmt.Errorf("failed to parse stream URI: %w", err)
			return result
		}

		fullStreamURL := fmt.Sprintf("%s://%s:%s@%s%s", parsedURI.Scheme, client.Camera.Username, client.Camera.Password, parsedURI.Host, parsedURI.RequestURI())
		result.StreamURL = fullStreamURL

		return result
	}

	// Prepare the new configuration
	newConfig := models.EncoderConfig{
		Resolution: closestResolution,
		Quality:    currentConfig.Quality, // Keep the current quality
		FPS:        config.FPS,
		Bitrate:    config.Bitrate,
	}
	log.Printf("Prepared new config for camera %s: %+v", cameraID, newConfig)

	// Set the new encoder config
	log.Printf("Setting new encoder config for camera %s", cameraID)
	if err := camera.SetEncoderConfig(client, configToken, currentConfig, newConfig); err != nil {
		log.Printf("Failed to set encoder config for %s: %v", cameraID, err)
		result.Error = fmt.Errorf("failed to set encoder config: %w", err)
		return result
	}
	log.Printf("Successfully applied config for camera %s", cameraID)

	// Get stream URI for later validation
	streamURI, err := client.GetStreamURI(profileToken)
	if err != nil {
		log.Printf("Failed to get stream URI for %s: %v", cameraID, err)
		result.Error = fmt.Errorf("failed to get stream URI: %w", err)
		return result
	}

	// The ONVIF GetStreamUri typically doesn't include credentials. Embed them for FFmpeg validation.
	parsedURI, err := url.Parse(streamURI)
	if err != nil {
		log.Printf("Failed to parse stream URI for %s: %v", cameraID, err)
		result.Error = fmt.Errorf("failed to parse stream URI: %w", err)
		return result
	}

	// Construct the URL with embedded credentials
	fullStreamURL := fmt.Sprintf("%s://%s:%s@%s%s", parsedURI.Scheme, client.Camera.Username, client.Camera.Password, parsedURI.Host, parsedURI.RequestURI())

	// Mark this camera as successfully configured
	result.Success = true
	result.AppliedConfig = map[string]interface{}{
		"resolution": map[string]int{
			"width":  closestResolution.Width,
			"height": closestResolution.Height,
		},
		"fps":     config.FPS,
		"bitrate": config.Bitrate,
	}
	result.ResolutionAdjusted = config.Width != closestResolution.Width || config.Height != closestResolution.Height
	result.StreamURL = fullStreamURL

	return result
}

func (cs *CameraService) validateCameraStream(streamURL string, config *ConfigData) *ValidationResult {
	validationResult, err := ffmpeg.ValidateStream(streamURL, config.Width, config.Height, config.FPS, config.Bitrate)
	if err != nil {
		return &ValidationResult{
			IsValid:         false,
			Error:           err.Error(),
			ExpectedWidth:   config.Width,
			ExpectedHeight:  config.Height,
			ExpectedFPS:     config.FPS,
			ExpectedBitrate: config.Bitrate,
		}
	}

	// Determine validation status based on business rules
	// Resolution mismatch = failure, FPS/bitrate mismatch = warning
	resolutionMatches := validationResult.ActualWidth > 0 && validationResult.ActualHeight > 0 &&
		validationResult.ActualWidth == validationResult.ExpectedWidth &&
		validationResult.ActualHeight == validationResult.ExpectedHeight

	// Override validation result: resolution mismatch = failure, others = warning
	overrideIsValid := resolutionMatches // Only consider valid if resolution matches

	return &ValidationResult{
		IsValid:         overrideIsValid, // Use our override logic
		ExpectedWidth:   validationResult.ExpectedWidth,
		ExpectedHeight:  validationResult.ExpectedHeight,
		ExpectedFPS:     validationResult.ExpectedFPS,
		ExpectedBitrate: validationResult.ExpectedBitrate,
		ActualWidth:     validationResult.ActualWidth,
		ActualHeight:    validationResult.ActualHeight,
		ActualFPS:       validationResult.ActualFPS,
		ActualBitrate:   validationResult.ActualBitrate,
		Error:           validationResult.Error,
	}
}
