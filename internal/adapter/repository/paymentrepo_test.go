package repository

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ijlik/dating-user/pkg/constant"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreatePayment(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbx := sqlx.NewDb(db, "postgres")
	repo := NewUserRepo(dbx)

	// Set up the input data for creating a payment
	UserID := "user_id_1"
	amount := 100.0
	identifier := "payment_identifier"
	paymentMethod := constant.PAYMENT_METHOD_GOOGLE_WALLET
	paymentData := "credit_card_data"
	status := constant.PAYMENT_STATUS_SUCCESS

	// Set up the expected query and result for CreatePayment
	createPaymentQueryMock := "INSERT INTO payments \\(user_id, amount, identifier, payment_method, payment_data, status, created_at\\) VALUES \\(\\$1, \\$2, \\$3, \\$4, \\$5, \\$6, CURRENT_TIMESTAMP\\)"
	mock.ExpectExec(createPaymentQueryMock).
		WithArgs(UserID, amount, identifier, paymentMethod, paymentData, status).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Call the CreatePayment function
	ctx := context.Background()
	req := &Payment{
		UserID:        UserID,
		Amount:        float32(amount),
		Identifier:    identifier,
		PaymentMethod: paymentMethod,
		PaymentData:   paymentData,
		Status:        status,
	}
	err = repo.CreatePayment(ctx, req)
	assert.NoError(t, err)
}
