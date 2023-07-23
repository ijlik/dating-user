package repository

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ijlik/dating-user/pkg/constant"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCreateOneTimePasswordLog(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbx := sqlx.NewDb(db, "postgres")
	repo := NewUserRepo(dbx)

	// Set up the input data for creating a one-time password log
	UserID := "user_id_1"
	otpType := constant.OTP_TYPE_EMAIL
	code := "123456"
	status := constant.OTP_STATUS_UNUSED
	createdAt := time.Now().UTC()

	// Set up the expected query and result for CreateOneTimePasswordLog
	updatePreviousUnusedOtpToExpiredMock := "UPDATE one_time_password_logs SET status = 'EXPIRED' WHERE status = 'UNUSED' AND user_id = \\$1"
	mock.ExpectExec(updatePreviousUnusedOtpToExpiredMock).
		WithArgs(UserID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	createOneTimePasswordLogMock := "INSERT INTO one_time_password_logs \\(user_id, onetime_password_type, code, status, created_at\\) VALUES \\(\\$1, \\$2, \\$3, \\$4, \\$5\\)"
	mock.ExpectExec(createOneTimePasswordLogMock).
		WithArgs(UserID, otpType, code, status, createdAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Call the CreateOneTimePasswordLog function
	ctx := context.Background()
	req := &OneTimePasswordLog{
		UserID:              UserID,
		OneTimePasswordType: otpType,
		Code:                code,
		Status:              status,
		CreatedAt:           createdAt,
	}
	err = repo.CreateOneTimePasswordLog(ctx, req)
	assert.NoError(t, err)
}

func TestGetOneTimePasswordLogByUserAndType(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbx := sqlx.NewDb(db, "postgres")
	repo := NewUserRepo(dbx)

	// Set up the input data for getting a one-time password log by user ID and type
	UserID := "user_id_1"
	otpType := constant.OTP_TYPE_EMAIL

	// Set up the expected query and result for GetOneTimePasswordLogByUserAndType
	getOneTimePasswordLogByUserAndTypeMock := "SELECT id, user_id, onetime_password_type, code, status, created_at, otp_limit FROM one_time_password_logs WHERE user_id = \\$1 AND onetime_password_type = \\$2 ORDER BY created_at DESC LIMIT 1"
	mock.ExpectQuery(getOneTimePasswordLogByUserAndTypeMock).
		WithArgs(UserID, otpType).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "onetime_password_type", "code", "status", "created_at", "otp_limit"}).
			AddRow("log_id_1", UserID, otpType, "123456", constant.OTP_STATUS_UNUSED, time.Now().UTC(), 5))

	// Call the GetOneTimePasswordLogByUserAndType function
	ctx := context.Background()
	log, err := repo.GetOneTimePasswordLogByUserAndType(ctx, UserID, otpType.String())
	assert.NoError(t, err)
	assert.NotNil(t, log)
	assert.Equal(t, "log_id_1", log.ID)
	assert.Equal(t, UserID, log.UserID)
	assert.Equal(t, otpType, log.OneTimePasswordType)
	assert.Equal(t, "123456", log.Code)
	assert.Equal(t, constant.OTP_STATUS_UNUSED, log.Status)
	assert.Equal(t, 5, log.OTPLimit)
}

func TestUpdateStatusOneTimePasswordLog(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbx := sqlx.NewDb(db, "postgres")
	repo := NewUserRepo(dbx)

	// Set up the input data for updating the status of a one-time password log
	logID := "log_id_1"
	status := constant.OTP_STATUS_USED
	otpLimit := 0 // Assuming the OTP has been used

	// Set up the expected query and result for UpdateStatusOneTimePasswordLog
	updateStatusOtpLogQueryMock := "UPDATE one_time_password_logs SET status = \\$1, updated_at = \\$2, otp_limit = \\$3 WHERE id = \\$4"
	mock.ExpectExec(updateStatusOtpLogQueryMock).
		WithArgs(status, sqlmock.AnyArg(), otpLimit, logID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Call the UpdateStatusOneTimePasswordLog function
	ctx := context.Background()
	err = repo.UpdateStatusOneTimePasswordLog(ctx, status, logID, otpLimit)
	assert.NoError(t, err)
}

func TestGetCountOneTimePasswordByTime(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbx := sqlx.NewDb(db, "postgres")
	repo := NewUserRepo(dbx)

	// Set up the input data for querying the count of one-time passwords by time
	UserID := "user_id_1"
	otpType := "OTP_TYPE"
	startDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2023, 1, 10, 0, 0, 0, 0, time.UTC)

	// Set up the expected query and result for GetCountOneTimePasswordByTime
	getCountOTPbyTimeQueryMock := "SELECT count\\(\\*\\) FROM one_time_password_logs WHERE user_id = \\$1 AND onetime_password_type = \\$2 AND created_at BETWEEN \\$3 AND \\$4"
	mock.ExpectQuery(getCountOTPbyTimeQueryMock).
		WithArgs(UserID, otpType, startDate, endDate).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(5)) // Assuming 5 one-time passwords found within the given time range

	// Call the GetCountOneTimePasswordByTime function
	ctx := context.Background()
	count, err := repo.GetCountOneTimePasswordByTime(ctx, UserID, otpType, startDate, endDate)
	assert.NoError(t, err)
	assert.Equal(t, 5, count) // Expecting 5 one-time passwords found within the given time range
}
