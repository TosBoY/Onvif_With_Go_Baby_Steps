import { useState } from 'react';
import { 
  Dialog, 
  DialogTitle, 
  DialogContent, 
  DialogActions, 
  Typography, 
  Button, 
  Box, 
  Divider, 
  CircularProgress,
  Alert
} from '@mui/material';
import DeleteIcon from '@mui/icons-material/Delete';
import CloseIcon from '@mui/icons-material/Close';
import { deleteCamera } from '../services/api';

const CameraInfoDialog = ({ open, onClose, camera, onCameraDeleted }) => {
  const [isDeleting, setIsDeleting] = useState(false);
  const [error, setError] = useState('');
  
  const handleClose = () => {
    setError('');
    onClose();
  };
  
  const handleDeleteCamera = async () => {
    // Confirm deletion
    if (!window.confirm(`Are you sure you want to delete Camera ${camera.id}? This action cannot be undone.`)) {
      return;
    }
    
    setIsDeleting(true);
    setError('');
    
    try {
      await deleteCamera(camera.id);
      onCameraDeleted(camera.id);
      handleClose();
    } catch (err) {
      console.error('Failed to delete camera:', err);
      setError(`Failed to delete camera: ${err.message}`);
    } finally {
      setIsDeleting(false);
    }
  };
  
  // Only render if camera data is available
  if (!camera) return null;
  
  return (
    <Dialog 
      open={open} 
      onClose={handleClose} 
      maxWidth="sm" 
      fullWidth
      PaperProps={{
        sx: { borderRadius: 2 }
      }}
    >
      <DialogTitle sx={{ 
        display: 'flex', 
        justifyContent: 'space-between', 
        alignItems: 'center',
        borderBottom: '1px solid rgba(0,0,0,0.12)',
        pb: 1
      }}>
        <Typography variant="h6">Camera Information</Typography>
        <Button 
          onClick={handleClose} 
          color="inherit" 
          size="small" 
          startIcon={<CloseIcon />}
        >
          Close
        </Button>
      </DialogTitle>
      
      <DialogContent sx={{ pt: 3 }}>
        <Box sx={{ mb: 3 }}>
          <Typography variant="subtitle1" fontWeight="bold" gutterBottom>
            Camera Details
          </Typography>
          
          <Box sx={{ display: 'grid', gridTemplateColumns: '120px 1fr', rowGap: 1.5 }}>
            <Typography variant="body2" color="text.secondary">ID:</Typography>
            <Typography variant="body2">{camera.id}</Typography>
            
            <Typography variant="body2" color="text.secondary">IP Address:</Typography>
            <Typography variant="body2">{camera.ip}</Typography>
            
            <Typography variant="body2" color="text.secondary">Username:</Typography>
            <Typography variant="body2">{camera.username}</Typography>
            
            <Typography variant="body2" color="text.secondary">Type:</Typography>
            <Typography variant="body2">{camera.isFake ? 'Simulated Camera' : 'ONVIF Camera'}</Typography>
          </Box>
        </Box>
        
        <Divider sx={{ my: 2 }} />
        
        {error && (
          <Alert severity="error" sx={{ mb: 2 }}>
            {error}
          </Alert>
        )}
        
        <Box sx={{ mt: 2 }}>
          <Typography variant="subtitle2" color="error" gutterBottom>
            Danger Zone
          </Typography>
          <Typography variant="body2" color="text.secondary" paragraph>
            Deleting this camera will remove it from the system permanently. This action cannot be undone.
          </Typography>
        </Box>
      </DialogContent>
      
      <DialogActions sx={{ p: 2, pt: 0, justifyContent: 'flex-start' }}>
        <Button
          variant="contained"
          color="error"
          startIcon={isDeleting ? <CircularProgress size={16} color="inherit" /> : <DeleteIcon />}
          onClick={handleDeleteCamera}
          disabled={isDeleting}
        >
          {isDeleting ? 'Deleting...' : 'Delete Camera'}
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default CameraInfoDialog;
