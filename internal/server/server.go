package server

import (
	"log"
	"net/http"
	"time"

	"task_scheduler/internal/api"
	help "task_scheduler/internal/helpers"
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

	mux.Handle("/", http.FileServer(http.Dir(webDir)))
	api.Init(mux)


	server := Server{
		Log: logger,
		Serv: http.Server{
			Addr: ":" + help.GetPort(defPort),
			Handler: mux,
			ErrorLog: logger,
			ReadTimeout:  5 * time.Second,
		},
	}

	return &server
}