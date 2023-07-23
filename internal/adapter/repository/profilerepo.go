package repository

import (
	"context"
	"database/sql"
	"time"
)

const createProfileQuery = `INSERT INTO profiles (user_id, created_at) VALUES ($1, CURRENT_TIMESTAMP) RETURNING id`

func (r *repo) CreateProfile(
	ctx context.Context,
	UserID string,
) (*Profile, error) {
	var id string
	if err := r.conn.QueryRowContext(
		ctx,
		createProfileQuery,
		UserID,
	).Scan(&id); err != nil {
		return nil, err
	}

	return &Profile{
		ID:                  id,
		UserID:              UserID,
		Name:                sql.NullString{},
		BirthDate:           sql.NullTime{},
		Gender:              sql.NullString{},
		Photos:              sql.NullString{},
		Hobby:               sql.NullString{},
		Interest:            sql.NullString{},
		Location:            sql.NullString{},
		IsPremium:           false,
		IsPremiumValidUntil: sql.NullTime{},
		DailySwapQuota:      10,
		CreatedAt:           time.Now().UTC(),
		UpdatedAt:           sql.NullTime{},
	}, nil
}

const getProfileByUserIDQuery = `SELECT id, user_id, name, birth_date, gender, photos, hobby, interest, location, is_premium, is_premium_valid_until, daily_swap_quota, created_at, updated_at FROM profiles WHERE user_id = $1 LIMIT 1`

func (r *repo) GetProfileByUserID(
	ctx context.Context,
	UserID string,
) (*Profile, error) {
	var data Profile
	err := r.conn.GetContext(
		ctx,
		&data,
		getProfileByUserIDQuery,
		UserID,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &data, nil
}

const updateBasicInfoProfileQuery = `UPDATE profiles SET name = $2, birth_date = $3, gender = $4 WHERE id = $1`

func (r *repo) UpdateBasicInfoProfile(
	ctx context.Context,
	req *UpdateProfileInfo,
) error {
	if _, err := r.conn.ExecContext(
		ctx,
		updateBasicInfoProfileQuery,
		req.RowData()...,
	); err != nil {
		return err
	}

	return nil
}

const updatePhotosProfileQuery = `UPDATE profiles SET photos = $2 WHERE id = $1`

func (r *repo) UpdatePhotosProfile(
	ctx context.Context,
	req *UpdatePhotos,
) error {
	if _, err := r.conn.ExecContext(
		ctx,
		updatePhotosProfileQuery,
		req.RowData()...,
	); err != nil {
		return err
	}

	return nil
}

const updateHobbyAndInterestProfileQuery = `UPDATE profiles SET hobby = $2, interest = $3 WHERE id = $1`

func (r *repo) UpdateHobbyAndInterestProfile(
	ctx context.Context,
	req *UpdateHobbyAndInterest,
) error {
	if _, err := r.conn.ExecContext(
		ctx,
		updateHobbyAndInterestProfileQuery,
		req.RowData()...,
	); err != nil {
		return err
	}

	return nil
}

const updateLocationProfileQuery = `UPDATE profiles SET location = $2 WHERE id = $1`

func (r *repo) UpdateLocationProfile(
	ctx context.Context,
	req *UpdateLocation,
) error {
	if _, err := r.conn.ExecContext(
		ctx,
		updateLocationProfileQuery,
		req.RowData()...,
	); err != nil {
		return err
	}

	return nil
}

const updatePremiumStatusProfileQuery = `UPDATE profiles SET is_premium = $2, is_premium_valid_until = $3, daily_swap_quota = $4 WHERE id = $1`

func (r *repo) UpdatePremiumStatusProfile(
	ctx context.Context,
	req *UpdatePremiumStatus,
) error {
	if _, err := r.conn.ExecContext(
		ctx,
		updatePremiumStatusProfileQuery,
		req.RowData()...,
	); err != nil {
		return err
	}

	return nil
}
