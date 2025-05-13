import React, { useState, useEffect } from 'react';
import {
  Dialog,
  DialogContent,
  DialogTitle,
  Box,
  Typography,
  Alert,
  Button,
  IconButton,
  CircularProgress,
  Snackbar,
  Paper,
} from '@mui/material';
import CloseIcon from '@mui/icons-material/Close';
import ContentCopyIcon from '@mui/icons-material/ContentCopy';
import api from '../services/api';

const CameraDetailsPopup = ({ open, onClose, camera, selectedProfile, selectedConfig }) => {
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [streamUrl, setStreamUrl] = useState('');
  const [configDetails, setConfigDetails] = useState(null);
  const [deviceInfo, setDeviceInfo] = useState(null);
  const [copySuccess, setCopySuccess] = useState(false);
  const [showDeleteConfirm, setShowDeleteConfirm] = useState(false);

  useEffect(() => {
    if (open && selectedProfile && !camera.isFake) {
      loadData();
    }
  }, [open, selectedProfile, camera.isFake]);

  const handleCopyClick = async () => {
    try {
      await navigator.clipboard.writeText(streamUrl);
      setCopySuccess(true);
      setTimeout(() => setCopySuccess(false), 3000);
    } catch (err) {
      console.error('Failed to copy text: ', err);
    }
  };

  const loadData = async () => {
    setLoading(true);
    setError('');
    try {
      if (!camera.isFake) {
        const url = await api.getStreamUrl(selectedProfile);
        setStreamUrl(url);

        if (selectedConfig) {
          const configData = await api.getSingleConfig(selectedConfig);
          setConfigDetails(configData);
        }

        const deviceData = await api.getDeviceInfo();
        setDeviceInfo(deviceData);
      }
      setLoading(false);
    } catch (err) {
      console.error('Error loading camera details:', err);
      setError('Failed to load camera details');
      setLoading(false);
    }
  };

  const handleLaunchVLC = async () => {
    if (camera.isFake) {
      setError('Cannot launch VLC for a simulated camera');
      return;
    }
    try {
      const response = await api.launchVLC(selectedProfile);
      console.log('VLC response:', response);
    } catch (error) {
      console.error('Error launching VLC:', error);
      setError('Failed to launch VLC');
    }
  };
  const handleDelete = async () => {
    try {
      await api.deleteCamera(camera.id);
      setShowDeleteConfirm(false);
      onClose('deleted');
    } catch (error) {
      setError('Failed to delete camera. Please try again.');
    }
  };

  return (
    <>
      <Dialog open={open} onClose={() => onClose()} maxWidth="md" fullWidth>
        <DialogTitle>
          <Box display="flex" justifyContent="space-between" alignItems="center">
            <Typography variant="h6">Camera Details - {camera?.ip}</Typography>
            <Box>
              <Button
                variant="contained"
                color="error"
                size="small"
                onClick={() => setShowDeleteConfirm(true)}
                sx={{ mr: 1 }}
              >
                Delete Camera
              </Button>
              <IconButton onClick={() => onClose()} size="small">
                <CloseIcon />
              </IconButton>
            </Box>
          </Box>
        </DialogTitle>
        <DialogContent>
          {camera.isFake ? (
            <Alert severity="info" sx={{ mb: 4 }}>
              This is a simulated camera. Stream URL and detailed information are not available.
            </Alert>
          ) : loading ? (
            <Box display="flex" justifyContent="center" my={4}>
              <CircularProgress />
            </Box>
          ) : error ? (
            <Alert severity="error">{error}</Alert>
          ) : (
            <Box>
              {/* Stream Section */}
              <Box mb={4}>
                <Typography variant="h6" color="primary" gutterBottom>
                  Camera Stream
                </Typography>
                <Paper
                  variant="outlined"
                  sx={{
                    p: 2,
                    mb: 2,
                    backgroundColor: 'rgba(0, 0, 0, 0.2)',
                    borderRadius: 1,
                    fontFamily: 'monospace',
                    position: 'relative',
                    '&:hover .copy-button': {
                      opacity: 1,
                    },
                  }}
                >
                  <Box sx={{ 
                    maxWidth: '100%', 
                    overflowX: 'auto',
                    whiteSpace: 'nowrap',
                    color: '#e6e6e6',
                    pr: 4
                  }}>
                    {streamUrl}
                  </Box>
                  <IconButton
                    className="copy-button"
                    onClick={handleCopyClick}
                    sx={{
                      position: 'absolute',
                      right: 8,
                      top: '50%',
                      transform: 'translateY(-50%)',
                      opacity: 0.6,
                      transition: 'opacity 0.2s',
                      backgroundColor: 'rgba(255, 255, 255, 0.1)',
                      '&:hover': {
                        backgroundColor: 'rgba(255, 255, 255, 0.2)',
                      },
                    }}
                    size="small"
                  >
                    <ContentCopyIcon fontSize="small" />
                  </IconButton>
                </Paper>
                <Box sx={{ display: 'flex', gap: 2, alignItems: 'center' }}>
                  <Button
                    variant="contained"
                    color="primary"
                    onClick={handleLaunchVLC}
                    disabled={camera.isFake}
                  >
                    Open in VLC
                  </Button>
                  {camera.isFake && (
                    <Typography variant="caption" color="text.secondary">
                      VLC streaming is not available for simulated cameras
                    </Typography>
                  )}
                </Box>
              </Box>

              {/* Configuration Section */}
              {!camera.isFake && configDetails && (
                <Box mb={4}>
                  <Typography variant="h6" color="primary" gutterBottom>
                    Current Configuration
                  </Typography>
                  <Typography variant="body2">
                    Resolution: {configDetails.Width || configDetails.width || 'N/A'} x {configDetails.Height || configDetails.height || 'N/A'}
                  </Typography>
                  <Typography variant="body2">
                    Frame Rate: {configDetails.RateControl?.FrameRateLimit || configDetails.FrameRate || configDetails.frameRate || 'N/A'} fps
                  </Typography>
                  <Typography variant="body2">
                    Bit Rate: {configDetails.RateControl?.BitrateLimit || configDetails.BitRate || configDetails.bitRate || 'N/A'} kbps
                  </Typography>
                  <Typography variant="body2">
                    GOP Length: {configDetails.H264?.GovLength || configDetails.govLength || configDetails.GovLength || 'N/A'}
                  </Typography>
                  <Typography variant="body2">
                    H264 Profile: {configDetails.H264?.H264Profile || configDetails.h264Profile || configDetails.H264Profile || 'N/A'}
                  </Typography>
                </Box>
              )}

              {/* Device Info Section */}
              <Box>
                <Typography variant="h6" color="primary" gutterBottom>
                  Device Information
                </Typography>
                {!camera.isFake && deviceInfo ? (
                  <>
                    <Typography variant="body2">
                      Manufacturer: {deviceInfo.manufacturer}
                    </Typography>
                    <Typography variant="body2">
                      Model: {deviceInfo.model}
                    </Typography>
                    <Typography variant="body2">
                      Firmware Version: {deviceInfo.firmwareVersion}
                    </Typography>
                    <Typography variant="body2">
                      Serial Number: {deviceInfo.serialNumber}
                    </Typography>
                    <Typography variant="body2">
                      Hardware ID: {deviceInfo.hardwareId}
                    </Typography>
                  </>
                ) : (
                  <>
                    <Typography variant="body2">
                      IP Address: {camera.ip}
                    </Typography>
                    <Typography variant="body2">
                      Type: Simulated Camera
                    </Typography>
                    <Typography variant="body2">
                      Status: Active
                    </Typography>
                  </>
                )}
              </Box>
            </Box>
          )}
        </DialogContent>
      </Dialog>

      <Dialog
        open={showDeleteConfirm}
        onClose={() => setShowDeleteConfirm(false)}
        maxWidth="xs"
        fullWidth
      >
        <DialogTitle>Confirm Delete</DialogTitle>
        <DialogContent>
          <Typography>
            Are you sure you want to delete this camera? This action cannot be undone.
          </Typography>
          <Box sx={{ display: 'flex', justifyContent: 'flex-end', mt: 2, gap: 1 }}>
            <Button onClick={() => setShowDeleteConfirm(false)}>Cancel</Button>
            <Button variant="contained" color="error" onClick={handleDelete}>
              Delete
            </Button>
          </Box>
        </DialogContent>
      </Dialog>

      <Snackbar
        open={copySuccess}
        autoHideDuration={3000}
        onClose={() => setCopySuccess(false)}
        message="Stream URL copied to clipboard"
        anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}
      />
    </>
  );
};

export default CameraDetailsPopup;
