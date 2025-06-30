# ONVIF Manager - User Guide

ONVIF Manager is a comprehensive tool for managing, configuring, and validating ONVIF-compatible IP cameras. This guide explains how to set up, build, and use the ONVIF Manager through its command-line interface (CLI) and web application.

## Table of Contents
- [Prerequisites](#prerequisites)
- [Building the Application](#building-the-application)
  - [Standard Build](#standard-build)
  - [Building with Embedded Frontend](#building-with-embedded-frontend)
- [Usage](#usage)
  - [CLI Mode](#cli-mode)
  - [Web Application Mode](#web-application-mode)
  - [API Server Mode](#api-server-mode)
- [Command Reference](#command-reference)
- [CSV File Formats](#csv-file-formats)
- [Examples](#examples)

## Prerequisites

Before you can build and run ONVIF Manager, ensure you have the following prerequisites installed:

1. **Go** (version 1.20 or later)
   - Download from: https://go.dev/dl/
   - Verify with: `go version`

2. **C Compiler** (for CGO dependencies)
   - Windows: MinGW or MSVC with Windows SDK
   - Linux: GCC (`sudo apt install build-essential`)
   - macOS: Xcode Command Line Tools (`xcode-select --install`)

3. **FFmpeg** (for stream validation)
   - Windows: https://ffmpeg.org/download.html#build-windows
   - Linux: `sudo apt install ffmpeg`
   - macOS: `brew install ffmpeg`
   - Verify with: `ffmpeg -version`

4. **VLC Media Player** (for camera stream viewing)
   - Windows: https://www.videolan.org/vlc/download-windows.html
   - Linux: `sudo apt install vlc`
   - macOS: `brew install --cask vlc` or download from https://www.videolan.org/vlc/download-macosx.html
   - Verify with: `vlc --version`

Note: The web frontend is already built and included in the repository, so Node.js and npm are not required unless you want to modify the frontend.

## Building the Application

The ONVIF Manager application can be built on both Windows and Linux platforms. The resulting binary will include both the backend API and the pre-built web frontend.

### Building on Windows

```bash
# Navigate to the project directory
cd onvif_manager

# Build the application
go build -o onvif-manager.exe cmd/app/main.go

# Verify the build
onvif-manager.exe --help
```

### Building on Linux

```bash
# Navigate to the project directory
cd onvif_manager

# Build the application
go build -o onvif-manager cmd/app/main.go

# Make the binary executable
chmod +x onvif-manager

# Verify the build
./onvif-manager --help
```

The frontend components are already compiled and included in the repository, so you don't need to build them separately.

Note: If you've made changes to the frontend code in the main_front directory, you'll need to rebuild it using npm and copy the output to the appropriate directory in the onvif_manager project before building.

## Usage

### CLI Mode

Run the application in CLI mode to perform specific operations:

```bash
# Apply camera configuration
onvif-manager.exe config apply cameras.csv config_1080p.csv

# Other commands (hidden in simplified workflow)
# onvif-manager.exe list
# onvif-manager.exe select cameras.csv
# onvif-manager.exe export results.csv
```

### Web Application Mode

Run the application as a web application:

```bash
onvif-manager.exe web
```

This starts a web server at http://localhost:8090 with both the frontend interface and API endpoints.

### API Server Mode

Run only the API server without the frontend:

```bash
onvif-manager.exe server
```

This starts the API server at http://localhost:8090, which can be accessed by any custom frontend or through API calls.

## Command Reference

### Main Commands

- **web**: Start combined web application (frontend + API) on port 8090
  ```
  onvif-manager.exe web
  ```

- **server**: Start API server only on port 8090
  ```
  onvif-manager.exe server
  ```

- **config apply**: Import cameras and apply configuration
  ```
  onvif-manager.exe config apply [camera-csv] [config-csv]
  ```



## CSV File Formats

### Camera CSV Format
```
id,ip,port,url,username,password
1,192.168.10.101,80,,admin,MySecurePass1
2,192.168.10.102,0,,operator,CameraPass#2
3,192.168.10.103,80,,administrator,Secure$Pass3
4,192.168.10.104,0,,root,TestDevice#4
```

### Configuration CSV Format
```
width,height,fps,bitrate
1920,1080,30,4000
```

### Validation Results CSV Format
```
cam_id,cam_ip,result,reso_expected,reso_actual,fps_expected,fps_actual,notes
1,192.168.1.100,PASS,1920x1080,1920x1080,30,30.00,All parameters match expected values
2,192.168.1.101,FAIL,1920x1080,1280x720,30,25.00,Resolution mismatch
```

## Examples

### Example 1: Configuring Multiple Cameras

```bash
# Configure cameras with 1080p settings
onvif-manager.exe config apply examples/cameras.csv examples/config_1080p.csv
```

### Example 2: Starting the Web Application

```bash
# Start the web application
onvif-manager.exe web
```

Then navigate to http://localhost:8090 in your web browser.

---

For additional help or to report issues, please contact your system administrator or refer to the project documentation.

© 2025 ONVIF Manager Project
