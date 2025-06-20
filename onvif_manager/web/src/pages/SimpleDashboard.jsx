import React from 'react';
import { Container, Typography, Box } from '@mui/material';

// Simplified Dashboard for testing
const SimpleDashboard = () => {
  return (
    <Container maxWidth="lg">
      <Box sx={{ my: 4 }}>
        <Typography variant="h4" component="h1" gutterBottom>
          Simple Camera Dashboard
        </Typography>
        <Typography variant="body1">
          This is a simplified dashboard for testing purposes.
        </Typography>
      </Box>
    </Container>
  );
};

export default SimpleDashboard;
