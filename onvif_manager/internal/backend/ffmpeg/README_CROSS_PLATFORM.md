# RTSP Analyzer - Cross-Platform Setup

The RTSP Analyzer uses FFmpeg libraries through CGO and now supports Windows, Linux, and macOS.

## Platform-Specific Setup

### Windows
1. **Install FFmpeg**:
   - Download FFmpeg dev libraries from https://ffmpeg.org/download.html
   - Extract to `C:\ffmpeg\` (or adjust paths in build script)
   - Ensure you have the dev/shared libraries, not just the executables

2. **Install Build Tools**:
   - Install MSYS2/MinGW64 or TDM-GCC
   - Ensure `gcc` is in your PATH

3. **Build**:
   ```cmd
   build_go.bat
   ```

### Linux (Ubuntu/Debian)
1. **Install FFmpeg Development Libraries**:
   ```bash
   sudo apt-get update
   sudo apt-get install libavformat-dev libavcodec-dev libavutil-dev
   sudo apt-get install build-essential  # for gcc
   ```

2. **Build**:
   ```bash
   chmod +x build_cross_platform.sh
   ./build_cross_platform.sh
   ```

### Linux (CentOS/RHEL/Fedora)
1. **Install FFmpeg Development Libraries**:
   ```bash
   # CentOS/RHEL (with EPEL)
   sudo yum install epel-release
   sudo yum install ffmpeg-devel gcc

   # Fedora
   sudo dnf install ffmpeg-devel gcc
   ```

2. **Build**:
   ```bash
   chmod +x build_cross_platform.sh
   ./build_cross_platform.sh
   ```

### macOS
1. **Install FFmpeg**:
   ```bash
   # Using Homebrew
   brew install ffmpeg

   # Using MacPorts
   sudo port install ffmpeg
   ```

2. **Build**:
   ```bash
   chmod +x build_cross_platform.sh
   ./build_cross_platform.sh
   ```

## CGO Directives

The source file now includes platform-specific CGO directives:

```go
/*
#cgo windows CFLAGS: -IC:/ffmpeg/include
#cgo windows LDFLAGS: -LC:/ffmpeg/lib -lavformat -lavcodec -lavutil
#cgo linux CFLAGS: -I/usr/include -I/usr/local/include
#cgo linux LDFLAGS: -lavformat -lavcodec -lavutil
#cgo darwin CFLAGS: -I/usr/local/include -I/opt/homebrew/include
#cgo darwin LDFLAGS: -L/usr/local/lib -L/opt/homebrew/lib -lavformat -lavcodec -lavutil
#cgo pkg-config: libavformat libavcodec libavutil
*/
```

## Testing

After building, test with:

```bash
# Linux/macOS
./rtsp_analyzer_linux "rtsp://username:password@camera_ip:port/path"

# Windows
rtsp_analyzer_go.exe "rtsp://username:password@camera_ip:port/path"
```

## Troubleshooting

### Common Issues:

1. **FFmpeg not found**: Make sure FFmpeg development libraries are installed
2. **CGO compilation errors**: Verify that gcc/build tools are installed
3. **pkg-config errors**: Install pkg-config (`sudo apt-get install pkg-config` on Ubuntu)
4. **Library linking errors**: Check that FFmpeg libraries are in the system library path

### Checking FFmpeg Installation:

```bash
# Check if pkg-config finds FFmpeg
pkg-config --exists libavformat libavcodec libavutil && echo "Found" || echo "Not found"

# Check FFmpeg version
pkg-config --modversion libavformat

# Check include paths
pkg-config --cflags libavformat libavcodec libavutil

# Check library paths
pkg-config --libs libavformat libavcodec libavutil
```

## Cross-Compilation

To cross-compile for different platforms:

```bash
# Linux to Windows (requires mingw-w64)
GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc go build

# Note: Cross-compilation with CGO requires platform-specific toolchains
```
