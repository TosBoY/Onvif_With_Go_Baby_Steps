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
import ValidationResults from '../components/ValidationResults'; // Add this import
import Loading from '../components/Loading';
import ConnectionStatus from '../components/ConnectionStatus';
import { getCameras } from '../services/api';

const Dashboard = () => {
  const [cameras, setCameras] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [selectedCamera, setSelectedCamera] = useState(null);
  const [selectedCameras, setSelectedCameras] = useState([]);
  
  // Add validation state
  const [validationResults, setValidationResults] = useState(null);
  const [appliedConfig, setAppliedConfig] = useState(null);
  const [configSuccess, setConfigSuccess] = useState('');
  
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
  // Add function to handle configuration results
  const handleConfigurationApplied = (result) => {
    // Handle both single result and array of results
    setValidationResults(result.validation);
    setAppliedConfig(result.appliedConfig);
    
    // Count successful validations
    const successCount = Array.isArray(result.validation) 
      ? result.validation.filter(v => v.isValid).length 
      : (result.validation?.isValid ? 1 : 0);
    
    const totalCount = Array.isArray(result.validation) ? result.validation.length : 1;
    
    setConfigSuccess(`Configuration applied successfully! ${successCount} of ${totalCount} cameras validated successfully.`);
    
    // Clear success message after 5 seconds
    setTimeout(() => {
      setConfigSuccess('');
    }, 5000);
  };

  // Add function to clear validation results
  const handleClearValidation = () => {
    setValidationResults(null);
    setAppliedConfig(null);
    setConfigSuccess('');
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

        {configSuccess && (
          <Alert severity="success" sx={{ mb: 4 }} onClose={() => setConfigSuccess('')}>
            {configSuccess}
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
            >
              <Typography variant="h5" component="h2" gutterBottom>
                Camera Settings
              </Typography>
              <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
                Set resolution and frame rate, then apply to selected cameras.
              </Typography>
              <CameraConfigPanel 
                selectedCamera={selectedCamera} 
                selectedCameras={selectedCameras}
                cameras={cameras}
                onConfigurationApplied={handleConfigurationApplied} // Add this prop
                onClearValidation={handleClearValidation} // Add this prop
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
            >
              <Typography variant="h5" component="h2" gutterBottom>
                Camera List ({selectedCameras.length} selected)
              </Typography>
              <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
                Select cameras to apply configuration to. Multiple cameras can be selected.
              </Typography>
              
              <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
                <Button
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
                <Box sx={{ maxHeight: '400px', overflow: 'auto', py: 1 }}>
                  {cameras.map((camera) => (
                    <Box key={camera.id} sx={{ mb: 1 }}> {/* Increased margin-bottom back to 1 */}
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

        {/* Validation Results Panel - Full Width Below */}
        {validationResults && (
          <Box sx={{ 
            width: '100%',
            maxWidth: '1400px',
            mx: 'auto',
            mt: 3
          }}>
            <Paper elevation={3} sx={{ p: 3 }}>
              <ValidationResults 
                validation={validationResults} 
                appliedConfig={appliedConfig}
                onClear={handleClearValidation}
              />
            </Paper>
          </Box>
        )}
      </Box>
    </Container>
  );
};

export default Dashboard;
