package onvif_test

import (
	"encoding/xml"
	"fmt"
	"io"

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
