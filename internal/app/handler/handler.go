package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/slavkluev/gophermart/internal/app"
	"github.com/slavkluev/gophermart/internal/app/model"
)

type UserRepository interface {
	Create(ctx context.Context, user model.User) error
	GetByLogin(ctx context.Context, login string) (model.User, error)
}

type OrderRepository interface {
	Create(ctx context.Context, order model.Order) error
	GetByUserID(ctx context.Context, userID uint64) ([]model.Order, error)
	GetByNumber(ctx context.Context, number string) (model.Order, error)
}

type WithdrawalRepository interface {
	Create(ctx context.Context, withdrawal model.Withdrawal) error
	GetByUserID(ctx context.Context, userID uint64) ([]model.Withdrawal, error)
}

type CookieAuthenticator interface {
	SetCookie(w http.ResponseWriter, login string) error
}

type PointAccrualService interface {
	Accrue(order string)
}

type Middleware interface {
	Handle(next http.HandlerFunc) http.HandlerFunc
}

type Handler struct {
	*chi.Mux
	baseURL              string
	userRepository       UserRepository
	orderRepository      OrderRepository
	withdrawalRepository WithdrawalRepository
	cookieAuthenticator  CookieAuthenticator
	pointAccrualService  PointAccrualService
}

func NewHandler(
	baseURL string,
	userRepository UserRepository,
	orderRepository OrderRepository,
	withdrawalRepository WithdrawalRepository,
	cookieAuthenticator CookieAuthenticator,
	pointAccrualService PointAccrualService,
	authenticator Middleware,
	middlewares []Middleware,
) *Handler {
	h := &Handler{
		Mux:                  chi.NewMux(),
		baseURL:              baseURL,
		userRepository:       userRepository,
		orderRepository:      orderRepository,
		withdrawalRepository: withdrawalRepository,
		cookieAuthenticator:  cookieAuthenticator,
		pointAccrualService:  pointAccrualService,
	}

	h.Post("/api/user/register", applyMiddlewares(h.Register(), middlewares))
	h.Post("/api/user/login", applyMiddlewares(h.Login(), middlewares))

	h.Post("/api/user/orders", authenticator.Handle(applyMiddlewares(h.CreateOrder(), middlewares)))
	h.Get("/api/user/orders", authenticator.Handle(applyMiddlewares(h.GetOrders(), middlewares)))
	h.Get("/api/user/balance", authenticator.Handle(applyMiddlewares(h.GetBalance(), middlewares)))
	h.Post("/api/user/balance/withdraw", authenticator.Handle(applyMiddlewares(h.Withdraw(), middlewares)))
	h.Get("/api/user/balance/withdrawals", authenticator.Handle(applyMiddlewares(h.GetWithdrawals(), middlewares)))

	return h
}

func applyMiddlewares(handler http.HandlerFunc, middlewares []Middleware) http.HandlerFunc {
	for _, middleware := range middlewares {
		handler = middleware.Handle(handler)
	}

	return handler
}

func (h *Handler) getAuthUser(r *http.Request) (model.User, error) {
	login, ok := app.LoginFromContext(r.Context())
	if !ok {
		return model.User{}, errors.New("unauthorized")
	}

	user, err := h.userRepository.GetByLogin(r.Context(), login)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}
