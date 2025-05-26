package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

// Default camera connection details
const (
	defaultCameraIP = "192.168.1.12"
	defaultTimeout  = 5 * time.Second
)

// Common ONVIF ports
var commonPorts = []int{80, 8000, 8080, 8081, 554, 10000, 5000}

// Common ONVIF paths
var commonPaths = []string{
	"/onvif/device_service",
	"/onvif/services",
	"/onvif",
	"/onvif/service",
	"/device_service",
	"/onvif/media_service",
}

func main() {
	// Define flags for camera connection parameters
	ipPtr := flag.String("ip", defaultCameraIP, "Camera IP address")
	portPtr := flag.Int("port", 0, "Camera port (0 to test common ports)")
	flag.Parse()

	fmt.Println("📹 ONVIF Camera Troubleshooter 📹")
	fmt.Printf("Testing connectivity to camera at %s\n", *ipPtr)

	// Step 1: Check if the IP is reachable via ping
	fmt.Println("\n🔍 Step 1: Testing if the camera IP is reachable...")
	if pingHost(*ipPtr) {
		fmt.Printf("✅ Ping to %s successful - the device is online and reachable\n", *ipPtr)
	} else {
		fmt.Printf("❌ Could not ping %s - the device may be offline or does not respond to ping\n", *ipPtr)
		fmt.Println("\n⚠️  Next steps to try:")
		fmt.Println("   1. Verify that the camera is powered on")
		fmt.Println("   2. Check if the IP address is correct")
		fmt.Println("   3. Ensure there are no firewall rules blocking ICMP traffic")
		fmt.Println("   4. Try connecting directly to the camera's network (if applicable)")
		fmt.Println("\n   Note: Some devices are configured to ignore ping requests for security reasons.")
		fmt.Println("         Let's continue with port scanning even if ping failed.")
	}

	// Step 2: Port scanning
	fmt.Println("\n🔍 Step 2: Scanning for open ports commonly used by ONVIF cameras...")
	
	portsToTest := commonPorts
	if *portPtr > 0 {
		portsToTest = []int{*portPtr}
	}
	
	openPorts := []int{}
	for _, port := range portsToTest {
		fmt.Printf("   Testing port %d... ", port)
		if isPortOpen(*ipPtr, port) {
			fmt.Println("✅ OPEN")
			openPorts = append(openPorts, port)
		} else {
			fmt.Println("❌ CLOSED")
		}
	}
	
	if len(openPorts) == 0 {
		fmt.Println("\n❌ No common ONVIF ports are open on this device.")
		fmt.Println("\n⚠️  Next steps to try:")
		fmt.Println("   1. Verify the camera's network settings")
		fmt.Println("   2. Check if the camera has ONVIF services enabled")
		fmt.Println("   3. Check if there are any firewall rules blocking access")
		os.Exit(1)
	} else {
		fmt.Printf("\n✅ Found %d open ports: %v\n", len(openPorts), openPorts)
	}

	// Step 3: Try to access ONVIF endpoints on open ports
	fmt.Println("\n🔍 Step 3: Testing ONVIF endpoints on open ports...")
	
	foundEndpoints := []string{}
	for _, port := range openPorts {
		for _, path := range commonPaths {
			url := fmt.Sprintf("http://%s:%d%s", *ipPtr, port, path)
			fmt.Printf("   Testing %s... ", url)
			
			status, err := getHTTPStatus(url)
			if err == nil && (status == 200 || status == 400 || status == 401 || status == 500) {
				// Status codes 200, 400, 401, or 500 might indicate a valid ONVIF endpoint
				fmt.Printf("✅ POSSIBLE ONVIF ENDPOINT (Status: %d)\n", status)
				foundEndpoints = append(foundEndpoints, url)
			} else if err != nil {
				fmt.Printf("❌ ERROR: %v\n", err)
			} else {
				fmt.Printf("❌ Status: %d\n", status)
			}
		}
	}
	
	fmt.Println("\n===== TROUBLESHOOTING SUMMARY =====")
	if len(foundEndpoints) > 0 {
		fmt.Printf("\n✅ Found %d potential ONVIF endpoints:\n", len(foundEndpoints))
		for i, endpoint := range foundEndpoints {
			fmt.Printf("   %d. %s\n", i+1, endpoint)
		}
		
		fmt.Println("\n🔧 Try modifying your get_h264_options.go file with these endpoints:")
		for _, endpoint := range foundEndpoints {
			parts := strings.Split(endpoint, "/")
			hostPort := parts[2]
			hostPortParts := strings.Split(hostPort, ":")
			host := hostPortParts[0]
			port := 80
			if len(hostPortParts) > 1 {
				fmt.Sscanf(hostPortParts[1], "%d", &port)
			}
			
			fmt.Printf("\n   go run cmd/get_h264_options.go -ip %s -port %d\n", host, port)
		}
	} else {
		fmt.Println("\n❌ Could not find any responsive ONVIF endpoints.")
		fmt.Println("\n⚠️  Possible issues:")
		fmt.Println("   1. The camera might not support ONVIF")
		fmt.Println("   2. ONVIF services might be disabled on the camera")
		fmt.Println("   3. The camera might require authentication for initial connections")
		fmt.Println("   4. There might be network restrictions preventing access")
		
		fmt.Println("\n📋 Recommendations:")
		fmt.Println("   1. Check the camera's documentation to confirm ONVIF support and settings")
		fmt.Println("   2. Verify the camera's IP address and port settings")
		fmt.Println("   3. Ensure ONVIF services are enabled in the camera's settings")
		fmt.Println("   4. Try accessing the camera's web interface to verify connectivity")
	}
}

// pingHost attempts to ping the specified host
func pingHost(host string) bool {
	var cmd *exec.Cmd
	
	if runtime.GOOS == "windows" {
		cmd = exec.Command("ping", "-n", "1", "-w", "1000", host)
	} else {
		cmd = exec.Command("ping", "-c", "1", "-W", "1", host)
	}
	
	err := cmd.Run()
	return err == nil
}

// isPortOpen checks if a TCP port is open
func isPortOpen(host string, port int) bool {
	address := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", address, defaultTimeout)
	
	if err != nil {
		return false
	}
	
	defer conn.Close()
	return true
}

// getHTTPStatus gets the HTTP status code for a URL
func getHTTPStatus(url string) (int, error) {
	client := &http.Client{
		Timeout: defaultTimeout,
	}
	
	resp, err := client.Get(url)
	if err != nil {
		return 0, err
	}
	
	defer resp.Body.Close()
	return resp.StatusCode, nil
}