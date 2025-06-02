import { useState } from 'react';
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
  { width: 640, height: 480, label: '640x480 (SD)' },
  { width: 1280, height: 720, label: '1280x720 (HD)' },
  { width: 1920, height: 1080, label: '1920x1080 (Full HD)' },
  { width: 2560, height: 1440, label: '2560x1440 (2K)' },
  { width: 3840, height: 2160, label: '3840x2160 (4K)' },
];

// Common FPS options
const fpsOptions = [15, 25, 30, 60];

const CameraConfigPanel = ({ selectedCamera, selectedCameras = [], cameras = [] }) => {
  const [width, setWidth] = useState(1280);
  const [height, setHeight] = useState(720);
  const [fps, setFps] = useState(30);
  const [isLoading, setIsLoading] = useState(false);
  const [result, setResult] = useState({ success: false, message: null });

  const handleResolutionChange = (event) => {
    const selectedResolution = JSON.parse(event.target.value);
    setWidth(selectedResolution.width);
    setHeight(selectedResolution.height);
  };

  const handleFpsChange = (event) => {
    setFps(event.target.value);
  };  const handleApplyConfig = async () => {
    if (selectedCameras.length === 0) {
      setResult({ 
        success: false, 
        message: 'Please select at least one camera from the list to apply configuration' 
      });
      return;
    }
    
    setIsLoading(true);
    setResult({ success: false, message: null });
    
    const results = [];
    const errors = [];
    
    // Apply configuration to all selected cameras
    for (const cameraId of selectedCameras) {
      try {
        await applyConfig(cameraId, width, height, fps);
        results.push(cameraId);
      } catch (error) {
        errors.push({ cameraId, error: error.message });
        console.error(`Error applying configuration to camera ${cameraId}:`, error);
      }
    }
    
    // Show results
    if (errors.length === 0) {
      setResult({ 
        success: true, 
        message: `Configuration applied successfully to ${results.length} camera(s): ${results.join(', ')}` 
      });
    } else if (results.length > 0) {
      setResult({ 
        success: false, 
        message: `Partial success: Applied to cameras ${results.join(', ')}. Failed: ${errors.map(e => e.cameraId).join(', ')}` 
      });
    } else {
      setResult({ 
        success: false, 
        message: `Failed to apply configuration to all cameras. Errors: ${errors.map(e => `${e.cameraId}: ${e.error}`).join('; ')}` 
      });
    }
    
    setIsLoading(false);
  };  return (
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
            select
            label="Frame Rate (FPS)"
            value={fps}
            onChange={handleFpsChange}
            fullWidth
            margin="normal"
          >
            {fpsOptions.map((option) => (
              <MenuItem key={option} value={option}>
                {option} FPS
              </MenuItem>
            ))}
          </TextField>
        </Grid>      </Grid>

      <Box sx={{ mt: 3, display: 'flex', justifyContent: 'flex-end' }}>
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
              Applying...
            </>
          ) : `Apply to ${selectedCameras.length} Camera(s)`}
        </Button>
      </Box>
      
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
