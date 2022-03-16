package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/slavkluev/gophermart/internal/app/model"
	"github.com/theplant/luhn"
)

func (h *Handler) Withdraw() http.HandlerFunc {
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

		b, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var withdrawal model.Withdrawal
		err = json.Unmarshal(b, &withdrawal)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		orderInt, err := strconv.Atoi(withdrawal.Order)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !luhn.Valid(orderInt) {
			http.Error(w, "invalid order number", http.StatusUnprocessableEntity)
			return
		}

		if user.Balance < withdrawal.Sum {
			http.Error(w, "insufficient balance", http.StatusPaymentRequired)
			return
		}

		err = h.userRepository.DecreaseBalanceByUserID(r.Context(), user.ID, withdrawal.Sum)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		withdrawal.ProcessedAt = time.Now()
		withdrawal.UserID = user.ID
		err = h.withdrawalRepository.Create(r.Context(), withdrawal)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		return
	}
}
