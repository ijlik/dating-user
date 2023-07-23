package repository

import (
	"context"
	"database/sql"
	"github.com/ijlik/dating-user/pkg/constant"
	"time"
)

const updatePreviousUnusedOtpToExpired = `UPDATE one_time_password_logs SET status = 'EXPIRED' WHERE status = 'UNUSED' AND user_id = $1;`

const createOneTimePasswordLog = `INSERT INTO one_time_password_logs (user_id, onetime_password_type, code, status, created_at) VALUES ($1, $2, $3, $4, $5)`

func (r *repo) CreateOneTimePasswordLog(
	ctx context.Context,
	req *OneTimePasswordLog,
) error {
	if _, err := r.conn.ExecContext(
		ctx,
		updatePreviousUnusedOtpToExpired,
		req.UserID,
	); err != nil {
		return err
	}

	if _, err := r.conn.ExecContext(
		ctx,
		createOneTimePasswordLog,
		req.RowData()...,
	); err != nil {
		return err
	}

	return nil
}

const getOneTimePasswordLogByUserAndTypeQuery = `SELECT id, user_id, onetime_password_type, code, status, created_at, otp_limit FROM one_time_password_logs WHERE user_id = $1 AND onetime_password_type = $2 ORDER BY created_at DESC LIMIT 1`

func (r *repo) GetOneTimePasswordLogByUserAndType(
	ctx context.Context,
	UserID, otpType string,
) (*OneTimePasswordLog, error) {
	var data OneTimePasswordLog
	err := r.conn.GetContext(
		ctx,
		&data,
		getOneTimePasswordLogByUserAndTypeQuery,
		UserID, otpType,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &data, nil

}

const updateStatusOtpLogQuery = `UPDATE one_time_password_logs SET status = $1, updated_at = $2, otp_limit = $3 WHERE id = $4`

func (r *repo) UpdateStatusOneTimePasswordLog(
	ctx context.Context,
	status constant.OneTimeLogStatus, id string, otpLimit int,
) error {
	tag, err := r.conn.ExecContext(
		ctx,
		updateStatusOtpLogQuery,
		status,
		time.Now().UTC(),
		otpLimit,
		id,
	)

	if err != nil {
		return err
	}

	return checkTagInt(tag, "UpdateStatusOneTimePasswordLog")
}

const getCountOTPbyTimeQuery = `SELECT count(*) FROM one_time_password_logs WHERE user_id = $1 AND onetime_password_type = $2 AND created_at BETWEEN $3 AND $4`

func (r *repo) GetCountOneTimePasswordByTime(
	ctx context.Context,
	UserID,
	otpType string,
	startDate, endDate time.Time,
) (int, error) {
	var count int
	err := r.conn.GetContext(
		ctx,
		&count,
		getCountOTPbyTimeQuery,
		UserID,
		otpType,
		startDate,
		endDate,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}

		return 0, err
	}

	return count, nil
}
