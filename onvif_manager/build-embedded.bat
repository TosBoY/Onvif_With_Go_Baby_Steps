@echo off
echo 🔨 Building ONVIF Manager with Embedded Frontend
echo.

echo 📦 Step 1: Building React frontend...
cd ..\main_front
call npm run build
if %errorlevel% neq 0 (
    echo ❌ Frontend build failed!
    exit /b 1
)

echo 📂 Step 2: Copying frontend files...
cd ..\onvif_manager
xcopy "..\main_front\dist" "cmd\app\web" /E /I /Y
if %errorlevel% neq 0 (
    echo ❌ Failed to copy frontend files!
    exit /b 1
)

echo 🔨 Step 3: Building Go binary...
go build -o onvif-manager-embedded.exe cmd/app/main.go cmd/app/webserver.go
if %errorlevel% neq 0 (
    echo ❌ Go build failed!
    exit /b 1
)

echo.
echo ✅ Build completed successfully!
echo.
echo 🚀 Usage:
echo   onvif-manager-embedded.exe web       - Start web application (frontend + API)
echo   onvif-manager-embedded.exe server    - Start API server only
echo   onvif-manager-embedded.exe list      - CLI commands
echo.
echo 🌐 Web application will be available at: http://localhost:8090
echo 🔌 API endpoints available at: http://localhost:8090/api
