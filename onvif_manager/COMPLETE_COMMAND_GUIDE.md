# ONVIF Manager - Complete Command Guide

A comprehensive tool for managing ONVIF cameras with an ultra-simplified CLI workflow, API server, and web application capabilities. This version uses in-memory camera management without requiring persistent configuration files and features a single-command approach to import cameras and apply configuration.

## 🚀 Quick Start - Ultra Simplified CLI Workflow

```bash
# The only CLI command you need
./onvif-manager.exe config apply examples/cameras.csv examples/config_1080p.csv
```

This single command:
1. Imports all cameras from cameras.csv (with details like IP, username, password)
2. Applies configuration settings from config_1080p.csv
3. Prompts to export validation results to CSV file

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
# 2. Copy files: xcopy "..\main_front\dist" "internal\webserver\web" /E /I /Y
# 3. Build binary: go build -o onvif-manager-embedded.exe cmd/app/main.go
```

## 🌐 Server Modes

### Web Application Mode (Recommended)
Starts combined frontend and API server on port 8090:
```bash
./onvif-manager.exe web
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

## 📋 CLI Command

### The Only Command You Need

```bash
./onvif-manager.exe config apply cameras.csv config_1080p.csv
```

**Process:**
1. Imports cameras from first CSV file
2. Loads configuration from second CSV file
3. Applies configuration to all imported cameras
4. Validates results
5. Prompts to export results to CSV

**First CSV (cameras.csv) Format:**
```csv
ip,username,password,port,url,isfake
192.168.1.12,admin,admin123,80,,false
192.168.1.30,admin,admin123,80,,false
192.168.1.31,admin,admin123,80,,true
```

**Second CSV (config_1080p.csv) Format:**
```csv
width,height,fps,bitrate
1920,1080,30,4096
```

### Server Commands

#### 3. Start Web Server
Start the combined web interface and API server:
```bash
./onvif-manager.exe web
```
**Features:**
- Web interface available at http://localhost:8090
- API endpoints at http://localhost:8090/api/*
- Interactive camera management and configuration

#### 4. Start API Server Only
Start just the API server without the web interface:
```bash
./onvif-manager.exe server
```
**Features:**
- API endpoints available at http://localhost:8090/*
- Programmatic access for custom clients

### Results Format

When the command completes and you choose to export validation results, the CSV will have this format:

**Config CSV Format (Second Parameter):**
```csv
width,height,fps,bitrate
1920,1080,30,4096
```

### Results Management

After running the command, you'll be prompted to export validation results:
```
💾 Do you want to export validation results to CSV? (y/N):
```

If you select 'y', you'll be prompted to enter a filename:
```
📄 Enter output filename (default: validation_results_20250624_153000.csv):
```

**Output CSV Format:**
```csv
cam_id,cam_ip,result,reso_expected,reso_actual,fps_expected,fps_actual
1,192.168.1.12,PASS,1920x1080,1920x1080,30,30.00
2,192.168.1.30,FAIL,1920x1080,1280x720,30,25.00
```

## 📁 CSV File Examples

### Adding Cameras to the System

Cameras can be added to the system in multiple ways:

**Option 1: Using the CLI Command**
```bash
./onvif-manager.exe config apply cameras.csv config_1080p.csv
```
This will import cameras and immediately apply configuration in one step.

**Option 2: Using the Web Interface**
1. Navigate to http://localhost:8090 after starting with `./onvif-manager.exe web`
2. Use the "Add Camera" function to add cameras individually
3. Or use the "Import CSV" function with a CSV file

**CSV Format for Camera Import:**
```csv
ip,username,password,port,url,isfake
192.168.1.12,admin,admin123,80,,false
192.168.1.30,admin,admin123,80,,false
192.168.1.31,admin,admin123,80,,true
```

Required columns:
- `ip`: Camera IP address
- `username`: Camera login username

Optional columns:
- `port`: ONVIF port (default: 80)
- `url`: Custom stream URL (leave empty for auto-detection)
- `password`: Camera login password
- `isfake`: Set to `true` for simulated cameras (default: false)

### CSV File Types

#### Camera CSV File (First Parameter)

Used as the first parameter to `config apply` command.

**Example: `cameras.csv`**
```csv
ip,username,password,port,url,isfake
192.168.1.12,admin,admin123,80,,false
192.168.1.30,admin,admin123,80,,false
192.168.1.31,admin,admin123,80,,true
```

#### Camera Selection Files

Used to select cameras that are already in the system for operations like configuration application.

**Example: `cameras_to_configure.csv`**
```csv
ip
192.168.1.12
192.168.1.30
192.168.1.31
192.168.1.45
```

Note: Selection files don't add cameras to the system; they only specify which existing cameras to operate on.

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

### Workflow 1: Simple One-Command Operation

```bash
# Single command to import cameras and apply configuration
./onvif-manager.exe config apply examples/cameras.csv examples/config_1080p.csv
```

This command will:
- Import all cameras from cameras.csv
- Apply configuration from config_1080p.csv to all imported cameras
- Validate the configuration on each camera
- Prompt to export results to CSV if desired

### Workflow 2: Different Configuration Profiles

```bash
# Apply 1080p configuration to first group
./onvif-manager.exe config apply group1_cameras.csv config_1080p.csv

# Apply 720p configuration to second group
./onvif-manager.exe config apply group2_cameras.csv config_720p.csv

# Apply 4K configuration to third group
./onvif-manager.exe config apply group3_cameras.csv config_4k.csv
```

### Workflow 3: Web Interface Management

```bash
# Start the web server
./onvif-manager.exe web

# Access web interface at http://localhost:8090
# - Import cameras using the web interface
# - Apply configurations through the UI
# - View and export results
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

### Working with Fake Cameras
You can create simulated cameras for testing by setting `isfake` to `true` in your CSV:

```csv
ip,username,password,isfake
192.168.1.100,admin,admin123,true
192.168.1.101,admin,admin123,true
```

Fake cameras will:
- Always report successful configuration
- Simulate expected validation results
- Not attempt to make actual network connections

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
./onvif-manager.exe export full_validation_results.csv
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

- **Simplified CLI Flow**: The main CLI workflow is a single command that imports cameras and applies configuration
- **In-Memory Storage**: All camera and configuration data is stored in memory only
- **No Persistence**: Camera data will be lost when the application is restarted
- **Single Process Execution**: The `config apply` command handles the entire workflow in a single process
- **Camera CSV Format**: The first CSV must contain full camera details (IP, username, password, etc.)
- **Configuration CSV Format**: The second CSV must contain configuration parameters (width, height, fps, bitrate)
- **Simulated Cameras**: Cameras marked as `IsFake: true` will show successful results without actual network calls
- **File Paths**: All CSV file paths are relative to the executable location
- **Port Configuration**: Default port is 8090 for both web and API modes
- **Web Alternative**: For interactive management, use `./onvif-manager.exe web` and access http://localhost:8090

## 🔀 CLI Usage Approaches

Due to the in-memory design, we've simplified the CLI to a single command workflow:

### Simplified CLI Approach (Recommended)
The best approach is to use one command that does everything:

```bash
./onvif-manager.exe config apply cameras.csv config_1080p.csv
```

This single command:
1. Imports all cameras from the first CSV file
2. Applies configuration settings from the second CSV file 
3. Validates the applied configuration
4. Prompts to export results to a CSV file

### Web Server Alternative
For more interactive management, you can use the web interface:

```bash
# Start the web server
./onvif-manager.exe web
```

Then access the web interface at http://localhost:8090 to:
- Add/import cameras
- Apply configurations
- View and export results

### API Endpoints
The API endpoints can also be accessed programmatically:
```
POST   http://localhost:8090/api/cameras     # Add or import cameras
GET    http://localhost:8090/api/cameras     # List cameras  
POST   http://localhost:8090/api/config      # Apply configuration
```

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
