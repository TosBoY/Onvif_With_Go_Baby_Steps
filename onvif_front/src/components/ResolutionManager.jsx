import { useState } from 'react';
import { 
  Box, 
  FormControl, 
  InputLabel, 
  Select, 
  MenuItem, 
  TextField, 
  Button, 
  Alert,
  Typography,
} from '@mui/material';
import api from '../services/api';
import ValidationResultDisplay from './ValidationResultDisplay';

// Standard resolution definitions
const STANDARD_RESOLUTIONS = [
  { label: '480p', width: 640, height: 480 },    // SD
  { label: '720p', width: 1280, height: 720 },   // HD
  { label: '1080p', width: 1920, height: 1080 }, // Full HD
  { label: '2K', width: 2048, height: 1080 },    // 2K DCI
  { label: '1440p', width: 2560, height: 1440 }, // QHD
  { label: '4K', width: 3840, height: 2160 },    // 4K UHD
];

const ResolutionManager = ({ selectedCameras }) => {
  const [selectedResolution, setSelectedResolution] = useState('');
  const [frameRate, setFrameRate] = useState('25');
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const [validationResults, setValidationResults] = useState([]);
  const [applyingConfig, setApplyingConfig] = useState(false);

  const handleResolutionChange = (event) => {
    setSelectedResolution(event.target.value);
  };

  const handleFrameRateChange = (event) => {
    setFrameRate(event.target.value);
  };

  const handleApplyConfig = async () => {
    if (selectedCameras.length === 0) {
      setError('Please select at least one camera');
      return;
    }

    if (selectedResolution === '') {
      setError('Please select a resolution');
      return;
    }

    const resolution = STANDARD_RESOLUTIONS[selectedResolution];
    if (!resolution) {
      setError('Invalid resolution selected');
      return;
    }

    setApplyingConfig(true);
    setError('');
    setSuccess('');
    setValidationResults([]);

    try {
      const result = await api.changeResolutionSimple({
        cameraIds: selectedCameras,
        width: resolution.width,
        height: resolution.height,
        frameRate: parseInt(frameRate) || 25,
        bitRate: 4096, // Default bitrate
        gopLength: 30,  // Default GOP length
        h264Profile: 'Main' // Default profile
      });

      setValidationResults(result.results || []);

      const successCount = result.results?.filter(r => r.success).length || 0;
      if (successCount > 0) {
        setSuccess(`Successfully updated ${successCount} camera(s)`);
      } else {
        setError('Failed to update any cameras');
      }
    } catch (err) {
      console.error('Error updating camera settings:', err);
      setError('Error updating camera settings: ' + (err.message || 'Unknown error'));
    } finally {
      setApplyingConfig(false);
    }
  };

  return (
    <Box sx={{ 
      mt: 2,
      width: '100%',
      maxWidth: '100%',
      display: 'flex',
      flexDirection: 'column',
      '& > *': { width: '100%', maxWidth: '100%' }
    }}>
      <Box sx={{ mb: 2, width: '100%', maxWidth: '100%' }}>
        <FormControl fullWidth sx={{ mb: 2 }}>
          <InputLabel id="resolution-select-label">Resolution</InputLabel>
          <Select
            labelId="resolution-select-label"
            id="resolution-select"
            value={selectedResolution}
            label="Resolution"
            onChange={handleResolutionChange}
          >
            {STANDARD_RESOLUTIONS.map((res, index) => (
              <MenuItem key={index} value={index}>
                {res.label} ({res.width}x{res.height})
              </MenuItem>
            ))}
          </Select>
        </FormControl>

        <FormControl fullWidth sx={{ mb: 2 }}>
          <TextField
            label="Frame Rate (fps)"
            type="number"
            value={frameRate}
            onChange={handleFrameRateChange}
            InputProps={{ 
              inputProps: { min: 1, max: 60 }
            }}
          />
        </FormControl>

        <Button
          variant="contained"
          color="primary"
          onClick={handleApplyConfig}
          disabled={selectedCameras.length === 0 || selectedResolution === '' || applyingConfig}
          fullWidth
        >
          {applyingConfig ? 'Applying...' : 'Apply Settings'}
        </Button>
      </Box>

      <Box sx={{ width: '100%', maxWidth: '100%' }}>
        {error && (
          <Alert severity="error" sx={{ mb: 2 }}>
            {error}
          </Alert>
        )}
        
        {success && (
          <Alert severity="success" sx={{ mb: 2 }}>
            {success}
          </Alert>
        )}

        {validationResults.length > 0 && (
          <Box sx={{ width: '100%', maxWidth: '100%', overflowX: 'hidden' }}>
            <ValidationResultDisplay results={validationResults} />
          </Box>
        )}
      </Box>
    </Box>
  );
};

export default ResolutionManager;