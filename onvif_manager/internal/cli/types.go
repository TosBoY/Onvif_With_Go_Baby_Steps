package cli

import "onvif_manager/pkg/models"

// ImportResult represents the result of importing cameras from CSV
type ImportResult struct {
	Message      string            `json:"message"`
	TotalRows    int               `json:"totalRows"`
	SuccessCount int               `json:"successCount"`
	ErrorCount   int               `json:"errorCount"`
	Results      []ImportRowResult `json:"results"`
}

// ImportRowResult represents the result of importing a single camera row
type ImportRowResult struct {
	Row      int            `json:"row"`
	Success  bool           `json:"success"`
	Error    string         `json:"error,omitempty"`
	CameraID string         `json:"cameraId,omitempty"`
	Camera   *models.Camera `json:"camera,omitempty"`
	Data     []string       `json:"data,omitempty"`
}

// SelectionResult represents the result of selecting cameras from CSV
type SelectionResult struct {
	Message           string           `json:"message"`
	TotalRows         int              `json:"totalRows"`
	SelectedCameraIDs []string         `json:"selectedCameraIds"`
	SelectedCameras   []models.Camera  `json:"selectedCameras"`
	MatchedCount      int              `json:"matchedCount"`
	UnmatchedIPs      []string         `json:"unmatchedIPs"`
	UnmatchedCount    int              `json:"unmatchedCount"`
	InvalidRows       []InvalidRowInfo `json:"invalidRows"`
	InvalidRowCount   int              `json:"invalidRowCount"`
}

// InvalidRowInfo represents information about invalid rows
type InvalidRowInfo struct {
	Row   int      `json:"row"`
	Error string   `json:"error"`
	Data  []string `json:"data"`
}

// ConfigData represents configuration data imported from CSV
type ConfigData struct {
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	FPS      int    `json:"fps"`
	Bitrate  int    `json:"bitrate"`
	Encoding string `json:"encoding"`
}

// SavedConfig represents the persistent configuration stored in saved_config.json
type SavedConfig struct {
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	FPS         int    `json:"fps"`
	Bitrate     int    `json:"bitrate"`
	Encoding    string `json:"encoding"`
	LastUpdated string `json:"lastUpdated"`
	Source      string `json:"source"` // "csv", "manual", "default"
}

// ToConfigData converts SavedConfig to ConfigData for applying to cameras
func (sc *SavedConfig) ToConfigData() *ConfigData {
	return &ConfigData{
		Width:    sc.Width,
		Height:   sc.Height,
		FPS:      sc.FPS,
		Bitrate:  sc.Bitrate,
		Encoding: sc.Encoding,
	}
}

// ValidationResults represents the overall validation results
type ValidationResults struct {
	CameraResults     map[string]*CameraResult     `json:"cameraResults"`
	ValidationResults map[string]*ValidationResult `json:"validationResults"`
}

// CameraResult represents the result of applying configuration to a single camera
type CameraResult struct {
	CameraID           string                 `json:"cameraId"`
	Success            bool                   `json:"success"`
	Error              error                  `json:"error,omitempty"`
	AppliedConfig      map[string]interface{} `json:"appliedConfig,omitempty"`
	ResolutionAdjusted bool                   `json:"resolutionAdjusted"`
	ProfileToken       string                 `json:"profileToken,omitempty"`
	StreamURL          string                 `json:"streamUrl,omitempty"`
}

// ValidationResult represents the result of validating a camera stream
type ValidationResult struct {
	IsValid          bool    `json:"isValid"`
	ExpectedWidth    int     `json:"expectedWidth"`
	ExpectedHeight   int     `json:"expectedHeight"`
	ExpectedFPS      int     `json:"expectedFPS"`
	ExpectedBitrate  int     `json:"expectedBitrate"`
	ExpectedEncoding string  `json:"expectedEncoding"`
	ActualWidth      int     `json:"actualWidth"`
	ActualHeight     int     `json:"actualHeight"`
	ActualFPS        float64 `json:"actualFPS"`
	ActualBitrate    int     `json:"actualBitrate"`
	ActualEncoding   string  `json:"actualEncoding"`
	Error            string  `json:"error,omitempty"`
	Message          string  `json:"message,omitempty"`
}

// Summary represents a summary of operations
type Summary struct {
	TotalCameras     int `json:"totalCameras"`
	SuccessfulCams   int `json:"successfulCams"`
	FailedCams       int `json:"failedCams"`
	ValidatedCams    int `json:"validatedCams"`
	PassedValidation int `json:"passedValidation"`
	FailedValidation int `json:"failedValidation"`
}
