# ONVIF Camera Management System

A complete solution for managing ONVIF-compatible IP cameras with both backend and frontend components.

## Project Structure

- **onvif_back**: Go-based backend server for communicating with ONVIF cameras
- **onvif_front**: React-based frontend interface
- **onvif_docs**: Documentation resources for ONVIF protocols
- **onvif_test**: Test utilities and examples for ONVIF functionality

## Features

- Camera discovery and connection
- Video stream management and viewing through VLC
- Resolution and configuration management (change resolution, bitrate, etc.)
- Device information display
- Real-time stream viewing via RTSP
- Integration with VLC media player for stream playback
- Comprehensive camera configuration options

## Running Locally (Windows)

### Prerequisites

- Go 1.21 or later installed
- Node.js and npm installed
- VLC Media Player installed (for stream playback)

### Start the Backend

```bash
# Navigate to the backend directory
cd onvif_back

# Start the Go server
go run main.go
```

The backend server will start on `http://localhost:8090`.

### Start the Frontend

```bash
# In a new terminal, navigate to the frontend directory
cd onvif_front

# Install dependencies (if not already installed)
npm install

# Start the development server
npm run dev
```

The frontend will be available at `http://localhost:5173`.

### Access the Application

Open your web browser and navigate to:
```
http://localhost:5173
```

## Camera Configuration

The application is pre-configured to connect to a camera with the following default settings:

```
IP Address: 192.168.1.12
Username: admin
Password: admin123
```

To change these settings, modify the `cameraConfig` variable in `onvif_back/main.go`.

## Using the Application

1. **Connect to Camera**: The application will automatically attempt to connect to the configured camera on startup.

2. **View Profiles**: Select from available camera profiles to access different streams.

3. **Modify Settings**: Change resolution, frame rate, bit rate, and other encoding settings.

4. **Launch Stream**: Click "Play Stream in VLC" to open the current stream in VLC media player.

## Running on Raspberry Pi

The repository includes scripts for running on a Raspberry Pi:

```bash
# To start both backend and frontend services
./run_onvif.sh

# To stop the services
./stop_onvif.sh
```

Refer to the Pi deployment documentation for details on setting up for headless environments.

## Troubleshooting

- **VLC Not Starting**: Ensure VLC is correctly installed and accessible in your system PATH
- **Connection Issues**: Verify camera IP, username, password, and network connectivity
- **Stream Not Playing**: Some cameras require specific RTSP transport settings (TCP vs UDP)

## Development Notes

- Backend uses the [use-go/onvif](https://github.com/use-go/onvif) library for ONVIF protocol communication
- Frontend is built with React and Material UI
- API communication between frontend and backend is handled via Axios

## License

This project is available for internal use.