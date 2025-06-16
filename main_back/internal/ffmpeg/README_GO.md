# Go RTSP Stream Analyzer

This package provides a Go wrapper around the C RTSP analyzer using CGO. It allows you to analyze RTSP streams and extract codec, resolution, and frame rate information.

## Files

- `rtsp_analyzer.go` - Go package with CGO wrapper functions
- `test_analyzer.go` - Command-line test program
- `example_usage.go` - Example showing integration into applications
- `build_go.bat` - Windows build script
- `rtsp_analyzer.c` - Original C implementation

## Prerequisites

1. **FFmpeg Development Libraries**
   - Download FFmpeg development libraries for Windows
   - Extract to `C:\ffmpeg` (or update paths in build script)
   - Ensure you have both include files and lib files

2. **GCC Compiler**
   - Install MSYS2/MinGW64 or another GCC distribution
   - Add gcc to your system PATH

3. **Go with CGO Support**
   - Go 1.11+ with CGO enabled
   - Set `CGO_ENABLED=1` environment variable

## Building

1. **Using the build script:**
   ```cmd
   cd d:\VNG\test\main_onvif\main_back\internal\ffmpeg
   build_go.bat
   ```

2. **Manual build:**
   ```cmd
   set CGO_ENABLED=1
   set CGO_CFLAGS=-IC:\ffmpeg\include
   set CGO_LDFLAGS=-LC:\ffmpeg\lib -lavformat -lavcodec -lavutil
   go build -o test_analyzer.exe test_analyzer.go
   ```

## Usage

### Command Line Tool
```cmd
test_analyzer.exe "rtsp://admin:admin123@192.168.1.100:554/live"
```

### In Go Code
```go
package main

import (
    "fmt"
    "main_back/internal/ffmpeg"
)

func main() {
    info, err := ffmpeg.AnalyzeRTSPStream("rtsp://your_camera_url")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    fmt.Printf("Codec: %s\n", info.Codec)
    fmt.Printf("Resolution: %s\n", info.GetStreamResolution())
    fmt.Printf("Frame Rate: %s\n", info.GetFrameRate())
}
```

## API Reference

### Types

```go
type StreamInfo struct {
    Codec      string  `json:"codec"`
    Width      int     `json:"width"`
    Height     int     `json:"height"`
    FPS        float64 `json:"fps"`
    Success    bool    `json:"success"`
    ErrorMsg   string  `json:"error_msg,omitempty"`
}
```

### Functions

- `AnalyzeRTSPStream(rtspURL string) (*StreamInfo, error)` - Main analysis function
- `GetStreamResolution() string` - Returns formatted resolution (e.g., "1920x1080")
- `GetFrameRate() string` - Returns formatted frame rate (e.g., "30.00 fps")
- `String() string` - Returns complete formatted string
- `IsHighDefinition() bool` - Returns true if 720p or higher
- `IsFullHD() bool` - Returns true if 1080p or higher
- `Is4K() bool` - Returns true if 2160p or higher

## Integration with Camera Management

You can integrate this into your camera management system:

```go
// In your camera package
import "main_back/internal/ffmpeg"

func (c *Camera) AnalyzeStream() error {
    info, err := ffmpeg.AnalyzeRTSPStream(c.RTSPUrl)
    if err != nil {
        return err
    }
    
    c.Codec = info.Codec
    c.Width = info.Width
    c.Height = info.Height
    c.FPS = info.FPS
    
    return nil
}
```

## Error Handling

The package provides detailed error messages for common issues:
- Network connection failures
- Authentication errors
- Stream format issues
- FFmpeg library errors

## Performance Notes

- The analysis typically takes 1-5 seconds depending on network latency
- Uses TCP transport for better reliability
- Configured with 5-second connection timeout
- Optimized for low latency analysis

## Troubleshooting

1. **CGO Compilation Errors**
   - Ensure FFmpeg libraries are properly installed
   - Check that GCC is in your PATH
   - Verify CGO_CFLAGS and CGO_LDFLAGS paths

2. **Runtime Errors**
   - Check RTSP URL format
   - Verify camera credentials
   - Ensure camera is accessible from your network

3. **Performance Issues**
   - Check network connectivity to camera
   - Verify camera supports the requested transport method
   - Consider adjusting timeout values in the C code
