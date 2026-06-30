package server

import (
	"os"
	"log"
	"net/http"
	"time"

	"task_scheduler/internal/api"
)

const(
	defPort = "7540"
	webDir = "./web"
)

type Server struct {
	Log *log.Logger
	Serv http.Server
}

func NewRout(logger *log.Logger) *Server {
	mux := http.NewServeMux()

	mux.Handle("GET /", http.FileServer(http.Dir(webDir)))
	api.Init(mux)

	server := Server{
		Log: logger,
		Serv: http.Server{
			Addr: ":" + getPort(),
			Handler: mux,
			ErrorLog: logger,
			ReadTimeout:  5 * time.Second,
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