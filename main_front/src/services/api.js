import axios from 'axios';

// Use environment variable, proxy path for development, or default URL
const isDevelopment = import.meta.env.DEV;
const API_BASE_URL = isDevelopment 
  ? '/api' // Use Vite's proxy in development
  : (import.meta.env.VITE_API_BASE_URL || 'http://localhost:8090');

// Create an axios instance with default config
const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
  timeout: 10000, // 10 second timeout
});

// API methods
export const getCameras = async () => {
  try {
    console.log('Making API request to:', `${API_BASE_URL}/cameras`);
    const response = await api.get('/cameras');
    console.log('API Response:', response);
    // Backend returns cameras array directly, not wrapped in an object
    return response.data;
  } catch (error) {
    console.error('Error fetching cameras:', error);
    if (error.code === 'ECONNREFUSED') {
      throw new Error('Cannot connect to backend server. Please make sure it is running on port 8090.');
    } else if (error.response) {
      throw new Error(`Server error: ${error.response.status} ${error.response.statusText}`);
    } else if (error.request) {
      throw new Error('Network error: No response from server');
    } else {
      throw new Error(`Request error: ${error.message}`);
    }
  }
};

export const applyConfig = async (cameraId, width, height, fps) => {
  try {
    console.log('Applying config:', { cameraId, width, height, fps });
    const response = await api.post('/apply-config', {
      cameraId,
      width,
      height,
      fps
    });
    return response.data;
  } catch (error) {
    console.error('Error applying configuration:', error);
    if (error.response) {
      throw new Error(`Server error: ${error.response.status} ${error.response.statusText}`);
    } else if (error.request) {
      throw new Error('Network error: No response from server');
    } else {
      throw new Error(`Request error: ${error.message}`);
    }
  }
};

// Test connection to backend server
export const testConnection = async () => {
  try {
    const response = await api.get('/cameras');
    return { 
      success: true,
      message: 'Connected to backend'
    };
  } catch (error) {
    console.error('Backend connection test failed:', error);
    return {
      success: false,
      message: `Connection failed: ${error.message}`
    };
  }
};

export default api;
