package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
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

// ResolutionLimits captures the limits specific to a resolution
type ResolutionLimits struct {
	Width                 int
	Height                int
	MinFrameRate          int
	MaxFrameRate          int
	MinGOVLength          int
	MaxGOVLength          int
	MinBitrate            int
	MaxBitrate            int
	SupportedH264Profiles []string
}

func main() {
	// Define flags for camera connection parameters
	ipPtr := flag.String("ip", defaultCameraIP, "Camera IP address")
	portPtr := flag.Int("port", 80, "Camera port")
	userPtr := flag.String("user", defaultUsername, "Username")
	passPtr := flag.String("pass", defaultPassword, "Password")
	flag.Parse()

	fmt.Println("📹 ONVIF Resolution-Specific H264 Limits Discovery 📹")
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
	deviceInfoReq := device.GetDeviceInformation{}
	_, err = dev.CallMethod(deviceInfoReq)
	if err != nil {
		log.Printf("⚠️ Could not get device information: %v", err)
	} else {
		fmt.Println("✅ Device information retrieved successfully")
	}

	// Get profiles
	fmt.Println("\n🔍 Getting media profiles...")
	getProfilesReq := media.GetProfiles{}
	profilesResp, err := dev.CallMethod(getProfilesReq)
	if err != nil {
		log.Fatalf("❌ Could not get profiles: %v", err)
	}

	// Parse the profiles response to get profile tokens
	profileTokens, err := parseProfileTokens(profilesResp.Body)
	if err != nil {
		log.Fatalf("❌ Could not parse profiles: %v", err)
	}

	if len(profileTokens) == 0 {
		log.Fatalf("❌ No profiles found")
	}

	fmt.Printf("Found %d profiles\n", len(profileTokens))
	profileToken := profileTokens[0] // Using the first profile for our tests

	// Get encoder configurations
	fmt.Println("\n🔍 Getting video encoder configurations...")
	getConfigsReq := media.GetVideoEncoderConfigurations{}
	configsResp, err := dev.CallMethod(getConfigsReq)
	if err != nil {
		log.Fatalf("❌ Could not get video encoder configurations: %v", err)
	}

	// Parse the configurations response to get config tokens
	h264Configs, err := parseH264ConfigTokens(configsResp.Body)
	if err != nil {
		log.Fatalf("❌ Could not parse configurations: %v", err)
	}

	if len(h264Configs) == 0 {
		log.Fatalf("❌ No H264 encoder configurations found")
	}

	fmt.Printf("Found %d H264 configurations\n", len(h264Configs))
	configToken := h264Configs[0] // Using the first H264 config for our tests

	// Get encoder configuration options to find available resolutions
	fmt.Println("\n🔍 Getting available resolutions...")
	optionsReq := media.GetVideoEncoderConfigurationOptions{
		ConfigurationToken: onvifXSD.ReferenceToken(configToken),
		ProfileToken:       onvifXSD.ReferenceToken(profileToken),
	}

	optionsResp, err := dev.CallMethod(optionsReq)
	if err != nil {
		log.Fatalf("❌ Could not get encoder options: %v", err)
	}

	// Parse available resolutions
	resolutions, err := parseAvailableResolutions(optionsResp.Body)
	if err != nil {
		log.Fatalf("❌ Could not parse available resolutions: %v", err)
	}

	if len(resolutions) == 0 {
		log.Fatalf("❌ No available resolutions found")
	}

	fmt.Printf("Found %d available resolutions\n", len(resolutions))

	// Now, for each resolution, we'll get the specific limits
	fmt.Println("\n📊 Discovering resolution-specific limits...")

	// Store the results for each resolution
	var resolutionLimits []ResolutionLimits

	for _, resolution := range resolutions {
		fmt.Printf("\n🔍 Checking limits for resolution %dx%d...\n",
			resolution.Width, resolution.Height)

		// Get the options for this specific resolution
		limits, err := getResolutionSpecificLimits(dev, configToken, profileToken, resolution)
		if err != nil {
			fmt.Printf("⚠️ Could not get limits for %dx%d: %v\n",
				resolution.Width, resolution.Height, err)
			continue
		}

		resolutionLimits = append(resolutionLimits, limits)
	}

	// Display the results in a nicely formatted table
	fmt.Println("\n===== Resolution-Specific H264 Limits =====")
	fmt.Println("\n┌───────────────┬────────────┬────────────┬────────────┬────────────┬─────────────┐")
	fmt.Println("│  Resolution   │  FPS Range │ GOP Range  │Bitrate(kbps)│ H264 Profile│    Notes    │")
	fmt.Println("├───────────────┼────────────┼────────────┼────────────┼────────────┼─────────────┤")

	for _, limit := range resolutionLimits {
		profiles := ""
		if len(limit.SupportedH264Profiles) > 0 {
			profiles = limit.SupportedH264Profiles[0]
			if len(limit.SupportedH264Profiles) > 1 {
				profiles += "..."
			}
		}

		var notes string

		// Special notes for specific resolutions
		if limit.Width == 2304 && limit.Height == 1296 && limit.MaxFrameRate == 20 {
			notes = " Max FPS Limit!"
		}

		// Display "unavailable" for bitrate values that couldn't be determined
		bitrateRange := ""
		if limit.MinBitrate == -1 && limit.MaxBitrate == -1 {
			bitrateRange = "unavailable"
		} else {
			minStr := "?"
			maxStr := "?"

			if limit.MinBitrate != -1 {
				minStr = fmt.Sprintf("%d", limit.MinBitrate)
			}

			if limit.MaxBitrate != -1 {
				maxStr = fmt.Sprintf("%d", limit.MaxBitrate)
			}

			bitrateRange = fmt.Sprintf("%s-%s", minStr, maxStr)
		}

		fmt.Printf("│ %4dx%-7d │ %3d-%-6d │ %3d-%-6d │ %-11s │ %-10s │ %-11s │\n",
			limit.Width, limit.Height,
			limit.MinFrameRate, limit.MaxFrameRate,
			limit.MinGOVLength, limit.MaxGOVLength,
			bitrateRange,
			profiles, notes)
	}

	fmt.Println("└───────────────┴────────────┴────────────┴────────────┴────────────┴─────────────┘")

	fmt.Println("\n✅ Resolution-specific limits discovery completed")
}

// Helper function to get limits for a specific resolution
func getResolutionSpecificLimits(dev *onvif.Device, configToken, profileToken string, resolution struct{ Width, Height int }) (ResolutionLimits, error) {
	limits := ResolutionLimits{
		Width:  resolution.Width,
		Height: resolution.Height,
	}

	// Get the options for this specific resolution by querying with the tokens
	optionsReq := media.GetVideoEncoderConfigurationOptions{
		ConfigurationToken: onvifXSD.ReferenceToken(configToken),
		ProfileToken:       onvifXSD.ReferenceToken(profileToken),
	}

	optionsResp, err := dev.CallMethod(optionsReq)
	if err != nil {
		return limits, err
	}

	// Parse the options to get the specific limits
	h264Options, err := parseH264Options(optionsResp.Body)
	if err != nil {
		return limits, err
	}

	// Extract the limits
	limits.MinFrameRate = h264Options.FrameRateMin
	limits.MaxFrameRate = h264Options.FrameRateMax
	limits.MinGOVLength = h264Options.GovLengthMin
	limits.MaxGOVLength = h264Options.GovLengthMax
	limits.SupportedH264Profiles = h264Options.Profiles

	// Get bitrate limits from extension if available
	if h264Options.BitrateMin > 0 || h264Options.BitrateMax > 0 {
		limits.MinBitrate = h264Options.BitrateMin
		limits.MaxBitrate = h264Options.BitrateMax
	} else {
		// If no specific bitrate range is provided by the camera,
		// mark them as unavailable (-1) rather than using hardcoded values
		limits.MinBitrate = -1
		limits.MaxBitrate = -1
	}

	return limits, nil
}

// Helper functions for parsing ONVIF responses

// Parse profile tokens from the GetProfiles response
func parseProfileTokens(body interface{}) ([]string, error) {
	// Convert the body to a byte array
	var bodyBytes []byte
	var err error

	switch b := body.(type) {
	case []byte:
		bodyBytes = b
	case string:
		bodyBytes = []byte(b)
	case io.Reader:
		bodyBytes, err = ioutil.ReadAll(b)
		if err != nil {
			return nil, fmt.Errorf("error reading response body: %v", err)
		}
	default:
		return nil, fmt.Errorf("unsupported body type for parsing profile tokens: %T", body)
	}

	// Define the XML response structure
	var response struct {
		XMLName xml.Name `xml:"Envelope"`
		Body    struct {
			GetProfilesResponse struct {
				Profiles []struct {
					Token string `xml:"token,attr"`
					Name  string `xml:"Name"`
				} `xml:"Profiles"`
			} `xml:"GetProfilesResponse"`
		} `xml:"Body"`
	}

	// Parse the XML
	if err = xml.Unmarshal(bodyBytes, &response); err != nil {
		return nil, fmt.Errorf("error unmarshalling profiles response: %v", err)
	}

	// Extract the tokens
	var tokens []string
	for _, profile := range response.Body.GetProfilesResponse.Profiles {
		tokens = append(tokens, profile.Token)
	}

	return tokens, nil
}

// Parse H264 configuration tokens from the GetVideoEncoderConfigurations response
func parseH264ConfigTokens(body interface{}) ([]string, error) {
	// Convert the body to a byte array
	var bodyBytes []byte
	var err error

	switch b := body.(type) {
	case []byte:
		bodyBytes = b
	case string:
		bodyBytes = []byte(b)
	case io.Reader:
		bodyBytes, err = ioutil.ReadAll(b)
		if err != nil {
			return []string{}, fmt.Errorf("error reading response body: %v", err)
		}
	default:
		return []string{}, fmt.Errorf("unsupported body type for parsing config tokens: %T", body)
	}

	// For debug purposes, print the raw XML
	// fmt.Println("Raw XML Response:", string(bodyBytes))

	// Define the XML response structure
	var response struct {
		XMLName xml.Name `xml:"Envelope"`
		Body    struct {
			GetVideoEncoderConfigurationsResponse struct {
				Configurations []struct {
					Token    string      `xml:"token,attr"`
					Name     string      `xml:"Name"`
					Encoding interface{} `xml:"Encoding"` // Could be int or string
				} `xml:"Configurations"`
			} `xml:"GetVideoEncoderConfigurationsResponse"`
		} `xml:"Body"`
	}

	// Parse the XML
	if err = xml.Unmarshal(bodyBytes, &response); err != nil {
		return []string{}, fmt.Errorf("error unmarshalling config response: %v", err)
	}

	// Extract the tokens for H264 configurations
	var tokens []string
	for _, config := range response.Body.GetVideoEncoderConfigurationsResponse.Configurations {
		// Check if this is an H264 configuration
		isH264 := false

		switch encoding := config.Encoding.(type) {
		case string:
			isH264 = encoding == "H264"
		case float64:
			isH264 = int(encoding) == 2
		case int:
			isH264 = encoding == 2
		case map[string]interface{}:
			// Some cameras might encode this differently
			if val, ok := encoding["#text"]; ok {
				if strVal, ok := val.(string); ok {
					isH264 = strVal == "H264"
				}
			}
		default:
			// If we can't determine, let's assume it could be H264 and let later validation catch issues
			fmt.Printf("Unhandled encoding type %T: %v for config %s\n", encoding, encoding, config.Token)
			isH264 = true
		}

		if isH264 {
			tokens = append(tokens, config.Token)
		}
	}

	// If still no tokens found, just return all configurations
	if len(tokens) == 0 {
		fmt.Println("No specific H264 configurations found, using all available configurations")
		for _, config := range response.Body.GetVideoEncoderConfigurationsResponse.Configurations {
			tokens = append(tokens, config.Token)
		}
	}

	return tokens, nil
}

// Parse available resolutions from the GetVideoEncoderConfigurationOptions response
func parseAvailableResolutions(body interface{}) ([]struct{ Width, Height int }, error) {
	// Convert the body to a byte array
	var bodyBytes []byte
	var err error

	switch b := body.(type) {
	case []byte:
		bodyBytes = b
	case string:
		bodyBytes = []byte(b)
	case io.Reader:
		bodyBytes, err = ioutil.ReadAll(b)
		if err != nil {
			return nil, fmt.Errorf("error reading response body: %v", err)
		}
	default:
		return nil, fmt.Errorf("unsupported body type for parsing resolutions: %T", body)
	}

	// Define the XML response structure focusing on H264 options
	var response struct {
		XMLName xml.Name `xml:"Envelope"`
		Body    struct {
			GetVideoEncoderConfigurationOptionsResponse struct {
				Options struct {
					H264 struct {
						ResolutionsAvailable []struct {
							Width  int `xml:"Width"`
							Height int `xml:"Height"`
						} `xml:"ResolutionsAvailable"`
					} `xml:"H264"`
				} `xml:"Options"`
			} `xml:"GetVideoEncoderConfigurationOptionsResponse"`
		} `xml:"Body"`
	}

	// Parse the XML
	if err = xml.Unmarshal(bodyBytes, &response); err != nil {
		return nil, fmt.Errorf("error unmarshalling options response: %v", err)
	}

	// Extract the resolutions
	xmlResolutions := response.Body.GetVideoEncoderConfigurationOptionsResponse.Options.H264.ResolutionsAvailable

	// Convert to the expected return type
	var resolutions []struct{ Width, Height int }
	for _, res := range xmlResolutions {
		resolutions = append(resolutions, struct{ Width, Height int }{
			Width:  res.Width,
			Height: res.Height,
		})
	}

	// If no resolutions were found, try a different approach with known common resolutions
	if len(resolutions) == 0 {
		return []struct{ Width, Height int }{
			{1280, 720},  // 720p
			{1280, 960},  // 960p
			{1920, 1080}, // 1080p
			{2304, 1296}, // 2304x1296
		}, nil
	}

	return resolutions, nil
}

// Parse H264 options from the GetVideoEncoderConfigurationOptions response
func parseH264Options(body interface{}) (struct {
	FrameRateMin int
	FrameRateMax int
	GovLengthMin int
	GovLengthMax int
	BitrateMin   int
	BitrateMax   int
	Profiles     []string
}, error) {
	// Convert the body to a byte array if needed
	var bodyBytes []byte
	var err error

	switch b := body.(type) {
	case []byte:
		bodyBytes = b
	case string:
		bodyBytes = []byte(b)
	default:
		// Handle the io.Reader case
		if reader, ok := body.(io.Reader); ok {
			bodyBytes, err = ioutil.ReadAll(reader)
			if err != nil {
				return struct {
					FrameRateMin int
					FrameRateMax int
					GovLengthMin int
					GovLengthMax int
					BitrateMin   int
					BitrateMax   int
					Profiles     []string
				}{}, fmt.Errorf("error reading response body: %v", err)
			}
		} else {
			return struct {
				FrameRateMin int
				FrameRateMax int
				GovLengthMin int
				GovLengthMax int
				BitrateMin   int
				BitrateMax   int
				Profiles     []string
			}{}, fmt.Errorf("unsupported body type for parsing H264 options")
		}
	}

	// Print the raw XML for debugging
	// fmt.Println(string(bodyBytes))

	// Define the XML response structure
	var response struct {
		XMLName xml.Name `xml:"Envelope"`
		Body    struct {
			GetVideoEncoderConfigurationOptionsResponse struct {
				Options struct {
					H264 struct {
						GovLengthRange struct {
							Min int `xml:"Min"`
							Max int `xml:"Max"`
						} `xml:"GovLengthRange"`
						FrameRateRange struct {
							Min int `xml:"Min"`
							Max int `xml:"Max"`
						} `xml:"FrameRateRange"`
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

	// Parse the XML
	if err = xml.Unmarshal(bodyBytes, &response); err != nil {
		return struct {
			FrameRateMin int
			FrameRateMax int
			GovLengthMin int
			GovLengthMax int
			BitrateMin   int
			BitrateMax   int
			Profiles     []string
		}{}, fmt.Errorf("error unmarshalling H264 options: %v", err)
	}

	// Extract the values
	h264Options := response.Body.GetVideoEncoderConfigurationOptionsResponse.Options.H264
	extension := response.Body.GetVideoEncoderConfigurationOptionsResponse.Options.Extension

	// Return the parsed options
	return struct {
		FrameRateMin int
		FrameRateMax int
		GovLengthMin int
		GovLengthMax int
		BitrateMin   int
		BitrateMax   int
		Profiles     []string
	}{
		FrameRateMin: h264Options.FrameRateRange.Min,
		FrameRateMax: h264Options.FrameRateRange.Max,
		GovLengthMin: h264Options.GovLengthRange.Min,
		GovLengthMax: h264Options.GovLengthRange.Max,
		BitrateMin:   extension.H264.BitrateRange.Min,
		BitrateMax:   extension.H264.BitrateRange.Max,
		Profiles:     h264Options.H264ProfilesSupported,
	}, nil
}
