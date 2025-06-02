import React from 'react';
import { List, ListItem, ListItemText, Checkbox } from '@mui/material';

function CameraConfigDisplay({ cameraList, selectedCameras, onCameraSelectionChange }) {
  const toggleSelection = (cameraId) => {
    const updated = selectedCameras.includes(cameraId)
      ? selectedCameras.filter(id => id !== cameraId)
      : [...selectedCameras, cameraId];
    onCameraSelectionChange(updated);
  };

  return (
    <List>
      {cameraList.map(camera => (
        <ListItem key={camera.id} button onClick={() => toggleSelection(camera.id)}>
          <Checkbox checked={selectedCameras.includes(camera.id)} />
          <ListItemText primary={camera.name || `Camera ${camera.id}`} />
        </ListItem>
      ))}
    </List>
  );
}

export default CameraConfigDisplay;
