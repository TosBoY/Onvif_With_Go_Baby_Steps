import React, { useState, useEffect, useMemo, useCallback } from 'react';
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
  Snackbar,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  MenuItem,
  FormControl,
  InputLabel,
  Select,
  Pagination,
  Stack
} from '@mui/material';
import {
  Refresh as RefreshIcon,
  PlayArrow as PlayIcon,
  Circle as CircleIcon,
  Info as InfoIcon,
  Settings as SettingsIcon,
  CheckCircle as CheckCircleIcon,
  Error as ErrorIcon,
  Warning as WarningIcon
} from '@mui/icons-material';
import { loadCameraList, checkSingleCamera, configureSingleCamera, validateSingleCamera, launchVLC } from '../services/api';

const CameraStatus = () => {
  const [cameras, setCameras] = useState([]);
  const [cameraStatuses, setCameraStatuses] = useState({});
  const [cameraValidations, setCameraValidations] = useState({});
  const [loading, setLoading] = useState(false);
  const [refreshing, setRefreshing] = useState(false);
  const [checking, setChecking] = useState({});
  const [validating, setValidating] = useState({});
  
  // Pagination state
  const [currentPage, setCurrentPage] = useState(1);
  const [pageJumpValue, setPageJumpValue] = useState('');
  const camerasPerPage = 50;
  
  const [error, setError] = useState(null);
  const [snackbar, setSnackbar] = useState({ open: false, message: '', severity: 'info' });
  
  // Configuration dialog state
  const [configDialog, setConfigDialog] = useState({
    open: false,
    cameraId: '',
    cameraInfo: null,
    loading: false
  });
  const [configForm, setConfigForm] = useState({
    resolution: '1920x1080',
    fps: 30,
    bitrate: 2000,
    encoding: 'h264'
  });

  // Optimized pagination helpers with memoization
  const { totalPages, startIndex, endIndex, paginatedCameras } = useMemo(() => {
    const total = Math.ceil(cameras.length / camerasPerPage);
    const start = (currentPage - 1) * camerasPerPage;
    const end = start + camerasPerPage;
    const paginated = cameras.slice(start, end);
    
    return {
      totalPages: total,
      startIndex: start,
      endIndex: end,
      paginatedCameras: paginated
    };
  }, [cameras, currentPage, camerasPerPage]);

  const handlePageChange = useCallback((event, value) => {
    setCurrentPage(value);
    setPageJumpValue(''); // Clear jump input when page changes
  }, []);

  const handlePageJump = useCallback((event) => {
    if (event.key === 'Enter') {
      const pageNum = parseInt(pageJumpValue);
      if (pageNum >= 1 && pageNum <= totalPages) {
        setCurrentPage(pageNum);
        setPageJumpValue('');
      } else {
        showSnackbar(`Please enter a page number between 1 and ${totalPages}`, 'warning');
      }
    }
  }, [pageJumpValue, totalPages]);

  // Load cameras from CSV
  const loadCameras = async () => {
    try {
      setLoading(true);
      setError(null);
      
      // Load camera list from CSV
      const result = await loadCameraList();
      setCameras(result.cameras || []);
      
      // Reset pagination to first page
      setCurrentPage(1);
      setPageJumpValue('');
      
      // Initialize camera statuses as unknown
      const initialStatuses = {};
      (result.cameras || []).forEach(camera => {
        initialStatuses[camera.cameraId] = {
          status: 'unknown',
          lastChecked: null
        };
      });
      setCameraStatuses(initialStatuses);
      
    } catch (err) {
      console.error('Error loading cameras:', err);
      setError(err.message || 'Failed to load cameras');
      setCameras([]);
      setCameraStatuses({});
    } finally {
      setLoading(false);
    }
  };

  // Check a single camera
  const checkCamera = async (cameraId) => {
    try {
      setChecking(prev => ({ ...prev, [cameraId]: true }));
      
      const result = await checkSingleCamera(cameraId);
      setCameraStatuses(prev => ({
        ...prev,
        [cameraId]: {
          ...result,
          lastChecked: new Date().toLocaleTimeString()
        }
      }));
      
      showSnackbar(`Camera ${cameraId} checked successfully`, 'success');
    } catch (err) {
      console.error(`Error checking camera ${cameraId}:`, err);
      setCameraStatuses(prev => ({
        ...prev,
        [cameraId]: {
          status: 'error',
          error: err.message,
          lastChecked: new Date().toLocaleTimeString()
        }
      }));
      showSnackbar(`Failed to check camera ${cameraId}: ${err.message}`, 'error');
    } finally {
      setChecking(prev => ({ ...prev, [cameraId]: false }));
    }
  };

  // Validate a single camera using RTSP analyzer
  const validateCamera = async (cameraId) => {
    try {
      setValidating(prev => ({ ...prev, [cameraId]: true }));
      
      const result = await validateSingleCamera(cameraId);
      setCameraValidations(prev => ({
        ...prev,
        [cameraId]: {
          ...result,
          lastValidated: new Date().toLocaleTimeString()
        }
      }));
      
      const status = result.isValid ? 'success' : 'warning';
      const message = result.isValid 
        ? `Camera ${cameraId} validation passed` 
        : `Camera ${cameraId} validation completed with issues`;
      showSnackbar(message, status);
    } catch (err) {
      console.error(`Error validating camera ${cameraId}:`, err);
      setCameraValidations(prev => ({
        ...prev,
        [cameraId]: {
          isValid: false,
          error: err.message,
          lastValidated: new Date().toLocaleTimeString()
        }
      }));
      showSnackbar(`Failed to validate camera ${cameraId}: ${err.message}`, 'error');
    } finally {
      setValidating(prev => ({ ...prev, [cameraId]: false }));
    }
  };

  // Check all cameras
  const checkAllCameras = useCallback(async () => {
    setRefreshing(true);
    try {
      // Check all cameras, not just the ones on current page
      const checkPromises = cameras.map(camera => checkCamera(camera.cameraId));
      await Promise.all(checkPromises);
      showSnackbar(`All ${cameras.length} cameras checked`, 'success');
    } catch (err) {
      showSnackbar('Some cameras failed to check', 'warning');
    } finally {
      setRefreshing(false);
    }
  }, [cameras]);

  // Check cameras on current page only
  const checkCurrentPageCameras = useCallback(async () => {
    setRefreshing(true);
    try {
      const checkPromises = paginatedCameras.map(camera => checkCamera(camera.cameraId));
      await Promise.all(checkPromises);
      showSnackbar(`Page ${currentPage} cameras checked (${paginatedCameras.length} cameras)`, 'success');
    } catch (err) {
      showSnackbar('Some cameras on this page failed to check', 'warning');
    } finally {
      setRefreshing(false);
    }
  }, [paginatedCameras, currentPage]);

  // Refresh camera list from CSV
  const handleRefreshList = async () => {
    setRefreshing(true);
    try {
      await loadCameras();
      showSnackbar('Camera list refreshed successfully', 'success');
    } catch (err) {
      showSnackbar('Failed to refresh camera list', 'error');
    } finally {
      setRefreshing(false);
    }
  };

  // Open configuration dialog
  const handleOpenConfig = async (cameraId) => {
    setConfigDialog({
      open: true,
      cameraId,
      cameraInfo: null,
      loading: true
    });

    try {
      // Get current camera info/config
      const result = await checkSingleCamera(cameraId);
      setConfigDialog(prev => ({
        ...prev,
        cameraInfo: result,
        loading: false
      }));

      // Set form values from current config
      if (result.currentConfig) {
        const width = result.currentConfig.resolution?.width || 1920;
        const height = result.currentConfig.resolution?.height || 1080;
        setConfigForm({
          resolution: `${width}x${height}`,
          fps: result.currentConfig.fps || 30,
          bitrate: result.currentConfig.bitrate || 2000,
          encoding: result.currentConfig.encoding || 'h264'
        });
      }
    } catch (err) {
      setConfigDialog(prev => ({
        ...prev,
        loading: false
      }));
      showSnackbar(`Failed to load camera config: ${err.message}`, 'error');
    }
  };

  // Close configuration dialog
  const handleCloseConfig = () => {
    setConfigDialog({
      open: false,
      cameraId: '',
      cameraInfo: null,
      loading: false
    });
    setConfigForm({
      resolution: '1920x1080',
      fps: 30,
      bitrate: 2000,
      encoding: 'h264'
    });
  };

  // Handle config form changes
  const handleConfigFormChange = (field, value) => {
    setConfigForm(prev => ({
      ...prev,
      [field]: value
    }));
  };

  // Submit configuration
  const handleSubmitConfig = async () => {
    try {
      setConfigDialog(prev => ({ ...prev, loading: true }));

      // Parse resolution from "1920x1080" format
      const [width, height] = configForm.resolution.split('x').map(Number);

      const result = await configureSingleCamera(
        configDialog.cameraId,
        width,
        height,
        parseInt(configForm.fps),
        parseInt(configForm.bitrate),
        configForm.encoding
      );
      
      showSnackbar(`Camera ${configDialog.cameraId} configured successfully`, 'success');
      handleCloseConfig();

      // Refresh the camera status to get the actual current configuration
      await checkCamera(configDialog.cameraId);
      
    } catch (err) {
      showSnackbar(`Failed to configure camera: ${err.message}`, 'error');
    } finally {
      setConfigDialog(prev => ({ ...prev, loading: false }));
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

  const handlePageJumpChange = useCallback((event) => {
    const value = event.target.value;
    if (value === '' || /^\d+$/.test(value)) {
      setPageJumpValue(value);
    }
  }, []);

  // Create a reusable pagination component with quick page jump
  const PaginationComponent = useMemo(() => {
    if (totalPages <= 1) return null;
    
    return (
      <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', my: 2, gap: 2 }}>
        <Stack spacing={1} alignItems="center">
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
            <Pagination
              count={totalPages}
              page={currentPage}
              onChange={handlePageChange}
              color="primary"
              size="large"
              showFirstButton
              showLastButton
              siblingCount={1}
              boundaryCount={1}
            />
            {totalPages > 5 && (
              <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                <Typography variant="body2" color="text.secondary">
                  Go to:
                </Typography>
                <TextField
                  size="small"
                  value={pageJumpValue}
                  onChange={handlePageJumpChange}
                  onKeyPress={handlePageJump}
                  placeholder={`1-${totalPages}`}
                  sx={{ width: 80 }}
                  inputProps={{
                    style: { textAlign: 'center' }
                  }}
                />
              </Box>
            )}
          </Box>
          <Typography variant="body2" color="text.secondary" textAlign="center">
            Showing {startIndex + 1}-{Math.min(endIndex, cameras.length)} of {cameras.length} cameras
          </Typography>
        </Stack>
      </Box>
    );
  }, [totalPages, currentPage, handlePageChange, startIndex, endIndex, cameras.length, pageJumpValue, handlePageJumpChange, handlePageJump]);

  const getStatusChip = (cameraId) => {
    const status = cameraStatuses[cameraId];
    if (!status) {
      return <Chip size="small" label="Unknown" color="default" />;
    }

    switch (status.status) {
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
      case 'error':
        return (
          <Chip
            size="small"
            label="Error"
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

  const getResolution = (cameraId) => {
    const status = cameraStatuses[cameraId];
    if (!status || !status.currentConfig?.resolution) {
      return 'N/A';
    }
    const res = status.currentConfig.resolution;
    return `${res.width}x${res.height}`;
  };

  const getEncoding = (cameraId) => {
    const status = cameraStatuses[cameraId];
    return status?.currentConfig?.encoding || 'N/A';
  };

  const getFPS = (cameraId) => {
    const status = cameraStatuses[cameraId];
    return status?.currentConfig?.fps || 'N/A';
  };

  const getBitrate = (cameraId) => {
    const status = cameraStatuses[cameraId];
    return status?.currentConfig?.bitrate ? `${status.currentConfig.bitrate} kbps` : 'N/A';
  };

  const getErrorTooltip = (cameraId) => {
    const status = cameraStatuses[cameraId];
    return status?.error || '';
  };

  // Validate camera configuration using RTSP analyzer
  const validateCameraConfig = (cameraId) => {
    const validation = cameraValidations[cameraId];
    if (!validation) {
      return { isValid: null, message: 'Click to validate stream' };
    }

    if (validation.error && !validation.isValid) {
      return { isValid: false, message: validation.error };
    }

    const result = validation.validationResult;
    if (!result) {
      return { isValid: false, message: 'No validation data available' };
    }

    if (result.isValid) {
      return { isValid: true, message: 'Stream validation passed' };
    } else {
      return { isValid: false, message: result.error || 'Stream validation failed' };
    }
  };

  const getValidationChip = (cameraId) => {
    const validation = validateCameraConfig(cameraId);
    const isValidating = validating[cameraId];
    
    if (isValidating) {
      return (
        <Tooltip title="Validating stream..." arrow>
          <Chip
            size="small"
            label="Validating"
            color="default"
            icon={<CircularProgress size={16} />}
          />
        </Tooltip>
      );
    }

    if (validation.isValid === null) {
      return (
        <Tooltip title={validation.message} arrow>
          <Chip
            size="small"
            label="Validate"
            color="default"
            icon={<WarningIcon />}
            onClick={() => validateCamera(cameraId)}
            style={{ cursor: 'pointer' }}
          />
        </Tooltip>
      );
    }

    if (validation.isValid) {
      return (
        <Tooltip title={validation.message} arrow>
          <Chip
            size="small"
            label="Valid"
            color="success"
            icon={<CheckCircleIcon />}
            onClick={() => validateCamera(cameraId)}
            style={{ cursor: 'pointer' }}
          />
        </Tooltip>
      );
    } else {
      return (
        <Tooltip title={validation.message} arrow>
          <Chip
            size="small"
            label="Invalid"
            color="error"
            icon={<ErrorIcon />}
            onClick={() => validateCamera(cameraId)}
            style={{ cursor: 'pointer' }}
          />
        </Tooltip>
      );
    }
  };

  // Load cameras on component mount
  useEffect(() => {
    loadCameras();
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
        <Button variant="contained" onClick={loadCameras}>
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
            <Box display="flex" gap={2}>
              <Button
                variant="outlined"
                startIcon={refreshing ? <CircularProgress size={20} color="inherit" /> : <RefreshIcon />}
                onClick={handleRefreshList}
                disabled={refreshing}
              >
                {refreshing ? 'Loading...' : 'Refresh List'}
              </Button>
              {cameras.length > camerasPerPage && (
                <Button
                  variant="outlined"
                  onClick={checkCurrentPageCameras}
                  disabled={paginatedCameras.length === 0 || refreshing}
                >
                  Check Page {currentPage}
                </Button>
              )}
              <Button
                variant="contained"
                onClick={checkAllCameras}
                disabled={cameras.length === 0 || refreshing}
              >
                Check All ({cameras.length})
              </Button>
            </Box>
          </Box>
          <Typography variant="body2" color="text.secondary" sx={{ mt: 1 }}>
            Load cameras from CSV file and check their status individually
          </Typography>
          
          {/* Camera Count Info */}
          <Box sx={{ mt: 2 }}>
            <Typography variant="body2" color="text.secondary">
              {cameras.length} cameras loaded from CSV
            </Typography>
          </Box>
        </CardContent>
      </Card>

      {/* Pagination outside of the camera status monitor box */}
      {PaginationComponent}

      {cameras.length === 0 ? (
        <Alert severity="info">
          No cameras found. Add some cameras first to monitor their status.
        </Alert>
      ) : (
        <Box>
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
                  <TableCell>Validation</TableCell>
                  <TableCell align="center">Actions</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {paginatedCameras.map((camera) => {
                  const status = cameraStatuses[camera.cameraId];
                  const hasError = status?.status === 'offline' || status?.status === 'error';
                  const isChecking = checking[camera.cameraId];
                  
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
                          {getStatusChip(camera.cameraId)}
                          {hasError && status?.error && (
                            <Tooltip title={getErrorTooltip(camera.cameraId)} arrow>
                              <InfoIcon color="error" fontSize="small" />
                            </Tooltip>
                          )}
                          {status?.lastChecked && (
                            <Typography variant="caption" color="text.secondary" sx={{ ml: 1 }}>
                              Last checked: {status.lastChecked}
                            </Typography>
                          )}
                        </Box>
                      </TableCell>
                      <TableCell>
                        <Typography variant="body2">
                          {getResolution(camera.cameraId)}
                        </Typography>
                      </TableCell>
                      <TableCell>
                        <Typography variant="body2">
                          {getEncoding(camera.cameraId)}
                        </Typography>
                      </TableCell>
                      <TableCell>
                        <Typography variant="body2">
                          {getFPS(camera.cameraId)}
                        </Typography>
                      </TableCell>
                      <TableCell>
                        <Typography variant="body2">
                          {getBitrate(camera.cameraId)}
                        </Typography>
                      </TableCell>
                      <TableCell>
                        {getValidationChip(camera.cameraId)}
                      </TableCell>
                      <TableCell align="center">
                        <Box display="flex" gap={1}>
                          <Tooltip title="Check Camera Status">
                            <IconButton
                              color="info"
                              onClick={() => checkCamera(camera.cameraId)}
                              disabled={isChecking}
                              size="small"
                            >
                              {isChecking ? <CircularProgress size={16} /> : <RefreshIcon />}
                            </IconButton>
                          </Tooltip>
                          <Tooltip title="Configure Camera">
                            <IconButton
                              color="secondary"
                              onClick={() => handleOpenConfig(camera.cameraId)}
                              size="small"
                            >
                              <SettingsIcon />
                            </IconButton>
                          </Tooltip>
                          <Tooltip title="Launch VLC Player">
                            <IconButton
                              color="primary"
                              onClick={() => handleLaunchVLC(camera.cameraId)}
                              disabled={hasError}
                              size="small"
                            >
                              <PlayIcon />
                            </IconButton>
                          </Tooltip>
                        </Box>
                      </TableCell>
                    </TableRow>
                  );
                })}
              </TableBody>
            </Table>
          </TableContainer>
          
          {/* Bottom Pagination */}
          {PaginationComponent}
        </Box>
      )}

      {/* Configuration Dialog */}
      <Dialog
        open={configDialog.open}
        onClose={handleCloseConfig}
        maxWidth="sm"
        fullWidth
      >
        <DialogTitle>
          Configure Camera: {configDialog.cameraId}
        </DialogTitle>
        <DialogContent>
          {configDialog.loading ? (
            <Box display="flex" justifyContent="center" alignItems="center" p={3}>
              <CircularProgress />
              <Typography sx={{ ml: 2 }}>Loading camera configuration...</Typography>
            </Box>
          ) : (
            <Box sx={{ pt: 2 }}>
              <FormControl fullWidth sx={{ mb: 2 }}>
                <InputLabel>Resolution</InputLabel>
                <Select
                  value={configForm.resolution}
                  label="Resolution"
                  onChange={(e) => handleConfigFormChange('resolution', e.target.value)}
                >
                  <MenuItem value="320x240">320x240 (QVGA)</MenuItem>
                  <MenuItem value="640x480">640x480 (VGA)</MenuItem>
                  <MenuItem value="720x480">720x480 (NTSC)</MenuItem>
                  <MenuItem value="720x576">720x576 (PAL)</MenuItem>
                  <MenuItem value="800x600">800x600 (SVGA)</MenuItem>
                  <MenuItem value="1024x768">1024x768 (XGA)</MenuItem>
                  <MenuItem value="1280x720">1280x720 (720p)</MenuItem>
                  <MenuItem value="1280x960">1280x960 (SXGA)</MenuItem>
                  <MenuItem value="1280x1024">1280x1024 (SXGA)</MenuItem>
                  <MenuItem value="1600x1200">1600x1200 (UXGA)</MenuItem>
                  <MenuItem value="1920x1080">1920x1080 (1080p)</MenuItem>
                  <MenuItem value="2048x1536">2048x1536 (QXGA)</MenuItem>
                  <MenuItem value="2560x1440">2560x1440 (1440p)</MenuItem>
                  <MenuItem value="3840x2160">3840x2160 (4K)</MenuItem>
                </Select>
              </FormControl>
              
              <Box display="flex" gap={2} mb={2}>
                <TextField
                  label="Frame Rate (FPS)"
                  type="number"
                  value={configForm.fps}
                  onChange={(e) => handleConfigFormChange('fps', e.target.value)}
                  sx={{ flex: 1 }}
                  inputProps={{ min: 1, max: 60 }}
                />
                <TextField
                  label="Bitrate (kbps)"
                  type="number"
                  value={configForm.bitrate}
                  onChange={(e) => handleConfigFormChange('bitrate', e.target.value)}
                  sx={{ flex: 1 }}
                  inputProps={{ min: 100, max: 50000 }}
                />
              </Box>

              <FormControl fullWidth sx={{ mb: 2 }}>
                <InputLabel>Encoding</InputLabel>
                <Select
                  value={configForm.encoding}
                  label="Encoding"
                  onChange={(e) => handleConfigFormChange('encoding', e.target.value)}
                >
                  <MenuItem value="h264">H.264</MenuItem>
                  <MenuItem value="h265">H.265</MenuItem>
                  <MenuItem value="mjpeg">MJPEG</MenuItem>
                </Select>
              </FormControl>

              {configDialog.cameraInfo && (
                <Box sx={{ mt: 2, p: 2, bgcolor: 'grey.800', borderRadius: 1 }}>
                  <Typography variant="subtitle2" gutterBottom sx={{ color: 'white' }}>
                    Current Configuration:
                  </Typography>
                  <Typography variant="body2" sx={{ color: 'grey.300' }}>
                    Resolution: {getResolution(configDialog.cameraId)} | 
                    FPS: {getFPS(configDialog.cameraId)} | 
                    Bitrate: {getBitrate(configDialog.cameraId)} | 
                    Encoding: {getEncoding(configDialog.cameraId)}
                  </Typography>
                </Box>
              )}
            </Box>
          )}
        </DialogContent>
        <DialogActions>
          <Button onClick={handleCloseConfig}>
            Cancel
          </Button>
          <Button
            onClick={handleSubmitConfig}
            variant="contained"
            disabled={configDialog.loading}
          >
            {configDialog.loading ? <CircularProgress size={20} /> : 'Apply Configuration'}
          </Button>
        </DialogActions>
      </Dialog>

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
