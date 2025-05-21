import React, { useState, useEffect } from 'react';
import {
  Box,
  Typography,
  FormGroup,
  FormControlLabel,
  Checkbox,
  IconButton,
  Button,
  Paper,
} from '@mui/material';
import InfoIcon from '@mui/icons-material/Info';
import AddIcon from '@mui/icons-material/Add';
import axios from 'axios';
import CameraDetailsPopup from './CameraDetailsPopup';
import AddCameraDialog from './AddCameraDialog';

const CameraConfigDisplay = ({ selectedCameras, onCameraSelectionChange }) => {
  const [cameras, setCameras] = useState([]);
  const [openPopup, setOpenPopup] = useState(false);
  const [openAddDialog, setOpenAddDialog] = useState(false);
  const [selectedCamera, setSelectedCamera] = useState(null);
  const [refreshTrigger, setRefreshTrigger] = useState(0);

  useEffect(() => {
    axios.get('/api/cameras')
      .then(response => setCameras(response.data))
      .catch(error => console.error('Error fetching cameras:', error));
  }, [refreshTrigger]);

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

  const handlePopupClose = () => {
    setOpenPopup(false);
    setSelectedCamera(null);
  };

  const handleAddCamera = () => {
    setOpenAddDialog(true);
  };

  const handleAddDialogClose = (added) => {
    setOpenAddDialog(false);
    if (added) {
      setRefreshTrigger(prev => prev + 1);
    }
  };

  return (
    <Box>
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
        <Button
          variant="outlined"
          onClick={handleSelectDeselectAll}
          sx={{ mr: 1 }}
        >
          {selectedCameras.length === cameras.length ? 'Deselect All' : 'Select All'}
        </Button>
        <Button
          variant="contained"
          startIcon={<AddIcon />}
          onClick={handleAddCamera}
        >
          Add Camera
        </Button>
      </Box>

      <Paper variant="outlined" sx={{ p: 2 }}>
        <FormGroup>
          {cameras.map((camera) => (
            <Box
              key={camera.id}
              sx={{
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'space-between',
                py: 0.5
              }}
            >
              <FormControlLabel
                control={
                  <Checkbox
                    checked={selectedCameras.includes(camera.id)}
                    onChange={() => handleCheckboxChange(camera.id)}
                  />
                }
                label={
                  <Typography>
                    {camera.isFake ? '(Simulated) ' : ''}
                    Camera {camera.id} - {camera.ip}
                  </Typography>
                }
              />
              <IconButton
                size="small"
                onClick={() => handleInfoClick(camera)}
                disabled={camera.isFake}
              >
                <InfoIcon />
              </IconButton>
            </Box>
          ))}
        </FormGroup>
      </Paper>

      {selectedCamera && (
        <CameraDetailsPopup
          open={openPopup}
          onClose={handlePopupClose}
          camera={selectedCamera}
        />
      )}

      <AddCameraDialog
        open={openAddDialog}
        onClose={handleAddDialogClose}
      />
    </Box>
  );
};

export default CameraConfigDisplay;
