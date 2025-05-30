import { useState } from 'react';
import { 
  Typography, 
  Paper, 
  CssBaseline,
  ThemeProvider,
  createTheme,
  Box,
} from '@mui/material';
import ResolutionManager from './components/ResolutionManager';
import CameraConfigDisplay from './components/CameraConfigDisplay';
import './App.css';

// Create a dark theme
const darkTheme = createTheme({
  palette: {
    mode: 'dark',
    primary: {
      main: '#90caf9',
    },
    secondary: {
      main: '#ce93d8',
    },
    background: {
      default: '#121212',
      paper: '#1e1e1e',
    },
  },
});

const App = () => {
  const [selectedCameras, setSelectedCameras] = useState([]);

  const handleCameraSelectionChange = (cameras) => {
    setSelectedCameras(cameras);
  };

  return (
    <ThemeProvider theme={darkTheme}>
      <CssBaseline />
      <Box sx={{ 
        width: '100%',
        maxWidth: '100vw', 
        minHeight: '100vh',
        p: 4,
        backgroundColor: 'background.default',
        overflow: 'hidden'
      }}>
        <Typography variant="h4" gutterBottom align="center">
          ONVIF Camera Control
        </Typography>
        
        <Box sx={{ 
          display: 'flex',
          width: '100%',
          maxWidth: '1300px',
          mx: 'auto',
          gap: 3,
          flexWrap: { xs: 'wrap', md: 'nowrap' }
        }}>
          {/* Camera Settings - Fixed width */}
          <Box sx={{ 
            width: { xs: '100%', md: '550px' },
            flexShrink: 0,
            flexGrow: 0,
          }}>
            <Paper 
              elevation={3} 
              sx={{ 
                p: 2,
                height: '100%',
                overflow: 'hidden'
              }}
            >
              <Typography variant="h5" component="h2" gutterBottom>
                Camera Settings
              </Typography>
              <Box sx={{ 
                width: '100%',
                overflow: 'hidden'
              }}>
                <ResolutionManager selectedCameras={selectedCameras} />
              </Box>
            </Paper>
          </Box>
          
          {/* Camera List - Fixed width */}
          <Box sx={{ 
            width: { xs: '100%', md: '550px' },
            flexShrink: 0,
            flexGrow: 0,
          }}>
            <Paper 
              elevation={3} 
              sx={{ 
                p: 2,
                height: '100%',
                overflow: 'hidden'
              }}
            >
              <Typography variant="h5" component="h2" gutterBottom>
                Camera List
              </Typography>
              <Box sx={{ 
                width: '100%',
                overflow: 'hidden'
              }}>
                <CameraConfigDisplay 
                  selectedCameras={selectedCameras} 
                  onCameraSelectionChange={handleCameraSelectionChange}
                />
              </Box>
            </Paper>
          </Box>
        </Box>
      </Box>
    </ThemeProvider>
  );
};

export default App;