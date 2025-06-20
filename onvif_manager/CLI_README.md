# ONVIF Camera Management CLI

A command line interface for managing ONVIF cameras, applying configurations, and validating streams.

## Installation & Build

```bash
cd onvif_manager
go build -o cli.exe .\cmd\app\main.go
```

## Available Commands

### 1. List Cameras (`list`)
Lists all cameras currently configured in the system.

```bash
.\cli.exe list
```

**Output:**
- Shows camera ID, IP, Port, Username, URL, and whether it's a simulated camera
- Displays total count of cameras

### 2. Select Cameras (`select`)
Select cameras from a CSV file containing IP addresses.

```bash
.\cli.exe select cameras_to_select.csv
```

**CSV Format:**
```csv
ip
192.168.1.12
192.168.1.30
192.168.1.31
```

**Output:**
- Shows which cameras were matched/found
- Lists unmatched IPs (cameras not in system)
- Reports invalid rows

### 3. Apply Configuration (`config apply`)
Apply configuration to selected cameras using two CSV files:
1. Camera selection CSV (IPs to configure)
2. Configuration CSV (settings to apply)

```bash
.\cli.exe config apply cameras_to_configure.csv config_settings.csv
```

**Camera Selection CSV Format:**
```csv
ip
192.168.1.12
192.168.1.30
```

**Configuration CSV Format:**
```csv
width,height,fps,bitrate
1920,1080,30,4096
```

**Process:**
1. Loads cameras from selection CSV
2. Loads configuration from config CSV
3. Asks for confirmation before applying
4. Applies configuration to each camera
5. Validates the applied configuration
6. Offers to export validation results

### 4. Export Results (`export`)
Export the last validation results to a CSV file.

```bash
.\cli.exe export validation_results.csv
```

**Output CSV Format:**
```csv
cam_id,cam_ip,result,reso_expected,reso_actual,fps_expected,fps_actual
1,192.168.1.12,PASS,1920x1080,1920x1080,30,30.00
2,192.168.1.30,FAIL,1920x1080,1280x720,30,25.00
```

## CSV File Examples

### Camera Selection CSV (`cameras_to_configure.csv`)
```csv
ip
192.168.1.12
192.168.1.30
192.168.1.31
```

### Configuration CSV (`config_1080p.csv`)
```csv
width,height,fps,bitrate
1920,1080,30,4096
```

### Alternative Configuration CSV (`config_720p.csv`)
```csv
width,height,fps,bitrate
1280,720,25,2048
```

## Usage Examples

### Basic Workflow

1. **List available cameras:**
   ```bash
   .\cli.exe list
   ```

2. **Apply configuration to specific cameras:**
   ```bash
   .\cli.exe config apply cameras_to_configure.csv config_1080p.csv
   ```
   - The CLI will ask for confirmation before applying
   - It will show progress and results
   - You can choose to export results to CSV

3. **Export results separately (if needed):**
   ```bash
   .\cli.exe export my_validation_results.csv
   ```

### Advanced Usage

**Select cameras first to verify which will be configured:**
```bash
.\cli.exe select cameras_to_configure.csv
```

**Then apply configuration:**
```bash
.\cli.exe config apply cameras_to_configure.csv config_settings.csv
```

## Configuration Options

- **Width/Height**: Target resolution (e.g., 1920x1080, 1280x720)
- **FPS**: Target frame rate (e.g., 30, 25, 15)
- **Bitrate**: Target bitrate in kbps (optional - set to 0 or omit for auto)

## Validation Results

The CLI validates each configuration by:
1. Checking if the camera responds
2. Verifying the applied resolution matches expected
3. Checking FPS and bitrate (warnings if mismatched)
4. Marking results as PASS/FAIL/WARNING

**Result Types:**
- **PASS**: All settings applied and validated successfully
- **WARNING**: Settings applied but FPS/bitrate differs (resolution matches)
- **FAIL**: Configuration failed or resolution doesn't match

## Error Handling

The CLI provides detailed error messages for:
- Network connectivity issues
- ONVIF protocol errors
- Invalid CSV formats
- Missing cameras
- Configuration failures

## Notes

- The CLI works with both real and simulated cameras
- Simulated cameras (marked as `IsFake: true`) will show successful results without actual network calls
- All paths are relative to where the CLI executable is run
- Configuration files should be in the same directory as the executable
