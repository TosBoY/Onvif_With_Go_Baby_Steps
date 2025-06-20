# RTSP Stream Analyzer

This tool analyzes RTSP streams and provides information about:
- Video resolution
- Frame rate (FPS)
- Video and audio codecs
- Other stream details

## Prerequisites

1. FFmpeg development libraries must be installed
2. GCC compiler or equivalent

## Installation

Make sure you have FFmpeg installed and the development headers available.
The program expects to find FFmpeg libraries and headers in standard locations.

For Windows:
- FFmpeg binaries should be in C:\ffmpeg\bin
- Include files should be in C:\ffmpeg\include
- Library files should be in C:\ffmpeg\lib

## Building

### Using the batch file (Windows):

```
build.bat
```

### Using make:

```
make
```

### Manual compilation:

```
gcc -Wall -o rtsp_analyzer rtsp_analyzer.c -lavformat -lavcodec -lavutil -lswscale -lavdevice -lswresample
```

## Usage

```
rtsp_analyzer.exe <rtsp_url>
```

Example:
```
rtsp_analyzer.exe rtsp://username:password@camera-ip:port/path
rtsp_analyzer.exe rtsp://admin:admin123@192.168.1.100:554/live
```

## Output

The program will output details about the RTSP stream including:

- Overall stream information
- Video resolution
- Frame rate
- Video codec
- Audio codec (if available)
- Bit rates
- Pixel format

## Troubleshooting

1. If you get "cannot find -l(library)" errors, make sure your FFmpeg development libraries are properly installed and linked.

2. If you get "header file not found" errors, check that FFmpeg development headers are in your compiler's include path.

3. For runtime errors about missing DLLs, ensure that FFmpeg DLLs are in your PATH or in the same directory as the executable.
