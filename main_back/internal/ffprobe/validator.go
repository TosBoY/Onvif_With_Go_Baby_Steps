package ffprobe

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type FFProbeStream struct {
	CodecType    string `json:"codec_type"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	RFrameRate   string `json:"r_frame_rate"`   // Real frame rate (as backup)
	AvgFrameRate string `json:"avg_frame_rate"` // Average frame rate (preferred)
	BitRate      string `json:"bit_rate"`       // Bit rate in bits per second
	CodecName    string `json:"codec_name"`
}

type FFProbeResult struct {
	Streams []FFProbeStream `json:"streams"`
}

type ValidationResult struct {
	IsValid        bool    `json:"isValid"`
	ExpectedWidth  int     `json:"expectedWidth"`
	ExpectedHeight int     `json:"expectedHeight"`
	ExpectedFPS    int     `json:"expectedFPS"`
	ActualWidth    int     `json:"actualWidth"`
	ActualHeight   int     `json:"actualHeight"`
	ActualFPS      float64 `json:"actualFPS"`
	Error          string  `json:"error,omitempty"`
}

// ValidateStreamDetailed returns detailed validation results
func ValidateStream(rtspURL string, expectedWidth, expectedHeight, expectedFPS int) (*ValidationResult, error) {
	// Initialize the result with expected values
	result := &ValidationResult{
		IsValid:        false,
		ExpectedWidth:  expectedWidth,
		ExpectedHeight: expectedHeight,
		ExpectedFPS:    expectedFPS,
	}

	workingDir, err := os.Getwd()
	if err != nil {
		result.Error = fmt.Sprintf("failed to get working directory: %v", err)
		return result, nil
	}

	// Navigate to the ffprobe directory from working directory
	ffprobeDir := filepath.Join(workingDir, "..", "..", "ffprobe")

	// Determine OS-specific ffprobe binary name
	ffprobeBinaryName := ""
	switch runtime.GOOS {
	case "windows":
		ffprobeBinaryName = "ffprobe.exe"
	case "linux":
		ffprobeBinaryName = "ffprobe_linux"
	case "darwin":
		ffprobeBinaryName = "ffprobe_darwin"
	default:
		result.Error = fmt.Sprintf("unsupported OS: %s", runtime.GOOS)
		return result, nil
	}

	ffprobePath := filepath.Join(ffprobeDir, ffprobeBinaryName)

	// Validate binary exists
	if _, err := os.Stat(ffprobePath); os.IsNotExist(err) {
		result.Error = fmt.Sprintf("ffprobe binary not found at path: %s", ffprobePath)
		return result, nil
	}

	// Give the camera some time to apply settings before validation
	time.Sleep(1 * time.Second)

	cmd := exec.Command(ffprobePath,
		"-v", "quiet",
		"-print_format", "json",
		"-show_streams",
		"-select_streams", "v:0", // Only select the video stream (index 0)
		rtspURL)

	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			fmt.Printf("FFprobe stderr: %s\n", string(exitErr.Stderr))
		}
		fmt.Printf("FFprobe execution failed: %v\n", err)
		return nil, fmt.Errorf("stream info failed: %v", err)
	}

	if len(output) == 0 {
		fmt.Println("FFprobe returned empty output")
		return nil, fmt.Errorf("ffprobe returned no data")
	}

	// Parse the JSON output
	var ffprobeResult FFProbeResult
	if err := json.Unmarshal(output, &ffprobeResult); err != nil {
		result.Error = fmt.Sprintf("failed to parse ffprobe output: %v", err)
		return result, nil
	}

	// Find the video stream
	if len(ffprobeResult.Streams) == 0 {
		result.Error = "no video stream found in ffprobe output"
		return nil, fmt.Errorf("no streams found in video")
	}

	for _, stream := range ffprobeResult.Streams {
		if stream.CodecType == "video" {
			// Parse frame rate
			actualFrameRate, err := parseFrameRate(&stream)
			if err != nil {
				fmt.Printf("Failed to parse frame rate: %v\n", err)
				actualFrameRate = 0
			}

			result.ActualWidth = stream.Width
			result.ActualHeight = stream.Height
			result.ActualFPS = actualFrameRate

			fmt.Printf("Found video stream: %dx%d @ %.2f fps, codec: %s\n",
				result.ActualWidth, result.ActualHeight, result.ActualFPS, stream.CodecName)
		}
	}

	// Check if validation passes
	resolutionMatch := result.ActualWidth > 0 && result.ActualHeight > 0 &&
		result.ActualWidth == result.ExpectedWidth && result.ActualHeight == result.ExpectedHeight

	// Only consider FPS match if we have a valid FPS value
	fpsMatch := result.ActualFPS > 0 && int(result.ActualFPS+0.5) == result.ExpectedFPS

	// Only consider valid if we have all the necessary information
	result.IsValid = resolutionMatch && fpsMatch

	if !result.IsValid {
		var errors []string

		// Only report resolution mismatch if we have actual values
		if result.ActualWidth > 0 && result.ActualHeight > 0 {
			if !resolutionMatch {
				errors = append(errors, fmt.Sprintf("resolution mismatch: got %dx%d, expected %dx%d",
					result.ActualWidth, result.ActualHeight, result.ExpectedWidth, result.ExpectedHeight))
			}
		} else {
			errors = append(errors, "failed to detect actual resolution")
		}

		// Only report FPS mismatch if we have an actual FPS value
		if result.ActualFPS > 0 {
			if !fpsMatch {
				errors = append(errors, fmt.Sprintf("FPS mismatch: got %.2f, expected %d", result.ActualFPS, result.ExpectedFPS))
			}
		} else {
			errors = append(errors, "failed to detect actual FPS")
		}

		// Set the error message
		if len(errors) > 0 {
			result.Error = strings.Join(errors, "; ")
		}
	}

	return result, nil
}

// parseFrameRate converts a frame rate string (like "30/1") to float64.
// Copied from the previous implementation (onvif_back/lib/ffprobe.go).
func parseFrameRate(stream *FFProbeStream) (float64, error) {
	// Try avg_frame_rate first as it's more accurate for variable frame rate streams
	if stream.AvgFrameRate != "" && stream.AvgFrameRate != "0/0" {
		return evaluateFrameRateExpression(stream.AvgFrameRate)
	}

	// Fall back to r_frame_rate if avg_frame_rate is not available
	if stream.RFrameRate != "" && stream.RFrameRate != "0/0" {
		return evaluateFrameRateExpression(stream.RFrameRate)
	}

	return 0, fmt.Errorf("no valid frame rate found in stream")
}

// evaluateFrameRateExpression evaluates a frame rate expression string (like "30/1").
// Copied from the previous implementation (onvif_back/lib/ffprobe.go).
func evaluateFrameRateExpression(expr string) (float64, error) {
	parts := strings.Split(expr, "/")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid frame rate format: %s", expr)
	}

	num, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return 0, fmt.Errorf("invalid numerator in frame rate: %s", err)
	}

	den, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return 0, fmt.Errorf("invalid denominator in frame rate: %s", err)
	}

	if den == 0 {
		return 0, fmt.Errorf("zero denominator in frame rate: %s", expr)
	}

	return num / den, nil
}
