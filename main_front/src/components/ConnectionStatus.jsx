import { useState, useEffect } from 'react';
import { Alert, Box, Button } from '@mui/material';
// Using only essential MUI icons
import { Refresh as RefreshIcon } from '@mui/icons-material';
import { testConnection } from '../services/api';

const ConnectionStatus = ({ onRefresh }) => {
  const [status, setStatus] = useState({ loading: true, success: false, message: 'Checking connection...' });

  const checkConnection = async () => {
    setStatus({ loading: true, success: false, message: 'Checking connection...' });
    const result = await testConnection();
    setStatus({ loading: false, success: result.success, message: result.message });
  };

  useEffect(() => {
    checkConnection();
  }, []);

  return (
    <Box sx={{ mb: 3 }}>
      <Alert 
        severity={status.loading ? 'info' : status.success ? 'success' : 'error'}
        action={
          <Button 
            color="inherit" 
            size="small" 
            onClick={() => {
              checkConnection();
              if (onRefresh) onRefresh();
            }}
            startIcon={<RefreshIcon />}
          >
            Retry
          </Button>
        }
      >
        Backend Server: {status.message}
      </Alert>
    </Box>
  );
};

export default ConnectionStatus;
