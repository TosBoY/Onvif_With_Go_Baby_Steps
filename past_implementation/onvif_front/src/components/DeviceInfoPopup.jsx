import { useState, useEffect } from 'react';
import { 
  Box, 
  Typography, 
  CircularProgress,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableRow,
  Paper,
  Dialog,
  DialogTitle,
  DialogContent
} from '@mui/material';
import api from '../services/api';

const DeviceInfoPopup = ({ open, onClose }) => {
  const [deviceInfo, setDeviceInfo] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  useEffect(() => {
    if (open) {
      fetchDeviceInfo();
    }
  }, [open]);

  const fetchDeviceInfo = async () => {
    try {
      setLoading(true);
      setError('');
      const data = await api.getDeviceInfo();
      setDeviceInfo(data);
      setLoading(false);
    } catch (err) {
      setError('Failed to load device information');
      setLoading(false);
    }
  };

  return (
    <Dialog
      open={open}
      onClose={onClose}
      maxWidth="sm"
    >
      <DialogTitle>Device Information</DialogTitle>
      <DialogContent>
        <Box sx={{ p: 1, minWidth: 300 }}>
          {loading && (
            <Box sx={{ display: 'flex', justifyContent: 'center', my: 2 }}>
              <CircularProgress size={24} />
            </Box>
          )}
          
          {error && (
            <Typography color="error" sx={{ my: 2 }}>
              {error}
            </Typography>
          )}
          
          {deviceInfo && !loading && (
            <TableContainer component={Paper} sx={{ mt: 1, bgcolor: 'background.default' }}>
              <Table size="small">
                <TableBody>
                  <TableRow>
                    <TableCell component="th" sx={{ fontWeight: 'bold' }}>Manufacturer</TableCell>
                    <TableCell>{deviceInfo.manufacturer}</TableCell>
                  </TableRow>
                  <TableRow>
                    <TableCell component="th" sx={{ fontWeight: 'bold' }}>Model</TableCell>
                    <TableCell>{deviceInfo.model}</TableCell>
                  </TableRow>
                  <TableRow>
                    <TableCell component="th" sx={{ fontWeight: 'bold' }}>Firmware Version</TableCell>
                    <TableCell>{deviceInfo.firmwareVersion}</TableCell>
                  </TableRow>
                  <TableRow>
                    <TableCell component="th" sx={{ fontWeight: 'bold' }}>Serial Number</TableCell>
                    <TableCell>{deviceInfo.serialNumber}</TableCell>
                  </TableRow>
                  <TableRow>
                    <TableCell component="th" sx={{ fontWeight: 'bold' }}>Hardware ID</TableCell>
                    <TableCell>{deviceInfo.hardwareId}</TableCell>
                  </TableRow>
                </TableBody>
              </Table>
            </TableContainer>
          )}
        </Box>
      </DialogContent>
    </Dialog>
  );
};

export default DeviceInfoPopup;