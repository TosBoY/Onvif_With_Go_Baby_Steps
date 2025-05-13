package validator

import (
	"fmt"
	"time"

	lib "onvif_test2/lib"
)

type ValidationResult struct {
	IsValid        bool    `json:"isValid"`
	Message        string  `json:"message"`
	ActualWidth    int     `json:"actualWidth"`
	ActualHeight   int     `json:"actualHeight"`
	ActualFPS      float64 `json:"actualFPS"`
	ExpectedWidth  int     `json:"expectedWidth"`
	ExpectedHeight int     `json:"expectedHeight"`
	ExpectedFPS    float64 `json:"expectedFPS"`
	Error          string  `json:"error,omitempty"`
}

// ValidateVideoConfig checks if the actual stream parameters match the expected configuration
func ValidateVideoConfig(streamURL string, expectedConfig VideoConfig) (*ValidationResult, error) {
	// Give the camera a moment to apply the new configuration
	time.Sleep(2 * time.Second)

	// Try up to 3 times with increasing delays
	for attempt := 1; attempt <= 3; attempt++ {
		fmt.Printf("Validation attempt %d for stream %s\n", attempt, streamURL)

		valid, err := lib.ValidateStreamConfig(
			streamURL,
			expectedConfig.Width,
			expectedConfig.Height,
			float64(expectedConfig.FrameRate),
		)

		// Prepare base result with expected values
		result := &ValidationResult{
			ExpectedWidth:  expectedConfig.Width,
			ExpectedHeight: expectedConfig.Height,
			ExpectedFPS:    float64(expectedConfig.FrameRate),
		}

		if err == nil && valid {
			result.IsValid = true
			result.Message = "Stream configuration matches expected values"
			result.ActualWidth = expectedConfig.Width
			result.ActualHeight = expectedConfig.Height
			result.ActualFPS = float64(expectedConfig.FrameRate)
			return result, nil
		}

		if err != nil {
			fmt.Printf("Validation error on attempt %d: %v\n", attempt, err)

			result.IsValid = false
			result.Message = fmt.Sprintf("Configuration mismatch on attempt %d", attempt)
			result.Error = err.Error()

			errStr := err.Error()

			// Look for various error formats in the error string
			if _, err2 := fmt.Sscanf(errStr, "resolution mismatch: got %dx%d", &result.ActualWidth, &result.ActualHeight); err2 != nil {
				fmt.Printf("Failed to parse actual resolution: %v\n", err2)
			}
			if _, err2 := fmt.Sscanf(errStr, "frame rate mismatch: got %f", &result.ActualFPS); err2 != nil {
				fmt.Printf("Failed to parse actual frame rate: %v\n", err2)
			}

			// If this is the last attempt, return the result
			if attempt == 3 {
				fmt.Printf("Validation failed after 3 attempts: %v\n", err)
				return result, nil
			}

			// Wait longer before the next attempt
			delay := time.Duration(attempt) * 2 * time.Second
			fmt.Printf("Waiting %v before next attempt\n", delay)
			time.Sleep(delay)
			continue
		}
	}

	return &ValidationResult{
		IsValid:        false,
		Message:        "Failed to validate stream configuration after 3 attempts",
		ExpectedWidth:  expectedConfig.Width,
		ExpectedHeight: expectedConfig.Height,
		ExpectedFPS:    float64(expectedConfig.FrameRate),
	}, nil
}

type VideoConfig struct {
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	FrameRate   int    `json:"frameRate"`
	BitRate     int    `json:"bitRate"`
	GopLength   int    `json:"gopLength"`
	H264Profile string `json:"h264Profile"`
}
