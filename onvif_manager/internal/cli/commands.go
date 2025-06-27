package cli

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "onvif-manager",
	Short: "ONVIF Camera Management CLI",
	Long:  `A command line interface for managing ONVIF cameras, applying configurations, and validating streams.`,
}

// webCmd represents the web server command
var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Start combined web application (frontend + API) on port 8090",
	Long:  `Start combined web application with both frontend and API endpoints on port 8090.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Print information about the web application
		fmt.Println("🌐 Starting ONVIF Manager Web Application")
		fmt.Println("📱 Frontend will be available at http://localhost:8090")
		fmt.Println("🔌 API endpoints will be available at http://localhost:8090/api")
		fmt.Println("")

		// Since we can't directly call webserver.StartWebServer due to import cycle,
		// we'll just exit here. The actual handling is done in webserver.StartApp
		os.Exit(0)
	},
}

// serverCmd represents the API server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start API server only on port 8090",
	Long:  `Start API server only without the web interface on port 8090.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Print information about the API server
		fmt.Println("🚀 Starting ONVIF Manager API Server")
		fmt.Println("📊 API endpoints will be available at http://localhost:8090")
		fmt.Println("")

		// Since we can't directly call webserver.StartAPIServer due to import cycle,
		// we'll just exit here. The actual handling is done in webserver.StartApp
		os.Exit(0)
	},
}

var cameraService *CameraService

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return RootCmd.Execute()
}

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import [csv-file]",
	Short: "Import cameras from CSV file",
	Long:  `Import cameras into the system from a CSV file containing camera details.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runImportCameras(args[0])
	},
}

func init() {
	cameraService = NewCameraService()
	// Register only the simplified workflow commands
	RootCmd.AddCommand(webCmd)    // Add web command
	RootCmd.AddCommand(serverCmd) // Add server command
	RootCmd.AddCommand(configCmd) // Keep config command

	// These commands are no longer exposed in the simplified workflow
	// RootCmd.AddCommand(listCmd)
	// RootCmd.AddCommand(selectCmd)
	// RootCmd.AddCommand(exportCmd)
	// RootCmd.AddCommand(importCmd)
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all cameras in the system",
	Long:  `Display a list of all cameras currently configured in the system.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runListCameras()
	},
}

// selectCmd represents the select command
var selectCmd = &cobra.Command{
	Use:   "select [csv-file]",
	Short: "Select cameras from CSV file",
	Long:  `Select cameras for configuration by providing a CSV file containing IP addresses.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runSelectCameras(args[0])
	},
}

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Apply configuration to cameras",
	Long:  `Apply configuration settings to selected cameras using CSV files.`,
}

var applyConfigCmd = &cobra.Command{
	Use:   "apply [camera-csv] [config-csv]",
	Short: "Import cameras and apply configuration",
	Long:  `Import cameras from first CSV file and apply configuration from second CSV file in a single operation.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runApplyConfig(args[0], args[1])
	},
}

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export [output-file]",
	Short: "Export validation results to CSV",
	Long:  `Export the last validation results to a CSV file.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runExportResults(args[0])
	},
}

// configShowCmd shows the current saved configuration
var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current saved configuration",
	Long:  `Display the current saved configuration values.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runShowConfig()
	},
}

// configSetCmd sets the saved configuration manually
var configSetCmd = &cobra.Command{
	Use:   "set [width] [height] [fps] [bitrate]",
	Short: "Set saved configuration manually",
	Long:  `Set the saved configuration values manually by providing width, height, fps, and bitrate.`,
	Args:  cobra.ExactArgs(4),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runSetConfig(args[0], args[1], args[2], args[3])
	},
}

// configImportCmd imports configuration from CSV and saves it
var configImportCmd = &cobra.Command{
	Use:   "import [config-csv]",
	Short: "Import configuration from CSV file",
	Long:  `Import configuration from CSV file and save it as the current saved configuration.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runImportConfig(args[0])
	},
}

// applyToSelectedCmd applies saved config to selected cameras
var applyToSelectedCmd = &cobra.Command{
	Use:   "apply-to [camera-csv]",
	Short: "Apply saved configuration to selected cameras",
	Long:  `Apply the current saved configuration to cameras selected from CSV file.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runApplyToSelected(args[0])
	},
}

func init() {
	// Only add the apply command for the simplified workflow
	configCmd.AddCommand(applyConfigCmd)

	// These commands are no longer exposed in the simplified workflow
	// but kept in the code for compatibility with existing scripts
	// configCmd.AddCommand(configShowCmd)
	// configCmd.AddCommand(configSetCmd)
	// configCmd.AddCommand(configImportCmd)
	// configCmd.AddCommand(applyToSelectedCmd)
}

// Global variable to store last validation results
var lastValidationResults *ValidationResults

// runListCameras lists all cameras in the system
func runListCameras() error {
	fmt.Println("📹 Loading camera list...")

	cameras, err := cameraService.GetCameraList()
	if err != nil {
		return fmt.Errorf("failed to load camera list: %w", err)
	}

	if len(cameras) == 0 {
		fmt.Println("No cameras found in the system.")
		return nil
	}

	fmt.Printf("\n📋 Found %d camera(s):\n\n", len(cameras))
	fmt.Printf("%-15s %-15s %-8s %-20s %-12s\n", "ID", "IP", "Port", "Username", "URL")
	fmt.Println(strings.Repeat("-", 70))

	for _, camera := range cameras {
		fmt.Printf("%-15s %-15s %-8d %-20s %-12s\n",
			camera.ID, camera.IP, camera.Port, camera.Username, camera.URL)
	}

	fmt.Printf("\nTotal: %d cameras\n", len(cameras))
	return nil
}

// runSelectCameras selects cameras from CSV file
func runSelectCameras(csvFile string) error {
	fmt.Printf("📂 Loading camera selection from: %s\n", csvFile)

	result, err := cameraService.SelectCamerasFromCSV(csvFile)
	if err != nil {
		return fmt.Errorf("failed to select cameras: %w", err)
	}

	fmt.Printf("\n✅ %s\n", result.Message)
	fmt.Printf("📊 Selection Summary:\n")
	fmt.Printf("   • Total rows processed: %d\n", result.TotalRows)
	fmt.Printf("   • Cameras matched: %d\n", result.MatchedCount)
	fmt.Printf("   • IPs not found: %d\n", result.UnmatchedCount)
	fmt.Printf("   • Invalid rows: %d\n", result.InvalidRowCount)

	if len(result.SelectedCameras) > 0 {
		fmt.Printf("\n🎯 Selected cameras:\n")
		for _, camera := range result.SelectedCameras {
			fmt.Printf("   • %s - %s\n", camera.ID, camera.IP)
		}
	}

	if len(result.UnmatchedIPs) > 0 {
		fmt.Printf("\n⚠️  IPs not found in system:\n")
		for _, ip := range result.UnmatchedIPs {
			fmt.Printf("   • %s\n", ip)
		}
	}

	if len(result.InvalidRows) > 0 {
		fmt.Printf("\n❌ Invalid rows:\n")
		for _, row := range result.InvalidRows {
			fmt.Printf("   • Row %d: %s\n", row.Row, row.Error)
		}
	}

	return nil
}

// runApplyConfig imports cameras and applies configuration in one workflow
func runApplyConfig(cameraCSV, configCSV string) error {
	// Step 1: Import cameras from the first CSV file
	fmt.Printf("📂 Importing cameras from: %s\n", cameraCSV)
	importResult, err := cameraService.ImportCamerasFromCSV(cameraCSV)
	if err != nil {
		return fmt.Errorf("failed to import cameras: %w", err)
	}

	if importResult.SuccessCount == 0 {
		return fmt.Errorf("no cameras were successfully imported from CSV file")
	}

	fmt.Printf("✅ Imported %d cameras\n", importResult.SuccessCount)

	// Extract camera IDs from import result
	var cameraIDs []string
	for _, result := range importResult.Results {
		if result.Success {
			cameraIDs = append(cameraIDs, result.CameraID)
		}
	}

	// Step 2: Load configuration
	fmt.Printf("📂 Loading configuration from: %s\n", configCSV)
	config, err := cameraService.ImportConfigFromCSV(configCSV)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	fmt.Printf("⚙️  Configuration loaded: %dx%d, %d FPS, %d kbps\n",
		config.Width, config.Height, config.FPS, config.Bitrate)
	// Step 3: Confirm with user
	fmt.Printf("\n🤔 Do you want to apply this configuration to %d cameras? (y/N): ", len(cameraIDs))
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	response := strings.ToLower(strings.TrimSpace(scanner.Text()))

	if response != "y" && response != "yes" {
		fmt.Println("❌ Configuration application cancelled.")
		return nil
	}
	// Step 4: Apply configuration
	fmt.Printf("\n🔧 Applying configuration to cameras...\n")

	// Note: No need to call EnsureCamerasInitialized as cameras are already initialized during import

	validation, err := cameraService.ApplyConfigToCameras(cameraIDs, config)
	if err != nil {
		return fmt.Errorf("failed to apply configuration: %w", err)
	}

	// Store results for potential export
	lastValidationResults = validation

	// Step 5: Display results
	fmt.Printf("\n📊 Configuration Results:\n")

	successCount := 0
	failureCount := 0
	validationPassCount := 0
	validationFailCount := 0

	for cameraID, result := range validation.CameraResults {
		status := "❌ FAILED"
		if result.Success {
			status = "✅ SUCCESS"
			successCount++
		} else {
			failureCount++
		}

		fmt.Printf("   • Camera %s: %s", cameraID, status)
		if !result.Success && result.Error != nil {
			fmt.Printf(" - %s", result.Error.Error())
		}
		fmt.Println()

		if validationResult, exists := validation.ValidationResults[cameraID]; exists {
			if validationResult.IsValid {
				fmt.Printf("      Validation: ✅ PASSED\n")
				validationPassCount++
			} else {
				fmt.Printf("      Validation: ❌ FAILED - %s\n", validationResult.Error)
				validationFailCount++
			}
		}
	}

	fmt.Printf("\n📈 Summary:\n")
	fmt.Printf("   • Configuration: %d success, %d failed\n", successCount, failureCount)
	fmt.Printf("   • Validation: %d passed, %d failed\n", validationPassCount, validationFailCount)

	// Step 6: Offer to export results
	if lastValidationResults != nil {
		fmt.Printf("\n💾 Do you want to export validation results to CSV? (y/N): ")
		scanner.Scan()
		response = strings.ToLower(strings.TrimSpace(scanner.Text()))
		if response == "y" || response == "yes" {
			defaultFilename := generateTimestampedFilename("validation_results.csv")
			fmt.Printf("📄 Enter output filename (default: %s): ", defaultFilename)
			scanner.Scan()
			filename := strings.TrimSpace(scanner.Text())
			if filename == "" {
				filename = defaultFilename
			}

			return runExportResults(filename)
		}
	}

	return nil
}

// runExportResults exports validation results to CSV
func runExportResults(outputFile string) error {
	if lastValidationResults == nil {
		return fmt.Errorf("no validation results available to export. Run 'config apply' first")
	}

	fmt.Printf("💾 Exporting validation results to: %s\n", outputFile)

	err := cameraService.ExportValidationToCSV(lastValidationResults, outputFile)
	if err != nil {
		return fmt.Errorf("failed to export results: %w", err)
	}

	fmt.Printf("✅ Results exported successfully to %s\n", outputFile)
	return nil
}

// runShowConfig displays the current saved configuration
func runShowConfig() error {
	configService := NewConfigService()
	config, err := configService.LoadSavedConfig()
	if err != nil {
		return fmt.Errorf("failed to load saved configuration: %w", err)
	}

	fmt.Printf("\n📂 Current Saved Configuration:\n")
	fmt.Printf("   • Resolution: %dx%d\n", config.Width, config.Height)
	fmt.Printf("   • Frame Rate: %d FPS\n", config.FPS)
	fmt.Printf("   • Bitrate: %d kbps\n", config.Bitrate)
	fmt.Printf("   • Encoding: %s\n", config.Encoding)
	fmt.Printf("   • Last Updated: %s\n", config.LastUpdated)
	fmt.Printf("   • Source: %s\n", config.Source)

	return nil
}

// runSetConfig sets the saved configuration manually
func runSetConfig(widthStr, heightStr, fpsStr, bitrateStr string) error {
	return runSetConfigWithEncoding(widthStr, heightStr, fpsStr, bitrateStr, "H264")
}

// runSetConfigWithEncoding sets the saved configuration manually with encoding
func runSetConfigWithEncoding(widthStr, heightStr, fpsStr, bitrateStr, encoding string) error {
	width, err := strconv.Atoi(widthStr)
	if err != nil {
		return fmt.Errorf("invalid width value: %w", err)
	}

	height, err := strconv.Atoi(heightStr)
	if err != nil {
		return fmt.Errorf("invalid height value: %w", err)
	}

	fps, err := strconv.Atoi(fpsStr)
	if err != nil {
		return fmt.Errorf("invalid fps value: %w", err)
	}

	bitrate, err := strconv.Atoi(bitrateStr)
	if err != nil {
		return fmt.Errorf("invalid bitrate value: %w", err)
	}

	// Use default encoding if empty
	if encoding == "" {
		encoding = "H264"
	}

	configService := NewConfigService()
	if err := configService.ValidateConfig(width, height, fps, bitrate, encoding); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	err = configService.UpdateManually(width, height, fps, bitrate, encoding)
	if err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}
	fmt.Printf("✅ Configuration updated successfully:\n")
	fmt.Printf("   • Resolution: %dx%d\n", width, height)
	fmt.Printf("   • Frame Rate: %d FPS\n", fps)
	fmt.Printf("   • Bitrate: %d kbps\n", bitrate)
	fmt.Printf("   • Encoding: %s\n", encoding)
	return nil
}

// runImportConfig imports configuration from CSV and saves it
func runImportConfig(configCSV string) error {
	fmt.Printf("📂 Importing configuration from: %s\n", configCSV)

	configService := NewConfigService()
	configData, err := cameraService.ImportConfigFromCSV(configCSV)
	if err != nil {
		return fmt.Errorf("failed to import configuration: %w", err)
	}

	fmt.Printf("⚙️  Configuration imported: %dx%d, %d FPS, %d kbps\n",
		configData.Width, configData.Height, configData.FPS, configData.Bitrate)

	// Save as current configuration
	err = configService.ImportFromConfigData(configData, "csv")
	if err != nil {
		return fmt.Errorf("failed to save imported configuration: %w", err)
	}

	fmt.Println("✅ Configuration imported and saved successfully.")
	return nil
}

// runApplyToSelected applies saved config to selected cameras
func runApplyToSelected(cameraCSV string) error {
	// Step 1: Select cameras
	fmt.Printf("📂 Loading camera selection from: %s\n", cameraCSV)
	selection, err := cameraService.SelectCamerasFromCSV(cameraCSV)
	if err != nil {
		return fmt.Errorf("failed to select cameras: %w", err)
	}

	if len(selection.SelectedCameraIDs) == 0 {
		return fmt.Errorf("no cameras selected from CSV file")
	}

	fmt.Printf("✅ Selected %d cameras\n", len(selection.SelectedCameraIDs))

	// Step 2: Get current saved configuration
	configService := NewConfigService()
	savedConfig, err := configService.LoadSavedConfig()
	if err != nil {
		return fmt.Errorf("failed to load saved configuration: %w", err)
	}

	fmt.Printf("⚙️  Current saved configuration: %dx%d, %d FPS, %d kbps\n",
		savedConfig.Width, savedConfig.Height, savedConfig.FPS, savedConfig.Bitrate)

	// Step 3: Confirm with user
	fmt.Printf("\n🤔 Do you want to apply this configuration to %d cameras? (y/N): ", len(selection.SelectedCameraIDs))
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	response := strings.ToLower(strings.TrimSpace(scanner.Text()))

	if response != "y" && response != "yes" {
		fmt.Println("❌ Configuration application cancelled.")
		return nil
	}
	// Step 4: Apply configuration
	fmt.Printf("\n🔧 Applying configuration to cameras...\n")

	// Ensure cameras are initialized
	if err := cameraService.EnsureCamerasInitialized(); err != nil {
		return fmt.Errorf("failed to initialize cameras: %w", err)
	}

	validation, err := cameraService.ApplyConfigToCameras(selection.SelectedCameraIDs, savedConfig.ToConfigData())
	if err != nil {
		return fmt.Errorf("failed to apply configuration: %w", err)
	}

	// Store results for potential export
	lastValidationResults = validation

	// Step 5: Display results
	fmt.Printf("\n📊 Configuration Results:\n")

	successCount := 0
	failureCount := 0
	validationPassCount := 0
	validationFailCount := 0
	for cameraID, result := range validation.CameraResults {
		status := "❌ FAILED"
		if result.Success {
			status = "✅ SUCCESS"
			successCount++
		} else {
			failureCount++
		}

		fmt.Printf("   • Camera %s: %s", cameraID, status)
		if !result.Success && result.Error != nil {
			fmt.Printf(" - %s", result.Error.Error())
		}
		fmt.Println()

		if validationResult, exists := validation.ValidationResults[cameraID]; exists {
			if validationResult.IsValid {
				fmt.Printf("      Validation: ✅ PASSED\n")
				validationPassCount++
			} else {
				fmt.Printf("      Validation: ❌ FAILED - %s\n", validationResult.Error)
				validationFailCount++
			}
		}
	}

	fmt.Printf("\n📈 Summary:\n")
	fmt.Printf("   • Configuration: %d success, %d failed\n", successCount, failureCount)
	fmt.Printf("   • Validation: %d passed, %d failed\n", validationPassCount, validationFailCount)

	// Step 6: Offer to export results
	if lastValidationResults != nil {
		fmt.Printf("\n💾 Do you want to export validation results to CSV? (y/N): ")
		scanner.Scan()
		response = strings.ToLower(strings.TrimSpace(scanner.Text()))
		if response == "y" || response == "yes" {
			defaultFilename := generateTimestampedFilename("validation_results.csv")
			fmt.Printf("📄 Enter output filename (default: %s): ", defaultFilename)
			scanner.Scan()
			filename := strings.TrimSpace(scanner.Text())
			if filename == "" {
				filename = defaultFilename
			}

			return runExportResults(filename)
		}
	}

	return nil
}

// generateTimestampedFilename creates a filename with timestamp
func generateTimestampedFilename(baseName string) string {
	timestamp := time.Now().Format("20060102_150405")
	if strings.Contains(baseName, ".") {
		parts := strings.Split(baseName, ".")
		extension := parts[len(parts)-1]
		nameWithoutExt := strings.Join(parts[:len(parts)-1], ".")
		return fmt.Sprintf("%s_%s.%s", nameWithoutExt, timestamp, extension)
	}
	return fmt.Sprintf("%s_%s", baseName, timestamp)
}

// askForConfirmation asks user for yes/no confirmation
func askForConfirmation(message string) bool {
	fmt.Printf("%s (y/N): ", message)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	response := strings.ToLower(strings.TrimSpace(scanner.Text()))
	return response == "y" || response == "yes"
}

// runImportCameras imports cameras from a CSV file
func runImportCameras(csvFile string) error {
	fmt.Printf("📂 Importing cameras from: %s\n", csvFile)

	result, err := cameraService.ImportCamerasFromCSV(csvFile)
	if err != nil {
		return fmt.Errorf("failed to import cameras: %w", err)
	}

	fmt.Printf("\n✅ %s\n", result.Message)
	fmt.Printf("📊 Import Summary:\n")
	fmt.Printf("   • Total rows processed: %d\n", result.TotalRows)
	fmt.Printf("   • Cameras added successfully: %d\n", result.SuccessCount)
	fmt.Printf("   • Errors: %d\n", result.ErrorCount)

	if result.SuccessCount > 0 {
		fmt.Printf("\n✨ Successfully added cameras:\n")
		for _, rowResult := range result.Results {
			if rowResult.Success {
				fmt.Printf("   • Camera ID: %s - %s\n", rowResult.CameraID, rowResult.Camera.IP)
			}
		}
	}

	if result.ErrorCount > 0 {
		fmt.Printf("\n⚠️ Import errors:\n")
		for _, rowResult := range result.Results {
			if !rowResult.Success {
				fmt.Printf("   • Row %d: %s\n", rowResult.Row, rowResult.Error)
			}
		}
	}

	return nil
}
