package handlers

import(
	"net/http"
	"time"
	"strings"

	"task_scheduler/internal/service"
)
// можно переделать ошибки в объявляемые через var, как в nextdate.go и проверить в остальных пакетах
func NextHandler(w http.ResponseWriter, req *http.Request){
	now := time.Now()
	currentDate:= req.URL.Query().Get("now")
	if !(strings.TrimSpace(currentDate) == ""){
		nowParse, err := time.Parse(service.DateFormat, currentDate)
		if err != nil {
			http.Error(w, "incorrect now format", http.StatusBadRequest)
			return
		}
		now = nowParse
	}

	startDate := req.URL.Query().Get("date")
	if strings.TrimSpace(startDate) == ""{
		http.Error(w, "incorrect date format", http.StatusBadRequest)
		return
	}
	
	repeat := req.URL.Query().Get("repeat")
	if strings.TrimSpace(repeat) == ""{
		http.Error(w, "incorrect repeat format", http.StatusBadRequest)
		return
	}

	nextDate, err := service.NextDate(now, startDate, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write([]byte(nextDate)) // здесь подумать, возможно стоит сразу заголовок записать тут же
}