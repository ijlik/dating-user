package repository

import (
	"database/sql"
	"time"
)

type Profile struct {
	ID                  string         `db:"id"`
	UserID              string         `db:"user_id"`
	Name                sql.NullString `db:"name"`
	BirthDate           sql.NullTime   `db:"birth_date"`
	Gender              sql.NullString `db:"gender"`
	Photos              sql.NullString `db:"photos"`
	Hobby               sql.NullString `db:"hobby"`
	Interest            sql.NullString `db:"interest"`
	Location            sql.NullString `db:"location"`
	IsPremium           bool           `db:"is_premium"`
	IsPremiumValidUntil sql.NullTime   `db:"is_premium_valid_until"`
	DailySwapQuota      int            `db:"daily_swap_quota"`
	CreatedAt           time.Time      `db:"created_at"`
	UpdatedAt           sql.NullTime   `db:"updated_at"`
}

type UpdateProfileInfo struct {
	ID        string       `db:"id"`
	Name      string       `db:"name"`
	BirthDate sql.NullTime `db:"birth_date"`
	Gender    string       `db:"gender"`
}

func (u *UpdateProfileInfo) RowData() []interface{} {
	var data = []interface{}{
		u.ID,
		u.Name,
		u.BirthDate,
		u.Gender,
	}
	return data
}

type UpdatePhotos struct {
	ID     string `db:"id"`
	Photos string `db:"photos"`
}

func (u *UpdatePhotos) RowData() []interface{} {
	var data = []interface{}{
		u.ID,
		u.Photos,
	}
	return data
}

type UpdateHobbyAndInterest struct {
	ID       string `db:"id"`
	Hobby    string `db:"hobby"`
	Interest string `db:"interest"`
}

func (u *UpdateHobbyAndInterest) RowData() []interface{} {
	var data = []interface{}{
		u.ID,
		u.Hobby,
		u.Interest,
	}
	return data
}

type UpdateLocation struct {
	ID       string `db:"id"`
	Location string `db:"location"`
}

func (u *UpdateLocation) RowData() []interface{} {
	var data = []interface{}{
		u.ID,
		u.Location,
	}
	return data
}

type UpdatePremiumStatus struct {
	ID                  string       `db:"id"`
	IsPremium           bool         `db:"is_premium"`
	IsPremiumValidUntil sql.NullTime `db:"is_premium_valid_until"`
	DailySwapQuota      int8         `db:"daily_swap_quota"`
}

func (u *UpdatePremiumStatus) RowData() []interface{} {
	var data = []interface{}{
		u.ID,
		u.IsPremium,
		u.IsPremiumValidUntil,
		u.DailySwapQuota,
	}
	return data
}

func (p *Profile) GetBirthDate() string {
	if !p.BirthDate.Valid {
		return ""
	}
	return p.BirthDate.Time.String()
}

func (p *Profile) GetIsPremiumValidUntil() time.Time {
	return p.IsPremiumValidUntil.Time
}

func (p *Profile) GetUpdatedAt() time.Time {
	return p.UpdatedAt.Time
}
