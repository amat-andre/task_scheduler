package main

import (
	"log"

	    "task_scheduler/internal/db"
	srv "task_scheduler/internal/server"
   // "task_scheduler/internal/service"
   
)


func main() {
    
    logger := log.New(log.Writer(), "", log.LstdFlags|log.Lshortfile) // может логер не создавать отдельный а просто log использовать
	server := srv.NewRout(logger)

    err := db.Init()
    if err != nil {
		logger.Fatal(err)
	} else {
        defer db.Close() // может так? нужно узнать logger.Fatal(err) останавливает приложение или нет, от этого зависит наличие else
    }
    

	err = server.Serv.ListenAndServe()
	if err != nil {
		logger.Fatal(err)
	}

}  