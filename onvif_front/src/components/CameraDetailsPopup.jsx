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
  DialogActions,
} from '@mui/material';
import CloseIcon from '@mui/icons-material/Close';
import ContentCopyIcon from '@mui/icons-material/ContentCopy';
import DeleteIcon from '@mui/icons-material/Delete';
import api from '../services/api';

const CameraDetailsPopup = ({ open, onClose, camera }) => {
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [streamUrl, setStreamUrl] = useState('');
  const [configDetails, setConfigDetails] = useState(null);
  const [deviceInfo, setDeviceInfo] = useState(null);
  const [copySuccess, setCopySuccess] = useState(false);
  const [launching, setLaunching] = useState(false);
  const [deleting, setDeleting] = useState(false);

  useEffect(() => {
    if (open && !camera.isFake) {
      loadData();
    } else if (open && camera.isFake) {
      // Reset states for fake camera
      setLoading(false);
      setError('');
      setStreamUrl('');
      setConfigDetails(null);
      setDeviceInfo(null);
    }
  }, [open, camera.isFake]);

  const handleCopyClick = async () => {
    try {
      await navigator.clipboard.writeText(streamUrl);
      setCopySuccess(true);
    } catch (err) {
      console.error('Failed to copy text: ', err);
    }
  };

  const handleLaunchVLC = async () => {
    setLaunching(true);
    try {
      await api.launchVLC(camera.id);
    } catch (err) {
      console.error('Failed to launch VLC:', err);
      setError('Failed to launch VLC. Is VLC installed on the server?');
    } finally {
      setLaunching(false);
    }
  };

  const handleDelete = async () => {
    if (!window.confirm(`Are you sure you want to delete Camera ${camera.id}?`)) {
      return;
    }
    
    setDeleting(true);
    try {
      await api.deleteCamera(camera.id);
      onClose();
      // Force page refresh to update camera list
      window.location.reload();
    } catch (err) {
      console.error('Failed to delete camera:', err);
      setError('Failed to delete camera: ' + err.message);
    } finally {
      setDeleting(false);
    }
  };

  const loadData = async () => {
    setLoading(true);
    setError('');
    try {
      if (!camera.isFake) {
        const details = await api.getCameraDetailsSimple(camera.id);
        setStreamUrl(details.streamUrl);
        setConfigDetails(details.config);
        setDeviceInfo(details.deviceInfo);
      }
    } catch (err) {
      console.error('Error loading camera details:', err);
      setError('Failed to load camera details. ' + err.message);
    } finally {
      setLoading(false);
    }
  };

  const renderFakeCameraInfo = () => {
    return (
      <Box sx={{ mt: 2 }}>
        <Typography variant="h6" gutterBottom>Simulated Camera Information</Typography>
        <Paper sx={{ p: 2 }}>
          <Typography variant="body2" sx={{ mb: 1 }}>
            Camera ID: {camera.id}
          </Typography>
          <Typography variant="body2" sx={{ mb: 1 }}>
            IP Address: {camera.ip}
          </Typography>
          <Typography variant="body2" sx={{ mb: 1 }}>
            Type: Simulated ONVIF Camera
          </Typography>
          <Alert severity="info" sx={{ mt: 2 }}>
            This is a simulated camera. Streaming and detailed configuration are not available.
          </Alert>
        </Paper>
      </Box>
    );
  };

  const renderConfigDetails = () => {
    if (!configDetails) return null;

    return (
      <Box sx={{ mt: 2 }}>
        <Typography variant="h6" gutterBottom>Configuration Details</Typography>
        <Paper sx={{ p: 2 }}>
          <Typography variant="body2" sx={{ mb: 1 }}>
            Resolution: {configDetails.Resolution?.Width || configDetails.Width}x{configDetails.Resolution?.Height || configDetails.Height}
          </Typography>
          <Typography variant="body2" sx={{ mb: 1 }}>
            Frame Rate: {configDetails.RateControl?.FrameRateLimit || configDetails.FrameRate || 0} fps
          </Typography>
          <Typography variant="body2" sx={{ mb: 1 }}>
            Bit Rate: {configDetails.RateControl?.BitRateLimit || configDetails.BitRate || 0} kbps
          </Typography>
          <Typography variant="body2" sx={{ mb: 1 }}>
            Encoding: {configDetails.Encoding || 'H264'}
          </Typography>
          {configDetails.H264?.GovLength || configDetails.GovLength ? (
            <Typography variant="body2" sx={{ mb: 1 }}>
              GOP Length: {configDetails.H264?.GovLength || configDetails.GovLength}
            </Typography>
          ) : null}
        </Paper>
      </Box>
    );
  };

  const renderDeviceInfo = () => {
    if (!deviceInfo) return null;

    return (
      <Box sx={{ mt: 2 }}>
        <Typography variant="h6" gutterBottom>Device Information</Typography>
        <Paper sx={{ p: 2 }}>
          <Typography variant="body2" sx={{ mb: 1 }}>Manufacturer: {deviceInfo.Manufacturer}</Typography>
          <Typography variant="body2" sx={{ mb: 1 }}>Model: {deviceInfo.Model}</Typography>
          <Typography variant="body2" sx={{ mb: 1 }}>Firmware Version: {deviceInfo.FirmwareVersion}</Typography>
          <Typography variant="body2" sx={{ mb: 1 }}>Serial Number: {deviceInfo.SerialNumber}</Typography>
          <Typography variant="body2">Hardware ID: {deviceInfo.HardwareId}</Typography>
        </Paper>
      </Box>
    );
  };

  return (
    <>
      <Dialog
        open={open}
        onClose={onClose}
        maxWidth="md"
        fullWidth
      >
        <DialogTitle>
          <Box display="flex" justifyContent="space-between" alignItems="center">
            <Typography variant="h6">
              Camera {camera.id} Details {camera.isFake ? '(Simulated)' : ''}
            </Typography>
            <IconButton
              aria-label="close"
              onClick={onClose}
              size="small"
            >
              <CloseIcon />
            </IconButton>
          </Box>
        </DialogTitle>
        <DialogContent>
          {loading ? (
            <Box display="flex" justifyContent="center" my={3}>
              <CircularProgress />
            </Box>
          ) : error ? (
            <Alert severity="error">{error}</Alert>
          ) : (
            <>
              {camera.isFake ? (
                renderFakeCameraInfo()
              ) : (
                <>
                  {streamUrl && (
                    <Box sx={{ mt: 2 }}>
                      <Typography variant="h6" gutterBottom>Stream URL</Typography>
                      <Paper sx={{ p: 2, display: 'flex', alignItems: 'center', gap: 2 }}>
                        <Typography variant="body2" sx={{ flexGrow: 1, wordBreak: 'break-all' }}>
                          {streamUrl}
                        </Typography>
                        <IconButton onClick={handleCopyClick} size="small">
                          <ContentCopyIcon />
                        </IconButton>
                      </Paper>
                      <Box sx={{ mt: 2, display: 'flex', gap: 2 }}>
                        <Button
                          variant="contained"
                          onClick={handleLaunchVLC}
                          disabled={launching}
                        >
                          {launching ? 'Launching...' : 'Open in VLC'}
                        </Button>
                      </Box>
                    </Box>
                  )}
                  {renderConfigDetails()}
                  {renderDeviceInfo()}
                </>
              )}
            </>
          )}
        </DialogContent>
        <DialogActions sx={{ p: 2, pt: 0 }}>
          <Button
            variant="contained"
            color="error"
            startIcon={<DeleteIcon />}
            onClick={handleDelete}
            disabled={deleting}
          >
            {deleting ? 'Deleting...' : 'Delete Camera'}
          </Button>
        </DialogActions>
      </Dialog>
      <Snackbar
        open={copySuccess}
        autoHideDuration={3000}
        onClose={() => setCopySuccess(false)}
        message="Stream URL copied to clipboard"
      />
    </>
  );
};

export default CameraDetailsPopup;
