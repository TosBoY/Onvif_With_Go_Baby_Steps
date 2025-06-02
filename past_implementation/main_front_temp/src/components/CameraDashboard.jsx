import React, { useState, useEffect } from 'react';
import api from '../services/api';

function CameraDashboard() {
  const [cameras, setCameras] = useState([]);
  const [selectedCamera, setSelectedCamera] = useState(null);
  const [resolutions, setResolutions] = useState([]);
  const [streamUrl, setStreamUrl] = useState('');

  useEffect(() => {
    const fetchCameras = async () => {
      try {
        const data = await api.getCameras();
        setCameras(data);
      } catch (error) {
        console.error('Failed to fetch cameras:', error);
      }
    };
    fetchCameras();
  }, []);

  const handleCameraSelect = async (cameraId) => {
    try {
      const [info, res] = await Promise.all([
        api.getCameraInfo(cameraId),
        api.getResolutions(cameraId)
      ]);
      setSelectedCamera(info);
      setResolutions(res);
    } catch (error) {
      console.error('Failed to fetch camera details:', error);
    }
  };

  const handleApplyConfig = async (width, height, fps) => {
    if (!selectedCamera) return;
    try {
      await api.applyConfig(selectedCamera.id, width, height, fps);
      alert('Configuration applied successfully!');
    } catch (error) {
      console.error('Failed to apply config:', error);
      alert('Failed to apply configuration');
    }
  };

  const handleGetStream = async () => {
    if (!selectedCamera) return;
    try {
      const { streamUrl } = await api.getStreamUrl(selectedCamera.id);
      setStreamUrl(streamUrl);
    } catch (error) {
      console.error('Failed to get stream URL:', error);
    }
  };

  return (
    <div>
      <h1>Camera Dashboard</h1>
      <div>
        <h2>Cameras</h2>
        <ul>
          {cameras.map(camera => (
            <li key={camera.id} onClick={() => handleCameraSelect(camera.id)}>
              {camera.name}
            </li>
          ))}
        </ul>
      </div>
      
      {selectedCamera && (
        <div>
          <h2>Camera Details</h2>
          <p>Name: {selectedCamera.name}</p>
          <p>IP: {selectedCamera.ip}</p>
          
          <h3>Resolutions</h3>
          <ul>
            {resolutions.map((res, index) => (
              <li key={index}>
                {res.width}x{res.height} @ {res.fps}fps
                <button onClick={() => handleApplyConfig(res.width, res.height, res.fps)}>
                  Apply
                </button>
              </li>
            ))}
          </ul>
          
          <button onClick={handleGetStream}>Get Stream URL</button>
          {streamUrl && (
            <div>
              <h3>Stream URL</h3>
              <a href={streamUrl} target="_blank" rel="noopener noreferrer">
                {streamUrl}
              </a>
            </div>
          )}
        </div>
      )}
    </div>
  );
}

export default CameraDashboard; 