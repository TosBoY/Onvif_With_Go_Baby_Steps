import { AppBar, Toolbar, Typography, Box } from '@mui/material';
import { Videocam as VideocamIcon } from '@mui/icons-material';

const Header = () => {
  return (
    <AppBar position="static" sx={{ mb: 4 }}>
      <Toolbar>
        <VideocamIcon sx={{ mr: 2 }} />
        <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
          ONVIF Camera Management
        </Typography>
      </Toolbar>
    </AppBar>
  );
};

export default Header;
