package main

import (
	"bytes"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/use-go/onvif"
	"github.com/use-go/onvif/device"
	"github.com/use-go/onvif/media"
	onvifXSD "github.com/use-go/onvif/xsd/onvif"
)

// Default camera connection details
const (
	defaultCameraIP = "192.168.1.31"
	defaultUsername = "admin"
	defaultPassword = "Admin123"
	defaultTimeout  = 10 * time.Second
)

// Camera represents an ONVIF camera connection with its credentials
type Camera struct {
	IP       string
	Port     int
	Username string
	Password string
	Device   *onvif.Device
}

// NewCamera creates a new Camera instance
func NewCamera(ip string, port int, username, password string) *Camera {
	return &Camera{
		IP:       ip,
		Port:     port,
		Username: username,
		Password: password,
	}
}

// Connect establishes a connection to the camera
func (c *Camera) Connect() error {
	dev, err := onvif.NewDevice(onvif.DeviceParams{
		Xaddr:    fmt.Sprintf("%s:%d", c.IP, c.Port),
		Username: c.Username,
		Password: c.Password,
	})

	if err != nil {
		return fmt.Errorf("failed to connect to camera: %v", err)
	}

	c.Device = dev
	return nil
}

// GetDevice returns the underlying ONVIF device
func (c *Camera) GetDevice() *onvif.Device {
	return c.Device
}

// Response structures for parsing XML responses
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

type ProfilesResponse struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    struct {
		GetProfilesResponse struct {
			Profiles []struct {
				Name  string `xml:"Name"`
				Token string `xml:"token,attr"`
			} `xml:"Profiles"`
		} `xml:"GetProfilesResponse"`
	} `xml:"Body"`
}

// CompatibleVideoEncoderConfigurationsResponse structure
type CompatibleVideoEncoderConfigurationsResponse struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    struct {
		GetCompatibleVideoEncoderConfigurationsResponse struct {
			Configurations []struct {
				Token               string      `xml:"token,attr"`
				Name                string      `xml:"Name"`
				UseCount            int         `xml:"UseCount"`
				Encoding            interface{} `xml:"Encoding"` // Could be string or int
				GuaranteedFrameRate bool        `xml:"GuaranteedFrameRate"`
				Resolution          struct {
					Width  int `xml:"Width"`
					Height int `xml:"Height"`
				} `xml:"Resolution"`
				Quality     float64 `xml:"Quality"`
				RateControl struct {
					FrameRateLimit   int `xml:"FrameRateLimit"`
					EncodingInterval int `xml:"EncodingInterval"`
					BitrateLimit     int `xml:"BitrateLimit"`
				} `xml:"RateControl"`
				H264 struct {
					GovLength   int    `xml:"GovLength"`
					H264Profile string `xml:"H264Profile"`
				} `xml:"H264"`
				MPEG4 struct {
					GovLength    int    `xml:"GovLength"`
					Mpeg4Profile string `xml:"Mpeg4Profile"`
				} `xml:"MPEG4"`
				Multicast struct {
					Address struct {
						Type        string `xml:"Type"`
						IPv4Address string `xml:"IPv4Address"`
						IPv6Address string `xml:"IPv6Address"`
					} `xml:"Address"`
					Port      int  `xml:"Port"`
					TTL       int  `xml:"TTL"`
					AutoStart bool `xml:"AutoStart"`
				} `xml:"Multicast"`
				SessionTimeout string `xml:"SessionTimeout"`
			} `xml:"Configurations"`
		} `xml:"GetCompatibleVideoEncoderConfigurationsResponse"`
	} `xml:"Body"`
}

// Helper function to format XML
func formatXML(input []byte) (string, error) {
	var buf strings.Builder
	decoder := xml.NewDecoder(bytes.NewReader(input))
	encoder := xml.NewEncoder(&buf)
	encoder.Indent("", "  ")

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
		if err := encoder.EncodeToken(token); err != nil {
			return "", err
		}
	}

	if err := encoder.Flush(); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func main() {
	// Define flags for camera connection parameters
	ipPtr := flag.String("ip", defaultCameraIP, "Camera IP address")
	portPtr := flag.Int("port", 80, "Camera port")
	userPtr := flag.String("user", defaultUsername, "Username")
	passPtr := flag.String("pass", defaultPassword, "Password")
	profileTokenPtr := flag.String("profile", "", "Optional profile token")
	flag.Parse()

	fmt.Println("📹 ONVIF Get Compatible Video Encoder Configurations 📹")
	fmt.Printf("Connecting to camera at %s:%d...\n", *ipPtr, *portPtr)

	// Create and connect to the camera
	camera := NewCamera(*ipPtr, *portPtr, *userPtr, *passPtr)
	err := camera.Connect()
	if err != nil {
		log.Fatalf("❌ Could not connect to the camera: %v", err)
	}
	fmt.Println("✅ Connected to camera successfully")

	// Get the underlying device object for API calls
	dev := camera.GetDevice()

	// Get device information
	fmt.Println("\n🔍 Getting device information...")
	deviceInfoRequest := device.GetDeviceInformation{}
	deviceInfoResp, err := dev.CallMethod(deviceInfoRequest)
	if err != nil {
		log.Printf("⚠️ Could not get device information: %v", err)
	} else {
		// Read and parse the response
		rawXML, _ := ioutil.ReadAll(deviceInfoResp.Body)
		var deviceInfo DeviceInformationResponse
		if err := xml.Unmarshal(rawXML, &deviceInfo); err != nil {
			log.Printf("⚠️ Could not parse device information: %v", err)
		} else {
			info := deviceInfo.Body.GetDeviceInformationResponse
			fmt.Println("\n===== Device Information =====")
			fmt.Printf("Manufacturer: %s\n", info.Manufacturer)
			fmt.Printf("Model: %s\n", info.Model)
			fmt.Printf("Firmware Version: %s\n", info.FirmwareVersion)
			fmt.Printf("Serial Number: %s\n", info.SerialNumber)
			fmt.Printf("Hardware ID: %s\n", info.HardwareId)
		}
	}

	// Get profiles to retrieve available token information
	fmt.Println("\n🔍 Getting media profiles...")
	getProfilesRequest := media.GetProfiles{}
	profilesResp, err := dev.CallMethod(getProfilesRequest)
	if err != nil {
		log.Fatalf("❌ Could not get profiles: %v", err)
	}

	// Read and parse the profiles response
	rawProfilesXML, _ := ioutil.ReadAll(profilesResp.Body)
	var profilesInfo ProfilesResponse
	if err := xml.Unmarshal(rawProfilesXML, &profilesInfo); err != nil {
		log.Fatalf("❌ Could not parse profiles: %v", err)
	}

	profiles := profilesInfo.Body.GetProfilesResponse.Profiles
	if len(profiles) == 0 {
		log.Fatalf("❌ No profiles found")
	}

	// Print all profiles
	fmt.Println("\n===== Available Profiles =====")
	for i, profile := range profiles {
		fmt.Printf("%d. %s (Token: %s)\n", i+1, profile.Name, profile.Token)
	}

	// Determine which profile token to use
	profileToken := *profileTokenPtr

	// If no profile token is provided, let user choose
	if profileToken == "" {
		if len(profiles) == 1 {
			// Use the only available profile
			profileToken = profiles[0].Token
			fmt.Printf("\nUsing the only available profile: %s (Token: %s)\n",
				profiles[0].Name, profileToken)
		} else {
			// Ask user to select a profile
			fmt.Println("\nPlease select a profile by number:")
			var choice int
			fmt.Scanln(&choice)

			if choice < 1 || choice > len(profiles) {
				log.Fatalf("❌ Invalid choice")
			}

			profileToken = profiles[choice-1].Token
			fmt.Printf("\nUsing profile: %s (Token: %s)\n",
				profiles[choice-1].Name, profileToken)
		}
	}

	// Create the GetCompatibleVideoEncoderConfigurations request
	fmt.Printf("\n🔍 Getting compatible video encoder configurations for profile '%s'...\n", profileToken)

	// Create a request for GetCompatibleVideoEncoderConfigurations
	request := media.GetCompatibleVideoEncoderConfigurations{
		ProfileToken: onvifXSD.ReferenceToken(profileToken),
	}

	// Call the API method
	response, err := dev.CallMethod(request)

	// Read and parse the response
	rawXML, _ := ioutil.ReadAll(response.Body)

	// Uncomment to see raw XML for debugging
	//fmt.Println("\nRaw XML Response:")
	//formattedXML, _ := formatXML(rawXML)
	//fmt.Println(formattedXML)

	var compatibleInfo CompatibleVideoEncoderConfigurationsResponse
	if err := xml.Unmarshal(rawXML, &compatibleInfo); err != nil {
		log.Fatalf("❌ Could not parse compatible video encoder configurations: %v", err)
	}

	configs := compatibleInfo.Body.GetCompatibleVideoEncoderConfigurationsResponse.Configurations

	// Display the results
	if len(configs) == 0 {
		fmt.Println("\n❌ No compatible video encoder configurations found for this profile")
	} else {
		fmt.Printf("\n===== Compatible Video Encoder Configurations for Profile '%s' =====\n", profileToken)
		fmt.Printf("Found %d compatible configurations\n", len(configs))

		for i, config := range configs {
			fmt.Printf("\n--- Configuration %d ---\n", i+1)
			fmt.Printf("Token: %s\n", config.Token)
			fmt.Printf("Name: %s\n", config.Name)
			fmt.Printf("Use Count: %d\n", config.UseCount)

			// Handle encoding which can be string or int
			encoding := "Unknown"
			switch v := config.Encoding.(type) {
			case string:
				encoding = v
			case float64:
				switch int(v) {
				case 0:
					encoding = "JPEG"
				case 1:
					encoding = "MPEG4"
				case 2:
					encoding = "H264"
				}
			case int:
				switch v {
				case 0:
					encoding = "JPEG"
				case 1:
					encoding = "MPEG4"
				case 2:
					encoding = "H264"
				}
			}

			fmt.Printf("Encoding: %s\n", encoding)
			fmt.Printf("Resolution: %dx%d\n", config.Resolution.Width, config.Resolution.Height)
			fmt.Printf("Quality: %.1f\n", config.Quality)
			fmt.Printf("Guaranteed Frame Rate: %v\n", config.GuaranteedFrameRate)

			// Rate control info
			fmt.Printf("Frame Rate Limit: %d\n", config.RateControl.FrameRateLimit)
			fmt.Printf("Encoding Interval: %d\n", config.RateControl.EncodingInterval)
			fmt.Printf("Bitrate Limit: %d kbps\n", config.RateControl.BitrateLimit)

			// H264-specific parameters if applicable
			if encoding == "H264" && config.H264.GovLength > 0 {
				fmt.Printf("GOV Length: %d\n", config.H264.GovLength)
				fmt.Printf("H264 Profile: %s\n", config.H264.H264Profile)
			}

			// MPEG4-specific parameters if applicable
			if encoding == "MPEG4" && config.MPEG4.GovLength > 0 {
				fmt.Printf("GOV Length: %d\n", config.MPEG4.GovLength)
				fmt.Printf("MPEG4 Profile: %s\n", config.MPEG4.Mpeg4Profile)
			}

			// Multicast configuration
			fmt.Printf("Multicast Address Type: %s\n", config.Multicast.Address.Type)
			if config.Multicast.Address.IPv4Address != "" {
				fmt.Printf("Multicast IPv4: %s\n", config.Multicast.Address.IPv4Address)
			}
			if config.Multicast.Address.IPv6Address != "" {
				fmt.Printf("Multicast IPv6: %s\n", config.Multicast.Address.IPv6Address)
			}
			fmt.Printf("Multicast Port: %d\n", config.Multicast.Port)
			fmt.Printf("Multicast TTL: %d\n", config.Multicast.TTL)
			fmt.Printf("Multicast AutoStart: %v\n", config.Multicast.AutoStart)

			// Session timeout
			fmt.Printf("Session Timeout: %s\n", config.SessionTimeout)
		}
	}

	fmt.Println("\n✅ Operation completed")
}
