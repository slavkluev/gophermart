package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"sort"
)

func (h *Handler) GetOrders() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := h.getAuthUser(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		orders, err := h.orderRepository.GetByUserID(r.Context(), user.ID)
		if err != nil {
			if errors.As(err, &sql.ErrNoRows) {
				w.WriteHeader(http.StatusNoContent)
				return
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		sort.Slice(orders, func(i, j int) bool {
			return orders[i].UploadedAt.Before(orders[j].UploadedAt)
		})

		res, err := json.Marshal(orders)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(res)
	}
}
