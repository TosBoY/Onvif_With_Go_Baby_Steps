package onvif_back

import (
	"encoding/xml"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/use-go/onvif/device"
	xonvif "github.com/use-go/onvif/xsd/onvif"
)

// GetAllCapabilities retrieves all capabilities from an ONVIF camera
func GetAllCapabilities(camera *Camera) (*xonvif.Capabilities, string, error) {
	// Create request with empty category to get all capabilities
	req := device.GetCapabilities{
		Category: "", // Empty category to get all capabilities
	}

	// Call the method
	resp, err := camera.Device.CallMethod(req)
	if err != nil {
		return nil, "", fmt.Errorf("error calling GetCapabilities: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("error reading response body: %v", err)
	}

	// Save raw XML response for debugging
	rawResponse := string(body)

	// Parse the response
	var capabilitiesResp GetCapabilitiesResponse
	if err := xml.Unmarshal(body, &capabilitiesResp); err != nil {
		return nil, rawResponse, fmt.Errorf("error parsing response: %v", err)
	}

	return &capabilitiesResp.Capabilities, rawResponse, nil
}

// GetAllCapabilitiesWithCategory retrieves capabilities from an ONVIF camera with a specific category
// Some cameras don't handle empty category correctly, so we need this variant
func GetAllCapabilitiesWithCategory(camera *Camera, category string) (*xonvif.Capabilities, string, error) {
	// Create request with the specified category
	req := device.GetCapabilities{
		Category: xonvif.CapabilityCategory(category), // Convert string to proper type
	}

	// Call the method
	resp, err := camera.Device.CallMethod(req)
	if err != nil {
		return nil, "", fmt.Errorf("error calling GetCapabilities with category %s: %v", category, err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("error reading response body: %v", err)
	}

	// Save raw XML response for debugging
	rawResponse := string(body)

	// Check if the response contains a SOAP fault
	if isSoapFault(rawResponse) {
		return nil, rawResponse, fmt.Errorf("received SOAP fault from camera")
	}

	// Parse the response
	var capabilitiesResp GetCapabilitiesResponse
	if err := xml.Unmarshal(body, &capabilitiesResp); err != nil {
		return nil, rawResponse, fmt.Errorf("error parsing response: %v", err)
	}

	return &capabilitiesResp.Capabilities, rawResponse, nil
}

// Helper function to check if a response contains a SOAP fault
func isSoapFault(xmlResponse string) bool {
	return strings.Contains(xmlResponse, "<s:Fault>") ||
		strings.Contains(xmlResponse, "<Fault>")
}

// DisplayCapabilities prints all camera capabilities in a structured format
func DisplayCapabilities(capabilities *xonvif.Capabilities) {
	fmt.Println("\n╔═════════════════════════════════════════════╗")
	fmt.Println("║        ONVIF CAMERA CAPABILITIES             ║")
	fmt.Println("╚═════════════════════════════════════════════╝")

	// Analytics capabilities
	if capabilities.Analytics.XAddr != "" {
		fmt.Println("\n[ANALYTICS CAPABILITIES]")
		fmt.Printf("XAddr: %s\n", capabilities.Analytics.XAddr)
		fmt.Printf("Rule Support: %t\n", capabilities.Analytics.RuleSupport)
		fmt.Printf("Analytics Module Support: %t\n", capabilities.Analytics.AnalyticsModuleSupport)
	}

	// Device capabilities
	if capabilities.Device.XAddr != "" {
		fmt.Println("\n[DEVICE CAPABILITIES]")
		fmt.Printf("XAddr: %s\n", capabilities.Device.XAddr)

		// System capabilities
		fmt.Println("  [System Capabilities]")
		fmt.Printf("  Discovery Resolve: %t\n", capabilities.Device.System.DiscoveryResolve)
		fmt.Printf("  Discovery Bye: %t\n", capabilities.Device.System.DiscoveryBye)
		fmt.Printf("  Remote Discovery: %t\n", capabilities.Device.System.RemoteDiscovery)
		fmt.Printf("  System Backup: %t\n", capabilities.Device.System.SystemBackup)
		fmt.Printf("  System Logging: %t\n", capabilities.Device.System.SystemLogging)
		fmt.Printf("  Firmware Upgrade: %t\n", capabilities.Device.System.FirmwareUpgrade)
		fmt.Printf("  Support Information: %t\n", capabilities.Device.System.SupportedVersions)

		// IO capabilities
		fmt.Println("  [IO Capabilities]")
		fmt.Printf("  Input Connectors: %d\n", capabilities.Device.IO.InputConnectors)
		fmt.Printf("  Relay Outputs: %d\n", capabilities.Device.IO.RelayOutputs)

		// Network capabilities
		fmt.Println("  [Network Capabilities]")
		fmt.Printf("  IP Filter: %t\n", capabilities.Device.Network.IPFilter)
		fmt.Printf("  Zero Configuration: %t\n", capabilities.Device.Network.ZeroConfiguration)
		fmt.Printf("  IP Version6: %t\n", capabilities.Device.Network.IPVersion6)
		fmt.Printf("  Dynamic DNS: %t\n", capabilities.Device.Network.DynDNS)
	}

	// Events capabilities
	if capabilities.Events.XAddr != "" {
		fmt.Println("\n[EVENT CAPABILITIES]")
		fmt.Printf("XAddr: %s\n", capabilities.Events.XAddr)
		fmt.Printf("WSSubscriptionPolicySupport: %t\n", capabilities.Events.WSSubscriptionPolicySupport)
		fmt.Printf("WSPullPointSupport: %t\n", capabilities.Events.WSPullPointSupport)
		fmt.Printf("WSPausableSubscriptionManagerInterfaceSupport: %t\n",
			capabilities.Events.WSPausableSubscriptionManagerInterfaceSupport)
	}

	// Imaging capabilities
	if capabilities.Imaging.XAddr != "" {
		fmt.Println("\n[IMAGING CAPABILITIES]")
		fmt.Printf("XAddr: %s\n", capabilities.Imaging.XAddr)
	}

	// Media capabilities
	if capabilities.Media.XAddr != "" {
		fmt.Println("\n[MEDIA CAPABILITIES]")
		fmt.Printf("XAddr: %s\n", capabilities.Media.XAddr)

		// Streaming capabilities
		fmt.Println("  [Streaming Capabilities]")
		fmt.Printf("  RTP Multicast: %t\n", capabilities.Media.StreamingCapabilities.RTPMulticast)
		fmt.Printf("  RTP/RTSP/TCP: %t\n", capabilities.Media.StreamingCapabilities.RTP_RTSP_TCP)
		fmt.Printf("  RTP/TCP: %t\n", capabilities.Media.StreamingCapabilities.RTP_TCP)
	}

	// PTZ capabilities
	if capabilities.PTZ.XAddr != "" {
		fmt.Println("\n[PTZ CAPABILITIES]")
		fmt.Printf("XAddr: %s\n", capabilities.PTZ.XAddr)
	}

	// Extension capabilities
	if !isEmptyStruct(capabilities.Extension) {
		fmt.Println("\n[EXTENSION CAPABILITIES]")

		// Only include extension sections if they have values
		if capabilities.Extension.DeviceIO.XAddr != "" {
			fmt.Println("  [Device IO Capabilities]")
			fmt.Printf("  XAddr: %s\n", capabilities.Extension.DeviceIO.XAddr)
			fmt.Printf("  Video Sources: %d\n", capabilities.Extension.DeviceIO.VideoSources)
			fmt.Printf("  Video Outputs: %d\n", capabilities.Extension.DeviceIO.VideoOutputs)
			fmt.Printf("  Audio Sources: %d\n", capabilities.Extension.DeviceIO.AudioSources)
			fmt.Printf("  Audio Outputs: %d\n", capabilities.Extension.DeviceIO.AudioOutputs)
			fmt.Printf("  Relay Outputs: %d\n", capabilities.Extension.DeviceIO.RelayOutputs)
		}

		if capabilities.Extension.Display.XAddr != "" {
			fmt.Println("  [Display Capabilities]")
			fmt.Printf("  XAddr: %s\n", capabilities.Extension.Display.XAddr)
		}

		if capabilities.Extension.Recording.XAddr != "" {
			fmt.Println("  [Recording Capabilities]")
			fmt.Printf("  XAddr: %s\n", capabilities.Extension.Recording.XAddr)
			fmt.Printf("  Media Profile Source: %t\n", capabilities.Extension.Recording.MediaProfileSource)
			fmt.Printf("  Dynamic Recordings: %t\n", capabilities.Extension.Recording.DynamicRecordings)
			fmt.Printf("  Dynamic Tracks: %t\n", capabilities.Extension.Recording.DynamicTracks)
		}

		if capabilities.Extension.Search.XAddr != "" {
			fmt.Println("  [Search Capabilities]")
			fmt.Printf("  XAddr: %s\n", capabilities.Extension.Search.XAddr)
			fmt.Printf("  Metadata Search: %t\n", capabilities.Extension.Search.MetadataSearch)
		}

		if capabilities.Extension.Replay.XAddr != "" {
			fmt.Println("  [Replay Capabilities]")
			fmt.Printf("  XAddr: %s\n", capabilities.Extension.Replay.XAddr)
		}
	}
}

// Helper function to check if a struct is effectively empty
func isEmptyStruct(obj interface{}) bool {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// If it's not a struct, don't consider it empty
	if v.Kind() != reflect.Struct {
		return false
	}

	// Check all fields to see if any have non-zero values
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if !field.IsZero() {
			return false
		}
	}

	return true
}
