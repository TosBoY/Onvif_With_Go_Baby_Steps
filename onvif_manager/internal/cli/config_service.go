package cli

import (
	"fmt"
	"time"
)

// In-memory default config
var inMemoryConfig = SavedConfig{
	Width:       1920,
	Height:      1080,
	FPS:         25,
	Bitrate:     4096,
	Encoding:    "H264",
	Source:      "default",
	LastUpdated: time.Now().Format("2006-01-02 15:04:05"),
}

// ConfigService handles saved configuration operations
type ConfigService struct{}

// NewConfigService creates a new config service instance
func NewConfigService() *ConfigService {
	return &ConfigService{}
}

// LoadSavedConfig returns the current in-memory config
func (cs *ConfigService) LoadSavedConfig() (*SavedConfig, error) {
	return cs.GetDefaultConfig(), nil
}

// SaveConfig is a no-op in stateless mode
func (cs *ConfigService) SaveConfig(config *SavedConfig) error {
	// Update the in-memory config instead of saving to file
	inMemoryConfig = *config
	inMemoryConfig.LastUpdated = time.Now().Format("2006-01-02 15:04:05")
	return nil
}

// ImportFromConfigData updates the in-memory config
func (cs *ConfigService) ImportFromConfigData(configData *ConfigData, source string) error {
	inMemoryConfig = SavedConfig{
		Width:       configData.Width,
		Height:      configData.Height,
		FPS:         configData.FPS,
		Bitrate:     configData.Bitrate,
		Encoding:    configData.Encoding,
		Source:      source,
		LastUpdated: time.Now().Format("2006-01-02 15:04:05"),
	}
	return nil
}

// UpdateManually updates the saved config with manual values
func (cs *ConfigService) UpdateManually(width, height, fps, bitrate int, encoding string) error {
	savedConfig := &SavedConfig{
		Width:    width,
		Height:   height,
		FPS:      fps,
		Bitrate:  bitrate,
		Encoding: encoding,
		Source:   "manual",
	}

	return cs.SaveConfig(savedConfig)
}

// GetDefaultConfig returns the current in-memory configuration
func (cs *ConfigService) GetDefaultConfig() *SavedConfig {
	// Return a copy of the in-memory config to avoid external modifications
	config := inMemoryConfig
	return &config
}

// ValidateConfig validates configuration values
func (cs *ConfigService) ValidateConfig(width, height, fps, bitrate int, encoding string) error {
	if width <= 0 {
		return fmt.Errorf("width must be greater than 0")
	}
	if height <= 0 {
		return fmt.Errorf("height must be greater than 0")
	}
	if fps <= 0 {
		return fmt.Errorf("fps must be greater than 0")
	}
	if bitrate < 0 {
		return fmt.Errorf("bitrate must be 0 or greater")
	}

	// Validate encoding if provided
	if encoding != "" {
		validEncodings := []string{"H264", "H265", "MJPEG", "HEVC"}
		isValidEncoding := false
		for _, validEnc := range validEncodings {
			if encoding == validEnc {
				isValidEncoding = true
				break
			}
		}
		if !isValidEncoding {
			return fmt.Errorf("encoding '%s' is not supported. Valid options: H264, H265, HEVC, MJPEG", encoding)
		}
	}

	// Additional reasonable limits
	if width > 7680 { // 8K width
		return fmt.Errorf("width %d seems too large (max recommended: 7680)", width)
	}
	if height > 4320 { // 8K height
		return fmt.Errorf("height %d seems too large (max recommended: 4320)", height)
	}
	if fps > 120 {
		return fmt.Errorf("fps %d seems too high (max recommended: 120)", fps)
	}
	if bitrate > 50000 { // 50 Mbps
		return fmt.Errorf("bitrate %d kbps seems too high (max recommended: 50000)", bitrate)
	}

	return nil
}
