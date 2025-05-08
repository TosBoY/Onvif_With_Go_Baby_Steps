package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/use-go/onvif/media"
	"github.com/use-go/onvif/xsd"
)

// VideoResolution represents a video resolution with width and height
type VideoResolution struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

// H264Options represents the available H264 encoding options
type H264Options struct {
	ResolutionOptions  []VideoResolution `json:"resolutions"`
	FrameRateOptions   []int             `json:"frameRates"`
	EncodingIntervals  []int             `json:"encodingIntervals"`
	BitrateOptions     []int             `json:"bitrates"`
	H264ProfileOptions []string          `json:"h264Profiles"`
}

// VideoEncoderConfiguration represents a video encoder configuration
type VideoEncoderConfiguration struct {
	Name        string         `json:"name"`
	Token       string         `json:"token"`
	Encoding    string         `json:"encoding"`
	Width       int            `json:"width"`
	Height      int            `json:"height"`
	Quality     float64        `json:"quality"`
	FrameRate   int            `json:"frameRate"`
	BitRate     int            `json:"bitRate"`
	GovLength   int            `json:"govLength"`
	H264Profile string         `json:"h264Profile,omitempty"`
	Multicast   *MulticastInfo `json:"multicast,omitempty"`
}

// MulticastInfo represents multicast configuration information
type MulticastInfo struct {
	Address   string `json:"address"`
	Port      int    `json:"port"`
	TTL       int    `json:"ttl"`
	AutoStart bool   `json:"autoStart"`
}

// Profile represents a media profile
type Profile struct {
	Name               string                     `json:"name"`
	Token              string                     `json:"token"`
	VideoSourceConfig  *VideoSourceConfiguration  `json:"videoSourceConfig,omitempty"`
	VideoEncoderConfig *VideoEncoderConfiguration `json:"videoEncoderConfig,omitempty"`
}

// VideoSourceConfiguration represents a video source configuration
type VideoSourceConfiguration struct {
	Name        string `json:"name"`
	Token       string `json:"token"`
	SourceToken string `json:"sourceToken"`
}

// GetAllProfiles retrieves all media profiles from the camera
func GetAllProfiles(camera *Camera) ([]Profile, error) {
	request := media.GetProfiles{}
	response, err := camera.Device.CallMethod(request)
	if err != nil {
		return nil, fmt.Errorf("failed to get profiles: %v", err)
	}
	defer response.Body.Close()

	// Read the response body
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// Parse the response
	var envelope struct {
		Body struct {
			GetProfilesResponse struct {
				Profiles []struct {
					Name              xsd.String `xml:"Name"`
					Token             xsd.String `xml:"token,attr"`
					VideoSourceConfig struct {
						Name        xsd.String `xml:"Name"`
						Token       xsd.String `xml:"token,attr"`
						SourceToken xsd.String `xml:"SourceToken"`
					} `xml:"VideoSourceConfiguration"`
					VideoEncoderConfig struct {
						Name       xsd.String `xml:"Name"`
						Token      xsd.String `xml:"token,attr"`
						Encoding   xsd.String `xml:"Encoding"`
						Resolution struct {
							Width  xsd.Int `xml:"Width"`
							Height xsd.Int `xml:"Height"`
						} `xml:"Resolution"`
						Quality     xsd.Float `xml:"Quality"`
						RateControl struct {
							FrameRateLimit xsd.Int `xml:"FrameRateLimit"`
							BitrateLimit   xsd.Int `xml:"BitrateLimit"`
						} `xml:"RateControl"`
						H264 struct {
							GovLength   xsd.Int    `xml:"GovLength"`
							H264Profile xsd.String `xml:"H264Profile"`
						} `xml:"H264"`
						Multicast struct {
							Address struct {
								Type        xsd.String `xml:"Type"`
								IPv4Address xsd.String `xml:"IPv4Address"`
							} `xml:"Address"`
							Port      xsd.Int     `xml:"Port"`
							TTL       xsd.Int     `xml:"TTL"`
							AutoStart xsd.Boolean `xml:"AutoStart"`
						} `xml:"Multicast"`
					} `xml:"VideoEncoderConfiguration"`
				} `xml:"Profiles"`
			} `xml:"GetProfilesResponse"`
		} `xml:"Body"`
	}

	if err := xml.Unmarshal(responseBody, &envelope); err != nil {
		return nil, fmt.Errorf("failed to parse profiles XML: %v", err)
	}

	// Convert to our Profile struct
	profiles := make([]Profile, 0)
	for _, p := range envelope.Body.GetProfilesResponse.Profiles {
		profile := Profile{
			Name:  string(p.Name),
			Token: string(p.Token),
		}

		// Add video source config if available
		if p.VideoSourceConfig.Token != "" {
			profile.VideoSourceConfig = &VideoSourceConfiguration{
				Name:        string(p.VideoSourceConfig.Name),
				Token:       string(p.VideoSourceConfig.Token),
				SourceToken: string(p.VideoSourceConfig.SourceToken),
			}
		}

		// Add video encoder config if available
		if p.VideoEncoderConfig.Token != "" {
			profile.VideoEncoderConfig = &VideoEncoderConfiguration{
				Name:        string(p.VideoEncoderConfig.Name),
				Token:       string(p.VideoEncoderConfig.Token),
				Encoding:    string(p.VideoEncoderConfig.Encoding),
				Width:       int(p.VideoEncoderConfig.Resolution.Width),
				Height:      int(p.VideoEncoderConfig.Resolution.Height),
				Quality:     float64(p.VideoEncoderConfig.Quality),
				FrameRate:   int(p.VideoEncoderConfig.RateControl.FrameRateLimit),
				BitRate:     int(p.VideoEncoderConfig.RateControl.BitrateLimit),
				GovLength:   int(p.VideoEncoderConfig.H264.GovLength),
				H264Profile: string(p.VideoEncoderConfig.H264.H264Profile),
			}

			// Add multicast info if available
			if p.VideoEncoderConfig.Multicast.Address.IPv4Address != "" {
				profile.VideoEncoderConfig.Multicast = &MulticastInfo{
					Address:   string(p.VideoEncoderConfig.Multicast.Address.IPv4Address),
					Port:      int(p.VideoEncoderConfig.Multicast.Port),
					TTL:       int(p.VideoEncoderConfig.Multicast.TTL),
					AutoStart: bool(p.VideoEncoderConfig.Multicast.AutoStart),
				}
			}
		}

		profiles = append(profiles, profile)
	}

	return profiles, nil
}

// GetAllVideoEncoderConfigurations retrieves all video encoder configurations from the camera
func GetAllVideoEncoderConfigurations(camera *Camera) ([]VideoEncoderConfiguration, error) {
	request := media.GetVideoEncoderConfigurations{}
	response, err := camera.Device.CallMethod(request)
	if err != nil {
		return nil, fmt.Errorf("failed to get video encoder configurations: %v", err)
	}
	defer response.Body.Close()

	// Read the response body
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var envelope struct {
		Body struct {
			GetVideoEncoderConfigurationsResponse struct {
				Configurations []struct {
					Name       xsd.String `xml:"Name"`
					Token      xsd.String `xml:"token,attr"`
					Encoding   xsd.String `xml:"Encoding"`
					Resolution struct {
						Width  xsd.Int `xml:"Width"`
						Height xsd.Int `xml:"Height"`
					} `xml:"Resolution"`
					Quality     xsd.Float `xml:"Quality"`
					RateControl struct {
						FrameRateLimit xsd.Int `xml:"FrameRateLimit"`
						BitrateLimit   xsd.Int `xml:"BitrateLimit"`
					} `xml:"RateControl"`
					H264 struct {
						GovLength   xsd.Int    `xml:"GovLength"`
						H264Profile xsd.String `xml:"H264Profile"`
					} `xml:"H264"`
					Multicast struct {
						Address struct {
							Type        xsd.String `xml:"Type"`
							IPv4Address xsd.String `xml:"IPv4Address"`
						} `xml:"Address"`
						Port      xsd.Int     `xml:"Port"`
						TTL       xsd.Int     `xml:"TTL"`
						AutoStart xsd.Boolean `xml:"AutoStart"`
					} `xml:"Multicast"`
				} `xml:"Configurations"`
			} `xml:"GetVideoEncoderConfigurationsResponse"`
		} `xml:"Body"`
	}

	if err := xml.Unmarshal(responseBody, &envelope); err != nil {
		return nil, fmt.Errorf("failed to parse video encoder configurations XML: %v", err)
	}

	configs := make([]VideoEncoderConfiguration, 0)
	for _, c := range envelope.Body.GetVideoEncoderConfigurationsResponse.Configurations {
		config := VideoEncoderConfiguration{
			Name:        string(c.Name),
			Token:       string(c.Token),
			Encoding:    string(c.Encoding),
			Width:       int(c.Resolution.Width),
			Height:      int(c.Resolution.Height),
			Quality:     float64(c.Quality),
			FrameRate:   int(c.RateControl.FrameRateLimit),
			BitRate:     int(c.RateControl.BitrateLimit),
			GovLength:   int(c.H264.GovLength),
			H264Profile: string(c.H264.H264Profile),
		}

		// Add multicast info if available
		if c.Multicast.Address.IPv4Address != "" {
			config.Multicast = &MulticastInfo{
				Address:   string(c.Multicast.Address.IPv4Address),
				Port:      int(c.Multicast.Port),
				TTL:       int(c.Multicast.TTL),
				AutoStart: bool(c.Multicast.AutoStart),
			}
		}

		configs = append(configs, config)
	}

	return configs, nil
}

// GetVideoEncoderConfiguration retrieves a specific video encoder configuration by token
func GetVideoEncoderConfiguration(camera *Camera, token string) (*VideoEncoderConfiguration, error) {
	// For v0.0.9, we need to create a custom SOAP request for GetVideoEncoderConfiguration
	soapEnvelope := `<?xml version="1.0" encoding="UTF-8"?>
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope" 
                   xmlns:trt="http://www.onvif.org/ver10/media/wsdl">
  <SOAP-ENV:Body>
    <trt:GetVideoEncoderConfiguration>
      <trt:ConfigurationToken>` + token + `</trt:ConfigurationToken>
    </trt:GetVideoEncoderConfiguration>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

	// Send the SOAP request as a custom HTTP request
	response, err := camera.sendCustomSOAPRequest(soapEnvelope)
	if err != nil {
		return nil, fmt.Errorf("failed to get video encoder configuration: %v", err)
	}
	defer response.Body.Close()

	// Read the response body
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var envelope struct {
		Body struct {
			GetVideoEncoderConfigurationResponse struct {
				Configuration struct {
					Name       xsd.String `xml:"Name"`
					Token      xsd.String `xml:"token,attr"`
					Encoding   xsd.String `xml:"Encoding"`
					Resolution struct {
						Width  xsd.Int `xml:"Width"`
						Height xsd.Int `xml:"Height"`
					} `xml:"Resolution"`
					Quality     xsd.Float `xml:"Quality"`
					RateControl struct {
						FrameRateLimit xsd.Int `xml:"FrameRateLimit"`
						BitrateLimit   xsd.Int `xml:"BitrateLimit"`
					} `xml:"RateControl"`
					H264 struct {
						GovLength   xsd.Int    `xml:"GovLength"`
						H264Profile xsd.String `xml:"H264Profile"`
					} `xml:"H264"`
					Multicast struct {
						Address struct {
							Type        xsd.String `xml:"Type"`
							IPv4Address xsd.String `xml:"IPv4Address"`
						} `xml:"Address"`
						Port      xsd.Int     `xml:"Port"`
						TTL       xsd.Int     `xml:"TTL"`
						AutoStart xsd.Boolean `xml:"AutoStart"`
					} `xml:"Multicast"`
				} `xml:"Configuration"`
			} `xml:"GetVideoEncoderConfigurationResponse"`
		} `xml:"Body"`
	}

	if err := xml.Unmarshal(responseBody, &envelope); err != nil {
		return nil, fmt.Errorf("failed to parse video encoder configuration XML: %v", err)
	}

	c := envelope.Body.GetVideoEncoderConfigurationResponse.Configuration
	config := &VideoEncoderConfiguration{
		Name:        string(c.Name),
		Token:       string(c.Token),
		Encoding:    string(c.Encoding),
		Width:       int(c.Resolution.Width),
		Height:      int(c.Resolution.Height),
		Quality:     float64(c.Quality),
		FrameRate:   int(c.RateControl.FrameRateLimit),
		BitRate:     int(c.RateControl.BitrateLimit),
		GovLength:   int(c.H264.GovLength),
		H264Profile: string(c.H264.H264Profile),
	}

	// Add multicast info if available
	if c.Multicast.Address.IPv4Address != "" {
		config.Multicast = &MulticastInfo{
			Address:   string(c.Multicast.Address.IPv4Address),
			Port:      int(c.Multicast.Port),
			TTL:       int(c.Multicast.TTL),
			AutoStart: bool(c.Multicast.AutoStart),
		}
	}

	return config, nil
}

// GetVideoEncoderOptions retrieves the video encoder options for a specific configuration
func GetVideoEncoderOptions(camera *Camera, configToken string, profileToken string) (interface{}, error) {
	// For v0.0.9, we need to create a custom SOAP request for GetVideoEncoderConfigurationOptions
	soapEnvelope := `<?xml version="1.0" encoding="UTF-8"?>
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope" 
                   xmlns:trt="http://www.onvif.org/ver10/media/wsdl">
  <SOAP-ENV:Body>
    <trt:GetVideoEncoderConfigurationOptions>`

	if configToken != "" {
		soapEnvelope += `
      <trt:ConfigurationToken>` + configToken + `</trt:ConfigurationToken>`
	}

	if profileToken != "" {
		soapEnvelope += `
      <trt:ProfileToken>` + profileToken + `</trt:ProfileToken>`
	}

	soapEnvelope += `
    </trt:GetVideoEncoderConfigurationOptions>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

	// Send the SOAP request
	response, err := camera.sendCustomSOAPRequest(soapEnvelope)
	if err != nil {
		return nil, fmt.Errorf("failed to get video encoder options: %v", err)
	}

	// Read the body content for parsing
	defer response.Body.Close()
	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	return string(bodyBytes), nil
}

// ParseH264Options parses the raw response body to extract H.264 encoding options
func ParseH264Options(responseBody interface{}) *H264Options {
	// Convert response body to string for parsing
	bodyStr := fmt.Sprintf("%v", responseBody)

	options := &H264Options{
		ResolutionOptions:  make([]VideoResolution, 0),
		FrameRateOptions:   make([]int, 0),
		EncodingIntervals:  make([]int, 0),
		BitrateOptions:     make([]int, 0),
		H264ProfileOptions: make([]string, 0),
	}

	// Parse width and height options (simplistic approach)
	widthMatches := strings.Split(bodyStr, "Width:")
	heightMatches := strings.Split(bodyStr, "Height:")

	// Extract resolutions
	// Note: This is a naive implementation and might need adaptation
	// based on the actual XML structure
	if len(widthMatches) > 1 && len(heightMatches) > 1 {
		for i := 1; i < len(widthMatches) && i < len(heightMatches); i++ {
			wStr := widthMatches[i]
			hStr := heightMatches[i]

			var width, height int
			fmt.Sscanf(wStr, "%d", &width)
			fmt.Sscanf(hStr, "%d", &height)

			if width > 0 && height > 0 {
				options.ResolutionOptions = append(options.ResolutionOptions, VideoResolution{
					Width:  width,
					Height: height,
				})
			}
		}
	}

	// Extract frame rate range
	if strings.Contains(bodyStr, "FrameRateRange") {
		// This is simplified - actual implementation would parse Min/Max values
		options.FrameRateOptions = []int{1, 5, 10, 15, 20, 25, 30}
	}

	// Extract encoding interval range
	if strings.Contains(bodyStr, "EncodingIntervalRange") {
		// This is simplified - actual implementation would parse Min/Max values
		options.EncodingIntervals = []int{1, 5, 10, 15, 20, 25, 30}
	}

	// Extract bitrate options
	if strings.Contains(bodyStr, "BitrateRange") {
		// This is simplified - actual implementation would parse Min/Max values
		options.BitrateOptions = []int{
			64000, 128000, 256000, 512000, 1000000, 2000000, 4000000, 8000000,
		}
	}

	// Extract H264 profile options
	if strings.Contains(bodyStr, "H264ProfilesSupported") {
		if strings.Contains(bodyStr, "Baseline") {
			options.H264ProfileOptions = append(options.H264ProfileOptions, "Baseline")
		}
		if strings.Contains(bodyStr, "Main") {
			options.H264ProfileOptions = append(options.H264ProfileOptions, "Main")
		}
		if strings.Contains(bodyStr, "High") {
			options.H264ProfileOptions = append(options.H264ProfileOptions, "High")
		}
		if strings.Contains(bodyStr, "Extended") {
			options.H264ProfileOptions = append(options.H264ProfileOptions, "Extended")
		}
	}

	return options
}

// SetVideoEncoderConfiguration sets a video encoder configuration
func SetVideoEncoderConfiguration(
	camera *Camera,
	token, name string,
	width, height int,
	frameRate, bitRate, govLength int,
	h264Profile string) error {

	// First get the current configuration to preserve settings we're not changing
	config, err := GetVideoEncoderConfiguration(camera, token)
	if err != nil {
		return fmt.Errorf("failed to get current configuration: %v", err)
	}

	// For v0.0.9, we need to use a custom struct that matches the expected XML format
	// Create a custom SOAP request for SetVideoEncoderConfiguration
	soapEnvelope := `<?xml version="1.0" encoding="UTF-8"?>
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope" 
                   xmlns:trt="http://www.onvif.org/ver10/media/wsdl"
                   xmlns:tt="http://www.onvif.org/ver10/schema">
  <SOAP-ENV:Body>
    <trt:SetVideoEncoderConfiguration>
      <trt:Configuration token="` + token + `">
        <tt:Name>` + name + `</tt:Name>
        <tt:UseCount>1</tt:UseCount>
        <tt:Encoding>H264</tt:Encoding>
        <tt:Resolution>
          <tt:Width>` + fmt.Sprintf("%d", width) + `</tt:Width>
          <tt:Height>` + fmt.Sprintf("%d", height) + `</tt:Height>
        </tt:Resolution>
        <tt:Quality>` + fmt.Sprintf("%.1f", config.Quality) + `</tt:Quality>
        <tt:RateControl>
          <tt:FrameRateLimit>` + fmt.Sprintf("%d", frameRate) + `</tt:FrameRateLimit>
          <tt:EncodingInterval>` + fmt.Sprintf("%d", govLength) + `</tt:EncodingInterval>
          <tt:BitrateLimit>` + fmt.Sprintf("%d", bitRate) + `</tt:BitrateLimit>
        </tt:RateControl>
        <tt:H264>
          <tt:GovLength>` + fmt.Sprintf("%d", govLength) + `</tt:GovLength>
          <tt:H264Profile>` + h264Profile + `</tt:H264Profile>
        </tt:H264>
        <tt:SessionTimeout>PT10S</tt:SessionTimeout>
      </trt:Configuration>
      <trt:ForcePersistence>true</trt:ForcePersistence>
    </trt:SetVideoEncoderConfiguration>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

	// Send the raw SOAP request
	resp, err := camera.sendCustomSOAPRequest(soapEnvelope)
	if err != nil {
		return fmt.Errorf("failed to set video encoder configuration: %v", err)
	}
	defer resp.Body.Close()

	return nil
}

// GetStreamURI retrieves the stream URI for a specific profile
func GetStreamURI(camera *Camera, profileToken string) (string, error) {
	// For v0.0.9, we need to create a custom SOAP request for GetStreamUri
	soapEnvelope := `<?xml version="1.0" encoding="UTF-8"?>
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope" 
                   xmlns:trt="http://www.onvif.org/ver10/media/wsdl"
                   xmlns:tt="http://www.onvif.org/ver10/schema">
  <SOAP-ENV:Body>
    <trt:GetStreamUri>
      <trt:StreamSetup>
        <tt:Stream>RTP-Unicast</tt:Stream>
        <tt:Transport>
          <tt:Protocol>RTSP</tt:Protocol>
        </tt:Transport>
      </trt:StreamSetup>
      <trt:ProfileToken>` + profileToken + `</trt:ProfileToken>
    </trt:GetStreamUri>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

	// Send the SOAP request
	response, err := camera.sendCustomSOAPRequest(soapEnvelope)
	if err != nil {
		return "", fmt.Errorf("failed to get stream URI: %v", err)
	}
	defer response.Body.Close()

	// Read the response body
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	var envelope struct {
		Body struct {
			GetStreamUriResponse struct {
				MediaUri struct {
					Uri                 string `xml:"Uri"`
					InvalidAfterConnect bool   `xml:"InvalidAfterConnect"`
					InvalidAfterReboot  bool   `xml:"InvalidAfterReboot"`
					Timeout             string `xml:"Timeout"`
				} `xml:"MediaUri"`
			} `xml:"GetStreamUriResponse"`
		} `xml:"Body"`
	}

	if err := xml.Unmarshal(responseBody, &envelope); err != nil {
		return "", fmt.Errorf("failed to parse stream URI XML: %v", err)
	}

	return envelope.Body.GetStreamUriResponse.MediaUri.Uri, nil
}

// Helper method to send custom SOAP requests
func (camera *Camera) sendCustomSOAPRequest(soapEnvelope string) (*http.Response, error) {
	// Build the SOAP request URL
	endpoint := fmt.Sprintf("http://%s:%d/onvif/media", camera.IP, camera.Port)

	// Create a new HTTP request
	req, err := http.NewRequest("POST", endpoint, bytes.NewBufferString(soapEnvelope))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %v", err)
	}

	// Set headers
	req.SetBasicAuth(camera.Username, camera.Password)
	req.Header.Set("Content-Type", "application/soap+xml; charset=utf-8")
	req.Header.Set("SOAPAction", "http://www.onvif.org/ver10/media/wsdl")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %v", err)
	}

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP error %d: %s", resp.StatusCode, body)
	}

	return resp, nil
}
