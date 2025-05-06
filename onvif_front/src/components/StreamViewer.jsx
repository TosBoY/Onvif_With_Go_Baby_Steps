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
  Tooltip,
  Link
} from '@mui/material';
import ContentCopyIcon from '@mui/icons-material/ContentCopy';
import api from '../services/api';

const StreamViewer = ({ profileToken }) => {
  const [streamUrl, setStreamUrl] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const [copied, setCopied] = useState(false);

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
    if (!streamUrl) {
      setError('Stream URL not available');
      return;
    }
    
    // Create a VLC protocol URL that will open VLC on the user's computer
    const vlcUrl = `vlc://${streamUrl}`;
    
    // Open the URL in a new window, which should trigger VLC on the client's computer
    window.open(vlcUrl, '_blank');
    
    setSuccess('VLC request sent to your browser. If VLC doesn\'t open, check if your browser supports the vlc:// protocol');
  };
  
  const copyToClipboard = () => {
    navigator.clipboard.writeText(streamUrl).then(
      () => {
        setCopied(true);
        setTimeout(() => setCopied(false), 2000);
      },
      () => {
        setError('Failed to copy to clipboard');
      }
    );
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
            <Box sx={{ display: 'flex', alignItems: 'center' }}>
              <TextField
                label="Stream URL"
                value={streamUrl}
                fullWidth
                InputProps={{ readOnly: true }}
                variant="outlined"
                margin="normal"
              />
              <Button 
                onClick={copyToClipboard}
                sx={{ ml: 1, mt: 1 }}
                variant="outlined"
                startIcon={<ContentCopyIcon />}
              >
                {copied ? 'Copied!' : 'Copy'}
              </Button>
            </Box>
            
            <Typography variant="body2" color="text.secondary" sx={{ mt: 1, mb: 2 }}>
              You can copy this URL to use in media players like VLC.
            </Typography>
            
            <Box sx={{ display: 'flex', gap: 2, mt: 2 }}>
              <Button
                variant="contained"
                color="primary"
                onClick={handleLaunchVLC}
                disabled={!streamUrl || loading}
                fullWidth
              >
                Open in VLC on Your Computer
              </Button>
              
              <Link 
                href={`vlc://${streamUrl}`}
                underline="none"
                sx={{ width: '100%' }}
              >
                <Button
                  variant="outlined"
                  color="secondary"
                  disabled={!streamUrl}
                  fullWidth
                >
                  Direct VLC Link
                </Button>
              </Link>
            </Box>
          </Box>
        )}
        
        <Typography variant="caption" display="block" sx={{ mt: 1, textAlign: 'center' }}>
          Note: The "Open in VLC" button will attempt to launch VLC on your computer.
          This feature requires your browser to support the vlc:// protocol.
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
      
      <Snackbar open={copied} autoHideDuration={2000} onClose={() => setCopied(false)}>
        <Alert onClose={() => setCopied(false)} severity="success" sx={{ width: '100%' }}>
          Stream URL copied to clipboard!
        </Alert>
      </Snackbar>
    </Card>
  );
};

export default StreamViewer;