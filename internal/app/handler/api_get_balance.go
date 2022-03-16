package handler

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) GetBalance() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		login, err := h.cookieAuthenticator.GetLogin(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		user, err := h.userRepository.GetByLogin(r.Context(), login)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		res, err := json.Marshal(user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(res)
	}
}
