import React from 'react';
import { Button, Box, Typography } from '@mui/material';
import { exportValidationCSV } from '../services/api';

const CSVExportTest = () => {
  const testExport = async () => {
    // Sample validation data for testing
    const sampleValidation = {
      "camera1": {
        "isValid": true,
        "expectedWidth": 1920,
        "expectedHeight": 1080,
        "expectedFPS": 30,
        "actualWidth": 1920,
        "actualHeight": 1080,
        "actualFPS": 29.97
      },
      "camera2": {
        "isValid": false,
        "expectedWidth": 1280,
        "expectedHeight": 720,
        "expectedFPS": 25,
        "actualWidth": 1920,
        "actualHeight": 1080,
        "actualFPS": 30.0
      },
      "camera3": {
        "isValid": true,
        "expectedWidth": 640,
        "expectedHeight": 480,
        "expectedFPS": 15,
        "actualWidth": 640,
        "actualHeight": 480,
        "actualFPS": 15.12
      }
    };

    try {
      await exportValidationCSV(sampleValidation);
      console.log('CSV export test successful!');
    } catch (error) {
      console.error('CSV export test failed:', error);
      alert('CSV export test failed: ' + error.message);
    }
  };

  return (
    <Box sx={{ p: 2, border: '1px dashed #ccc', borderRadius: 1, mb: 2 }}>
      <Typography variant="subtitle2" gutterBottom>
        CSV Export Test
      </Typography>
      <Button 
        variant="outlined" 
        size="small" 
        onClick={testExport}
        sx={{ fontSize: '0.8rem' }}
      >
        Test CSV Export
      </Button>
    </Box>
  );
};

export default CSVExportTest;
