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

// Default resolutions to use as fallback when API returns empty array
const DEFAULT_RESOLUTIONS = [
  { Width: 1920, Height: 1080 },
  { Width: 1280, Height: 720 },
  { Width: 640, Height: 480 },
  { Width: 320, Height: 240 }
];

const ResolutionManager = ({ configToken, profileToken, refreshCameraInfo }) => {
  const [options, setOptions] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const [autoLaunchVLC, setAutoLaunchVLC] = useState(false);
  const [usingFallbackResolutions, setUsingFallbackResolutions] = useState(false);
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
        console.log("Resolution data received:", data);
        
        // Check if we need to use fallback resolutions
        if (!data.ResolutionsAvailable || data.ResolutionsAvailable.length === 0) {
          console.log("No resolutions available from API, using fallbacks");
          data.ResolutionsAvailable = DEFAULT_RESOLUTIONS;
          setUsingFallbackResolutions(true);
        } else {
          setUsingFallbackResolutions(false);
        }
        
        setOptions(data);
        
        // Set default values when options are loaded with proper null checks
        if (data && data.ResolutionsAvailable && data.ResolutionsAvailable.length > 0) {
          const defaultIndex = 0;
          const defaultRes = data.ResolutionsAvailable[defaultIndex];
          
          // Set a default frame rate based on available options
          let defaultFrameRate = 1;
          if (data.frameRates && data.frameRates.length > 0) {
            // Try to select a middle frame rate value from the available options
            // or the first one if there's only one option
            const midIndex = Math.floor(data.frameRates.length / 2);
            defaultFrameRate = data.frameRates[midIndex];
          } else if (data.FrameRateRange && data.FrameRateRange.Min) {
            defaultFrameRate = data.FrameRateRange.Min;
          }
          
          setFormData({
            resolution: defaultIndex,
            frameRate: defaultFrameRate,
            bitRate: 4096, // Default bitrate
            govLength: data.GovLengthRange && data.GovLengthRange.Min ? data.GovLengthRange.Min : 1,
            h264Profile: data.H264ProfilesSupported && data.H264ProfilesSupported.length > 0 ? data.H264ProfilesSupported[0] : 'Main',
            width: defaultRes.Width,
            height: defaultRes.Height
          });
        } else {
          // No resolution data available, set default values
          setError('No resolution options available for this profile/config combination');
          setFormData({
            resolution: '',
            frameRate: 1,
            bitRate: 4096,
            govLength: 1,
            h264Profile: 'Main'
          });
        }
        
        setLoading(false);
      } catch (err) {
        console.error('Failed to load resolution options:', err);
        setError('Failed to load resolution options');
        setLoading(false);
      }
    };

    fetchOptions();
  }, [configToken, profileToken]);

  const handleResolutionChange = (event) => {
    const selectedResIndex = event.target.value;
    if (options && options.ResolutionsAvailable && selectedResIndex !== undefined && 
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

    if (!options || !options.ResolutionsAvailable || !options.ResolutionsAvailable[formData.resolution]) {
      setError('Invalid resolution selection');
      return;
    }

    const selectedRes = options.ResolutionsAvailable[formData.resolution];

    const configData = {
      configToken,
      profileToken,
      width: selectedRes.Width,
      height: selectedRes.Height,
      frameRate: parseInt(formData.frameRate, 10) || (options.FrameRateRange && options.FrameRateRange.Min ? options.FrameRateRange.Min : 1),
      bitRate: parseInt(formData.bitRate, 10) || 4096,
      gopLength: parseInt(formData.govLength, 10) || (options.GovLengthRange && options.GovLengthRange.Min ? options.GovLengthRange.Min : 1),
      h264Profile: formData.h264Profile || (options.H264ProfilesSupported && options.H264ProfilesSupported.length > 0 ? options.H264ProfilesSupported[0] : 'Main')
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

  // Safe accessors for frame rate and range values with fallbacks
  const getFrameRateMin = () => {
    // First try to get from frameRates array
    if (options.frameRates && options.frameRates.length > 0) {
      return Math.min(...options.frameRates);
    }
    // Then try from FrameRateRange (original format)
    if (options.FrameRateRange && options.FrameRateRange.Min !== undefined) {
      return options.FrameRateRange.Min;
    }
    // Default fallback
    return 1;
  };

  const getFrameRateMax = () => {
    // First try to get from frameRates array
    if (options.frameRates && options.frameRates.length > 0) {
      return Math.max(...options.frameRates);
    }
    // Then try from FrameRateRange (original format)
    if (options.FrameRateRange && options.FrameRateRange.Max !== undefined) {
      return options.FrameRateRange.Max;
    }
    // Default fallback
    return 30;
  };
  
  // Helper function to get available H264 profiles
  const getH264Profiles = () => {
    // First try the 'h264Profiles' array (from Pi component)
    if (options.h264Profiles && options.h264Profiles.length > 0) {
      return options.h264Profiles;
    }
    // Then try the 'H264ProfilesSupported' array (original format)
    if (options.H264ProfilesSupported && options.H264ProfilesSupported.length > 0) {
      return options.H264ProfilesSupported;
    }
    // Default fallback
    return ["Baseline", "Main", "High"];
  };

  const hasSpecificFrameRateValues = () => {
    return options.frameRates && options.frameRates.length > 0;
  };

  // Safely get range values with fallbacks using our helper functions
  const frameRateMin = getFrameRateMin();
  const frameRateMax = getFrameRateMax();
  const h264Profiles = getH264Profiles();
  const govLengthMin = options.GovLengthRange && options.GovLengthRange.Min ? options.GovLengthRange.Min : 1;
  const govLengthMax = options.GovLengthRange && options.GovLengthRange.Max ? options.GovLengthRange.Max : 60;

  return (
    <Card variant="outlined" sx={{ mb: 3 }}>
      <CardContent>
        <Typography variant="h6" gutterBottom>
          Change Resolution Settings
        </Typography>
        
        {usingFallbackResolutions && (
          <Alert severity="info" sx={{ mb: 2 }}>
            No resolution options were returned by the camera. Using common resolutions instead.
            These may not all be supported by your device.
          </Alert>
        )}

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
          label={`Frame Rate (${frameRateMin}-${frameRateMax} fps)`}
          type="number"
          name="frameRate"
          value={formData.frameRate}
          onChange={handleInputChange}
          InputProps={{ inputProps: { min: frameRateMin, max: frameRateMax } }}
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
          label={`GOP Length (${govLengthMin}-${govLengthMax})`}
          type="number"
          name="govLength"
          value={formData.govLength}
          onChange={handleInputChange}
          InputProps={{ inputProps: { min: govLengthMin, max: govLengthMax } }}
          fullWidth
          sx={{ mb: 2 }}
        />

        {/* H264 profile selector */}
        {h264Profiles && h264Profiles.length > 0 && (
          <FormControl fullWidth sx={{ mb: 2 }}>
            <InputLabel id="h264-profile-label">H264 Profile</InputLabel>
            <Select
              labelId="h264-profile-label"
              id="h264-profile"
              name="h264Profile"
              value={formData.h264Profile || h264Profiles[0]}
              label="H264 Profile"
              onChange={handleInputChange}
            >
              {h264Profiles.map((profile) => (
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