package onvif_test

import (
	"encoding/xml"

	"github.com/use-go/onvif/xsd/onvif"
)

// Common types used across the application

// Profile represents an ONVIF media profile
type Profile struct {
	Token string
	Name  string
}

// Resolution represents video resolution dimensions
type Resolution struct {
	Width  int
	Height int
}

// Range represents a min/max range for encoder settings
type Range struct {
	Min int
	Max int
}

// VideoEncoderConfig represents video encoder configuration settings
type VideoEncoderConfig struct {
	Token       string
	Name        string
	UseCount    int
	Encoding    string
	Width       int
	Height      int
	FrameRate   int
	BitRate     int
	GovLength   int
	Quality     float64
	H264Profile string
}

// AudioEncoderConfig represents audio encoder configuration settings
type AudioEncoderConfig struct {
	Token      string
	Name       string
	UseCount   int
	Encoding   string
	Bitrate    int
	SampleRate int
}

// AudioDecoderConfig represents audio decoder configuration settings
type AudioDecoderConfig struct {
	Token    string
	Name     string
	UseCount int
}

// ProfilesResponse represents the structure of the GetProfiles response
type ProfilesResponse struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    struct {
		GetProfilesResponse struct {
			Profiles []struct {
				Token              string `xml:"token,attr"`
				VideoEncoderConfig struct {
					Resolution struct {
						Width  int `xml:"Width"`
						Height int `xml:"Height"`
					} `xml:"Resolution"`
				} `xml:"VideoEncoderConfiguration"`
			} `xml:"Profiles"`
		} `xml:"GetProfilesResponse"`
	} `xml:"Body"`
}

// StreamUriResponse represents the structure of the GetStreamUri response
type StreamUriResponse struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    struct {
		GetStreamUriResponse struct {
			MediaUri struct {
				Uri string `xml:"Uri"`
			} `xml:"MediaUri"`
		} `xml:"GetStreamUriResponse"`
	} `xml:"Body"`
}

// VideoEncoderConfigurationOptionsResponse is the structure for encoder options
type VideoEncoderConfigurationOptionsResponse struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    struct {
		GetVideoEncoderConfigurationOptionsResponse struct {
			Options struct {
				QualityRange struct {
					Min int `xml:"Min"`
					Max int `xml:"Max"`
				} `xml:"QualityRange"`

				// JPEG options fields
				JPEG struct {
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
				} `xml:"JPEG"`

				// MPEG4 options fields
				MPEG4 struct {
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
					Mpeg4ProfilesSupported []string `xml:"Mpeg4ProfilesSupported"`
				} `xml:"MPEG4"`

				// H264 options fields
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
					// Extension fields if needed
				} `xml:"Extension"`
			} `xml:"Options"`
		} `xml:"GetVideoEncoderConfigurationOptionsResponse"`
	} `xml:"Body"`
}

// AudioEncoderConfigurationOptionsResponse represents audio encoder options
type AudioEncoderConfigurationOptionsResponse struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    struct {
		GetAudioEncoderConfigurationOptionsResponse struct {
			Options struct {
				Encoding    string `xml:"Encoding"`
				BitrateList struct {
					Items []int `xml:"Items"`
				} `xml:"BitrateList"`
				SampleRateList struct {
					Items []int `xml:"Items"`
				} `xml:"SampleRateList"`
				Options []struct {
					Encoding    string `xml:"Encoding"`
					BitrateList struct {
						Items []int `xml:"Items"`
					} `xml:"BitrateList"`
					SampleRateList struct {
						Items []int `xml:"Items"`
					} `xml:"SampleRateList"`
				} `xml:"Options"`
				EncodingOptions []struct {
					Encoding    string `xml:"Encoding"`
					BitrateList struct {
						Items []int `xml:"Items"`
					} `xml:"BitrateList"`
					SampleRateList struct {
						Items []int `xml:"Items"`
					} `xml:"SampleRateList"`
				} `xml:"EncodingOptions"`
			} `xml:"Options"`
		} `xml:"GetAudioEncoderConfigurationOptionsResponse"`
	} `xml:"Body"`
}

// AudioDecoderConfigOptionsResponse is the structure for audio decoder options
type AudioDecoderConfigOptionsResponse struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    struct {
		GetAudioDecoderConfigurationOptionsResponse struct {
			Options struct {
				AACDecOptions struct {
					Bitrate struct {
						Items []int `xml:"Items"`
					} `xml:"Bitrate"`
					SampleRateRange struct {
						Items []int `xml:"Items"`
					} `xml:"SampleRateRange"`
				} `xml:"AACDecOptions"`
				G711DecOptions struct {
					Bitrate struct {
						Items []int `xml:"Items"`
					} `xml:"Bitrate"`
					SampleRateRange struct {
						Items []int `xml:"Items"`
					} `xml:"SampleRateRange"`
				} `xml:"G711DecOptions"`
				G726DecOptions struct {
					Bitrate struct {
						Items []int `xml:"Items"`
					} `xml:"Bitrate"`
					SampleRateRange struct {
						Items []int `xml:"Items"`
					} `xml:"SampleRateRange"`
				} `xml:"G726DecOptions"`
			} `xml:"Options"`
		} `xml:"GetAudioDecoderConfigurationOptionsResponse"`
	} `xml:"Body"`
}

// CompatibleAudioDecoderConfigsResponse for compatible audio decoders
type CompatibleAudioDecoderConfigsResponse struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    struct {
		GetCompatibleAudioDecoderConfigurationsResponse struct {
			Configurations []struct {
				Token    string `xml:"token,attr"`
				Name     string `xml:"Name"`
				UseCount int    `xml:"UseCount"`
			} `xml:"Configurations"`
		} `xml:"GetCompatibleAudioDecoderConfigurationsResponse"`
	} `xml:"Body"`
}

// GetCompatibleVideoEncoderConfigurationsResponse for compatible video encoders
type GetCompatibleVideoEncoderConfigurationsResponse struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    struct {
		GetCompatibleVideoEncoderConfigurationsResponse struct {
			Configurations []VideoEncoderConfig `xml:"Configurations"`
		} `xml:"GetCompatibleVideoEncoderConfigurationsResponse"`
	} `xml:"Body"`
}

// VideoEncoderConfigResponse represents GetVideoEncoderConfigurations response
type VideoEncoderConfigResponse struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    struct {
		GetVideoEncoderConfigurationsResponse struct {
			Configurations []struct {
				Token      string `xml:"token,attr"`
				Resolution struct {
					Width  int `xml:"Width"`
					Height int `xml:"Height"`
				} `xml:"Resolution"`
			} `xml:"Configurations"`
		} `xml:"GetVideoEncoderConfigurationsResponse"`
	} `xml:"Body"`
}

// H264Options contains options specific to H.264 encoders
type H264Options struct {
	ResolutionsAvailable  []Resolution
	GovLengthRange        Range
	FrameRateRange        Range
	EncodingIntervalRange Range
	H264ProfilesSupported []string
}

// JpegOptions contains options specific to JPEG encoders
type JpegOptions struct {
	ResolutionsAvailable  []Resolution
	FrameRateRange        Range
	EncodingIntervalRange Range
}

// Mpeg4Options contains options specific to MPEG-4 encoders
type Mpeg4Options struct {
	ResolutionsAvailable   []Resolution
	GovLengthRange         Range
	FrameRateRange         Range
	EncodingIntervalRange  Range
	Mpeg4ProfilesSupported []string
}

// VideoEncoderOptions contains all encoder options
type VideoEncoderOptions struct {
	H264         *H264Options
	JPEG         *JpegOptions
	MPEG4        *Mpeg4Options
	QualityRange Range
}

type GetCapabilities struct {
	XMLName  string                   `xml:"tds:GetCapabilities"`
	Category onvif.CapabilityCategory `xml:"tds:Category,omitempty"`
}

// GetCapabilitiesResponse represents the response from a GetCapabilities request
type GetCapabilitiesResponse struct {
	XMLName      xml.Name           `xml:"Envelope"`
	Capabilities onvif.Capabilities `xml:"Body>GetCapabilitiesResponse>Capabilities"`
}

// AnalyticsCapabilities represents analytics service capabilities
type AnalyticsCapabilities struct {
	XAddr                     string `xml:"XAddr,omitempty"`
	RuleSupport               bool   `xml:"RuleSupport,omitempty"`
	AnalyticsModuleSupport    bool   `xml:"AnalyticsModuleSupport,omitempty"`
	SupportedAnalyticsModules string `xml:"SupportedAnalyticsModules,omitempty"`
	SupportedRules            string `xml:"SupportedRules,omitempty"`
	SupportedSyntaxes         string `xml:"SupportedSyntaxes,omitempty"`
}

// DeviceCapabilities represents device service capabilities
type DeviceCapabilities struct {
	XAddr     string `xml:"XAddr,omitempty"`
	Network   string `xml:"Network,omitempty"`
	System    string `xml:"System,omitempty"`
	IO        string `xml:"IO,omitempty"`
	Security  string `xml:"Security,omitempty"`
	Extension string `xml:"Extension,omitempty"`
}

// EventCapabilities represents event service capabilities
type EventCapabilities struct {
	XAddr                                         string `xml:"XAddr,omitempty"`
	WSSubscriptionPolicySupport                   bool   `xml:"WSSubscriptionPolicySupport,omitempty"`
	WSPullPointSupport                            bool   `xml:"WSPullPointSupport,omitempty"`
	WSPausableSubscriptionManagerInterfaceSupport bool   `xml:"WSPausableSubscriptionManagerInterfaceSupport,omitempty"`
	MaxNotificationProducers                      int    `xml:"MaxNotificationProducers,omitempty"`
	MaxPullPoints                                 int    `xml:"MaxPullPoints,omitempty"`
	PersistentNotificationStorage                 bool   `xml:"PersistentNotificationStorage,omitempty"`
}

// ImagingCapabilities represents imaging service capabilities
type ImagingCapabilities struct {
	XAddr string `xml:"XAddr,omitempty"`
}

// MediaCapabilities represents media service capabilities
type MediaCapabilities struct {
	XAddr                 string `xml:"XAddr,omitempty"`
	StreamingCapabilities string `xml:"StreamingCapabilities,omitempty"`
	SnapshotUri           bool   `xml:"SnapshotUri,omitempty"`
	Rotation              bool   `xml:"Rotation,omitempty"`
	VideoSourceMode       bool   `xml:"VideoSourceMode,omitempty"`
	OSD                   bool   `xml:"OSD,omitempty"`
}

// PTZCapabilities represents PTZ service capabilities
type PTZCapabilities struct {
	XAddr string `xml:"XAddr,omitempty"`
}

// CapabilitiesExtension represents capability extensions
type CapabilitiesExtension struct {
	DeviceIO        string `xml:"DeviceIO,omitempty"`
	Display         string `xml:"Display,omitempty"`
	Recording       string `xml:"Recording,omitempty"`
	Search          string `xml:"Search,omitempty"`
	Replay          string `xml:"Replay,omitempty"`
	Receiver        string `xml:"Receiver,omitempty"`
	AnalyticsDevice string `xml:"AnalyticsDevice,omitempty"`
}
