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
  // Direct approach without useState/useEffect to reduce complexity
  console.log("CameraConfigDisplay props:", { selectedProfile, selectedConfig });
  console.log("CameraConfigDisplay cameraInfo first items:", { 
    firstProfile: cameraInfo?.profiles?.[0], 
    firstConfig: cameraInfo?.configs?.[0] 
  });
  
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
  
  // Find the profile and config objects directly from the props
  // Note: Using Token (capital T) to match the data structure from ProfileSelector and ConfigSelector
  const profileDetails = cameraInfo.profiles.find(p => p.Token === selectedProfile);
  const configDetails = cameraInfo.configs.find(c => c.Token === selectedConfig);
  
  // Debug output to help diagnose missing fields
  console.log("Found objects:", { profileDetails, configDetails });
  
  // Helper function to safely navigate objects and get values with different possible property names
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
            (Selected profile: {selectedProfile}, Selected config: {selectedConfig})
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
            Profile: {profileDetails.Name}
          </Typography>
          
          <TableContainer component={Paper} variant="outlined" sx={{ mb: 2 }}>
            <Table size="small">
              <TableBody>
                <TableRow>
                  <TableCell component="th" scope="row" sx={{ fontWeight: 'bold' }}>Profile Token</TableCell>
                  <TableCell>{profileDetails.Token}</TableCell>
                </TableRow>
                <TableRow>
                  <TableCell component="th" scope="row" sx={{ fontWeight: 'bold' }}>Profile Name</TableCell>
                  <TableCell>{profileDetails.Name}</TableCell>
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
                  <TableCell>{configDetails.Token}</TableCell>
                </TableRow>
                <TableRow>
                  <TableCell component="th" scope="row" sx={{ fontWeight: 'bold' }}>Encoding</TableCell>
                  <TableCell>{configDetails.Encoding}</TableCell>
                </TableRow>
                
                {/* Resolution - try different property paths */}
                {(configDetails.Resolution || configDetails.Width) && (
                  <TableRow>
                    <TableCell component="th" scope="row" sx={{ fontWeight: 'bold' }}>Resolution</TableCell>
                    <TableCell>
                      {getProperty(configDetails.Resolution || configDetails, ['Width', 'width'])} × {getProperty(configDetails.Resolution || configDetails, ['Height', 'height'])}
                    </TableCell>
                  </TableRow>
                )}
                
                <TableRow>
                  <TableCell component="th" scope="row" sx={{ fontWeight: 'bold' }}>Quality</TableCell>
                  <TableCell>{configDetails.Quality}</TableCell>
                </TableRow>
                
                {/* Frame Rate - try different property paths */}
                {(configDetails.RateControl || configDetails.FrameRate || configDetails.FrameRateLimit) && (
                  <TableRow>
                    <TableCell component="th" scope="row" sx={{ fontWeight: 'bold' }}>Frame Rate</TableCell>
                    <TableCell>
                      {getProperty(configDetails.RateControl || configDetails, ['FrameRateLimit', 'frameRateLimit', 'FrameRate', 'frameRate'])} fps
                    </TableCell>
                  </TableRow>
                )}
                
                {/* Bit Rate - try different property paths */}
                {(configDetails.RateControl || configDetails.Bitrate || configDetails.BitrateLimit) && (
                  <TableRow>
                    <TableCell component="th" scope="row" sx={{ fontWeight: 'bold' }}>Bit Rate</TableCell>
                    <TableCell>
                      {getProperty(configDetails.RateControl || configDetails, ['BitrateLimit', 'bitrateLimit', 'Bitrate', 'bitrate'])} kbps
                    </TableCell>
                  </TableRow>
                )}
                
                {/* GOP Length - try different property paths */}
                {(configDetails.H264 || configDetails.GovLength) && (
                  <TableRow>
                    <TableCell component="th" scope="row" sx={{ fontWeight: 'bold' }}>GOP Length</TableCell>
                    <TableCell>
                      {getProperty(configDetails.H264 || configDetails, ['GovLength', 'govLength'])}
                    </TableCell>
                  </TableRow>
                )}
                
                {/* H264 Profile - try different property paths */}
                {(configDetails.H264 || configDetails.H264Profile) && (
                  <TableRow>
                    <TableCell component="th" scope="row" sx={{ fontWeight: 'bold' }}>H264 Profile</TableCell>
                    <TableCell>
                      {getProperty(configDetails.H264 || configDetails, ['H264Profile', 'h264Profile'])}
                    </TableCell>
                  </TableRow>
                )}

                {/* Add debug row to see all available properties */}
                <TableRow>
                  <TableCell component="th" scope="row" sx={{ fontWeight: 'bold', color: 'gray' }}>All Properties</TableCell>
                  <TableCell sx={{ fontSize: '0.7rem', color: 'gray', wordBreak: 'break-all' }}>
                    {Object.keys(configDetails).join(', ')}
                  </TableCell>
                </TableRow>
              </TableBody>
            </Table>
          </TableContainer>
        </Box>
      </CardContent>
    </Card>
  );
};

export default CameraConfigDisplay;
