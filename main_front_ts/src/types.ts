export interface Camera {
  id: string;
  ip: string;
  isFake: boolean;
  // Add other relevant camera properties if known from backend config/models
  // e.g., name, location, etc.
}

export interface ApplyConfigPayload {
  cameraId: string;
  width: number;
  height: number;
  fps: number;
}

export interface ApplyConfigResponse {
  status: string; // Assuming the backend returns { "status": "configuration applied" }
}

// Add other types as needed based on backend API responses or frontend data structures 