package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// PiProxyClient handles communication with the Raspberry Pi ONVIF proxy
type PiProxyClient struct {
	BaseURL     string
	HTTPClient  *http.Client
	ConnectInfo ConnectInfo
}

// ConnectInfo contains connection information for the Pi proxy
type ConnectInfo struct {
	Address string `json:"address"`
	Port    int    `json:"port"`
}

// NewPiProxyClient creates a new client for the Pi proxy
func NewPiProxyClient(address string, port int) *PiProxyClient {
	return &PiProxyClient{
		BaseURL: fmt.Sprintf("http://%s:%d/api", address, port),
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		ConnectInfo: ConnectInfo{
			Address: address,
			Port:    port,
		},
	}
}

// GetCameraInfo retrieves camera information from the Pi proxy
func (c *PiProxyClient) GetCameraInfo() (map[string]interface{}, error) {
	log.Printf("Retrieving camera info from Pi proxy at %s/camera/info", c.BaseURL)
	resp, err := c.HTTPClient.Get(fmt.Sprintf("%s/camera/info", c.BaseURL))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Pi proxy: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Pi proxy returned error status %d: %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse camera info response: %v", err)
	}

	return result, nil
}

// GetResolutions retrieves available resolutions from the Pi proxy
func (c *PiProxyClient) GetResolutions(configToken, profileToken string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/camera/resolutions?configToken=%s&profileToken=%s",
		c.BaseURL, configToken, profileToken)

	log.Printf("Retrieving resolutions from Pi proxy: %s", url)
	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Pi proxy: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Pi proxy returned error status %d: %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse resolutions response: %v", err)
	}

	return result, nil
}

// ChangeResolution sends a request to change resolution via the Pi proxy
func (c *PiProxyClient) ChangeResolution(payload interface{}) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal resolution change payload: %v", err)
	}

	log.Printf("Changing resolution via Pi proxy: %s/camera/change-resolution", c.BaseURL)
	resp, err := c.HTTPClient.Post(
		fmt.Sprintf("%s/camera/change-resolution", c.BaseURL),
		"application/json",
		bytes.NewBuffer(jsonPayload),
	)
	if err != nil {
		return fmt.Errorf("failed to connect to Pi proxy: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Pi proxy returned error status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// GetStreamURL retrieves the streaming URL for a profile via the Pi proxy
func (c *PiProxyClient) GetStreamURL(profileToken string) (string, error) {
	url := fmt.Sprintf("%s/camera/stream-url?profileToken=%s", c.BaseURL, profileToken)

	log.Printf("Retrieving stream URL from Pi proxy: %s", url)
	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to connect to Pi proxy: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Pi proxy returned error status %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		StreamURL string `json:"streamUrl"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to parse stream URL response: %v", err)
	}

	return result.StreamURL, nil
}

// GetConfigDetails retrieves configuration details for a specific token via the Pi proxy
func (c *PiProxyClient) GetConfigDetails(configToken string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/camera/config?configToken=%s", c.BaseURL, configToken)

	log.Printf("Retrieving config details from Pi proxy: %s", url)
	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Pi proxy: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Pi proxy returned error status %d: %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse config details response: %v", err)
	}

	return result, nil
}

// GetDeviceInfo retrieves device information via the Pi proxy
func (c *PiProxyClient) GetDeviceInfo() (map[string]interface{}, error) {
	log.Printf("Retrieving device info from Pi proxy: %s/camera/device-info", c.BaseURL)
	resp, err := c.HTTPClient.Get(fmt.Sprintf("%s/camera/device-info", c.BaseURL))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Pi proxy: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Pi proxy returned error status %d: %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse device info response: %v", err)
	}

	return result, nil
}

// GetSystemStatus retrieves the Pi system status
func (c *PiProxyClient) GetSystemStatus() (map[string]interface{}, error) {
	log.Printf("Retrieving system status from Pi proxy: %s/system/status", c.BaseURL)
	resp, err := c.HTTPClient.Get(fmt.Sprintf("%s/system/status", c.BaseURL))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Pi proxy: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Pi proxy returned error status %d: %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse system status response: %v", err)
	}

	return result, nil
}

// IsConnected checks if we can connect to the Pi proxy
func (c *PiProxyClient) IsConnected() bool {
	_, err := c.GetSystemStatus()
	return err == nil
}
