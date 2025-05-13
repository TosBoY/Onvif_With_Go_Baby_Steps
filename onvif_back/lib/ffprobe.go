package onvif_test

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

// FFProbeStream represents the stream information from ffprobe output
type FFProbeStream struct {
	CodecType    string `json:"codec_type"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	RFrameRate   string `json:"r_frame_rate"`   // Real frame rate (as backup)
	AvgFrameRate string `json:"avg_frame_rate"` // Average frame rate (preferred)
	BitRate      string `json:"bit_rate"`       // Bit rate in bits per second
}

// FFProbeResult represents the ffprobe command output structure
type FFProbeResult struct {
	Streams []FFProbeStream `json:"streams"`
}

// GetFFProbePath returns the path to the appropriate ffprobe binary for the current OS
func GetFFProbePath() (string, error) {
	// Try to get the directory containing this source file
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("failed to get current file path")
	}

	// Go up from lib/ to the project root and then into ffprobe/
	projectRoot := filepath.Dir(filepath.Dir(thisFile)) // up from lib to project root
	ffprobeDir := filepath.Join(projectRoot, "ffprobe")

	var ffprobeName string
	switch runtime.GOOS {
	case "windows":
		ffprobeName = "ffprobe.exe"
	case "linux":
		ffprobeName = "ffprobe_linux"
	case "darwin":
		ffprobeName = "ffprobe_darwin"
	default:
		return "", fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	ffprobePath := filepath.Join(ffprobeDir, ffprobeName)
	fmt.Printf("Looking for ffprobe at: %s\n", ffprobePath)

	if _, err := os.Stat(ffprobePath); err != nil {
		if os.IsNotExist(err) {
			// Try current working directory as fallback
			cwd, _ := os.Getwd()
			ffprobePath = filepath.Join(cwd, "ffprobe", ffprobeName)
			if _, err := os.Stat(ffprobePath); err != nil {
				return "", fmt.Errorf("ffprobe binary (%s) not found in %s or %s",
					ffprobeName,
					filepath.Join(ffprobeDir, ffprobeName),
					ffprobePath)
			}
		} else {
			return "", fmt.Errorf("error accessing ffprobe binary: %v", err)
		}
	}

	return ffprobePath, nil
}

// EnsureFFProbeExecutable ensures the ffprobe binary has executable permissions
func EnsureFFProbeExecutable() error {
	fmt.Println("Ensuring ffprobe binary is executable...")
	ffprobePath, err := GetFFProbePath()
	if err != nil {
		return fmt.Errorf("failed to get ffprobe path: %v", err)
	}

	fmt.Printf("Found ffprobe at: %s\n", ffprobePath)

	// Check if file exists and is accessible
	fileInfo, err := os.Stat(ffprobePath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("ffprobe binary does not exist at path: %s", ffprobePath)
		}
		return fmt.Errorf("error accessing ffprobe binary at %s: %v", ffprobePath, err)
	}

	// Check if it's a regular file
	if !fileInfo.Mode().IsRegular() {
		return fmt.Errorf("ffprobe path %s is not a regular file", ffprobePath)
	}

	// Only change permissions on Unix-like systems
	if runtime.GOOS != "windows" {
		currentMode := fileInfo.Mode()
		if currentMode&0111 == 0 { // Check if file is executable
			fmt.Printf("Setting executable permissions on %s\n", ffprobePath)
			err := os.Chmod(ffprobePath, 0755)
			if err != nil {
				return fmt.Errorf("failed to set executable permissions on %s: %v", ffprobePath, err)
			}
		}
	}

	// Try to execute ffprobe with -version to verify it works
	fmt.Println("Verifying ffprobe execution...")
	cmd := exec.Command(ffprobePath, "-version")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("ffprobe binary at %s failed to execute: %v\nOutput: %s", ffprobePath, err, string(output))
	} else {
		fmt.Printf("ffprobe version check succeeded: %s\n", string(output))
	}

	return nil
}

// ValidateStreamConfig checks if the stream configuration matches the expected values
func ValidateStreamConfig(streamURL string, expectedWidth, expectedHeight int, expectedFrameRate float64) (bool, error) {
	ffprobePath, err := GetFFProbePath()
	if err != nil {
		fmt.Printf("FFprobe path error: %v\n", err)
		return false, fmt.Errorf("failed to get ffprobe: %v", err)
	}

	// Ensure ffprobe is executable
	if err := EnsureFFProbeExecutable(); err != nil {
		fmt.Printf("FFprobe executable check failed: %v\n", err)
		return false, fmt.Errorf("ffprobe executable check failed: %v", err)
	}

	fmt.Printf("Validating stream %s with ffprobe at %s\n", streamURL, ffprobePath)
	fmt.Printf("Expected resolution: %dx%d @ %f fps\n", expectedWidth, expectedHeight, expectedFrameRate)

	// Run ffprobe command
	cmd := exec.Command(ffprobePath,
		"-v", "quiet",
		"-print_format", "json",
		"-show_streams",
		streamURL)

	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			fmt.Printf("FFprobe stderr: %s\n", string(exitErr.Stderr))
		}
		fmt.Printf("FFprobe execution failed: %v\n", err)
		return false, fmt.Errorf("stream validation failed: %v", err)
	}

	if len(output) == 0 {
		fmt.Println("FFprobe returned empty output")
		return false, fmt.Errorf("ffprobe returned no data")
	}

	fmt.Printf("FFprobe raw output: %s\n", string(output))

	// Parse the JSON output
	var result FFProbeResult
	if err := json.Unmarshal(output, &result); err != nil {
		fmt.Printf("Failed to parse FFprobe JSON output: %v\nOutput was: %s\n", err, string(output))
		return false, fmt.Errorf("failed to parse ffprobe output: %v", err)
	}

	if len(result.Streams) == 0 {
		fmt.Println("No streams found in FFprobe output")
		return false, fmt.Errorf("no streams found in video")
	}

	// Find the video stream
	for _, stream := range result.Streams {
		if stream.CodecType == "video" {
			fmt.Printf("Found video stream: %dx%d", stream.Width, stream.Height)

			// Parse frame rate
			actualFrameRate, err := parseFrameRate(&stream)
			if err != nil {
				fmt.Printf("Failed to parse frame rate: %v\n", err)
				return false, nil
			}
			fmt.Printf(" @ %f fps\n", actualFrameRate)

			// Check dimensions
			if stream.Width != expectedWidth || stream.Height != expectedHeight {
				fmt.Printf("Resolution mismatch: got %dx%d, expected %dx%d\n",
					stream.Width, stream.Height, expectedWidth, expectedHeight)
				return false, nil
			}

			// Check frame rate with some tolerance (0.1)
			if math.Abs(actualFrameRate-expectedFrameRate) > 0.1 {
				fmt.Printf("Frame rate mismatch: got %f, expected %f\n",
					actualFrameRate, expectedFrameRate)
				return false, nil
			}

			fmt.Println("Stream validation successful")
			return true, nil
		}
	}

	fmt.Println("No video stream found")
	return false, fmt.Errorf("no video stream found")
}

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
		return 0, fmt.Errorf("division by zero in frame rate")
	}

	return num / den, nil
}
