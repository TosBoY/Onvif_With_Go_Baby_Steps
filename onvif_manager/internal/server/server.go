package server

import (
	"log"
	"onvif_manager/internal/backend/api"
	"onvif_manager/internal/backend/camera"
	"onvif_manager/internal/backend/config"
)

func Run() {
	cams, err := config.LoadCameraList()
	if err != nil {
		log.Fatalf("Failed to load camera config: %v", err)
	}

	err = camera.InitializeAllCameras(cams)
	if err != nil {
		log.Fatalf("Failed to connect to cameras: %v", err)
	}

	api.StartServer(":8090")
}
