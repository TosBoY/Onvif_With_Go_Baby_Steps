import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { ThemeProvider, createTheme, CssBaseline, Box } from '@mui/material';
import Header from './components/Header';
import Dashboard from './pages/Dashboard';
import SimpleDashboard from './pages/SimpleDashboard';
import TestDashboard from './components/TestDashboard';
import DebugDashboard from './components/DebugDashboard';
import './App.css'

// Create a theme instance
const theme = createTheme({
  palette: {
    primary: {
      main: '#2c3e50',
    },
    secondary: {
      main: '#3498db',
    },
    success: {
      main: '#2ecc71',
    },
    warning: {
      main: '#f39c12',
    },
    error: {
      main: '#e74c3c',
    }
  },
});

function App() {
  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <Router>
        <Box sx={{ display: 'flex', flexDirection: 'column', minHeight: '100vh' }}>
          <Header />
          <Box component="main" sx={{ flexGrow: 1, py: 3 }}>
            <Routes>
              <Route path="/" element={<Dashboard />} />
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
