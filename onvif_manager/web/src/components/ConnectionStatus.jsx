import { useState, useEffect, useRef } from 'react';
import { Alert, Box, Typography } from '@mui/material';
import { testConnection } from '../services/api';

const ConnectionStatus = ({ onRefresh }) => {
  const [status, setStatus] = useState({ loading: true, success: false, message: 'Checking connection...' });
  const [retryCount, setRetryCount] = useState(0);
  const [showReconnectedMessage, setShowReconnectedMessage] = useState(false);
  const [showComponent, setShowComponent] = useState(true);
  const intervalRef = useRef(null);
  const timeoutRef = useRef(null);
  const reconnectedTimeoutRef = useRef(null);
  const hideComponentTimeoutRef = useRef(null);
  const checkConnection = async (isAutoCheck = false) => {
    if (!isAutoCheck) {
      setStatus({ loading: true, success: false, message: 'Checking connection...' });
      setShowReconnectedMessage(false);
      setShowComponent(true);
    }
    
    try {
      const result = await testConnection();
      console.log('Connection check result:', result, 'isAutoCheck:', isAutoCheck, 'retryCount:', retryCount);
      
      if (result.success) {
        // Connection is successful
        if (isAutoCheck && retryCount > 0) {
          // This was an auto-check and we were previously disconnected - show reconnected message
          console.log('Connection restored after being down');
          setStatus({ loading: false, success: true, message: 'Reconnected to backend' });
          setShowReconnectedMessage(true);
          setShowComponent(true);
          setRetryCount(0);
          
          // Auto-refresh cameras
          if (onRefresh) {
            console.log('Calling onRefresh due to reconnection');
            onRefresh();
          }
            // Clear reconnected message after 3 seconds and show normal connected status
          if (reconnectedTimeoutRef.current) {
            clearTimeout(reconnectedTimeoutRef.current);
          }
          reconnectedTimeoutRef.current = setTimeout(() => {
            setShowReconnectedMessage(false);
            setStatus({ loading: false, success: true, message: 'Connected to backend' });
            
            // Hide the entire component after another 2 seconds of showing "Connected to backend"
            if (hideComponentTimeoutRef.current) {
              clearTimeout(hideComponentTimeoutRef.current);
            }
            hideComponentTimeoutRef.current = setTimeout(() => {
              setShowComponent(false);
            }, 2000);
          }, 3000);        } else {
          // Either initial check or manual check - normal success handling
          setStatus({ loading: false, success: true, message: result.message });
          
          // If this is an auto-check and we have errors to clear, call onRefresh
          if (isAutoCheck && onRefresh) {
            console.log('Auto-check successful, calling onRefresh to clear any errors');
            onRefresh();
          }
          
          // Always set a timeout to hide the component for successful connections
          if (hideComponentTimeoutRef.current) {
            clearTimeout(hideComponentTimeoutRef.current);
          }
          hideComponentTimeoutRef.current = setTimeout(() => {
            setShowComponent(false);
          }, 3000);
        }
      } else {
        // Connection failed
        console.log('Connection failed:', result.message);
        setStatus({ loading: false, success: false, message: result.message });
        setShowComponent(true);
        setShowReconnectedMessage(false);
      }
      
      return result.success;
    } catch (error) {
      console.error('Connection check error:', error);
      setStatus({ loading: false, success: false, message: 'Connection failed' });
      setShowReconnectedMessage(false);
      setShowComponent(true);
      return false;
    }
  };
  const startAutoRetry = () => {
    // Clear any existing intervals/timeouts
    if (intervalRef.current) {
      console.log('Clearing existing interval');
      clearInterval(intervalRef.current);
      intervalRef.current = null;
    }
    if (timeoutRef.current) clearTimeout(timeoutRef.current);

    console.log('Starting auto-retry with 5 second interval');
    // Start checking every 5 seconds when connection is lost
    intervalRef.current = setInterval(async () => {
      console.log('Auto-retry check triggered, current retryCount:', retryCount);
      const isConnected = await checkConnection(true);
      
      if (isConnected) {
        // Connection restored, stop auto-checking
        console.log('Connection restored, stopping auto-retry');
        if (intervalRef.current) {
          clearInterval(intervalRef.current);
          intervalRef.current = null;
        }
        // Reset retry count is handled in checkConnection
      } else {
        console.log('Still disconnected, incrementing retry count');
        setRetryCount(prev => prev + 1);
      }
    }, 5000);
  };
  useEffect(() => {
    const initialCheck = async () => {
      console.log('Initial connection check');
      const isConnected = await checkConnection(false);
      if (!isConnected) {
        console.log('Initial check failed, starting auto-retry');
        setRetryCount(1);
        startAutoRetry();
      } else {
        console.log('Initial check successful');
      }
    };

    initialCheck();

    // Cleanup on unmount
    return () => {
      console.log('ConnectionStatus cleanup');
      if (intervalRef.current) clearInterval(intervalRef.current);
      if (timeoutRef.current) clearTimeout(timeoutRef.current);
      if (reconnectedTimeoutRef.current) clearTimeout(reconnectedTimeoutRef.current);
      if (hideComponentTimeoutRef.current) clearTimeout(hideComponentTimeoutRef.current);
    };
  }, []);

  // Start auto-retry when connection fails (but avoid duplicate intervals)
  useEffect(() => {
    if (!status.success && !status.loading && retryCount === 0 && !intervalRef.current) {
      console.log('Connection failed, starting auto-retry from useEffect');
      setRetryCount(1);
      setShowComponent(true);
      startAutoRetry();
    }
  }, [status.success, status.loading]);
  const getStatusMessage = () => {
    if (status.loading) {
      return status.message;
    }
    
    if (status.success) {
      if (showReconnectedMessage) {
        return 'Reconnected to backend';
      }
      return status.message;
    }
    
    if (retryCount > 0) {
      return `${status.message} (Auto-retry attempt ${retryCount})`;
    }
    
    return status.message;
  };
  return (
    showComponent ? (
      <Box sx={{ mb: 3 }}>
        <Alert 
          severity={status.loading ? 'info' : status.success ? 'success' : 'error'}
        >
          <Typography variant="body2">
            Backend Server: {getStatusMessage()}
          </Typography>
          {!status.success && !status.loading && !showReconnectedMessage && (
            <Typography variant="caption" sx={{ display: 'block', mt: 0.5, opacity: 0.8 }}>
              Automatically checking connection every 5 seconds...
            </Typography>
          )}
        </Alert>
      </Box>
    ) : null
  );
};

export default ConnectionStatus;
