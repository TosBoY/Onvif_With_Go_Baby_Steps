# ONVIF Manager - Complete Command Guide

A comprehensive tool for managing ONVIF cameras with CLI interface, API server, and web application capabilities.

## 🚀 Quick Start

### Build Options

**Option 1: Standard Build (CLI + API Server)**
```bash
cd onvif_manager
go build -o onvif-manager.exe cmd/app/main.go
```

**Option 2: Embedded Build (CLI + API Server + Web Frontend)**
```bash
cd onvif_manager
# Run the automated build script
build-embedded.bat

# Or manual build:
# 1. Build frontend: cd ../main_front && npm run build
# 2. Copy files: xcopy "..\main_front\dist" "cmd\app\web" /E /I /Y
# 3. Build binary: go build -o onvif-manager-embedded.exe cmd/app/main.go cmd/app/webserver.go
```

## 🌐 Server Modes

### Web Application Mode (Recommended)
Starts combined frontend and API server on port 8090:
```bash
./onvif-manager-embedded.exe web
```
- **Frontend**: http://localhost:8090
- **API**: http://localhost:8090/api/*
- Full GUI interface for camera management

### API Server Only Mode
Starts API server only on port 8090:
```bash
./onvif-manager.exe server
```
- **API**: http://localhost:8090/*
- For programmatic access or custom frontends

## 📋 CLI Commands

### Camera Management

#### 1. List All Cameras
Display all cameras currently configured in the system:
```bash
./onvif-manager.exe list
```
**Output:**
- Camera ID, IP, Port, Username, URL, IsFake status
- Total count of cameras

#### 2. Select Cameras from CSV
Select specific cameras using a CSV file containing IP addresses:
```bash
./onvif-manager.exe select cameras_to_select.csv
```

**CSV Format:**
```csv
ip
192.168.1.12
192.168.1.30
192.168.1.31
```

**Output:**
- Shows matched cameras and unmatched IPs
- Reports invalid rows and statistics

### Configuration Management

#### 3. Apply Configuration (Two CSV Files)
Apply configuration using separate camera selection and config files:
```bash
./onvif-manager.exe config apply cameras_to_configure.csv config_settings.csv
```

**Camera Selection CSV (`cameras_to_configure.csv`):**
```csv
ip
192.168.1.12
192.168.1.30
192.168.1.31
```

**Configuration CSV (`config_settings.csv`):**
```csv
width,height,fps,bitrate
1920,1080,30,4096
```

#### 4. Show Current Saved Configuration
Display the currently saved configuration:
```bash
./onvif-manager.exe config show
```

#### 5. Set Configuration Manually
Set configuration values manually:
```bash
./onvif-manager.exe config set 1920 1080 30 4096
```
**Parameters:** width height fps bitrate

#### 6. Import Configuration from CSV
Import and save configuration from a CSV file:
```bash
./onvif-manager.exe config import config_1080p.csv
```

**Config CSV Format:**
```csv
width,height,fps,bitrate
1920,1080,30,4096
```

#### 7. Apply Saved Configuration to Selected Cameras
Apply the current saved configuration to cameras from CSV:
```bash
./onvif-manager.exe config apply-to cameras_to_configure.csv
```

### Results Management

#### 8. Export Validation Results
Export the last validation results to CSV:
```bash
./onvif-manager.exe export validation_results.csv
```

**Output CSV Format:**
```csv
cam_id,cam_ip,result,reso_expected,reso_actual,fps_expected,fps_actual
1,192.168.1.12,PASS,1920x1080,1920x1080,30,30.00
2,192.168.1.30,FAIL,1920x1080,1280x720,30,25.00
```

## 📁 CSV File Examples

### Camera Selection Files

**Example: `cameras_to_configure.csv`**
```csv
ip
192.168.1.12
192.168.1.30
192.168.1.31
192.168.1.45
```

### Configuration Files

**Example: `config_1080p.csv`**
```csv
width,height,fps,bitrate
1920,1080,30,4096
```

**Example: `config_720p.csv`**
```csv
width,height,fps,bitrate
1280,720,25,2048
```

**Example: `config_4k.csv`**
```csv
width,height,fps,bitrate
3840,2160,30,8192
```

## 🔄 Complete Workflow Examples

### Workflow 1: Direct Configuration Application
```bash
# 1. List available cameras
./onvif-manager.exe list

# 2. Apply configuration directly
./onvif-manager.exe config apply cameras_to_configure.csv config_1080p.csv

# 3. Export results (optional)
./onvif-manager.exe export results_20240623.csv
```

### Workflow 2: Using Saved Configuration
```bash
# 1. Import and save configuration
./onvif-manager.exe config import config_1080p.csv

# 2. Verify saved configuration
./onvif-manager.exe config show

# 3. Apply to different camera groups
./onvif-manager.exe config apply-to group1_cameras.csv
./onvif-manager.exe config apply-to group2_cameras.csv

# 4. Export results
./onvif-manager.exe export validation_results.csv
```

### Workflow 3: Camera Selection and Testing
```bash
# 1. Test camera selection first
./onvif-manager.exe select cameras_to_test.csv

# 2. If selection looks good, apply configuration
./onvif-manager.exe config apply cameras_to_test.csv config_settings.csv
```

### Workflow 4: Manual Configuration
```bash
# 1. Set configuration manually
./onvif-manager.exe config set 1920 1080 25 3072

# 2. Apply to cameras
./onvif-manager.exe config apply-to cameras.csv

# 3. Check and export results
./onvif-manager.exe export manual_config_results.csv
```

## ⚙️ Configuration Parameters

| Parameter | Description | Example Values | Notes |
|-----------|-------------|----------------|-------|
| **Width** | Video width in pixels | 1920, 1280, 3840 | Must match camera capabilities |
| **Height** | Video height in pixels | 1080, 720, 2160 | Must match camera capabilities |
| **FPS** | Frames per second | 30, 25, 15, 10 | Higher values require more bandwidth |
| **Bitrate** | Bitrate in kbps | 4096, 2048, 8192 | Set to 0 for auto/default |

## 📊 Validation Results

The system validates each configuration by:

1. **Connectivity Check**: Verifies camera responds to ONVIF commands
2. **Resolution Validation**: Confirms applied resolution matches expected
3. **Performance Check**: Validates FPS and bitrate settings
4. **Error Reporting**: Provides detailed error messages for failures

**Result Types:**
- **PASS**: All settings applied and validated successfully
- **WARNING**: Settings applied but some parameters differ (resolution matches)
- **FAIL**: Configuration failed or critical parameters don't match

## 🔧 Advanced Usage

### Environment Setup
```bash
# Ensure cameras.json is configured
cp config/cameras.json.example config/cameras.json
# Edit cameras.json with your camera details
```

### Batch Operations
```bash
# Apply different configs to different groups
./onvif-manager.exe config apply group1.csv config_1080p.csv
./onvif-manager.exe config apply group2.csv config_720p.csv
./onvif-manager.exe config apply group3.csv config_4k.csv
```

### Testing and Validation
```bash
# Test camera connectivity first
./onvif-manager.exe select all_cameras.csv

# Apply test configuration
./onvif-manager.exe config apply all_cameras.csv test_config.csv

# Export detailed results
./onvif-manager.exe export full_validation_$(date +%Y%m%d).csv
```

## 🚨 Error Handling

Common issues and solutions:

**Network Issues:**
- Ensure cameras are accessible on the network
- Check IP addresses in CSV files
- Verify ONVIF port accessibility (usually 80 or 8080)

**Configuration Errors:**
- Verify camera supports requested resolution/fps
- Check bitrate compatibility
- Ensure CSV format is correct

**File Issues:**
- Verify CSV files exist and are readable
- Check CSV headers match expected format
- Ensure no extra spaces or special characters

## 📱 Web Interface

When using the embedded build with `web` mode:

1. **Camera Management**: Add, remove, and configure cameras via GUI
2. **Bulk Operations**: Import cameras and configurations via CSV upload
3. **Live Validation**: Real-time configuration validation and results
4. **Export Features**: Download validation results as CSV
5. **VLC Integration**: Launch VLC for stream viewing

## 🛠️ Development and Building

### Frontend Development
```bash
cd main_front
npm install
npm run dev          # Development mode with hot reload
npm run build        # Production build
```

### Backend Development
```bash
cd onvif_manager
go mod tidy
go run cmd/app/main.go [commands]
```

### Creating New Builds
```bash
# Update embedded frontend
cd onvif_manager
./build-embedded.bat

# Standard build
go build -o onvif-manager.exe cmd/app/main.go
```

## 📝 Notes

- **Simulated Cameras**: Cameras marked as `IsFake: true` will show successful results without actual network calls
- **File Paths**: All CSV file paths are relative to the executable location
- **Port Configuration**: Default port is 8090 for both web and API modes
- **CORS**: API server includes CORS headers for cross-origin requests
- **Validation Storage**: Last validation results are automatically stored for export

## 🆘 Help and Support

### Get Command Help
```bash
./onvif-manager.exe help
./onvif-manager.exe help config
./onvif-manager.exe config help apply
```

### Usage Summary
```bash
# Show all available commands
./onvif-manager.exe
```

This guide covers all functionality of the ONVIF Manager tool in both CLI and web modes.
