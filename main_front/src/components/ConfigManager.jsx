import React from 'react';
import { Typography } from '@mui/material';

function ResolutionManager({ selectedCameras }) {
  return (
    <div>
      {selectedCameras.length === 0 ? (
        <Typography variant="body1">No camera selected.</Typography>
      ) : (
        <Typography variant="body1">
          Selected Camera IDs: {selectedCameras.join(', ')}
        </Typography>
      )}
    </div>
  );
}

export default ResolutionManager;
