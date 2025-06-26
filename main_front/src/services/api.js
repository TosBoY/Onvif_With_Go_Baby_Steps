import axios from 'axios';

// Use environment variable, proxy path for development, or default URL
const isDevelopment = import.meta.env.DEV;
const API_BASE_URL = isDevelopment 
  ? '/api' // Use Vite's proxy in development
  : '/api'; // Use relative path for production (embedded mode)

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

export const applyConfig = async (cameraIds, width, height, fps, bitrate, encoding) => {
  try {
    // Handle both single camera ID (string) and multiple camera IDs (array)
    const isBatchMode = Array.isArray(cameraIds);
    const payload = isBatchMode 
      ? { cameraIds, width, height, fps, bitrate, encoding }
      : { cameraId: cameraIds, width, height, fps, bitrate, encoding };
      
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

// Get camera detailed information
export const getCameraInfo = async (cameraId) => {
  try {
    console.log(`Getting camera info for ID: ${cameraId}`);
    const response = await api.get(`/cameras/${cameraId}/info`);
    return response.data;
  } catch (error) {
    console.error(`Error getting camera info for ${cameraId}:`, error);
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
    console.log(`Launching VLC for camera ID: ${cameraId}`);
    const response = await api.post('/vlc', { cameraId });
    return response.data;
  } catch (error) {
    console.error(`Error launching VLC for ${cameraId}:`, error);
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
    }  }
};

// Export validation results as CSV
export const exportValidationCSV = async (validation, configurationErrors = [], cameraOrder = []) => {
  try {
    const response = await api.post('/export-validation-csv', { 
      validation,
      configurationErrors,
      cameraOrder
    }, {
      responseType: 'blob',
    });
    
    // Create a blob from the response
    const blob = new Blob([response.data], { type: 'text/csv' });
    
    // Create a temporary URL for the blob
    const url = window.URL.createObjectURL(blob);
      // Create a temporary anchor element and trigger download
    const a = document.createElement('a');
    a.href = url;
    
    // Generate local timestamp in format YYYY-MM-DDTHH-mm-ss
    const now = new Date();
    const localTimestamp = 
      `${now.getFullYear()}-` + 
      `${String(now.getMonth() + 1).padStart(2, '0')}-` +
      `${String(now.getDate()).padStart(2, '0')}T` +
      `${String(now.getHours()).padStart(2, '0')}-` +
      `${String(now.getMinutes()).padStart(2, '0')}-` +
      `${String(now.getSeconds()).padStart(2, '0')}`;
    
    a.download = `validation_results_${localTimestamp}.csv`;
    document.body.appendChild(a);
    a.click();
    
    // Clean up
    document.body.removeChild(a);
    window.URL.revokeObjectURL(url);
    
    return { success: true };
  } catch (error) {
    console.error('Error exporting CSV:', error);
    throw new Error(`Failed to export CSV: ${error.message}`);
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

// Import configuration from CSV
export const importConfigCSV = async (file) => {
  try {
    const formData = new FormData();
    formData.append('csvFile', file);

    const response = await api.post('/import-config-csv', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
    return response.data;
  } catch (error) {
    console.error('Error importing config CSV:', error);
    if (error.response) {
      throw new Error(error.response.data.error || `Server error: ${error.response.status}`);
    } else if (error.request) {
      throw new Error('Network error: No response from server');
    } else {
      throw new Error(`Request error: ${error.message}`);
    }
  }
};

// Check all cameras status from CSV
export const checkAllCameras = async () => {
  try {
    console.log('Checking all cameras from CSV');
    const response = await api.get('/check-all-cams');
    return response.data;
  } catch (error) {
    console.error('Error checking all cameras:', error);
    if (error.response) {
      throw new Error(`Server error: ${error.response.status} ${error.response.statusText}`);
    } else if (error.request) {
      throw new Error('Network error: No response from server');
    } else {
      throw new Error(`Request error: ${error.message}`);
    }
  }
};

export default api;
