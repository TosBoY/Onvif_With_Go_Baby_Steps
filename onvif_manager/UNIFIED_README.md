# ONVIF Manager - CLI and API Application

A command-line interface (CLI) application for ONVIF camera management, with an optional API server for integration with other applications.

## 🚀 Quick Start

### Download and Run
1. Download the `onvif-manager.exe` binary
2. Choose your preferred mode:

**CLI Mode** (with CLI commands):
```bash
./onvif-manager.exe list
./onvif-manager.exe config show
./onvif-manager.exe config apply cameras.csv config.csv
```

**API Server Mode** (for integration with other apps):
```bash
./onvif-manager.exe server
# Starts API server at http://localhost:8090
```

## 📋 Application Features

### 🖥️ CLI Operations
Perfect for automation, scripting, and advanced users.

**Benefits**: 
- Fast execution
- Scriptable operations
- Batch processing
- SSH-friendly
- No GUI dependencies

**Examples**:
```bash
# List all cameras
./onvif-manager.exe list

# Show saved configuration
./onvif-manager.exe config show

# Set configuration manually
./onvif-manager.exe config set 1920 1080 30 4096

# Apply configuration to cameras from CSV
./onvif-manager.exe config apply-to cameras.csv

# Full workflow with two CSV files
./onvif-manager.exe config apply cameras.csv config.csv
```

### 🌐 API Server Mode
Provides REST API access to all backend functionality for integration with other applications.

**Activation**: `./onvif-manager.exe server`
**Access**: http://localhost:8090
**Benefits**:
- REST API endpoints
- Integration with custom frontends
- Programmatic access
- Cross-platform compatibility

**API Endpoints**:
- `GET /cameras` - List all cameras
- `POST /cameras` - Add new camera
- `DELETE /cameras/{id}` - Remove camera
- `POST /cameras/import-csv` - Import cameras from CSV
- `POST /import-config-csv` - Import configuration from CSV
- `POST /apply-config` - Apply configuration to cameras
- `POST /export-validation-csv` - Export validation results
- `POST /vlc` - Launch VLC for camera stream

## 🔧 Technical Details

### Architecture
- **CLI Application**: Primary interface for command-line operations
- **API Server**: Optional REST API server for integration
- **Shared Backend**: Core functionality shared between CLI and API
- **Camera Management**: Automatic camera connection and initialization
- **Persistent Configuration**: Saved configuration persists between sessions
- **CSV Integration**: Import/export functionality for batch operations

### File Structure
```
onvif_manager/
├── cmd/app/main.go           # Application entry point
├── internal/
│   ├── cli/                  # CLI functionality
│   └── backend/              # Shared backend logic
│       ├── api/              # REST API handlers (for server mode)
│       ├── camera/           # Camera management
│       ├── config/           # Configuration handling
│       ├── ffmpeg/           # Video validation
│       └── vlc/              # Media player integration
├── config/                   # Configuration files
├── examples/                 # Example CSV files
└── onvif-manager.exe         # Built binary
```

## � Migration from Previous Versions

If you were previously using separate CLI and web applications:

### From CLI-only setup:
- Replace your CLI executable with this binary
- All CLI commands work exactly the same
- No changes needed to existing scripts

### From Web-based setup:
- Migrate to CLI workflows using CSV import/export
- All camera configurations remain compatible
- Existing data formats supported

### Configuration Compatibility
- Uses the same `cameras.json` format
- Same `saved_config.json` for persistent settings  
- CSV formats remain unchanged
- All existing workflows continue to work

## 🛠️ Development and Building

### Prerequisites
- Go 1.24.2+
- Required dependencies in `go.mod`

### Building
```bash
go build -o onvif-manager.exe cmd/app/main.go
```

## 🔍 Troubleshooting

### CLI Issues
- **Commands not working**: Ensure you're passing arguments
- **Cameras not found**: Run initialization commands first
- **Permission errors**: Check file permissions for config directory
- **Connection timeouts**: Verify network connectivity and camera credentials

### General Issues
- **Config file errors**: Verify `cameras.json` format and permissions
- **Missing features**: Ensure all dependencies are installed
- **Performance issues**: Consider reducing camera count for testing
- **CSV format errors**: Check example files for proper formatting

## 📝 Examples and Use Cases

### Automated Deployment
```bash
# Script for automated camera configuration
./onvif-manager.exe config import production_config.csv
./onvif-manager.exe config apply-to production_cameras.csv
./onvif-manager.exe export results_$(date +%Y%m%d).csv
```

### Batch Configuration
```bash
# Configure multiple cameras with single command
./onvif-manager.exe config apply cameras.csv config.csv

# Export results for verification
./onvif-manager.exe export validation_results.csv
```

### Camera Discovery and Management
```bash
# List all configured cameras
./onvif-manager.exe list

# Show current configuration
./onvif-manager.exe config show

# Test camera connectivity
./onvif-manager.exe select cameras.csv
```

This CLI-focused approach provides efficient camera management for automation and scripting scenarios while maintaining compatibility with existing workflows.
