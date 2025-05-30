import axios from 'axios';
import type { Camera, ApplyConfigPayload, ApplyConfigResponse } from '../types';

const API_BASE_URL = '/api'; // Assuming your backend is proxied under /api, adjust if necessary

const api = {
  /**
   * Fetches the list of cameras from the backend.
   * Corresponds to GET /cameras.
   * @returns A promise resolving with an array of Camera objects.
   */
  getCameras: async (): Promise<Camera[]> => {
    try {
      const response = await axios.get<Camera[]>(`${API_BASE_URL}/cameras`);
      // Check if the received data is an array before returning
      if (Array.isArray(response.data)) {
        return response.data;
      } else {
        console.error('Received data for /cameras is not an array:', response.data);
        // Depending on expected backend behavior, you might throw an error
        // or return an empty array. Returning empty array to prevent frontend errors.
        return [];
      }
    } catch (error) {
      console.error('Error fetching cameras:', error);
      // It's better to throw the error so the calling component can handle it
      throw error;
    }
  },

  /**
   * Applies a new configuration to a camera via the backend.
   * Corresponds to POST /apply-config.
   * @param configData - The configuration payload.
   * @returns A promise resolving with the ApplyConfigResponse.
   */
  applyConfig: async (configData: ApplyConfigPayload): Promise<ApplyConfigResponse> => {
    try {
      const response = await axios.post<ApplyConfigResponse>(`${API_BASE_URL}/apply-config`, configData);
      return response.data;
    } catch (error) {
      console.error('Error applying config:', error);
      throw error;
    }
  },

  // Add other API calls here as you implement more backend endpoints
};

export default api; 