package validator

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"time"

	lib "onvif_back/lib"
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

		// Prepare base result with expected values
		result := &ValidationResult{
			ExpectedWidth:  expectedConfig.Width,
			ExpectedHeight: expectedConfig.Height,
			ExpectedFPS:    float64(expectedConfig.FrameRate),
		}

		// Get the actual stream configuration directly from ffprobe
		ffprobeResult, err := lib.GetStreamInfo(streamURL)
		if err == nil && ffprobeResult != nil {
			// Found a valid result from ffprobe, use it directly
			result.ActualWidth = ffprobeResult.Width
			result.ActualHeight = ffprobeResult.Height
			result.ActualFPS = ffprobeResult.FrameRate

			// Check if values match with some tolerance for FPS
			resolutionMatches := result.ActualWidth == result.ExpectedWidth &&
				result.ActualHeight == result.ExpectedHeight
			fpsMatches := math.Abs(result.ActualFPS-result.ExpectedFPS) < 0.1

			result.IsValid = resolutionMatches && fpsMatches
			if result.IsValid {
				result.Message = "Stream configuration matches expected values"
			} else {
				result.Message = "Stream configuration doesn't match expected values"
				if !resolutionMatches {
					result.Error = fmt.Sprintf("Resolution mismatch: got %dx%d, expected %dx%d",
						result.ActualWidth, result.ActualHeight, result.ExpectedWidth, result.ExpectedHeight)
				} else if !fpsMatches {
					result.Error = fmt.Sprintf("Frame rate mismatch: got %.1f fps, expected %.1f fps",
						result.ActualFPS, result.ExpectedFPS)
				}
			}

			fmt.Printf("Validation result from direct ffprobe: %+v\n", result)
			return result, nil
		}

		// Fall back to the original validation method if direct extraction failed
		valid, err := lib.ValidateStreamConfig(
			streamURL,
			expectedConfig.Width,
			expectedConfig.Height,
			float64(expectedConfig.FrameRate),
		)

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

			// Parse resolution from error message
			resMatcher := regexp.MustCompile(`got (\d+)x(\d+), expected`)
			if matches := resMatcher.FindStringSubmatch(errStr); len(matches) == 3 {
				if width, err := strconv.Atoi(matches[1]); err == nil {
					result.ActualWidth = width
				}
				if height, err := strconv.Atoi(matches[2]); err == nil {
					result.ActualHeight = height
				}
				fmt.Printf("Parsed actual resolution: %dx%d\n", result.ActualWidth, result.ActualHeight)
			} else {
				fmt.Printf("Failed to parse resolution from: %s\n", errStr)
			}

			// Parse frame rate from error message
			fpsMatcher := regexp.MustCompile(`@ (\d+\.?\d*) fps`)
			if matches := fpsMatcher.FindStringSubmatch(errStr); len(matches) == 2 {
				if fps, err := strconv.ParseFloat(matches[1], 64); err == nil {
					result.ActualFPS = fps
				}
				fmt.Printf("Parsed actual frame rate: %.2f fps\n", result.ActualFPS)
			} else {
				fmt.Printf("Failed to parse frame rate from: %s\n", errStr)
			}

			// Debug the complete validation result
			fmt.Printf("Validation result after parsing: %+v\n", result)

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
