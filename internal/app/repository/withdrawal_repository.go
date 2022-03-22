package repository

import (
	"context"
	"database/sql"
	"errors"
	"sort"

	"github.com/slavkluev/gophermart/internal/app/model"
)

var (
	ErrInsufficientBalance = errors.New("insufficient balance")
)

type WithdrawalRepository struct {
	db *sql.DB
}

func CreateWithdrawalRepository(db *sql.DB) *WithdrawalRepository {
	return &WithdrawalRepository{
		db: db,
	}
}

func (r *WithdrawalRepository) Create(ctx context.Context, withdrawal model.Withdrawal) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var balance float64
	row := tx.QueryRowContext(ctx, `SELECT balance FROM "user" WHERE id = $1`, withdrawal.UserID)
	err = row.Scan(&balance)
	if err != nil {
		return err
	}

	if balance < withdrawal.Sum {
		return ErrInsufficientBalance
	}

	createWithdrawalStatement := `INSERT INTO withdrawal ("order", sum, processed_at, user_id) VALUES ($1, $2, $3, $4)`
	_, err = tx.ExecContext(ctx, createWithdrawalStatement, withdrawal.Order, withdrawal.Sum, withdrawal.ProcessedAt, withdrawal.UserID)
	if err != nil {
		return err
	}

	updateBalanceStatement := `UPDATE "user" SET balance = balance - $1, withdrawn = withdrawn + $1 WHERE id = $2`
	_, err = tx.ExecContext(ctx, updateBalanceStatement, withdrawal.Sum, withdrawal.UserID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *WithdrawalRepository) GetByUserID(ctx context.Context, userID uint64) ([]model.Withdrawal, error) {
	var withdrawals []model.Withdrawal

	rows, err := r.db.QueryContext(ctx, `SELECT id, "order", sum, processed_at, user_id FROM withdrawal WHERE user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var withdrawal model.Withdrawal
		err := rows.Scan(&withdrawal.ID, &withdrawal.Order, &withdrawal.Sum, &withdrawal.ProcessedAt, &withdrawal.UserID)
		if err != nil {
			return nil, err
		}

		withdrawals = append(withdrawals, withdrawal)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	if len(withdrawals) == 0 {
		return nil, sql.ErrNoRows
	}

	sort.Slice(withdrawals, func(i, j int) bool {
		return withdrawals[i].ProcessedAt.Before(withdrawals[j].ProcessedAt)
	})

	return withdrawals, nil
}
