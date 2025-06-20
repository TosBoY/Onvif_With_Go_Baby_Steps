@echo off
REM Build Go RTSP Analyzer with CGO

echo Building Go RTSP Analyzer with CGO...

REM Check if gcc is available
where gcc >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo ERROR: GCC not found in PATH
    echo Make sure you have MSYS2/MinGW or another GCC distribution installed and added to PATH
    exit /b 1
)

REM Set CGO environment variables
set CGO_ENABLED=1
set CC=gcc

REM Note: CGO flags are now defined in the source file rtsp_analyzer.go
REM with platform-specific directives. These environment variables
REM are kept for compatibility but may be overridden by source directives.
set FFMPEG_INCLUDE=C:\ffmpeg\include
set FFMPEG_LIB=C:\ffmpeg\lib

REM Set CGO flags (these may be overridden by source file directives)
set CGO_CFLAGS=-I%FFMPEG_INCLUDE%
set CGO_LDFLAGS=-L%FFMPEG_LIB% -lavformat -lavcodec -lavutil

echo Building rtsp_analyzer.go...
go build -o rtsp_analyzer_go.exe rtsp_analyzer.go

if %ERRORLEVEL% NEQ 0 (
    echo Build failed!
    exit /b 1
)

@REM echo Building test_analyzer.go...
@REM go build -o test_analyzer.exe test_analyzer.go

if %ERRORLEVEL% NEQ 0 (
    echo Test build failed!
    exit /b 1
)

echo.
echo Build successful!
echo.
echo Usage: test_analyzer.exe "rtsp://username:password@camera_ip:port/path"
echo Example: test_analyzer.exe "rtsp://admin:admin123@192.168.1.100:554/live"
echo.
echo The Go package can also be imported and used in other Go programs:
echo   import "main_back/internal/ffmpeg"
echo   info, err := ffmpeg.AnalyzeRTSPStream("rtsp://...")
