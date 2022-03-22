package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
)

func (h *Handler) GetWithdrawals() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := h.getAuthUser(r)
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
