import { useState, useEffect, useRef } from 'react';
import { 
  Container, 
  Typography, 
  Box, 
  Alert,
  Button,
  Paper,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  FormControl,
  InputLabel,
  OutlinedInput,
  InputAdornment,
  IconButton,
  Snackbar,
  Divider
} from '@mui/material';
import { 
  Refresh as RefreshIcon,
  Add as AddIcon,
  Visibility as VisibilityIcon,
  VisibilityOff as VisibilityOffIcon,
  Upload as UploadIcon
} from '@mui/icons-material';
import CameraCard from '../components/CameraCard';
import CameraConfigPanel from '../components/CameraConfigPanel';
import ValidationResults from '../components/ValidationResults'; // Add this import
import Loading from '../components/Loading';
import ConnectionStatus from '../components/ConnectionStatus';
import { getCameras, addNewCamera } from '../services/api';

const Dashboard = () => {
  const [cameras, setCameras] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [showError, setShowError] = useState(false);
  const [selectedCamera, setSelectedCamera] = useState(null);
  const [selectedCameras, setSelectedCameras] = useState([]);
    // Add validation state
  const [validationResults, setValidationResults] = useState(null);
  const [appliedConfig, setAppliedConfig] = useState(null);
  const [configSuccess, setConfigSuccess] = useState('');
  
  // Add camera dialog state
  const [addCameraDialogOpen, setAddCameraDialogOpen] = useState(false);
  const [newCameraIP, setNewCameraIP] = useState('');
  const [newCameraPort, setNewCameraPort] = useState('');
  const [newCameraURL, setNewCameraURL] = useState('');
  const [newCameraUsername, setNewCameraUsername] = useState('');
  const [newCameraPassword, setNewCameraPassword] = useState('');
  const [newCameraIsFake, setNewCameraIsFake] = useState(false);  const [showPassword, setShowPassword] = useState(false);
  const [addingCamera, setAddingCamera] = useState(false);
  const [addCameraError, setAddCameraError] = useState('');
  
  // CSV import state
  const [csvFile, setCsvFile] = useState(null);
  const [importingCsv, setImportingCsv] = useState(false);
  const [csvImportResult, setCsvImportResult] = useState(null);
  
  // Ref for error auto-hide timeout
  const errorTimeoutRef = useRef(null);
  
  console.log('Dashboard rendering with state:', { cameras, loading, error, selectedCamera, selectedCameras });
  const handleSuccessfulReconnection = () => {
    // Clear any existing error immediately when reconnection is successful
    if (error && showError) {
      setShowError(false);
      if (errorTimeoutRef.current) {
        clearTimeout(errorTimeoutRef.current);
        errorTimeoutRef.current = null;
      }
    }
    // Also call fetchCameras to refresh the camera list
    fetchCameras();
  };

  const fetchCameras = async () => {
    setLoading(true);
    setError(null);
    setShowError(false);
    
    // Clear any existing error timeout
    if (errorTimeoutRef.current) {
      clearTimeout(errorTimeoutRef.current);
      errorTimeoutRef.current = null;
    }
    
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
      const errorMessage = err.message || 'Failed to load cameras. Please check if the backend server is running.';
      setError(errorMessage);
      setShowError(true);
      setCameras([]); // Reset cameras on error
      
      // Auto-hide error after 8 seconds
      errorTimeoutRef.current = setTimeout(() => {
        setShowError(false);
      }, 8000);
    } finally {
      setLoading(false);
    }
  };

  const handleCameraDeleted = (cameraId) => {
    // Remove the camera from state
    setCameras(prevCameras => prevCameras.filter(camera => camera.id !== cameraId));
    
    // Remove from selected cameras if it was selected
    if (selectedCameras.includes(cameraId)) {
      setSelectedCameras(prevSelected => prevSelected.filter(id => id !== cameraId));
    }
    
    // Clear main selected camera if it was the deleted one
    if (selectedCamera?.id === cameraId) {
      setSelectedCamera(null);
    }
    
    // Show success message
    setConfigSuccess(`Camera ${cameraId} has been successfully deleted`);
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
      setSelectedCameras([...selectedCameras, camera.id]);      setSelectedCamera(camera); // Set as main selected camera for reference
    }
  };  // Add function to handle configuration results
  const handleConfigurationApplied = (result) => {
    console.log('Configuration result received:', result);
    
    // Handle both single result and array of results
    setValidationResults(result.validation);
    setAppliedConfig(result.appliedConfig);
    
    // Count successful validations
    const successCount = Array.isArray(result.validation) 
      ? result.validation.filter(v => v.isValid).length 
      : (result.validation?.isValid ? 1 : 0);
    
    const totalCount = Array.isArray(result.validation) ? result.validation.length : 1;
    
    setConfigSuccess(`Configuration applied successfully! ${successCount} of ${totalCount} cameras validated successfully.`);
    
    // No need to manually clear the message - the Snackbar will auto-hide
  };
  // Add function to clear validation results
  const handleClearValidation = () => {
    setValidationResults(null);
    setAppliedConfig(null);
    setConfigSuccess('');  };
  
  // Add camera dialog functions
  const handleAddCameraDialogOpen = () => {
    setAddCameraDialogOpen(true);
    setNewCameraIP('');
    setNewCameraPort('');
    setNewCameraURL('');
    setNewCameraUsername('');
    setNewCameraPassword('');
    setNewCameraIsFake(false);
    setAddCameraError('');
  };
    const handleAddCameraDialogClose = () => {
    setAddCameraDialogOpen(false);
    // Clear CSV import state
    setCsvFile(null);
    setCsvImportResult(null);
    setAddCameraError('');
  };
  
  const handleAddCamera = async () => {
    // Validate inputs
    if (!newCameraIP) {
      setAddCameraError('Camera IP address is required');
      return;
    }
    
    if (!newCameraUsername) {
      setAddCameraError('Username is required');
      return;
    }
    
    setAddingCamera(true);
    setAddCameraError('');
      try {
      // Convert port to number, default to 0 if empty
      const portValue = newCameraPort ? parseInt(newCameraPort, 10) : 0;
      const newCamera = await addNewCamera(newCameraIP, portValue, newCameraURL, newCameraUsername, newCameraPassword, newCameraIsFake);
      console.log('New camera added:', newCamera);
        // Close the dialog and refresh camera list
      setAddCameraDialogOpen(false);
      
      // Show success message
      setConfigSuccess(`New camera added successfully with ID: ${newCamera.id}`);
      
      // Clear success message after 6 seconds (matches Snackbar autoHideDuration)
      // This is handled automatically by the Snackbar's autoHideDuration
      
      // Refresh camera list
      fetchCameras();
    } catch (err) {
      console.error('Error adding camera:', err);
      setAddCameraError(err.message || 'Failed to add camera. Please check your inputs.');
    } finally {
      setAddingCamera(false);
    }
  };

  // CSV Import Functions
  const handleCsvFileChange = (event) => {
    const file = event.target.files[0];
    if (file && file.type === 'text/csv') {
      setCsvFile(file);
      setCsvImportResult(null);
      setAddCameraError('');
    } else if (file) {
      setAddCameraError('Please select a valid CSV file');
      setCsvFile(null);
    }
  };

  const handleCsvImport = async () => {
    if (!csvFile) {
      setAddCameraError('Please select a CSV file first');
      return;
    }

    setImportingCsv(true);
    setAddCameraError('');
    setCsvImportResult(null);

    try {
      const formData = new FormData();
      formData.append('csvFile', csvFile);

      const response = await fetch('/api/cameras/import-csv', {
        method: 'POST',
        body: formData,
      });

      const result = await response.json();

      if (response.ok || response.status === 206) { // Success or partial success
        setCsvImportResult(result);
        
        // Show success message
        setConfigSuccess(result.message);
        
        // Refresh camera list
        await fetchCameras();
        
        // If all cameras imported successfully, close dialog
        if (result.errorCount === 0) {
          setTimeout(() => {
            setAddCameraDialogOpen(false);
            setCsvFile(null);
            setCsvImportResult(null);
          }, 2000);
        }
      } else {
        setAddCameraError(result.message || 'CSV import failed');
      }
    } catch (error) {
      console.error('Error importing CSV:', error);
      setAddCameraError('Failed to import CSV. Please check your file and try again.');
    } finally {
      setImportingCsv(false);
    }
  };

  const handleClearCsvImport = () => {
    setCsvFile(null);
    setCsvImportResult(null);
    setAddCameraError('');
  };

  useEffect(() => {
    fetchCameras();
    
    // Cleanup timeout on unmount
    return () => {
      if (errorTimeoutRef.current) {
        clearTimeout(errorTimeoutRef.current);
      }
    };
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
          <ConnectionStatus onRefresh={handleSuccessfulReconnection} />        {error && showError && (
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
                <Box>
                  <Button
                    variant="outlined"
                    onClick={handleSelectDeselectAll}
                    disabled={cameras.length === 0}
                    sx={{ mr: 1 }}
                  >
                    {selectedCameras.length === cameras.length ? 'Deselect All' : 'Select All'}
                  </Button>
                  <Button 
                    variant="contained" 
                    color="primary"
                    startIcon={<AddIcon />}
                    onClick={handleAddCameraDialogOpen}
                  >
                    Add Camera
                  </Button>
                </Box>
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
                  {cameras.map((camera) => (                    <Box key={camera.id} sx={{ mb: 1 }}> {/* Increased margin-bottom back to 1 */}
                      <CameraCard 
                        camera={camera} 
                        isSelected={selectedCameras.includes(camera.id)}
                        onSelect={handleCameraSelect}
                        compact={true}
                        onCameraDeleted={handleCameraDeleted}
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
          </Box>        )}
      </Box>
      
      {/* Add Camera Dialog */}
      <Dialog 
        open={addCameraDialogOpen} 
        onClose={handleAddCameraDialogClose}
        maxWidth="sm"
        fullWidth
      >
        <DialogTitle>Add New Camera</DialogTitle>        <DialogContent>
          {addCameraError && (
            <Alert severity="error" sx={{ mb: 2, mt: 1 }}>
              {addCameraError}
            </Alert>
          )}
          
          {csvImportResult && (
            <Alert 
              severity={csvImportResult.errorCount === 0 ? "success" : csvImportResult.successCount > 0 ? "warning" : "error"} 
              sx={{ mb: 2, mt: 1 }}
            >
              <Typography variant="body2" sx={{ fontWeight: 'bold' }}>
                {csvImportResult.message}
              </Typography>
              {csvImportResult.errorCount > 0 && (
                <Typography variant="caption" sx={{ display: 'block', mt: 1 }}>
                  {csvImportResult.successCount} successful, {csvImportResult.errorCount} failed
                </Typography>
              )}
            </Alert>
          )}

          {/* CSV Import Section */}
          <Paper variant="outlined" sx={{ p: 2, mb: 3, bgcolor: 'action.hover' }}>
            <Typography variant="subtitle2" gutterBottom sx={{ fontWeight: 'bold' }}>
              Bulk Import from CSV
            </Typography>
            <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
              Upload a CSV file to add multiple cameras at once. Required columns: ip, username
            </Typography>
            
            <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 2 }}>
              <input
                type="file"
                accept=".csv"
                onChange={handleCsvFileChange}
                style={{ display: 'none' }}
                id="csv-file-input"
              />
              <label htmlFor="csv-file-input">
                <Button variant="outlined" component="span" startIcon={<UploadIcon />}>
                  Choose CSV File
                </Button>
              </label>
              {csvFile && (
                <Typography variant="body2" sx={{ flex: 1 }}>
                  {csvFile.name}
                </Typography>
              )}
              {csvFile && (
                <Button size="small" onClick={handleClearCsvImport}>
                  Clear
                </Button>
              )}
            </Box>
            
            <Button
              variant="contained"
              onClick={handleCsvImport}
              disabled={!csvFile || importingCsv}
              startIcon={importingCsv ? null : <UploadIcon />}
              fullWidth
            >
              {importingCsv ? 'Importing...' : 'Import CSV'}
            </Button>
          </Paper>

          <Divider sx={{ my: 2 }}>
            <Typography variant="caption" color="text.secondary">
              OR ADD SINGLE CAMERA
            </Typography>
          </Divider>
            <TextField
            autoFocus
            margin="dense"
            id="camera-ip"
            label="Camera IP Address"
            type="text"
            fullWidth
            variant="outlined"
            value={newCameraIP}
            onChange={(e) => setNewCameraIP(e.target.value)}
            placeholder="192.168.1.100"
            helperText="Enter the IP address of your ONVIF camera"
            sx={{ mb: 2, mt: 1 }}
          />
          
          <TextField
            margin="dense"
            id="camera-port"
            label="Port (Optional)"
            type="number"
            fullWidth
            variant="outlined"
            value={newCameraPort}
            onChange={(e) => setNewCameraPort(e.target.value)}
            placeholder="80"
            helperText="Enter the port number (default: 80 if left empty)"
            sx={{ mb: 2 }}
          />
          
          <TextField
            margin="dense"
            id="camera-url"
            label="Service URL (Optional)"
            type="text"
            fullWidth
            variant="outlined"
            value={newCameraURL}
            onChange={(e) => setNewCameraURL(e.target.value)}
            placeholder="onvif/media_service"
            helperText="Enter the ONVIF service path (default: onvif/media_service if left empty)"
            sx={{ mb: 2 }}
          />
          
          <TextField
            margin="dense"
            id="camera-username"
            label="Username"
            type="text"
            fullWidth
            variant="outlined"
            value={newCameraUsername}
            onChange={(e) => setNewCameraUsername(e.target.value)}
            placeholder="admin"
            helperText="Enter the username for camera authentication"
            sx={{ mb: 2 }}
          />
            <FormControl variant="outlined" fullWidth sx={{ mb: 1 }}>
            <InputLabel htmlFor="camera-password">Password</InputLabel>
            <OutlinedInput
              id="camera-password"
              type={showPassword ? 'text' : 'password'}
              value={newCameraPassword}
              onChange={(e) => setNewCameraPassword(e.target.value)}
              endAdornment={
                <InputAdornment position="end">
                  <IconButton
                    onClick={() => setShowPassword(!showPassword)}
                    edge="end"
                  >
                    {showPassword ? <VisibilityOffIcon /> : <VisibilityIcon />}
                  </IconButton>
                </InputAdornment>
              }
              label="Password"
            />
          </FormControl>
          <Typography variant="caption" color="text.secondary">
            Enter the password for camera authentication (leave empty if not required)
          </Typography>
          
          <Box sx={{ 
            display: 'flex', 
            alignItems: 'center', 
            mt: 2, 
            p: 1.5, 
            border: '1px solid',
            borderColor: newCameraIsFake ? 'primary.light' : 'divider',
            borderRadius: 1,
            bgcolor: newCameraIsFake ? 'action.hover' : 'transparent'
          }}>
            <FormControl fullWidth>
              <Box sx={{ display: 'flex', alignItems: 'center' }}>
                <input
                  type="checkbox"
                  id="camera-fake"
                  checked={newCameraIsFake}
                  onChange={(e) => setNewCameraIsFake(e.target.checked)}
                  style={{ 
                    marginRight: '12px', 
                    width: '18px', 
                    height: '18px', 
                    accentColor: newCameraIsFake ? '#1976d2' : undefined 
                  }}
                />
                <Typography 
                  onClick={() => setNewCameraIsFake(!newCameraIsFake)} 
                  sx={{ 
                    cursor: 'pointer', 
                    color: newCameraIsFake ? 'primary.main' : 'text.primary',
                    fontWeight: newCameraIsFake ? '500' : 'normal',
                    flex: 1
                  }}
                >
                  Use as fake camera (for testing)
                </Typography>
              </Box>
              {newCameraIsFake && (
                <Typography variant="caption" sx={{ mt: 1, display: 'block', ml: 3.5, color: 'info.main' }}>
                  Fake cameras allow you to test the application without connecting to a physical ONVIF camera.
                  The system will skip authentication and connection validation for this device.
                </Typography>
              )}
            </FormControl>          </Box>
        </DialogContent>
        <DialogActions sx={{ px: 3, pb: 2 }}>
          <Button 
            onClick={handleAddCameraDialogClose} 
            color="primary"
            variant="outlined"
            sx={{ mr: 1 }}
          >
            Cancel
          </Button>
          <Button 
            onClick={handleAddCamera} 
            variant="contained" 
            color="primary"
            disabled={addingCamera}
            startIcon={addingCamera ? null : <AddIcon size="small" />}
          >
            {addingCamera ? 'Adding Camera...' : 'Add Camera'}
          </Button>
        </DialogActions>
      </Dialog>
        {/* Success message Snackbar (appears at the bottom of the screen) */}
      <Snackbar
        open={!!configSuccess}
        autoHideDuration={6000}
        onClose={() => setConfigSuccess('')}
        anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}
        sx={{ mb: 2 }}
      >
        <Alert 
          onClose={() => setConfigSuccess('')}
          severity="success"
          variant="filled"
          sx={{ width: '100%', boxShadow: 3 }}
        >
          {configSuccess}
        </Alert>
      </Snackbar>
    </Container>
  );
};

export default Dashboard;
