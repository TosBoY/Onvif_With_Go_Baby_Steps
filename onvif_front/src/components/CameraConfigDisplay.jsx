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
  Tooltip,
} from '@mui/material';
import InfoIcon from '@mui/icons-material/Info';
import PlayArrowIcon from '@mui/icons-material/PlayArrow';
import AddIcon from '@mui/icons-material/Add';
import axios from 'axios';
import CameraDetailsPopup from './CameraDetailsPopup';
import AddCameraDialog from './AddCameraDialog';
import api from '../services/api';

const CameraConfigDisplay = ({ selectedCameras, onCameraSelectionChange }) => {
  const [cameras, setCameras] = useState([]);
  const [openPopup, setOpenPopup] = useState(false);
  const [openAddDialog, setOpenAddDialog] = useState(false);
  const [selectedCamera, setSelectedCamera] = useState(null);
  const [refreshTrigger, setRefreshTrigger] = useState(0);
  const [launchingVLC, setLaunchingVLC] = useState(null);

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

  const handleLaunchVLC = async (cameraId) => {
    setLaunchingVLC(cameraId);
    try {
      await api.launchVLC(cameraId);
    } catch (error) {
      console.error('Error launching VLC:', error);
    } finally {
      setLaunchingVLC(null);
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

      <Paper variant="outlined" sx={{ p: 2, minWidth: '400px' }}>
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
                sx={{ flexGrow: 1 }}
              />
              <Box sx={{ display: 'flex', gap: 1 }}>
                <Tooltip title="Open in VLC">
                  <IconButton
                    size="small"
                    onClick={() => handleLaunchVLC(camera.id)}
                    disabled={camera.isFake || launchingVLC === camera.id}
                  >
                    <PlayArrowIcon />
                  </IconButton>
                </Tooltip>
                <Tooltip title="Camera Info">
                  <IconButton
                    size="small"
                    onClick={() => handleInfoClick(camera)}
                  >
                    <InfoIcon />
                  </IconButton>
                </Tooltip>
              </Box>
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
