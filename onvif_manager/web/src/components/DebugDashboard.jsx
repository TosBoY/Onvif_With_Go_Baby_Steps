import { useState, useEffect } from 'react';
import { 
  Container, 
  Typography, 
  Box, 
  Alert,
  Button,
  Paper
} from '@mui/material';
import { Refresh as RefreshIcon } from '@mui/icons-material';

const DebugDashboard = () => {
  const [step, setStep] = useState(1);

  console.log('DebugDashboard rendering, step:', step);

  const renderStep = () => {
    switch(step) {
      case 1:
        return (
          <Paper sx={{ p: 3 }}>
            <Typography variant="h6">Step 1: Basic rendering works</Typography>
            <Button onClick={() => setStep(2)} variant="contained" sx={{ mt: 2 }}>
              Test Loading Component
            </Button>
          </Paper>
        );
      case 2:
        return (
          <Paper sx={{ p: 3 }}>
            <Typography variant="h6">Step 2: Loading component test</Typography>
            <Box display="flex" justifyContent="center" sx={{ my: 2 }}>
              {/* Simple loading indicator without external component */}
              <Typography>Loading simulation...</Typography>
            </Box>
            <Button onClick={() => setStep(3)} variant="contained">
              Test API Import
            </Button>
          </Paper>
        );
      case 3:
        return (
          <Paper sx={{ p: 3 }}>
            <Typography variant="h6">Step 3: API service test</Typography>
            <Typography>getCameras function exists: {typeof getCameras}</Typography>
            <Button onClick={() => setStep(4)} variant="contained" sx={{ mt: 2 }}>
              Test Components
            </Button>
          </Paper>
        );
      default:
        return (
          <Paper sx={{ p: 3 }}>
            <Typography variant="h6">All basic tests passed!</Typography>
            <Button onClick={() => setStep(1)} variant="outlined">
              Reset
            </Button>
          </Paper>
        );
    }
  };

  // Step 3 import test
  let getCameras;
  try {
    getCameras = require('../services/api').getCameras;
  } catch (e) {
    getCameras = 'Import failed';
  }

  return (
    <Container maxWidth="lg">
      <Box sx={{ my: 4 }}>
        <Typography variant="h4" component="h1" gutterBottom>
          Debug Dashboard
        </Typography>
        {renderStep()}
      </Box>
    </Container>
  );
};

export default DebugDashboard;
