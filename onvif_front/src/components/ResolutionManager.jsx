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

// Standard resolution definitions
const STANDARD_RESOLUTIONS = [
  { label: '480p', width: 640, height: 480 },    // SD
  { label: '720p', width: 1280, height: 720 },   // HD
  { label: '1080p', width: 1920, height: 1080 }, // Full HD
  { label: '2K', width: 2048, height: 1080 },    // 2K DCI
  { label: '1440p', width: 2560, height: 1440 }, // QHD
  { label: '4K', width: 3840, height: 2160 },    // 4K UHD
  { label: '5K', width: 5120, height: 2880 },    // 5K
  { label: '8K', width: 7680, height: 4320 },    // 8K UHD
];

// Function to find closest standard resolution
const findClosestResolution = (availableResolutions) => {
  const result = [];
  const seen = new Set(); // To prevent duplicate resolutions

  for (const stdRes of STANDARD_RESOLUTIONS) {
    let closestRes = null;
    let minDiff = Number.MAX_SAFE_INTEGER;
    let bestAspectRatioDiff = Number.MAX_SAFE_INTEGER;

    for (const res of availableResolutions) {
      const targetArea = stdRes.width * stdRes.height;
      const resArea = res.Width * res.Height;
      const areaDiff = Math.abs(targetArea - resArea);
      
      // Calculate aspect ratio difference
      const targetRatio = stdRes.width / stdRes.height;
      const resRatio = res.Width / res.Height;
      const aspectRatioDiff = Math.abs(targetRatio - resRatio);

      // Only consider if aspect ratio is similar enough (within 10%)
      if (aspectRatioDiff < 0.1 && (areaDiff < minDiff || 
         (areaDiff === minDiff && aspectRatioDiff < bestAspectRatioDiff))) {
        minDiff = areaDiff;
        bestAspectRatioDiff = aspectRatioDiff;
        closestRes = res;
      }
    }

    if (closestRes) {
      const key = `${closestRes.Width}x${closestRes.Height}`;
      if (!seen.has(key)) {
        seen.add(key);
        result.push({
          ...closestRes,
          label: stdRes.label,
          originalWidth: stdRes.width,
          originalHeight: stdRes.height
        });
      }
    }
  }

  return result;
};

const ResolutionManager = ({ configToken, profileToken, refreshCameraInfo, selectedCameras }) => {
  const [options, setOptions] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const [autoLaunchVLC, setAutoLaunchVLC] = useState(false);
  const [applyingConfig, setApplyingConfig] = useState(false);
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
        
        if (!data.ResolutionsAvailable || data.ResolutionsAvailable.length === 0) {
          console.log("No resolutions available from API, using standard resolutions");
          // Convert standard resolutions to match API format
          data.ResolutionsAvailable = STANDARD_RESOLUTIONS.map(res => ({
            Width: res.width,
            Height: res.height,
            label: res.label
          }));
        } else {
          // Find closest matches to standard resolutions
          data.ResolutionsAvailable = findClosestResolution(data.ResolutionsAvailable);
        }
        
        setOptions(data);
        
        // Set default values
        if (data.ResolutionsAvailable && data.ResolutionsAvailable.length > 0) {
          // Try to find 1080p as default, otherwise use first available
          const defaultIndex = data.ResolutionsAvailable.findIndex(res => 
            res.label === '1080p' || (res.Width === 1920 && res.Height === 1080)
          ) || 0;
          
          const defaultRes = data.ResolutionsAvailable[defaultIndex];
          
          // Get default frame rate
          let defaultFrameRate = 30; // Common default
          if (data.frameRates && data.frameRates.length > 0) {
            const midIndex = Math.floor(data.frameRates.length / 2);
            defaultFrameRate = data.frameRates[midIndex];
          } else if (data.FrameRateRange && data.FrameRateRange.Min) {
            defaultFrameRate = data.FrameRateRange.Min;
          }

          setFormData({
            resolution: defaultIndex,
            frameRate: defaultFrameRate,
            bitRate: 4096,
            govLength: data.GovLengthRange?.Min || 1,
            h264Profile: data.h264Profiles?.[0] || data.H264ProfilesSupported?.[0] || 'Main',
            width: defaultRes.Width,
            height: defaultRes.Height
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
    if (options?.ResolutionsAvailable?.[selectedResIndex]) {
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
    // Validate camera selection first
    if (selectedCameras.length === 0) {
      setError('Please select at least one camera to apply the configuration');
      return;
    }

    // Validate resolution selection
    if (formData.resolution === undefined || formData.resolution === null || formData.resolution === '') {
      setError('Please select a resolution');
      return;
    }

    if (!options || !options.ResolutionsAvailable || !options.ResolutionsAvailable[formData.resolution]) {
      setError('Invalid resolution selection');
      return;
    }

    setApplyingConfig(true);
    setError('');
    setSuccess('');

    const selectedRes = options.ResolutionsAvailable[formData.resolution];
    const baseConfigData = {
      configToken,
      profileToken,
      width: selectedRes.Width,
      height: selectedRes.Height,
      frameRate: parseInt(formData.frameRate, 10) || frameRateMin,
      bitRate: parseInt(formData.bitRate, 10) || 4096,
      gopLength: parseInt(formData.govLength, 10) || govLengthMin,
      h264Profile: formData.h264Profile || (h264Profiles.length > 0 ? h264Profiles[0] : 'Main')
    };

    const successfulCameras = [];
    const failedCameras = [];

    try {
      // Update each selected camera in sequence
      for (const cameraId of selectedCameras) {
        try {
          const cameraConfigData = {
            ...baseConfigData,
            cameraId
          };
          await api.changeResolution(cameraConfigData);
          successfulCameras.push(cameraId);
        } catch (err) {
          console.error(`Failed to update camera ${cameraId}:`, err);
          failedCameras.push(cameraId);
        }
      }

      // Set appropriate success/error message based on results
      if (successfulCameras.length > 0) {
        const successMessage = successfulCameras.length === selectedCameras.length
          ? `Configuration successfully applied to all ${successfulCameras.length} cameras`
          : `Configuration applied to ${successfulCameras.length} out of ${selectedCameras.length} cameras`;
        
        setSuccess(successMessage);
        
        // Refresh displayed information
        console.log("Configuration updated, refreshing display...");
        refreshCameraInfo();
        
        // Additional refresh after delay to ensure update is captured
        setTimeout(() => {
          console.log("Calling refresh again after timeout");
          refreshCameraInfo();
        }, 500);
        
        // Handle VLC auto-launch if enabled and at least one camera was updated successfully
        if (autoLaunchVLC) {
          try {
            const vlcResponse = await api.launchVLC(profileToken);
            setSuccess(`${successMessage} and ${vlcResponse.message.toLowerCase()}`);
          } catch (vlcError) {
            console.error('Error launching VLC:', vlcError);
            setSuccess(`${successMessage} but failed to launch VLC`);
          }
        }
      }

      // If any cameras failed, show error
      if (failedCameras.length > 0) {
        const errorMessage = successfulCameras.length === 0
          ? 'Failed to apply configuration to all cameras'
          : `Failed to apply configuration to ${failedCameras.length} camera(s)`;
        setError(errorMessage);
      }

    } catch (err) {
      console.error('Error updating camera configurations:', err);
      setError('Failed to update camera configurations: ' + err.message);
    } finally {
      setApplyingConfig(false);
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

  // Helper function for getting GOP length min value
  const getGovLengthMin = () => {
    // First check if we have a GovLengthRange with Min (our new backend format)
    if (options.GovLengthRange && options.GovLengthRange.Min !== undefined) {
      return options.GovLengthRange.Min;
    }
    // Then try 'encodingIntervals' array for min value (from Pi component)
    if (options.encodingIntervals && options.encodingIntervals.length > 0) {
      return Math.min(...options.encodingIntervals);
    }
    // Default fallback
    return 1;
  };

  // Helper function for getting GOP length max value
  const getGovLengthMax = () => {
    // First check if we have a GovLengthRange with Max (our new backend format)
    if (options.GovLengthRange && options.GovLengthRange.Max !== undefined) {
      return options.GovLengthRange.Max;
    }
    // Then try 'encodingIntervals' array for max value (from Pi component)
    if (options.encodingIntervals && options.encodingIntervals.length > 0) {
      return Math.max(...options.encodingIntervals);
    }
    // Default fallback
    return 60;
  };

  const hasSpecificFrameRateValues = () => {
    return options.frameRates && options.frameRates.length > 0;
  };

  // Safely get range values with fallbacks using our helper functions
  const frameRateMin = getFrameRateMin();
  const frameRateMax = getFrameRateMax();
  const h264Profiles = getH264Profiles();
  const govLengthMin = getGovLengthMin();
  const govLengthMax = getGovLengthMax();

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
                {res.label || `${res.Width}x${res.Height}`} 
                {res.label && (res.Width !== res.originalWidth || res.Height !== res.originalHeight) && 
                  ` (closest match: ${res.Width}x${res.Height})`}
              </MenuItem>
            ))}
          </Select>
        </FormControl>

        {/* Selected resolution preview with standard label */}
        {formData.width && formData.height && (
          <Typography variant="subtitle2" sx={{ mb: 2, color: 'text.secondary' }}>
            Selected: {
              options.ResolutionsAvailable[formData.resolution]?.label || 
              `${formData.width}x${formData.height}`
            }
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

        {/* Apply Configuration button with camera count */}
        <Button 
          variant="contained" 
          color="primary" 
          onClick={handleSubmit}
          disabled={selectedCameras.length === 0 || applyingConfig}
          fullWidth
          sx={{
            backgroundColor: selectedCameras.length === 0 ? 'action.disabledBackground' : 'primary.main',
            '&:disabled': {
              backgroundColor: 'action.disabledBackground',
              color: 'text.disabled'
            }
          }}
        >
          {applyingConfig ? (
            'Applying Configuration...'
          ) : selectedCameras.length === 0 ? (
            'Select Cameras to Apply Configuration'
          ) : (
            `Apply Configuration to ${selectedCameras.length} Camera${selectedCameras.length === 1 ? '' : 's'}`
          )}
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