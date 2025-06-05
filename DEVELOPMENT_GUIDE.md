# Development Guide – ONVIF Camera Management System

## Project Structure

- **main_back/**: Go backend (API, camera logic, config)
- **main_front/**: React frontend (UI, config, API calls)

## Quick Start

```cmd
# From project root
start_system.bat
```

### Manual Start

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

- Frontend: http://localhost:5173
- Backend API: http://localhost:8090

## Development Workflow

### Frontend (main_front)
The frontend is built with React and uses Material-UI for the component library.

**Key Components:**
- **Dashboard.jsx**: Main control interface with camera listing and configuration panels
- **CameraCard.jsx**: Display for each camera with controls (VLC and Info buttons)
- **CameraConfigPanel.jsx**: Panel for configuring camera resolution and FPS settings
- **CameraInfoDialog.jsx**: Dialog for displaying camera information and delete option
- **ValidationResults.jsx**: Shows validation results after applying camera configurations

**Core Files:**
- `src/components/`: React components for UI elements
- `src/pages/`: Page-level components for routing
- `src/services/api.js`: API client for backend communication
- `src/App.jsx`: Main application with routing

### Backend (main_back)
The backend is built with Go and handles ONVIF camera communication and configuration.

**Key Components:**
- **handlers.go**: API endpoints for camera operations (list, add, delete, configure)
- **manager.go**: Core camera management functionality
- **client.go**: ONVIF client implementation
- **validator.go**: Stream validation using FFprobe

**Core Directories:**
- `cmd/backend/`: Application entry point
- `internal/api/`: API handlers and routes
- `internal/camera/`: Camera management logic
- `internal/ffprobe/`: Stream validation
- `pkg/models/`: Data models
- `config/cameras.json`: Camera configuration storage

## Testing

- **Frontend:** Use browser dev tools, React component tests, or `/debug` route.
- **Backend:** Use curl/Postman for API, Go tests for logic.

## Debugging

### Frontend
- Check browser console for errors and debug output
- Use the Network tab to inspect API requests and responses
- Examine React component state using React DevTools
- Use `console.log` for tracking data flow between components

### Backend
- Check terminal output for detailed logging information
- Use `go run cmd/backend/main.go 2>&1 | tee backend.log` to capture logs to a file
- Set breakpoints in code when using an IDE like VSCode or GoLand
- Check logs for camera configuration and validation details

## Camera Configuration Process

The system implements a two-phase approach to camera configuration:

1. **Phase 1 - Apply Configuration**: 
   - Applies settings to all selected cameras
   - Determines the closest supported resolution for each camera
   - Prepares stream URLs for validation

2. **Phase 2 - Validation**:
   - After a brief pause (1 second), validates all cameras in the original order
   - Uses FFprobe to check if the applied settings are functioning correctly
   - For fake cameras, simulates successful validation

This batch processing approach ensures efficient configuration of multiple cameras.

## Maintenance

- The camera configuration is stored in `main_back/config/cameras.json`
- Backup this file regularly to preserve camera settings
- Update the frontend and backend dependencies periodically
- Check for ONVIF library updates for new camera compatibility

---

_Last updated: June 5, 2025_
