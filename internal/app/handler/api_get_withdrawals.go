package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"sort"
)

func (h *Handler) GetWithdrawals() http.HandlerFunc {
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

		withdrawals, err := h.withdrawalRepository.GetByUserID(r.Context(), user.ID)
		if err != nil {
			if errors.As(err, &sql.ErrNoRows) {
				w.WriteHeader(http.StatusNoContent)
				return
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		sort.Slice(withdrawals, func(i, j int) bool {
			return withdrawals[i].ProcessedAt.Before(withdrawals[j].ProcessedAt)
		})

		res, err := json.Marshal(withdrawals)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(res)
	}
}
