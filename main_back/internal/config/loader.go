package config

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"main_back/pkg/models"
)

// LoadCameraList reads a CSV file containing camera configurations.
// It determines the path to the config file relative to the executable's location.
func LoadCameraList() ([]models.Camera, error) {
	workingDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get working directory: %w", err)
	}

	configPath := filepath.Join(workingDir, "..", "..", "config", "cameras.csv")

	// Open the CSV file
	file, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open camera config file %s: %w", configPath, err)
	}
	defer file.Close()

	// Create a new CSV reader
	csvReader := csv.NewReader(file)

	// Read the header row first
	header, err := csvReader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV header: %w", err)
	}

	// Create a map of column indices for each field
	columnMap := make(map[string]int)
	for i, column := range header {
		columnMap[strings.TrimSpace(strings.ToLower(column))] = i
	}

	// Verify all required columns exist
	requiredColumns := []string{"id", "ip", "port", "url", "username", "password", "isfake"}
	for _, col := range requiredColumns {
		if _, exists := columnMap[col]; !exists {
			return nil, fmt.Errorf("required column '%s' not found in CSV header", col)
		}
	}

	// Read all records and convert to cameras
	var cameras []models.Camera
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading CSV record: %w", err)
		}

		// Parse port as integer
		port, err := strconv.Atoi(record[columnMap["port"]])
		if err != nil {
			return nil, fmt.Errorf("invalid port value for camera ID %s: %w", record[columnMap["id"]], err)
		}

		// Parse isFake as boolean
		isFake, err := strconv.ParseBool(record[columnMap["isfake"]])
		if err != nil {
			return nil, fmt.Errorf("invalid isFake value for camera ID %s: %w", record[columnMap["id"]], err)
		}

		camera := models.Camera{
			ID:       record[columnMap["id"]],
			IP:       record[columnMap["ip"]],
			Port:     port,
			URL:      record[columnMap["url"]],
			Username: record[columnMap["username"]],
			Password: record[columnMap["password"]],
			IsFake:   isFake,
		}
		cameras = append(cameras, camera)
	}

	return cameras, nil
}
