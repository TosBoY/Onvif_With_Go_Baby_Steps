import React, { useState, useEffect } from 'react';
import {
  Card,
  CardContent,
  Typography,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableRow,
  Paper,
  Box,
  Alert
} from '@mui/material';

const CameraConfigDisplay = ({ selectedProfile, selectedConfig, cameraInfo }) => {
  const [lastUpdate, setLastUpdate] = useState(Date.now());

  // Simple effect to track updates
  useEffect(() => {
    setLastUpdate(Date.now());
    console.log("CameraConfigDisplay updated");
  }, [cameraInfo, selectedProfile, selectedConfig]);
  
  if (!cameraInfo || !selectedProfile || !selectedConfig) {
    return (
      <Card variant="outlined" sx={{ mb: 3, mt: 3, border: '1px solid #444' }}>
        <CardContent>
          <Alert severity="info">
            Select both a profile and configuration to view details
          </Alert>
        </CardContent>
      </Card>
    );
  }
  
  // Find the profile and config objects with case-insensitive token comparison
  const findByToken = (array, tokenToFind) => {
    if (!array || !Array.isArray(array) || !tokenToFind) return null;
    return array.find(item => {
      const itemToken = item.Token || item.token;
      return itemToken && itemToken.toLowerCase() === tokenToFind.toLowerCase();
    });
  };
  
  const profileDetails = findByToken(cameraInfo.profiles, selectedProfile);
  const configDetails = findByToken(cameraInfo.configs, selectedConfig);
  
  // Helper function to safely navigate objects and get values
  const getProperty = (obj, propertyNames) => {
    if (!obj) return null;
    for (const name of propertyNames) {
      if (obj[name] !== undefined) return obj[name];
    }
    return null;
  };

  if (!profileDetails || !configDetails) {
    return (
      <Card variant="outlined" sx={{ mb: 3, mt: 3, border: '1px solid #444' }}>
        <CardContent>
          <Alert severity="warning">
            Could not find details for the selected profile or configuration.
          </Alert>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card variant="outlined" sx={{ mb: 3, border: '2px solid #555', boxShadow: 2 }}>
      <CardContent>
        <Box sx={{ mb: 3 }}>
          <Typography variant="subtitle1" color="primary" fontWeight="bold">
            Profile: {profileDetails.Name || profileDetails.name}
          </Typography>
          
          <TableContainer component={Paper} variant="outlined" sx={{ mb: 2 }}>
            <Table size="small">
              <TableBody>
                <TableRow>
                  <TableCell component="th" scope="row" sx={{ fontWeight: 'bold' }}>Profile Token</TableCell>
                  <TableCell>{profileDetails.Token || profileDetails.token}</TableCell>
                </TableRow>
                <TableRow>
                  <TableCell component="th" scope="row" sx={{ fontWeight: 'bold' }}>Profile Name</TableCell>
                  <TableCell>{profileDetails.Name || profileDetails.name}</TableCell>
                </TableRow>
                <TableRow>
                  <TableCell component="th" scope="row" sx={{ fontWeight: 'bold' }}>Fixed</TableCell>
                  <TableCell>{profileDetails.fixed ? 'Yes' : 'No'}</TableCell>
                </TableRow>
              </TableBody>
            </Table>
          </TableContainer>
        </Box>

        <Box sx={{ mt: 3 }}>
          <Typography variant="subtitle1" color="primary" fontWeight="bold">
            Video Configuration: {configDetails.Name || configDetails.name}
          </Typography>
          
          <TableContainer component={Paper} variant="outlined">
            <Table size="small">
              <TableBody>
                <TableRow>
                  <TableCell component="th" scope="row" sx={{ fontWeight: 'bold' }}>Config Token</TableCell>
                  <TableCell>{configDetails.Token || configDetails.token}</TableCell>
                </TableRow>
                <TableRow>
                  <TableCell component="th" scope="row" sx={{ fontWeight: 'bold' }}>Encoding</TableCell>
                  <TableCell>{configDetails.Encoding || configDetails.encoding}</TableCell>
                </TableRow>
                
                {/* Resolution - try different property paths */}
                {(configDetails.Resolution || configDetails.Width || configDetails.width) && (
                  <TableRow>
                    <TableCell component="th" scope="row" sx={{ fontWeight: 'bold' }}>Resolution</TableCell>
                    <TableCell>
                      {getProperty(configDetails.Resolution || configDetails, ['Width', 'width'])} × {getProperty(configDetails.Resolution || configDetails, ['Height', 'height'])}
                    </TableCell>
                  </TableRow>
                )}
                
                <TableRow>
                  <TableCell component="th" scope="row" sx={{ fontWeight: 'bold' }}>Quality</TableCell>
                  <TableCell>{configDetails.Quality || configDetails.quality}</TableCell>
                </TableRow>
                
                {/* Frame Rate - try different property paths */}
                {(configDetails.RateControl || configDetails.FrameRate || configDetails.FrameRateLimit || configDetails.frameRate || configDetails.frameRateLimit) && (
                  <TableRow>
                    <TableCell component="th" scope="row" sx={{ fontWeight: 'bold' }}>Frame Rate</TableCell>
                    <TableCell>
                      {getProperty(configDetails.RateControl || configDetails, ['FrameRateLimit', 'frameRateLimit', 'FrameRate', 'frameRate'])} fps
                    </TableCell>
                  </TableRow>
                )}
                
                {/* Bit Rate - try different property paths */}
                {(configDetails.RateControl || configDetails.Bitrate || configDetails.BitrateLimit || configDetails.bitrate || configDetails.bitrateLimit) && (
                  <TableRow>
                    <TableCell component="th" scope="row" sx={{ fontWeight: 'bold' }}>Bit Rate</TableCell>
                    <TableCell>
                      {getProperty(configDetails.RateControl || configDetails, ['BitrateLimit', 'bitrateLimit', 'Bitrate', 'bitrate'])} kbps
                    </TableCell>
                  </TableRow>
                )}
                
                {/* GOP Length - try different property paths */}
                {(configDetails.H264 || configDetails.GovLength || configDetails.govLength) && (
                  <TableRow>
                    <TableCell component="th" scope="row" sx={{ fontWeight: 'bold' }}>GOP Length</TableCell>
                    <TableCell>
                      {getProperty(configDetails.H264 || configDetails, ['GovLength', 'govLength'])}
                    </TableCell>
                  </TableRow>
                )}
                
                {/* H264 Profile - try different property paths */}
                {(configDetails.H264 || configDetails.H264Profile || configDetails.h264Profile) && (
                  <TableRow>
                    <TableCell component="th" scope="row" sx={{ fontWeight: 'bold' }}>H264 Profile</TableCell>
                    <TableCell>
                      {getProperty(configDetails.H264 || configDetails, ['H264Profile', 'h264Profile'])}
                    </TableCell>
                  </TableRow>
                )}
              </TableBody>
            </Table>
          </TableContainer>
        </Box>
      </CardContent>
    </Card>
  );
};

export default CameraConfigDisplay;
