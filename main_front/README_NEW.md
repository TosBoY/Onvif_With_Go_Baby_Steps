# ONVIF Camera Management Frontend ✅

A modern React-based frontend for managing ONVIF cameras with configuration capabilities.

## 🎉 Status: READY TO USE

### ✅ Recent Fixes Completed
- **White Screen Issue**: Fixed syntax error in Dashboard component
- **Component Rendering**: All components now render correctly
- **Error Handling**: Comprehensive error boundaries and validation
- **API Integration**: Full backend communication with CORS support

## 🚀 Quick Start

### Option 1: Automated Startup (Recommended)
```cmd
# From project root directory: d:\VNG\test\main_onvif\
start_system.bat
```

### Option 2: Manual Startup

#### Start Backend Server
```cmd
cd d:\VNG\test\main_onvif\main_back
go run cmd/backend/main.go
```
**Backend URL**: http://localhost:8090

#### Start Frontend Development Server
```cmd
cd d:\VNG\test\main_onvif\main_front
npm run dev
```
**Frontend URL**: http://localhost:5173

## 🎯 Application Features

### Main Dashboard (`/`)
- **Camera List Display**: Shows all cameras from backend API
- **Camera Selection**: Click any camera card to select it
- **Configuration Panel**: Adjust resolution and FPS for selected cameras
- **Real-time Status**: Backend connection monitoring
- **Error Recovery**: Automatic retry mechanisms

### User Interface
- Modern Material-UI design with animations
- Responsive grid layout for camera cards
- Color-coded status indicators (Connected/Simulation)
- Loading states with progress indicators
- Success/error alerts with user-friendly messages

### Test Routes (for debugging)
- `/debug` - Step-by-step component testing
- `/test` - Simple test dashboard
- `/simple` - Minimal functionality test

## 🔧 Technical Implementation

### Core Components
- `Dashboard.jsx` - Main camera management interface ✅
- `CameraCard.jsx` - Individual camera display with selection ✅
- `CameraConfigPanel.jsx` - Camera configuration interface ✅
- `ConnectionStatus.jsx` - Backend connectivity monitoring ✅
- `Header.jsx` - Application header with branding ✅

### API Integration
- `api.js` - Backend communication service ✅
- Automatic error handling and retry logic
- Development proxy configuration
- CORS support for cross-origin requests

### Configuration Files
- `vite.config.js` - Development server with backend proxy ✅
- `package.json` - Dependencies and scripts ✅
- `.env` - Environment variables for API URL

## 📋 Available Scripts

```cmd
npm run dev          # Start development server
npm run build        # Build for production
npm run preview      # Preview production build
npm run start:backend    # Start Go backend server
npm run start:full      # Start both backend and frontend
```

## 🔌 API Endpoints

### Backend Communication
- `GET /cameras` - Fetch camera list with status
- `POST /apply-config` - Apply camera configuration

### Configuration Format
```json
{
  "cameraId": "1",
  "width": 1920,
  "height": 1080,
  "fps": 30
}
```

## 🛠️ Development Setup

### Prerequisites
- Node.js 16+ 
- Go 1.19+
- Git

### Installation
```cmd
cd d:\VNG\test\main_onvif\main_front
npm install
```

### Development
```cmd
npm run dev
```

## 🐛 Troubleshooting

### ✅ Fixed Issues

#### White Screen Problem (RESOLVED)
- **Issue**: Dashboard showing blank page
- **Cause**: Missing newline between conditional renders in Dashboard.jsx
- **Status**: FIXED ✅

#### Component Import Errors (RESOLVED)
- **Issue**: Components failing to load
- **Cause**: Syntax errors in conditional logic
- **Status**: FIXED ✅

### Current Known Issues
- None reported ✅

### Debug Steps
1. Check browser console for JavaScript errors
2. Verify backend is running on port 8090
3. Test connection at http://localhost:8090/cameras
4. Use `/debug` route for component testing

## 📁 Project Structure

```
main_front/
├── src/
│   ├── components/      # UI components ✅
│   │   ├── CameraCard.jsx
│   │   ├── CameraConfigPanel.jsx
│   │   ├── ConnectionStatus.jsx
│   │   ├── Header.jsx
│   │   └── Loading.jsx
│   ├── pages/          # Page components ✅
│   │   ├── Dashboard.jsx
│   │   └── SimpleDashboard.jsx
│   ├── services/       # API services ✅
│   │   └── api.js
│   └── assets/        # Static files
├── public/            # Public assets
├── vite.config.js    # Vite configuration ✅
├── package.json      # Dependencies ✅
└── README.md         # This file
```

## 🎨 UI Features

### Camera Management
- Grid layout with responsive design
- Hover animations and visual feedback
- Status indicators (Connected/Simulation)
- Click-to-select functionality

### Configuration Panel
- Dropdown resolution presets
- FPS selection options
- Real-time configuration preview
- Apply button with loading states

### Error Handling
- Connection status monitoring
- Retry mechanisms for failed requests
- User-friendly error messages
- Graceful fallbacks

## 🔗 Related Documentation

- [Backend Documentation](../main_back/README.md)
- [Troubleshooting Guide](./TROUBLESHOOTING.md)
- [Material-UI Docs](https://mui.com/)
- [Vite Documentation](https://vitejs.dev/)

---

**🎯 Status**: Production Ready ✅  
**🔄 Last Updated**: December 2024  
**👨‍💻 Tested**: Windows Environment with cmd.exe
