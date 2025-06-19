# CSV Camera Import Feature

## Overview

The CSV import feature allows you to bulk add multiple cameras to the system by uploading a CSV file with camera configuration data.

## API Endpoint

**POST** `/cameras/import-csv`

### Request Format

This endpoint accepts a multipart form upload with the following field:
- **csvFile**: The CSV file containing camera data

### CSV Format

The CSV file must include a header row with the following columns:

#### Required Columns
- **ip**: Camera IP address (required)
- **username**: Camera username (required)

#### Optional Columns
- **port**: Camera port (default: 80)
- **url**: Camera URL path (default: empty)
- **password**: Camera password (default: empty)
- **isfake**: Whether camera is simulated (true/false, default: false)

### Sample CSV

```csv
ip,port,url,username,password,isfake
192.168.1.10,80,,admin,admin123,false
192.168.1.11,8080,/onvif/device,user,password,false
192.168.1.12,80,,administrator,pass123,false
192.168.1.100,80,,fake_user,fake_pass,true
```

## Response Format

The endpoint returns a JSON response with the following structure:

```json
{
  "message": "CSV import completed: 3 cameras added successfully, 1 errors",
  "totalRows": 4,
  "successCount": 3,
  "errorCount": 1,
  "results": [
    {
      "row": 2,
      "success": true,
      "cameraId": "1",
      "camera": {
        "id": "1",
        "ip": "192.168.1.10",
        "port": 80,
        "url": "",
        "username": "admin",
        "password": "admin123",
        "isFake": false
      }
    },
    {
      "row": 3,
      "success": false,
      "error": "Camera with IP 192.168.1.11 already exists",
      "data": ["192.168.1.11", "8080", "/onvif/device", "user", "password", "false"]
    }
  ]
}
```

## HTTP Status Codes

- **201 Created**: All cameras added successfully
- **206 Partial Content**: Some cameras added, some failed
- **400 Bad Request**: No cameras added (all failed) or invalid CSV format

## Features

1. **Flexible Column Order**: Columns can be in any order as long as headers match
2. **Case Insensitive Headers**: Column names are case-insensitive
3. **Default Values**: Missing optional fields use sensible defaults
4. **Detailed Results**: Response includes success/failure status for each row
5. **Error Handling**: Invalid rows are skipped with error details
6. **Row Tracking**: Each result includes the original CSV row number

## Error Handling

Common errors include:
- Missing required columns (ip, username)
- Invalid port numbers
- Duplicate camera IP addresses
- Empty CSV files
- Invalid CSV format

Each error includes the row number and original data for easy identification.

## Usage Example

### cURL Example

```bash
curl -X POST "http://localhost:8090/cameras/import-csv" \
     -F "csvFile=@cameras.csv" \
     -H "Content-Type: multipart/form-data"
```

### JavaScript Example

```javascript
const formData = new FormData();
formData.append('csvFile', csvFile); // csvFile is a File object

const response = await fetch('/api/cameras/import-csv', {
  method: 'POST',
  body: formData
});

const result = await response.json();
console.log(result);
```

## Validation Rules

1. **IP Address**: Must be provided and non-empty
2. **Username**: Must be provided and non-empty
3. **Port**: Must be a valid integer (defaults to 80 if invalid)
4. **IsFake**: Accepts true/false, 1/0, yes/no (case insensitive)

## File Size Limit

The current file size limit is 10MB. This should accommodate thousands of camera entries.
