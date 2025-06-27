import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { ThemeProvider, createTheme, CssBaseline, Box } from '@mui/material';
import Header from './components/Header';
import Dashboard from './pages/Dashboard';
import CameraStatus from './pages/CameraStatus';
import SimpleDashboard from './pages/SimpleDashboard';
import TestDashboard from './components/TestDashboard';
import DebugDashboard from './components/DebugDashboard';
import './App.css'

const darkTheme = createTheme({
  palette: {
    mode: 'dark',
    primary: {
      main: '#90caf9',
    },
    secondary: {
      main: '#f48fb1',
    },
    background: {
      default: '#121212',
      paper: '#1e1e1e',
    },
    text: {
      primary: '#ffffff',
      secondary: '#b0bec5',
    },
  },
});

function App() {
  return (
    <ThemeProvider theme={darkTheme}>
      <CssBaseline />
      <Router>
        <Box sx={{ display: 'flex', flexDirection: 'column', minHeight: '100vh' }}>
          <Header />
          <Box component="main" sx={{ flexGrow: 1, py: 3 }}>
            <Routes>
              <Route path="/" element={<Dashboard />} />
              <Route path="/status" element={<CameraStatus />} />
              <Route path="/debug" element={<DebugDashboard />} />
              <Route path="/test" element={<TestDashboard />} />
              <Route path="/simple" element={<SimpleDashboard />} />
            </Routes>
          </Box>
        </Box>
      </Router>
    </ThemeProvider>
  );
}

export default App;
