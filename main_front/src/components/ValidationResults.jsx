import React from 'react';
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
  Stack
} from '@mui/material';
import CheckCircleIcon from '@mui/icons-material/CheckCircle';
import ErrorIcon from '@mui/icons-material/Error';

const ValidationResults = ({ validation, appliedConfig }) => {
  if (!validation) {
    return null;
  }
  // Check if this is a single validation result or multiple
  const isMultiple = Array.isArray(validation);
  
  // If single validation, convert to array for uniform processing
  const validations = isMultiple ? validation : [validation];
  
  // Function to check if resolution matches
  const resolutionMatches = (v) => {
    // If this is a fake camera, consider it matching if isValid is true
    if (v.isFake && v.isValid) return true;
    
    return v.actualWidth && v.actualHeight && 
           v.actualWidth === v.expectedWidth && 
           v.actualHeight === v.expectedHeight;
  };
  
  // Function to check if FPS matches
  const fpsMatches = (v) => {
    // If this is a fake camera, consider it matching if isValid is true
    if (v.isFake && v.isValid) return true;
    
    return v.actualFPS && Math.abs(v.actualFPS - v.expectedFPS) < 1;
  };

  // Function to check if bitrate matches (if expected bitrate was provided)
  const bitrateMatches = (v) => {
    // If this is a fake camera, consider it matching if isValid is true
    if (v.isFake && v.isValid) return true;
    
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
  
  // Separate successful and failed validations
  const successfulValidations = validations.filter(v => v.isValid);
  const failedValidations = validations.filter(v => !v.isValid);
  
  // For multi-camera view, separate by type of issue
  const resolutionFailures = failedValidations.filter(v => !resolutionMatches(v));
  const warningValidations = validations.filter(v => {
    // Include cameras that are valid but have FPS/bitrate mismatches
    if (v.isValid && (!fpsMatches(v) || !bitrateMatches(v))) {
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
    
    // Determine severity based on new business rules
    // Resolution mismatch = error, FPS/bitrate mismatch = warning
    const hasFpsOrBitrateWarnings = (fpsMismatch || bitrateMismatch) && validation.isValid;
    const alertSeverity = !validation.isValid ? 'error' : hasFpsOrBitrateWarnings ? 'warning' : 'success';

    return (
      <Box sx={{ mb: 3 }}>        <Alert 
          severity={alertSeverity}
          icon={validation.isValid ? 
            (hasFpsOrBitrateWarnings ? <ErrorIcon color="warning" sx={{ mr: 1 }} /> : getStatusIcon(true)) : 
            <ErrorIcon color="error" sx={{ mr: 1 }} />}
          sx={{ mb: 2 }}
        >
          <Box>
            <Typography variant="body1" sx={{ fontWeight: 'bold' }}>
              {!validation.isValid ? 'Validation Failed' : 
               hasFpsOrBitrateWarnings ? 'Configuration Applied with Warnings' : 
               'Validation Successful'}
            </Typography>
            {validation.error && (
              <Typography variant="body2" sx={{ mt: 0.5 }}>
                {!validation.isValid ? 
                 'Resolution mismatch detected - configuration failed.' :
                 'FPS or bitrate differences detected but within acceptable limits.'}
              </Typography>
            )}
          </Box>
        </Alert>        
        {(!validation.isValid || hasFpsOrBitrateWarnings) && (
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
            </CardContent>
          </Card>
        )}
          {validation.isValid && !hasFpsOrBitrateWarnings && (
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
  };

  return (
    <Box sx={{ mt: 3 }}>
      <Typography variant="h6" gutterBottom>
        Stream Validation Results
      </Typography>      {/* Show summary if multiple validations */}
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
          )}
        </Grid>
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