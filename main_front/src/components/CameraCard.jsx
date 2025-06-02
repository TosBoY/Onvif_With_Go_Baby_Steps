import { 
  Card, 
  CardContent, 
  Typography, 
  Button, 
  CardMedia,
  CardActions, 
  Chip,
  Box
} from '@mui/material';
import { Videocam as VideocamIcon, Settings as SettingsIcon } from '@mui/icons-material';

const CameraCard = ({ camera, isSelected, onSelect, compact = false }) => {
  const getStatusColor = () => {
    if (camera.isFake) return 'warning';
    return 'success';
  };

  if (compact) {
    return (
      <Card sx={{ 
        display: 'flex',
        alignItems: 'center',
        transition: 'all 0.2s ease-in-out',
        border: isSelected ? '2px solid #3498db' : '1px solid #e0e0e0',
        '&:hover': {
          backgroundColor: '#f5f5f5',
          boxShadow: 2
        }
      }}>
        <Box sx={{
          width: 60,
          height: 60,
          backgroundColor: isSelected ? '#3498db' : '#2c3e50',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          transition: 'background-color 0.3s'
        }}>
          <VideocamIcon sx={{ fontSize: 30, color: '#ecf0f1' }} />
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
            </Box>
            <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
              <Chip 
                label={camera.isFake ? 'Simulation' : 'Connected'} 
                color={getStatusColor()} 
                size="small" 
              />              <Button 
                size="small" 
                variant={isSelected ? "contained" : "outlined"}
                onClick={() => onSelect(camera)}
                color={isSelected ? "success" : "primary"}
              >
                {isSelected ? '✓ Selected' : 'Select'}
              </Button>
            </Box>
          </Box>        </CardContent>
      </Card>
    );
  }

  // Original card design for non-compact mode
  return (
    <Card sx={{ 
      height: '100%', 
      display: 'flex', 
      flexDirection: 'column',
      transition: 'all 0.2s ease-in-out',
      border: isSelected ? '2px solid #3498db' : 'none',
      '&:hover': {
        transform: 'translateY(-5px)',
        boxShadow: 5
      }
    }}>
      <CardMedia
        component="div"
        sx={{
          height: 140,
          backgroundColor: isSelected ? '#3498db' : '#2c3e50',
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
          }}
        >
          Camera {camera.id}
          <Chip 
            label={camera.isFake ? 'Simulation' : 'Connected'} 
            color={getStatusColor()} 
            size="small" 
          />
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
