import React, { useState } from 'react';
import {
  Box,
  Typography,
  Paper,
  Alert,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Chip,
  Grid,
  Divider,
  Card,
  CardContent,
  List,
  ListItem,
  Stack,
  Button,
  CircularProgress
} from '@mui/material';
import CheckCircleIcon from '@mui/icons-material/CheckCircle';
import ErrorIcon from '@mui/icons-material/Error';
import DownloadIcon from '@mui/icons-material/Download';
import { exportValidationCSV } from '../services/api';

const ValidationResults = ({ validation, appliedConfig, configurationErrors }) => {
  const [isExporting, setIsExporting] = useState(false);

  if (!validation && !configurationErrors) {
    return null;
  }// CSV Export Function
  const handleExportCSV = async () => {
    setIsExporting(true);
    try {
      // Get camera IDs in the correct order - includes both validation results and configuration errors
      const allCameraIds = [];
      
      // Add IDs from validation results
      if (validations && validations.length > 0) {
        validations.forEach(v => {
          if (v.cameraId && !allCameraIds.includes(v.cameraId)) {
            allCameraIds.push(v.cameraId);
          }
        });
      }
      
      // Add IDs from configuration errors
      if (configurationErrors && configurationErrors.length > 0) {
        configurationErrors.forEach(e => {
          if (e.cameraId && !allCameraIds.includes(e.cameraId)) {
            allCameraIds.push(e.cameraId);
          }
        });
      }
      
      await exportValidationCSV(validation, configurationErrors, allCameraIds);
    } catch (error) {
      console.error('Error exporting CSV:', error);
      alert('Failed to export CSV. Please try again.');
    } finally {
      setIsExporting(false);
    }
  };// Check if this is a single validation result or multiple
  const isMultiple = Array.isArray(validation);
  
  // If single validation, convert to array for uniform processing
  const validations = validation ? (isMultiple ? validation : [validation]) : [];
  
  // Function to check if resolution matches
  const resolutionMatches = (v) => {
    return v.actualWidth && v.actualHeight && 
           v.actualWidth === v.expectedWidth && 
           v.actualHeight === v.expectedHeight;
  };
  
  // Function to check if FPS matches
  const fpsMatches = (v) => {
    return v.actualFPS && Math.abs(v.actualFPS - v.expectedFPS) < 1;
  };

  // Function to check if bitrate matches (if expected bitrate was provided)
  const bitrateMatches = (v) => {
    // If no expected bitrate was provided, consider it matching
    if (!v.expectedBitrate || v.expectedBitrate === 0) return true;
    
    // If actual bitrate is available, check if it's close (within 10% tolerance)
    if (v.actualBitrate) {
      const tolerance = v.expectedBitrate * 0.1; // 10% tolerance
      return Math.abs(v.actualBitrate - v.expectedBitrate) <= tolerance;
    }
    
    // If actual bitrate is not available, we can't determine matching
    return true; // Don't flag as mismatch if we can't measure it
  };
  
  // Function to check if encoding matches (if expected encoding was provided)
  const encodingMatches = (v) => {
    // If no expected encoding was provided, consider it matching
    if (!v.expectedEncoding) return true;
    
    // If actual encoding is available, do case-insensitive comparison
    if (v.actualEncoding) {
      return v.actualEncoding.toLowerCase().includes(v.expectedEncoding.toLowerCase());
    }
    
    // If actual encoding is not available, we can't determine matching
    return true; // Don't flag as mismatch if we can't measure it
  };
  
  // Separate successful and failed validations
  const successfulValidations = validations.filter(v => v.isValid);
  const failedValidations = validations.filter(v => !v.isValid);
  
  // For multi-camera view, separate by type of issue
  const resolutionFailures = failedValidations.filter(v => !resolutionMatches(v));
  const warningValidations = validations.filter(v => {
    // Include cameras that are valid but have FPS/bitrate/encoding mismatches
    if (v.isValid && (!fpsMatches(v) || !bitrateMatches(v) || !encodingMatches(v))) {
      return true;
    }
    return false;
  });

  const getStatusIcon = (isValid) => {
    return isValid ? (
      <CheckCircleIcon color="success" sx={{ mr: 1 }} />
    ) : (
      <ErrorIcon color="error" sx={{ mr: 1 }} />    );
  };
  
  // Single validation result renderer with improved layout
  const renderSingleValidation = (validation) => {
    const resolutionMismatch = !resolutionMatches(validation);
    const fpsMismatch = !fpsMatches(validation);
    const bitrateMismatch = !bitrateMatches(validation);
    const encodingMismatch = !encodingMatches(validation);
    
    // Determine severity based on new business rules
    // Resolution mismatch = error, FPS/bitrate/encoding mismatch = warning
    const hasWarnings = (fpsMismatch || bitrateMismatch || encodingMismatch) && validation.isValid;
    const alertSeverity = !validation.isValid ? 'error' : hasWarnings ? 'warning' : 'success';

    return (
      <Box sx={{ mb: 3 }}>        <Alert 
          severity={alertSeverity}
          icon={validation.isValid ? 
            (hasWarnings ? <ErrorIcon color="warning" sx={{ mr: 1 }} /> : getStatusIcon(true)) : 
            <ErrorIcon color="error" sx={{ mr: 1 }} />}
          sx={{ mb: 2 }}
        >
          <Box>
            <Typography variant="body1" sx={{ fontWeight: 'bold' }}>                {!validation.isValid ? 'Validation Failed' : 
               hasWarnings ? 'Configuration Applied with Warnings' : 
               'Validation Successful'}
            </Typography>
            {validation.error && (
              <Typography variant="body2" sx={{ mt: 0.5 }}>
                {!validation.isValid ? 
                 'Resolution mismatch detected - configuration failed.' :
                 'FPS, bitrate or encoding differences detected but within acceptable limits.'}
              </Typography>
            )}
          </Box>
        </Alert>        
        {(!validation.isValid || hasWarnings) && (
          <Card variant="outlined" sx={{ mb: 2 }}>
            <CardContent sx={{ pb: '16px !important' }}>              <Typography variant="subtitle1" gutterBottom sx={{ fontWeight: 'bold' }}>
                {!validation.isValid ? 'Configuration Failed - Details' : 'Configuration Warnings - Details'}
              </Typography>
              
              {resolutionMismatch && (
                <Box sx={{ mb: 2 }}>                  <Typography variant="body2" color="text.secondary" sx={{ fontWeight: 'medium' }}>
                    Resolution Adjusted:
                  </Typography>
                  <Grid container spacing={2} sx={{ mt: 0.5 }}>
                    {validation.requestedWidth && (
                      <Grid item xs={4}>
                        <Typography variant="body2" color="text.secondary">
                          Requested:
                        </Typography>
                        <Typography variant="body1">
                          {validation.requestedWidth}×{validation.requestedHeight}
                        </Typography>
                      </Grid>
                    )}
                    <Grid item xs={4}>
                      <Typography variant="body2" color="text.secondary">
                        Expected:
                      </Typography>
                      <Typography variant="body1">
                        {validation.expectedWidth}×{validation.expectedHeight}
                      </Typography>
                    </Grid>
                    <Grid item xs={4}>
                      <Typography variant="body2" color="text.secondary">
                        Actual:
                      </Typography>
                      <Typography variant="body1">
                        {validation.actualWidth || 0}×{validation.actualHeight || 0}
                      </Typography>
                    </Grid>
                  </Grid>
                </Box>
              )}
                {fpsMismatch && (
                <Box sx={{ mb: 2 }}>
                  <Typography variant="body2" color="text.secondary" sx={{ fontWeight: 'medium' }}>
                    Frame Rate Adjusted:
                  </Typography>
                  <Grid container spacing={2} sx={{ mt: 0.5 }}>
                    <Grid item xs={6}>
                      <Typography variant="body2" color="text.secondary">
                        Expected:
                      </Typography>
                      <Typography variant="body1">
                        {validation.expectedFPS} fps
                      </Typography>
                    </Grid>
                    <Grid item xs={6}>
                      <Typography variant="body2" color="text.secondary">
                        Actual:
                      </Typography>
                      <Typography variant="body1">
                        {validation.actualFPS ? validation.actualFPS.toFixed(1) : 0} fps
                      </Typography>
                    </Grid>
                  </Grid>
                </Box>
              )}

              {bitrateMismatch && validation.expectedBitrate && validation.expectedBitrate > 0 && (
                <Box>
                  <Typography variant="body2" color="text.secondary" sx={{ fontWeight: 'medium' }}>
                    Bitrate Adjusted:
                  </Typography>
                  <Grid container spacing={2} sx={{ mt: 0.5 }}>
                    <Grid item xs={6}>
                      <Typography variant="body2" color="text.secondary">
                        Expected:
                      </Typography>
                      <Typography variant="body1">
                        {validation.expectedBitrate} kbps
                      </Typography>
                    </Grid>
                    <Grid item xs={6}>
                      <Typography variant="body2" color="text.secondary">
                        Actual:
                      </Typography>
                      <Typography variant="body1">
                        {validation.actualBitrate ? `${validation.actualBitrate} kbps` : 'Auto/Unknown'}
                      </Typography>
                    </Grid>
                  </Grid>
                </Box>
              )}
              
              {encodingMismatch && validation.expectedEncoding && (
                <Box>
                  <Typography variant="body2" color="text.secondary" sx={{ fontWeight: 'medium' }}>
                    Encoding Difference:
                  </Typography>
                  <Grid container spacing={2} sx={{ mt: 0.5 }}>
                    <Grid item xs={6}>
                      <Typography variant="body2" color="text.secondary">
                        Expected:
                      </Typography>
                      <Typography variant="body1">
                        {validation.expectedEncoding || 'Auto'}
                      </Typography>
                    </Grid>
                    <Grid item xs={6}>
                      <Typography variant="body2" color="text.secondary">
                        Actual:
                      </Typography>
                      <Typography variant="body1">
                        {validation.actualEncoding || 'Unknown'}
                      </Typography>
                    </Grid>
                  </Grid>
                </Box>
              )}
            </CardContent>
          </Card>
        )}
          {validation.isValid && !hasWarnings && (
          <Alert severity="success" sx={{ mb: 2 }}>
            <Typography variant="body2">
              All parameters match exactly as expected. Camera is configured perfectly.
            </Typography>
          </Alert>
        )}
      </Box>
    );
  };
    // Improved compact list item for failed camera validation
  const renderFailedCamera = (v, index) => {
    const resolutionMismatch = !resolutionMatches(v);
    const fpsMismatch = !fpsMatches(v);
    const bitrateMismatch = !bitrateMatches(v);
    const hasRequestedDimensions = v.requestedWidth && v.requestedHeight;
    const actualWidth = v.actualWidth || 0;
    const actualHeight = v.actualHeight || 0;
    const actualFPS = v.actualFPS || 0;
    
    return (
      <ListItem 
        key={`failed-${index}`}
        divider={index < failedValidations.length - 1}
        sx={{ px: 2, py: 1.5 }}
      >
        <Box sx={{ width: '100%' }}>
          <Typography variant="subtitle1" sx={{ fontWeight: 'medium', mb: 1 }}>
            {v.cameraId || `Camera ${index + 1}`}
          </Typography>
          
          <Stack 
            direction="row" 
            spacing={2} 
            sx={{ 
              flexWrap: 'wrap', 
              '& > *': { mb: 1, mr: 1 } 
            }}
          >
            {resolutionMismatch && (
              <Box>
                <Stack direction="row" spacing={1} alignItems="center">
                  <Typography variant="body2" color="text.secondary" sx={{ fontSize: '0.75rem' }}>
                    Resolution:
                  </Typography>
                  {hasRequestedDimensions && (
                    <Chip 
                      size="small" 
                      color="default" 
                      label={`Requested: ${v.requestedWidth}×${v.requestedHeight}`}
                      variant="outlined"
                      sx={{ height: '22px', '& .MuiChip-label': { px: 1, py: 0 } }}
                    />
                  )}
                  <Chip 
                    size="small" 
                    color="warning" 
                    label={`Expected: ${v.expectedWidth}×${v.expectedHeight}`}
                    variant="outlined"
                    sx={{ height: '22px', '& .MuiChip-label': { px: 1, py: 0 } }}
                  />                  <Chip 
                    size="small" 
                    color="info" 
                    label={`Actual: ${actualWidth}×${actualHeight}`} 
                    variant="outlined"
                    sx={{ height: '22px', '& .MuiChip-label': { px: 1, py: 0 } }}
                  />
                </Stack>
                {(!actualWidth || !actualHeight) && (
                  <Typography variant="caption" color="error.main">                    Unable to detect resolution from stream
                  </Typography>
                )}
              </Box>
            )}
              {fpsMismatch && (
              <Box>
                <Stack direction="row" spacing={1} alignItems="center">
                  <Typography variant="body2" color="text.secondary" sx={{ fontSize: '0.75rem' }}>
                    FPS:
                  </Typography>
                  <Chip 
                    size="small" 
                    color="default" 
                    label={`Expected: ${v.expectedFPS}`}
                    variant="outlined"
                    sx={{ height: '22px', '& .MuiChip-label': { px: 1, py: 0 } }}
                  />
                  <Chip 
                    size="small" 
                    color="info" 
                    label={`Actual: ${actualFPS.toFixed(1)}`}
                    variant="outlined"
                    sx={{ height: '22px', '& .MuiChip-label': { px: 1, py: 0 } }}
                  />
                </Stack>
                {(!actualFPS) && (
                  <Typography variant="caption" color="error.main">
                    Unable to detect FPS from stream
                  </Typography>
                )}
              </Box>
            )}

            {bitrateMismatch && v.expectedBitrate && v.expectedBitrate > 0 && (
              <Box>
                <Stack direction="row" spacing={1} alignItems="center">
                  <Typography variant="body2" color="text.secondary" sx={{ fontSize: '0.75rem' }}>
                    Bitrate:
                  </Typography>
                  <Chip 
                    size="small" 
                    color="default" 
                    label={`Expected: ${v.expectedBitrate} kbps`}
                    variant="outlined"
                    sx={{ height: '22px', '& .MuiChip-label': { px: 1, py: 0 } }}
                  />
                  <Chip 
                    size="small" 
                    color="info" 
                    label={`Actual: ${v.actualBitrate ? `${v.actualBitrate} kbps` : 'Auto'}`}
                    variant="outlined"
                    sx={{ height: '22px', '& .MuiChip-label': { px: 1, py: 0 } }}
                  />
                </Stack>
                {(!v.actualBitrate && v.expectedBitrate > 0) && (
                  <Typography variant="caption" color="warning.main">
                    Bitrate validation not available from stream analysis
                  </Typography>
                )}
              </Box>
            )}
          </Stack>          {v.error && !v.error.includes("resolution mismatch: got") && !v.error.includes("FPS mismatch: got") && !v.error.includes("bitrate mismatch: got") && 
           !v.error.includes("RESOLUTION MISMATCH") && !v.error.includes("FPS DIFFERENCE") && !v.error.includes("BITRATE DIFFERENCE") && (
            <Typography variant="caption" color="text.secondary" sx={{ display: 'block', mt: 0.5 }}>
              Note: {v.error}
            </Typography>
          )}
        </Box>
      </ListItem>
    );
  };  return (
    <Box sx={{ mt: 3 }}>
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
        <Typography variant="h6">
          {validation && validation.length > 0 ? 'Stream Validation Results' : 'Configuration Results'}
        </Typography>        {((validation && validation.length > 0) || (configurationErrors && configurationErrors.length > 0)) && (
          <Button
            variant="outlined"
            startIcon={isExporting ? <CircularProgress size={16} /> : <DownloadIcon />}
            onClick={handleExportCSV}
            disabled={isExporting}
            size="small"
            sx={{ minWidth: '120px' }}
          >
            {isExporting ? 'Exporting...' : 'Export CSV'}
          </Button>
        )}
      </Box>      {/* Configuration Errors Section */}
      {configurationErrors && configurationErrors.length > 0 && (
        <Box sx={{ mb: 3 }}>
          <Typography variant="h6" gutterBottom sx={{ display: 'flex', alignItems: 'center', color: 'error.main' }}>
            <ErrorIcon sx={{ mr: 1 }} />
            Configuration Errors ({configurationErrors.length} camera{configurationErrors.length > 1 ? 's' : ''})
          </Typography>
          <Alert severity="error" sx={{ mb: 2 }}>
            <Typography variant="body1" sx={{ fontWeight: 'bold' }}>
              The following cameras failed during configuration and could not proceed to validation:
            </Typography>
          </Alert>
          <Paper variant="outlined">
            <List disablePadding>
              {configurationErrors.map((errorItem, index) => (
                <ListItem 
                  key={index}
                  sx={{ 
                    px: 2, 
                    py: 1.5,
                    borderBottom: index < configurationErrors.length - 1 ? '1px solid' : 'none',
                    borderColor: 'divider'
                  }}
                >
                  <Box sx={{ width: '100%' }}>
                    <Box sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
                      <ErrorIcon color="error" sx={{ mr: 1, fontSize: 20 }} />
                      <Typography variant="subtitle2" sx={{ fontWeight: 'bold' }}>
                        Camera ID: {errorItem.cameraId}
                      </Typography>
                    </Box>
                    <Typography variant="body2" color="text.secondary" sx={{ ml: 3 }}>
                      <strong>Error:</strong> {errorItem.error}
                    </Typography>
                  </Box>
                </ListItem>
              ))}
            </List>
          </Paper>
        </Box>
      )}

      {/* Show validation results only if there are validations */}
      {validation && validations.length > 0 && (
        <>
          {/* Show summary if multiple validations */}
          {isMultiple && (
            <Alert 
              severity={resolutionFailures.length > 0 ? 'error' : (warningValidations.length > 0 ? 'warning' : 'success')} 
              sx={{ mb: 2 }}
            >
              <Typography variant="body1">
                {resolutionFailures.length === 0 && warningValidations.length === 0 && 
                 `${successfulValidations.length} camera(s) configured successfully with exact settings.`}
                {resolutionFailures.length > 0 && 
                 `${resolutionFailures.length} camera(s) failed due to resolution incompatibility.`}
                {warningValidations.length > 0 && 
                 ` ${warningValidations.length} camera(s) configured with adjusted FPS/bitrate settings.`}
                {successfulValidations.length > 0 && (resolutionFailures.length > 0 || warningValidations.length > 0) && 
                 ` ${successfulValidations.length} camera(s) configured perfectly.`}
              </Typography>
            </Alert>
          )}

          {/* If single validation, render it directly */}
          {!isMultiple && renderSingleValidation(validation)}
        </>
      )}

      {/* If multiple validations, show in separate sections */}
      {isMultiple && (
        <Grid container spacing={3}>
          {/* Successful Validations Section */}
          {successfulValidations.length > 0 && (
            <Grid item xs={12}>
              <Box sx={{ mb: 3 }}>
                <Typography variant="subtitle1" gutterBottom sx={{ fontWeight: 'bold', color: 'success.main' }}>
                  Passed Cameras
                </Typography>
                <Alert severity="success">
                  <Typography variant="body2">
                    {successfulValidations.map(v => v.cameraId || "Unknown").join(", ")}
                  </Typography>
                </Alert>
              </Box>
            </Grid>
          )}          {/* Resolution Failure Section */}
          {resolutionFailures.length > 0 && (
            <Grid item xs={12}>
              <Typography variant="subtitle1" gutterBottom sx={{ fontWeight: 'bold', color: 'error.main' }}>
                Failed Cameras (Resolution Incompatible)
              </Typography>
              <Paper variant="outlined">
                <List disablePadding>
                  {resolutionFailures.slice(0, 10).map((v, index) => renderFailedCamera(v, index))}
                  {resolutionFailures.length > 10 && (
                    <ListItem sx={{ px: 2, py: 1.5, bgcolor: 'action.hover' }}>
                      <Typography variant="body2" color="text.secondary" sx={{ fontStyle: 'italic' }}>
                        ... and {resolutionFailures.length - 10} more camera(s) failed
                      </Typography>
                    </ListItem>
                  )}
                </List>
              </Paper>
            </Grid>
          )}          {/* Warning Validations Section - Cameras with FPS/Bitrate adjustments */}
          {warningValidations.length > 0 && (
            <Grid item xs={12}>
              <Typography variant="subtitle1" gutterBottom sx={{ fontWeight: 'bold', color: 'warning.main' }}>
                Cameras with Adjusted Settings (FPS/Bitrate)
              </Typography>
              <Paper variant="outlined">
                <List disablePadding>
                  {warningValidations.slice(0, 10).map((v, index) => renderFailedCamera(v, index))}
                  {warningValidations.length > 10 && (
                    <ListItem sx={{ px: 2, py: 1.5, bgcolor: 'action.hover' }}>
                      <Typography variant="body2" color="text.secondary" sx={{ fontStyle: 'italic' }}>
                        ... and {warningValidations.length - 10} more camera(s) with warnings
                      </Typography>
                    </ListItem>
                  )}
                </List>
              </Paper>
            </Grid>
          )}        </Grid>
      )}

      {appliedConfig && (
        <Box sx={{ mt: 2 }}>
          <Divider sx={{ my: 2 }} />
          <Typography variant="subtitle2" gutterBottom>
            Applied Configuration:
          </Typography>
          <Typography variant="body2" color="text.secondary">
            Resolution: {appliedConfig.resolution.width}×{appliedConfig.resolution.height}, 
            FPS: {appliedConfig.fps}
          </Typography>
        </Box>
      )}
    </Box>
  );
};

export default ValidationResults;