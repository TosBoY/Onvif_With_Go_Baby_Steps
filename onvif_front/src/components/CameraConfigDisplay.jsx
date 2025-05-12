import React, { useState, useEffect } from 'react';
import {
  Card,
  CardContent,
  Typography,
  FormGroup,
  FormControlLabel,
  Checkbox,
  IconButton,
  Box,
  Button,
} from '@mui/material';
import InfoIcon from '@mui/icons-material/Info';
import axios from 'axios';
import CameraDetailsPopup from './CameraDetailsPopup';

const CameraConfigDisplay = ({ selectedProfile, selectedConfig, selectedCameras, onCameraSelectionChange }) => {
  const [cameras, setCameras] = useState([]);
  const [openPopup, setOpenPopup] = useState(false);
  const [selectedCamera, setSelectedCamera] = useState(null);

  useEffect(() => {
    axios.get('/api/cameras')
      .then(response => setCameras(response.data))
      .catch(error => console.error('Error fetching cameras:', error));
  }, []);

  const handleCheckboxChange = (cameraId) => {
    const newSelection = selectedCameras.includes(cameraId)
      ? selectedCameras.filter(id => id !== cameraId)
      : [...selectedCameras, cameraId];
    onCameraSelectionChange(newSelection);
  };

  const handleSelectDeselectAll = () => {
    if (selectedCameras.length === cameras.length) {
      onCameraSelectionChange([]);
    } else {
      onCameraSelectionChange(cameras.map(camera => camera.id));
    }
  };

  const handleInfoClick = (camera) => {
    setSelectedCamera(camera);
    setOpenPopup(true);
  };

  return (
    <div>
      <Card variant="outlined" sx={{ mb: 3, border: '2px solid #555', boxShadow: 2 }}>
        <CardContent>            
          <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
            <Typography variant="h6" color="primary" fontWeight="bold">
              Camera List
            </Typography>
            <Button
              variant="outlined"
              size="small"
              onClick={handleSelectDeselectAll}
            >
              {selectedCameras.length === cameras.length ? 'Deselect All' : 'Select All'}
            </Button>
          </Box>
          <FormGroup>
            {cameras.map(camera => (
              <Box 
                key={camera.id} 
                sx={{ 
                  display: 'flex', 
                  alignItems: 'center', 
                  justifyContent: 'space-between',
                  mb: 1 
                }}
              >
                <FormControlLabel
                  sx={{ flex: 1 }}
                  control={
                    <Checkbox
                      checked={selectedCameras.includes(camera.id)}
                      onChange={() => handleCheckboxChange(camera.id)}
                    />
                  }
                  label={`${camera.ip} (Fake: ${camera.isFake ? 'Yes' : 'No'})`}
                />
                <IconButton
                  size="small"
                  onClick={() => handleInfoClick(camera)}
                >
                  <InfoIcon />
                </IconButton>
              </Box>
            ))}
          </FormGroup>
        </CardContent>
      </Card>

      {selectedCamera && (
        <CameraDetailsPopup
          open={openPopup}
          onClose={() => setOpenPopup(false)}
          camera={selectedCamera}
          selectedProfile={selectedProfile}
          selectedConfig={selectedConfig}
        />
      )}
    </div>
  );
};

export default CameraConfigDisplay;
