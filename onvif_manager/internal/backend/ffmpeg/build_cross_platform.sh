#!/bin/bash
# Cross-platform build script for rtsp_analyzer

echo "Building RTSP Analyzer for multiple platforms..."

# Set CGO environment
export CGO_ENABLED=1

# Function to check if FFmpeg is installed
check_ffmpeg() {
    if command -v pkg-config >/dev/null 2>&1; then
        if pkg-config --exists libavformat libavcodec libavutil; then
            echo "✓ FFmpeg found via pkg-config"
            return 0
        fi
    fi
    
    # Check common installation paths
    if [ -d "/usr/include/libavformat" ] || [ -d "/usr/local/include/libavformat" ]; then
        echo "✓ FFmpeg headers found in system paths"
        return 0
    fi
    
    echo "✗ FFmpeg not found. Please install FFmpeg development libraries."
    echo "  Ubuntu/Debian: sudo apt-get install libavformat-dev libavcodec-dev libavutil-dev"
    echo "  CentOS/RHEL: sudo yum install ffmpeg-devel"
    echo "  macOS: brew install ffmpeg"
    return 1
}

# Detect OS and build accordingly
case "$(uname -s)" in
    Linux*)
        echo "Building for Linux..."
        if ! check_ffmpeg; then
            exit 1
        fi
        go build -o rtsp_analyzer_linux rtsp_analyzer.go
        ;;
    Darwin*)
        echo "Building for macOS..."
        if ! check_ffmpeg; then
            exit 1
        fi
        go build -o rtsp_analyzer_macos rtsp_analyzer.go
        ;;
    CYGWIN*|MINGW*|MSYS*)
        echo "Building for Windows..."
        # Assume FFmpeg is in C:\ffmpeg (Windows build handled separately)
        go build -o rtsp_analyzer_windows.exe rtsp_analyzer.go
        ;;
    *)
        echo "Unsupported OS: $(uname -s)"
        exit 1
        ;;
esac

echo "Build completed successfully!"
