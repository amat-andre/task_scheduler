package handlers

import (
	"net/http"

	"task_scheduler/internal/db"
	help "task_scheduler/internal/helpers"
)
 
type TasksResp struct {
    Tasks []*db.Task `json:"tasks"`
}

func TasksHandler(w http.ResponseWriter, req *http.Request){
	/*
	метод возможно так проверть
	if req.Method != http.MethodGet {
		help.WriteJSON(w, http.StatusMethodNotAllowed, "Method not allowed") 
		return
	}
	*/
	defaultLimit := 50
	tasks, err := db.Tasks(defaultLimit)
    if err != nil {
       	help.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
        return
    }

    help.WriteJSON(w, http.StatusOK, TasksResp{ 
        Tasks: tasks,
    })
}