package main

import (
	"log"

	"task_scheduler/internal/db"
	srv "task_scheduler/internal/server"
)

func main() {
	logger := log.New(log.Writer(), "", log.LstdFlags|log.Lshortfile)
	server := srv.NewRout(logger)

	err := db.Init()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	err = server.Serv.ListenAndServe()
	if err != nil {
		logger.Println(err)
	}
}
