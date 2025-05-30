import { useState, useEffect } from 'react';
import { FormControl, InputLabel, Select, MenuItem, Box, Typography } from '@mui/material';

const ConfigSelector = ({ configs, onChange, selectedConfig }) => {
  const [config, setConfig] = useState(selectedConfig || '');

  useEffect(() => {
    if (selectedConfig) {
      setConfig(selectedConfig);
    }
  }, [selectedConfig]);

  const handleChange = (event) => {
    const value = event.target.value;
    setConfig(value);
    onChange(value);
  };

  return (
    <Box sx={{ mb: 3 }}>
      <Typography variant="h6" sx={{ mb: 1 }}>Video Encoder Configuration</Typography>
      <FormControl fullWidth>
        <InputLabel id="config-select-label">Select Configuration</InputLabel>
        <Select
          labelId="config-select-label"
          id="config-select"
          value={config}
          label="Select Configuration"
          onChange={handleChange}
        >
          {configs && configs.map((c) => (
            <MenuItem key={c.token || c.Token} value={c.token || c.Token}>
              {c.name || c.Name || 'Unnamed Config'} ({c.width || c.Width || '?'}x{c.height || c.Height || '?'}, {c.encoding || c.Encoding || '?'})
            </MenuItem>
          ))}
        </Select>
      </FormControl>
    </Box>
  );
};

export default ConfigSelector;