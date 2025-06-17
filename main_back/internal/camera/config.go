package camera

import (
	"fmt"
	"main_back/pkg/models"

	"github.com/videonext/onvif/profiles/media"
)

// GetProfilesAndConfigs returns all profile tokens and config tokens.
func GetProfilesAndConfigs(client *CameraClient) (profileTokens, configTokens []string, err error) {
	resp, err := client.Media.GetProfiles(&media.GetProfiles{})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get profiles: %w", err)
	}

	for _, profile := range resp.Profiles {
		profileTokens = append(profileTokens, string(profile.Token))
		if profile.VideoEncoderConfiguration.Token != "" {
			configTokens = append(configTokens, string(profile.VideoEncoderConfiguration.Token))
		}
	}
	if len(profileTokens) == 0 || len(configTokens) == 0 {
		return nil, nil, fmt.Errorf("no usable profiles or config tokens found")
	}

	return profileTokens, configTokens, nil
}

// GetCurrentEncoderOptions returns the available encoder options for a given config.
func GetCurrentEncoderOptions(client *CameraClient, profileToken, configToken string) (models.EncoderOption, error) {
	resp, err := client.Media.GetVideoEncoderConfigurationOptions(&media.GetVideoEncoderConfigurationOptions{
		ConfigurationToken: media.ReferenceToken(configToken),
		ProfileToken:       media.ReferenceToken(profileToken),
	})
	if err != nil {
		return models.EncoderOption{}, fmt.Errorf("failed to get encoder options: %w", err)
	}
	var resolutions []models.Resolution
	for _, res := range resp.Options.H264.ResolutionsAvailable {
		resolutions = append(resolutions, models.Resolution{
			Width:  int(res.Width),
			Height: int(res.Height),
		})
	}
	var fpsList []int
	// Iterate from min to max frame rate
	for fps := resp.Options.H264.FrameRateRange.Min; fps <= resp.Options.H264.FrameRateRange.Max; fps++ {
		fpsList = append(fpsList, int(fps))
	}
	var bitrateList []int
	// Check if bitrate range is available in the extension
	if resp.Options.Extension.H264.BitrateRange.Min > 0 && resp.Options.Extension.H264.BitrateRange.Max > 0 {
		// Use only the actual min and max values from camera's range
		minBitrate := int(resp.Options.Extension.H264.BitrateRange.Min)
		maxBitrate := int(resp.Options.Extension.H264.BitrateRange.Max)

		// Only provide min and max bitrate options
		bitrateList = append(bitrateList, minBitrate)
		if minBitrate != maxBitrate {
			bitrateList = append(bitrateList, maxBitrate)
		}
	}

	return models.EncoderOption{
		Resolutions: resolutions,
		FPSOptions:  fpsList,
		Bitrate:     bitrateList,
	}, nil
}

// GetCurrentConfig retrieves the actual current video encoder configuration.
func GetCurrentConfig(client *CameraClient, configToken string) (models.EncoderConfig, error) {
	resp, err := client.Media.GetVideoEncoderConfiguration(&media.GetVideoEncoderConfiguration{
		ConfigurationToken: media.ReferenceToken(configToken),
	})
	if err != nil {
		return models.EncoderConfig{}, fmt.Errorf("failed to get encoder config: %w", err)
	}

	cfg := resp.Configuration
	return models.EncoderConfig{
		Resolution: models.Resolution{
			Width:  int(cfg.Resolution.Width),
			Height: int(cfg.Resolution.Height),
		},
		Quality: int(cfg.Quality),
		FPS:     int(cfg.RateControl.FrameRateLimit),
		Bitrate: int(cfg.RateControl.BitrateLimit),
	}, nil
}

// SetEncoderConfig updates the camera's encoder configuration.
func SetEncoderConfig(client *CameraClient, configToken string, config models.EncoderConfig, input models.EncoderConfig) error {

	resp, _ := client.Media.GetVideoEncoderConfiguration(&media.GetVideoEncoderConfiguration{
		ConfigurationToken: media.ReferenceToken(configToken),
	})

	cfg := resp.Configuration

	// Use existing values if input values are not provided (0)
	if input.Quality == 0 {
		input.Quality = config.Quality
	}
	if input.FPS == 0 {
		input.FPS = config.FPS
	}
	if input.Bitrate == 0 {
		// Default to a reasonable bitrate based on resolution if not provided
		if config.Bitrate > 0 {
			input.Bitrate = config.Bitrate
		} else {
			// Calculate default bitrate based on resolution (rough estimate)
			pixels := input.Resolution.Width * input.Resolution.Height
			if pixels > 2000000 { // 1080p+
				input.Bitrate = 8192
			} else if pixels > 1000000 { // 720p+
				input.Bitrate = 4096
			} else {
				input.Bitrate = 2048
			}
		}
	}

	cfg.Resolution.Width = int32(input.Resolution.Width)
	cfg.Resolution.Height = int32(input.Resolution.Height)
	cfg.RateControl.FrameRateLimit = int32(input.FPS)
	cfg.Quality = float32(input.Quality)
	cfg.RateControl.BitrateLimit = int32(input.Bitrate)

	req := &media.SetVideoEncoderConfiguration{
		Configuration:    cfg,
		ForcePersistence: true,
	}

	_, err := client.Media.SetVideoEncoderConfiguration(req)
	if err != nil {
		return fmt.Errorf("failed to set video encoder config: %w", err)
	}
	return nil
}
