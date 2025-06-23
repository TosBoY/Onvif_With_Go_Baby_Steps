package main

import (
	"main_back/internal/api"
)

func main() {
	// cams, err := config.LoadCameraList()
	// if err != nil {
	// 	log.Fatalf("Failed to load camera config: %v", err)
	// }

	// err = camera.InitializeAllCameras(cams)
	// if err != nil {
	// 	log.Fatalf("Failed to connect to cameras: %v", err)
	// }

	api.StartServer(":8090")
}
