package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/use-go/onvif"
	ondevice "github.com/use-go/onvif/device"
)

const (
	cameraIP = "192.168.1.12" // Replace with your camera's IP
	username = "admin"        // Replace with your camera's username
	password = "admin123"     // Replace with your camera's password
)

// Define structs to parse the ONVIF GetCapabilities response
type Envelope struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    Body     `xml:"Body"`
}

type Body struct {
	GetCapabilitiesResponse GetCapabilitiesResponse `xml:"GetCapabilitiesResponse"`
}

type GetCapabilitiesResponse struct {
	Capabilities Capabilities `xml:"Capabilities"`
}

type Capabilities struct {
	Analytics Analytics `xml:"Analytics"`
	Device    Device    `xml:"Device"`
	Events    Events    `xml:"Events"`
	Imaging   Imaging   `xml:"Imaging"`
	Media     Media     `xml:"Media"`
	Extension Extension `xml:"Extension"`
}

type Analytics struct {
	XAddr                  string `xml:"XAddr"`
	RuleSupport            bool   `xml:"RuleSupport"`
	AnalyticsModuleSupport bool   `xml:"AnalyticsModuleSupport"`
}

type Device struct {
	XAddr    string   `xml:"XAddr"`
	Network  Network  `xml:"Network"`
	System   System   `xml:"System"`
	IO       IO       `xml:"IO"`
	Security Security `xml:"Security"`
}

type Network struct {
	IPFilter          bool `xml:"IPFilter"`
	ZeroConfiguration bool `xml:"ZeroConfiguration"`
	IPVersion6        bool `xml:"IPVersion6"`
	DynDNS            bool `xml:"DynDNS"`
	Extension         struct {
		Dot11Configuration bool `xml:"Dot11Configuration"`
	} `xml:"Extension"`
}

type System struct {
	DiscoveryResolve  bool `xml:"DiscoveryResolve"`
	DiscoveryBye      bool `xml:"DiscoveryBye"`
	RemoteDiscovery   bool `xml:"RemoteDiscovery"`
	SystemBackup      bool `xml:"SystemBackup"`
	SystemLogging     bool `xml:"SystemLogging"`
	FirmwareUpgrade   bool `xml:"FirmwareUpgrade"`
	SupportedVersions []struct {
		Major int `xml:"Major"`
		Minor int `xml:"Minor"`
	} `xml:"SupportedVersions"`
	Extension struct {
		HttpFirmwareUpgrade    bool `xml:"HttpFirmwareUpgrade"`
		HttpSystemBackup       bool `xml:"HttpSystemBackup"`
		HttpSystemLogging      bool `xml:"HttpSystemLogging"`
		HttpSupportInformation bool `xml:"HttpSupportInformation"`
	} `xml:"Extension"`
}

type IO struct {
	InputConnectors int `xml:"InputConnectors"`
	RelayOutputs    int `xml:"RelayOutputs"`
	Extension       struct {
		Auxiliary bool `xml:"Auxiliary"`
	} `xml:"Extension"`
}

type Security struct {
	TLS11                bool `xml:"TLS1.1"`
	TLS12                bool `xml:"TLS1.2"`
	OnboardKeyGeneration bool `xml:"OnboardKeyGeneration"`
	AccessPolicyConfig   bool `xml:"AccessPolicyConfig"`
	X509Token            bool `xml:"X.509Token"`
	SAMLToken            bool `xml:"SAMLToken"`
	KerberosToken        bool `xml:"KerberosToken"`
	RELToken             bool `xml:"RELToken"`
	Extension            struct {
		TLS10     bool `xml:"TLS1.0"`
		Extension struct {
			Dot1X              bool `xml:"Dot1X"`
			SupportedEAPMethod int  `xml:"SupportedEAPMethod"`
			RemoteUserHandling bool `xml:"RemoteUserHandling"`
		} `xml:"Extension"`
	} `xml:"Extension"`
}

type Events struct {
	XAddr                                         string `xml:"XAddr"`
	WSSubscriptionPolicySupport                   bool   `xml:"WSSubscriptionPolicySupport"`
	WSPullPointSupport                            bool   `xml:"WSPullPointSupport"`
	WSPausableSubscriptionManagerInterfaceSupport bool   `xml:"WSPausableSubscriptionManagerInterfaceSupport"`
}

type Imaging struct {
	XAddr string `xml:"XAddr"`
}

type Media struct {
	XAddr                 string `xml:"XAddr"`
	StreamingCapabilities struct {
		RTPMulticast bool `xml:"RTPMulticast"`
		RTP_TCP      bool `xml:"RTP_TCP"`
		RTP_RTSP_TCP bool `xml:"RTP_RTSP_TCP"`
	} `xml:"StreamingCapabilities"`
	Extension struct {
		ProfileCapabilities struct {
			MaximumNumberOfProfiles int `xml:"MaximumNumberOfProfiles"`
		} `xml:"ProfileCapabilities"`
	} `xml:"Extension"`
}

type Extension struct {
	DeviceIO struct {
		XAddr        string `xml:"XAddr"`
		VideoSources int    `xml:"VideoSources"`
		VideoOutputs int    `xml:"VideoOutputs"`
		AudioSources int    `xml:"AudioSources"`
		AudioOutputs int    `xml:"AudioOutputs"`
		RelayOutputs int    `xml:"RelayOutputs"`
	} `xml:"DeviceIO"`
	Extensions struct {
		TelexCapabilities struct {
			XAddr                 string `xml:"XAddr"`
			TimeOSDSupport        bool   `xml:"TimeOSDSupport"`
			TitleOSDSupport       bool   `xml:"TitleOSDSupport"`
			PTZ3DZoomSupport      bool   `xml:"PTZ3DZoomSupport"`
			PTZAuxSwitchSupport   bool   `xml:"PTZAuxSwitchSupport"`
			MotionDetectorSupport bool   `xml:"MotionDetectorSupport"`
			TamperDetectorSupport bool   `xml:"TamperDetectorSupport"`
		} `xml:"TelexCapabilities"`
	} `xml:"Extensions"`
}

// Device information structs
type DeviceInformationResponse struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    struct {
		GetDeviceInformationResponse struct {
			Manufacturer    string `xml:"Manufacturer"`
			Model           string `xml:"Model"`
			FirmwareVersion string `xml:"FirmwareVersion"`
			SerialNumber    string `xml:"SerialNumber"`
			HardwareId      string `xml:"HardwareId"`
		} `xml:"GetDeviceInformationResponse"`
	} `xml:"Body"`
}

func main() {
	// Create a new ONVIF device
	device, err := onvif.NewDevice(onvif.DeviceParams{
		Xaddr:    "192.168.1.12:80",
		Username: username,
		Password: password,
	})
	if err != nil {
		log.Fatalf("Failed to connect to the device: %v", err)
	}

	// Successfully connected to the device
	fmt.Println("Connected to the device: ", device)

	// Get the capabilities of the device
	capabilitiesRequest := ondevice.GetCapabilities{Category: "All"}
	response, err := device.CallMethod(capabilitiesRequest)
	if err != nil {
		log.Fatalf("Failed to get capabilities: %v", err)
	}

	// Read the raw XML response
	rawXML, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("Failed to read raw XML response: %v", err)
	}

	// Parse the XML into our struct
	var envelope Envelope
	if err := xml.Unmarshal(rawXML, &envelope); err != nil {
		log.Fatalf("Failed to parse XML: %v", err)
	}

	// Print the structured capabilities
	printCapabilities(envelope.Body.GetCapabilitiesResponse.Capabilities)

	// Get device information using the library's built-in struct
	deviceInfoRequest := ondevice.GetDeviceInformation{}
	deviceInfoResponse, err := device.CallMethod(deviceInfoRequest)
	if err != nil {
		log.Fatalf("Failed to get device information: %v", err)
	}

	// Read the raw XML response for device information
	rawDeviceInfoXML, err := ioutil.ReadAll(deviceInfoResponse.Body)
	if err != nil {
		log.Fatalf("Failed to read device information response: %v", err)
	}

	// Parse the XML into our struct
	var deviceInfo DeviceInformationResponse
	if err := xml.Unmarshal(rawDeviceInfoXML, &deviceInfo); err != nil {
		log.Fatalf("Failed to parse device information XML: %v", err)
	}

	// Print the device information
	printDeviceInformation(deviceInfo.Body.GetDeviceInformationResponse)
}

// Helper function to format XML
func formatXML(input []byte) (string, error) {
	var buffer bytes.Buffer
	decoder := xml.NewDecoder(bytes.NewReader(input))
	encoder := xml.NewEncoder(&buffer)
	encoder.Indent("", "  ")

	for {
		token, err := decoder.Token()
		if err != nil {
			break
		}
		if err := encoder.EncodeToken(token); err != nil {
			return "", err
		}
	}

	if err := encoder.Flush(); err != nil {
		return "", err
	}
	return buffer.String(), nil
}

// Print capabilities in a hierarchical, human-readable format
func printCapabilities(c Capabilities) {
	fmt.Println("\n======== CAMERA CAPABILITIES ========")

	fmt.Println("\n📹 ANALYTICS:")
	fmt.Printf("  • XAddr: %s\n", c.Analytics.XAddr)
	fmt.Printf("  • Rule Support: %t\n", c.Analytics.RuleSupport)
	fmt.Printf("  • Analytics Module Support: %t\n", c.Analytics.AnalyticsModuleSupport)

	fmt.Println("\n🔧 DEVICE:")
	fmt.Printf("  • XAddr: %s\n", c.Device.XAddr)

	fmt.Println("  • Network:")
	fmt.Printf("    - IP Filter: %t\n", c.Device.Network.IPFilter)
	fmt.Printf("    - Zero Configuration: %t\n", c.Device.Network.ZeroConfiguration)
	fmt.Printf("    - IPv6: %t\n", c.Device.Network.IPVersion6)
	fmt.Printf("    - DynDNS: %t\n", c.Device.Network.DynDNS)
	fmt.Printf("    - Dot11 Configuration: %t\n", c.Device.Network.Extension.Dot11Configuration)

	fmt.Println("  • System:")
	fmt.Printf("    - Discovery Resolve: %t\n", c.Device.System.DiscoveryResolve)
	fmt.Printf("    - Discovery Bye: %t\n", c.Device.System.DiscoveryBye)
	fmt.Printf("    - Remote Discovery: %t\n", c.Device.System.RemoteDiscovery)
	fmt.Printf("    - System Backup: %t\n", c.Device.System.SystemBackup)
	fmt.Printf("    - System Logging: %t\n", c.Device.System.SystemLogging)
	fmt.Printf("    - Firmware Upgrade: %t\n", c.Device.System.FirmwareUpgrade)

	fmt.Println("    - Supported Versions:")
	for _, v := range c.Device.System.SupportedVersions {
		fmt.Printf("      > %d.%d\n", v.Major, v.Minor)
	}

	fmt.Println("    - HTTP Extensions:")
	fmt.Printf("      > HTTP Firmware Upgrade: %t\n", c.Device.System.Extension.HttpFirmwareUpgrade)
	fmt.Printf("      > HTTP System Backup: %t\n", c.Device.System.Extension.HttpSystemBackup)
	fmt.Printf("      > HTTP System Logging: %t\n", c.Device.System.Extension.HttpSystemLogging)
	fmt.Printf("      > HTTP Support Information: %t\n", c.Device.System.Extension.HttpSupportInformation)

	fmt.Println("  • IO:")
	fmt.Printf("    - Input Connectors: %d\n", c.Device.IO.InputConnectors)
	fmt.Printf("    - Relay Outputs: %d\n", c.Device.IO.RelayOutputs)
	fmt.Printf("    - Auxiliary: %t\n", c.Device.IO.Extension.Auxiliary)

	fmt.Println("  • Security:")
	fmt.Printf("    - TLS 1.0: %t\n", c.Device.Security.Extension.TLS10)
	fmt.Printf("    - TLS 1.1: %t\n", c.Device.Security.TLS11)
	fmt.Printf("    - TLS 1.2: %t\n", c.Device.Security.TLS12)
	fmt.Printf("    - Onboard Key Generation: %t\n", c.Device.Security.OnboardKeyGeneration)
	fmt.Printf("    - Access Policy Config: %t\n", c.Device.Security.AccessPolicyConfig)
	fmt.Printf("    - X.509 Token: %t\n", c.Device.Security.X509Token)
	fmt.Printf("    - SAML Token: %t\n", c.Device.Security.SAMLToken)
	fmt.Printf("    - Kerberos Token: %t\n", c.Device.Security.KerberosToken)
	fmt.Printf("    - REL Token: %t\n", c.Device.Security.RELToken)
	fmt.Printf("    - Dot1X: %t\n", c.Device.Security.Extension.Extension.Dot1X)
	fmt.Printf("    - Supported EAP Method: %d\n", c.Device.Security.Extension.Extension.SupportedEAPMethod)
	fmt.Printf("    - Remote User Handling: %t\n", c.Device.Security.Extension.Extension.RemoteUserHandling)

	fmt.Println("\n📢 EVENTS:")
	fmt.Printf("  • XAddr: %s\n", c.Events.XAddr)
	fmt.Printf("  • WS Subscription Policy Support: %t\n", c.Events.WSSubscriptionPolicySupport)
	fmt.Printf("  • WS Pull Point Support: %t\n", c.Events.WSPullPointSupport)
	fmt.Printf("  • WS Pausable Subscription Manager Interface Support: %t\n", c.Events.WSPausableSubscriptionManagerInterfaceSupport)

	fmt.Println("\n🎨 IMAGING:")
	fmt.Printf("  • XAddr: %s\n", c.Imaging.XAddr)

	fmt.Println("\n🎬 MEDIA:")
	fmt.Printf("  • XAddr: %s\n", c.Media.XAddr)
	fmt.Println("  • Streaming Capabilities:")
	fmt.Printf("    - RTP Multicast: %t\n", c.Media.StreamingCapabilities.RTPMulticast)
	fmt.Printf("    - RTP TCP: %t\n", c.Media.StreamingCapabilities.RTP_TCP)
	fmt.Printf("    - RTP RTSP TCP: %t\n", c.Media.StreamingCapabilities.RTP_RTSP_TCP)
	fmt.Printf("  • Maximum Number of Profiles: %d\n", c.Media.Extension.ProfileCapabilities.MaximumNumberOfProfiles)

	fmt.Println("\n📎 EXTENSIONS:")
	fmt.Println("  • DeviceIO:")
	fmt.Printf("    - XAddr: %s\n", c.Extension.DeviceIO.XAddr)
	fmt.Printf("    - Video Sources: %d\n", c.Extension.DeviceIO.VideoSources)
	fmt.Printf("    - Video Outputs: %d\n", c.Extension.DeviceIO.VideoOutputs)
	fmt.Printf("    - Audio Sources: %d\n", c.Extension.DeviceIO.AudioSources)
	fmt.Printf("    - Audio Outputs: %d\n", c.Extension.DeviceIO.AudioOutputs)
	fmt.Printf("    - Relay Outputs: %d\n", c.Extension.DeviceIO.RelayOutputs)

	fmt.Println("  • Telex Capabilities:")
	fmt.Printf("    - XAddr: %s\n", c.Extension.Extensions.TelexCapabilities.XAddr)
	fmt.Printf("    - Time OSD Support: %t\n", c.Extension.Extensions.TelexCapabilities.TimeOSDSupport)
	fmt.Printf("    - Title OSD Support: %t\n", c.Extension.Extensions.TelexCapabilities.TitleOSDSupport)
	fmt.Printf("    - PTZ 3D Zoom Support: %t\n", c.Extension.Extensions.TelexCapabilities.PTZ3DZoomSupport)
	fmt.Printf("    - PTZ Aux Switch Support: %t\n", c.Extension.Extensions.TelexCapabilities.PTZAuxSwitchSupport)
	fmt.Printf("    - Motion Detector Support: %t\n", c.Extension.Extensions.TelexCapabilities.MotionDetectorSupport)
	fmt.Printf("    - Tamper Detector Support: %t\n", c.Extension.Extensions.TelexCapabilities.TamperDetectorSupport)
}

// Print device information in a human-readable format
func printDeviceInformation(info struct {
	Manufacturer    string `xml:"Manufacturer"`
	Model           string `xml:"Model"`
	FirmwareVersion string `xml:"FirmwareVersion"`
	SerialNumber    string `xml:"SerialNumber"`
	HardwareId      string `xml:"HardwareId"`
}) {
	fmt.Println("\n======== DEVICE INFORMATION ========")
	fmt.Printf("Manufacturer: %s\n", info.Manufacturer)
	fmt.Printf("Model: %s\n", info.Model)
	fmt.Printf("Firmware Version: %s\n", info.FirmwareVersion)
	fmt.Printf("Serial Number: %s\n", info.SerialNumber)
	fmt.Printf("Hardware ID: %s\n", info.HardwareId)
}
