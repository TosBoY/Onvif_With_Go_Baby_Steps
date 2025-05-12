import React, { useState, useEffect } from 'react';
import {
  Card,
  CardContent,
  Typography,
  Button,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Alert,
} from '@mui/material';
import axios from 'axios';

const CameraSelector = () => {
  const [cameras, setCameras] = useState([]);
  const [selectedCamera, setSelectedCamera] = useState('');
  const [message, setMessage] = useState(null);

  useEffect(() => {
    axios.get('/api/cameras')
      .then(response => setCameras(response.data))
      .catch(error => console.error('Error fetching cameras:', error));
  }, []);

  const handleCameraChange = (event) => {
    setSelectedCamera(event.target.value);
  };

  const handleSelect = () => {
    if (!selectedCamera) {
      setMessage({ type: 'error', text: 'Please select a camera first' });
      return;
    }

    axios.post('/api/select-camera', { cameraId: selectedCamera })
      .then(() => {
        setMessage({ type: 'success', text: 'Camera selected successfully' });
        // Reload the page after 2 seconds to refresh all components
        setTimeout(() => {
          window.location.reload();
        }, 2000);
      })
      .catch(error => {
        setMessage({ 
          type: 'error', 
          text: error.response?.data || 'Error selecting camera'
        });
      });
  };

  return (
    <Card variant="outlined" sx={{ mb: 3, border: '2px solid #555', boxShadow: 2 }}>
      <CardContent>
        <Typography variant="h6" color="primary" fontWeight="bold" gutterBottom>
          Select Camera
        </Typography>
        
        <FormControl fullWidth sx={{ mb: 2 }}>
          <InputLabel>Camera</InputLabel>
          <Select
            value={selectedCamera}
            label="Camera"
            onChange={handleCameraChange}
          >
            {cameras.map(camera => (
              <MenuItem 
                key={camera.id} 
                value={camera.id}
                disabled={camera.isFake}
              >
                {camera.ip} {camera.isFake ? '(Fake)' : ''}
              </MenuItem>
            ))}
          </Select>
        </FormControl>

        {message && (
          <Alert 
            severity={message.type} 
            sx={{ mb: 2 }}
            onClose={() => setMessage(null)}
          >
            {message.text}
          </Alert>
        )}

        <Button
          variant="contained"
          color="primary"
          onClick={handleSelect}
          fullWidth
        >
          Select Camera
        </Button>
      </CardContent>
    </Card>
  );
};

export default CameraSelector;
