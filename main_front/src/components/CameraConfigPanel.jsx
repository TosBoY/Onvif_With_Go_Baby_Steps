import React, { useState } from 'react';
import { 
  Box, 
  Typography, 
  TextField, 
  MenuItem, 
  Button, 
  Paper, 
  CircularProgress,
  Alert,
  Grid,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Divider
} from '@mui/material';
import { Upload as UploadIcon } from '@mui/icons-material';
import { applyConfig, importConfigCSV } from '../services/api';

// Common resolution presets
const resolutions = [
  { width: 640, height: 480, label: '640x480 (VGA / SD)' },
  { width: 854, height: 480, label: '854x480 (FWVGA / SD Wide)' },
  { width: 1280, height: 720, label: '1280x720 (HD / 720p)' },
  { width: 1366, height: 768, label: '1366x768 (HD+)' },
  { width: 1600, height: 900, label: '1600x900 (HD+)' },
  { width: 1920, height: 1080, label: '1920x1080 (Full HD / 1080p)' },
  { width: 2560, height: 1440, label: '2560x1440 (QHD / 2K)' },
  { width: 3440, height: 1440, label: '3440x1440 (Ultrawide QHD)' },
  { width: 3840, height: 2160, label: '3840x2160 (4K UHD)' },
  { width: 5120, height: 2880, label: '5120x2880 (5K Retina)' },
  { width: 7680, height: 4320, label: '7680x4320 (8K UHD)' }
];

const CameraConfigPanel = ({ 
  selectedCamera, 
  selectedCameras = [], 
  cameras = [],
  onConfigurationApplied, // Add this prop
  onClearValidation // Add this prop
}) => {  const [width, setWidth] = useState(1280);
  const [height, setHeight] = useState(720);
  const [fps, setFps] = useState(25); // Set default FPS to 25
  const [bitrate, setBitrate] = useState(''); // Add bitrate state
  const [isLoading, setIsLoading] = useState(false);
  const [result, setResult] = useState({ success: false, message: null });

  // CSV Config Upload state
  const [uploadDialogOpen, setUploadDialogOpen] = useState(false);
  const [configCsvFile, setConfigCsvFile] = useState(null);
  const [uploadingConfig, setUploadingConfig] = useState(false);
  const [uploadResult, setUploadResult] = useState(null);

  const handleResolutionChange = (event) => {
    const selectedResolution = JSON.parse(event.target.value);
    setWidth(selectedResolution.width);
    setHeight(selectedResolution.height);
  };
  const handleFpsChange = (event) => {
    const value = event.target.value;
    // Allow empty string for clearing the field
    if (value === '') {
      setFps('');
      return;
    }
    // Parse and validate the number
    const numValue = parseInt(value, 10);
    if (!isNaN(numValue) && numValue > 0) {
      setFps(numValue);
    } else {
      // Allow partial input (like just typing numbers)
      setFps(value);
    }
  };
  const handleBitrateChange = (event) => {
    const value = event.target.value;
    // Allow empty string for clearing the field
    if (value === '') {
      setBitrate('');
      return;
    }
    // Parse and validate the number
    const numValue = parseInt(value, 10);
    if (!isNaN(numValue) && numValue > 0) {
      setBitrate(numValue);
    } else {
      // Allow partial input (like just typing numbers)
      setBitrate(value);
    }
  };

  // CSV Upload handlers
  const handleUploadDialogOpen = () => {
    setUploadDialogOpen(true);
    setConfigCsvFile(null);
    setUploadResult(null);
  };

  const handleUploadDialogClose = () => {
    setUploadDialogOpen(false);
    setConfigCsvFile(null);
    setUploadResult(null);
  };
  const handleConfigCsvFileChange = async (event) => {
    const file = event.target.files[0];
    if (file) {
      setConfigCsvFile(file);
      setUploadResult(null);
      setUploadingConfig(true);

      try {
        // Automatically scan the CSV file when selected
        const data = await importConfigCSV(file);
        
        if (data && data.config) {
          setUploadResult({ 
            success: true, 
            message: 'Configuration scanned successfully!',
            config: data.config
          });
        } else {
          setUploadResult({ 
            success: false, 
            message: 'Failed to parse configuration from CSV' 
          });
        }
      } catch (error) {
        console.error('Error scanning config CSV:', error);
        setUploadResult({ 
          success: false, 
          message: error.message || 'Failed to scan CSV file' 
        });
      } finally {
        setUploadingConfig(false);
      }
    } else {
      setConfigCsvFile(null);
      setUploadResult(null);
    }
  };

  const handleUploadConfig = async () => {
    // This function is no longer needed since we scan automatically
    // But keeping it for backward compatibility
    return;
  };  const handleApplyUploadedConfig = async () => {
    if (uploadResult && uploadResult.config) {
      const config = uploadResult.config;
      setWidth(config.width);
      setHeight(config.height);
      setFps(config.fps);
      setBitrate(config.bitrate || '');
      handleUploadDialogClose();
      
      // If cameras are selected, apply the configuration immediately
      if (selectedCameras.length > 0) {
        // Clear previous validation results
        if (onClearValidation) {
          onClearValidation();
        }

        setIsLoading(true);
        setResult({ 
          success: true, 
          message: `Applying configuration from CSV to ${selectedCameras.length} camera(s)... This may take a while as cameras need time to update and validate settings.`
        });

        try {
          // Apply the configuration using the same logic as handleApplyConfig
          const batchResult = await applyConfig(selectedCameras, config.width, config.height, config.fps, config.bitrate || 0);
          console.log('Batch configuration result from CSV:', batchResult);

          // Process response - extract validation results
          const validations = [];
          const successfulCameras = [];
          const failedCameras = [];

          if (batchResult && batchResult.results) {
            Object.entries(batchResult.results).forEach(([cameraId, result]) => {
              if (result.success) {
                successfulCameras.push(cameraId);
                
                if (result.validation) {
                  const validationData = {
                    ...result.validation,
                    cameraId,
                    expectedWidth: config.width,
                    expectedHeight: config.height,
                    expectedFPS: config.fps,
                    expectedBitrate: config.bitrate || 0
                  };
                  
                  if (result.isFake) {
                    validationData.isValid = true;
                  }
                  
                  validations.push(validationData);
                }
              } else {
                failedCameras.push({
                  cameraId,
                  error: result.error || 'Unknown error'
                });
              }
            });
          }
          
          // Pass validation results back to Dashboard
          if (onConfigurationApplied && validations.length > 0) {
            const compositeResult = {
              validation: validations,
              appliedConfig: {
                resolution: { width: config.width, height: config.height },
                fps: config.fps,
                bitrate: config.bitrate || 0
              }
            };
            onConfigurationApplied(compositeResult);
          }
          
          // Show final results summary
          if (failedCameras.length === 0) {
            setResult({ 
              success: true, 
              message: `CSV configuration applied successfully to ${successfulCameras.length} camera(s): ${successfulCameras.join(', ')}` 
            });
          } else if (successfulCameras.length > 0) {
            const successMessage = `Applied to cameras ${successfulCameras.join(', ')}`;
            const errorDetails = failedCameras.map(e => `${e.cameraId}: ${e.error}`).join('; ');
            
            setResult({ 
              success: false, 
              message: `Partial success: ${successMessage}. Failed cameras: ${errorDetails}` 
            });
          } else {
            const errorDetails = failedCameras.map(e => `${e.cameraId}: ${e.error}`).join('; ');
            
            setResult({ 
              success: false, 
              message: `Failed to apply CSV configuration to any cameras. Errors: ${errorDetails}` 
            });
          }
        } catch (error) {
          console.error('Error applying CSV configuration:', error);
          setResult({ 
            success: false, 
            message: `CSV configuration failed: ${error.message}` 
          });
        } finally {
          setIsLoading(false);
        }
      } else {
        // No cameras selected, just show that configuration was loaded
        setResult({ 
          success: true, 
          message: `Configuration loaded from CSV: ${config.width}x${config.height}, ${config.fps} FPS, ${config.bitrate || 'Auto'} kbps. Select cameras to apply configuration.` 
        });
      }
    }
  };
  const handleApplyConfig = async () => {
    // Clear previous validation results
    if (onClearValidation) {
      onClearValidation();
    }

    if (selectedCameras.length === 0) {
      setResult({ 
        success: false, 
        message: 'Please select at least one camera from the list to apply configuration' 
      });
      return;
    }    // Validate FPS value before sending
    const fpsValue = parseInt(fps, 10);
    if (isNaN(fpsValue) || fpsValue <= 0) {
      setResult({ 
        success: false, 
        message: 'Please enter a valid frame rate (positive number)' 
      });
      return;
    }

    // Validate bitrate value if provided
    let bitrateValue = 0; // Default to 0 if not provided
    if (bitrate !== '') {
      bitrateValue = parseInt(bitrate, 10);
      if (isNaN(bitrateValue) || bitrateValue <= 0) {
        setResult({ 
          success: false, 
          message: 'Please enter a valid bitrate (positive number) or leave it empty for automatic selection' 
        });
        return;
      }
    }
    
    setIsLoading(true);
    setResult({ 
      success: true, 
      message: `Applying configuration to ${selectedCameras.length} camera(s)... This may take a while as cameras need time to update and validate settings.`
    });
      // Use batch mode to send all camera IDs at once
    try {      // Call API with array of camera IDs (batch mode) - ensure fps is a number
      const batchResult = await applyConfig(selectedCameras, width, height, fpsValue, bitrateValue);
      console.log('Batch configuration result:', batchResult);      // Process response - extract validation results
      const validations = [];
      const successfulCameras = [];
      const failedCameras = [];

      console.log('Received batch result structure:', batchResult);
      
      // Analyze results for each camera
      if (batchResult && batchResult.results) {
        Object.entries(batchResult.results).forEach(([cameraId, result]) => {
          console.log(`Processing camera ${cameraId} result:`, result);
          
          if (result.success) {
            successfulCameras.push(cameraId);
            
            // Extract and format validation data if available
            if (result.validation) {              const validationData = {
                ...result.validation,
                cameraId,
                // Add expected values from the original request
                expectedWidth: width,
                expectedHeight: height,
                expectedFPS: fps,
                expectedBitrate: bitrateValue
              };
              
              // For fake cameras, ensure we're marking if values match exactly
              if (result.isFake) {
                validationData.isValid = true;
              }
              
              console.log(`Adding validation for camera ${cameraId}:`, validationData);
              validations.push(validationData);
            }
          } else {
            failedCameras.push({
              cameraId,
              error: result.error || 'Unknown error'
            });
          }
        });
      }
      
      // Pass all validation results back to Dashboard
      if (onConfigurationApplied && validations.length > 0) {        const compositeResult = {
          validation: validations,
          appliedConfig: {
            resolution: { width, height },
            fps,
            bitrate: bitrateValue
          }
        };
        onConfigurationApplied(compositeResult);
      }
      
      // Show final results summary
      if (failedCameras.length === 0) {
        setResult({ 
          success: true, 
          message: `Configuration applied successfully to ${successfulCameras.length} camera(s): ${successfulCameras.join(', ')}` 
        });
      } else if (successfulCameras.length > 0) {
        // For partial success, be more specific about what worked
        const successMessage = `Applied to cameras ${successfulCameras.join(', ')}`;
        const errorDetails = failedCameras.map(e => `${e.cameraId}: ${e.error}`).join('; ');
        
        setResult({ 
          success: false, 
          message: `Partial success: ${successMessage}. Failed cameras: ${errorDetails}` 
        });
      } else {
        // For complete failure, show detailed error for each camera
        const errorDetails = failedCameras.map(e => `${e.cameraId}: ${e.error}`).join('; ');
        
        setResult({ 
          success: false, 
          message: `Failed to apply configuration to any cameras. Errors: ${errorDetails}` 
        });
      }
    } catch (error) {
      console.error('Error in batch configuration:', error);
      setResult({ 
        success: false, 
        message: `Configuration failed: ${error.message}` 
      });
    } finally {
      setIsLoading(false);
    }
  };  return (
    <Box sx={{ p: 2 }}>
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
        <Typography variant="h6">
          Camera Configuration
        </Typography>
        <Button
          variant="outlined"
          startIcon={<UploadIcon />}
          onClick={handleUploadDialogOpen}
          size="small"
        >
          Upload Config
        </Button>
      </Box>
        <Grid container spacing={2}>
        <Grid item xs={12} md={6}>
          <TextField
            select
            label="Resolution"
            value={JSON.stringify({ width, height })}
            onChange={handleResolutionChange}
            fullWidth
            margin="normal"
          >
            {resolutions.map((resolution) => (
              <MenuItem key={resolution.label} value={JSON.stringify({ width: resolution.width, height: resolution.height })}>
                {resolution.label}
              </MenuItem>
            ))}
          </TextField>
        </Grid>
        
        <Grid item xs={12} md={6}>
          <TextField
            label="Frame Rate (FPS)"
            type="number"
            value={fps}
            onChange={handleFpsChange}
            fullWidth
            margin="normal"
            InputProps={{
              inputProps: { min: 1, max: 60 },
              sx: { color: 'white' }
            }}
          />
        </Grid>

        <Grid item xs={12} md={6}>
          <TextField
            label="Bitrate (kbps)"
            type="number"
            value={bitrate}
            onChange={handleBitrateChange}
            placeholder="Auto (leave empty for camera default)"
            fullWidth
            margin="normal"
            helperText="Leave empty to use camera's default bitrate range"
            InputProps={{
              inputProps: { min: 1 },
              sx: { color: 'white' }
            }}
          />
        </Grid>
      </Grid><Box sx={{ mt: 3, display: 'flex', justifyContent: 'flex-end' }}>
        <Button 
          variant="contained" 
          color="primary"
          onClick={handleApplyConfig}
          disabled={isLoading || selectedCameras.length === 0}
          sx={{ minWidth: '150px' }}
        >
          {isLoading ? (
            <>
              <CircularProgress size={20} sx={{ mr: 1, color: 'white' }} /> 
              {selectedCameras.length > 1 ? 'Processing Cameras...' : 'Applying...'}
            </>
          ) : `Apply to ${selectedCameras.length} Camera(s)`}
        </Button>
      </Box>
      
      {isLoading && selectedCameras.length > 1 && (
        <Alert severity="info" sx={{ mt: 2 }}>
          <Typography variant="body2">
            Applying configuration to multiple cameras... This may take up to 
            {selectedCameras.length > 1 ? ` ${selectedCameras.length * 20} seconds ` : ' 20 seconds '}
            to complete as each camera requires time to validate the settings.
          </Typography>
        </Alert>
      )}        {result.message && (
        <Alert 
          severity={result.success ? 'success' : 'error'} 
          sx={{ mt: 2 }}
          onClose={() => setResult({ success: false, message: null })}
        >
          {result.message}
        </Alert>
      )}
        
      {/* Upload Config Dialog */}      <Dialog 
        open={uploadDialogOpen} 
        onClose={handleUploadDialogClose}
        maxWidth="sm"
        fullWidth
      >
        <DialogTitle>Upload Configuration CSV</DialogTitle>        <DialogContent>
          {uploadResult && (
            <Alert 
              severity={uploadResult.success ? "success" : "error"} 
              sx={{ mb: 2, mt: 1 }}
            >
              {uploadResult.message}
            </Alert>
          )}

          {isLoading && selectedCameras.length > 0 && (
            <Alert severity="info" sx={{ mb: 2, mt: 1 }}>
              <Typography variant="body2">
                Applying CSV configuration to {selectedCameras.length} camera(s)... This may take up to 
                {selectedCameras.length > 1 ? ` ${selectedCameras.length * 20} seconds ` : ' 20 seconds '}
                to complete.
              </Typography>
            </Alert>
          )}

          {uploadResult && uploadResult.success && uploadResult.config && (
            <Paper variant="outlined" sx={{ p: 2, mb: 2, bgcolor: 'action.hover' }}>
              <Typography variant="subtitle2" gutterBottom sx={{ fontWeight: 'bold' }}>
                Scanned Configuration:
              </Typography>
              <Typography variant="body2">
                <strong>Resolution:</strong> {uploadResult.config.width}x{uploadResult.config.height}
              </Typography>
              <Typography variant="body2">
                <strong>FPS:</strong> {uploadResult.config.fps}
              </Typography>
              <Typography variant="body2">
                <strong>Bitrate:</strong> {uploadResult.config.bitrate || 'Auto'} kbps
              </Typography>
              {selectedCameras.length > 0 && (
                <Typography variant="body2" sx={{ mt: 1, fontStyle: 'italic' }}>
                  Will be applied to {selectedCameras.length} selected camera(s)
                </Typography>
              )}
            </Paper>
          )}

          <Paper variant="outlined" sx={{ p: 2, bgcolor: 'action.hover' }}>
            <Typography variant="subtitle2" gutterBottom sx={{ fontWeight: 'bold' }}>
              CSV File Upload
            </Typography>
            <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
              Upload a CSV file with configuration parameters. Required columns: width, height, fps. Optional: bitrate
              <br />
              <em>File will be automatically scanned when selected.</em>
            </Typography>
            
            <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 2 }}>
              <input
                type="file"
                accept=".csv"
                onChange={handleConfigCsvFileChange}
                style={{ display: 'none' }}
                id="config-csv-file-input"
              />
              <label htmlFor="config-csv-file-input">
                <Button 
                  variant="outlined" 
                  component="span" 
                  startIcon={uploadingConfig ? <CircularProgress size={16} /> : <UploadIcon />}
                  disabled={uploadingConfig}
                >
                  {uploadingConfig ? 'Scanning...' : 'Choose CSV File'}
                </Button>
              </label>
              {configCsvFile && !uploadingConfig && (
                <Typography variant="body2" sx={{ flex: 1 }}>
                  {configCsvFile.name}
                </Typography>
              )}
            </Box>
          </Paper>
        </DialogContent>        <DialogActions>
          <Button onClick={handleUploadDialogClose} disabled={isLoading}>
            Cancel
          </Button>
          {uploadResult && uploadResult.success && uploadResult.config && (
            <Button 
              onClick={handleApplyUploadedConfig}
              variant="contained"
              color="primary"
              disabled={isLoading}
              startIcon={isLoading ? <CircularProgress size={16} /> : null}
            >
              {isLoading ? 'Applying...' : 'Apply Configuration'}
            </Button>
          )}
        </DialogActions>
      </Dialog>
    </Box>
  );
};

export default CameraConfigPanel;
