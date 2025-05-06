# ONVIF Camera Management System

A complete solution for managing ONVIF-compatible IP cameras with both backend and frontend components.

## Project Structure

- **onvif_back**: Go-based backend server for communicating with ONVIF cameras
- **onvif_front**: React-based frontend interface
- **onvif_docs**: Documentation resources for ONVIF protocols
- **onvif_test**: Test utilities and examples for ONVIF functionality

## Features

- Camera discovery and management
- Video stream viewing
- Resolution and configuration management
- Device information display

## Getting Started

To run the application:

```bash
./run_onvif.sh
```

To stop the application:

```bash
./stop_onvif.sh
```

## Requirements

- Go (for backend)
- Node.js and npm (for frontend)
- Network access to ONVIF-compatible cameras