import { useState, useEffect } from 'react';
import { 
  Container, 
  Typography, 
  Paper, 
  Box, 
  CircularProgress, 
  Alert,
  Grid,
  CssBaseline,
  ThemeProvider,
  createTheme,
} from '@mui/material';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import ProfileSelector from './components/ProfileSelector';
import ConfigSelector from './components/ConfigSelector';
import ResolutionManager from './components/ResolutionManager';
import CameraConfigDisplay from './components/CameraConfigDisplay';
import api from './services/api';
import './App.css';

// Create a dark theme with less contrasting text
const darkTheme = createTheme({
  palette: {
    mode: 'dark',
    primary: {
      main: '#90caf9', // A lighter blue for dark mode
    },
    secondary: {
      main: '#ce93d8', // Purple
    },
    background: {
      default: '#121212',
      paper: '#1e1e1e',
    },
    text: {
      primary: 'rgba(255, 255, 255, 0.85)', // Slightly transparent white
      secondary: 'rgba(176, 190, 197, 0.8)', // Slightly transparent gray-blue
    },
  },
  components: {
    MuiButton: {
      styleOverrides: {
        root: {
          borderRadius: 8,
        },
      },
    },
    MuiPaper: {
      styleOverrides: {
        root: {
          borderRadius: 8,
        },
      },
    },
  },
});

// Main camera control component
const CameraControlPanel = () => {
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [cameraInfo, setCameraInfo] = useState(null);
  const [selectedProfile, setSelectedProfile] = useState('');
  const [selectedConfig, setSelectedConfig] = useState('');
  const [selectedCameras, setSelectedCameras] = useState([]);

  const fetchCameraInfo = async () => {
    try {
      setLoading(true);
      const data = await api.getCameraInfo();
      setCameraInfo(data);
      setLoading(false);
    } catch (err) {
      setError('Failed to load camera information. Please make sure the backend server is running.');
      setLoading(false);
    }
  };

  const refreshConfigDisplay = async () => {
    if (!selectedConfig) return;
    try {
      const updatedConfig = await api.getSingleConfig(selectedConfig);
      if (cameraInfo && cameraInfo.configs) {
        const updatedConfigs = cameraInfo.configs.map(config => {
          const configToken = config.Token || config.token;
          if (configToken && configToken.toLowerCase() === selectedConfig.toLowerCase()) {
            return { 
              ...config,
              ...updatedConfig,
              _updated: Date.now()
            };
          }
          return config;
        });
        setCameraInfo({
          ...cameraInfo,
          configs: updatedConfigs
        });
      }
    } catch (err) {
      console.error("Error refreshing config display:", err);
    }
  };

  useEffect(() => {
    fetchCameraInfo();
  }, []);

  const handleProfileChange = (profileToken) => {
    setSelectedProfile(profileToken);
  };

  const handleConfigChange = (configToken) => {
    setSelectedConfig(configToken);
  };

  const handleCameraSelectionChange = (cameras) => {
    setSelectedCameras(cameras);
  };

  return (
    <Container maxWidth="xl" sx={{ py: 4 }}>
      <Typography 
        variant="h4" 
        component="h1" 
        sx={{ 
          mb: 3, 
          textAlign: 'center',
          color: 'rgba(144, 202, 249, 0.9)',
          fontWeight: 'bold',
          backgroundColor: 'rgba(144, 202, 249, 0.08)',
          padding: '12px',
          borderRadius: '8px',
          boxShadow: '0 2px 8px rgba(0,0,0,0.3)'
        }}
      >
        ONVIF Camera Control
      </Typography>
      
      {loading && (
        <Box sx={{ display: 'flex', justifyContent: 'center', my: 4 }}>
          <CircularProgress />
        </Box>
      )}

      {error && <Alert severity="error" sx={{ mb: 3 }}>{error}</Alert>}

      {!loading && cameraInfo && (
        <Grid container spacing={3}>
          <Grid item xs={12} md={7}>
            <Paper elevation={3} sx={{ p: 3, height: '100%' }}>
              <Typography variant="h5" component="h2" sx={{ mb: 2, color: 'rgba(144, 202, 249, 0.9)' }}>
                Camera Configuration
              </Typography>
              
              <ProfileSelector 
                profiles={cameraInfo.profiles} 
                onChange={handleProfileChange}
                selectedProfile={selectedProfile}
              />
              
              <ConfigSelector 
                configs={cameraInfo.configs} 
                onChange={handleConfigChange}
                selectedConfig={selectedConfig}
              />
              
              {selectedProfile && selectedConfig ? (
                <ResolutionManager 
                  configToken={selectedConfig}
                  profileToken={selectedProfile}
                  refreshCameraInfo={refreshConfigDisplay}
                  selectedCameras={selectedCameras}
                />
              ) : (
                <Alert severity="info" sx={{ mt: 2 }}>
                  Please select both a profile and a configuration to adjust camera settings
                </Alert>
              )}

            </Paper>
          </Grid>
          <Grid item xs={12} md={5}>
            <Paper elevation={3} sx={{ p: 3, height: '100%' }}>
              <Typography variant="h5" component="h2" sx={{ mb: 2, color: 'rgba(144, 202, 249, 0.9)' }}>
                Camera List
              </Typography>
              <CameraConfigDisplay 
                selectedProfile={selectedProfile}
                selectedConfig={selectedConfig}
                selectedCameras={selectedCameras}
                onCameraSelectionChange={handleCameraSelectionChange}
              />
            </Paper>
          </Grid>
        </Grid>
      )}
    </Container>
  );
};

function App() {
  return (
    <ThemeProvider theme={darkTheme}>
      <CssBaseline />
      <Router>
        <Routes>
          <Route path="/" element={<CameraControlPanel />} />
        </Routes>
      </Router>
    </ThemeProvider>
  );
}

export default App;
