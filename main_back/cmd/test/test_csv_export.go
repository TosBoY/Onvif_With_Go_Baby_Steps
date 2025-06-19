package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func main() {
	// Test CSV export with sample validation data
	testCSVExport()
}

func testCSVExport() {
	log.Println("Testing CSV export functionality...")

	// Sample validation data (mimicking the structure from apply-config response)
	sampleValidation := map[string]interface{}{
		"camera1": map[string]interface{}{
			"isValid":        true,
			"expectedWidth":  1920.0,
			"expectedHeight": 1080.0,
			"expectedFPS":    30.0,
			"actualWidth":    1920.0,
			"actualHeight":   1080.0,
			"actualFPS":      29.97,
		},
		"camera2": map[string]interface{}{
			"isValid":        false,
			"expectedWidth":  1280.0,
			"expectedHeight": 720.0,
			"expectedFPS":    25.0,
			"actualWidth":    1920.0,
			"actualHeight":   1080.0,
			"actualFPS":      30.0,
		},
		"camera3": map[string]interface{}{
			"isValid":        true,
			"expectedWidth":  640.0,
			"expectedHeight": 480.0,
			"expectedFPS":    15.0,
			"actualWidth":    640.0,
			"actualHeight":   480.0,
			"actualFPS":      15.12,
		},
	}

	// Prepare request payload
	payload := map[string]interface{}{
		"validation": sampleValidation,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Fatalf("Failed to marshal payload: %v", err)
	}

	// Make request to CSV export endpoint
	resp, err := http.Post("http://localhost:8080/export-validation-csv", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check response
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Fatalf("Request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Read and display CSV content
	csvContent, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response: %v", err)
	}

	log.Println("CSV Export successful!")
	log.Println("Generated CSV content:")
	log.Println(strings.Repeat("=", 50))
	fmt.Print(string(csvContent))
	log.Println(strings.Repeat("=", 50))

	// Verify CSV structure
	lines := strings.Split(strings.TrimSpace(string(csvContent)), "\n")
	if len(lines) < 4 { // Header + 3 data rows
		log.Printf("Warning: Expected at least 4 lines, got %d", len(lines))
	}

	// Check header
	expectedHeader := "cam_id,result,reso_expected,reso_actual,fps_expected,fps_actual"
	if lines[0] != expectedHeader {
		log.Printf("Warning: Header mismatch. Expected: %s, Got: %s", expectedHeader, lines[0])
	} else {
		log.Println("✓ Header format is correct")
	}

	// Check data rows
	expectedRows := 3
	actualRows := len(lines) - 1 // Excluding header
	if actualRows == expectedRows {
		log.Printf("✓ Correct number of data rows: %d", actualRows)
	} else {
		log.Printf("Warning: Expected %d data rows, got %d", expectedRows, actualRows)
	}

	log.Println("CSV export test completed!")
}
