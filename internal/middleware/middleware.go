package middleware

import (
	"net/http"

	"task_scheduler/internal/auth"
	help "task_scheduler/internal/helpers"
)

func Middleware(servPass string, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if servPass != "" {
			cookie, err := req.Cookie("token")
			if err != nil || cookie.Value == "" {
				help.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "Authentication required"})
				return
			}
			token := cookie.Value

			if !auth.ValidateToken(token, servPass) {
				help.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "Authentication required"})
				return
			}
		}
		next(w, req)
	})
}
