package helpers

import (
	"os"
	"encoding/json"
	"errors"
	"net/http"
)

var (
	ErrorIncorrectDate = errors.New("incorrect date format")
	ErrorTaskNotFound = errors.New("identifier task not found")
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
	if errors.Is(err, ErrorTaskNotFound) {
		WriteJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
}

func GetPort(defPort string) string {
	if port := os.Getenv("TODO_PORT"); port != "" {
		return port
	}
	return defPort
}

func GetFileDB(defFile string) string {
	if dbFile := os.Getenv("TODO_DBFILE"); dbFile != "" {
		return dbFile
	}
	return defFile
}