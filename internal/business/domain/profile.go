package domain

import (
	"time"
)

type Profile struct {
	ID                  string    `json:"id"`
	UserID              string    `json:"user_id"`
	Name                string    `json:"name"`
	BirthDate           string    `json:"birth_date"`
	Gender              string    `json:"gender"`
	Photos              []string  `json:"photos"`
	Hobby               []string  `json:"hobby"`
	Interest            []string  `json:"interest"`
	Location            *Location `json:"location"`
	IsPremium           bool      `json:"is_premium"`
	IsPremiumValidUntil time.Time `json:"is_premium_valid_until"`
	DailySwapQuota      int       `json:"daily_swap_quota"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
	User                *User     `json:"user"`
}
