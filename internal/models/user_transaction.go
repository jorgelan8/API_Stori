package models

import (
	"time"
)

// UserTransaction represents a user transaction in the system
type UserTransaction struct {
	ID       int       `json:"id"`
	UserID   int       `json:"user_id"`
	Amount   float64   `json:"amount"`
	DateTime time.Time `json:"datetime"`
}
