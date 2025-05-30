import { useState, useEffect } from 'react';
import { FormControl, InputLabel, Select, MenuItem, Box, Typography } from '@mui/material';

const ProfileSelector = ({ profiles, onChange, selectedProfile }) => {
  const [profile, setProfile] = useState(selectedProfile || '');

  useEffect(() => {
    if (selectedProfile) {
      setProfile(selectedProfile);
    }
  }, [selectedProfile]);

  const handleChange = (event) => {
    const value = event.target.value;
    setProfile(value);
    onChange(value);
  };

  const getProfileName = (profile) => {
    // Try to access name with various possible paths
    return profile.Name || 
           (profile.name) || 
           'Unnamed Profile';
  };

  const getProfileToken = (profile) => {
    // Try to access token with various possible paths
    return profile.Token || profile.token || 'Unknown';
  };

  const getVideoSourceInfo = (profile) => {
    try {
      // Try to get video source configuration info
      const videoSource = profile.VideoSourceConfiguration || profile.VideoSource || profile.videoSourceConfiguration;
      if (videoSource) {
        const sourceName = videoSource.Name || videoSource.name || '';
        const sourceToken = videoSource.SourceToken || videoSource.sourceToken || '';
        return ` - Source: ${sourceName || sourceToken}`;
      }
    } catch (error) {
      console.log('Error getting video source info:', error);
    }
    return '';
  };

  return (
    <Box sx={{ mb: 3 }}>
      <Typography variant="h6" sx={{ mb: 1 }}>Camera Profile</Typography>
      <FormControl fullWidth>
        <InputLabel id="profile-select-label">Select Profile</InputLabel>
        <Select
          labelId="profile-select-label"
          id="profile-select"
          value={profile}
          label="Select Profile"
          onChange={handleChange}
        >
          {profiles && profiles.map((p) => (
            <MenuItem key={getProfileToken(p)} value={getProfileToken(p)}>
              {getProfileName(p)}{getVideoSourceInfo(p)} (Token: {getProfileToken(p)})
            </MenuItem>
          ))}
        </Select>
      </FormControl>
    </Box>
  );
};

export default ProfileSelector;