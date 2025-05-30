package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	lib "onvif_back/lib"

	"github.com/videonext/onvif/profiles/media2"
	"github.com/videonext/onvif/soap"
)

const (
	cameraIP   = "192.168.1.30"
	username   = "admin"
	password   = "Admin123"
	mediaXAddr = "http://" + cameraIP + "/onvif/media_service"
)

func main() {
	fmt.Println("🎥 ONVIF Dual-Mode Resolution Changer")

	camera := lib.NewCamera(cameraIP, 80, username, password)
	err := camera.Connect()
	if err != nil {
		log.Printf("Ver10 connection failed: %v", err)
		useVer20Fallback()
		return
	}

	configs, err := lib.GetAllVideoEncoderConfigurations(camera)
	if err != nil || len(configs) == 0 {
		log.Println("Ver10 failed to get configs. Falling back to ver20.")
		useVer20Fallback()
		return
	}

	profiles, err := lib.GetAllProfiles(camera)
	if err != nil || len(profiles) == 0 {
		log.Println("Ver10 failed to get profiles. Falling back to ver20.")
		useVer20Fallback()
		return
	}

	config := configs[0]
	profile := profiles[0]

	options, err := lib.GetVideoEncoderOptions(camera, config.Token, profile.Token)
	if err != nil {
		log.Println("Ver10 failed to get options. Falling back to ver20.")
		useVer20Fallback()
		return
	}

	h264 := lib.ParseH264Options(options)
	if len(h264.ResolutionsAvailable) == 0 {
		log.Println("No resolutions found in ver10. Falling back to ver20.")
		useVer20Fallback()
		return
	}

	fmt.Println("Available Resolutions:")
	for i, res := range h264.ResolutionsAvailable {
		fmt.Printf("%d. %dx%d\n", i+1, res.Width, res.Height)
	}

	fmt.Print("Choose resolution: ")
	idx := readIntInput(1, len(h264.ResolutionsAvailable)) - 1
	selectedRes := h264.ResolutionsAvailable[idx]

	fmt.Printf("Frame Rate Range: %d - %d\n", h264.FrameRateRange.Min, h264.FrameRateRange.Max)
	fmt.Print("Enter frame rate: ")
	fps := readIntRangeInput(h264.FrameRateRange.Min, h264.FrameRateRange.Max)

	fmt.Print("Enter bitrate (kbps): ")
	bitrate := readIntInput(256, 20000)

	fmt.Println("\nAttempting to apply settings using ver10...")
	err = lib.SetVideoEncoderConfiguration(camera, config.Token, config.Name,
		selectedRes.Width, selectedRes.Height, fps, bitrate, config.GovLength, config.H264Profile)

	if err != nil {
		fmt.Printf("❌ ver10 SetVideoEncoderConfiguration returned error: %v\n", err)
	} else {
		updated, verifyErr := lib.GetVideoEncoderConfiguration(camera, config.Token)
		if verifyErr != nil {
			fmt.Printf("⚠️  Could not verify config after applying: %v\n", verifyErr)
		} else if updated.Width == selectedRes.Width &&
			updated.Height == selectedRes.Height &&
			updated.FrameRate == fps {
			fmt.Println("✅ Successfully applied using ver10.")
			return
		} else {
			fmt.Printf("❌ Ver10 config mismatch:\nExpected: %dx%d @%dfps\nGot:      %dx%d @%dfps\n",
				selectedRes.Width, selectedRes.Height, fps,
				updated.Width, updated.Height, updated.FrameRate)
		}
	}

	fmt.Println("🔁 Falling back to ver20...")
	applyVer20(selectedRes.Width, selectedRes.Height, fps, int32(config.GovLength))
}

func applyVer20(width, height int, fps int, gov int32) {
	client := soap.NewClient(soap.WithTimeout(5 * time.Second))
	client.AddHeader(soap.NewWSSSecurityHeader(username, password, time.Now()))
	mediaService := media2.NewMedia2(client, mediaXAddr)

	profiles, err := mediaService.GetProfiles(&media2.GetProfiles{})
	if err != nil || len(profiles.Profiles) == 0 {
		log.Fatalf("Ver20 failed to get profiles: %v", err)
	}

	configs, err := mediaService.GetVideoEncoderConfigurations(&media2.GetConfiguration{})
	if err != nil || len(configs.Configurations) == 0 {
		log.Fatalf("Ver20 failed to get configs: %v", err)
	}

	selectedConfig := configs.Configurations[0]
	selectedConfig.Resolution.Width = int32(width)
	selectedConfig.Resolution.Height = int32(height)
	selectedConfig.RateControl.FrameRateLimit = float32(fps)
	selectedConfig.GovLength = gov
	selectedConfig.Multicast.AutoStart = true

	_, err = mediaService.SetVideoEncoderConfiguration(&media2.SetVideoEncoderConfiguration{
		Configuration: selectedConfig,
	})
	if err != nil {
		log.Fatalf("Failed to apply using ver20: %v", err)
	}
	fmt.Println("✅ Successfully applied using ver20.")
}

func useVer20Fallback() {
	applyVer20(1920, 1080, 25, 30)
}

func readIntInput(min, max int) int {
	reader := bufio.NewReader(os.Stdin)
	for {
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		value, err := strconv.Atoi(input)
		if err != nil || value < min || value > max {
			fmt.Printf("Enter a number between %d and %d: ", min, max)
			continue
		}
		return value
	}
}

func readIntRangeInput(min, max int) int {
	fmt.Printf("Enter a value between %d and %d: ", min, max)
	return readIntInput(min, max)
}
