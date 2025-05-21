import { useState } from 'react';
import { 
  Container, 
  Typography, 
  Paper, 
  Grid,
  CssBaseline,
  ThemeProvider,
  createTheme,
} from '@mui/material';
import ResolutionManager from './components/ResolutionManager';
import CameraConfigDisplay from './components/CameraConfigDisplay';
import './App.css';

// Create a dark theme
const darkTheme = createTheme({
  palette: {
    mode: 'dark',
    primary: {
      main: '#90caf9',
    },
    secondary: {
      main: '#ce93d8',
    },
    background: {
      default: '#121212',
      paper: '#1e1e1e',
    },
  },
});

const App = () => {
  const [selectedCameras, setSelectedCameras] = useState([]);

  const handleCameraSelectionChange = (cameras) => {
    setSelectedCameras(cameras);
  };

  return (
    <ThemeProvider theme={darkTheme}>
      <CssBaseline />
      <Container maxWidth="xl" sx={{ py: 4 }}>
        <Typography variant="h4" gutterBottom>
          ONVIF Camera Control
        </Typography>
        <Grid container spacing={3}>
          <Grid item xs={12} md={7}>
            <Paper elevation={3} sx={{ p: 2 }}>
              <Typography variant="h5" component="h2" gutterBottom>
                Camera Settings
              </Typography>
              <ResolutionManager selectedCameras={selectedCameras} />
            </Paper>
          </Grid>
          <Grid item xs={12} md={5}>
            <Paper elevation={3} sx={{ p: 2 }}>
              <Typography variant="h5" component="h2" gutterBottom>
                Camera List
              </Typography>
              <CameraConfigDisplay 
                selectedCameras={selectedCameras} 
                onCameraSelectionChange={handleCameraSelectionChange}
              />
            </Paper>
          </Grid>
        </Grid>
      </Container>
    </ThemeProvider>
  );
};

export default App;