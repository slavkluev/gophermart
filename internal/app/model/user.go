package model

type User struct {
	ID           uint64  `json:"-"`
	Login        string  `json:"-"`
	PasswordHash string  `json:"-"`
	Balance      float64 `json:"current"`
	Withdrawn    float64 `json:"withdrawn"`
}
