package service

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ijlik/dating-user/internal/adapter/repository"
	"github.com/ijlik/dating-user/internal/business/domain"
	"github.com/ijlik/dating-user/pkg/constant"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	mocktest "github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestCreatePayment(t *testing.T) {
	// Create a mock DB connection
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbx := sqlx.NewDb(db, "postgres")
	repo := repository.NewUserRepo(dbx)
	mockPaymentService := &PaymentServiceMock{Mock: mocktest.Mock{}, repo: repo}
	svc := NewPaymentService(mockPaymentService, repo)

	// Set up input data for the CreatePayment function
	ctx := context.Background()
	UserID := "user_id_1"
	amount := 120.0
	identifier := "payment_identifier"
	paymentMethod := constant.PAYMENT_METHOD_GOOGLE_WALLET
	paymentData := "payment_data"
	status := constant.PAYMENT_STATUS_SUCCESS

	// Set up the expected profile data retrieved from the repository
	expectedProfile := &repository.Profile{
		ID:                  "profile_id_1",
		UserID:              UserID,
		Name:                sql.NullString{String: "John Smith", Valid: true},
		BirthDate:           sql.NullTime{Time: time.Now().AddDate(-25, 0, 0), Valid: true},
		Gender:              sql.NullString{String: "Male", Valid: true},
		Photos:              sql.NullString{String: "photo1.jpg", Valid: true},
		Hobby:               sql.NullString{String: "swimming", Valid: true},
		Interest:            sql.NullString{String: "cooking", Valid: true},
		Location:            sql.NullString{String: "San Francisco", Valid: true},
		IsPremium:           true,
		IsPremiumValidUntil: sql.NullTime{Time: time.Now().AddDate(1, 0, 0), Valid: true},
		DailySwapQuota:      10,
		CreatedAt:           time.Now(),
		UpdatedAt:           sql.NullTime{Time: time.Now(), Valid: true},
	}

	// Set up the input data for updating premium status
	profileID := "profile_id_1"
	isPremium := true
	times, _ := time.Parse("2006-01-02 15:04:05", "2023-07-21 14:30:00")
	premiumValidUntil := times
	dailySwapQuota := -1

	// Mock the GetProfileByUserID function to return the expected profile data
	getProfileByUserIDQueryMock := "SELECT id, user_id, name, birth_date, gender, photos, hobby, interest, location, is_premium, is_premium_valid_until, daily_swap_quota, created_at, updated_at FROM profiles WHERE user_id = \\$1 LIMIT 1"
	rows := sqlmock.NewRows([]string{"id", "user_id", "name", "birth_date", "gender", "photos", "hobby", "interest", "location", "is_premium", "is_premium_valid_until", "daily_swap_quota", "created_at", "updated_at"}).
		AddRow(expectedProfile.ID, expectedProfile.UserID, expectedProfile.Name, expectedProfile.BirthDate, expectedProfile.Gender, expectedProfile.Photos, expectedProfile.Hobby, expectedProfile.Interest, expectedProfile.Location, expectedProfile.IsPremium, expectedProfile.IsPremiumValidUntil, expectedProfile.DailySwapQuota, expectedProfile.CreatedAt, expectedProfile.UpdatedAt)
	mock.ExpectQuery(getProfileByUserIDQueryMock).WithArgs(UserID).WillReturnRows(rows)

	// Mock the CreatePayment function to return nil (success)
	createPaymentQueryMock := "INSERT INTO payments \\(user_id, amount, identifier, payment_method, payment_data, status, created_at\\) VALUES \\(\\$1, \\$2, \\$3, \\$4, \\$5, \\$6, CURRENT_TIMESTAMP\\)"
	mock.ExpectExec(createPaymentQueryMock).
		WithArgs(UserID, amount, identifier, paymentMethod, paymentData, status).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Set up the expected query and result for UpdatePremiumStatusProfile
	updatePremiumStatusProfileQueryMock := "UPDATE profiles SET is_premium = \\$2, is_premium_valid_until = \\$3, daily_swap_quota = \\$4 WHERE id = \\$1"
	mock.ExpectExec(updatePremiumStatusProfileQueryMock).
		WithArgs(profileID, isPremium, premiumValidUntil, dailySwapQuota).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Call the CreatePayment function
	req := &domain.PaymentRequest{
		Amount:        float32(amount),
		Identifier:    identifier,
		PaymentMethod: paymentMethod,
		PaymentData:   paymentData,
	}

	mockPaymentService.Mock.On("CreatePayment", req, UserID).Return(nil)
	err = svc.paymentService.CreatePayment(ctx, req, UserID)

	// Check for any errors
	assert.NoError(t, err)
	assert.Equal(t, req.Amount, float32(120))
	assert.Equal(t, req.Identifier, "payment_identifier")
	assert.Equal(t, req.PaymentMethod, constant.PAYMENT_METHOD_GOOGLE_WALLET)
	assert.Equal(t, req.PaymentData, "payment_data")
	// Check that the mock expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}
