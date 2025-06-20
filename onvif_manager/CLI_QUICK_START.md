# ONVIF Manager CLI - Quick Start Guide

## Overview
The ONVIF Manager CLI provides a comprehensive command-line interface for managing ONVIF cameras with persistent configuration support.

## Available Commands

### 1. List Cameras
```bash
onvif-manager list
```
Shows all cameras in the system with their details.

### 2. Select Cameras from CSV
```bash
onvif-manager select cameras.csv
```
Selects cameras from a CSV file containing IP addresses. The CSV format should be:
```csv
ip
192.168.1.12
192.168.1.30
192.168.1.31
```

### 3. Configuration Management

#### Show Current Saved Configuration
```bash
onvif-manager config show
```

#### Set Configuration Manually
```bash
onvif-manager config set 1920 1080 30 4096
```
Arguments: width height fps bitrate

#### Import Configuration from CSV
```bash
onvif-manager config import config.csv
```
CSV format should be:
```csv
width,height,fps,bitrate
1920,1080,30,4096
```

#### Apply Configuration

##### Apply with both camera and config CSV files
```bash
onvif-manager config apply cameras.csv config.csv
```

##### Apply saved configuration to selected cameras
```bash
onvif-manager config apply-to cameras.csv
```

### 4. Export Results
```bash
onvif-manager export results.csv
```
Exports validation results after applying configuration.

## Workflow Examples

### Example 1: Simple Configuration Application
```bash
# 1. Set configuration manually
onvif-manager config set 1280 720 25 2048

# 2. Apply to selected cameras
onvif-manager config apply-to cameras.csv

# 3. Export results
onvif-manager export validation_results.csv
```

### Example 2: Import and Apply from CSV
```bash
# 1. Import configuration from CSV
onvif-manager config import high_quality_config.csv

# 2. View imported configuration
onvif-manager config show

# 3. Apply to cameras
onvif-manager config apply-to target_cameras.csv
```

### Example 3: Direct Application
```bash
# Apply specific config to specific cameras in one command
onvif-manager config apply cameras.csv config.csv
```

## Features

- **Persistent Configuration**: The system maintains a saved configuration that persists between CLI sessions
- **Source Tracking**: Tracks whether config came from CSV import, manual entry, or default values
- **User Confirmation**: Interactive prompts before applying potentially disruptive changes
- **Comprehensive Results**: Detailed success/failure reporting with validation results
- **Export Capability**: Export validation results to CSV for analysis
- **Robust Error Handling**: Clear error messages and helpful guidance

## File Locations

- **Camera Configuration**: `config/cameras.json`
- **Saved Configuration**: `config/saved_config.json`
- **Example Files**: `examples/` directory

The CLI automatically finds the correct config directory whether you run it with `go run` or as a compiled executable.

## CSV File Examples

All example CSV files are available in the `examples/` directory:
- `cameras_to_configure.csv` - Example camera selection
- `config_1080p.csv` - 1080p configuration
- `config_720p.csv` - 720p configuration
- `new_cameras.csv` - Additional camera examples
