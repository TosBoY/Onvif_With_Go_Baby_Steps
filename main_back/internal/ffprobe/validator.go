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

	"main_back/pkg/models"
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

// ValidateStream runs ffprobe using an OS-specific binary from a relative path and checks if the actual stream config matches the expected one.
// It uses a path relative to the executable and incorporates parsing logic from the previous implementation.
func ValidateStream(rtspURL string, expected models.EncoderConfig) (bool, error) {
	// Determine the executable directory
	execPath, err := os.Executable()
	if err != nil {
		return false, fmt.Errorf("failed to get executable path: %w", err)
	}
	execDir := filepath.Dir(execPath)

	// Base path to the ffprobe directory relative to the backend executable (assuming executable is in main_back)
	// This navigates up one directory from main_back to the project root and then into the ffprobe folder.
	baseFFprobeDir := filepath.Join(execDir, "..", "ffprobe")

	// Determine OS-specific ffprobe binary name
	osName := runtime.GOOS
	ffprobeBinaryName := ""

	switch osName {
	case "windows":
		ffprobeBinaryName = "ffprobe.exe"
	case "linux":
		ffprobeBinaryName = "ffprobe_linux"
	case "darwin":
		ffprobeBinaryName = "ffprobe_darwin"
	default:
		return false, fmt.Errorf("unsupported operating system: %s", osName)
	}

	// Construct the full relative path to the ffprobe binary
	ffprobePath := filepath.Join(baseFFprobeDir, ffprobeBinaryName)

	// Check if the ffprobe binary exists at the constructed path
	if _, err := os.Stat(ffprobePath); os.IsNotExist(err) {
		return false, fmt.Errorf("ffprobe binary not found at path: %s", ffprobePath)
	}

	// Give the camera some time to apply settings before validation
	time.Sleep(5 * time.Second)

	cmd := exec.Command(ffprobePath,
		"-print_format", "json",
		"-show_streams",
		"-select_streams", "v:0", // Only select the video stream (index 0)
		rtspURL)

	output, err := cmd.Output()
	if err != nil {
		// Include command output and stderr in the error for easier debugging, similar to old implementation
		exitError, ok := err.(*exec.ExitError)
		if ok {
			return false, fmt.Errorf("ffprobe command failed: %w\nStdout: %s\nStderr: %s", err, output, exitError.Stderr)
		}
		return false, fmt.Errorf("ffprobe command failed: %w\nStdout: %s", err, output)
	}

	// Check for empty output
	if len(output) == 0 {
		return false, fmt.Errorf("ffprobe returned no data")
	}

	// Parse the JSON output
	var result FFProbeResult
	if err := json.Unmarshal(output, &result); err != nil {
		return false, fmt.Errorf("failed to parse ffprobe output: %w", err)
	}

	// Find the video stream (should be the first if -select_streams v:0 worked)
	if len(result.Streams) == 0 {
		return false, fmt.Errorf("no video stream found in ffprobe output")
	}

	stream := result.Streams[0] // Assuming the first stream is the video stream due to -select_streams v:0

	// Parse frame rate using the helper function from the old implementation
	actualFPS, err := parseFrameRate(&stream)
	if err != nil {
		return false, fmt.Errorf("failed to parse frame rate: %w", err)
	}

	// Perform validation against expected configuration
	match := stream.Width == expected.Resolution.Width &&
		stream.Height == expected.Resolution.Height &&
		(int(actualFPS+0.5) == expected.FPS) // allow rounding

	if !match {
		var discrepancies []string
		if stream.Width != expected.Resolution.Width || stream.Height != expected.Resolution.Height {
			discrepancies = append(discrepancies, fmt.Sprintf("resolution mismatch: got %dx%d, expected %dx%d", stream.Width, stream.Height, expected.Resolution.Width, expected.Resolution.Height))
		}
		if int(actualFPS+0.5) != expected.FPS {
			discrepancies = append(discrepancies, fmt.Sprintf("FPS mismatch: got %.2f, expected %d", actualFPS, expected.FPS))
		}

		return false, fmt.Errorf("stream mismatch: %s", strings.Join(discrepancies, ", "))
	}

	// If validation passes
	return true, nil
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
