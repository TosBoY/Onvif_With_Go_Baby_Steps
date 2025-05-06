import { useState, useEffect } from 'react';
import { 
  Box, 
  Card, 
  CardContent, 
  Typography, 
  FormControl, 
  InputLabel, 
  Select, 
  MenuItem, 
  TextField, 
  Button, 
  Alert,
  Snackbar,
  FormControlLabel,
  Checkbox
} from '@mui/material';
import api from '../services/api';

const ResolutionManager = ({ configToken, profileToken, refreshCameraInfo }) => {
  const [options, setOptions] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const [autoLaunchVLC, setAutoLaunchVLC] = useState(false);
  const [formData, setFormData] = useState({
    resolution: '',
    frameRate: '',
    bitRate: '',
    govLength: '',
    h264Profile: ''
  });

  useEffect(() => {
    const fetchOptions = async () => {
      if (!configToken || !profileToken) return;
      
      setLoading(true);
      try {
        const data = await api.getResolutions(configToken, profileToken);
        setOptions(data);
        
        // Set default values when options are loaded
        if (data.ResolutionsAvailable && data.ResolutionsAvailable.length > 0) {
          const defaultIndex = 0;
          const defaultRes = data.ResolutionsAvailable[defaultIndex];
          
          setFormData({
            resolution: defaultIndex,
            frameRate: data.FrameRateRange.Min,
            bitRate: 4096, // Default bitrate
            govLength: data.GovLengthRange.Min,
            h264Profile: data.H264ProfilesSupported ? data.H264ProfilesSupported[0] : 'Main',
            width: defaultRes.Width,
            height: defaultRes.Height
          });
        }
        
        setLoading(false);
      } catch (err) {
        setError('Failed to load resolution options');
        setLoading(false);
      }
    };

    fetchOptions();
  }, [configToken, profileToken]);

  const handleResolutionChange = (event) => {
    const selectedResIndex = event.target.value;
    if (options.ResolutionsAvailable && selectedResIndex !== undefined && 
        selectedResIndex !== null && options.ResolutionsAvailable[selectedResIndex]) {
      const selectedRes = options.ResolutionsAvailable[selectedResIndex];
      
      setFormData(prevData => ({
        ...prevData,
        resolution: selectedResIndex,
        width: selectedRes.Width,
        height: selectedRes.Height
      }));
    }
  };

  const handleInputChange = (event) => {
    const { name, value } = event.target;
    setFormData(prevData => ({
      ...prevData,
      [name]: value
    }));
  };

  const handleSubmit = async () => {
    if (formData.resolution === undefined || formData.resolution === null || formData.resolution === '') {
      setError('Please select a resolution');
      return;
    }

    if (!options.ResolutionsAvailable || !options.ResolutionsAvailable[formData.resolution]) {
      setError('Invalid resolution selection');
      return;
    }

    const selectedRes = options.ResolutionsAvailable[formData.resolution];

    const configData = {
      configToken,
      profileToken,
      width: selectedRes.Width,
      height: selectedRes.Height,
      frameRate: parseInt(formData.frameRate, 10) || options.FrameRateRange.Min,
      bitRate: parseInt(formData.bitRate, 10) || 4096,
      gopLength: parseInt(formData.govLength, 10) || options.GovLengthRange.Min,
      h264Profile: formData.h264Profile || (options.H264ProfilesSupported ? options.H264ProfilesSupported[0] : 'Main')
    };

    try {
      // First, update the camera configuration
      await api.changeResolution(configData);
      setSuccess('Camera configuration updated successfully');
      
      // Refresh camera info to show updated configuration
      refreshCameraInfo();
      
      // If auto-launch VLC is checked, launch VLC
      if (autoLaunchVLC) {
        try {
          setSuccess('Applying configuration and launching VLC...');
          const vlcResponse = await api.launchVLC(profileToken);
          console.log('VLC launch response:', vlcResponse);
          setSuccess(`Configuration updated and ${vlcResponse.message.toLowerCase()}`);
        } catch (vlcError) {
          console.error('Error launching VLC:', vlcError);
          setSuccess('Configuration updated but failed to launch VLC');
        }
      }
    } catch (err) {
      console.error('Error updating camera configuration:', err);
      setError('Failed to update camera configuration');
    }
  };

  const handleCloseAlert = () => {
    setError('');
    setSuccess('');
  };

  if (loading) {
    return <Typography>Loading options...</Typography>;
  }

  if (!options) {
    return <Typography>Select a valid profile and configuration token first</Typography>;
  }

  return (
    <Card variant="outlined" sx={{ mb: 3 }}>
      <CardContent>
        <Typography variant="h6" gutterBottom>
          Change Resolution Settings
        </Typography>

        {/* Resolution selector */}
        <FormControl fullWidth sx={{ mb: 2 }}>
          <InputLabel id="resolution-select-label">Resolution</InputLabel>
          <Select
            labelId="resolution-select-label"
            id="resolution-select"
            value={formData.resolution !== undefined && formData.resolution !== null ? formData.resolution : ''}
            label="Resolution"
            onChange={handleResolutionChange}
          >
            {options.ResolutionsAvailable && options.ResolutionsAvailable.map((res, index) => (
              <MenuItem key={index} value={index}>
                {res.Width} x {res.Height}
              </MenuItem>
            ))}
          </Select>
        </FormControl>

        {/* Selected resolution preview */}
        {formData.width && formData.height && (
          <Typography variant="subtitle2" sx={{ mb: 2, color: 'text.secondary' }}>
            Selected: {formData.width} x {formData.height}
          </Typography>
        )}

        {/* Frame rate input */}
        <TextField
          label={`Frame Rate (${options.FrameRateRange.Min}-${options.FrameRateRange.Max} fps)`}
          type="number"
          name="frameRate"
          value={formData.frameRate}
          onChange={handleInputChange}
          InputProps={{ inputProps: { min: options.FrameRateRange.Min, max: options.FrameRateRange.Max } }}
          fullWidth
          sx={{ mb: 2 }}
        />

        {/* Bit rate input */}
        <TextField
          label="Bit Rate (kbps)"
          type="number"
          name="bitRate"
          value={formData.bitRate}
          onChange={handleInputChange}
          InputProps={{ inputProps: { min: 256, max: 20000 } }}
          fullWidth
          sx={{ mb: 2 }}
        />

        {/* GOP length input */}
        <TextField
          label={`GOP Length (${options.GovLengthRange.Min}-${options.GovLengthRange.Max})`}
          type="number"
          name="govLength"
          value={formData.govLength}
          onChange={handleInputChange}
          InputProps={{ inputProps: { min: options.GovLengthRange.Min, max: options.GovLengthRange.Max } }}
          fullWidth
          sx={{ mb: 2 }}
        />

        {/* H264 profile selector */}
        {options.H264ProfilesSupported && options.H264ProfilesSupported.length > 0 && (
          <FormControl fullWidth sx={{ mb: 2 }}>
            <InputLabel id="h264-profile-label">H264 Profile</InputLabel>
            <Select
              labelId="h264-profile-label"
              id="h264-profile"
              name="h264Profile"
              value={formData.h264Profile || options.H264ProfilesSupported[0]}
              label="H264 Profile"
              onChange={handleInputChange}
            >
              {options.H264ProfilesSupported.map((profile) => (
                <MenuItem key={profile} value={profile}>
                  {profile}
                </MenuItem>
              ))}
            </Select>
          </FormControl>
        )}

        {/* Auto-launch VLC checkbox */}
        <FormControlLabel
          control={
            <Checkbox
              checked={autoLaunchVLC}
              onChange={(e) => setAutoLaunchVLC(e.target.checked)}
              name="autoLaunchVLC"
            />
          }
          label="Automatically launch VLC after applying configuration"
          sx={{ mb: 2 }}
        />

        <Button 
          variant="contained" 
          color="primary" 
          onClick={handleSubmit}
          fullWidth
        >
          Apply Configuration
        </Button>
      </CardContent>

      <Snackbar open={!!error} autoHideDuration={6000} onClose={handleCloseAlert}>
        <Alert onClose={handleCloseAlert} severity="error" sx={{ width: '100%' }}>
          {error}
        </Alert>
      </Snackbar>

      <Snackbar open={!!success} autoHideDuration={6000} onClose={handleCloseAlert}>
        <Alert onClose={handleCloseAlert} severity="success" sx={{ width: '100%' }}>
          {success}
        </Alert>
      </Snackbar>
    </Card>
  );
};

export default ResolutionManager;