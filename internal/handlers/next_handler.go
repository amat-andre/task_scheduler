package handlers

import (
	"net/http"
	"strings"
	"time"

	"task_scheduler/internal/db"
	help "task_scheduler/internal/helpers"
	"task_scheduler/internal/service"
)

func NextHandler(w http.ResponseWriter, req *http.Request) {
	now := time.Now()
	currentDate := strings.TrimSpace(req.URL.Query().Get("now"))
	if !(currentDate == "") {
		nowParse, err := time.Parse(db.DateFormat, currentDate)
		if err != nil {
			help.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "incorrect now format"})
			return
		}
		now = nowParse
	}

	startDate := strings.TrimSpace(req.URL.Query().Get("date"))
	if startDate == "" {
		help.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "incorrect date format"})
		return
	}

	repeat := strings.TrimSpace(req.URL.Query().Get("repeat"))
	if repeat == "" {
		help.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "incorrect repeat format"})
		return
	}

	nextDate, err := service.NextDate(now, startDate, repeat)
	if err != nil {
		help.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	w.Write([]byte(nextDate))
}
