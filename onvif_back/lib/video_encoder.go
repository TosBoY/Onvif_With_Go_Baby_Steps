package onvif_test

import (
	"encoding/xml"
	"fmt"
	"io"
	"strings"

	"github.com/use-go/onvif/media"
	mxsd "github.com/use-go/onvif/xsd"
	xsd "github.com/use-go/onvif/xsd/onvif"
)

// GetAllVideoEncoderConfigurations retrieves all video encoder configurations
func GetAllVideoEncoderConfigurations(c *Camera) ([]VideoEncoderConfig, error) {
	if c.Device == nil {
		return nil, fmt.Errorf("camera not connected")
	}

	// Create request
	req := media.GetVideoEncoderConfigurations{}

	// Call the method
	resp, err := c.Device.CallMethod(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Define response structure
	var configResp struct {
		Body struct {
			GetVideoEncoderConfigurationsResponse struct {
				Configurations []struct {
					Token      string `xml:"token,attr"`
					Name       string `xml:"Name"`
					UseCount   int    `xml:"UseCount"`
					Encoding   string `xml:"Encoding"`
					Resolution struct {
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
				} `xml:"Configurations"`
			} `xml:"GetVideoEncoderConfigurationsResponse"`
		} `xml:"Body"`
	}

	// Parse the response
	if err := xml.Unmarshal(body, &configResp); err != nil {
		return nil, err
	}

	// Convert to our type
	var configs []VideoEncoderConfig
	for _, c := range configResp.Body.GetVideoEncoderConfigurationsResponse.Configurations {
		configs = append(configs, VideoEncoderConfig{
			Token:       c.Token,
			Name:        c.Name,
			UseCount:    c.UseCount,
			Encoding:    c.Encoding,
			Width:       c.Resolution.Width,
			Height:      c.Resolution.Height,
			FrameRate:   c.RateControl.FrameRateLimit,
			BitRate:     c.RateControl.BitrateLimit,
			GovLength:   c.H264.GovLength,
			Quality:     c.Quality,
			H264Profile: c.H264.H264Profile,
		})
	}

	return configs, nil
}

// GetVideoEncoderConfiguration gets a specific video encoder configuration by token
func GetVideoEncoderConfiguration(c *Camera, token string) (*VideoEncoderConfig, error) {
	if c.Device == nil {
		return nil, fmt.Errorf("camera not connected")
	}

	// Create request
	req := media.GetVideoEncoderConfiguration{
		ConfigurationToken: xsd.ReferenceToken(token),
	}

	// Call the method
	resp, err := c.Device.CallMethod(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Define response structure
	var configResp struct {
		Body struct {
			GetVideoEncoderConfigurationResponse struct {
				Configuration struct {
					Token      string `xml:"token,attr"`
					Name       string `xml:"Name"`
					UseCount   int    `xml:"UseCount"`
					Encoding   string `xml:"Encoding"`
					Resolution struct {
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
				} `xml:"Configuration"`
			} `xml:"GetVideoEncoderConfigurationResponse"`
		} `xml:"Body"`
	}

	// Parse the response
	if err := xml.Unmarshal(body, &configResp); err != nil {
		return nil, err
	}

	// Get configuration from response
	config := configResp.Body.GetVideoEncoderConfigurationResponse.Configuration

	// Create our config type
	return &VideoEncoderConfig{
		Token:       config.Token,
		Name:        config.Name,
		UseCount:    config.UseCount,
		Encoding:    config.Encoding,
		Width:       config.Resolution.Width,
		Height:      config.Resolution.Height,
		FrameRate:   config.RateControl.FrameRateLimit,
		BitRate:     config.RateControl.BitrateLimit,
		GovLength:   config.H264.GovLength,
		Quality:     config.Quality,
		H264Profile: config.H264.H264Profile,
	}, nil
}

// GetVideoEncoderOptions gets available encoder options using both config token and profile token
func GetVideoEncoderOptions(c *Camera, configToken, profileToken string) (*VideoEncoderConfigurationOptionsResponse, error) {
	if c.Device == nil {
		return nil, fmt.Errorf("camera not connected")
	}

	// Create request with both tokens
	req := media.GetVideoEncoderConfigurationOptions{
		ConfigurationToken: xsd.ReferenceToken(configToken),
		ProfileToken:       xsd.ReferenceToken(profileToken),
	}

	// Call the method
	resp, err := c.Device.CallMethod(req)
	if err != nil {
		return nil, fmt.Errorf("error calling GetVideoEncoderConfigurationOptions: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	// Parse the response
	var optionsResp VideoEncoderConfigurationOptionsResponse
	if err := xml.Unmarshal(body, &optionsResp); err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %v", err)
	}

	return &optionsResp, nil
}

// GetH264Profiles retrieves supported H264 profiles using different tokens
func GetH264Profiles(c *Camera, configToken, profileToken string) ([]string, error) {
	if c.Device == nil {
		return nil, fmt.Errorf("camera not connected")
	}

	// Create request based on provided tokens
	req := media.GetVideoEncoderConfigurationOptions{}
	if configToken != "" {
		req.ConfigurationToken = xsd.ReferenceToken(configToken)
	}
	if profileToken != "" {
		req.ProfileToken = xsd.ReferenceToken(profileToken)
	}

	// Call the method
	resp, err := c.Device.CallMethod(req)
	if err != nil {
		return nil, fmt.Errorf("error calling GetVideoEncoderConfigurationOptions: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	// Parse the response
	var optionsResp VideoEncoderConfigurationOptionsResponse
	if err := xml.Unmarshal(body, &optionsResp); err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %v", err)
	}

	// Return supported H264 profiles
	return optionsResp.Body.GetVideoEncoderConfigurationOptionsResponse.Options.H264.H264ProfilesSupported, nil
}

// SetVideoEncoderConfiguration changes an encoder configuration
func SetVideoEncoderConfiguration(
	c *Camera,
	configToken string,
	configName string,
	width int,
	height int,
	frameRate int,
	bitRate int,
	govLength int,
	h264Profile string) error {

	if c.Device == nil {
		return fmt.Errorf("camera not connected")
	}

	// Create the configuration request
	setConfigRequest := media.SetVideoEncoderConfiguration{
		Configuration: xsd.VideoEncoderConfiguration{
			ConfigurationEntity: xsd.ConfigurationEntity{
				Token: xsd.ReferenceToken(configToken),
				Name:  xsd.Name(configName),
			},
			Encoding: "H264",
			Resolution: xsd.VideoResolution{
				Width:  mxsd.Int(width),
				Height: mxsd.Int(height),
			},
			Quality: 6.0, // Standard quality value
			RateControl: xsd.VideoRateControl{
				FrameRateLimit:   mxsd.Int(frameRate),
				EncodingInterval: mxsd.Int(1), // Standard encoding interval
				BitrateLimit:     mxsd.Int(bitRate),
			},
			H264: xsd.H264Configuration{
				GovLength:   mxsd.Int(govLength),
				H264Profile: xsd.H264Profile(h264Profile),
			},
			Multicast: xsd.MulticastConfiguration{
				Address: xsd.IPAddress{
					Type:        "IPv4",
					IPv4Address: "224.1.0.0",
				},
				Port:      0,
				TTL:       5,
				AutoStart: false,
			},
			SessionTimeout: "PT60S",
		},
		ForcePersistence: mxsd.Boolean(true),
	}

	// Call the method
	setConfigResp, err := c.Device.CallMethod(setConfigRequest)
	if err != nil {
		return fmt.Errorf("error setting video encoder configuration: %v", err)
	}
	defer setConfigResp.Body.Close()

	// Read response body to check if successful
	body, err := io.ReadAll(setConfigResp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %v", err)
	}

	// Check if the response has a fault element
	if ContainsFault(body) {
		return fmt.Errorf("server returned an error in response to SetVideoEncoderConfiguration")
	}

	return nil
}

// SetH264Profile attempts to change only the H264 profile
func SetH264Profile(c *Camera, configToken, profile string) (bool, error) {
	if c.Device == nil {
		return false, fmt.Errorf("camera not connected")
	}

	// First, get the current configuration
	config, err := GetVideoEncoderConfiguration(c, configToken)
	if err != nil {
		return false, fmt.Errorf("failed to get configuration: %v", err)
	}

	// Create the updated configuration with the new H264 profile
	setConfigReq := media.SetVideoEncoderConfiguration{
		Configuration: xsd.VideoEncoderConfiguration{
			ConfigurationEntity: xsd.ConfigurationEntity{
				Token:    xsd.ReferenceToken(configToken),
				Name:     xsd.Name(config.Name),
				UseCount: config.UseCount,
			},
			Encoding: xsd.VideoEncoding(config.Encoding),
			Resolution: xsd.VideoResolution{
				Width:  mxsd.Int(config.Width),
				Height: mxsd.Int(config.Height),
			},
			Quality: config.Quality,
			RateControl: xsd.VideoRateControl{
				FrameRateLimit:   mxsd.Int(config.FrameRate),
				EncodingInterval: mxsd.Int(1), // Standard value
				BitrateLimit:     mxsd.Int(config.BitRate),
			},
			H264: xsd.H264Configuration{
				GovLength:   mxsd.Int(config.GovLength),
				H264Profile: xsd.H264Profile(profile),
			},
		},
		ForcePersistence: mxsd.Boolean(true),
	}

	// Call the method to update the configuration
	setConfigResp, err := c.Device.CallMethod(setConfigReq)
	if err != nil {
		return false, nil // Consider failed request as "not supported"
	}
	defer setConfigResp.Body.Close()

	// Check if the update was successful
	body, err := io.ReadAll(setConfigResp.Body)
	if err != nil {
		return false, fmt.Errorf("failed to read response: %v", err)
	}

	// Check if the response contains a fault element
	if ContainsFault(body) {
		return false, nil // Profile not supported
	}

	return true, nil // Profile supported
}

// Add logging to ParseH264Options to debug empty resolutions
func ParseH264Options(options interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	// Check if we received any options
	if options == nil {
		fmt.Println("ParseH264Options: No options provided (nil)")
		return result
	}

	// Log raw options for debugging
	fmt.Printf("ParseH264Options: Raw options received: %+v\n", options)

	// Try to extract the available resolutions
	if videoOptions, ok := options.(map[string]interface{}); ok {
		// Extract available resolutions if present
		if resolutions, ok := videoOptions["ResolutionsAvailable"].([]interface{}); ok {
			result["ResolutionsAvailable"] = resolutions
			fmt.Printf("ParseH264Options: Found %d resolutions\n", len(resolutions))
		} else {
			fmt.Println("ParseH264Options: No ResolutionsAvailable found or invalid type")
			// For debugging: Try to find any resolution-like fields in the response
			for k, v := range videoOptions {
				if strings.Contains(strings.ToLower(k), "resolution") {
					fmt.Printf("ParseH264Options: Found resolution-like field: %s = %+v\n", k, v)
				}
			}
			result["ResolutionsAvailable"] = []interface{}{}
		}

		// Extract frame rate range if present
		if frameRateRange, ok := videoOptions["FrameRateRange"].(map[string]interface{}); ok {
			result["FrameRateRange"] = frameRateRange
			fmt.Printf("ParseH264Options: Found frame rate range: min=%v, max=%v\n",
				frameRateRange["Min"], frameRateRange["Max"])
		} else {
			fmt.Println("ParseH264Options: No FrameRateRange found or invalid type")
			// Check for other frame rate fields
			for k, v := range videoOptions {
				if strings.Contains(strings.ToLower(k), "frame") && strings.Contains(strings.ToLower(k), "rate") {
					fmt.Printf("ParseH264Options: Found frame rate field: %s = %+v\n", k, v)
				}
			}
		}

		// Extract frame rates directly from the response for compatibility
		frameRates := extractFrameRates(videoOptions)
		result["frameRates"] = frameRates

		// Extract available H.264 profiles if present
		if h264Profiles, ok := videoOptions["H264ProfilesSupported"].([]interface{}); ok {
			result["H264ProfilesSupported"] = h264Profiles
			fmt.Printf("ParseH264Options: Found %d H264 profiles\n", len(h264Profiles))
		} else {
			// Try alternative fields and format strings as needed
			h264Profiles := extractH264Profiles(videoOptions)
			result["H264ProfilesSupported"] = h264Profiles
		}

		// Extract encoding intervals
		if encodingIntervals, ok := videoOptions["EncodingIntervalRange"].(map[string]interface{}); ok {
			min := 1
			max := 30

			if minVal, ok := encodingIntervals["Min"].(float64); ok {
				min = int(minVal)
			}
			if maxVal, ok := encodingIntervals["Max"].(float64); ok {
				max = int(maxVal)
			}

			intervals := make([]int, 0)
			step := (max - min) / 6
			if step < 1 {
				step = 1
			}

			for i := min; i <= max; i += step {
				intervals = append(intervals, i)
			}
			if len(intervals) > 0 && intervals[len(intervals)-1] != max {
				intervals = append(intervals, max)
			}

			result["encodingIntervals"] = intervals
		} else {
			// Default encoding intervals
			result["encodingIntervals"] = []int{1, 5, 10, 15, 20, 25, 30}
		}

		// Include error message if resolutions are empty but other data exists
		if len(result["ResolutionsAvailable"].([]interface{})) == 0 && len(frameRates) > 0 {
			fmt.Println("ParseH264Options: WARNING - Got frame rates but no resolutions!")

			// Add common resolutions as fallback if the camera doesn't provide them
			fmt.Println("ParseH264Options: Adding fallback resolutions")
			fallbackResolutions := []map[string]interface{}{
				{"Width": 1920, "Height": 1080},
				{"Width": 1280, "Height": 720},
				{"Width": 640, "Height": 480},
				{"Width": 320, "Height": 240},
			}

			// Convert to generic interface slice
			resolutionsInterface := make([]interface{}, len(fallbackResolutions))
			for i, res := range fallbackResolutions {
				resolutionsInterface[i] = res
			}

			result["ResolutionsAvailable"] = resolutionsInterface
		}
	}

	// Log the final parsed result
	fmt.Printf("ParseH264Options: Final result: %+v\n", result)

	return result
}

// Helper to extract frame rates from various possible sources in the response
func extractFrameRates(options map[string]interface{}) []int {
	// Default values if we can't find anything
	defaultRates := []int{1, 5, 10, 15, 20, 25, 30}

	// Try to find frame rate information
	if frameRateRange, ok := options["FrameRateRange"].(map[string]interface{}); ok {
		min := 1
		max := 30

		if minVal, ok := frameRateRange["Min"].(float64); ok {
			min = int(minVal)
		}
		if maxVal, ok := frameRateRange["Max"].(float64); ok {
			max = int(maxVal)
		}

		// Generate a reasonable set of frame rates
		rates := make([]int, 0)
		step := (max - min) / 6
		if step < 1 {
			step = 1
		}

		for i := min; i <= max; i += step {
			rates = append(rates, i)
		}
		if len(rates) > 0 && rates[len(rates)-1] != max {
			rates = append(rates, max)
		}

		return rates
	}

	// Check for other possible frame rate fields
	for k, v := range options {
		if strings.Contains(strings.ToLower(k), "framerate") {
			fmt.Printf("extractFrameRates: Found field %s with value %+v\n", k, v)
			// Try to parse based on value type
			// Not implemented - would need to handle specific cases
		}
	}

	return defaultRates
}

// Helper to extract H264 profiles from various possible sources in the response
func extractH264Profiles(options map[string]interface{}) []string {
	// Default values if we can't find anything
	defaultProfiles := []string{"Baseline", "Main", "High"}

	// Try various field names that might contain profile information
	possibleFields := []string{
		"H264ProfilesSupported", "H264ProfileSupported",
		"ProfilesSupported", "ProfileSupported",
		"H264Profiles", "H264Profile",
	}

	for _, field := range possibleFields {
		if profiles, ok := options[field]; ok {
			fmt.Printf("extractH264Profiles: Found field %s with value %+v\n", field, profiles)

			// Try to convert to string slice based on type
			switch v := profiles.(type) {
			case []interface{}:
				strProfiles := make([]string, len(v))
				for i, p := range v {
					if str, ok := p.(string); ok {
						strProfiles[i] = str
					} else {
						// Convert to string if possible
						strProfiles[i] = fmt.Sprintf("%v", p)
					}
				}
				return strProfiles

			case []string:
				return v

			case string:
				// Single string value - might be comma separated or a single profile
				if strings.Contains(v, ",") {
					parts := strings.Split(v, ",")
					for i, p := range parts {
						parts[i] = strings.TrimSpace(p)
					}
					return parts
				}
				return []string{v}
			}
		}
	}

	return defaultProfiles
}
