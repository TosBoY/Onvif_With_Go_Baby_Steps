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
  };    const handleApplyConfig = async () => {
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
      message: `Applying configuration to ${selectedCameras.length} camera(s)... This may take a while as each camera needs time to update settings.`
    });
    
    const results = [];
    const errors = [];
    const validations = [];
    
    // Show progressive updates
    const updateProgress = (processed, total) => {
      setResult({
        success: true,
        message: `Processing cameras: ${processed}/${total} complete... Please wait as validation can take up to 20 seconds per camera.`
      });
    };
    
    // Apply configuration to all selected cameras
    for (let i = 0; i < selectedCameras.length; i++) {
      const cameraId = selectedCameras[i];
      updateProgress(i, selectedCameras.length);
      
      try {
        // Update timeout dynamically based on number of cameras
        const result = await applyConfig(cameraId, width, height, fps);
        console.log(`Camera ${cameraId} configuration result:`, result);
        
        // Store validation results for each camera
        if (result && result.validation) {
          // Add camera ID to the validation result for identification
          validations.push({
            ...result.validation,
            cameraId
          });
        }
        
        results.push(cameraId);
      } catch (error) {
        console.error(`Error applying configuration to camera ${cameraId}:`, error);
        // Provide more context in the error message
        const errorMsg = error.message || "Unknown error";
        errors.push({ 
          cameraId, 
          error: errorMsg.includes("timeout") 
            ? `Timeout - camera may need more time to apply settings (${errorMsg})` 
            : errorMsg 
        });
      }
      
      // Brief pause between camera configurations to avoid overloading
      if (i < selectedCameras.length - 1) {
        await new Promise(resolve => setTimeout(resolve, 500));
      }
    }
      // Pass all validation results back to Dashboard, even if there are errors
    // This ensures we show whatever validation data we were able to collect
    if (onConfigurationApplied && validations.length > 0) {
      // Create composite result with all validations
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
    if (errors.length === 0) {
      setResult({ 
        success: true, 
        message: `Configuration applied successfully to ${results.length} camera(s): ${results.join(', ')}` 
      });
    } else if (results.length > 0) {
      // For partial success, be more specific about what worked
      const successMessage = `Applied to cameras ${results.join(', ')}`;
      const errorDetails = errors.map(e => `${e.cameraId}: ${e.error}`).join('; ');
      
      setResult({ 
        success: false, 
        message: `Partial success: ${successMessage}. Failed cameras: ${errorDetails}` 
      });
    } else {
      // For complete failure, show detailed error for each camera
      const errorDetails = errors.map(e => `${e.cameraId}: ${e.error}`).join('; ');
      
      setResult({ 
        success: false, 
        message: `Failed to apply configuration to any cameras. Errors: ${errorDetails}` 
      });
    }
    
    setIsLoading(false);
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
        <Box sx={{ mt: 3 }}>
        <Typography variant="h6">Configuration Summary</Typography>
        <Box sx={{ mt: 1, p: 2, bgcolor: '#f5f5f5', borderRadius: 1 }}>
          <Typography variant="body1">
            <strong>Selected Cameras:</strong> {selectedCameras.length > 0 ? selectedCameras.join(', ') : 'None selected'}
          </Typography>
          <Typography variant="body1">
            <strong>Resolution:</strong> {width}×{height}
          </Typography>
          <Typography variant="body1">
            <strong>Frame Rate:</strong> {fps} FPS
          </Typography>
          {selectedCameras.length === 0 && (
            <Typography variant="body2" color="text.secondary" sx={{ mt: 1 }}>
              Select cameras from the list to apply configuration
            </Typography>
          )}
        </Box>
      </Box>
    </Box>
  );
};

export default CameraConfigPanel;
