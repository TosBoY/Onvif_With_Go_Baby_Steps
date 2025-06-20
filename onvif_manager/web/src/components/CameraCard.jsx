// filepath: d:\VNG\test\main_onvif\main_front\src\components\CameraCard.jsx
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
  PlayArrow as PlayArrowIcon,
  Info as InfoIcon 
} from '@mui/icons-material';
import { useState } from 'react';
import { launchVLC } from '../services/api';
import CameraInfoDialog from './CameraInfoDialog';

const CameraCard = ({ camera, isSelected, onSelect, compact = false, onCameraDeleted }) => {
  const [infoDialogOpen, setInfoDialogOpen] = useState(false);
  
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
  };
  
  const handleInfoClick = (e) => {
    e.stopPropagation(); // Prevent triggering camera selection
    setInfoDialogOpen(true);
  };

  if (compact) {      
    return (        
      <Card sx={{ 
        display: 'flex',
        alignItems: 'center',
        transition: 'all 0.2s ease-in-out',
        border: 'none',
        minHeight: '40px',  // Increased from 32px for less thin appearance
        maxHeight: '52px',  // Increased from 40px
        boxShadow: '0 1px 2px rgba(0,0,0,0.12)',  // Slightly stronger shadow
        borderRadius: 1     
      }}>        
        <Box sx={{ display: 'flex', alignItems: 'center', pl: 1 }}>
          <Checkbox
            checked={isSelected}
            onChange={() => onSelect(camera)}
            color="primary"
            size="small"
            sx={{ p: '4px', mr: 0 }}  // Increased padding from 2px to 4px
          />
        </Box>
        <CardContent sx={{ flexGrow: 1, py: 0.75, px: 1.5, '&:last-child': { pb: 0.75 } }}>  
          <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', width: '100%' }}>
            <Box sx={{ mr: 2, flexGrow: 1, overflow: 'hidden', textOverflow: 'ellipsis' }}>              
              <Typography variant="subtitle1" component="div" sx={{ lineHeight: 1.2, fontSize: '0.875rem', m: 0, fontWeight: 500, whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis' }}>  
                Camera {camera.id}
              </Typography>
              <Typography variant="body2" color="text.secondary" sx={{ lineHeight: 1.2, fontSize: '0.775rem', m: 0, mt: '2px', whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis' }}>  
                {camera.ip}
              </Typography>
            </Box>
            <Box sx={{ display: 'flex', alignItems: 'center', flexShrink: 0, gap: 0.75 }}>              
              <Chip 
                label={camera.isFake ? 'Simulation' : 'Connected'} 
                color={getStatusColor()} 
                size="small"
                sx={{ 
                  height: '20px', 
                  '& .MuiChip-label': { px: 0.75, fontSize: '0.7rem', py: 0 },
                  mr: camera.isFake ? 0 : 0.5
                }}
              />
              {!camera.isFake && (
                <Tooltip title="Launch VLC with stream">
                  <IconButton                    
                    color="primary" 
                    size="small"
                    onClick={handleLaunchVLC}
                    sx={{ padding: '4px', minWidth: '24px', minHeight: '24px' }}
                  >
                    <PlayArrowIcon sx={{ fontSize: '1rem' }} />
                  </IconButton>
                </Tooltip>
              )}              <Tooltip title="Camera Info">
                <IconButton                    
                  color="primary" 
                  size="small"
                  onClick={handleInfoClick}
                  sx={{ padding: '4px', minWidth: '24px', minHeight: '24px' }}
                >
                  <InfoIcon sx={{ fontSize: '1rem' }} />
                </IconButton>
              </Tooltip>
            </Box>
          </Box>
        </CardContent>
        
        <CameraInfoDialog
          open={infoDialogOpen}
          onClose={() => setInfoDialogOpen(false)}
          camera={camera}
          onCameraDeleted={onCameraDeleted}
        />
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
      border: 'none'
    }}>
      <CardMedia
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
          }}
        >
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
            )}            <Tooltip title="Camera Info">
              <IconButton 
                color="primary"
                size="small"
                onClick={handleInfoClick}
              >
                <InfoIcon />
              </IconButton>
            </Tooltip>
          </Box>
        </Typography>
        <Typography variant="body2" color="text.secondary">
          IP Address: {camera.ip}
        </Typography>
      </CardContent>
      <CardActions>
        <Button 
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
      
      <CameraInfoDialog
        open={infoDialogOpen}
        onClose={() => setInfoDialogOpen(false)}
        camera={camera}
        onCameraDeleted={onCameraDeleted}
      />
    </Card>
  );
};

export default CameraCard;
