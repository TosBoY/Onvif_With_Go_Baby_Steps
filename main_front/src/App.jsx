import { useState, useEffect } from 'react';
import {
  Typography,
  Paper,
  CssBaseline,
  ThemeProvider,
  createTheme,
  Box
} from '@mui/material';
import axios from 'axios';
import ConfigManager from './components/ConfigManager';
import CameraConfigDisplay from './components/CameraConfigDisplay';
import './App.css';

const darkTheme = createTheme({
  palette: {
    mode: 'dark',
    primary: { main: '#90caf9' },
    secondary: { main: '#ce93d8' },
    background: {
      default: '#121212',
      paper: '#1e1e1e'
    },
  },
});

function App() {
  const [cameraList, setCameraList] = useState([]);
  const [selectedCameras, setSelectedCameras] = useState([]);

  useEffect(() => {
    axios.get('http://localhost:8080/api/cameras') // Adjust backend URL as needed
      .then(res => setCameraList(res.data))
      .catch(err => console.error('Error fetching camera list:', err));
  }, []);

  const handleCameraSelectionChange = (newSelection) => {
    setSelectedCameras(newSelection);
  };

  return (
    <ThemeProvider theme={darkTheme}>
      <CssBaseline />
      <Box sx={{ width: '100%', maxWidth: '100vw', minHeight: '100vh', p: 4, backgroundColor: 'background.default', overflow: 'hidden' }}>
        <Typography variant="h4" gutterBottom align="center">
          ONVIF Camera Control
        </Typography>
        <Box sx={{ display: 'flex', width: '100%', maxWidth: '1300px', mx: 'auto', gap: 3, flexWrap: { xs: 'wrap', md: 'nowrap' } }}>
          <Box sx={{ width: { xs: '100%', md: '550px' }, flexShrink: 0, flexGrow: 0 }}>
            <Paper elevation={3} sx={{ p: 2, height: '100%', overflow: 'hidden' }}>
              <Typography variant="h5" component="h2" gutterBottom>
                Camera Settings
              </Typography>
              <ConfigManager selectedCameras={selectedCameras} />
            </Paper>
          </Box>
          <Box sx={{ width: { xs: '100%', md: '550px' }, flexShrink: 0, flexGrow: 0 }}>
            <Paper elevation={3} sx={{ p: 2, height: '100%', overflow: 'hidden' }}>
              <Typography variant="h5" component="h2" gutterBottom>
                Camera List
              </Typography>
              <CameraConfigDisplay
                cameraList={cameraList}
                selectedCameras={selectedCameras}
                onCameraSelectionChange={handleCameraSelectionChange}
              />
            </Paper>
          </Box>
        </Box>
      </Box>
    </ThemeProvider>
  );
}

export default App;
