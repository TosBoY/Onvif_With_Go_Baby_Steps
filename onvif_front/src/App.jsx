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
  Divider,
  Button,
  Tooltip
} from '@mui/material';
import { BrowserRouter as Router, Routes, Route, useNavigate } from 'react-router-dom';
import InfoIcon from '@mui/icons-material/Info';
import ProfileSelector from './components/ProfileSelector';
import ConfigSelector from './components/ConfigSelector';
import ResolutionManager from './components/ResolutionManager';
import StreamViewer from './components/StreamViewer';
import CameraConfigDisplay from './components/CameraConfigDisplay';
import DeviceInfoPage from './pages/DeviceInfoPage';
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
      // Using less contrasting text colors
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

// Main camera control component - extracted from main App component
const CameraControlPanel = () => {
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [cameraInfo, setCameraInfo] = useState(null);
  const [displayConfig, setDisplayConfig] = useState(null);
  const [selectedProfile, setSelectedProfile] = useState('');
  const [selectedConfig, setSelectedConfig] = useState('');
  const navigate = useNavigate();

  // Function to fetch camera information
  const fetchCameraInfo = async () => {
    try {
      setLoading(true);
      const data = await api.getCameraInfo();
      setCameraInfo(data);
      
      console.log("Camera Info Data Structure:", {
        profiles: data.profiles,
        configs: data.configs,
        firstProfileToken: data.profiles?.[0]?.token,
        firstConfigToken: data.configs?.[0]?.token,
        profilesCount: data.profiles?.length,
        configsCount: data.configs?.length
      });
      
      setLoading(false);
    } catch (err) {
      setError('Failed to load camera information. Please make sure the backend server is running.');
      setLoading(false);
    }
  };

  // Function to refresh only the configuration display without affecting form inputs
  const refreshConfigDisplay = async () => {
    if (!selectedConfig) return;
    
    try {
      console.log("Refreshing the config display for:", selectedConfig);
      const updatedConfig = await api.getSingleConfig(selectedConfig);
      console.log("Received updated config:", updatedConfig);
      
      // Only update the specific configuration without reloading everything
      if (cameraInfo && cameraInfo.configs) {
        // Create a new configs array with the updated config
        const updatedConfigs = cameraInfo.configs.map(config => {
          // Case-insensitive token comparison for robustness
          const configToken = config.Token || config.token;
          if (configToken && configToken.toLowerCase() === selectedConfig.toLowerCase()) {
            // Merge the existing config with the updated values
            console.log("Updating config:", configToken);
            return { 
              ...config,               // Keep existing properties
              ...updatedConfig,        // Add updated properties
              _updated: Date.now()     // Add timestamp to force detection of change
            };
          }
          return config;  // Return unchanged configs
        });
        
        // Create a new cameraInfo object to trigger re-render, only changing configs
        const newCameraInfo = {
          ...cameraInfo,
          configs: updatedConfigs
        };
        
        console.log("Updated config in camera info state");
        setCameraInfo(newCameraInfo);
      }
    } catch (err) {
      console.error("Error refreshing config display:", err);
    }
  };

  // Function to refresh camera info - can be called after updating configuration
  const refreshCameraInfo = () => {
    fetchCameraInfo();
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

  const handleDeviceInfoClick = () => {
    navigate('/device-info');
  };

  return (
    <Container maxWidth="xl" sx={{ py: 4 }}>
      <Typography 
        variant="h4" 
        component="h1" 
        sx={{ 
          mb: 3, 
          textAlign: 'center',
          color: 'rgba(144, 202, 249, 0.9)', // Slightly less bright blue
          fontWeight: 'bold',
          backgroundColor: 'rgba(144, 202, 249, 0.08)', // More subtle background
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
        <Grid container spacing={2}>
          {/* Main panel - Camera Configuration */}
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
                <>
                  <ResolutionManager 
                    configToken={selectedConfig}
                    profileToken={selectedProfile}
                    refreshCameraInfo={refreshConfigDisplay}
                  />
                  
                  {/* Debug values being passed */}
                  {console.log("App.jsx - Values being passed to CameraConfigDisplay:", { 
                    selectedProfile, 
                    selectedConfig, 
                    cameraInfo: cameraInfo ? {
                      profilesCount: cameraInfo.profiles?.length,
                      configsCount: cameraInfo.configs?.length
                    } : null
                  })}
                </>
              ) : (
                <Alert severity="info" sx={{ mt: 2 }}>
                  Please select both a profile and a configuration to adjust camera settings
                </Alert>
              )}
            </Paper>
          </Grid>
          
          {/* Side panel - Stream Viewer */}
          <Grid item xs={12} md={5}>
            <Paper elevation={3} sx={{ p: 3, height: '100%' }}>
              <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
                <Typography variant="h5" component="h2" sx={{ color: 'rgba(144, 202, 249, 0.9)' }}>
                  Camera Stream
                </Typography>
                
                <Tooltip title="View device information">
                  <Button 
                    variant="contained" 
                    color="primary" 
                    startIcon={<InfoIcon />} 
                    onClick={handleDeviceInfoClick}
                  >
                    Device Info
                  </Button>
                </Tooltip>
              </Box>
              
              {selectedProfile ? (
                <StreamViewer profileToken={selectedProfile} />
              ) : (
                <Alert severity="info">
                  Please select a profile to view the stream
                </Alert>
              )}

              {selectedProfile && selectedConfig && (
                <Box sx={{ mt: 4, p: 2, bgcolor: 'rgba(144, 202, 249, 0.05)', borderRadius: 2 }}>
                  <Typography variant="h6" sx={{ mb: 2, color: 'primary.main' }}>
                    Current Camera Configuration Details
                  </Typography>
                  <CameraConfigDisplay 
                    selectedConfig={selectedConfig}
                    selectedProfile={selectedProfile}
                    cameraInfo={cameraInfo}
                  />
                </Box>
              )}
            </Paper>
          </Grid>
        </Grid>
      )}
    </Container>
  );
};

// Main App component now just handles routing
function App() {
  return (
    <ThemeProvider theme={darkTheme}>
      <CssBaseline />
      <Router>
        <Routes>
          <Route path="/" element={<CameraControlPanel />} />
          <Route path="/device-info" element={<DeviceInfoPage />} />
        </Routes>
      </Router>
    </ThemeProvider>
  );
}

export default App;
