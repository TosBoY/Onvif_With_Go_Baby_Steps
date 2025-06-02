# ONVIF Camera Management Frontend

A modern React-based frontend for managing ONVIF cameras with a clean Material-UI interface.

## Features

- 📹 View and manage ONVIF cameras
- ⚙️ Configure camera settings (resolution, FPS)
- 🎨 Modern Material-UI design
- 📱 Responsive layout
- 🔄 Real-time configuration updates

## Prerequisites

- Node.js (v16 or higher)
- npm or yarn
- Backend server running on port 8090

## Installation

1. Install dependencies:
```bash
npm install
```

2. Start the development server:
```bash
npm run dev
```

3. Open your browser and navigate to `http://localhost:5173`

## Backend Integration

This frontend is designed to work with the Go backend server. Make sure the backend is running on `http://localhost:8090` before using the application.

To start the backend:
```bash
cd ../main_back
go run cmd/backend/main.go
```

## Configuration

The frontend expects the backend to be running on `http://localhost:8090`. You can modify this by creating a `.env.local` file:

```
VITE_API_BASE_URL=http://your-backend-url
```

## Troubleshooting

If you encounter issues connecting to the backend, please refer to the [Troubleshooting Guide](./TROUBLESHOOTING.md).
