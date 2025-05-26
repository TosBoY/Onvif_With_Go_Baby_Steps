package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/use-go/onvif"
)

// Default camera connection details
const (
	defaultCameraIP   = "192.168.1.12"
	defaultUsername   = "admin"
	defaultPassword   = "admin123"
	defaultTimeout    = 5 * time.Second
	defaultDevicePath = "/onvif/device_service"
)

// Common ONVIF paths to try
var onvifPaths = []string{
	"/onvif/device_service",
	"/onvif/services",
	"/onvif",
	"/onvif/service",
	"/device_service",
	"/onvif/media_service",
}

// Common ports for ONVIF devices
var commonPorts = []int{80, 8000, 8080, 8081, 554, 10000, 5000}

func main() {
	// Define flags for camera connection parameters
	ipPtr := flag.String("ip", defaultCameraIP, "Camera IP address")
	usernamePtr := flag.String("user", defaultUsername, "Username")
	passwordPtr := flag.String("pass", defaultPassword, "Password")
	flag.Parse()

	fmt.Println("📹 ONVIF Camera Connection Troubleshooter 📹")

	// Step 1: Test basic network connectivity to the camera IP
	fmt.Printf("\n🔍 Step 1: Testing network connectivity to %s...\n", *ipPtr)

	pingResult := pingHost(*ipPtr)
	if pingResult {
		fmt.Printf("✅ Ping to %s successful - the device is online\n", *ipPtr)
	} else {
		fmt.Printf("⚠️ Could not ping %s - the device may be offline or not responding to pings\n", *ipPtr)
		fmt.Println("   Many devices are configured to ignore ping requests for security reasons.")
		fmt.Println("   Continuing with port scanning...")
	}

	// Step 2: Scan for open ports
	fmt.Printf("\n🔍 Step 2: Scanning for common open ports on %s...\n", *ipPtr)
	openPorts := []int{}

	for _, port := range commonPorts {
		fmt.Printf("   Testing port %d... ", port)
		if isPortOpen(*ipPtr, port) {
			fmt.Println("✅ OPEN")
			openPorts = append(openPorts, port)
		} else {
			fmt.Println("❌ CLOSED")
		}
	}

	if len(openPorts) == 0 {
		fmt.Println("\n❌ No common ports are open on this device. The camera may be:")
		fmt.Println("   1. Offline or unreachable")
		fmt.Println("   2. Using non-standard ports")
		fmt.Println("   3. Blocked by a firewall")
		promptToContinue()
	} else {
		fmt.Printf("\n✅ Found %d open ports: %v\n", len(openPorts), openPorts)
	}

	// Step 3: Test ONVIF connectivity using different paths and ports
	fmt.Println("\n🔍 Step 3: Testing ONVIF connectivity with different paths and ports...")

	foundWorkingCombination := false
	var workingPort int
	var workingPath string

	for _, port := range openPorts {
		for _, path := range onvifPaths {
			fmt.Printf("\nTrying connection with port %d and path '%s'\n", port, path)
			result := testOnvifConnection(*ipPtr, port, path, *usernamePtr, *passwordPtr)
			if result {
				fmt.Printf("✅ SUCCESS! Connection established using port %d and path '%s'\n", port, path)
				foundWorkingCombination = true
				workingPort = port
				workingPath = path
				break
			}
		}
		if foundWorkingCombination {
			break
		}
	}

	// Step 4: Compare with the working interactive resolution change settings
	fmt.Println("\n🔍 Step 4: Comparing with settings from interactive_resolution_change.go...")

	fmt.Println("Settings in interactive resolution change:")
	fmt.Println("   - URL format: http://{IP}:{PORT}")
	fmt.Println("   - No explicit path provided in the NewCamera call")
	fmt.Println("   - Using Camera wrapper from onvif_test2/lib")

	fmt.Println("\nSettings in get_h264_options.go:")
	fmt.Println("   - URL format: http://{IP}:{PORT}/onvif/device_service")
	fmt.Println("   - Explicitly using '/onvif/device_service' path")
	fmt.Println("   - Using onvif.NewDevice directly instead of the Camera wrapper")

	// Step 5: Provide recommendations
	fmt.Println("\n===== CONNECTION TROUBLESHOOTING SUMMARY =====")

	if foundWorkingCombination {
		fmt.Println("\n✅ Good news! A working connection was found with the following parameters:")
		fmt.Printf("   - IP: %s\n", *ipPtr)
		fmt.Printf("   - Port: %d\n", workingPort)
		fmt.Printf("   - Path: %s\n", workingPath)
		fmt.Printf("   - Username: %s\n", *usernamePtr)
		fmt.Printf("   - Password: %s\n", strings.Repeat("*", len(*passwordPtr)))

		fmt.Println("\n🔧 Fix for get_h264_options.go:")
		fmt.Println("   You need to update the xaddr in get_h264_options.go to use the correct URL format:")
		fmt.Printf("   xaddr := fmt.Sprintf(\"http://%%s:%%d%s\", *ipPtr, %d)\n", workingPath, workingPort)

		fmt.Println("\n📋 Try running the following command:")
		fmt.Printf("   go run cmd\\get_h264_options.go -ip %s -port %d\n", *ipPtr, workingPort)
	} else {
		fmt.Println("\n❌ Could not establish a working ONVIF connection.")
		fmt.Println("\nRecommendations:")
		fmt.Println("1. Verify the camera supports ONVIF and that ONVIF services are enabled")
		fmt.Println("2. Check firewall settings that might be blocking connections")
		fmt.Println("3. Try with different credentials if authentication might be the issue")
		fmt.Println("4. Look at the working interactive_resolution_change.go code and see how it connects")
		fmt.Println("5. Examine the onvif_test2/lib/client.go file to see how the connection is established there")
	}

	// Step 6: Try to create a connection using the same method as interactive_resolution_change.go
	fmt.Println("\n🔍 Step 6: Attempting to connect using method from interactive_resolution_change.go...")

	// This requires the onvif_test2/lib package, so providing instructions instead
	fmt.Println("\nTo test using the same connection method as interactive_resolution_change.go:")
	fmt.Println("1. Create a new file in onvif_test/cmd/test_connection.go")
	fmt.Println("2. Import the onvif_test2/lib package (you'll need to update your go.mod)")
	fmt.Println("3. Use the lib.NewCamera() and camera.Connect() methods")

	fmt.Println("\n📝 Sample code:")
	fmt.Println(`
package main

import (
	"fmt"
	"log"
	
	lib "onvif_test2/lib"  // You'll need to update module paths
)

func main() {
	camera := lib.NewCamera("192.168.1.12", 80, "admin", "admin123")
	err := camera.Connect()
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	fmt.Println("Connected successfully!")
}
`)
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

// testOnvifConnection attempts to connect to the camera using the specified parameters
func testOnvifConnection(ip string, port int, path string, username, password string) bool {
	// Build the camera endpoint URL
	xaddr := fmt.Sprintf("http://%s:%d%s", ip, port, path)

	fmt.Printf("Attempting to connect to: %s\n", xaddr)
	fmt.Printf("Using credentials: %s / %s\n", username, strings.Repeat("*", len(password)))

	// Create a new ONVIF device connection with timeout
	_, err := onvif.NewDevice(onvif.DeviceParams{
		Xaddr:    xaddr,
		Username: username,
		Password: password,
	})

	if err != nil {
		fmt.Printf("❌ Connection failed: %v\n", err)
		return false
	}

	fmt.Println("✅ Device connection established successfully")
	return true
}

// promptToContinue asks the user if they want to continue
func promptToContinue() {
	fmt.Print("\nPress Enter to continue troubleshooting anyway...")
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')
}
