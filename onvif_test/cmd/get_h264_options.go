package main

import (
	"bytes"
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/use-go/onvif"
	"github.com/use-go/onvif/device"
	"github.com/use-go/onvif/media"
	onvifXSD "github.com/use-go/onvif/xsd/onvif"
)

// Default camera connection details
const (
	defaultCameraIP = "192.168.1.12"
	defaultUsername = "admin"
	defaultPassword = "admin123"
	defaultTimeout  = 10 * time.Second
)

// Camera represents an ONVIF camera connection with its credentials
// Based on the working implementation in onvif_back/lib/client.go
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
	// This is the key difference - notice how the URL is formed without protocol or path
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

// XML wrapper structs for parsing responses
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

// VideoEncoderConfigurationsResponse has been updated to handle both string and integer encoding values
type VideoEncoderConfigurationsResponse struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    struct {
		GetVideoEncoderConfigurationsResponse struct {
			Configurations []struct {
				Name       string      `xml:"Name"`
				Token      string      `xml:"token,attr"`
				Encoding   interface{} `xml:"Encoding"` // Changed from int to interface{} to handle both string and int
				Quality    float64     `xml:"Quality"`
				Resolution struct {
					Width  int `xml:"Width"`
					Height int `xml:"Height"`
				} `xml:"Resolution"`
				RateControl struct {
					FrameRateLimit   int `xml:"FrameRateLimit"`
					EncodingInterval int `xml:"EncodingInterval"`
					BitrateLimit     int `xml:"BitrateLimit"`
				} `xml:"RateControl"`
			} `xml:"Configurations"`
		} `xml:"GetVideoEncoderConfigurationsResponse"`
	} `xml:"Body"`
}

type VideoEncoderConfigurationOptionsResponse struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    struct {
		GetVideoEncoderConfigurationOptionsResponse struct {
			Options struct {
				QualityRange struct {
					Min float64 `xml:"Min"`
					Max float64 `xml:"Max"`
				} `xml:"QualityRange"`
				H264 struct {
					ResolutionsAvailable []struct {
						Width  int `xml:"Width"`
						Height int `xml:"Height"`
					} `xml:"ResolutionsAvailable"`
					GovLengthRange struct {
						Min int `xml:"Min"`
						Max int `xml:"Max"`
					} `xml:"GovLengthRange"`
					FrameRateRange struct {
						Min int `xml:"Min"`
						Max int `xml:"Max"`
					} `xml:"FrameRateRange"`
					EncodingIntervalRange struct {
						Min int `xml:"Min"`
						Max int `xml:"Max"`
					} `xml:"EncodingIntervalRange"`
					H264ProfilesSupported []string `xml:"H264ProfilesSupported"`
				} `xml:"H264"`
				Extension struct {
					H264 struct {
						BitrateRange struct {
							Min int `xml:"Min"`
							Max int `xml:"Max"`
						} `xml:"BitrateRange"`
					} `xml:"H264"`
				} `xml:"Extension"`
			} `xml:"Options"`
		} `xml:"GetVideoEncoderConfigurationOptionsResponse"`
	} `xml:"Body"`
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

func main() {
	// Define flags for camera connection parameters
	ipPtr := flag.String("ip", defaultCameraIP, "Camera IP address")
	portPtr := flag.Int("port", 80, "Camera port")
	userPtr := flag.String("user", defaultUsername, "Username")
	passPtr := flag.String("pass", defaultPassword, "Password")
	configTokenPtr := flag.String("config", "", "Optional video encoder configuration token")
	profileTokenPtr := flag.String("profile", "", "Optional profile token")
	flag.Parse()

	fmt.Println("📹 ONVIF H264 Video Encoder Configuration Options 📹")
	fmt.Printf("Connecting to camera at %s:%d...\n", *ipPtr, *portPtr)

	// Create and connect to the camera using the same approach as in interactive_resolution_change.go
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

	// Get encoder configurations
	fmt.Println("\n🔍 Getting video encoder configurations...")
	getConfigsRequest := media.GetVideoEncoderConfigurations{}
	configsResp, err := dev.CallMethod(getConfigsRequest)
	if err != nil {
		log.Fatalf("❌ Could not get video encoder configurations: %v", err)
	}

	// Read and parse the configurations response
	rawConfigsXML, _ := ioutil.ReadAll(configsResp.Body)
	var configsInfo VideoEncoderConfigurationsResponse
	if err := xml.Unmarshal(rawConfigsXML, &configsInfo); err != nil {
		log.Fatalf("❌ Could not parse video encoder configurations: %v", err)
	}

	configs := configsInfo.Body.GetVideoEncoderConfigurationsResponse.Configurations
	if len(configs) == 0 {
		log.Fatalf("❌ No video encoder configurations found")
	}

	fmt.Println("\n===== Available Video Encoder Configurations =====")
	for i, config := range configs {
		fmt.Printf("%d. %s (Token: %s)\n", i+1, config.Name, config.Token)

		encoding := "Unknown"
		switch v := config.Encoding.(type) {
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
		case string:
			encoding = v
		}

		fmt.Printf("   Encoding: %s\n", encoding)
		fmt.Printf("   Resolution: %dx%d\n", config.Resolution.Width, config.Resolution.Height)
		fmt.Printf("   Bitrate: %d\n", config.RateControl.BitrateLimit)
		fmt.Printf("   Frame Rate: %d\n", config.RateControl.FrameRateLimit)
		fmt.Printf("   Quality: %.1f\n", config.Quality)
		fmt.Println()
	}

	// Determine which configuration and profile tokens to use
	configToken := *configTokenPtr
	profileToken := *profileTokenPtr

	// If no config token is provided, look for H264 configurations or use the first one
	if configToken == "" {
		for _, config := range configs {
			// Check for H264 encoding whether as int (2) or string ("H264")
			isH264 := false

			switch v := config.Encoding.(type) {
			case float64:
				isH264 = int(v) == 2
			case int:
				isH264 = v == 2
			case string:
				isH264 = v == "H264"
			}

			if isH264 {
				configToken = config.Token
				fmt.Printf("Found H264 configuration: %s (Token: %s)\n", config.Name, configToken)
				break
			}
		}

		// If still no H264 config found, use the first one
		if configToken == "" && len(configs) > 0 {
			configToken = configs[0].Token
			fmt.Printf("No H264 configuration found, using first available: %s (Token: %s)\n",
				configs[0].Name, configToken)
		}
	}

	// If no profile token is provided, use the first one
	if profileToken == "" && len(profiles) > 0 {
		profileToken = profiles[0].Token
		fmt.Printf("No profile token specified, using first one: %s (Token: %s)\n",
			profiles[0].Name, profileToken)
	}

	// Exit if we don't have the tokens we need
	if configToken == "" {
		fmt.Println("❌ No configuration token available")
		os.Exit(1)
	}

	if profileToken == "" {
		fmt.Println("❌ No profile token available")
		os.Exit(1)
	}

	// Create request to get video encoder configuration options
	fmt.Printf("\n🔍 Getting H264 options for config '%s' and profile '%s'...\n",
		configToken, profileToken)

	// Create the request with the correct reference token types
	request := media.GetVideoEncoderConfigurationOptions{
		ConfigurationToken: onvifXSD.ReferenceToken(configToken),
		ProfileToken:       onvifXSD.ReferenceToken(profileToken),
	}

	// Call the API method directly
	optionsResp, err := dev.CallMethod(request)
	if err != nil {
		log.Fatalf("❌ Failed to get video encoder configuration options: %v", err)
	}

	// Read and parse the options response
	rawOptionsXML, _ := ioutil.ReadAll(optionsResp.Body)

	// For debugging, you can uncomment this to print the raw XML
	// formattedXML, _ := formatXML(rawOptionsXML)
	// fmt.Println("Raw XML Response:", formattedXML)

	var optionsInfo VideoEncoderConfigurationOptionsResponse
	if err := xml.Unmarshal(rawOptionsXML, &optionsInfo); err != nil {
		log.Fatalf("❌ Could not parse video encoder configuration options: %v", err)
	}

	// Extract H264 options
	h264Options := optionsInfo.Body.GetVideoEncoderConfigurationOptionsResponse.Options.H264

	// Display H264 options
	fmt.Println("\n===== H264 Options =====")

	// Display available resolutions
	fmt.Println("\n🖥️ Available Resolutions:")
	if len(h264Options.ResolutionsAvailable) > 0 {
		for i, res := range h264Options.ResolutionsAvailable {
			fmt.Printf("%d. %dx%d\n", i+1, res.Width, res.Height)
		}
	} else {
		fmt.Println("   No resolution information available")
	}

	// Display GOP options
	fmt.Println("\n⏱️ GOP (Group of Pictures) Length Range:")
	fmt.Printf("   Min: %d, Max: %d\n",
		h264Options.GovLengthRange.Min, h264Options.GovLengthRange.Max)

	// Display frame rate options
	fmt.Println("\n🎞️ Frame Rate Range:")
	fmt.Printf("   Min: %d FPS, Max: %d FPS\n",
		h264Options.FrameRateRange.Min, h264Options.FrameRateRange.Max)

	// Display encoding interval options
	fmt.Println("\n⌛ Encoding Interval Range:")
	fmt.Printf("   Min: %d, Max: %d\n",
		h264Options.EncodingIntervalRange.Min, h264Options.EncodingIntervalRange.Max)

	// Display H264 profiles
	fmt.Println("\n📊 Supported H264 Profiles:")
	if len(h264Options.H264ProfilesSupported) > 0 {
		for i, profile := range h264Options.H264ProfilesSupported {
			fmt.Printf("%d. %s\n", i+1, profile)
		}
	} else {
		fmt.Println("   No H264 profile information available")
	}

	// Check for extended H264 options
	extension := optionsInfo.Body.GetVideoEncoderConfigurationOptionsResponse.Options.Extension
	fmt.Println("\n===== Extended H264 Options (if available) =====")

	// Display bitrate range if available
	fmt.Println("\n💾 Bitrate Range:")
	if extension.H264.BitrateRange.Max > 0 {
		fmt.Printf("   Min: %d kbps, Max: %d kbps\n",
			extension.H264.BitrateRange.Min, extension.H264.BitrateRange.Max)
	} else {
		fmt.Println("   No bitrate range information available from ONVIF API")
		fmt.Println("   Note: Check the camera's web interface for resolution-specific bitrate limits")
	}

	fmt.Println("\n✅ All H264 options retrieved successfully")
	fmt.Println("\nThis program reports only the information available through the ONVIF API.")
	fmt.Println("For more detailed, resolution-specific limitations, please consult your")
	fmt.Println("camera's web interface or manufacturer documentation.")
}
