# CSV Import Feature - Frontend Integration

## Overview

The CSV import feature has been integrated into the "Add New Camera" dialog, allowing users to bulk import multiple cameras via CSV upload directly from the UI.

## Location

The CSV import functionality is located in the **Add New Camera** dialog, accessible via the "Add Camera" button in the main dashboard.

## Features

### 1. **File Upload Interface**
- Drag-and-drop style file input
- CSV file validation
- File name display
- Clear/reset functionality

### 2. **Import Process**
- Real-time upload progress
- Detailed import results
- Success/error/warning alerts
- Automatic camera list refresh

### 3. **Results Display**
- Success count and error count
- Detailed error messages
- Color-coded alerts (green/yellow/red)
- Partial success handling

## User Interface

### CSV Import Section
```
┌─────────────────────────────────────┐
│ Bulk Import from CSV                │
│ Upload a CSV file to add multiple   │
│ cameras at once...                  │
│                                     │
│ [Choose CSV File] filename.csv [Clear] │
│ [Import CSV]                        │
└─────────────────────────────────────┘
            OR ADD SINGLE CAMERA
┌─────────────────────────────────────┐
│ Camera IP Address                   │
│ Port (Optional)                     │
│ Username                           │
│ etc...                             │
└─────────────────────────────────────┘
```

### Import Results
- **Success**: Green alert with "X cameras added successfully"
- **Partial**: Yellow alert with "X successful, Y failed"
- **Error**: Red alert with error details

## CSV Format

Users need to provide a CSV file with the following structure:

```csv
ip,port,url,username,password,isfake
192.168.1.12,80,,admin,admin123,false
192.168.1.30,0,,admin,Admin123,false
192.168.1.31,80,,admin,Admin123,false
192.168.1.15,0,,admin,Admin123,true
```

### Required Columns
- **ip**: Camera IP address
- **username**: Camera username

### Optional Columns
- **port**: Camera port (default: 80)
- **url**: Camera URL path (default: empty)
- **password**: Camera password (default: empty)
- **isfake**: Fake camera flag (default: false)

## User Workflow

1. **Open Dialog**: Click "Add Camera" button
2. **Choose Import Method**: 
   - Use CSV import for bulk operations
   - Use manual form for single cameras
3. **Upload CSV**: Click "Choose CSV File" and select file
4. **Import**: Click "Import CSV" button
5. **Review Results**: Check success/error alerts
6. **Continue or Close**: Add more cameras or close dialog

## Error Handling

### File Validation
- Only CSV files accepted
- File size limit enforced by backend (10MB)
- Invalid file types show error message

### Import Errors
- Row-level error reporting
- Detailed error messages
- Original data preservation for debugging

### Network Errors
- Connection failure handling
- Timeout protection
- User-friendly error messages

## State Management

The feature adds the following state variables:
- `csvFile`: Currently selected CSV file
- `importingCsv`: Import process status
- `csvImportResult`: Import results data

## Integration Points

### API Integration
- Uses `/api/cameras/import-csv` endpoint
- Multipart form upload
- Proper error handling

### UI Integration
- Seamless integration with existing dialog
- Consistent styling with Material-UI theme
- Responsive design for all screen sizes

## Benefits

1. **Bulk Operations**: Import dozens of cameras at once
2. **Time Saving**: No need to add cameras individually
3. **Error Recovery**: Failed imports don't affect successful ones
4. **User Feedback**: Clear indication of what succeeded/failed
5. **Seamless Integration**: Part of existing workflow

## Future Enhancements

1. **Template Download**: Provide CSV template download
2. **Drag-and-Drop**: Direct file drop support
3. **Preview Mode**: Show import preview before execution
4. **Progress Indicator**: Real-time import progress
5. **Export Integration**: Round-trip CSV export/import workflow
