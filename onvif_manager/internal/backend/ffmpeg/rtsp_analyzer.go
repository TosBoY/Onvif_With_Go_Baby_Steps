package ffmpeg

/*
#cgo windows CFLAGS: -IC:/ffmpeg/include
#cgo windows LDFLAGS: -LC:/ffmpeg/lib -lavformat -lavcodec -lavutil
#cgo linux CFLAGS: -I/usr/include -I/usr/local/include
#cgo linux LDFLAGS: -lavformat -lavcodec -lavutil
#cgo linux pkg-config: libavformat libavcodec libavutil
#cgo darwin CFLAGS: -I/usr/local/include -I/opt/homebrew/include
#cgo darwin LDFLAGS: -L/usr/local/lib -L/opt/homebrew/lib -lavformat -lavcodec -lavutil
#cgo darwin pkg-config: libavformat libavcodec libavutil

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <libavformat/avformat.h>
#include <libavcodec/avcodec.h>
#include <libavutil/avutil.h>

typedef struct {
    char codec[64];
    int width;
    int height;
    double fps;
    int bitrate;
    int success;
    char error_msg[256];
} StreamInfo;

StreamInfo analyze_rtsp_stream(const char* rtsp_url) {
    StreamInfo info = {0};
    AVFormatContext *format_ctx = NULL;
    int ret;

    // Initialize FFmpeg
    #if LIBAVCODEC_VERSION_INT < AV_VERSION_INT(58, 9, 100)
        av_register_all();
    #endif

    avformat_network_init();

    // RTSP options for low latency
    AVDictionary *options = NULL;
    av_dict_set(&options, "rtsp_transport", "tcp", 0);
    av_dict_set(&options, "max_delay", "500000", 0);
    av_dict_set(&options, "stimeout", "5000000", 0);

    // Open input
    ret = avformat_open_input(&format_ctx, rtsp_url, NULL, &options);
    av_dict_free(&options);

    if (ret < 0) {
        char err_buf[AV_ERROR_MAX_STRING_SIZE];
        av_strerror(ret, err_buf, sizeof(err_buf));
        snprintf(info.error_msg, sizeof(info.error_msg), "Could not open input: %s", err_buf);
        info.success = 0;
        return info;
    }

    // Retrieve stream information
    ret = avformat_find_stream_info(format_ctx, NULL);
    if (ret < 0) {
        char err_buf[AV_ERROR_MAX_STRING_SIZE];
        av_strerror(ret, err_buf, sizeof(err_buf));
        snprintf(info.error_msg, sizeof(info.error_msg), "Could not find stream info: %s", err_buf);
        info.success = 0;
        avformat_close_input(&format_ctx);
        return info;
    }

    // Find the first video stream
    for (unsigned int i = 0; i < format_ctx->nb_streams; i++) {
        AVStream *stream = format_ctx->streams[i];
        AVCodecParameters *codec_params = stream->codecpar;

        if (codec_params->codec_type == AVMEDIA_TYPE_VIDEO) {
            // Get codec name
            const char *codec_name = avcodec_get_name(codec_params->codec_id);
            strncpy(info.codec, codec_name, sizeof(info.codec) - 1);
            info.codec[sizeof(info.codec) - 1] = '\0';            // Get resolution
            info.width = codec_params->width;
            info.height = codec_params->height;

            // Get bitrate (convert from bits/sec to kbps)
            if (codec_params->bit_rate > 0) {
                info.bitrate = codec_params->bit_rate / 1000;
            } else {
                info.bitrate = 0;
            }

            // Get frame rate
            if (stream->avg_frame_rate.den && stream->avg_frame_rate.num) {
                info.fps = av_q2d(stream->avg_frame_rate);
            } else {
                info.fps = 0.0;
            }

            info.success = 1;
            break;
        }
    }

    if (!info.success) {
        snprintf(info.error_msg, sizeof(info.error_msg), "No video stream found in RTSP stream");
    }

    // Clean up
    avformat_close_input(&format_ctx);
    avformat_network_deinit();

    return info;
}
*/
import "C"

import (
	"fmt"
	"strings"
	"unsafe"
)

// StreamInfo represents the information extracted from an RTSP stream
type StreamInfo struct {
	Codec    string  `json:"codec"`
	Width    int     `json:"width"`
	Height   int     `json:"height"`
	FPS      float64 `json:"fps"`
	Bitrate  int     `json:"bitrate"` // in kbps
	Success  bool    `json:"success"`
	ErrorMsg string  `json:"error_msg,omitempty"`
}

// AnalyzeRTSPStream analyzes an RTSP stream and returns codec, resolution, and FPS information
func AnalyzeRTSPStream(rtspURL string) (*StreamInfo, error) {
	if rtspURL == "" {
		return nil, fmt.Errorf("RTSP URL cannot be empty")
	}

	// Convert Go string to C string
	cURL := C.CString(rtspURL)
	defer C.free(unsafe.Pointer(cURL))

	// Call the C function
	cInfo := C.analyze_rtsp_stream(cURL) // Convert C struct to Go struct
	info := &StreamInfo{
		Codec:    C.GoString(&cInfo.codec[0]),
		Width:    int(cInfo.width),
		Height:   int(cInfo.height),
		FPS:      float64(cInfo.fps),
		Bitrate:  int(cInfo.bitrate),
		Success:  int(cInfo.success) == 1,
		ErrorMsg: C.GoString(&cInfo.error_msg[0]),
	}

	if !info.Success {
		return info, fmt.Errorf("failed to analyze RTSP stream: %s", info.ErrorMsg)
	}

	return info, nil
}

// GetStreamResolution returns the resolution as a formatted string (e.g., "1920x1080")
func (s *StreamInfo) GetStreamResolution() string {
	if s.Width == 0 || s.Height == 0 {
		return "Unknown"
	}
	return fmt.Sprintf("%dx%d", s.Width, s.Height)
}

// GetFrameRate returns the frame rate as a formatted string
func (s *StreamInfo) GetFrameRate() string {
	if s.FPS == 0 {
		return "Unknown"
	}
	return fmt.Sprintf("%.2f fps", s.FPS)
}

// GetBitrate returns the bitrate as a formatted string
func (s *StreamInfo) GetBitrate() string {
	if s.Bitrate == 0 {
		return "Unknown"
	}
	return fmt.Sprintf("%d kbps", s.Bitrate)
}

// String returns a formatted string representation of the stream info
func (s *StreamInfo) String() string {
	if !s.Success {
		return fmt.Sprintf("Error: %s", s.ErrorMsg)
	}

	return fmt.Sprintf("Codec: %s, Resolution: %s, Frame Rate: %s, Bitrate: %s",
		s.Codec, s.GetStreamResolution(), s.GetFrameRate(), s.GetBitrate())
}

// IsHighDefinition returns true if the stream is HD (720p) or higher
func (s *StreamInfo) IsHighDefinition() bool {
	return s.Height >= 720
}

// IsFullHD returns true if the stream is Full HD (1080p) or higher
func (s *StreamInfo) IsFullHD() bool {
	return s.Height >= 1080
}

// Is4K returns true if the stream is 4K (2160p) or higher
func (s *StreamInfo) Is4K() bool {
	return s.Height >= 2160
}

type ValidationResult struct {
	IsValid         bool    `json:"isValid"`
	ExpectedWidth   int     `json:"expectedWidth"`
	ExpectedHeight  int     `json:"expectedHeight"`
	ExpectedFPS     int     `json:"expectedFPS"`
	ExpectedBitrate int     `json:"expectedBitrate"`
	ActualWidth     int     `json:"actualWidth"`
	ActualHeight    int     `json:"actualHeight"`
	ActualFPS       float64 `json:"actualFPS"`
	ActualBitrate   int     `json:"actualBitrate"`
	Error           string  `json:"error,omitempty"`
}

func ValidateStream(rtspURL string, expectedWidth, expectedHeight, expectedFPS, expectedBitrate int) (*ValidationResult, error) {
	result := &ValidationResult{
		ExpectedWidth:   expectedWidth,
		ExpectedHeight:  expectedHeight,
		ExpectedFPS:     expectedFPS,
		ExpectedBitrate: expectedBitrate,
	}

	streamInfo, err := AnalyzeRTSPStream(rtspURL)
	if err != nil {
		result.Error = fmt.Sprintf("Failed to analyze RTSP stream: %v", err)
		return result, nil
	}
	result.ActualWidth = streamInfo.Width
	result.ActualHeight = streamInfo.Height
	result.ActualFPS = streamInfo.FPS
	result.ActualBitrate = streamInfo.Bitrate

	if !streamInfo.Success {
		result.Error = streamInfo.ErrorMsg
		return result, nil
	}
	// Perform validation logic with new business rules:
	// Resolution mismatch = failure, FPS/bitrate mismatch = warning only
	resolutionMatch := result.ActualWidth > 0 && result.ActualHeight > 0 &&
		result.ActualWidth == result.ExpectedWidth && result.ActualHeight == result.ExpectedHeight
	// Only consider FPS match if we have a valid FPS value
	fpsMatch := result.ActualFPS > 0 && int(result.ActualFPS+0.5) == result.ExpectedFPS

	// Only consider bitrate match if expected bitrate was provided and we have actual bitrate
	bitrateMatch := true // Default to true if no expected bitrate
	if result.ExpectedBitrate > 0 {
		if result.ActualBitrate > 0 {
			// Allow 10% tolerance for bitrate comparison
			tolerance := float64(result.ExpectedBitrate) * 0.1
			diff := float64(result.ActualBitrate - result.ExpectedBitrate)
			if diff < 0 {
				diff = -diff
			}
			bitrateMatch = diff <= tolerance
		} else {
			bitrateMatch = false // Expected bitrate but couldn't detect actual
		}
	}

	// NEW BUSINESS LOGIC: Only resolution mismatch causes failure
	// FPS and bitrate mismatches are warnings only
	result.IsValid = resolutionMatch // Only require resolution to match for success
	// Generate error/warning messages with clear distinction
	if !result.IsValid || !fpsMatch || !bitrateMatch {
		var errors []string

		// Resolution mismatch = ERROR (causes failure)
		if !resolutionMatch {
			if result.ActualWidth > 0 && result.ActualHeight > 0 {
				errors = append(errors, fmt.Sprintf("RESOLUTION MISMATCH (ERROR): got %dx%d, expected %dx%d",
					result.ActualWidth, result.ActualHeight, result.ExpectedWidth, result.ExpectedHeight))
			} else {
				errors = append(errors, "RESOLUTION VALIDATION FAILED (ERROR): unable to detect actual resolution")
			}
		}

		// FPS mismatch = WARNING (does not cause failure)
		if !fpsMatch {
			if result.ActualFPS > 0 {
				errors = append(errors, fmt.Sprintf("FPS DIFFERENCE (WARNING): got %.2f fps, expected %d fps",
					result.ActualFPS, result.ExpectedFPS))
			} else {
				errors = append(errors, "FPS DETECTION FAILED (WARNING): unable to detect actual FPS")
			}
		}

		// Bitrate mismatch = WARNING (does not cause failure)
		if !bitrateMatch && result.ExpectedBitrate > 0 {
			if result.ActualBitrate > 0 {
				errors = append(errors, fmt.Sprintf("BITRATE DIFFERENCE (WARNING): got %d kbps, expected %d kbps",
					result.ActualBitrate, result.ExpectedBitrate))
			} else {
				errors = append(errors, "BITRATE DETECTION FAILED (WARNING): unable to detect actual bitrate")
			}
		}

		// Set the error message
		if len(errors) > 0 {
			result.Error = strings.Join(errors, "; ")
		}
	}

	return result, nil
}
