package onvif_test

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"regexp"
	"strings"
)

// ContainsFault checks if an XML response contains a SOAP fault
func ContainsFault(xmlData []byte) bool {
	// Simple string check for fault element
	return bytes.Contains(xmlData, []byte("<Fault>")) || bytes.Contains(xmlData, []byte("<fault>")) ||
		bytes.Contains(xmlData, []byte("<soap:Fault>")) || bytes.Contains(xmlData, []byte("<s:Fault>"))
}

// extractFaultString extracts the fault string from a SOAP fault response
func extractFaultString(body []byte) string {
	// Parse the XML to get the fault details
	var envelope struct {
		XMLName xml.Name `xml:"http://www.w3.org/2003/05/soap-envelope Envelope"`
		Body    struct {
			Fault struct {
				Code   string `xml:"Code>Value"`
				Reason struct {
					Text string `xml:"Text"`
				} `xml:"Reason"`
				Detail string `xml:"Detail"`
			} `xml:"Fault"`
		} `xml:"Body"`
	}

	if err := xml.Unmarshal(body, &envelope); err != nil {
		// If we can't parse the XML, try a simple regex approach
		re := regexp.MustCompile(`<faultstring.*?>(.*?)</faultstring>`)
		if matches := re.FindSubmatch(body); len(matches) > 1 {
			return string(matches[1])
		}
		return "unknown error (failed to parse fault response)"
	}

	if envelope.Body.Fault.Reason.Text != "" {
		return envelope.Body.Fault.Reason.Text
	}
	if envelope.Body.Fault.Detail != "" {
		return envelope.Body.Fault.Detail
	}
	if envelope.Body.Fault.Code != "" {
		return envelope.Body.Fault.Code
	}
	return "unknown error (no fault details found)"
}

// PrintFormattedXML formats and prints XML data
func PrintFormattedXML(xmlData []byte) (string, error) {
	// Re-encode with indentation
	var prettyXML bytes.Buffer
	decoder := xml.NewDecoder(bytes.NewReader(xmlData))
	encoder := xml.NewEncoder(&prettyXML)
	encoder.Indent("", "  ")

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("error decoding XML: %v", err)
		}
		err = encoder.EncodeToken(token)
		if err != nil {
			return "", fmt.Errorf("error encoding XML: %v", err)
		}
	}

	err := encoder.Flush()
	if err != nil {
		return "", fmt.Errorf("error flushing XML encoder: %v", err)
	}

	return prettyXML.String(), nil
}

// FormatXML formats XML byte array into a human-readable string with proper indentation
func FormatXML(input []byte) (string, error) {
	var buffer bytes.Buffer
	decoder := xml.NewDecoder(bytes.NewReader(input))
	encoder := xml.NewEncoder(&buffer)
	encoder.Indent("", "  ")

	for {
		token, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}
		if err := encoder.EncodeToken(token); err != nil {
			return "", err
		}
	}

	if err := encoder.Flush(); err != nil {
		return "", err
	}
	return buffer.String(), nil
}

// ParseH264Options parses the H264 options from the response
func ParseH264Options(optionsResp *VideoEncoderConfigurationOptionsResponse) *H264Options {
	fmt.Printf("ParseH264Options: Starting to parse response: %+v\n", optionsResp)

	if optionsResp == nil {
		fmt.Println("ParseH264Options: Received nil response")
		return &H264Options{}
	}

	h264 := optionsResp.Body.GetVideoEncoderConfigurationOptionsResponse.Options.H264
	fmt.Printf("ParseH264Options: H264 options data: %+v\n", h264)

	// Extract resolutions
	var resolutions []Resolution
	if len(h264.ResolutionsAvailable) > 0 {
		fmt.Printf("ParseH264Options: Found %d resolutions in XML\n", len(h264.ResolutionsAvailable))
		for _, res := range h264.ResolutionsAvailable {
			fmt.Printf("ParseH264Options: Adding resolution %dx%d\n", res.Width, res.Height)
			resolutions = append(resolutions, Resolution{
				Width:  res.Width,
				Height: res.Height,
			})
		}
	} else {
		fmt.Println("ParseH264Options: No resolutions found in XML")
	}

	return &H264Options{
		ResolutionsAvailable: resolutions,
		GovLengthRange: Range{
			Min: h264.GovLengthRange.Min,
			Max: h264.GovLengthRange.Max,
		},
		FrameRateRange: Range{
			Min: h264.FrameRateRange.Min,
			Max: h264.FrameRateRange.Max,
		},
		EncodingIntervalRange: Range{
			Min: h264.EncodingIntervalRange.Min,
			Max: h264.EncodingIntervalRange.Max,
		},
		H264ProfilesSupported: h264.H264ProfilesSupported,
	}
}

// ParseJpegOptions parses the JPEG options from the response
func ParseJpegOptions(optionsResp *VideoEncoderConfigurationOptionsResponse) *JpegOptions {
	jpeg := optionsResp.Body.GetVideoEncoderConfigurationOptionsResponse.Options.JPEG

	// Extract resolutions
	var resolutions []Resolution
	for _, res := range jpeg.ResolutionsAvailable {
		resolutions = append(resolutions, Resolution{
			Width:  res.Width,
			Height: res.Height,
		})
	}

	return &JpegOptions{
		ResolutionsAvailable: resolutions,
		FrameRateRange: Range{
			Min: jpeg.FrameRateRange.Min,
			Max: jpeg.FrameRateRange.Max,
		},
		EncodingIntervalRange: Range{
			Min: jpeg.EncodingIntervalRange.Min,
			Max: jpeg.EncodingIntervalRange.Max,
		},
	}
}

// ParseMpeg4Options parses the MPEG4 options from the response
func ParseMpeg4Options(optionsResp *VideoEncoderConfigurationOptionsResponse) *Mpeg4Options {
	mpeg4 := optionsResp.Body.GetVideoEncoderConfigurationOptionsResponse.Options.MPEG4

	// Extract resolutions
	var resolutions []Resolution
	for _, res := range mpeg4.ResolutionsAvailable {
		resolutions = append(resolutions, Resolution{
			Width:  res.Width,
			Height: res.Height,
		})
	}

	return &Mpeg4Options{
		ResolutionsAvailable: resolutions,
		GovLengthRange: Range{
			Min: mpeg4.GovLengthRange.Min,
			Max: mpeg4.GovLengthRange.Max,
		},
		FrameRateRange: Range{
			Min: mpeg4.FrameRateRange.Min,
			Max: mpeg4.FrameRateRange.Max,
		},
		EncodingIntervalRange: Range{
			Min: mpeg4.EncodingIntervalRange.Min,
			Max: mpeg4.EncodingIntervalRange.Max,
		},
		Mpeg4ProfilesSupported: mpeg4.Mpeg4ProfilesSupported,
	}
}

// CombineVideoEncoderOptions combines various encoder options into a single structure
func CombineVideoEncoderOptions(optionsResp *VideoEncoderConfigurationOptionsResponse) *VideoEncoderOptions {
	return &VideoEncoderOptions{
		H264:  ParseH264Options(optionsResp),
		JPEG:  ParseJpegOptions(optionsResp),
		MPEG4: ParseMpeg4Options(optionsResp),
		QualityRange: Range{
			Min: optionsResp.Body.GetVideoEncoderConfigurationOptionsResponse.Options.QualityRange.Min,
			Max: optionsResp.Body.GetVideoEncoderConfigurationOptionsResponse.Options.QualityRange.Max,
		},
	}
}

// ReadUserInput reads a string from user input with a prompt
func ReadUserInput(prompt string) string {
	var input string
	fmt.Print(prompt)
	fmt.Scanln(&input)
	return strings.TrimSpace(input)
}

// ReadIntInput reads an integer from user input with range validation
func ReadIntInput(prompt string, min, max int) int {
	for {
		input := ReadUserInput(prompt)

		// Check if input is an integer
		var value int
		_, err := fmt.Sscanf(input, "%d", &value)
		if err != nil {
			fmt.Println("❌ Invalid input. Please enter a number.")
			continue
		}

		// Check if input is within range
		if value < min || value > max {
			fmt.Printf("❌ Input must be between %d and %d.\n", min, max)
			continue
		}

		return value
	}
}

// ReadIntRangeInput reads an integer from user input with range validation and a prompt displaying the range
func ReadIntRangeInput(min, max int) int {
	fmt.Printf("Enter a value between %d and %d: ", min, max)
	return ReadIntInput("", min, max)
}
