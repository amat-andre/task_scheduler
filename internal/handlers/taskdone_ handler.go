package handlers

import (
	"net/http"
	"strings"
	"time"

	"task_scheduler/internal/db"
	help "task_scheduler/internal/helpers"
	"task_scheduler/internal/service"
)

func TaskDoneHandler(w http.ResponseWriter, req *http.Request) {
	id := strings.TrimSpace(req.URL.Query().Get("id"))
	if id == "" {
		help.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "id is empty"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		help.WriteErrorDB(w, err)
		return
	}

	if task.Repeat == "" {
		err = db.DeleteTask(id)
		if err != nil {
			help.WriteErrorDB(w, err)
			return
		}
		help.WriteJSON(w, http.StatusOK, map[string]string{})
		return
	}

	now := time.Now()
	nextDate, err := service.NextDate(now, task.Date, task.Repeat)
	if err != nil {
		help.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	err = db.UpdateDate(nextDate, id)
	if err != nil {
		help.WriteErrorDB(w, err)
		return
	}

	help.WriteJSON(w, http.StatusOK, map[string]string{})
}
