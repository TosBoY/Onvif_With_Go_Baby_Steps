package main

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/use-go/onvif"
)

func main() {
	ip := "192.168.1.12"
	port := 80
	username := "admin"
	password := "admin123"

	fmt.Printf("Testing connection to camera at %s:%d...\n", ip, port)

	// 1. Test TCP connection first
	fmt.Printf("1. Testing TCP connection to %s:%d...\n", ip, port)
	timeout := time.Second * 5
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", ip, port), timeout)
	if err != nil {
		fmt.Printf("❌ TCP connection failed: %v\n", err)
	} else {
		fmt.Println("✅ TCP connection successful")
		conn.Close()
	}

	// 2. Try HTTP GET request
	fmt.Printf("\n2. Testing HTTP connectivity...\n")
	client := &http.Client{Timeout: timeout}
	resp, err := client.Get(fmt.Sprintf("http://%s:%d", ip, port))
	if err != nil {
		fmt.Printf("❌ HTTP request failed: %v\n", err)
	} else {
		fmt.Printf("✅ HTTP request successful (Status: %s)\n", resp.Status)
		resp.Body.Close()
	}

	// 3. Try different ONVIF paths
	fmt.Printf("\n3. Testing ONVIF endpoints...\n")
	paths := []string{
		"/onvif/device_service",
		"/onvif/service",
		"/device_service",
		"",
	}

	for _, path := range paths {
		fullURL := fmt.Sprintf("http://%s:%d%s", ip, port, path)
		fmt.Printf("\nTrying ONVIF connection to: %s\n", fullURL)
		dev, err := onvif.NewDevice(onvif.DeviceParams{
			Xaddr:      fullURL,
			Username:   username,
			Password:   password,
			HttpClient: client,
		})
		if err != nil {
			fmt.Printf("❌ Connection failed: %v\n", err)
		} else {
			fmt.Printf("✅ Connection successful with path: %s\n", path)
			fmt.Println("Trying to get device information...")
			info, err := dev.GetInformation()
			if err != nil {
				fmt.Printf("❌ Failed to get device info: %v\n", err)
			} else {
				fmt.Printf("✅ Device Info:\n  Manufacturer: %s\n  Model: %s\n  Firmware: %s\n",
					info.Manufacturer, info.Model, info.FirmwareVersion)
			}
		}
	}
}
