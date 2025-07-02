# Development Guide – ONVIF Camera Management System

## Project Overview

This is a comprehensive ONVIF camera management system that provides multiple interfaces for discovering, configuring, and monitoring IP cameras. The system uses a development workflow where changes are prototyped in separate testbed applications and then integrated into the production-ready unified application.

## Project Structure

- **main_back/**: Development testbed - Go backend (REST API, ONVIF camera logic, stream validation)
- **main_front/**: Development testbed - React frontend (Web UI, camera management interface)
- **onvif_manager/**: Production application - Unified CLI, API server, and web interface
- **tests/**: Test suites and validation tools
- **past_implementation/**: Previous versions and experimental code

## Development Workflow

This project follows a **testbed-to-production** development approach:

1. **Prototype in Testbed**: New features and changes are first developed and tested in `main_back/` and `main_front/`
2. **Validate Changes**: Test the functionality thoroughly in the isolated testbed environment
3. **Integrate to Production**: Successful changes are then ported to the unified `onvif_manager/` application
4. **Deploy Production**: The `onvif_manager/` serves as the production-ready application

### Why This Approach?
- **Isolated Development**: Experiment with new features without affecting the production application
- **Risk Mitigation**: Test complex changes in a controlled environment before integration
- **Rapid Prototyping**: Faster iteration cycles during development
- **Clean Production Code**: Only proven, stable features make it to the production application

## Quick Start

### Production Application (onvif_manager)

Build and run the unified production application:
```cmd
cd onvif_manager
go build -o onvif-manager.exe cmd/app/main.go
```

Available modes:
```cmd
# Web application mode (embedded frontend + API on port 8090)
onvif-manager.exe web

# API server only mode (port 8090)
onvif-manager.exe server

# CLI mode with configuration management
onvif-manager.exe config apply cameras.csv config_1080p.csv

# Show help and available commands
onvif-manager.exe help
```

### Development Testbed (main_back / main_front)

For development and testing new features:

**Backend Testbed:**
```cmd
cd main_back
go run cmd/backend/main.go
```

**Frontend Testbed:**
```cmd
cd main_front
npm install
npm run dev
```

**Access Points:**
- **Production Web Interface**: http://localhost:8090 (onvif_manager web mode)
- **Production API**: http://localhost:8090/api (onvif_manager web/server mode)
- **Development Frontend**: http://localhost:5173 (main_front testbed)
- **Development Backend API**: http://localhost:8090 (main_back testbed)
- ONVIF Manager Web Interface: http://localhost:8090
- ONVIF Manager API: http://localhost:8090/api
- Legacy Frontend: http://localhost:5173 (if using legacy mode)
- Legacy Backend API: http://localhost:8090 (if using legacy mode)

## Development Workflow

### Testbed Development (main_back / main_front)

The testbed applications serve as the primary development environment where all new features and changes are prototyped.

**Backend Testbed (main_back):**
- **handlers.go**: REST API endpoints for all camera operations
- **manager.go**: Core camera management and ONVIF communication
- **client.go**: ONVIF client implementation with device discovery
- **resolution.go**: Camera resolution management and validation

**Frontend Testbed (main_front):**
Built with React 19 and Material-UI v7 for modern component styling.

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

### Production Integration (onvif_manager)

Once features are validated in the testbed, they are integrated into the production application.

**ONVIF Manager Structure:**
- **cmd/app/main.go**: Application entry point
- **internal/webserver/**: Web server with embedded frontend assets
- **internal/backend/**: Production backend implementation (ported from main_back)
- **internal/cli/**: Command-line interface using Cobra
- **internal/webserver/web/**: Embedded production frontend (built from main_front)
- **pkg/models/**: Shared data structures and models
- **examples/**: Sample CSV files for cameras and configurations

**Integration Process:**
1. **Test in Testbed**: Develop and validate features in main_back/main_front
2. **Port Backend Logic**: Copy validated backend code to onvif_manager/internal/backend/
3. **Build Frontend**: Build the React frontend and embed in onvif_manager/internal/webserver/web/
4. **Update CLI**: Add new CLI commands to onvif_manager/internal/cli/ if needed
5. **Test Production**: Verify the integrated functionality in onvif_manager
6. **Deploy**: The onvif_manager executable is ready for deployment

## API Endpoints

### Testbed API (main_back)
The testbed backend provides RESTful APIs for prototyping:

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

### Production API (onvif_manager)
The production application exposes similar endpoints under `/api` prefix:

- `GET /api/cameras` - List all cameras
- `POST /api/cameras` - Add new camera
- `DELETE /api/cameras/{id}` - Remove camera
- `GET /api/load-cam-list` - Load camera list
- `GET /api/check-single-cam/{id}` - Check camera status
- `POST /api/config-single-cam/{id}` - Configure camera
- `GET /api/validate-cam/{id}` - Validate camera stream
- `POST /api/cameras/import-csv` - Import cameras from CSV
- `POST /api/import-config-csv` - Import configuration CSV
- `POST /api/choose-cam-from-csv` - Select cameras from CSV
- `POST /api/apply-config` - Apply configuration to cameras
- `POST /api/export-validation-csv` - Export validation results
- `POST /api/vlc` - Launch VLC for camera stream

## CLI Interface (onvif_manager)

The production application includes a comprehensive CLI interface:

```cmd
# Show available commands
onvif-manager.exe help

# Web application modes
onvif-manager.exe web      # Start web interface + API
onvif-manager.exe server   # Start API server only

# Configuration management
onvif-manager.exe config apply [camera-csv] [config-csv]
onvif-manager.exe config show
onvif-manager.exe config set

# Export functionality
onvif-manager.exe export [output-file]
```

## Testing

### Testbed Testing (main_back / main_front)
Use the testbed environment for rapid development and testing:

**Frontend Testing:**
- **Browser DevTools**: Use console and network tabs for debugging
- **React DevTools**: Inspect component state and props
- **Debug Route**: Access `/debug` for development testing
- **Component Testing**: Individual component validation

**Backend Testing:**
- **API Testing**: Use curl, Postman, or similar tools for endpoint testing
- **Go Testing**: Run `go test ./...` in main_back directory
- **Integration Testing**: Full workflow testing with frontend and backend
- **ONVIF Testing**: Test with real cameras or ONVIF simulator

### Production Testing (onvif_manager)
Test the integrated production application:

**Web Interface Testing:**
- **Embedded Frontend**: Test the embedded web interface at http://localhost:8090
- **API Integration**: Verify API endpoints work with embedded frontend
- **CLI Testing**: Test command-line interface functionality

**Build Testing:**
```cmd
# Build production application
cd onvif_manager
go build -o onvif-manager.exe cmd/app/main.go

# Test web mode
onvif-manager.exe web

# Test CLI mode
onvif-manager.exe config apply examples/cameras.csv examples/config_1080p.csv
```

### Test Commands
```cmd
# Testbed backend tests
cd main_back
go test ./...

# Testbed frontend linting
cd main_front
npm run lint
npm run build

# Production application tests
cd onvif_manager
go test ./...
go build -o onvif-manager.exe cmd/app/main.go
```

## Debugging

### Testbed Debugging (main_back / main_front)
Use the testbed environment for detailed debugging:

**Frontend Debugging:**
- **Browser Console**: Check for JavaScript errors and debug output
- **Network Tab**: Inspect API requests and responses
- **React DevTools**: Examine component state and props
- **Source Maps**: Debug TypeScript/JSX code directly in browser
- **Console Logging**: Use `console.log` for data flow tracking

**Backend Debugging:**
- **Terminal Output**: Monitor detailed logging information
- **Log Files**: Capture logs with `go run cmd/backend/main.go 2>&1 | tee backend.log`
- **IDE Debugging**: Set breakpoints in VSCode or GoLand
- **API Debugging**: Test endpoints with curl or Postman
- **ONVIF Debugging**: Monitor ONVIF communication and errors

### Production Debugging (onvif_manager)
Debug the integrated production application:

**Application Modes:**
- **Web Mode**: Debug embedded frontend and API together
- **Server Mode**: Debug API-only functionality
- **CLI Mode**: Debug command-line operations with verbose output

**Debugging Commands:**
```cmd
# Debug web mode
onvif-manager.exe web

# Debug with verbose logging
go run cmd/app/main.go web

# Debug CLI operations
onvif-manager.exe config apply cameras.csv config.csv --verbose
```

### Common Issues
- **CORS Errors**: Ensure backend CORS is properly configured (check both testbed and production)
- **Camera Discovery**: Check network connectivity and ONVIF support
- **Stream Validation**: Verify FFmpeg/FFprobe installation and paths
- **Port Conflicts**: Ensure port 8090 is available for onvif_manager, port 5173 for main_front testbed
- **Integration Issues**: Test in testbed first, then verify integration in production
- **CLI Issues**: Use help commands to verify syntax and available options

## Camera Configuration Process

The system implements a sophisticated two-phase approach to camera configuration in both testbed and production:

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
- **CLI Support**: Command-line batch configuration via onvif_manager
- **CSV-based**: Configuration data managed through CSV files

This batch processing approach ensures efficient configuration of multiple cameras while maintaining system responsiveness in both development and production environments.

## Camera Status Monitoring

The Camera Status functionality is available in both testbed and production:

### Testbed Monitoring (main_front)
- **Real-time Status**: Check individual or all camera statuses
- **Pagination**: Browse cameras in groups of 50 for better performance
- **Stream Validation**: Validate RTSP streams using FFprobe
- **VLC Integration**: Launch VLC player for camera streams
- **Configuration Dialog**: Modify camera settings with real-time feedback

### Production Monitoring (onvif_manager web)
- **Embedded Interface**: Same monitoring capabilities in the unified application
- **API Integration**: Status monitoring through production API endpoints
- **CLI Monitoring**: Command-line status checking and validation

## CSV Data Management

The system uses CSV files for camera data persistence across all environments:

### Testbed CSV Management
- **Development Data**: Use main_back/internal/loader/cameras.csv for development
- **Rapid Testing**: Quick import/export for feature testing

### Production CSV Management
- **Example Files**: Use onvif_manager/examples/ for standard configurations:
  - `cameras.csv`: Sample camera definitions
  - `cameras_to_configure.csv`: Cameras ready for configuration
  - `config_1080p.csv`: 1080p configuration template
  - `config_720p.csv`: 720p configuration template
- **CLI Operations**: Import/export via command-line interface
- **Web Interface**: Upload and download CSV files through web UI

### CSV Features
- **Import/Export**: Import camera lists and export configurations
- **Auto-save**: Automatic saving of camera configurations
- **Backup**: Regular backups recommended for data preservation
- **Format**: Standardized CSV format for camera definitions
- **Templates**: Pre-configured templates for common resolutions

## Maintenance

### Development Workflow Maintenance
- **Testbed Sync**: Keep main_back and main_front in sync during active development
- **Integration Schedule**: Regularly integrate tested features from testbed to production
- **Version Control**: Use Git branches for feature development in testbed before merging

### Data Backup
- **Testbed Data**: Backup `main_back/internal/loader/cameras.csv` during development
- **Production Data**: Backup onvif_manager configuration and example files
- **System State**: Preserve camera settings and configurations across both environments
- **Version Control**: Use Git for code versioning and history

### Updates and Dependencies
- **Testbed Dependencies**: 
  - Frontend: Run `npm update` in main_front directory
  - Backend: Update Go modules with `go mod tidy` in main_back
- **Production Dependencies**: Update Go modules in onvif_manager with `go mod tidy`
- **ONVIF Library**: Check for updates to `github.com/videonext/onvif` in both environments
- **Security Updates**: Keep all dependencies up to date across testbed and production

### Performance Optimization
- **Pagination**: Camera lists are paginated at 50 items for optimal performance (both environments)
- **Batch Operations**: Use batch configuration for multiple cameras
- **Connection Pooling**: Backend manages ONVIF connections efficiently
- **Caching**: Status information is cached to reduce network load
- **Production Efficiency**: onvif_manager combines frontend and backend for reduced resource usage

### Troubleshooting
- **Development Issues**: Start debugging in testbed environment (main_back/main_front)
- **Production Issues**: Verify functionality works in testbed before investigating production
- **Connection Issues**: Verify network connectivity and firewall settings
- **Camera Discovery**: Ensure cameras are ONVIF compliant
- **Stream Problems**: Check FFmpeg/FFprobe installation and paths
- **Performance**: Monitor system resources during large camera operations
- **Integration Problems**: Compare testbed vs production behavior to isolate issues

### Deployment Strategy
1. **Develop in Testbed**: Use main_back/main_front for feature development
2. **Test Thoroughly**: Validate all functionality in testbed environment
3. **Build Production**: Integrate changes to onvif_manager
4. **Final Testing**: Test onvif_manager with integrated changes
5. **Deploy**: Use onvif_manager.exe for production deployment

---

_Last updated: July 2, 2025_

## Development Best Practices

### Testbed-First Development
1. **Always start new features in the testbed** (main_back/main_front)
2. **Test thoroughly** before integrating to production
3. **Use testbed for experimentation** and rapid prototyping
4. **Document changes** that need to be ported to production

### Integration Guidelines
1. **Code Review**: Review testbed changes before integration
2. **Incremental Integration**: Port changes in small, manageable chunks
3. **Testing**: Test each integration step in onvif_manager
4. **Rollback Plan**: Keep previous working versions for rollback if needed

### Git Workflow
- **Feature Branches**: Use separate branches for testbed development
- **Integration Branches**: Create integration branches when porting to production
- **Main Branch**: Keep main branch stable with working onvif_manager version
- **Commit Messages**: Clearly indicate whether changes are for testbed or production
