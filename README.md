# ONVIF Camera Management System

## Overview

This application provides an intuitive user interface for configuring and monitoring ONVIF cameras. The application consists of:

- **Backend**: Go-based API server for ONVIF camera communication
- **Frontend**: Modern React/TypeScript UI with Material-UI components

## Getting Started

### Quick Start

Run the application with a single command:

```bash
start_app.bat
```

This will start both backend and frontend servers.

### Manual Start

To start the backend server:

```bash
cd main_back
go run cmd/backend/main.go
```

To start the frontend development server:

```bash
cd onvif_frontend
npm install
npm run dev
```

## Accessing the Application

Once both servers are running, open your browser and navigate to:

```
http://localhost:3000
```

## Current Capabilities

### Device Information
- Device metadata retrieval (manufacturer, model, firmware version, serial number)
- Basic system status information

### Video Configuration
- Resolution configuration
- H264 profile settings
- Frame rate adjustment
- Bitrate control
- GoV (Group of Pictures) length settings

### Stream Management
- RTSP stream handling
- Profile selection and switching

### Basic Device Control
- Manual camera connection (IP, username, password)
- Basic device information retrieval
- Video configuration validation

### Network Features
- IPv4 support
- RTSP/RTP streaming
- Basic network configuration display

This system provides fundamental ONVIF camera management capabilities with a focus on video streaming and configuration.