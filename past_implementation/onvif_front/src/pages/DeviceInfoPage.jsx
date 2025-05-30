import { useState, useEffect } from 'react';
import { 
  Container,
  Box, 
  Typography, 
  CircularProgress,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableRow,
  Paper,
  Button,
  ThemeProvider,
  CssBaseline,
  createTheme
} from '@mui/material';
import ArrowBackIcon from '@mui/icons-material/ArrowBack';
import { useNavigate } from 'react-router-dom';
import api from '../services/api';

// Use the same theme as the main app for consistency
const darkTheme = createTheme({
  palette: {
    mode: 'dark',
    primary: {
      main: '#90caf9', // A lighter blue for dark mode
    },
    secondary: {
      main: '#ce93d8', // Purple
    },
    background: {
      default: '#121212',
      paper: '#1e1e1e',
    },
    text: {
      primary: 'rgba(255, 255, 255, 0.85)', // Slightly transparent white
      secondary: 'rgba(176, 190, 197, 0.8)', // Slightly transparent gray-blue
    },
  },
  components: {
    MuiButton: {
      styleOverrides: {
        root: {
          borderRadius: 8,
        },
      },
    },
    MuiPaper: {
      styleOverrides: {
        root: {
          borderRadius: 8,
        },
      },
    },
  },
});

const DeviceInfoPage = () => {
  const [deviceInfo, setDeviceInfo] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const navigate = useNavigate();

  useEffect(() => {
    fetchDeviceInfo();
  }, []);

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

  const handleBack = () => {
    navigate('/');
  };

  return (
    <ThemeProvider theme={darkTheme}>
      <CssBaseline />
      <Container maxWidth="md" sx={{ py: 4 }}>
        <Button
          variant="contained"
          color="primary"
          startIcon={<ArrowBackIcon />}
          onClick={handleBack}
          sx={{ mb: 3 }}
        >
          Back to Main
        </Button>
        
        <Paper elevation={3} sx={{ p: 3, mb: 3 }}>
          <Typography variant="h4" component="h1" sx={{ mb: 4, fontWeight: 'bold', color: 'primary.main' }}>
            Device Information
          </Typography>
          
          <Box sx={{ p: 1, minWidth: 300 }}>
            {loading && (
              <Box sx={{ display: 'flex', justifyContent: 'center', my: 4 }}>
                <CircularProgress size={40} />
              </Box>
            )}
            
            {error && (
              <Typography color="error" sx={{ my: 2 }}>
                {error}
              </Typography>
            )}
            
            {deviceInfo && !loading && (
              <TableContainer component={Paper} sx={{ mt: 1, bgcolor: 'background.paper' }}>
                <Table>
                  <TableBody>
                    <TableRow>
                      <TableCell component="th" sx={{ fontWeight: 'bold', width: '40%' }}>Manufacturer</TableCell>
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
        </Paper>
      </Container>
    </ThemeProvider>
  );
};

export default DeviceInfoPage;