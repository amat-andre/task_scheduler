package server

import (
	"log"
	"net/http"
	"os"
	"time"

	"task_scheduler/internal/api"
)

const (
	defPort = "7540"
	webDir  = "./web"
)

type Server struct {
	Log  *log.Logger
	Serv http.Server
}

func NewRout(logger *log.Logger) *Server {
	mux := http.NewServeMux()
	serverPass := getServPass()

	mux.Handle("GET /", http.FileServer(http.Dir(webDir)))
	api.Init(mux, serverPass)

	server := Server{
		Log: logger,
		Serv: http.Server{
			Addr:        ":" + getPort(),
			Handler:     mux,
			ErrorLog:    logger,
			ReadTimeout: 5 * time.Second,
		},
	}

	return &server
}

func getPort() string {
	if port := os.Getenv("TODO_PORT"); len(port) > 0 {
		return port
	}
	return defPort
}

func getServPass() string {
	if servPass := os.Getenv("TODO_PASSWORD"); len(servPass) > 0 {
		return servPass
	}
	return ""
}