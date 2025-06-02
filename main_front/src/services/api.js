import axios from 'axios';

const API_BASE_URL = '/api';

const api = {
  getCameras: async () => {
    try {
      const response = await axios.get(`${API_BASE_URL}/cameras`);
      return response.data;
    } catch (error) {
      console.error('Error fetching cameras:', error);
      throw error;
    }
  },

  applyConfig: async (cameraId, width, height, fps) => {
    try {
      const response = await axios.post(`${API_BASE_URL}/apply-config`, {
        cameraId,
        width,
        height,
        fps
      });
      return response.data;
    } catch (error) {
      console.error('Error applying config:', error);
      throw error;
    }
  },

  getCameraInfo: async (cameraId) => {
    try {
      const response = await axios.get(`${API_BASE_URL}/camera/info`, {
        params: { cameraId }
      });
      return response.data;
    } catch (error) {
      console.error('Error fetching camera info:', error);
      throw error;
    }
  },

  getResolutions: async (cameraId) => {
    try {
      const response = await axios.get(`${API_BASE_URL}/camera/resolutions`, {
        params: { cameraId }
      });
      return response.data;
    } catch (error) {
      console.error('Error fetching resolutions:', error);
      throw error;
    }
  },

  getStreamUrl: async (cameraId) => {
    try {
      const response = await axios.get(`${API_BASE_URL}/camera/stream-url`, {
        params: { cameraId }
      });
      return response.data;
    } catch (error) {
      console.error('Error fetching stream URL:', error);
      throw error;
    }
  }
};

export default api; 