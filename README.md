# ONVIF Camera Management System

## Overview

A full-stack system for configuring and monitoring ONVIF cameras.

- **Backend:** Go (Gin) API for ONVIF camera control and configuration.
- **Frontend:** React (Vite, Material-UI) for a modern, responsive UI.

## Features

- Camera discovery and listing
- Video configuration (resolution, FPS, bitrate, H264 profile)
- Stream management (RTSP, profile switching)
- Device info (manufacturer, model, firmware, serial)
- Bulk configuration for multiple cameras
- Basic network info display

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
  cmd/backend/main.go
  internal/api/
  internal/camera/
  config/cameras.json

main_front/
  src/
    components/
    pages/
    services/api.js
    App.jsx
```

## Development

- See `DEVELOPMENT_GUIDE.md` for detailed workflow, architecture, and troubleshooting.

## Notes

- Use `.gitignore` to avoid committing build artifacts, logs, and secrets.
- For production, move secrets to environment variables and enable authentication.

---

_Last updated: June 2025_