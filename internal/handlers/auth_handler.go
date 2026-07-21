package handlers

import (
	"encoding/json"
	"net/http"

	"task_scheduler/internal/auth"
	help "task_scheduler/internal/helpers"
)

type SignInRequest struct {
	Pass string `json:"password"`
}

func AuthHandler(servPass string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			help.WriteJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
			return
		}

		if servPass == "" {
			help.WriteJSON(w, http.StatusOK, map[string]string{"error": "authentication is not required"})
			return
		}

		var input SignInRequest
		decoder := json.NewDecoder(req.Body)

		err := decoder.Decode(&input)
		if err != nil {
			help.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}

		if input.Pass == "" {
			help.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "password is not entered"})
			return
		}

		if input.Pass != servPass {
			help.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid password"})
			return
		}

		token, err := auth.GenerateToken(input.Pass)
		if err != nil {
			help.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    token,
			HttpOnly: true,
		})

		help.WriteJSON(w, http.StatusOK, map[string]string{"token": token})
	}
}
