package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/slavkluev/gophermart/internal/app"
	"github.com/slavkluev/gophermart/internal/app/model"
)

func (h *Handler) Register() http.HandlerFunc {
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

		_, err = h.userRepository.GetByLogin(r.Context(), credentials.Login)
		if err == nil {
			http.Error(w, "login has already been taken", http.StatusConflict)
			return
		}

		newUser := model.User{
			Login:        credentials.Login,
			PasswordHash: app.Hash(credentials.Password),
		}

		err = h.userRepository.Create(r.Context(), newUser)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
