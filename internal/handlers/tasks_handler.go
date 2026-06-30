package handlers

import (
	"net/http"
	"strings"

	"task_scheduler/internal/db"
	help "task_scheduler/internal/helpers"
)
 
type TasksResp struct {
    Tasks []*db.Task `json:"tasks"`
}

func TasksHandler(w http.ResponseWriter, req *http.Request){
	defaultLimit := 50
	search := strings.TrimSpace(req.URL.Query().Get("search"))

	tasks, err := db.Tasks(search, defaultLimit)
    if err != nil {
       	help.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
        return
    }

    help.WriteJSON(w, http.StatusOK, TasksResp{ 
        Tasks: tasks,
    })
}