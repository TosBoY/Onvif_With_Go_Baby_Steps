import React, { useState } from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  FormControlLabel,
  Checkbox,
  Button,
  Alert,
  Box,
  IconButton,
  Typography,
} from '@mui/material';
import CloseIcon from '@mui/icons-material/Close';
import api from '../services/api';

const AddCameraDialog = ({ open, onClose }) => {
  const [formData, setFormData] = useState({
    ip: '',
    username: '',
    password: '',
    isFake: false
  });
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const handleInputChange = (e) => {
    const { name, value, checked } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: name === 'isFake' ? checked : value
    }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      await api.addCamera(formData);
      // Reset form
      setFormData({
        ip: '',
        username: '',
        password: '',
        isFake: false
      });
      onClose(true); // Pass true to indicate successful addition
    } catch (error) {
      setError(error.response?.data || 'Failed to add camera');
    } finally {
      setLoading(false);
    }
  };

  const handleClose = () => {
    // Reset form and error state when closing
    setFormData({
      ip: '',
      username: '',
      password: '',
      isFake: false
    });
    setError('');
    onClose(false);
  };

  return (
    <Dialog open={open} onClose={handleClose} maxWidth="sm" fullWidth>
      <DialogTitle>
        <Box display="flex" justifyContent="space-between" alignItems="center">
          <Typography variant="h6">Add New Camera</Typography>
          <IconButton onClick={handleClose} size="small">
            <CloseIcon />
          </IconButton>
        </Box>
      </DialogTitle>
      <form onSubmit={handleSubmit}>
        <DialogContent>
          {error && <Alert severity="error" sx={{ mb: 2 }}>{error}</Alert>}
          
          <TextField
            name="ip"
            label="IP Address"
            value={formData.ip}
            onChange={handleInputChange}
            fullWidth
            required
            margin="normal"
            placeholder="192.168.1.100"
          />
          
          <TextField
            name="username"
            label="Username"
            value={formData.username}
            onChange={handleInputChange}
            fullWidth
            required={!formData.isFake}
            margin="normal"
            disabled={formData.isFake}
          />
          
          <TextField
            name="password"
            label="Password"
            type="password"
            value={formData.password}
            onChange={handleInputChange}
            fullWidth
            required={!formData.isFake}
            margin="normal"
            disabled={formData.isFake}
          />
          
          <FormControlLabel
            control={
              <Checkbox
                name="isFake"
                checked={formData.isFake}
                onChange={handleInputChange}
              />
            }
            label="Simulated Camera"
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={handleClose}>Cancel</Button>
          <Button
            type="submit"
            variant="contained"
            color="primary"
            disabled={loading || (!formData.isFake && (!formData.username || !formData.password))}
          >
            {loading ? 'Adding...' : 'Add Camera'}
          </Button>
        </DialogActions>
      </form>
    </Dialog>
  );
};

export default AddCameraDialog;
