package onvif_test

import (
	"encoding/xml"
	"fmt"
	"io"

	"github.com/use-go/onvif/media"
	xonvif "github.com/use-go/onvif/xsd/onvif"
)

// GetAllProfiles retrieves all media profiles from the camera
func GetAllProfiles(c *Camera) ([]Profile, error) {
	if c.Device == nil {
		return nil, fmt.Errorf("camera not connected")
	}

	// Create request
	req := media.GetProfiles{}

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
	var profileResp struct {
		Body struct {
			GetProfilesResponse struct {
				Profiles []struct {
					Token string `xml:"token,attr"`
					Name  string `xml:"Name"`
				} `xml:"Profiles"`
			} `xml:"GetProfilesResponse"`
		} `xml:"Body"`
	}

	// Parse the response
	if err := xml.Unmarshal(body, &profileResp); err != nil {
		return nil, err
	}

	// Convert to our type
	var profiles []Profile
	for _, p := range profileResp.Body.GetProfilesResponse.Profiles {
		profiles = append(profiles, Profile{
			Token: p.Token,
			Name:  p.Name,
		})
	}

	return profiles, nil
}

// GetProfileDetails retrieves detailed profile information including resolutions
func GetProfileDetails(c *Camera) (*ProfilesResponse, error) {
	if c.Device == nil {
		return nil, fmt.Errorf("camera not connected")
	}

	// Call GetProfiles
	resp, err := c.Device.CallMethod(media.GetProfiles{})
	if err != nil {
		return nil, fmt.Errorf("error calling GetProfiles: %v", err)
	}
	defer resp.Body.Close()

	// Read the XML response body
	rawXML, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	// Parse XML into custom struct
	var profilesResponse ProfilesResponse
	if err := xml.Unmarshal(rawXML, &profilesResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return &profilesResponse, nil
}

// GetCompatibleVideoEncoderConfigurations retrieves compatible video encoder configurations for a profile
func GetCompatibleVideoEncoderConfigurations(c *Camera, profileToken string) ([]VideoEncoderConfig, error) {
	if c.Device == nil {
		return nil, fmt.Errorf("camera not connected")
	}

	// Create request
	req := media.GetCompatibleVideoEncoderConfigurations{
		ProfileToken: xonvif.ReferenceToken(profileToken),
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

	// Parse the response
	var compatibleResp GetCompatibleVideoEncoderConfigurationsResponse
	if err := xml.Unmarshal(body, &compatibleResp); err != nil {
		return nil, err
	}

	return compatibleResp.Body.GetCompatibleVideoEncoderConfigurationsResponse.Configurations, nil
}

// GetCompatibleAudioDecoderConfigurations retrieves compatible audio decoder configurations for a profile
func GetCompatibleAudioDecoderConfigurations(c *Camera, profileToken string) ([]AudioDecoderConfig, error) {
	if c.Device == nil {
		return nil, fmt.Errorf("camera not connected")
	}

	// Create request
	req := media.GetCompatibleAudioDecoderConfigurations{
		ProfileToken: xonvif.ReferenceToken(profileToken),
	}

	// Call the method
	resp, err := c.Device.CallMethod(req)
	if err != nil {
		return nil, fmt.Errorf("error calling GetCompatibleAudioDecoderConfigurations: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	// Parse the response
	var configResp CompatibleAudioDecoderConfigsResponse
	if err := xml.Unmarshal(body, &configResp); err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %v", err)
	}

	// Convert to our type
	var configs []AudioDecoderConfig
	for _, c := range configResp.Body.GetCompatibleAudioDecoderConfigurationsResponse.Configurations {
		configs = append(configs, AudioDecoderConfig{
			Token:    c.Token,
			Name:     c.Name,
			UseCount: c.UseCount,
		})
	}

	return configs, nil
}

// AddVideoEncoderConfiguration adds a video encoder configuration to a profile
func AddVideoEncoderConfiguration(c *Camera, profileToken string, configToken string) error {
	if c.Device == nil {
		return fmt.Errorf("camera not connected")
	}

	// Create the request
	req := media.AddVideoEncoderConfiguration{
		ProfileToken:       xonvif.ReferenceToken(profileToken),
		ConfigurationToken: xonvif.ReferenceToken(configToken),
	}

	// Call the method
	resp, err := c.Device.CallMethod(req)
	if err != nil {
		return fmt.Errorf("error adding video encoder configuration to profile: %v", err)
	}
	defer resp.Body.Close()

	// Check response for fault
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %v", err)
	}

	// Check for fault in response
	if ContainsFault(body) {
		return fmt.Errorf("server returned an error in response to AddVideoEncoderConfiguration")
	}

	return nil
}

// AddAudioEncoderConfiguration adds an audio encoder configuration to a profile
func AddAudioEncoderConfiguration(c *Camera, profileToken string, configToken string) error {
	if c.Device == nil {
		return fmt.Errorf("camera not connected")
	}

	// Create the request
	req := media.AddAudioEncoderConfiguration{
		ProfileToken:       xonvif.ReferenceToken(profileToken),
		ConfigurationToken: xonvif.ReferenceToken(configToken),
	}

	// Call the method
	resp, err := c.Device.CallMethod(req)
	if err != nil {
		return fmt.Errorf("error adding audio encoder configuration to profile: %v", err)
	}
	defer resp.Body.Close()

	// Check response for fault
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %v", err)
	}

	// Check for fault in response
	if ContainsFault(body) {
		return fmt.Errorf("server returned an error in response to AddAudioEncoderConfiguration")
	}

	return nil
}

// GetActiveProfile returns the currently active profile and its video encoder configuration
func GetActiveProfile(c *Camera) (string, string, error) {
	profiles, err := GetAllProfiles(c)
	if err != nil {
		return "", "", fmt.Errorf("failed to get profiles: %v", err)
	}

	if len(profiles) == 0 {
		return "", "", fmt.Errorf("no profiles found")
	}

	// Get detailed profile information
	details, err := GetProfileDetails(c)
	if err != nil {
		return "", "", fmt.Errorf("failed to get profile details: %v", err)
	}

	// Log all available profiles and their resolutions
	fmt.Println("\nAvailable profiles:")
	var highestResProfile struct {
		Token  string
		Width  int
		Height int
	}

	for _, p := range details.Body.GetProfilesResponse.Profiles {
		if enc := p.VideoEncoderConfig; enc.Resolution.Width > 0 {
			fmt.Printf("Profile: %s\n", p.Token)
			fmt.Printf("- Resolution: %dx%d\n", enc.Resolution.Width, enc.Resolution.Height)

			// Track the profile with the highest resolution
			if enc.Resolution.Width*enc.Resolution.Height > highestResProfile.Width*highestResProfile.Height {
				highestResProfile.Token = p.Token
				highestResProfile.Width = enc.Resolution.Width
				highestResProfile.Height = enc.Resolution.Height
			}
		}
	}

	if highestResProfile.Token == "" {
		highestResProfile.Token = profiles[0].Token
		fmt.Printf("\nNo profile with resolution found, using first profile: %s\n", highestResProfile.Token)
	} else {
		fmt.Printf("\nSelected profile %s with resolution %dx%d\n",
			highestResProfile.Token, highestResProfile.Width, highestResProfile.Height)
	}

	// Get the associated config token
	activeConfig, err := GetProfileVideoEncoderConfig(c, highestResProfile.Token)
	if err != nil {
		return "", "", fmt.Errorf("failed to get profile's video encoder configuration: %v", err)
	}

	fmt.Printf("Using config token: %s\n", activeConfig)

	// Get available video encoder options for this configuration
	options, err := GetVideoEncoderOptions(c, activeConfig, highestResProfile.Token)
	if err != nil {
		fmt.Printf("Warning: Failed to get video encoder options: %v\n", err)
	} else {
		// Log available resolutions
		h264 := options.Body.GetVideoEncoderConfigurationOptionsResponse.Options.H264
		fmt.Printf("\nSupported resolutions for this profile/config:\n")
		for _, res := range h264.ResolutionsAvailable {
			fmt.Printf("- %dx%d\n", res.Width, res.Height)
		}
	}

	return highestResProfile.Token, activeConfig, nil
}

// GetProfileVideoEncoderConfig gets the video encoder config that's actually in use by a profile
func GetProfileVideoEncoderConfig(c *Camera, profileToken string) (string, error) {
	if c.Device == nil {
		return "", fmt.Errorf("camera not connected")
	}

	// Create request to get profile configuration
	req := media.GetProfile{
		ProfileToken: xonvif.ReferenceToken(profileToken),
	}

	// Call the method
	resp, err := c.Device.CallMethod(req)
	if err != nil {
		return "", fmt.Errorf("error getting profile: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %v", err)
	}

	// Define response structure to get the video encoder configuration token
	var getProfileResp struct {
		Body struct {
			GetProfileResponse struct {
				Profile struct {
					VideoEncoderConfiguration struct {
						Token string `xml:"token,attr"`
					} `xml:"VideoEncoderConfiguration"`
				} `xml:"Profile"`
			} `xml:"GetProfileResponse"`
		} `xml:"Body"`
	}

	// Parse the response
	if err := xml.Unmarshal(body, &getProfileResp); err != nil {
		return "", fmt.Errorf("error parsing response: %v", err)
	}

	// Get the video encoder configuration token
	configToken := getProfileResp.Body.GetProfileResponse.Profile.VideoEncoderConfiguration.Token
	if configToken == "" {
		return "", fmt.Errorf("no video encoder configuration found in profile")
	}

	return configToken, nil
}
