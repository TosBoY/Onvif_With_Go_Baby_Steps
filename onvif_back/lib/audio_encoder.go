package onvif_test

import (
	"encoding/xml"
	"fmt"
	"io"

	"github.com/use-go/onvif/media"
	mxsd "github.com/use-go/onvif/xsd"
	xsd "github.com/use-go/onvif/xsd/onvif"
)

// GetAllAudioEncoderConfigurations retrieves all audio encoder configurations
func GetAllAudioEncoderConfigurations(c *Camera) ([]AudioEncoderConfig, error) {
	if c.Device == nil {
		return nil, fmt.Errorf("camera not connected")
	}

	// Create request
	req := media.GetAudioEncoderConfigurations{}

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
			GetAudioEncoderConfigurationsResponse struct {
				Configurations []struct {
					Token      string `xml:"token,attr"`
					Name       string `xml:"Name"`
					UseCount   int    `xml:"UseCount"`
					Encoding   string `xml:"Encoding"`
					Bitrate    int    `xml:"Bitrate"`
					SampleRate int    `xml:"SampleRate"`
				} `xml:"Configurations"`
			} `xml:"GetAudioEncoderConfigurationsResponse"`
		} `xml:"Body"`
	}

	// Parse the response
	if err := xml.Unmarshal(body, &configResp); err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %v", err)
	}

	// Convert to our type
	var configs []AudioEncoderConfig
	for _, c := range configResp.Body.GetAudioEncoderConfigurationsResponse.Configurations {
		configs = append(configs, AudioEncoderConfig{
			Token:      c.Token,
			Name:       c.Name,
			UseCount:   c.UseCount,
			Encoding:   c.Encoding,
			Bitrate:    c.Bitrate,
			SampleRate: c.SampleRate,
		})
	}

	return configs, nil
}

// GetAudioEncoderConfiguration gets a specific audio encoder configuration by token
func GetAudioEncoderConfiguration(c *Camera, token string) (*AudioEncoderConfig, error) {
	if c.Device == nil {
		return nil, fmt.Errorf("camera not connected")
	}

	// Create request
	req := media.GetAudioEncoderConfiguration{
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
			GetAudioEncoderConfigurationResponse struct {
				Configuration struct {
					Token      string `xml:"token,attr"`
					Name       string `xml:"Name"`
					UseCount   int    `xml:"UseCount"`
					Encoding   string `xml:"Encoding"`
					Bitrate    int    `xml:"Bitrate"`
					SampleRate int    `xml:"SampleRate"`
				} `xml:"Configuration"`
			} `xml:"GetAudioEncoderConfigurationResponse"`
		} `xml:"Body"`
	}

	// Parse the response
	if err := xml.Unmarshal(body, &configResp); err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %v", err)
	}

	// Get configuration from response
	config := configResp.Body.GetAudioEncoderConfigurationResponse.Configuration

	// Create our config type
	return &AudioEncoderConfig{
		Token:      config.Token,
		Name:       config.Name,
		UseCount:   config.UseCount,
		Encoding:   config.Encoding,
		Bitrate:    config.Bitrate,
		SampleRate: config.SampleRate,
	}, nil
}

// GetAudioEncoderOptions gets available encoder options using both config token and profile token
func GetAudioEncoderOptions(c *Camera, configToken, profileToken string) (*AudioEncoderConfigurationOptionsResponse, error) {
	if c.Device == nil {
		return nil, fmt.Errorf("camera not connected")
	}

	// Create request with both tokens
	req := media.GetAudioEncoderConfigurationOptions{
		ConfigurationToken: xsd.ReferenceToken(configToken),
		ProfileToken:       xsd.ReferenceToken(profileToken),
	}

	// Call the method
	resp, err := c.Device.CallMethod(req)
	if err != nil {
		return nil, fmt.Errorf("error calling GetAudioEncoderConfigurationOptions: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	// Parse the response
	var optionsResp AudioEncoderConfigurationOptionsResponse
	if err := xml.Unmarshal(body, &optionsResp); err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %v", err)
	}

	return &optionsResp, nil
}

// SetAudioEncoderConfiguration changes an audio encoder configuration
func SetAudioEncoderConfiguration(
	c *Camera,
	configToken string,
	configName string,
	encoding string,
	bitrate int,
	sampleRate int) error {

	if c.Device == nil {
		return fmt.Errorf("camera not connected")
	}

	// Create the configuration request
	setConfigRequest := media.SetAudioEncoderConfiguration{
		Configuration: xsd.AudioEncoderConfiguration{
			ConfigurationEntity: xsd.ConfigurationEntity{
				Token: xsd.ReferenceToken(configToken),
				Name:  xsd.Name(configName),
			},
			Encoding:   xsd.AudioEncoding(encoding),
			Bitrate:    bitrate,
			SampleRate: sampleRate,
			Multicast: xsd.MulticastConfiguration{
				Address: xsd.IPAddress{
					Type:        "IPv4",
					IPv4Address: "224.1.0.0",
				},
				Port:      40002,
				TTL:       64,
				AutoStart: false,
			},
			SessionTimeout: "PT60S",
		},
		ForcePersistence: mxsd.Boolean(true),
	}

	// Call the method
	setConfigResp, err := c.Device.CallMethod(setConfigRequest)
	if err != nil {
		return fmt.Errorf("error setting audio encoder configuration: %v", err)
	}
	defer setConfigResp.Body.Close()

	// Read response body to check if successful
	body, err := io.ReadAll(setConfigResp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %v", err)
	}

	// Check if the response has a fault element
	if ContainsFault(body) {
		return fmt.Errorf("server returned an error in response to SetAudioEncoderConfiguration")
	}

	return nil
}

// GetAllAudioDecoderConfigurations retrieves all audio decoder configurations
func GetAllAudioDecoderConfigurations(c *Camera) ([]AudioDecoderConfig, error) {
	if c.Device == nil {
		return nil, fmt.Errorf("camera not connected")
	}

	// Create request
	req := media.GetAudioDecoderConfigurations{}

	// Call the method
	resp, err := c.Device.CallMethod(req)
	if err != nil {
		return nil, fmt.Errorf("error calling GetAudioDecoderConfigurations: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	// Define response structure
	var configResp struct {
		Body struct {
			GetAudioDecoderConfigurationsResponse struct {
				Configurations []struct {
					Token    string `xml:"token,attr"`
					Name     string `xml:"Name"`
					UseCount int    `xml:"UseCount"`
				} `xml:"Configurations"`
			} `xml:"GetAudioDecoderConfigurationsResponse"`
		} `xml:"Body"`
	}

	// Parse the response
	if err := xml.Unmarshal(body, &configResp); err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %v", err)
	}

	// Convert to our type
	var configs []AudioDecoderConfig
	for _, c := range configResp.Body.GetAudioDecoderConfigurationsResponse.Configurations {
		configs = append(configs, AudioDecoderConfig{
			Token:    c.Token,
			Name:     c.Name,
			UseCount: c.UseCount,
		})
	}

	return configs, nil
}

// GetAudioDecoderOptions gets options for a specific audio decoder configuration
func GetAudioDecoderOptions(c *Camera, configToken, profileToken string) (*AudioDecoderConfigOptionsResponse, error) {
	if c.Device == nil {
		return nil, fmt.Errorf("camera not connected")
	}

	// Create request
	req := media.GetAudioDecoderConfigurationOptions{}

	// Add optional tokens if provided
	if configToken != "" {
		req.ConfigurationToken = xsd.ReferenceToken(configToken)
	}
	if profileToken != "" {
		req.ProfileToken = xsd.ReferenceToken(profileToken)
	}

	// Call the method
	resp, err := c.Device.CallMethod(req)
	if err != nil {
		return nil, fmt.Errorf("error calling GetAudioDecoderConfigurationOptions: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	// Parse the response
	var optionsResp AudioDecoderConfigOptionsResponse
	if err := xml.Unmarshal(body, &optionsResp); err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %v", err)
	}

	return &optionsResp, nil
}
