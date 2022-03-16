package repository

import (
	"context"
	"database/sql"
	"github.com/slavkluev/gophermart/internal/app/model"
)

type UserRepository struct {
	db *sql.DB
}

func CreateUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(ctx context.Context, user model.User) error {
	sqlStatement := `INSERT INTO "user" (login, password_hash) VALUES ($1, $2)`
	_, err := r.db.ExecContext(ctx, sqlStatement, user.Login, user.PasswordHash)
	return err
}

func (r *UserRepository) GetByLogin(ctx context.Context, login string) (model.User, error) {
	var user model.User

	row := r.db.QueryRowContext(ctx, `SELECT id, login, password_hash, balance, withdrawn FROM "user" WHERE login = $1`, login)
	err := row.Scan(&user.ID, &user.Login, &user.PasswordHash, &user.Balance, &user.Withdrawn)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (r *UserRepository) DecreaseBalanceByUserID(ctx context.Context, userID uint64, amount float64) error {
	sqlStatement := `UPDATE "user" SET balance = balance - $1, withdrawn = withdrawn + $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, sqlStatement, amount, userID)
	return err
}

func (r *UserRepository) IncreaseBalanceByUserID(ctx context.Context, userID uint64, amount float64) error {
	sqlStatement := `UPDATE "user" SET balance = balance + $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, sqlStatement, amount, userID)
	return err
}
