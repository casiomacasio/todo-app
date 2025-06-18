package domain

import (
	"time";
	"github.com/google/uuid";
)

type RefreshToken struct {
	Token     uuid.UUID `json:"token" db:"token" binding:"required"`
	UserID    int       `json:"user_id" db:"user_id" binding:"required"`
	IssuedAt  time.Time `json:"issued_at" db:"issued_at"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at" binding:"required"`
	Revoked   bool      `json:"revoked" db:"revoked"`
}