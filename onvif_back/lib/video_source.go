package onvif_back

import (
	"encoding/xml"
	"fmt"
	"io"

	"github.com/use-go/onvif/media"
	xsd "github.com/use-go/onvif/xsd/onvif"
)

// VideoSourceConfiguration represents a video source configuration
type VideoSourceConfiguration struct {
	Token       string
	Name        string
	UseCount    int
	ViewMode    string
	SourceToken string
	Bounds      IntRectangle
}

// IntRectangle represents a rectangle with position and size
type IntRectangle struct {
	X      int
	Y      int
	Width  int
	Height int
}

// IntRange represents a range of integer values
type IntRange struct {
	Min int
	Max int
}

// IntRectangleRange represents ranges for rectangle properties
type IntRectangleRange struct {
	XRange      IntRange
	YRange      IntRange
	WidthRange  IntRange
	HeightRange IntRange
}

// VideoSourceConfigurationOptions represents the available options for video source configuration
type VideoSourceConfigurationOptions struct {
	MaximumNumberOfProfiles int
	BoundsRange             IntRectangleRange
	VideoSourceTokens       []string
}

// VideoSourceConfigurationOptionsResponse is the structure for parsing the response
type VideoSourceConfigurationOptionsResponse struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    struct {
		GetVideoSourceConfigurationOptionsResponse struct {
			Options struct {
				MaximumNumberOfProfiles int `xml:"MaximumNumberOfProfiles,attr"`
				BoundsRange             struct {
					XRange struct {
						Min int `xml:"Min"`
						Max int `xml:"Max"`
					} `xml:"XRange"`
					YRange struct {
						Min int `xml:"Min"`
						Max int `xml:"Max"`
					} `xml:"YRange"`
					WidthRange struct {
						Min int `xml:"Min"`
						Max int `xml:"Max"`
					} `xml:"WidthRange"`
					HeightRange struct {
						Min int `xml:"Min"`
						Max int `xml:"Max"`
					} `xml:"HeightRange"`
				} `xml:"BoundsRange"`
				VideoSourceTokensAvailable []string `xml:"VideoSourceTokensAvailable"`
			} `xml:"Options"`
		} `xml:"GetVideoSourceConfigurationOptionsResponse"`
	} `xml:"Body"`
}

// DeviceEntity represents base device information
type DeviceEntity struct {
	Token string `xml:"token,attr"`
	Name  string `xml:"Name,omitempty"`
}

// VideoResolution represents the resolution of the video
type VideoResolution struct {
	Width  int `xml:"Width"`
	Height int `xml:"Height"`
}

// Rectangle represents a rectangle window for exposure
type Rectangle struct {
	Top    float64 `xml:"top,attr,omitempty"`
	Bottom float64 `xml:"bottom,attr,omitempty"`
	Left   float64 `xml:"left,attr,omitempty"`
	Right  float64 `xml:"right,attr,omitempty"`
}

// Window represents a rectangle window
type Window struct {
	Top    float64 `xml:"top,attr,omitempty"`
	Bottom float64 `xml:"bottom,attr,omitempty"`
	Left   float64 `xml:"left,attr,omitempty"`
	Right  float64 `xml:"right,attr,omitempty"`
}

// BacklightCompensation represents backlight compensation settings
type BacklightCompensation struct {
	Mode  string  `xml:"Mode,omitempty"`
	Level float64 `xml:"Level,omitempty"`
}

// Exposure represents exposure settings
type Exposure struct {
	Mode            string    `xml:"Mode,omitempty"`
	Priority        string    `xml:"Priority,omitempty"`
	Window          Rectangle `xml:"Window,omitempty"`
	MinExposureTime float64   `xml:"MinExposureTime,omitempty"`
	MaxExposureTime float64   `xml:"MaxExposureTime,omitempty"`
	MinGain         float64   `xml:"MinGain,omitempty"`
	MaxGain         float64   `xml:"MaxGain,omitempty"`
	MinIris         float64   `xml:"MinIris,omitempty"`
	MaxIris         float64   `xml:"MaxIris,omitempty"`
	ExposureTime    float64   `xml:"ExposureTime,omitempty"`
	Gain            float64   `xml:"Gain,omitempty"`
	Iris            float64   `xml:"Iris,omitempty"`
}

// FocusConfiguration represents focus settings
type FocusConfiguration struct {
	AutoFocusMode string  `xml:"AutoFocusMode,omitempty"`
	DefaultSpeed  float64 `xml:"DefaultSpeed,omitempty"`
	NearLimit     float64 `xml:"NearLimit,omitempty"`
	FarLimit      float64 `xml:"FarLimit,omitempty"`
}

// WideDynamicRange represents WDR settings
type WideDynamicRange struct {
	Mode  string  `xml:"Mode,omitempty"`
	Level float64 `xml:"Level,omitempty"`
}

// WhiteBalance represents white balance settings
type WhiteBalance struct {
	Mode   string  `xml:"Mode,omitempty"`
	CrGain float64 `xml:"CrGain,omitempty"`
	CbGain float64 `xml:"CbGain,omitempty"`
}

// ImagingSettings represents image configuration settings
type ImagingSettings struct {
	BacklightCompensation BacklightCompensation `xml:"BacklightCompensation,omitempty"`
	Brightness            float64               `xml:"Brightness,omitempty"`
	ColorSaturation       float64               `xml:"ColorSaturation,omitempty"`
	Contrast              float64               `xml:"Contrast,omitempty"`
	Exposure              Exposure              `xml:"Exposure,omitempty"`
	Focus                 FocusConfiguration    `xml:"Focus,omitempty"`
	IrCutFilter           string                `xml:"IrCutFilter,omitempty"`
	Sharpness             float64               `xml:"Sharpness,omitempty"`
	WideDynamicRange      WideDynamicRange      `xml:"WideDynamicRange,omitempty"`
	WhiteBalance          WhiteBalance          `xml:"WhiteBalance,omitempty"`
}

// ToneCompensation represents tone compensation settings
type ToneCompensation struct {
	Mode  string  `xml:"Mode,omitempty"`
	Level float64 `xml:"Level,omitempty"`
}

// Defogging represents defogging settings
type Defogging struct {
	Mode  string  `xml:"Mode,omitempty"`
	Level float64 `xml:"Level,omitempty"`
}

// NoiseReduction represents noise reduction settings
type NoiseReduction struct {
	Level float64 `xml:"Level,omitempty"`
}

// ImageStabilization represents image stabilization settings
type ImageStabilization struct {
	Mode  string  `xml:"Mode,omitempty"`
	Level float64 `xml:"Level,omitempty"`
}

// IrCutFilterAutoAdjustment represents IR cut filter auto adjustment
type IrCutFilterAutoAdjustment struct {
	BoundaryType   string  `xml:"BoundaryType,omitempty"`
	BoundaryOffset float64 `xml:"BoundaryOffset,omitempty"`
	ResponseTime   string  `xml:"ResponseTime,omitempty"`
}

// FocusConfiguration20 represents extended focus configuration
type FocusConfiguration20 struct {
	AutoFocusMode string  `xml:"AutoFocusMode,omitempty"`
	DefaultSpeed  float64 `xml:"DefaultSpeed,omitempty"`
	NearLimit     float64 `xml:"NearLimit,omitempty"`
	FarLimit      float64 `xml:"FarLimit,omitempty"`
	AFMode        string  `xml:"AFMode,omitempty"`
}

// WhiteBalance20 represents extended white balance settings
type WhiteBalance20 struct {
	Mode   string  `xml:"Mode,omitempty"`
	CrGain float64 `xml:"CrGain,omitempty"`
	CbGain float64 `xml:"CbGain,omitempty"`
}

// WideDynamicRange20 represents extended WDR settings
type WideDynamicRange20 struct {
	Mode  string  `xml:"Mode,omitempty"`
	Level float64 `xml:"Level,omitempty"`
}

// Exposure20 represents extended exposure settings
type Exposure20 struct {
	Mode            string    `xml:"Mode,omitempty"`
	Priority        string    `xml:"Priority,omitempty"`
	Window          Rectangle `xml:"Window,omitempty"`
	MinExposureTime float64   `xml:"MinExposureTime,omitempty"`
	MaxExposureTime float64   `xml:"MaxExposureTime,omitempty"`
	MinGain         float64   `xml:"MinGain,omitempty"`
	MaxGain         float64   `xml:"MaxGain,omitempty"`
	MinIris         float64   `xml:"MinIris,omitempty"`
	MaxIris         float64   `xml:"MaxIris,omitempty"`
	ExposureTime    float64   `xml:"ExposureTime,omitempty"`
	Gain            float64   `xml:"Gain,omitempty"`
	Iris            float64   `xml:"Iris,omitempty"`
}

// BacklightCompensation20 represents extended backlight compensation settings
type BacklightCompensation20 struct {
	Mode  string  `xml:"Mode,omitempty"`
	Level float64 `xml:"Level,omitempty"`
}

// ImagingSettings20 represents extended imaging settings
type ImagingSettings20 struct {
	BacklightCompensation     BacklightCompensation20     `xml:"BacklightCompensation,omitempty"`
	Brightness                float64                     `xml:"Brightness,omitempty"`
	ColorSaturation           float64                     `xml:"ColorSaturation,omitempty"`
	Contrast                  float64                     `xml:"Contrast,omitempty"`
	Exposure                  Exposure20                  `xml:"Exposure,omitempty"`
	Focus                     FocusConfiguration20        `xml:"Focus,omitempty"`
	IrCutFilter               string                      `xml:"IrCutFilter,omitempty"`
	Sharpness                 float64                     `xml:"Sharpness,omitempty"`
	WideDynamicRange          WideDynamicRange20          `xml:"WideDynamicRange,omitempty"`
	WhiteBalance              WhiteBalance20              `xml:"WhiteBalance,omitempty"`
	ImageStabilization        ImageStabilization          `xml:"ImageStabilization,omitempty"`
	IrCutFilterAutoAdjustment []IrCutFilterAutoAdjustment `xml:"IrCutFilterAutoAdjustment,omitempty"`
	ToneCompensation          ToneCompensation            `xml:"ToneCompensation,omitempty"`
	Defogging                 Defogging                   `xml:"Defogging,omitempty"`
	NoiseReduction            NoiseReduction              `xml:"NoiseReduction,omitempty"`
}

// VideoSourceExtension represents extensions to the video source
type VideoSourceExtension struct {
	Imaging   *ImagingSettings20 `xml:"Imaging,omitempty"`
	Extension interface{}        `xml:"Extension,omitempty"`
}

// GetVideoSources represents the ONVIF GetVideoSources request
type GetVideoSources struct {
	XMLName xml.Name `xml:"trt:GetVideoSources"`
}

// GetVideoSourcesResponse represents the ONVIF GetVideoSources response
type GetVideoSourcesResponse struct {
	VideoSources []VideoSource `xml:"VideoSources"`
}

// VideoSource represents a physical video input on an ONVIF device
type VideoSource struct {
	DeviceEntity
	Framerate  float64               `xml:"Framerate"`
	Resolution VideoResolution       `xml:"Resolution"`
	Imaging    *ImagingSettings      `xml:"Imaging,omitempty"`
	Extension  *VideoSourceExtension `xml:"Extension,omitempty"`
}

// VideoSourcesResponse is the structure for parsing the GetVideoSources response
type VideoSourcesResponse struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    struct {
		GetVideoSourcesResponse struct {
			VideoSources []struct {
				Token      string  `xml:"token,attr"`
				Name       string  `xml:"Name,omitempty"`
				Framerate  float64 `xml:"Framerate"`
				Resolution struct {
					Width  int `xml:"Width"`
					Height int `xml:"Height"`
				} `xml:"Resolution"`
				Imaging   *ImagingXMLElement   `xml:"Imaging,omitempty"`
				Extension *ExtensionXMLElement `xml:"Extension,omitempty"`
			} `xml:"VideoSources"`
		} `xml:"GetVideoSourcesResponse"`
	} `xml:"Body"`
}

// ImagingXMLElement represents the raw XML element for imaging settings
type ImagingXMLElement struct {
	BacklightCompensation *struct {
		Mode  string  `xml:"Mode,omitempty"`
		Level float64 `xml:"Level,omitempty"`
	} `xml:"BacklightCompensation,omitempty"`
	Brightness      *float64 `xml:"Brightness,omitempty"`
	ColorSaturation *float64 `xml:"ColorSaturation,omitempty"`
	Contrast        *float64 `xml:"Contrast,omitempty"`
	Exposure        *struct {
		Mode            string    `xml:"Mode,omitempty"`
		Priority        string    `xml:"Priority,omitempty"`
		Window          Rectangle `xml:"Window,omitempty"`
		MinExposureTime *float64  `xml:"MinExposureTime,omitempty"`
		MaxExposureTime *float64  `xml:"MaxExposureTime,omitempty"`
		MinGain         *float64  `xml:"MinGain,omitempty"`
		MaxGain         *float64  `xml:"MaxGain,omitempty"`
		MinIris         *float64  `xml:"MinIris,omitempty"`
		MaxIris         *float64  `xml:"MaxIris,omitempty"`
		ExposureTime    *float64  `xml:"ExposureTime,omitempty"`
		Gain            *float64  `xml:"Gain,omitempty"`
		Iris            *float64  `xml:"Iris,omitempty"`
	} `xml:"Exposure,omitempty"`
	Focus *struct {
		AutoFocusMode string   `xml:"AutoFocusMode,omitempty"`
		DefaultSpeed  *float64 `xml:"DefaultSpeed,omitempty"`
		NearLimit     *float64 `xml:"NearLimit,omitempty"`
		FarLimit      *float64 `xml:"FarLimit,omitempty"`
	} `xml:"Focus,omitempty"`
	IrCutFilter      *string  `xml:"IrCutFilter,omitempty"`
	Sharpness        *float64 `xml:"Sharpness,omitempty"`
	WideDynamicRange *struct {
		Mode  string   `xml:"Mode,omitempty"`
		Level *float64 `xml:"Level,omitempty"`
	} `xml:"WideDynamicRange,omitempty"`
	WhiteBalance *struct {
		Mode   string   `xml:"Mode,omitempty"`
		CrGain *float64 `xml:"CrGain,omitempty"`
		CbGain *float64 `xml:"CbGain,omitempty"`
	} `xml:"WhiteBalance,omitempty"`
}

// ExtensionXMLElement represents the raw XML element for extension settings
type ExtensionXMLElement struct {
	Imaging *struct {
		BacklightCompensation *struct {
			Mode  string   `xml:"Mode,omitempty"`
			Level *float64 `xml:"Level,omitempty"`
		} `xml:"BacklightCompensation,omitempty"`
		Brightness      *float64 `xml:"Brightness,omitempty"`
		ColorSaturation *float64 `xml:"ColorSaturation,omitempty"`
		Contrast        *float64 `xml:"Contrast,omitempty"`
		Exposure        *struct {
			Mode            string    `xml:"Mode,omitempty"`
			Priority        string    `xml:"Priority,omitempty"`
			Window          Rectangle `xml:"Window,omitempty"`
			MinExposureTime *float64  `xml:"MinExposureTime,omitempty"`
			MaxExposureTime *float64  `xml:"MaxExposureTime,omitempty"`
			MinGain         *float64  `xml:"MinGain,omitempty"`
			MaxGain         *float64  `xml:"MaxGain,omitempty"`
			MinIris         *float64  `xml:"MinIris,omitempty"`
			MaxIris         *float64  `xml:"MaxIris,omitempty"`
			ExposureTime    *float64  `xml:"ExposureTime,omitempty"`
			Gain            *float64  `xml:"Gain,omitempty"`
			Iris            *float64  `xml:"Iris,omitempty"`
		} `xml:"Exposure,omitempty"`
		Focus *struct {
			AutoFocusMode string   `xml:"AutoFocusMode,omitempty"`
			AFMode        string   `xml:"AFMode,omitempty"`
			DefaultSpeed  *float64 `xml:"DefaultSpeed,omitempty"`
			NearLimit     *float64 `xml:"NearLimit,omitempty"`
			FarLimit      *float64 `xml:"FarLimit,omitempty"`
		} `xml:"Focus,omitempty"`
		IrCutFilter      *string  `xml:"IrCutFilter,omitempty"`
		Sharpness        *float64 `xml:"Sharpness,omitempty"`
		WideDynamicRange *struct {
			Mode  string   `xml:"Mode,omitempty"`
			Level *float64 `xml:"Level,omitempty"`
		} `xml:"WideDynamicRange,omitempty"`
		WhiteBalance *struct {
			Mode   string   `xml:"Mode,omitempty"`
			CrGain *float64 `xml:"CrGain,omitempty"`
			CbGain *float64 `xml:"CbGain,omitempty"`
		} `xml:"WhiteBalance,omitempty"`
		ImageStabilization *struct {
			Mode  string   `xml:"Mode,omitempty"`
			Level *float64 `xml:"Level,omitempty"`
		} `xml:"ImageStabilization,omitempty"`
		IrCutFilterAutoAdjustment []struct {
			BoundaryType   string   `xml:"BoundaryType,omitempty"`
			BoundaryOffset *float64 `xml:"BoundaryOffset,omitempty"`
			ResponseTime   string   `xml:"ResponseTime,omitempty"`
		} `xml:"IrCutFilterAutoAdjustment,omitempty"`
		ToneCompensation *struct {
			Mode  string   `xml:"Mode,omitempty"`
			Level *float64 `xml:"Level,omitempty"`
		} `xml:"ToneCompensation,omitempty"`
		Defogging *struct {
			Mode  string   `xml:"Mode,omitempty"`
			Level *float64 `xml:"Level,omitempty"`
		} `xml:"Defogging,omitempty"`
		NoiseReduction *struct {
			Level *float64 `xml:"Level,omitempty"`
		} `xml:"NoiseReduction,omitempty"`
	} `xml:"Imaging,omitempty"`
	Extension interface{} `xml:"Extension,omitempty"`
}

// GetAllVideoSources retrieves all available video sources from the device
func GetAllVideoSources(c *Camera) ([]VideoSource, error) {
	if c.Device == nil {
		return nil, fmt.Errorf("camera not connected")
	}

	// Create request
	req := media.GetVideoSources{}

	// Call the method
	resp, err := c.Device.CallMethod(req)
	if err != nil {
		return nil, fmt.Errorf("error calling GetVideoSources: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	// Parse the response
	var sourceResp VideoSourcesResponse
	if err := xml.Unmarshal(body, &sourceResp); err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %v", err)
	}

	// Convert to our VideoSource type
	var sources []VideoSource
	for _, s := range sourceResp.Body.GetVideoSourcesResponse.VideoSources {
		source := VideoSource{
			DeviceEntity: DeviceEntity{
				Token: s.Token,
				Name:  s.Name,
			},
			Framerate: s.Framerate,
			Resolution: VideoResolution{
				Width:  s.Resolution.Width,
				Height: s.Resolution.Height,
			},
		}

		// Convert Imaging settings if present
		if s.Imaging != nil {
			imaging := &ImagingSettings{}

			// BacklightCompensation
			if s.Imaging.BacklightCompensation != nil {
				imaging.BacklightCompensation = BacklightCompensation{
					Mode: s.Imaging.BacklightCompensation.Mode,
				}
				if s.Imaging.BacklightCompensation.Level != 0 {
					imaging.BacklightCompensation.Level = s.Imaging.BacklightCompensation.Level
				}
			}

			// Basic settings
			if s.Imaging.Brightness != nil {
				imaging.Brightness = *s.Imaging.Brightness
			}
			if s.Imaging.ColorSaturation != nil {
				imaging.ColorSaturation = *s.Imaging.ColorSaturation
			}
			if s.Imaging.Contrast != nil {
				imaging.Contrast = *s.Imaging.Contrast
			}
			if s.Imaging.Sharpness != nil {
				imaging.Sharpness = *s.Imaging.Sharpness
			}
			if s.Imaging.IrCutFilter != nil {
				imaging.IrCutFilter = *s.Imaging.IrCutFilter
			}

			// Exposure
			if s.Imaging.Exposure != nil {
				imaging.Exposure = Exposure{
					Mode:     s.Imaging.Exposure.Mode,
					Priority: s.Imaging.Exposure.Priority,
					Window:   s.Imaging.Exposure.Window,
				}
				if s.Imaging.Exposure.MinExposureTime != nil {
					imaging.Exposure.MinExposureTime = *s.Imaging.Exposure.MinExposureTime
				}
				if s.Imaging.Exposure.MaxExposureTime != nil {
					imaging.Exposure.MaxExposureTime = *s.Imaging.Exposure.MaxExposureTime
				}
				if s.Imaging.Exposure.MinGain != nil {
					imaging.Exposure.MinGain = *s.Imaging.Exposure.MinGain
				}
				if s.Imaging.Exposure.MaxGain != nil {
					imaging.Exposure.MaxGain = *s.Imaging.Exposure.MaxGain
				}
				if s.Imaging.Exposure.MinIris != nil {
					imaging.Exposure.MinIris = *s.Imaging.Exposure.MinIris
				}
				if s.Imaging.Exposure.MaxIris != nil {
					imaging.Exposure.MaxIris = *s.Imaging.Exposure.MaxIris
				}
				if s.Imaging.Exposure.ExposureTime != nil {
					imaging.Exposure.ExposureTime = *s.Imaging.Exposure.ExposureTime
				}
				if s.Imaging.Exposure.Gain != nil {
					imaging.Exposure.Gain = *s.Imaging.Exposure.Gain
				}
				if s.Imaging.Exposure.Iris != nil {
					imaging.Exposure.Iris = *s.Imaging.Exposure.Iris
				}
			}

			// Focus
			if s.Imaging.Focus != nil {
				imaging.Focus = FocusConfiguration{
					AutoFocusMode: s.Imaging.Focus.AutoFocusMode,
				}
				if s.Imaging.Focus.DefaultSpeed != nil {
					imaging.Focus.DefaultSpeed = *s.Imaging.Focus.DefaultSpeed
				}
				if s.Imaging.Focus.NearLimit != nil {
					imaging.Focus.NearLimit = *s.Imaging.Focus.NearLimit
				}
				if s.Imaging.Focus.FarLimit != nil {
					imaging.Focus.FarLimit = *s.Imaging.Focus.FarLimit
				}
			}

			// WideDynamicRange
			if s.Imaging.WideDynamicRange != nil {
				imaging.WideDynamicRange = WideDynamicRange{
					Mode: s.Imaging.WideDynamicRange.Mode,
				}
				if s.Imaging.WideDynamicRange.Level != nil {
					imaging.WideDynamicRange.Level = *s.Imaging.WideDynamicRange.Level
				}
			}

			// WhiteBalance
			if s.Imaging.WhiteBalance != nil {
				imaging.WhiteBalance = WhiteBalance{
					Mode: s.Imaging.WhiteBalance.Mode,
				}
				if s.Imaging.WhiteBalance.CrGain != nil {
					imaging.WhiteBalance.CrGain = *s.Imaging.WhiteBalance.CrGain
				}
				if s.Imaging.WhiteBalance.CbGain != nil {
					imaging.WhiteBalance.CbGain = *s.Imaging.WhiteBalance.CbGain
				}
			}

			source.Imaging = imaging
		}

		// Convert Extension.Imaging if present
		if s.Extension != nil && s.Extension.Imaging != nil {
			extension := &VideoSourceExtension{}

			// Create ImagingSettings20
			imaging20 := &ImagingSettings20{}

			// BacklightCompensation
			if s.Extension.Imaging.BacklightCompensation != nil {
				imaging20.BacklightCompensation = BacklightCompensation20{
					Mode: s.Extension.Imaging.BacklightCompensation.Mode,
				}
				if s.Extension.Imaging.BacklightCompensation.Level != nil {
					imaging20.BacklightCompensation.Level = *s.Extension.Imaging.BacklightCompensation.Level
				}
			}

			// Basic settings
			if s.Extension.Imaging.Brightness != nil {
				imaging20.Brightness = *s.Extension.Imaging.Brightness
			}
			if s.Extension.Imaging.ColorSaturation != nil {
				imaging20.ColorSaturation = *s.Extension.Imaging.ColorSaturation
			}
			if s.Extension.Imaging.Contrast != nil {
				imaging20.Contrast = *s.Extension.Imaging.Contrast
			}
			if s.Extension.Imaging.Sharpness != nil {
				imaging20.Sharpness = *s.Extension.Imaging.Sharpness
			}
			if s.Extension.Imaging.IrCutFilter != nil {
				imaging20.IrCutFilter = *s.Extension.Imaging.IrCutFilter
			}

			// Exposure20
			if s.Extension.Imaging.Exposure != nil {
				imaging20.Exposure = Exposure20{
					Mode:     s.Extension.Imaging.Exposure.Mode,
					Priority: s.Extension.Imaging.Exposure.Priority,
					Window:   s.Extension.Imaging.Exposure.Window,
				}
				if s.Extension.Imaging.Exposure.MinExposureTime != nil {
					imaging20.Exposure.MinExposureTime = *s.Extension.Imaging.Exposure.MinExposureTime
				}
				if s.Extension.Imaging.Exposure.MaxExposureTime != nil {
					imaging20.Exposure.MaxExposureTime = *s.Extension.Imaging.Exposure.MaxExposureTime
				}
				if s.Extension.Imaging.Exposure.MinGain != nil {
					imaging20.Exposure.MinGain = *s.Extension.Imaging.Exposure.MinGain
				}
				if s.Extension.Imaging.Exposure.MaxGain != nil {
					imaging20.Exposure.MaxGain = *s.Extension.Imaging.Exposure.MaxGain
				}
				if s.Extension.Imaging.Exposure.MinIris != nil {
					imaging20.Exposure.MinIris = *s.Extension.Imaging.Exposure.MinIris
				}
				if s.Extension.Imaging.Exposure.MaxIris != nil {
					imaging20.Exposure.MaxIris = *s.Extension.Imaging.Exposure.MaxIris
				}
				if s.Extension.Imaging.Exposure.ExposureTime != nil {
					imaging20.Exposure.ExposureTime = *s.Extension.Imaging.Exposure.ExposureTime
				}
				if s.Extension.Imaging.Exposure.Gain != nil {
					imaging20.Exposure.Gain = *s.Extension.Imaging.Exposure.Gain
				}
				if s.Extension.Imaging.Exposure.Iris != nil {
					imaging20.Exposure.Iris = *s.Extension.Imaging.Exposure.Iris
				}
			}

			// Focus20
			if s.Extension.Imaging.Focus != nil {
				imaging20.Focus = FocusConfiguration20{
					AutoFocusMode: s.Extension.Imaging.Focus.AutoFocusMode,
					AFMode:        s.Extension.Imaging.Focus.AFMode,
				}
				if s.Extension.Imaging.Focus.DefaultSpeed != nil {
					imaging20.Focus.DefaultSpeed = *s.Extension.Imaging.Focus.DefaultSpeed
				}
				if s.Extension.Imaging.Focus.NearLimit != nil {
					imaging20.Focus.NearLimit = *s.Extension.Imaging.Focus.NearLimit
				}
				if s.Extension.Imaging.Focus.FarLimit != nil {
					imaging20.Focus.FarLimit = *s.Extension.Imaging.Focus.FarLimit
				}
			}

			// WideDynamicRange20
			if s.Extension.Imaging.WideDynamicRange != nil {
				imaging20.WideDynamicRange = WideDynamicRange20{
					Mode: s.Extension.Imaging.WideDynamicRange.Mode,
				}
				if s.Extension.Imaging.WideDynamicRange.Level != nil {
					imaging20.WideDynamicRange.Level = *s.Extension.Imaging.WideDynamicRange.Level
				}
			}

			// WhiteBalance20
			if s.Extension.Imaging.WhiteBalance != nil {
				imaging20.WhiteBalance = WhiteBalance20{
					Mode: s.Extension.Imaging.WhiteBalance.Mode,
				}
				if s.Extension.Imaging.WhiteBalance.CrGain != nil {
					imaging20.WhiteBalance.CrGain = *s.Extension.Imaging.WhiteBalance.CrGain
				}
				if s.Extension.Imaging.WhiteBalance.CbGain != nil {
					imaging20.WhiteBalance.CbGain = *s.Extension.Imaging.WhiteBalance.CbGain
				}
			}

			// ImageStabilization
			if s.Extension.Imaging.ImageStabilization != nil {
				imaging20.ImageStabilization = ImageStabilization{
					Mode: s.Extension.Imaging.ImageStabilization.Mode,
				}
				if s.Extension.Imaging.ImageStabilization.Level != nil {
					imaging20.ImageStabilization.Level = *s.Extension.Imaging.ImageStabilization.Level
				}
			}

			// IrCutFilterAutoAdjustment
			if len(s.Extension.Imaging.IrCutFilterAutoAdjustment) > 0 {
				imaging20.IrCutFilterAutoAdjustment = make([]IrCutFilterAutoAdjustment, len(s.Extension.Imaging.IrCutFilterAutoAdjustment))
				for i, adjustment := range s.Extension.Imaging.IrCutFilterAutoAdjustment {
					imaging20.IrCutFilterAutoAdjustment[i] = IrCutFilterAutoAdjustment{
						BoundaryType: adjustment.BoundaryType,
						ResponseTime: adjustment.ResponseTime,
					}
					if adjustment.BoundaryOffset != nil {
						imaging20.IrCutFilterAutoAdjustment[i].BoundaryOffset = *adjustment.BoundaryOffset
					}
				}
			}

			// ToneCompensation
			if s.Extension.Imaging.ToneCompensation != nil {
				imaging20.ToneCompensation = ToneCompensation{
					Mode: s.Extension.Imaging.ToneCompensation.Mode,
				}
				if s.Extension.Imaging.ToneCompensation.Level != nil {
					imaging20.ToneCompensation.Level = *s.Extension.Imaging.ToneCompensation.Level
				}
			}

			// Defogging
			if s.Extension.Imaging.Defogging != nil {
				imaging20.Defogging = Defogging{
					Mode: s.Extension.Imaging.Defogging.Mode,
				}
				if s.Extension.Imaging.Defogging.Level != nil {
					imaging20.Defogging.Level = *s.Extension.Imaging.Defogging.Level
				}
			}

			// NoiseReduction
			if s.Extension.Imaging.NoiseReduction != nil && s.Extension.Imaging.NoiseReduction.Level != nil {
				imaging20.NoiseReduction = NoiseReduction{
					Level: *s.Extension.Imaging.NoiseReduction.Level,
				}
			}

			extension.Imaging = imaging20
			source.Extension = extension
		}

		sources = append(sources, source)
	}

	return sources, nil
}

// GetAllVideoSourceConfigurations retrieves all video source configurations
func GetAllVideoSourceConfigurations(c *Camera) ([]VideoSourceConfiguration, error) {
	if c.Device == nil {
		return nil, fmt.Errorf("camera not connected")
	}

	// Create request
	req := media.GetVideoSourceConfigurations{}

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
			GetVideoSourceConfigurationsResponse struct {
				Configurations []struct {
					Token       string `xml:"token,attr"`
					Name        string `xml:"Name"`
					UseCount    int    `xml:"UseCount"`
					ViewMode    string `xml:"ViewMode"`
					SourceToken string `xml:"SourceToken"`
					Bounds      struct {
						X      int `xml:"x,attr"`
						Y      int `xml:"y,attr"`
						Width  int `xml:"width,attr"`
						Height int `xml:"height,attr"`
					} `xml:"Bounds"`
				} `xml:"Configurations"`
			} `xml:"GetVideoSourceConfigurationsResponse"`
		} `xml:"Body"`
	}

	// Parse the response
	if err := xml.Unmarshal(body, &configResp); err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %v", err)
	}

	// Convert to our type
	var configs []VideoSourceConfiguration
	for _, c := range configResp.Body.GetVideoSourceConfigurationsResponse.Configurations {
		configs = append(configs, VideoSourceConfiguration{
			Token:       c.Token,
			Name:        c.Name,
			UseCount:    c.UseCount,
			ViewMode:    c.ViewMode,
			SourceToken: c.SourceToken,
			Bounds: IntRectangle{
				X:      c.Bounds.X,
				Y:      c.Bounds.Y,
				Width:  c.Bounds.Width,
				Height: c.Bounds.Height,
			},
		})
	}

	return configs, nil
}

// GetVideoSourceConfiguration gets a specific video source configuration by token
func GetVideoSourceConfiguration(c *Camera, token string) (*VideoSourceConfiguration, error) {
	if c.Device == nil {
		return nil, fmt.Errorf("camera not connected")
	}

	// Create request
	req := media.GetVideoSourceConfiguration{
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
			GetVideoSourceConfigurationResponse struct {
				Configuration struct {
					Token       string `xml:"token,attr"`
					Name        string `xml:"Name"`
					UseCount    int    `xml:"UseCount"`
					ViewMode    string `xml:"ViewMode"`
					SourceToken string `xml:"SourceToken"`
					Bounds      struct {
						X      int `xml:"x,attr"`
						Y      int `xml:"y,attr"`
						Width  int `xml:"width,attr"`
						Height int `xml:"height,attr"`
					} `xml:"Bounds"`
				} `xml:"Configuration"`
			} `xml:"GetVideoSourceConfigurationResponse"`
		} `xml:"Body"`
	}

	// Parse the response
	if err := xml.Unmarshal(body, &configResp); err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %v", err)
	}

	config := configResp.Body.GetVideoSourceConfigurationResponse.Configuration

	// Create our config type
	return &VideoSourceConfiguration{
		Token:       config.Token,
		Name:        config.Name,
		UseCount:    config.UseCount,
		ViewMode:    config.ViewMode,
		SourceToken: config.SourceToken,
		Bounds: IntRectangle{
			X:      config.Bounds.X,
			Y:      config.Bounds.Y,
			Width:  config.Bounds.Width,
			Height: config.Bounds.Height,
		},
	}, nil
}

// GetVideoSourceConfigurationOptions gets available options for a video source configuration
func GetVideoSourceConfigurationOptions(c *Camera, configToken string, profileToken string) (*VideoSourceConfigurationOptionsResponse, error) {
	if c.Device == nil {
		return nil, fmt.Errorf("camera not connected")
	}

	// Create request
	req := media.GetVideoSourceConfigurationOptions{}

	// Add tokens if provided
	if configToken != "" {
		req.ConfigurationToken = xsd.ReferenceToken(configToken)
	}
	if profileToken != "" {
		req.ProfileToken = xsd.ReferenceToken(profileToken)
	}

	// Call the method
	resp, err := c.Device.CallMethod(req)
	if err != nil {
		return nil, fmt.Errorf("error calling GetVideoSourceConfigurationOptions: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	// Parse the response
	var optionsResp VideoSourceConfigurationOptionsResponse
	if err := xml.Unmarshal(body, &optionsResp); err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %v", err)
	}

	return &optionsResp, nil
}

// ParseVideoSourceConfigOptions converts the response into a more manageable structure
func ParseVideoSourceConfigOptions(optionsResp *VideoSourceConfigurationOptionsResponse) *VideoSourceConfigurationOptions {
	options := optionsResp.Body.GetVideoSourceConfigurationOptionsResponse.Options

	return &VideoSourceConfigurationOptions{
		MaximumNumberOfProfiles: options.MaximumNumberOfProfiles,
		BoundsRange: IntRectangleRange{
			XRange: IntRange{
				Min: options.BoundsRange.XRange.Min,
				Max: options.BoundsRange.XRange.Max,
			},
			YRange: IntRange{
				Min: options.BoundsRange.YRange.Min,
				Max: options.BoundsRange.YRange.Max,
			},
			WidthRange: IntRange{
				Min: options.BoundsRange.WidthRange.Min,
				Max: options.BoundsRange.WidthRange.Max,
			},
			HeightRange: IntRange{
				Min: options.BoundsRange.HeightRange.Min,
				Max: options.BoundsRange.HeightRange.Max,
			},
		},
		VideoSourceTokens: options.VideoSourceTokensAvailable,
	}
}

// SetVideoSourceConfiguration modifies a video source configuration
func SetVideoSourceConfiguration(
	c *Camera,
	configToken string,
	configName string,
	sourceToken string,
	x int,
	y int,
	width int,
	height int) error {

	if c.Device == nil {
		return fmt.Errorf("camera not connected")
	}

	// Create the configuration request
	setConfigRequest := media.SetVideoSourceConfiguration{
		Configuration: xsd.VideoSourceConfiguration{
			ConfigurationEntity: xsd.ConfigurationEntity{
				Token: xsd.ReferenceToken(configToken),
				Name:  xsd.Name(configName),
			},
			SourceToken: xsd.ReferenceToken(sourceToken),
			Bounds: xsd.IntRectangle{
				X:      x,
				Y:      y,
				Width:  width,
				Height: height,
			},
		},
		ForcePersistence: true,
	}

	// Call the method
	setConfigResp, err := c.Device.CallMethod(setConfigRequest)
	if err != nil {
		return fmt.Errorf("error setting video source configuration: %v", err)
	}
	defer setConfigResp.Body.Close()

	// Read response body to check if successful
	body, err := io.ReadAll(setConfigResp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %v", err)
	}

	// Check if the response has a fault element
	if ContainsFault(body) {
		return fmt.Errorf("server returned an error in response to SetVideoSourceConfiguration")
	}

	return nil
}

// GetRawVideoSourcesXML retrieves the raw XML response from GetVideoSources for debugging
func GetRawVideoSourcesXML(c *Camera) (string, error) {
	if c.Device == nil {
		return "", fmt.Errorf("camera not connected")
	}

	// Create request
	req := media.GetVideoSources{}

	// Call the method
	resp, err := c.Device.CallMethod(req)
	if err != nil {
		return "", fmt.Errorf("error calling GetVideoSources: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	return string(body), nil
}
