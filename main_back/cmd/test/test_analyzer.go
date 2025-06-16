package main

import (
	"fmt"
	"log"
	"os"

	"main_back/internal/ffmpeg"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <rtsp_url>\n", os.Args[0])
		os.Exit(1)
	}

	rtspURL := os.Args[1]

	fmt.Printf("Analyzing RTSP stream: %s\n", rtspURL)
	fmt.Println("===== RTSP Stream Analysis =====")

	// Analyze the RTSP stream
	streamInfo, err := ffmpeg.AnalyzeRTSPStream(rtspURL)
	if err != nil {
		log.Fatalf("Error analyzing stream: %v", err)
	}
	// Display the results
	fmt.Printf("Codec: %s\n", streamInfo.Codec)
	fmt.Printf("Resolution: %s\n", streamInfo.GetStreamResolution())
	fmt.Printf("Frame Rate: %s\n", streamInfo.GetFrameRate())
}
