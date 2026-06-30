package helpers

import (
	"encoding/json"
	"errors"
	"net/http"

	"task_scheduler/internal/db"
)

func WriteJSON(w http.ResponseWriter, statusCode int, data any) {
	js, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)

	w.Write(js)
}

func WriteErrorDB(w http.ResponseWriter, err error) {
	if errors.Is(err, db.ErrorTaskNotFound) {
		WriteJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
}