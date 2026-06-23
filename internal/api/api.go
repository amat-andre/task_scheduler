package api

import (
	"net/http"

hdl "task_scheduler/internal/handlers"
)

func Init(mux *http.ServeMux) {
    mux.HandleFunc("GET /api/nextdate", hdl.NextHandler) // cделать только Get обработчиком
	mux.HandleFunc("/api/task", hdl.TaskHandler) // cделать только Post обработчиком ** уже не надо, там свитч на выбор метода внутри
	mux.HandleFunc("POST /api/task/done", hdl.TaskDoneHandler) // post
	mux.HandleFunc("GET /api/tasks", hdl.TasksHandler) // обработчик для get запроса
} 
