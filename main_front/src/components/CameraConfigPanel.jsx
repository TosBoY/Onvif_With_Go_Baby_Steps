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
  Grid
} from '@mui/material';
import { applyConfig } from '../services/api';

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
}) => {
  const [width, setWidth] = useState(1280);
  const [height, setHeight] = useState(720);
  const [fps, setFps] = useState(25); // Set default FPS to 25
  const [isLoading, setIsLoading] = useState(false);
  const [result, setResult] = useState({ success: false, message: null });

  const handleResolutionChange = (event) => {
    const selectedResolution = JSON.parse(event.target.value);
    setWidth(selectedResolution.width);
    setHeight(selectedResolution.height);
  };

  const handleFpsChange = (event) => {
    const value = parseInt(event.target.value, 10);
    if (!isNaN(value) && value > 0) {
      setFps(value);
    }
  };  const handleApplyConfig = async () => {
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
    }
    
    setIsLoading(true);
    setResult({ 
      success: true, 
      message: `Applying configuration to ${selectedCameras.length} camera(s)... This may take a while as cameras need time to update and validate settings.`
    });
    
    // Use batch mode to send all camera IDs at once
    try {
      // Call API with array of camera IDs (batch mode)
      const batchResult = await applyConfig(selectedCameras, width, height, fps);
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
            if (result.validation) {
              const validationData = {
                ...result.validation,
                cameraId,
                // Add expected values from the original request
                expectedWidth: width,
                expectedHeight: height,
                expectedFPS: fps
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
      if (onConfigurationApplied && validations.length > 0) {
        const compositeResult = {
          validation: validations,
          appliedConfig: {
            resolution: { width, height },
            fps
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
  };
  return (
    <Box sx={{ p: 2 }}>
      <Typography variant="h6" gutterBottom>
        Camera Configuration
      </Typography>
      
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
        </Grid>      </Grid>      <Box sx={{ mt: 3, display: 'flex', justifyContent: 'flex-end' }}>
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
      )}
        {result.message && (
        <Alert 
          severity={result.success ? 'success' : 'error'} 
          sx={{ mt: 2 }}
          onClose={() => setResult({ success: false, message: null })}
        >
          {result.message}
        </Alert>
      )}
        {/* Configuration controls remain, but the summary section is removed */}
      </Box>
  );
};

export default CameraConfigPanel;
