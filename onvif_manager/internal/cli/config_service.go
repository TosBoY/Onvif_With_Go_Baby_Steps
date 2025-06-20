package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"onvif_manager/internal/backend/config"
)

// ConfigService handles saved configuration operations
type ConfigService struct{}

// NewConfigService creates a new config service instance
func NewConfigService() *ConfigService {
	return &ConfigService{}
}

// getSavedConfigPath returns the path to the saved config file
func (cs *ConfigService) getSavedConfigPath() (string, error) {
	configDir, err := config.FindConfigPath()
	if err != nil {
		return "", fmt.Errorf("failed to find config directory: %w", err)
	}

	// Get the directory containing cameras.json and use it for saved_config.json
	configDirPath := filepath.Dir(configDir)
	return filepath.Join(configDirPath, "saved_config.json"), nil
}

// LoadSavedConfig loads the saved configuration from file
func (cs *ConfigService) LoadSavedConfig() (*SavedConfig, error) {
	configPath, err := cs.getSavedConfigPath()
	if err != nil {
		return nil, err
	}

	// Check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Return default config if file doesn't exist
		return cs.GetDefaultConfig(), nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read saved config file: %w", err)
	}

	var savedConfig SavedConfig
	if err := json.Unmarshal(data, &savedConfig); err != nil {
		return nil, fmt.Errorf("failed to parse saved config: %w", err)
	}

	return &savedConfig, nil
}

// SaveConfig saves the configuration to file
func (cs *ConfigService) SaveConfig(config *SavedConfig) error {
	configPath, err := cs.getSavedConfigPath()
	if err != nil {
		return err
	}

	// Update timestamp and ensure directory exists
	config.LastUpdated = time.Now().Format("2006-01-02 15:04:05")

	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write saved config: %w", err)
	}

	return nil
}

// ImportFromConfigData saves a ConfigData as SavedConfig
func (cs *ConfigService) ImportFromConfigData(configData *ConfigData, source string) error {
	savedConfig := &SavedConfig{
		Width:   configData.Width,
		Height:  configData.Height,
		FPS:     configData.FPS,
		Bitrate: configData.Bitrate,
		Source:  source,
	}

	return cs.SaveConfig(savedConfig)
}

// UpdateManually updates the saved config with manual values
func (cs *ConfigService) UpdateManually(width, height, fps, bitrate int) error {
	savedConfig := &SavedConfig{
		Width:   width,
		Height:  height,
		FPS:     fps,
		Bitrate: bitrate,
		Source:  "manual",
	}

	return cs.SaveConfig(savedConfig)
}

// GetDefaultConfig returns a default configuration
func (cs *ConfigService) GetDefaultConfig() *SavedConfig {
	return &SavedConfig{
		Width:       1920,
		Height:      1080,
		FPS:         30,
		Bitrate:     4096,
		LastUpdated: time.Now().Format("2006-01-02 15:04:05"),
		Source:      "default",
	}
}

// ValidateConfig validates configuration values
func (cs *ConfigService) ValidateConfig(width, height, fps, bitrate int) error {
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
