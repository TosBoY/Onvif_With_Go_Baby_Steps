import { 
  Card, 
  CardContent, 
  Typography, 
  Button, 
  CardMedia,
  CardActions, 
  Chip,
  Box,
  Checkbox,
  IconButton,
  Tooltip
} from '@mui/material';
import { 
  Videocam as VideocamIcon, 
  Settings as SettingsIcon, 
  PlayArrow as PlayArrowIcon 
} from '@mui/icons-material';
import { launchVLC } from '../services/api';

const CameraCard = ({ camera, isSelected, onSelect, compact = false }) => {
  const getStatusColor = () => {
    if (camera.isFake) return 'warning';
    return 'success';
  };
  
  const handleLaunchVLC = async (e) => {
    e.stopPropagation(); // Prevent triggering camera selection
    try {
      const response = await launchVLC(camera.id);
      console.log('VLC launched successfully:', response);
      // Could show a success notification here
    } catch (error) {
      console.error('Failed to launch VLC:', error);
      // Could show an error notification here
    }
  };if (compact) {
    return (      <Card sx={{ 
        display: 'flex',
        alignItems: 'center',
        transition: 'all 0.2s ease-in-out',
        border: 'none'
      }}>
        <Box sx={{ display: 'flex', alignItems: 'center', pl: 1 }}>
          <Checkbox
            checked={isSelected}
            onChange={() => onSelect(camera)}
            color="primary"
          />
        </Box>
        <CardContent sx={{ flexGrow: 1, py: 1.5 }}>
          <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <Box>
              <Typography variant="h6" component="div">
                Camera {camera.id}
              </Typography>
              <Typography variant="body2" color="text.secondary">
                {camera.ip}
              </Typography>
            </Box>            <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
              <Chip 
                label={camera.isFake ? 'Simulation' : 'Connected'} 
                color={getStatusColor()} 
                size="small" 
              />
              {!camera.isFake && (
                <Tooltip title="Launch VLC with stream">
                  <IconButton 
                    color="primary" 
                    size="small"
                    onClick={handleLaunchVLC}
                  >
                    <PlayArrowIcon />
                  </IconButton>
                </Tooltip>
              )}
            </Box>
          </Box>
        </CardContent>
      </Card>
    );
  }

  // Original card design for non-compact mode
  return (    <Card sx={{ 
      height: '100%', 
      display: 'flex', 
      flexDirection: 'column',
      transition: 'all 0.2s ease-in-out',
      border: 'none'
    }}>      <CardMedia
        component="div"
        sx={{
          height: 140,
          backgroundColor: '#2c3e50',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          transition: 'background-color 0.3s'
        }}
      >
        <VideocamIcon sx={{ fontSize: 60, color: '#ecf0f1' }} />
      </CardMedia>
      <CardContent sx={{ flexGrow: 1 }}>
        <Typography 
          gutterBottom 
          variant="h5" 
          component="div" 
          sx={{ 
            display: 'flex', 
            justifyContent: 'space-between',
            alignItems: 'center' 
          }}        >
          Camera {camera.id}
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
            <Chip 
              label={camera.isFake ? 'Simulation' : 'Connected'} 
              color={getStatusColor()} 
              size="small" 
            />
            {!camera.isFake && (
              <Tooltip title="Launch VLC with stream">
                <IconButton 
                  color="primary"
                  size="small"
                  onClick={handleLaunchVLC}
                >
                  <PlayArrowIcon />
                </IconButton>
              </Tooltip>
            )}
          </Box>
        </Typography>
        <Typography variant="body2" color="text.secondary">
          IP Address: {camera.ip}
        </Typography>
      </CardContent>
      <CardActions>        <Button 
          size="small" 
          startIcon={<SettingsIcon />} 
          variant={isSelected ? "outlined" : "contained"}
          fullWidth
          onClick={() => onSelect(camera)}
          color={isSelected ? "success" : "primary"}
        >
          {isSelected ? '✓ Selected' : 'Select'}
        </Button>
      </CardActions>
    </Card>
  );
};

export default CameraCard;
