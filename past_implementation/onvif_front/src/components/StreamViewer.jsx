import { useState, useEffect } from 'react';
import { 
  Box,
  Card,
  CardContent,
  Typography,
  Button,
  TextField,
  Snackbar,
  Alert,
  CircularProgress,
  Tooltip
} from '@mui/material';
import api from '../services/api';

const StreamViewer = ({ profileToken }) => {
  const [streamUrl, setStreamUrl] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');

  useEffect(() => {
    const fetchStreamUrl = async () => {
      if (!profileToken) return;
      
      try {
        setLoading(true);
        const url = await api.getStreamUrl(profileToken);
        setStreamUrl(url);
        setLoading(false);
      } catch (err) {
        setError('Failed to fetch stream URL');
        setLoading(false);
      }
    };

    fetchStreamUrl();
  }, [profileToken]);

  const handleLaunchVLC = async () => {
    if (!profileToken) {
      setError('No profile selected');
      return;
    }

    try {
      setLoading(true);
      const response = await api.launchVLC(profileToken);
      setSuccess(response.message || 'VLC launched successfully');
      setLoading(false);
    } catch (err) {
      setError('Failed to launch VLC. Is VLC installed on the server?');
      setLoading(false);
    }
  };

  const handleCloseAlert = () => {
    setError('');
    setSuccess('');
  };

  return (
    <Card variant="outlined" sx={{ mb: 3 }}>
      <CardContent>
        <Typography variant="h6" gutterBottom>
          Camera Stream
        </Typography>
        
        {loading && (
          <Box sx={{ display: 'flex', justifyContent: 'center', my: 2 }}>
            <CircularProgress size={24} />
          </Box>
        )}
        
        {streamUrl && (
          <Box sx={{ mb: 2 }}>
            <TextField
              label="Stream URL"
              value={streamUrl}
              fullWidth
              InputProps={{ readOnly: true }}
              variant="outlined"
              margin="normal"
            />
            <Typography variant="body2" color="text.secondary" sx={{ mt: 1, mb: 2 }}>
              You can copy this URL to use in media players like VLC.
            </Typography>
          </Box>
        )}
        
        <Tooltip title="If VLC is already running, the stream will be added to the existing instance" arrow>
          <Button
            variant="contained"
            color="primary"
            onClick={handleLaunchVLC}
            disabled={!profileToken || loading}
            fullWidth
          >
            {loading ? "Working..." : "Play Stream in VLC"}
          </Button>
        </Tooltip>
        
        <Typography variant="caption" display="block" sx={{ mt: 1, textAlign: 'center' }}>
          Note: If VLC is already running, the stream will be added to the existing instance.
          Otherwise, a new VLC instance will be launched.
        </Typography>
      </CardContent>

      <Snackbar open={!!error} autoHideDuration={6000} onClose={handleCloseAlert}>
        <Alert onClose={handleCloseAlert} severity="error" sx={{ width: '100%' }}>
          {error}
        </Alert>
      </Snackbar>

      <Snackbar open={!!success} autoHideDuration={6000} onClose={handleCloseAlert}>
        <Alert onClose={handleCloseAlert} severity="success" sx={{ width: '100%' }}>
          {success}
        </Alert>
      </Snackbar>
    </Card>
  );
};

export default StreamViewer;