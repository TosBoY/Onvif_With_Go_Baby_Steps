import { useState, useEffect } from 'react';
import { 
  Container, 
  Typography, 
  Box, 
  Alert,
  Button,
  Paper
} from '@mui/material';
import { Refresh as RefreshIcon } from '@mui/icons-material';
import CameraCard from '../components/CameraCard';
import CameraConfigPanel from '../components/CameraConfigPanel';
import Loading from '../components/Loading';
import ConnectionStatus from '../components/ConnectionStatus';
import { getCameras } from '../services/api';

const Dashboard = () => {
  const [cameras, setCameras] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [selectedCamera, setSelectedCamera] = useState(null);
  const [selectedCameras, setSelectedCameras] = useState([]);
  
  console.log('Dashboard rendering with state:', { cameras, loading, error, selectedCamera, selectedCameras });

  const fetchCameras = async () => {
    setLoading(true);
    setError(null);
    
    try {
      console.log('Fetching cameras from API...');
      const data = await getCameras();
      console.log('Received camera data:', data);
      
      // Ensure we have an array of cameras
      if (Array.isArray(data)) {
        setCameras(data);
      } else if (data && Array.isArray(data.cameras)) {
        setCameras(data.cameras);
      } else {
        console.warn('Unexpected data format:', data);
        setCameras([]);
      }
    } catch (err) {
      console.error('Error fetching cameras:', err);
      setError(err.message || 'Failed to load cameras. Please check if the backend server is running.');
      setCameras([]); // Reset cameras on error
    } finally {
      setLoading(false);
    }
  };
  const handleSelectDeselectAll = () => {
    if (selectedCameras.length === cameras.length) {
      // Deselect all
      setSelectedCameras([]);
      setSelectedCamera(null);
    } else {
      // Select all
      setSelectedCameras(cameras.map(camera => camera.id));
      setSelectedCamera(cameras[0]); // Set first camera as the main selected one for config
    }
  };
  const handleCameraSelect = (camera) => {
    // Toggle camera selection in the list
    if (selectedCameras.includes(camera.id)) {
      setSelectedCameras(selectedCameras.filter(id => id !== camera.id));
      // If this was the main selected camera and we're deselecting it, clear it
      if (selectedCamera?.id === camera.id) {
        setSelectedCamera(null);
      }
    } else {
      setSelectedCameras([...selectedCameras, camera.id]);
      setSelectedCamera(camera); // Set as main selected camera for reference
    }
  };

  useEffect(() => {
    fetchCameras();
  }, []);

  if (loading) {
    return <Loading message="Loading cameras..." />;
  }

  return (
    <Container maxWidth="xl">
      <Box sx={{ my: 4 }}>
        <Typography variant="h4" component="h1" gutterBottom align="center">
          ONVIF Camera Control
        </Typography>
        
        <ConnectionStatus onRefresh={fetchCameras} />

        {error && (
          <Alert severity="error" sx={{ mb: 4 }}>
            {error}
          </Alert>
        )}
        
        <Box sx={{ 
          display: 'flex',
          width: '100%',
          maxWidth: '1400px',
          mx: 'auto',
          gap: 3,
          flexWrap: { xs: 'wrap', lg: 'nowrap' }
        }}>
          {/* Camera Settings Panel - Left Side */}
          <Box sx={{ 
            width: { xs: '100%', lg: '600px' },
            flexShrink: 0,
            flexGrow: 0,
          }}>
            <Paper 
              elevation={3} 
              sx={{ 
                p: 3,
                height: 'fit-content',
                minHeight: '500px'
              }}
            >              <Typography variant="h5" component="h2" gutterBottom>
                Camera Settings
              </Typography>              <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
                Set resolution and frame rate, then apply to selected cameras.
              </Typography>
              <CameraConfigPanel 
                selectedCamera={selectedCamera} 
                selectedCameras={selectedCameras}
                cameras={cameras}
              />
            </Paper>
          </Box>
          
          {/* Camera List Panel - Right Side */}
          <Box sx={{ 
            width: { xs: '100%', lg: '600px' },
            flexShrink: 0,
            flexGrow: 0,
          }}>
            <Paper 
              elevation={3} 
              sx={{ 
                p: 3,
                height: 'fit-content',
                minHeight: '500px'
              }}
            >              <Typography variant="h5" component="h2" gutterBottom>
                Camera List ({selectedCameras.length} selected)
              </Typography>              <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
                Select cameras to apply configuration to. Multiple cameras can be selected.
              </Typography>
              
              <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>                <Button
                  variant="outlined"
                  onClick={handleSelectDeselectAll}
                  disabled={cameras.length === 0}
                >
                  {selectedCameras.length === cameras.length ? 'Deselect All' : 'Select All'}
                </Button>
                <Button 
                  variant="outlined" 
                  startIcon={<RefreshIcon />}
                  onClick={fetchCameras}
                  size="small"
                >
                  Refresh
                </Button>
              </Box>

              {!error && cameras.length === 0 && (
                <Paper sx={{ p: 3, textAlign: 'center', bgcolor: '#f5f5f5' }}>
                  <Typography variant="h6" gutterBottom>
                    No cameras found
                  </Typography>
                  <Typography color="text.secondary" sx={{ mb: 2 }}>
                    Please check your configuration file and make sure cameras are properly set up.
                  </Typography>
                  <Button variant="outlined" onClick={fetchCameras} startIcon={<RefreshIcon />}>
                    Retry Loading
                  </Button>
                </Paper>
              )}
              
              {!error && cameras.length > 0 && (
                <Box sx={{ maxHeight: '400px', overflow: 'auto' }}>                  {cameras.map((camera) => (
                    <Box key={camera.id} sx={{ mb: 2 }}>
                      <CameraCard 
                        camera={camera} 
                        isSelected={selectedCameras.includes(camera.id)}
                        onSelect={handleCameraSelect}
                        compact={true}
                      />
                    </Box>
                  ))}
                </Box>
              )}
            </Paper>
          </Box>
        </Box>
      </Box>
    </Container>
  );
};

export default Dashboard;
