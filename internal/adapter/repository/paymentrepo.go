package repository

import (
	"context"
)

const createPaymentQuery = `INSERT INTO payments (user_id, amount, identifier, payment_method, payment_data, status, created_at) VALUES ($1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP)`

func (r *repo) CreatePayment(
	ctx context.Context,
	req *Payment,
) error {
	if _, err := r.conn.ExecContext(
		ctx,
		createPaymentQuery,
		req.RowData()...,
	); err != nil {
		return err
	}

	return nil
}
