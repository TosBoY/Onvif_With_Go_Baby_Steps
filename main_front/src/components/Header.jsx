import React, { useState, useEffect } from 'react';
import { AppBar, Toolbar, Typography, Box, Tabs, Tab } from '@mui/material';
import { Videocam as VideocamIcon } from '@mui/icons-material';
import { useNavigate, useLocation } from 'react-router-dom';

const Header = () => {
  const navigate = useNavigate();
  const location = useLocation();
  
  // Map paths to tab indices
  const getTabValue = (pathname) => {
    switch (pathname) {
      case '/': return 0;
      case '/status': return 1;
      case '/simple': return 2;
      case '/debug': return 3;
      case '/test': return 4;
      default: return 0;
    }
  };

  const [tabValue, setTabValue] = useState(getTabValue(location.pathname));

  // Update tab value when location changes
  useEffect(() => {
    setTabValue(getTabValue(location.pathname));
  }, [location.pathname]);

  const handleTabChange = (event, newValue) => {
    setTabValue(newValue);
    const paths = ['/', '/status', '/simple', '/debug', '/test'];
    navigate(paths[newValue]);
  };

  return (
    <AppBar position="static" sx={{ mb: 4 }}>
      <Toolbar>
        <VideocamIcon sx={{ mr: 2 }} />
        <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
          ONVIF Camera Management
        </Typography>
        <Box sx={{ ml: 'auto' }}>
          <Tabs 
            value={tabValue} 
            onChange={handleTabChange}
            textColor="inherit"
            indicatorColor="secondary"
          >
            <Tab label="Dashboard" />
            <Tab label="Camera Status" />
            <Tab label="Simple" />
            <Tab label="Debug" />
            <Tab label="Test" />
          </Tabs>
        </Box>
      </Toolbar>
    </AppBar>
  );
};

export default Header;
