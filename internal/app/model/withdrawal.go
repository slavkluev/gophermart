package model

import (
	"encoding/json"
	"time"
)

type Withdrawal struct {
	ID          uint64    `json:"-"`
	Order       string    `json:"order"`
	Sum         float64   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
	UserID      uint64    `json:"-"`
}

func (w Withdrawal) MarshalJSON() ([]byte, error) {
	type WithdrawalAlias Withdrawal

	aliasValue := struct {
		WithdrawalAlias
		ProcessedAt string `json:"processed_at"`
	}{
		WithdrawalAlias: WithdrawalAlias(w),
		ProcessedAt:     w.ProcessedAt.Format(time.RFC3339),
	}

	return json.Marshal(aliasValue)
}
