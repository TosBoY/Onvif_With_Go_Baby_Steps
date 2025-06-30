# Development Guide – ONVIF Camera Management System

## Project Overview

This is a comprehensive ONVIF camera management system that provides a web-based interface for discovering, configuring, and monitoring IP cameras. The system consists of a Go backend for ONVIF communication and a React frontend for user interaction.

## Project Structure

- **main_back/**: Go backend (REST API, ONVIF camera logic, stream validation)
- **main_front/**: React frontend (Web UI, camera management interface)
- **onvif_manager/**: Standalone ONVIF management tools
- **tests/**: Test suites and validation tools
- **past_implementation/**: Previous versions and experimental code

## Quick Start

**Backend:**
```cmd
cd main_back
go run cmd/backend/main.go
```

**Frontend:**
```cmd
cd main_front
npm install
npm run dev
```

**Access Points:**
- Frontend UI: http://localhost:5173
- Backend API: http://localhost:8090

## Development Workflow

### Frontend (main_front)
The frontend is built with React 19 and Material-UI v7 for modern component styling.

**Key Pages:**
- **Dashboard.jsx**: Main control interface with camera discovery, configuration, and CSV management
- **CameraStatus.jsx**: Camera status monitoring page with real-time status checking and pagination

**Key Components:**
- **CameraCard.jsx**: Individual camera display with VLC launch and info controls
- **CameraConfigPanel.jsx**: Batch configuration panel for camera resolution and settings
- **CameraInfoDialog.jsx**: Detailed camera information modal with delete functionality
- **ValidationResults.jsx**: Stream validation results display
- **Header.jsx**: Navigation header with page routing
- **ConnectionStatus.jsx**: Backend connection status indicator

**Core Files:**
- `src/components/`: Reusable React components
- `src/pages/`: Page-level components (Dashboard, CameraStatus)
- `src/services/api.js`: Backend API communication layer
- `src/App.jsx`: Main application with React Router routing
- `index.html`: Main HTML template (title: "ONVIF Camera Manager")

**Features:**
- Camera discovery and CSV import/export
- Real-time camera status monitoring with pagination (50 cameras per page)
- Batch camera configuration with validation
- VLC integration for camera stream viewing
- Responsive Material-UI design

### Backend (main_back)
The backend is built with Go 1.24+ and provides RESTful APIs for ONVIF camera management.

**Key Components:**
- **handlers.go**: REST API endpoints for all camera operations
- **manager.go**: Core camera management and ONVIF communication
- **client.go**: ONVIF client implementation with device discovery
- **resolution.go**: Camera resolution management and validation

**Core Directories:**
- `cmd/backend/`: Application entry point and main server
- `internal/api/`: HTTP handlers and route definitions
- `internal/camera/`: Camera management, ONVIF client, configuration logic
- `internal/ffmpeg/`: Stream analysis and validation using FFmpeg/FFprobe
- `internal/loader/`: CSV file handling for camera data
- `internal/vlc/`: VLC media player integration
- `pkg/models/`: Data structures and models

**Data Storage:**
- `internal/loader/cameras.csv`: Camera configuration database (CSV format)

**Dependencies:**
- `github.com/gorilla/mux`: HTTP routing
- `github.com/gorilla/handlers`: CORS and middleware
- `github.com/videonext/onvif`: ONVIF protocol implementation

## API Endpoints

The backend provides RESTful APIs for camera management:

- `GET /cameras` - List all cameras from CSV
- `POST /cameras` - Add new camera to CSV
- `DELETE /cameras/{id}` - Remove camera from CSV
- `POST /discover` - Discover ONVIF cameras on network
- `POST /check/{id}` - Check single camera status
- `POST /configure/{id}` - Configure single camera settings
- `POST /validate/{id}` - Validate camera stream using FFprobe
- `POST /vlc/{id}` - Launch VLC for camera stream
- `POST /export-csv` - Export camera data to CSV
- `POST /configure-batch` - Configure multiple cameras

## Testing

### Frontend Testing
- **Browser DevTools**: Use console and network tabs for debugging
- **React DevTools**: Inspect component state and props
- **Debug Route**: Access `/debug` for development testing
- **Component Testing**: Individual component validation

### Backend Testing
- **API Testing**: Use curl, Postman, or similar tools for endpoint testing
- **Go Testing**: Run `go test ./...` in main_back directory
- **Integration Testing**: Full workflow testing with frontend and backend
- **ONVIF Testing**: Test with real cameras or ONVIF simulator

### Test Commands
```cmd
# Backend tests
cd main_back
go test ./...

# Frontend linting
cd main_front
npm run lint

# Build verification
cd main_front
npm run build
```

## Debugging

### Frontend Debugging
- **Browser Console**: Check for JavaScript errors and debug output
- **Network Tab**: Inspect API requests and responses
- **React DevTools**: Examine component state and props
- **Source Maps**: Debug TypeScript/JSX code directly in browser
- **Console Logging**: Use `console.log` for data flow tracking

### Backend Debugging
- **Terminal Output**: Monitor detailed logging information
- **Log Files**: Capture logs with `go run cmd/backend/main.go 2>&1 | tee backend.log`
- **IDE Debugging**: Set breakpoints in VSCode or GoLand
- **API Debugging**: Test endpoints with curl or Postman
- **ONVIF Debugging**: Monitor ONVIF communication and errors

### Common Issues
- **CORS Errors**: Ensure backend CORS is properly configured
- **Camera Discovery**: Check network connectivity and ONVIF support
- **Stream Validation**: Verify FFmpeg/FFprobe installation
- **Port Conflicts**: Ensure ports 5173 (frontend) and 8090 (backend) are available

## Camera Configuration Process

The system implements a sophisticated two-phase approach to camera configuration:

### Phase 1 - Apply Configuration
- Applies settings to all selected cameras simultaneously
- Determines the closest supported resolution for each camera
- Configures FPS, bitrate, and encoding parameters
- Prepares stream URLs for subsequent validation

### Phase 2 - Stream Validation
- After a brief pause (1 second), validates all cameras in the original order
- Uses FFprobe to verify that the applied settings are functioning correctly
- Checks stream integrity, resolution, and encoding parameters
- For test/fake cameras, simulates successful validation

### Configuration Features
- **Batch Configuration**: Configure multiple cameras simultaneously
- **Resolution Mapping**: Automatically maps requested resolutions to camera capabilities
- **Stream Validation**: Real-time validation using FFmpeg/FFprobe
- **Status Monitoring**: Real-time status updates with pagination (50 cameras per page)
- **Error Handling**: Comprehensive error reporting and recovery

This batch processing approach ensures efficient configuration of multiple cameras while maintaining system responsiveness.

## Camera Status Monitoring

The Camera Status page provides comprehensive monitoring capabilities:

- **Real-time Status**: Check individual or all camera statuses
- **Pagination**: Browse cameras in groups of 50 for better performance
- **Stream Validation**: Validate RTSP streams using FFprobe
- **VLC Integration**: Launch VLC player for camera streams
- **Configuration Dialog**: Modify camera settings with real-time feedback

## CSV Data Management

The system uses CSV files for camera data persistence:

- **Import/Export**: Import camera lists and export configurations
- **Auto-save**: Automatic saving of camera configurations
- **Backup**: Regular backups recommended for data preservation
- **Format**: Standardized CSV format for camera definitions

## Maintenance

### Data Backup
- **Camera Configuration**: Backup `main_back/internal/loader/cameras.csv` regularly
- **System State**: Preserve camera settings and configurations
- **Version Control**: Use Git for code versioning and history

### Updates and Dependencies
- **Frontend Dependencies**: Run `npm update` in main_front directory
- **Backend Dependencies**: Update Go modules with `go mod tidy`
- **ONVIF Library**: Check for updates to `github.com/videonext/onvif`
- **Security Updates**: Keep all dependencies up to date

### Performance Optimization
- **Pagination**: Camera lists are paginated at 50 items for optimal performance
- **Batch Operations**: Use batch configuration for multiple cameras
- **Connection Pooling**: Backend manages ONVIF connections efficiently
- **Caching**: Status information is cached to reduce network load

### Troubleshooting
- **Connection Issues**: Verify network connectivity and firewall settings
- **Camera Discovery**: Ensure cameras are ONVIF compliant
- **Stream Problems**: Check FFmpeg/FFprobe installation and paths
- **Performance**: Monitor system resources during large camera operations

---

_Last updated: June 30, 2025_
