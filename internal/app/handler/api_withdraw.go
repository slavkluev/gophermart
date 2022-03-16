package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/slavkluev/gophermart/internal/app"
	"github.com/slavkluev/gophermart/internal/app/model"
	"github.com/slavkluev/gophermart/internal/app/repository"
)

func (h *Handler) Withdraw() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := h.getAuthUser(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		b, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		withdrawal := &model.Withdrawal{
			ProcessedAt: time.Now(),
			UserID:      user.ID,
		}
		err = json.Unmarshal(b, &withdrawal)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = app.CheckOrderNumber(withdrawal.Order)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}

		err = h.withdrawalRepository.Create(r.Context(), *withdrawal)
		if err != nil {
			if errors.As(err, &repository.ErrInsufficientBalance) {
				http.Error(w, "insufficient balance", http.StatusPaymentRequired)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
