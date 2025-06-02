import { useState, useEffect } from 'react';
import { 
  Container, 
  Typography, 
  Box, 
  Alert,
  Button,
  Paper
} from '@mui/material';

const TestDashboard = () => {
  const [status, setStatus] = useState('Loading...');

  useEffect(() => {
    console.log('TestDashboard mounted');
    setStatus('Ready');
  }, []);

  console.log('TestDashboard rendering with status:', status);

  return (
    <Container maxWidth="lg">
      <Box sx={{ my: 4 }}>
        <Typography variant="h4" component="h1" gutterBottom>
          Test Dashboard - {status}
        </Typography>
        <Paper sx={{ p: 3 }}>
          <Typography variant="body1">
            If you can see this, the basic components are working.
          </Typography>
          <Button variant="contained" sx={{ mt: 2 }}>
            Test Button
          </Button>
        </Paper>
      </Box>
    </Container>
  );
};

export default TestDashboard;
