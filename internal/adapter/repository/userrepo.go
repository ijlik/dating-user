package repository

import (
	"context"
	"database/sql"
	"github.com/ijlik/dating-user/pkg/constant"
)

const getUserByEmailQuery = `SELECT id, phone, email, status, onboarding_steps, created_at, updated_at FROM users WHERE email = $1 LIMIT 1`

func (r *repo) GetUserByEmail(
	ctx context.Context,
	email string,
) (*User, error) {
	var data User
	err := r.conn.GetContext(
		ctx,
		&data,
		getUserByEmailQuery,
		email,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &data, nil
}

const getUserByIdQuery = `SELECT id, phone, email, status, onboarding_steps, created_at, updated_at FROM users WHERE id = $1 LIMIT 1`

func (r *repo) GetUserById(
	ctx context.Context,
	UserID string,
) (*User, error) {
	var data User
	err := r.conn.GetContext(
		ctx,
		&data,
		getUserByIdQuery,
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

const createUserQuery = `INSERT INTO users (email, onboarding_steps, created_at) VALUES ($1, $2, CURRENT_TIMESTAMP) RETURNING id`

func (r *repo) CreateUser(
	ctx context.Context,
	req *CreateUser,
) (*User, error) {
	var id string
	if err := r.conn.QueryRowContext(
		ctx,
		createUserQuery,
		req.RowData()...,
	).Scan(&id); err != nil {
		return nil, err
	}

	return &User{
		ID:              id,
		Phone:           sql.NullString{},
		Email:           req.Email,
		Status:          constant.USER_STATUS_UNVERIFIED,
		OnboardingSteps: req.OnboardingSteps,
	}, nil
}

const updateOnboardingStepsUserQuery = `UPDATE users SET onboarding_steps = $2 WHERE id = $1`

func (r *repo) UpdateOnboardingStepsUser(
	ctx context.Context,
	req *UpdateOnboardingSteps,
) error {
	if _, err := r.conn.ExecContext(
		ctx,
		updateOnboardingStepsUserQuery,
		req.RowData()...,
	); err != nil {
		return err
	}

	return nil
}

const updateStatusUserQuery = `UPDATE users SET status = $2 WHERE id = $1`

func (r *repo) UpdateStatusUser(
	ctx context.Context,
	req *UpdateStatus,
) error {
	if _, err := r.conn.ExecContext(
		ctx,
		updateStatusUserQuery,
		req.RowData()...,
	); err != nil {
		return err
	}

	return nil
}
