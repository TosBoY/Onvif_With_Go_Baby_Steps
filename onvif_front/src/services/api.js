import axios from 'axios';

const api = {
  getCameraInfo: async () => {
    try {
      const response = await axios.get('/api/camera/info');
      return response.data;
    } catch (error) {
      console.error('Error fetching camera info:', error);
      throw error;
    }
  },

  addCamera: async (cameraData) => {
    try {
      const response = await axios.post('/api/cameras', cameraData);
      return response.data;
    } catch (error) {
      console.error('Error adding camera:', error);
      throw error;
    }
  },

  getResolutions: async (configToken, profileToken) => {
    try {
      const response = await axios.get('/api/camera/resolutions', {
        params: { configToken, profileToken }
      });
      return response.data;
    } catch (error) {
      console.error('Error fetching resolutions:', error);
      throw error;
    }
  },

  getResolutionsSimple: async () => {
    try {
      const response = await axios.get('/api/camera/resolutions-simple');
      return response.data;
    } catch (error) {
      console.error('Error fetching resolutions:', error);
      throw error;
    }
  },

  changeResolution: async (configData) => {
    try {
      // configData now includes cameraId and other configuration details
      const response = await axios.post('/api/camera/change-resolution', configData);
      return response.data;
    } catch (error) {
      console.error('Error changing resolution:', error);
      throw error;
    }
  },

  changeResolutionSimple: async (configData) => {
    try {
      const response = await axios.post('/api/camera/change-resolution-simple', configData);
      return response.data;
    } catch (error) {
      console.error('Error changing resolution:', error);
      throw error;
    }
  },
  
  getStreamUrl: async (profileToken) => {
    try {
      const response = await axios.get('/api/camera/stream-url', {
        params: { profileToken }
      });
      return response.data.streamUrl;
    } catch (error) {
      console.error('Error fetching stream URL:', error);
      throw error;
    }
  },
  
  launchVLC: async (profileToken) => {
    try {
      const response = await axios.post('/api/camera/launch-vlc', { profileToken });
      return response.data;
    } catch (error) {
      console.error('Error launching VLC:', error);
      throw error;
    }
  },
  
  getSingleConfig: async (configToken) => {
    try {
      const response = await axios.get('/api/camera/config', {
        params: { configToken }
      });
      return response.data;
    } catch (error) {
      console.error('Error fetching config details:', error);
      throw error;
    }
  },
  
  getDeviceInfo: async () => {
    try {
      const response = await axios.get('/api/camera/device-info');
      return response.data;
    } catch (error) {
      console.error('Error fetching device information:', error);
      throw error;
    }
  },

  getCameraDetailsSimple: async () => {
    try {
      const response = await axios.get('/api/camera/details-simple');
      return response.data;
    } catch (error) {
      console.error('Error fetching camera details:', error);
      throw error;
    }
  },

  deleteCamera: async (cameraId) => {
    try {
      const response = await axios.delete(`/api/cameras/${cameraId}`);
      return response.data;
    } catch (error) {
      console.error('Error deleting camera:', error);
      throw error;
    }
  },
};

export default api;