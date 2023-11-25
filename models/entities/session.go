package entities

import (
	"time"
)

type Session struct {
	ID        uint64
	CreatedAt time.Time
	UpdatedAt *time.Time
	UserID    uint64
	Token     SessionToken
}
