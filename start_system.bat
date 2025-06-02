@echo off
echo Starting ONVIF Camera Management System...
echo.

echo Starting backend server...
cd /d "d:\VNG\test\main_onvif\main_back"
start "Backend Server" cmd /k "go run cmd/backend/main.go"

echo Waiting for backend to start...
timeout /t 3 > nul

echo Starting frontend development server...
cd /d "d:\VNG\test\main_onvif\main_front"
start "Frontend Dev Server" cmd /k "npm run dev"

echo.
echo Both servers should be starting now.
echo Backend: http://localhost:8090
echo Frontend: http://localhost:5173
echo.
pause
