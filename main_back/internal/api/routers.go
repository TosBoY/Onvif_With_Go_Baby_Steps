package api

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func StartServer(addr string) {
	r := mux.NewRouter()
	RegisterRoutes(r)

	log.Println("Listening on", addr)
	err := http.ListenAndServe(addr, r)
	if err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
