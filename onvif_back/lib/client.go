package onvif_test

import (
	"fmt"
	
	"github.com/use-go/onvif"
)

// Camera represents an ONVIF camera connection with its credentials
type Camera struct {
	IP       string
	Port     int
	Username string
	Password string
	Device   *onvif.Device
}

// NewCamera creates a new Camera instance
func NewCamera(ip string, port int, username, password string) *Camera {
	return &Camera{
		IP:       ip,
		Port:     port,
		Username: username,
		Password: password,
	}
}

// Connect establishes a connection to the camera
func (c *Camera) Connect() error {
	dev, err := onvif.NewDevice(onvif.DeviceParams{
		Xaddr:    fmt.Sprintf("%s:%d", c.IP, c.Port),
		Username: c.Username,
		Password: c.Password,
	})
	
	if err != nil {
		return fmt.Errorf("failed to connect to camera: %v", err)
	}
	
	c.Device = dev
	return nil
}

// GetDevice returns the underlying ONVIF device
func (c *Camera) GetDevice() *onvif.Device {
	return c.Device
}

// DefaultCamera returns a Camera instance with default test credentials
func DefaultCamera() *Camera {
	return NewCamera("192.168.1.12", 80, "admin", "admin123")
}

// DefaultConnectedCamera returns a connected Camera instance with default test credentials
// or an error if the connection fails
func DefaultConnectedCamera() (*Camera, error) {
	camera := DefaultCamera()
	err := camera.Connect()
	if err != nil {
		return nil, err
	}
	return camera, nil
}