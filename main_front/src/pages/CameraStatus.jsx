import React, { useState, useEffect } from 'react';
import {
  Box,
  Paper,
  Typography,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Chip,
  IconButton,
  Button,
  Tooltip,
  CircularProgress,
  Alert,
  Card,
  CardContent,
  Snackbar
} from '@mui/material';
import {
  Refresh as RefreshIcon,
  PlayArrow as PlayIcon,
  Circle as CircleIcon,
  Info as InfoIcon
} from '@mui/icons-material';
import { checkAllCameras, launchVLC } from '../services/api';

const CameraStatus = () => {
  const [cameras, setCameras] = useState([]);
  const [loading, setLoading] = useState(false);
  const [refreshing, setRefreshing] = useState(false);
  const [error, setError] = useState(null);
  const [snackbar, setSnackbar] = useState({ open: false, message: '', severity: 'info' });
  const [summary, setSummary] = useState({ online: 0, offline: 0, partial: 0, unknown: 0 });

  // Load cameras and their status from CSV
  const loadCamerasAndStatus = async () => {
    try {
      setLoading(true);
      setError(null);
      
      // Use the new check all cameras endpoint that loads from CSV
      const result = await checkAllCameras();
      setCameras(result.cameras || []);
      setSummary(result.summary || { online: 0, offline: 0, partial: 0, unknown: 0 });
    } catch (err) {
      console.error('Error checking cameras:', err);
      setError(err.message || 'Failed to check cameras');
      setCameras([]);
      setSummary({ online: 0, offline: 0, partial: 0, unknown: 0 });
    } finally {
      setLoading(false);
    }
  };

  const handleCheckAll = async () => {
    setRefreshing(true);
    try {
      await loadCamerasAndStatus();
      showSnackbar('Camera check completed successfully', 'success');
    } catch (err) {
      showSnackbar('Failed to check cameras', 'error');
    } finally {
      setRefreshing(false);
    }
  };

  const handleLaunchVLC = async (cameraId) => {
    try {
      const result = await launchVLC(cameraId);
      showSnackbar(result.message || 'VLC launched successfully', 'success');
    } catch (err) {
      showSnackbar(`Failed to launch VLC: ${err.message}`, 'error');
    }
  };

  const showSnackbar = (message, severity = 'info') => {
    setSnackbar({ open: true, message, severity });
  };

  const handleCloseSnackbar = () => {
    setSnackbar({ ...snackbar, open: false });
  };

  const getStatusChip = (camera) => {
    if (!camera.status) {
      return <Chip size="small" label="Unknown" color="default" />;
    }

    switch (camera.status) {
      case 'online':
        return (
          <Chip
            size="small"
            label="Online"
            color="success"
            icon={<CircleIcon />}
          />
        );
      case 'partial':
        return (
          <Chip
            size="small"
            label="Partial"
            color="warning"
            icon={<CircleIcon />}
          />
        );
      case 'offline':
        return (
          <Chip
            size="small"
            label="Offline"
            color="error"
            icon={<CircleIcon />}
          />
        );
      default:
        return (
          <Chip
            size="small"
            label="Unknown"
            color="default"
            icon={<CircleIcon />}
          />
        );
    }
  };

  const getResolution = (camera) => {
    if (!camera.currentConfig?.resolution) {
      return 'N/A';
    }
    const res = camera.currentConfig.resolution;
    return `${res.width}x${res.height}`;
  };

  const getEncoding = (camera) => {
    return camera.currentConfig?.encoding || 'N/A';
  };

  const getFPS = (camera) => {
    return camera.currentConfig?.fps || 'N/A';
  };

  const getBitrate = (camera) => {
    return camera.currentConfig?.bitrate ? `${camera.currentConfig.bitrate} kbps` : 'N/A';
  };

  const getErrorTooltip = (camera) => {
    return camera.error || '';
  };

  // Load cameras and their status on component mount
  useEffect(() => {
    loadCamerasAndStatus();
  }, []);

  if (loading) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight="400px">
        <CircularProgress />
        <Typography variant="h6" sx={{ ml: 2 }}>
          Loading camera status...
        </Typography>
      </Box>
    );
  }

  if (error) {
    return (
      <Box p={3}>
        <Alert severity="error" sx={{ mb: 2 }}>
          {error}
        </Alert>
        <Button variant="contained" onClick={loadCamerasAndStatus}>
          Retry
        </Button>
      </Box>
    );
  }

  return (
    <Box p={3}>
      <Card sx={{ mb: 3 }}>
        <CardContent>
          <Box display="flex" justifyContent="space-between" alignItems="center">
            <Typography variant="h5" component="h1">
              Camera Status Monitor
            </Typography>
            <Button
              variant="contained"
              startIcon={refreshing ? <CircularProgress size={20} color="inherit" /> : <RefreshIcon />}
              onClick={handleCheckAll}
              disabled={refreshing}
            >
              {refreshing ? 'Checking...' : 'Check All Cameras'}
            </Button>
          </Box>
          <Typography variant="body2" color="text.secondary" sx={{ mt: 1 }}>
            Check all cameras from CSV file and view their current status and configuration
          </Typography>
          
          {/* Status Summary */}
          <Box sx={{ mt: 2, display: 'flex', gap: 2, flexWrap: 'wrap' }}>
            <Chip 
              size="small" 
              label={`Online: ${summary.online}`} 
              color="success" 
              variant="outlined"
            />
            <Chip 
              size="small" 
              label={`Offline: ${summary.offline}`} 
              color="error" 
              variant="outlined"
            />
            <Chip 
              size="small" 
              label={`Partial: ${summary.partial}`} 
              color="warning" 
              variant="outlined"
            />
            <Chip 
              size="small" 
              label={`Unknown: ${summary.unknown}`} 
              color="default" 
              variant="outlined"
            />
          </Box>
        </CardContent>
      </Card>

      {cameras.length === 0 ? (
        <Alert severity="info">
          No cameras found. Add some cameras first to monitor their status.
        </Alert>
      ) : (
        <TableContainer component={Paper}>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>Camera ID</TableCell>
                <TableCell>IP Address</TableCell>
                <TableCell>Status</TableCell>
                <TableCell>Resolution</TableCell>
                <TableCell>Encoding</TableCell>
                <TableCell>FPS</TableCell>
                <TableCell>Bitrate</TableCell>
                <TableCell>Type</TableCell>
                <TableCell align="center">Actions</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {cameras.map((camera) => {
                const hasError = camera.status === 'offline' || camera.status === 'unknown';
                
                return (
                  <TableRow key={camera.cameraId} hover>
                    <TableCell>
                      <Typography variant="body2" fontWeight="bold">
                        {camera.cameraId}
                      </Typography>
                    </TableCell>
                    <TableCell>
                      <Typography variant="body2">
                        {camera.ip}:{camera.port}
                      </Typography>
                    </TableCell>
                    <TableCell>
                      <Box display="flex" alignItems="center" gap={1}>
                        {getStatusChip(camera)}
                        {hasError && camera.error && (
                          <Tooltip title={getErrorTooltip(camera)} arrow>
                            <InfoIcon color="error" fontSize="small" />
                          </Tooltip>
                        )}
                      </Box>
                    </TableCell>
                    <TableCell>
                      <Typography variant="body2">
                        {getResolution(camera)}
                      </Typography>
                    </TableCell>
                    <TableCell>
                      <Typography variant="body2">
                        {getEncoding(camera)}
                      </Typography>
                    </TableCell>
                    <TableCell>
                      <Typography variant="body2">
                        {getFPS(camera)}
                      </Typography>
                    </TableCell>
                    <TableCell>
                      <Typography variant="body2">
                        {getBitrate(camera)}
                      </Typography>
                    </TableCell>
                    <TableCell>
                      <Chip
                        size="small"
                        label={camera.isFake ? 'Simulated' : 'Real'}
                        color={camera.isFake ? 'warning' : 'primary'}
                        variant="outlined"
                      />
                    </TableCell>
                    <TableCell align="center">
                      <Tooltip title="Launch VLC Player">
                        <IconButton
                          color="primary"
                          onClick={() => handleLaunchVLC(camera.cameraId)}
                          disabled={hasError}
                        >
                          <PlayIcon />
                        </IconButton>
                      </Tooltip>
                    </TableCell>
                  </TableRow>
                );
              })}
            </TableBody>
          </Table>
        </TableContainer>
      )}

      <Snackbar
        open={snackbar.open}
        autoHideDuration={6000}
        onClose={handleCloseSnackbar}
      >
        <Alert
          onClose={handleCloseSnackbar}
          severity={snackbar.severity}
          sx={{ width: '100%' }}
        >
          {snackbar.message}
        </Alert>
      </Snackbar>
    </Box>
  );
};

export default CameraStatus;
