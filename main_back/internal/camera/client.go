package camera

import (
	"fmt"
	"time"

	"main_back/pkg/models"

	"github.com/videonext/onvif/profiles/media"
	"github.com/videonext/onvif/soap"
)

type CameraClient struct {
	Camera models.Camera
	Client *soap.Client
	Media  media.Media // ONVIF Media service client
}

// NewCameraClient connects to the ONVIF camera and returns a usable CameraClient
func NewCameraClient(cam models.Camera) (*CameraClient, error) {
	// Initialize SOAP client with timeout
	client := soap.NewClient(soap.WithTimeout(5 * time.Second))
	client.AddHeader(soap.NewWSSSecurityHeader(cam.Username, cam.Password, time.Now()))

	// Determine port to use (default 80 if port is 0)
	port := cam.Port
	if port == 0 {
		port = 80
	}

	// Determine URL path to use (default onvif/media_service if URL is empty)
	urlPath := cam.URL
	if urlPath == "" {
		urlPath = "onvif/media_service"
	}

	// Build media service endpoint with port and URL path
	endpoint := fmt.Sprintf("http://%s:%d/%s", cam.IP, port, urlPath)

	// Initialize ONVIF media service
	mediaService := media.NewMedia(client, endpoint)

	return &CameraClient{
		Camera: cam,
		Client: client,
		Media:  mediaService,
	}, nil
}

// NewFakeCameraClient creates a simulated camera client for fake cameras
func NewFakeCameraClient(cam models.Camera) *CameraClient {
	return &CameraClient{
		Camera: cam,
		// For fake cameras, we don't need actual connections to ONVIF services
		Client: nil,
		Media:  nil,
	}
}

// GetStreamURI retrieves the RTSP stream URI for a given profile.
func (c *CameraClient) GetStreamURI(profileToken string) (string, error) {
	// Handle fake cameras with simulated data
	if c.Camera.IsFake {
		// For fake cameras, return a simulated RTSP URI
		return fmt.Sprintf("rtsp://%s/fake_stream", c.Camera.IP), nil
	}

	// For real cameras, proceed with actual ONVIF calls
	request := &media.GetStreamUri{
		StreamSetup: media.StreamSetup{
			Stream: media.StreamType("RTP-Unicast"),
			Transport: media.Transport{
				Protocol: media.TransportProtocol("RTSP"),
			},
		},
		ProfileToken: media.ReferenceToken(profileToken),
	}

	resp, err := c.Media.GetStreamUri(request)
	if err != nil {
		return "", fmt.Errorf("failed to get stream URI for profile %s: %w", profileToken, err)
	}

	if resp == nil || resp.MediaUri.Uri == "" {
		return "", fmt.Errorf("received empty or invalid stream URI response for profile %s", profileToken)
	}

	return string(resp.MediaUri.Uri), nil
}
