import React, { useEffect, useState } from 'react';
import { Container, Typography, Box, CircularProgress, Alert, List, ListItem, ListItemText, Paper, Checkbox, FormControlLabel, FormGroup, Button, TextField, FormControl, InputLabel, Select, MenuItem, Grid } from '@mui/material';
import api from './services/api';
import type { Camera, ApplyConfigPayload } from './types';

function App() {
  const [cameras, setCameras] = useState<Camera[]>([]);
  const [selectedCameraIds, setSelectedCameraIds] = useState<string[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [configLoading, setConfigLoading] = useState<boolean>(false);
  const [configError, setConfigError] = useState<string | null>(null);
  const [configSuccess, setConfigSuccess] = useState<string | null>(null);

  // State for configuration inputs
  const [configWidth, setConfigWidth] = useState<string>('');
  const [configHeight, setConfigHeight] = useState<string>('');
  const [configFps, setConfigFps] = useState<string>('25');

  // State for selected resolution preset
  const [selectedResolutionPreset, setSelectedResolutionPreset] = useState<string>('');

  // Define common resolution presets
  const resolutionPresets = [
    { label: 'Select Resolution', width: '', height: '' }, // Placeholder
    { label: '640x480 (VGA)', width: '640', height: '480' },
    { label: '800x600 (SVGA)', width: '800', height: '600' },
    { label: '1024x768 (XGA)', width: '1024', height: '768' },
    { label: '1280x720 (HD)', width: '1280', height: '720' },
    { label: '1280x1024 (SXGA)', width: '1280', height: '1024' },
    { label: '1366x768 (WXGA)', width: '1366', height: '768' },
    { label: '1440x900 (WXGA+)', width: '1440', height: '900' },
    { label: '1600x900 (HD+)', width: '1600', height: '900' },
    { label: '1600x1200 (UXGA)', width: '1600', height: '1200' },
    { label: '1680x1050 (WSXGA+)', width: '1680', height: '1050' },
    { label: '1920x1080 (FHD)', width: '1920', height: '1080' },
    { label: '1920x1200 (WUXGA)', width: '1920', height: '1200' },
    { label: '2560x1440 (QHD)', width: '2560', height: '1440' },
    { label: '2560x1600 (WQXGA)', width: '2560', height: '1600' },
    { label: '3840x2160 (4K UHD)', width: '3840', height: '2160' },
    { label: '4096x2160 (DCI 4K)', width: '4096', height: '2160' },
    { label: '7680x4320 (8K UHD)', width: '7680', height: '4320' },
  ];

  useEffect(() => {
    const fetchCameras = async () => {
      try {
        setLoading(true);
        setError(null);
        const data = await api.getCameras();
        setCameras(data);
      } catch (err) {
        console.error('Error fetching cameras:', err);
        setError('Failed to fetch cameras.');
      } finally {
        setLoading(false);
      }
    };

    fetchCameras();
  }, []);

  const handleCameraCheckboxChange = (cameraId: string) => {
    setSelectedCameraIds(prevSelected =>
      prevSelected.includes(cameraId)
        ? prevSelected.filter(id => id !== cameraId)
        : [...prevSelected, cameraId]
    );
  };

  const handleSelectAll = () => {
    if (selectedCameraIds.length === cameras.length) {
      setSelectedCameraIds([]);
    } else {
      setSelectedCameraIds(cameras.map(camera => camera.id));
    }
  };

  const handleResolutionPresetChange = (event: any) => {
    const presetLabel = event.target.value as string;
    setSelectedResolutionPreset(presetLabel);
    const selectedPreset = resolutionPresets.find(preset => preset.label === presetLabel);
    if (selectedPreset) {
      setConfigWidth(selectedPreset.width);
      setConfigHeight(selectedPreset.height);
    }
  };

  const handleApplyConfig = async () => {
    if (selectedCameraIds.length === 0) {
      setConfigError('Please select at least one camera.');
      setConfigSuccess(null);
      return;
    }

    setConfigLoading(true);
    setConfigError(null);
    setConfigSuccess(null);

    const payload: ApplyConfigPayload = {
      // These will be updated per camera in the loop
      cameraId: '',
      width: parseInt(configWidth, 10),
      height: parseInt(configHeight, 10),
      fps: parseInt(configFps, 10),
    };

    const results = [];

    for (const cameraId of selectedCameraIds) {
      try {
        // Update the payload with the current camera ID
        payload.cameraId = cameraId;
        const result = await api.applyConfig(payload);
        results.push({ cameraId, status: result.status, success: result.status === 'configuration applied' });
      } catch (err) {
        console.error(`Error applying config to camera ${cameraId}:`, err);
        results.push({ cameraId, status: 'Error', success: false });
      }
    }

    setConfigLoading(false);

    const successCount = results.filter(r => r.success).length;
    if (successCount === selectedCameraIds.length) {
      setConfigSuccess(`Successfully updated configuration for all ${successCount} selected camera(s).`);
    } else if (successCount > 0) {
      setConfigSuccess(`Successfully updated configuration for ${successCount} out of ${selectedCameraIds.length} selected camera(s). Some updates may have failed.`);
      setConfigError('Some configuration updates failed. Check console for details.');
    } else {
      setConfigError('Failed to update configuration for any selected cameras. Check console for details.');
    }
  };

  return (
    <Container maxWidth="md" sx={{ mt: 4, bgcolor: '#333', p: 3, borderRadius: 1 }}>
      <Typography variant="h4" gutterBottom component="h1" sx={{ color: 'white', textAlign: 'center', mb: 4 }}>
        ONVIF Camera Control
      </Typography>

      {/* Arrange sections side-by-side using Grid */}
      <Grid container spacing={4}>
        {/* Configuration Section - Now on the left */}
        <Grid item xs={12} md={6}>
          <Typography variant="h5" gutterBottom sx={{ color: 'white' }}>Camera Settings</Typography>
          <Paper elevation={2} sx={{ mt: 2, p: 2, bgcolor: '#444', color: 'white' }}>
            {/* Inner layout for form elements using Box/flexbox */}
            <Box display="flex" flexDirection="row" flexWrap="wrap" gap={2} alignItems="center">
              {/* Added Resolution Preset Select */}
              <Box sx={{ flexGrow: 1, minWidth: '180px' }}>
                <FormControl fullWidth>
                  <InputLabel id="resolution-preset-label" sx={{ color: '#bbb' }}>Resolution</InputLabel>
                  <Select
                    labelId="resolution-preset-label"
                    id="resolution-preset-select"
                    value={selectedResolutionPreset}
                    label="Resolution"
                    onChange={handleResolutionPresetChange}
                    sx={{ 
                      color: 'white', 
                      '.MuiOutlinedInput-notchedOutline': { borderColor: '#555' },
                      '&:hover .MuiOutlinedInput-notchedOutline': { borderColor: '#777' },
                      '&.Mui-focused .MuiOutlinedInput-notchedOutline': { borderColor: '#00aaff' },
                      '.MuiSvgIcon-root ': { fill: 'white' }
                    }}
                    MenuProps={{
                      PaperProps: {
                        sx: {
                          bgcolor: '#444',
                          color: 'white',
                        },
                      },
                    }}
                  >
                    {resolutionPresets.map((preset) => (
                      <MenuItem key={preset.label} value={preset.label} sx={{ '&:hover': { bgcolor: '#555' } }}>
                        {preset.label}
                      </MenuItem>
                    ))}
                  </Select>
                </FormControl>
              </Box>

              <Box sx={{ flexGrow: 1, minWidth: '120px' }}>
                <TextField
                  label="Frame Rate (fps)"
                  type="number"
                  fullWidth
                  value={configFps}
                  onChange={(e) => setConfigFps(e.target.value)}
                  InputProps={{ inputProps: { min: 1 }, sx: { color: 'white' } }}
                  InputLabelProps={{ sx: { color: '#bbb' } }}
                  sx={{
                    '.MuiOutlinedInput-notchedOutline': { borderColor: '#555' },
                    '&:hover .MuiOutlinedInput-notchedOutline': { borderColor: '#777' },
                    '&.Mui-focused .MuiOutlinedInput-notchedOutline': { borderColor: '#00aaff' },
                  }}
                />
              </Box>

              <Box sx={{ width: '100%', mt: { xs: 1, sm: 0 } }}>
                <Button
                  variant="contained"
                  color="primary"
                  onClick={handleApplyConfig}
                  disabled={selectedCameraIds.length === 0 || configLoading}
                  fullWidth
                  sx={{ bgcolor: '#00aaff', '&:hover': { bgcolor: '#0099cc' } }}
                >
                  {configLoading ? 'APPLYING...' : 'APPLY SETTINGS'}
                </Button>
              </Box>
            </Box>

            {configLoading && (
              <Box display="flex" justifyContent="center" sx={{ mt: 2 }}>
                <CircularProgress size={24} sx={{ color: '#00aaff' }} />
              </Box>
            )}

            {configError && (
              <Alert severity="error" sx={{ mt: 2, bgcolor: '#553c3c', color: 'white' }}>
                {configError}
              </Alert>
            )}

            {configSuccess && (
              <Alert severity="success" sx={{ mt: 2, bgcolor: '#3c553c', color: 'white' }}>
                {configSuccess}
              </Alert>
            )}

          </Paper>
        </Grid>

        {/* Camera List Section - Now on the right */}
        <Grid item xs={12} md={6}>
          <Typography variant="h5" gutterBottom sx={{ color: 'white' }}>Camera List</Typography>
          <Box display="flex" justifyContent="space-between" alignItems="center" sx={{ mb: 1 }}>
             <Button onClick={handleSelectAll} size="small" sx={{ color: '#00aaff' }}>
              {selectedCameraIds.length === cameras.length ? 'Deselect All' : 'SELECT ALL'}
            </Button>
            <Button variant="contained" size="small" startIcon={<i className="fas fa-plus"></i>} sx={{ bgcolor: '#00aaff', '&:hover': { bgcolor: '#0099cc' } }}>
              + ADD CAMERA
            </Button>
          </Box>

          {loading && (
            <Box display="flex" justifyContent="center" sx={{ mt: 2 }}>
              <CircularProgress sx={{ color: '#00aaff' }}/>
            </Box>
          )}

          {error && (
            <Alert severity="error" sx={{ mt: 2, bgcolor: '#553c3c', color: 'white' }}>
              {error}
            </Alert>
          )}

          {!loading && !error && cameras.length === 0 && (
            <Alert severity="info" sx={{ mt: 2, bgcolor: '#3c3c55', color: 'white' }}>
              No cameras found.
            </Alert>
          )}

          {!loading && !error && cameras.length > 0 && (
            <Paper elevation={2} sx={{ mt: 2, p: 2, bgcolor: '#444', color: 'white' }}>
              <FormGroup>
                {cameras.map(camera => (
                  <FormControlLabel
                    key={camera.id}
                    control={
                      <Checkbox
                        checked={selectedCameraIds.includes(camera.id)}
                        onChange={() => handleCameraCheckboxChange(camera.id)}
                        sx={{ color: '#00aaff', '&.Mui-checked': { color: '#00aaff' } }}
                      />
                    }
                    label={
                      <Typography sx={{ color: 'white' }}>
                        {camera.isFake ? '(Simulated) ' : ''}
                        Camera {camera.id} - {camera.ip}
                      </Typography>
                    }
                  />
                ))}
              </FormGroup>
            </Paper>
          )}
        </Grid>
      </Grid>
    </Container>
  );
}

export default App;
