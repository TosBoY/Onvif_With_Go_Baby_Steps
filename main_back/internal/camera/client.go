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

	// Build media service endpoint (typically /onvif/media_service)
	endpoint := fmt.Sprintf("http://%s/onvif/media_service", cam.IP)

	// Initialize ONVIF media service
	mediaService := media.NewMedia(client, endpoint)

	return &CameraClient{
		Camera: cam,
		Client: client,
		Media:  mediaService,
	}, nil
}

// GetStreamURI retrieves the RTSP stream URI for a given profile.
func (c *CameraClient) GetStreamURI(profileToken string) (string, error) {
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
