package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"regexp"
	"strconv"
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
	fmt.Println("ParseH264Options: Processing XML of length:", len(bodyStr))

	options := &H264Options{
		ResolutionOptions:  make([]VideoResolution, 0),
		FrameRateOptions:   make([]int, 0),
		EncodingIntervals:  make([]int, 0),
		BitrateOptions:     make([]int, 0),
		H264ProfileOptions: make([]string, 0),
	}

	// First try to parse the XML using proper XML parsing
	var envelope struct {
		XMLName xml.Name `xml:"Envelope"`
		Body    struct {
			GetVideoEncoderConfigurationOptionsResponse struct {
				Options struct {
					H264 struct {
						ResolutionsAvailable []struct {
							Width  int `xml:"Width"`
							Height int `xml:"Height"`
						} `xml:"ResolutionsAvailable"`
						FrameRateRange struct {
							Min int `xml:"Min"`
							Max int `xml:"Max"`
						} `xml:"FrameRateRange"`
						EncodingIntervalRange struct {
							Min int `xml:"Min"`
							Max int `xml:"Max"`
						} `xml:"EncodingIntervalRange"`
						BitrateRange struct {
							Min int `xml:"Min"`
							Max int `xml:"Max"`
						} `xml:"BitrateRange"`
						GovLengthRange struct {
							Min int `xml:"Min"`
							Max int `xml:"Max"`
						} `xml:"GovLengthRange"`
						H264ProfilesSupported []string `xml:"H264ProfilesSupported"`
					} `xml:"H264"`
				} `xml:"Options"`
			} `xml:"GetVideoEncoderConfigurationOptionsResponse"`
		} `xml:"Body"`
	}

	err := xml.Unmarshal([]byte(bodyStr), &envelope)
	if err == nil {
		// Successfully parsed the XML
		h264Options := envelope.Body.GetVideoEncoderConfigurationOptionsResponse.Options.H264

		// Extract resolutions
		for _, res := range h264Options.ResolutionsAvailable {
			options.ResolutionOptions = append(options.ResolutionOptions, VideoResolution{
				Width:  res.Width,
				Height: res.Height,
			})
		}

		// Extract frame rate range
		minFrameRate := h264Options.FrameRateRange.Min
		maxFrameRate := h264Options.FrameRateRange.Max
		fmt.Printf("ParseH264Options: Frame rate range: Min=%d, Max=%d\n", minFrameRate, maxFrameRate)
		if maxFrameRate > 0 {
			// Generate a more reasonable set of frame rate options within the range
			// Use a non-linear distribution to have more options at the lower end
			if maxFrameRate <= 30 {
				// For standard frame rates (up to 30 fps), provide common values
				for _, rate := range []int{1, 2, 5, 8, 10, 15, 20, 25, 30} {
					if rate >= minFrameRate && rate <= maxFrameRate {
						options.FrameRateOptions = append(options.FrameRateOptions, rate)
					}
				}
				// Always include max frame rate if it's not already in the list
				lastIncludedRate := 0
				if len(options.FrameRateOptions) > 0 {
					lastIncludedRate = options.FrameRateOptions[len(options.FrameRateOptions)-1]
				}
				if lastIncludedRate != maxFrameRate {
					options.FrameRateOptions = append(options.FrameRateOptions, maxFrameRate)
				}
			} else {
				// For high frame rates (>30), provide a wider range
				rates := []int{1, 5, 10, 15, 20, 25, 30, 50, 60}
				for _, rate := range rates {
					if rate >= minFrameRate && rate <= maxFrameRate {
						options.FrameRateOptions = append(options.FrameRateOptions, rate)
					}
				}
				// Include the max
				if maxFrameRate > rates[len(rates)-1] {
					options.FrameRateOptions = append(options.FrameRateOptions, maxFrameRate)
				}
			}
			// If no values were added (which shouldn't happen), add min and max
			if len(options.FrameRateOptions) == 0 {
				options.FrameRateOptions = append(options.FrameRateOptions, minFrameRate, maxFrameRate)
			}
		}

		// Extract encoding interval range and GOP length range
		// Note: GOP length is directly related to encoding interval in many cameras
		minInterval := h264Options.EncodingIntervalRange.Min
		maxInterval := h264Options.EncodingIntervalRange.Max

		minGovLength := h264Options.GovLengthRange.Min
		maxGovLength := h264Options.GovLengthRange.Max

		// If we have no explicit GOP length range, use encoding interval range
		if minGovLength == 0 && maxGovLength == 0 {
			minGovLength = minInterval
			maxGovLength = maxInterval
		}

		fmt.Printf("ParseH264Options: Encoding interval range: Min=%d, Max=%d\n", minInterval, maxInterval)
		fmt.Printf("ParseH264Options: GOP length range: Min=%d, Max=%d\n", minGovLength, maxGovLength)

		if maxInterval > 0 {
			// For encoding intervals, provide a reasonable set of options
			// These are typically small integers, so a linear distribution is fine
			commonIntervals := []int{1, 2, 4, 5, 10, 15, 20, 25, 30}
			for _, interval := range commonIntervals {
				if interval >= minInterval && interval <= maxInterval {
					options.EncodingIntervals = append(options.EncodingIntervals, interval)
				}
			}
			// Always include max interval if it's not already in the list
			lastIncludedInterval := 0
			if len(options.EncodingIntervals) > 0 {
				lastIncludedInterval = options.EncodingIntervals[len(options.EncodingIntervals)-1]
			}
			if lastIncludedInterval != maxInterval && maxInterval > 0 {
				options.EncodingIntervals = append(options.EncodingIntervals, maxInterval)
			}
		}

		// Extract bitrate range
		minBitrate := h264Options.BitrateRange.Min
		maxBitrate := h264Options.BitrateRange.Max
		fmt.Printf("ParseH264Options: Bitrate range: Min=%d, Max=%d\n", minBitrate, maxBitrate)

		if maxBitrate > 0 {
			// Generate a logarithmic bitrate distribution
			// This provides more options in the lower range and fewer in the higher range
			// which is more useful for users

			// Determine number of steps based on the range
			numSteps := 8
			if maxBitrate > 10000000 { // If max is >10Mbps, provide more options
				numSteps = 10
			}

			// For very low bitrate cameras, use a linear scale
			if maxBitrate <= 1000000 {
				step := maxBitrate / numSteps
				for i := 1; i <= numSteps; i++ {
					bitrate := minBitrate + step*i
					if bitrate > maxBitrate {
						bitrate = maxBitrate
					}
					options.BitrateOptions = append(options.BitrateOptions, bitrate)
				}
			} else {
				// For normal/high bitrate cameras, use a logarithmic scale
				// This gives more options at lower bitrates where differences matter more
				logMin := math.Log(float64(minBitrate))
				logMax := math.Log(float64(maxBitrate))
				step := (logMax - logMin) / float64(numSteps-1)

				for i := 0; i < numSteps; i++ {
					logVal := logMin + float64(i)*step
					bitrate := int(math.Round(math.Exp(logVal)))
					options.BitrateOptions = append(options.BitrateOptions, bitrate)
				}
			}

			// Always ensure the max bitrate is included
			if len(options.BitrateOptions) > 0 && options.BitrateOptions[len(options.BitrateOptions)-1] != maxBitrate {
				options.BitrateOptions = append(options.BitrateOptions, maxBitrate)
			}
		}

		// Extract H264 profiles
		options.H264ProfileOptions = h264Options.H264ProfilesSupported

		// If we got data through the XML parsing, return it
		if len(options.ResolutionOptions) > 0 {
			return options
		}
	}

	// If XML parsing failed, fall back to regex/string parsing
	fmt.Printf("ParseH264Options: XML parsing failed with error: %v. Falling back to regex parsing.\n", err)

	// Extract resolutions using regex
	resolutionRegex := `<tt:Width>(\d+)<\/tt:Width>\s*<tt:Height>(\d+)<\/tt:Height>`
	reResolution := regexp.MustCompile(resolutionRegex)
	resMatches := reResolution.FindAllStringSubmatch(bodyStr, -1)

	for _, match := range resMatches {
		if len(match) == 3 {
			width, _ := strconv.Atoi(match[1])
			height, _ := strconv.Atoi(match[2])
			if width > 0 && height > 0 {
				options.ResolutionOptions = append(options.ResolutionOptions, VideoResolution{
					Width:  width,
					Height: height,
				})
			}
		}
	}

	// Extract frame rate range
	frameRateRegex := `<tt:FrameRateRange>\s*<tt:Min>(\d+)<\/tt:Min>\s*<tt:Max>(\d+)<\/tt:Max>\s*<\/tt:FrameRateRange>`
	reFrameRate := regexp.MustCompile(frameRateRegex)
	frMatch := reFrameRate.FindStringSubmatch(bodyStr)

	if len(frMatch) == 3 {
		minFrameRate, _ := strconv.Atoi(frMatch[1])
		maxFrameRate, _ := strconv.Atoi(frMatch[2])
		fmt.Printf("ParseH264Options: Frame rate range from regex: Min=%d, Max=%d\n", minFrameRate, maxFrameRate)

		// Generate frame rate options using the same logic as above
		if maxFrameRate > 0 {
			if maxFrameRate <= 30 {
				for _, rate := range []int{1, 2, 5, 8, 10, 15, 20, 25, 30} {
					if rate >= minFrameRate && rate <= maxFrameRate {
						options.FrameRateOptions = append(options.FrameRateOptions, rate)
					}
				}
			} else {
				rates := []int{1, 5, 10, 15, 20, 25, 30, 50, 60}
				for _, rate := range rates {
					if rate >= minFrameRate && rate <= maxFrameRate {
						options.FrameRateOptions = append(options.FrameRateOptions, rate)
					}
				}
				// Include the max
				if maxFrameRate > rates[len(rates)-1] {
					options.FrameRateOptions = append(options.FrameRateOptions, maxFrameRate)
				}
			}
		}
	} else {
		// Default frame rates if we couldn't extract from the XML
		options.FrameRateOptions = []int{1, 5, 10, 15, 20, 25, 30}
	}

	// Extract encoding interval range
	encodingIntervalRegex := `<tt:EncodingIntervalRange>\s*<tt:Min>(\d+)<\/tt:Min>\s*<tt:Max>(\d+)<\/tt:Max>\s*<\/tt:EncodingIntervalRange>`
	reEncodingInterval := regexp.MustCompile(encodingIntervalRegex)
	eiMatch := reEncodingInterval.FindStringSubmatch(bodyStr)

	if len(eiMatch) == 3 {
		minInterval, _ := strconv.Atoi(eiMatch[1])
		maxInterval, _ := strconv.Atoi(eiMatch[2])
		fmt.Printf("ParseH264Options: Encoding interval range from regex: Min=%d, Max=%d\n", minInterval, maxInterval)

		// Generate encoding interval options
		commonIntervals := []int{1, 2, 4, 5, 10, 15, 20, 25, 30}
		for _, interval := range commonIntervals {
			if interval >= minInterval && interval <= maxInterval {
				options.EncodingIntervals = append(options.EncodingIntervals, interval)
			}
		}
		// Include the max
		if len(options.EncodingIntervals) > 0 && options.EncodingIntervals[len(options.EncodingIntervals)-1] != maxInterval && maxInterval > 0 {
			options.EncodingIntervals = append(options.EncodingIntervals, maxInterval)
		}
	} else {
		// Default encoding intervals
		options.EncodingIntervals = []int{1, 5, 10, 15, 20, 25, 30}
	}

	// Extract GOP length range - we can use this information if needed but not storing the match result
	govLengthRegex := `<tt:GovLengthRange>\s*<tt:Min>(\d+)<\/tt:Min>\s*<tt:Max>(\d+)<\/tt:Max>\s*<\/tt:GovLengthRange>`
	regexp.MustCompile(govLengthRegex) // Just compile but not storing the result since we're not using it

	// Extract bitrate range
	bitrateRegex := `<tt:BitrateRange>\s*<tt:Min>(\d+)<\/tt:Min>\s*<tt:Max>(\d+)<\/tt:Max>\s*<\/tt:BitrateRange>`
	reBitrate := regexp.MustCompile(bitrateRegex)
	brMatch := reBitrate.FindStringSubmatch(bodyStr)

	if len(brMatch) == 3 {
		minBitrate, _ := strconv.Atoi(brMatch[1])
		maxBitrate, _ := strconv.Atoi(brMatch[2])
		fmt.Printf("ParseH264Options: Bitrate range from regex: Min=%d, Max=%d\n", minBitrate, maxBitrate)

		// Generate bitrate options using the same logic as above
		if maxBitrate > 0 {
			numSteps := 8
			if maxBitrate > 10000000 {
				numSteps = 10
			}

			if maxBitrate <= 1000000 {
				step := maxBitrate / numSteps
				for i := 1; i <= numSteps; i++ {
					bitrate := minBitrate + step*i
					if bitrate > maxBitrate {
						bitrate = maxBitrate
					}
					options.BitrateOptions = append(options.BitrateOptions, bitrate)
				}
			} else {
				logMin := math.Log(float64(minBitrate))
				logMax := math.Log(float64(maxBitrate))
				step := (logMax - logMin) / float64(numSteps-1)

				for i := 0; i < numSteps; i++ {
					logVal := logMin + float64(i)*step
					bitrate := int(math.Round(math.Exp(logVal)))
					options.BitrateOptions = append(options.BitrateOptions, bitrate)
				}
			}
		}
	} else {
		// Default bitrate options
		options.BitrateOptions = []int{
			64000, 128000, 256000, 512000, 1000000, 2000000, 4000000, 8000000,
		}
	}

	// Extract H264 profile options
	h264ProfilesRegex := `<tt:H264ProfilesSupported>(.*?)<\/tt:H264ProfilesSupported>`
	reProfiles := regexp.MustCompile(h264ProfilesRegex)
	profilesMatches := reProfiles.FindAllStringSubmatch(bodyStr, -1)

	if len(profilesMatches) > 0 {
		for _, match := range profilesMatches {
			if len(match) == 2 {
				profile := strings.TrimSpace(match[1])
				if profile != "" {
					options.H264ProfileOptions = append(options.H264ProfileOptions, profile)
				}
			}
		}
	} else {
		// Try another approach for profiles
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

	// If we have no profiles, use defaults
	if len(options.H264ProfileOptions) == 0 {
		options.H264ProfileOptions = []string{"Baseline", "Main", "High"}
	}

	// If we have no resolutions, use defaults
	if len(options.ResolutionOptions) == 0 {
		options.ResolutionOptions = []VideoResolution{
			{Width: 1920, Height: 1080},
			{Width: 1280, Height: 720},
			{Width: 640, Height: 480},
			{Width: 320, Height: 240},
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
