# CSV Export Feature - Frontend Implementation

## Overview

The CSV export feature has been added to the ValidationResults component, allowing users to download validation results as a CSV file directly from the browser.

## Features

- **Export Button**: Located in the header of the validation results panel
- **Automatic Filename**: CSV files are named with timestamp (e.g., `validation_results_2025-06-19T10-30-45.csv`)
- **Real-time Data**: Always exports the current validation results displayed on screen
- **Loading State**: Shows loading indicator while generating/downloading CSV
- **Error Handling**: User-friendly error messages if export fails

## CSV Format

The exported CSV contains the following columns:
- `cam_id`: Camera identifier
- `cam_ip`: Camera IP address
- `result`: PASS, WARNING, or FAIL
- `reso_expected`: Expected resolution (e.g., "1920x1080")
- `reso_actual`: Actual detected resolution
- `fps_expected`: Expected frame rate
- `fps_actual`: Actual detected frame rate

### Result Values

- **PASS**: All parameters match expected values exactly
- **WARNING**: Resolution matches but FPS/bitrate differs (acceptable variance)
- **FAIL**: Resolution mismatch or camera connection failure

## Implementation Details

### Components Modified

1. **ValidationResults.jsx**
   - Added CSV export button in header
   - Integrated with API service
   - Loading states and error handling

2. **api.js**
   - Added `exportValidationCSV()` function
   - Handles blob download automatically

### API Integration

The frontend calls the backend endpoint `/export-validation-csv` with the current validation data and automatically triggers a file download.

### Usage

1. Run camera validation through the normal configuration process
2. View validation results in the ValidationResults component
3. Click "Export CSV" button in the validation results header
4. CSV file downloads automatically with timestamped filename

### Testing

Use the CSVExportTest component to test the export functionality with sample data:

```jsx
import CSVExportTest from '../components/CSVExportTest';

// Add to any page for testing
<CSVExportTest />
```

## Automatic Updates

The CSV export always reflects the current validation results - when new validation results come in from camera configuration, the export button will export the new data automatically.
