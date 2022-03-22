package handler

import (
	"database/sql"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/slavkluev/gophermart/internal/app"
	"github.com/slavkluev/gophermart/internal/app/model"
)

func (h *Handler) CreateOrder() http.HandlerFunc {
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

		number := string(b)
		err = app.CheckOrderNumber(number)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}

		newOrder := model.Order{
			Number:     number,
			Status:     model.NEW,
			UploadedAt: time.Now(),
			UserID:     user.ID,
		}

		order, err := h.orderRepository.GetByNumber(r.Context(), number)
		if err != nil {
			if errors.As(err, &sql.ErrNoRows) {
				err = h.orderRepository.Create(r.Context(), newOrder)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				h.pointAccrualService.Accrue(newOrder.Number)

				w.WriteHeader(http.StatusAccepted)
				return
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		if order.UserID != user.ID {
			w.WriteHeader(http.StatusConflict)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
