package api

import (
	"net/http"

	"task_scheduler/internal/auth"
	hdl "task_scheduler/internal/handlers"
	mdw "task_scheduler/internal/middleware"
)

func Init(mux *http.ServeMux) {
	servPass := auth.GetServPass()

	mux.HandleFunc("POST /api/signin", hdl.AuthHandler)
	mux.HandleFunc("GET /api/nextdate", hdl.NextHandler)

	mux.HandleFunc("GET /api/task", mdw.Middleware(servPass, hdl.GetTaskHandler))
	mux.HandleFunc("POST /api/task", mdw.Middleware(servPass, hdl.AddTaskHandler))
	mux.HandleFunc("PUT /api/task", mdw.Middleware(servPass, hdl.UpdateTaskHandler))
	mux.HandleFunc("DELETE /api/task", mdw.Middleware(servPass, hdl.DeleteTaskHandler))
 
	mux.HandleFunc("POST /api/task/done", mdw.Middleware(servPass, hdl.TaskDoneHandler))
	mux.HandleFunc("GET /api/tasks", mdw.Middleware(servPass, hdl.TasksHandler))
} 
