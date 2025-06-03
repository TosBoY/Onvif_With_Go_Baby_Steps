import React from 'react';
import {
  Box,
  Typography,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Alert,
} from '@mui/material';

const ValidationResultDisplay = ({ results }) => {
  if (!results || results.length === 0) {
    return null;
  }
  
  // Debug log to see the actual validation results
  console.log("Validation results:", JSON.stringify(results, null, 2));

  const validatedCameras = results.filter(r => r.success && r.validationResult?.isValid);
  const invalidConfigs = results.filter(r => r.validationResult && !r.validationResult.isValid);
  const failedCameras = results.filter(r => !r.validationResult);
  
  // Debug logs for each category
  console.log("Validated cameras:", validatedCameras);
  console.log("Invalid configs:", invalidConfigs);
  console.log("Failed cameras:", failedCameras);

  return (
    <Box sx={{ mt: 2, width: '100%', maxWidth: '100%', overflowX: 'auto' }}>
      <Typography variant="h6" gutterBottom>
        Configuration Validation Results
      </Typography>

      {/* Section for Passed Cameras */}
      {validatedCameras.length > 0 && (
        <Box sx={{ mb: 2 }}>
          <Alert severity="success" sx={{ mb: 1 }}>
            {validatedCameras.length} camera(s) validated successfully
          </Alert>
          <TableContainer component={Paper} variant="outlined" sx={{ width: '100%' }}>
            <Table size="small" sx={{ tableLayout: 'fixed' }}>
              <TableHead>
                <TableRow>
                  <TableCell width="30%">Camera ID</TableCell>
                  <TableCell width="70%">Status</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {validatedCameras.map((result) => (
                  <TableRow key={result.cameraId}>
                    <TableCell>{result.cameraId}</TableCell>
                    <TableCell sx={{ color: 'success.main' }}>
                      ✓ Validated Successfully
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </TableContainer>
        </Box>
      )}

      {/* Section for Failed Cameras */}
      {invalidConfigs.length > 0 && (
        <Box sx={{ mb: 2 }}>
          <Alert severity="warning" sx={{ mb: 1 }}>
            Configuration validation failed - settings don't match expected values
          </Alert>
          <TableContainer component={Paper} variant="outlined" sx={{ width: '100%' }}>
            <Table size="small" sx={{ tableLayout: 'fixed' }}>
              <TableHead>
                <TableRow>
                  <TableCell width="15%">Camera ID</TableCell>
                  <TableCell width="20%">Parameter</TableCell>
                  <TableCell width="25%">Expected</TableCell>
                  <TableCell width="25%">Actual</TableCell>
                  <TableCell width="15%">Status</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {invalidConfigs.map((result) => {
                  const v = result.validationResult;
                  if (!v) return null;

                  const resolutionMatches = v.actualWidth === v.expectedWidth && v.actualHeight === v.expectedHeight;
                  const fpsMatches = Math.abs(v.actualFPS - v.expectedFPS) < 0.1;

                  return [
                    <TableRow key={`${result.cameraId}-res`}>
                      <TableCell rowSpan={2}>{result.cameraId}</TableCell>
                      <TableCell>Resolution</TableCell>
                      <TableCell>{v.expectedWidth}x{v.expectedHeight}</TableCell>
                      <TableCell>{v.actualWidth}x{v.actualHeight}</TableCell>
                      <TableCell sx={{ color: resolutionMatches ? 'success.main' : 'error.main' }}>
                        {resolutionMatches ? '✓' : '✗'}
                      </TableCell>
                    </TableRow>,
                    <TableRow key={`${result.cameraId}-fps`}>
                      <TableCell>Frame Rate</TableCell>
                      <TableCell>{v.expectedFPS} fps</TableCell>
                      <TableCell>{v.actualFPS.toFixed(1)} fps</TableCell>
                      <TableCell sx={{ color: fpsMatches ? 'success.main' : 'error.main' }}>
                        {fpsMatches ? '✓' : '✗'}
                      </TableCell>
                    </TableRow>
                  ];
                }).flat()}
              </TableBody>
            </Table>
          </TableContainer>
        </Box>
      )}

      {failedCameras.length > 0 && (
        <Box sx={{ mb: 2 }}>
          <Alert severity="error" sx={{ mb: 1 }}>
            Validation check failed for {failedCameras.length} camera(s)
          </Alert>
          <TableContainer component={Paper} variant="outlined" sx={{ width: '100%' }}>
            <Table size="small" sx={{ tableLayout: 'fixed' }}>
              <TableHead>
                <TableRow>
                  <TableCell width="20%">Camera ID</TableCell>
                  <TableCell width="80%">Error</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {failedCameras.map((result) => (
                  <TableRow key={result.cameraId}>
                    <TableCell>{result.cameraId}</TableCell>
                    <TableCell sx={{ wordBreak: 'break-word' }}>{result.error}</TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </TableContainer>
        </Box>
      )}
    </Box>
  );
};

export default ValidationResultDisplay;
