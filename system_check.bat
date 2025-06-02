@echo off
echo ========================================
echo ONVIF Camera Management System - Status Check
echo ========================================
echo.

echo Checking project structure...

if exist "d:\VNG\test\main_onvif\main_back\cmd\backend\main.go" (
    echo ✓ Backend main.go found
) else (
    echo ✗ Backend main.go NOT found
)

if exist "d:\VNG\test\main_onvif\main_front\src\App.jsx" (
    echo ✓ Frontend App.jsx found
) else (
    echo ✗ Frontend App.jsx NOT found
)

if exist "d:\VNG\test\main_onvif\main_front\src\pages\Dashboard.jsx" (
    echo ✓ Dashboard component found
) else (
    echo ✗ Dashboard component NOT found
)

if exist "d:\VNG\test\main_onvif\main_back\config\cameras.json" (
    echo ✓ Camera configuration found
) else (
    echo ✗ Camera configuration NOT found
)

echo.
echo Checking dependencies...

cd /d "d:\VNG\test\main_onvif\main_front"
if exist "node_modules" (
    echo ✓ Frontend dependencies installed
) else (
    echo ✗ Frontend dependencies NOT installed
    echo   Run: npm install
)

cd /d "d:\VNG\test\main_onvif\main_back"
if exist "go.mod" (
    echo ✓ Backend go.mod found
) else (
    echo ✗ Backend go.mod NOT found
)

echo.
echo Testing backend connectivity...
echo Attempting to start backend server...

cd /d "d:\VNG\test\main_onvif\main_back"
timeout /t 1 >nul
echo Starting Go backend server...
start /B "Backend Test" cmd /c "go run cmd/backend/main.go & timeout /t 5"

timeout /t 3 >nul

echo Testing backend endpoint...
curl -s http://localhost:8090/cameras >nul 2>&1
if %ERRORLEVEL% == 0 (
    echo ✓ Backend server responding
) else (
    echo ? Backend server may not be responding
    echo   This is normal if cameras are not connected
)

echo.
echo ========================================
echo Status Check Complete
echo ========================================
echo.
echo To start the full system:
echo 1. Run: start_system.bat
echo 2. Open browser to: http://localhost:5173
echo.
pause
