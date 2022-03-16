package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/slavkluev/gophermart/internal/app"
	"github.com/slavkluev/gophermart/internal/app/model"
)

func (h *Handler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		credentials := model.Credentials{}
		if err := json.Unmarshal(b, &credentials); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		user, err := h.userRepository.GetByLogin(r.Context(), credentials.Login)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if app.Hash(credentials.Password) != user.PasswordHash {
			http.Error(w, "wrong password", http.StatusUnauthorized)
			return
		}

		err = h.cookieAuthenticator.SetCookie(w, credentials.Login)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
