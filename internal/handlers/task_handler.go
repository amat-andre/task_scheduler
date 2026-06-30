package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"task_scheduler/internal/db"
	help "task_scheduler/internal/helpers"
	"task_scheduler/internal/service"
)

func AddTaskHandler(w http.ResponseWriter, req *http.Request){
	var task db.Task
	decoder := json.NewDecoder(req.Body)

	err := decoder.Decode(&task)
	if err != nil {
		help.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	if strings.TrimSpace(task.Title) == "" {
		help.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "title is empty"})
		return
	}

	err = service.CheckDate(&task)
	if err != nil {
		help.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		help.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	help.WriteJSON(w, http.StatusCreated, map[string]string{"id": fmt.Sprint(id)})	
}

func GetTaskHandler(w http.ResponseWriter, req *http.Request){
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

	help.WriteJSON(w, http.StatusOK, task)
}

func UpdateTaskHandler(w http.ResponseWriter, req *http.Request){
	var task db.Task
	decoder := json.NewDecoder(req.Body)

	err := decoder.Decode(&task)
	if err != nil {
		help.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	if strings.TrimSpace(task.ID) == "" {
		help.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "id is empty"})
		return
	}
	if strings.TrimSpace(task.Title) == "" {
		help.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "title is empty"})
		return
	}

	err = service.CheckDate(&task)
	if err != nil {
		help.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	err = db.UpdateTask(&task)
	if err != nil {
		help.WriteErrorDB(w, err)
		return
	}

	help.WriteJSON(w, http.StatusOK, map[string]string{})
}

func DeleteTaskHandler(w http.ResponseWriter, req *http.Request){
	id := strings.TrimSpace(req.URL.Query().Get("id"))
	if id == "" {
		help.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "id is empty"})
		return
	}

	err := db.DeleteTask(id)
	if err != nil {
		help.WriteErrorDB(w, err)
		return
	}

	help.WriteJSON(w, http.StatusOK, map[string]string{})
}