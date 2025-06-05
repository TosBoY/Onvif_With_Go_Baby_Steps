# ONVIF Camera Management System

## Overview

A full-stack system for configuring and monitoring ONVIF IP cameras using the ONVIF protocol.

- **Backend:** Go API for ONVIF camera control and configuration with stream validation
- **Frontend:** React (Vite, Material-UI) for a modern, responsive UI

## Features

- Camera management (list, add, delete)
- Video configuration (resolution, FPS)
- VLC stream launching with direct RTSP URLs
- Camera information display with status
- Batch configuration for multiple cameras
- Two-phase configuration with validation
- Support for real and simulated (fake) cameras
- Stream validation with FFprobe

## Getting Started

### Quick Start

```cmd
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

## Project Structure

```
main_back/
  cmd/backend/main.go        # Application entry point
  internal/api/handlers.go   # API endpoint handlers
  internal/camera/manager.go # Camera management logic
  internal/ffprobe/          # Stream validation tool
  config/cameras.json        # Camera configuration storage

main_front/
  src/
    components/              # UI components
      CameraCard.jsx         # Individual camera display
      CameraConfigPanel.jsx  # Configuration interface
      CameraInfoDialog.jsx   # Camera info popup with delete option
      ValidationResults.jsx  # Stream validation display
    pages/
      Dashboard.jsx          # Main application page
    services/api.js          # Backend API client
    App.jsx                  # Main application with routing
```

## Key Features

1. **Multiple Camera Support**: Configure and manage any number of ONVIF-compatible cameras
2. **Batch Configuration**: Apply settings to multiple cameras in a single operation
3. **Two-Phase Process**: First apply settings to all cameras, then validate all cameras
4. **Stream Validation**: Verify actual stream parameters match requested configuration
5. **Fake Camera Support**: Test system functionality without physical cameras

## Development

See `DEVELOPMENT_GUIDE.md` for detailed workflow, architecture, and troubleshooting information.

---

_Last updated: June 5, 2025_