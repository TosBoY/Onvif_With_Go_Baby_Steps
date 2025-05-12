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

	// Mock functionality
	IsMock       bool
	MockProfiles []Profile
	MockConfigs  []VideoEncoderConfig
}

// NewCamera creates a new Camera instance
func NewCamera(ip string, port int, username, password string) *Camera {
	return &Camera{
		IP:       ip,
		Port:     port,
		Username: username,
		Password: password,
		IsMock:   false,
	}
}

// Connect establishes a connection to the camera
func (c *Camera) Connect() error {
	if c.IsMock {
		return nil // Mock cameras don't need real connection
	}

	xaddr := fmt.Sprintf("%s:%d", c.IP, c.Port)
	fmt.Printf("Attempting to connect to camera at %s\n", xaddr)

	dev, err := onvif.NewDevice(onvif.DeviceParams{
		Xaddr:    xaddr,
		Username: c.Username,
		Password: c.Password,
	})

	if err != nil {
		return fmt.Errorf("failed to connect to camera: %v", err)
	}

	c.Device = dev
	fmt.Printf("Successfully connected to camera at %s\n", xaddr)
	return nil
}

// SetMockMode enables or disables mock mode for the camera
func (c *Camera) SetMockMode(mock bool) {
	c.IsMock = mock
}

// SetMockData sets the mock profile and configuration data
func (c *Camera) SetMockData(profiles []Profile, configs []VideoEncoderConfig) {
	c.MockProfiles = profiles
	c.MockConfigs = configs
}

// GetAllProfiles returns either real or mock profiles based on mock mode
func (c *Camera) GetAllProfiles() ([]Profile, error) {
	if c.IsMock {
		return c.MockProfiles, nil
	}
	return GetAllProfiles(c)
}

// GetAllVideoEncoderConfigurations returns either real or mock configurations based on mock mode
func (c *Camera) GetAllVideoEncoderConfigurations() ([]VideoEncoderConfig, error) {
	if c.IsMock {
		return c.MockConfigs, nil
	}
	return GetAllVideoEncoderConfigurations(c)
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
