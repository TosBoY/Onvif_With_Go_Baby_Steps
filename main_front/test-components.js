// Simple test to verify components can be imported without errors
import React from 'react';

// Test imports
try {
  console.log('Testing component imports...');
  
  // Import Dashboard
  import('./src/pages/Dashboard.jsx').then(() => {
    console.log('✓ Dashboard component imported successfully');
  }).catch(err => {
    console.error('✗ Dashboard import failed:', err);
  });
  
  // Import other components
  import('./src/components/CameraCard.jsx').then(() => {
    console.log('✓ CameraCard component imported successfully');
  }).catch(err => {
    console.error('✗ CameraCard import failed:', err);
  });
  
  import('./src/components/CameraConfigPanel.jsx').then(() => {
    console.log('✓ CameraConfigPanel component imported successfully');
  }).catch(err => {
    console.error('✗ CameraConfigPanel import failed:', err);
  });
  
  import('./src/services/api.js').then(() => {
    console.log('✓ API service imported successfully');
  }).catch(err => {
    console.error('✗ API service import failed:', err);
  });
  
} catch (error) {
  console.error('Test failed:', error);
}

console.log('Component import test completed. Check console for results.');
