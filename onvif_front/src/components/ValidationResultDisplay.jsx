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

  const validatedCameras = results.filter(r => r.success && r.validationResult?.isValid);
  const invalidConfigs = results.filter(r => r.success && !r.validationResult?.isValid);
  const failedCameras = results.filter(r => !r.success);

  return (
    <Box sx={{ mt: 2 }}>
      <Typography variant="h6" gutterBottom>
        Configuration Validation Results
      </Typography>

      {validatedCameras.length > 0 && (
        <Alert severity="success" sx={{ mb: 2 }}>
          {validatedCameras.length} camera(s) validated successfully
        </Alert>
      )}

      {invalidConfigs.length > 0 && (
        <Box sx={{ mb: 2 }}>
          <Alert severity="warning" sx={{ mb: 1 }}>
            {invalidConfigs.length} camera(s) have configuration mismatches
          </Alert>
          <TableContainer component={Paper} variant="outlined">
            <Table size="small">
              <TableHead>
                <TableRow>
                  <TableCell>Camera ID</TableCell>
                  <TableCell>Parameter</TableCell>
                  <TableCell>Expected</TableCell>
                  <TableCell>Actual</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>                {invalidConfigs.map((result) => {
                  const v = result.validationResult;
                  if (!v) {
                    return (
                      <TableRow key={result.cameraId}>
                        <TableCell>{result.cameraId}</TableCell>
                        <TableCell colSpan={3}>Validation data not available</TableCell>
                      </TableRow>
                    );
                  }
                  return [
                    <TableRow key={`${result.cameraId}-res`}>
                      <TableCell rowSpan={2}>{result.cameraId}</TableCell>
                      <TableCell>Resolution</TableCell>
                      <TableCell>{v.expectedWidth}x{v.expectedHeight}</TableCell>
                      <TableCell>{v.actualWidth}x{v.actualHeight}</TableCell>
                    </TableRow>,
                    <TableRow key={`${result.cameraId}-fps`}>
                      <TableCell>Frame Rate</TableCell>
                      <TableCell>{v.expectedFPS} fps</TableCell>
                      <TableCell>{v.actualFPS ? v.actualFPS.toFixed(2) : 'N/A'} fps</TableCell>
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
            Failed to update {failedCameras.length} camera(s)
          </Alert>
          <TableContainer component={Paper} variant="outlined">
            <Table size="small">
              <TableHead>
                <TableRow>
                  <TableCell>Camera ID</TableCell>
                  <TableCell>Error</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {failedCameras.map((result) => (
                  <TableRow key={result.cameraId}>
                    <TableCell>{result.cameraId}</TableCell>
                    <TableCell>{result.error}</TableCell>
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
