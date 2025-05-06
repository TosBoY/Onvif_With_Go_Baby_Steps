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
            <MenuItem key={p.Token} value={p.Token}>
              {p.Name} (Token: {p.Token})
            </MenuItem>
          ))}
        </Select>
      </FormControl>
    </Box>
  );
};

export default ProfileSelector;