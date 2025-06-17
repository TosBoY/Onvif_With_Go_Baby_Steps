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
  timeout: 60000, // 60 second timeout - increased because camera operations take time
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

export const applyConfig = async (cameraIds, width, height, fps, bitrate) => {
  try {
    // Handle both single camera ID (string) and multiple camera IDs (array)
    const isBatchMode = Array.isArray(cameraIds);
    const payload = isBatchMode 
      ? { cameraIds, width, height, fps, bitrate }
      : { cameraId: cameraIds, width, height, fps, bitrate };
      
    console.log('Applying config:', payload);
    const response = await api.post('/apply-config', payload);
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

export const launchVLC = async (cameraId) => {
  try {
    console.log('Launching VLC for camera:', cameraId);
    const response = await api.post('/vlc', {
      cameraId
    });
    return response.data;
  } catch (error) {
    console.error('Error launching VLC:', error);
    if (error.response) {
      throw new Error(`Server error: ${error.response.status} ${error.response.statusText}`);
    } else if (error.request) {
      throw new Error('Network error: No response from server');
    } else {
      throw new Error(`Request error: ${error.message}`);
    }
  }
};

export const addNewCamera = async (ip, port, url, username, password, isFake = false) => {
  try {
    console.log('Adding new camera:', { ip, port, url, username, isFake });
    const response = await api.post('/cameras', {
      ip,
      port,
      url,
      username,
      password,
      isFake
    });
    return response.data;
  } catch (error) {
    console.error('Error adding new camera:', error);
    if (error.response) {
      throw new Error(`Server error: ${error.response.status} ${error.response.statusText}`);
    } else if (error.request) {
      throw new Error('Network error: No response from server');
    } else {
      throw new Error(`Request error: ${error.message}`);
    }
  }
};

export const deleteCamera = async (cameraId) => {
  try {
    console.log('Deleting camera with ID:', cameraId);
    const response = await api.delete(`/cameras/${cameraId}`);
    return response.data;
  } catch (error) {
    console.error('Error deleting camera:', error);
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
