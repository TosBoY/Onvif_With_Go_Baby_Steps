@echo off
REM Build and run RTSP analyzer

echo Building RTSP Analyzer...
IF NOT EXIST "C:\msys64\mingw64\bin\gcc.exe" (
    echo ERROR: GCC not found in C:\msys64\mingw64\bin
    echo Make sure you have MSYS2/MinGW or another GCC distribution installed
    exit /b 1
)

set INCLUDE_PATH=C:\ffmpeg\include
set LIB_PATH=C:\ffmpeg\lib

gcc -Wall -o rtsp_analyzer.exe rtsp_analyzer.c -I%INCLUDE_PATH% -L%LIB_PATH% -lavformat -lavcodec -lavutil

if %ERRORLEVEL% NEQ 0 (
    echo Build failed!
    exit /b 1
)

echo.
echo Build successful!
echo.
echo Usage: rtsp_analyzer.exe "rtsp://username:password@camera_ip:port/path"
echo Example: rtsp_analyzer.exe "rtsp://admin:admin123@192.168.1.100:554/live"
echo Note: URLs with special characters like '&' or '?' must be enclosed in double quotes
