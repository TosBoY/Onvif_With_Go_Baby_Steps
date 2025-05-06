package onvif_test

import (
	"encoding/xml"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
	
	"github.com/use-go/onvif/media"
	xsd "github.com/use-go/onvif/xsd/onvif"
)

// GetStreamURI retrieves the RTSP stream URI for a profile
func GetStreamURI(c *Camera, profileToken string) (string, error) {
	if c.Device == nil {
		return "", fmt.Errorf("camera not connected")
	}
	
	// Create the request
	getStreamUriRequest := media.GetStreamUri{
		StreamSetup: xsd.StreamSetup{
			Stream:    "RTP-Unicast",
			Transport: xsd.Transport{Protocol: "RTSP"},
		},
		ProfileToken: xsd.ReferenceToken(profileToken),
	}
	
	// Call the method
	getStreamUriResp, err := c.Device.CallMethod(getStreamUriRequest)
	if err != nil {
		return "", err
	}
	defer getStreamUriResp.Body.Close()
	
	// Read the response
	body, err := io.ReadAll(getStreamUriResp.Body)
	if err != nil {
		return "", err
	}
	
	// Extract the stream URI from XML
	return ExtractStreamURIFromXML(body)
}

// ExtractStreamURIFromXML parses GetStreamURI response XML to extract the URI
func ExtractStreamURIFromXML(xmlData []byte) (string, error) {
	var streamUriResp StreamUriResponse
	
	if err := xml.Unmarshal(xmlData, &streamUriResp); err != nil {
		return "", fmt.Errorf("error unmarshalling stream URI response: %v", err)
	}
	
	// Extract the URI
	streamUri := streamUriResp.Body.GetStreamUriResponse.MediaUri.Uri
	if streamUri == "" {
		return "", fmt.Errorf("no stream URI found in response")
	}
	
	return streamUri, nil
}

// RefreshStream attempts to refresh the RTSP stream
func RefreshStream(rtspURI string) error {
	// Method 1: Try to make a TCP connection to the RTSP server to wake it up
	rtspURL := strings.TrimPrefix(rtspURI, "rtsp://")
	parts := strings.Split(rtspURL, "/")
	if len(parts) == 0 {
		return fmt.Errorf("invalid RTSP URL format")
	}
	
	// Get the host:port part
	hostPort := parts[0]
	if !strings.Contains(hostPort, ":") {
		hostPort += ":554" // Default RTSP port
	}
	
	// Try to connect
	conn, err := net.DialTimeout("tcp", hostPort, 5*time.Second)
	if err != nil {
		return fmt.Errorf("failed to connect to RTSP server: %v", err)
	}
	defer conn.Close()
	
	return nil
}

// OpenStreamInVLC tries to open the stream in VLC media player
func OpenStreamInVLC(rtspURI string) error {
	// Close any existing VLC instances first
	if err := CloseVLCInstances(); err != nil {
		return fmt.Errorf("error closing VLC instances: %v", err)
	}
	
	// Command to open VLC with the stream
	var cmd *exec.Cmd
	
	switch runtime.GOOS {
	case "windows":
		// Try to find VLC in common install locations
		vlcPaths := []string{
			"C:\\Program Files\\VideoLAN\\VLC\\vlc.exe",
			"C:\\Program Files (x86)\\VideoLAN\\VLC\\vlc.exe",
		}
		
		var vlcPath string
		for _, path := range vlcPaths {
			if _, err := os.Stat(path); err == nil {
				vlcPath = path
				break
			}
		}
		
		if vlcPath == "" {
			return fmt.Errorf("VLC not found in common install locations")
		}
		
		cmd = exec.Command(vlcPath, rtspURI, "--no-video-title-show", "--rtsp-tcp")
	case "darwin": // macOS
		cmd = exec.Command("/Applications/VLC.app/Contents/MacOS/VLC", rtspURI, "--no-video-title-show", "--rtsp-tcp")
	default: // Linux and others
		cmd = exec.Command("vlc", rtspURI, "--no-video-title-show", "--rtsp-tcp")
	}
	
	// Start VLC in background
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start VLC: %v", err)
	}
	
	return nil
}

// CloseVLCInstances attempts to close all running VLC instances
func CloseVLCInstances() error {
	var cmd *exec.Cmd
	
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("taskkill", "/F", "/IM", "vlc.exe")
	case "darwin": // macOS
		cmd = exec.Command("killall", "VLC")
	default: // Linux and others
		cmd = exec.Command("killall", "vlc")
	}
	
	// We ignore errors here because it's okay if no VLC instances are running
	_ = cmd.Run()
	
	return nil
}

// ControlVLCHTTP attempts to control VLC via its HTTP interface
func ControlVLCHTTP(action string, vlcHTTPHost string, vlcHTTPPort int, password string) error {
	// VLC HTTP interface settings
	baseURL := fmt.Sprintf("http://%s:%d", vlcHTTPHost, vlcHTTPPort)
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	
	// Create the request
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/requests/status.xml?command=%s", baseURL, action), nil)
	if err != nil {
		return fmt.Errorf("error creating HTTP request: %v", err)
	}
	
	// Add authentication if configured
	if password != "" {
		req.SetBasicAuth("", password)
	}
	
	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending HTTP request to VLC: %v", err)
	}
	defer resp.Body.Close()
	
	// Check if successful
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("VLC HTTP interface returned status code %d", resp.StatusCode)
	}
	
	return nil
}