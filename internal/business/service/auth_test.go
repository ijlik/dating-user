package service

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ijlik/dating-user/internal/adapter/repository"
	"github.com/ijlik/dating-user/internal/business/domain"
	"github.com/ijlik/dating-user/pkg/auth"
	"github.com/ijlik/dating-user/pkg/constant"
	ctxsdk "github.com/ijlik/dating-user/pkg/context"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	mocktest "github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestResendOtp(t *testing.T) {
	// Create a mock DB connection
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbx := sqlx.NewDb(db, "postgres")
	repo := repository.NewUserRepo(dbx)
	mockAuthService := &AuthServiceMock{Mock: mocktest.Mock{}, repo: repo}
	svc := NewAuthService(mockAuthService, repo)

	// Define test data
	ctx := context.Background()
	email := "test@example.com"

	// Set up mock behavior for GetUserByEmail
	expectedData := &repository.User{
		ID:              "test_user_id",
		Phone:           sql.NullString{},
		Email:           email,
		Status:          constant.USER_STATUS_UNVERIFIED,
		OnboardingSteps: "pending",
		CreatedAt:       time.Now(),
	}

	getUserByEmailQueryMock := "SELECT id, phone, email, status, onboarding_steps, created_at, updated_at FROM users WHERE email = \\$1 LIMIT 1"

	rows := sqlmock.NewRows([]string{"id", "phone", "email", "status", "onboarding_steps", "created_at", "updated_at"}).
		AddRow(expectedData.ID, expectedData.Phone, expectedData.Email, expectedData.Status, expectedData.OnboardingSteps, expectedData.CreatedAt, nil)

	mock.ExpectQuery(getUserByEmailQueryMock).WithArgs(email).WillReturnRows(rows)

	// Set up mock behavior for CreateOneTimePasswordLog
	interval := 60
	UserID := expectedData.ID
	otpType := constant.OTP_TYPE_EMAIL
	startDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2023, 1, 10, 0, 0, 0, 0, time.UTC)
	code := "123456"
	status := constant.OTP_STATUS_UNUSED
	createdAt := time.Now().UTC()

	getCountOTPbyTimeQueryMock := "SELECT count\\(\\*\\) FROM one_time_password_logs WHERE user_id = \\$1 AND onetime_password_type = \\$2 AND created_at BETWEEN \\$3 AND \\$4"
	mock.ExpectQuery(getCountOTPbyTimeQueryMock).
		WithArgs(UserID, otpType, startDate, endDate).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(5))

	updatePreviousUnusedOtpToExpiredMock := "UPDATE one_time_password_logs SET status = 'EXPIRED' WHERE status = 'UNUSED' AND user_id = \\$1"
	mock.ExpectExec(updatePreviousUnusedOtpToExpiredMock).
		WithArgs(UserID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	createOneTimePasswordLogMock := "INSERT INTO one_time_password_logs \\(user_id, onetime_password_type, code, status, created_at\\) VALUES \\(\\$1, \\$2, \\$3, \\$4, \\$5\\)"
	mock.ExpectExec(createOneTimePasswordLogMock).
		WithArgs(UserID, otpType, code, status, createdAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	req := &domain.ResendOtpRequest{
		Email: email,
	}
	mockAuthService.Mock.On("ResendOtp", ctx, req).Return(&domain.ResendOtpResponse{
		Message:           fmt.Sprintf("OTP already sent to your email: %s.", req.Email),
		ResendOTPInterval: interval,
	})
	// Call the function being tested
	res, err := svc.authService.ResendOtp(ctx, req)

	// Assertions
	assert.Nil(t, err, "Expected no error")
	assert.NotNil(t, res, "Expected non-nil response")
	assert.Equal(t, fmt.Sprintf("OTP already sent to your email: %s.", req.Email), res.Message, "Response message mismatch")
	assert.Equal(t, interval, res.ResendOTPInterval, "ResendOTPInterval mismatch")
}

func TestLoginOrRegister(t *testing.T) {
	// Create a mock DB connection
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbx := sqlx.NewDb(db, "postgres")
	repo := repository.NewUserRepo(dbx)
	mockAuthService := &AuthServiceMock{Mock: mocktest.Mock{}, repo: repo}
	svc := NewAuthService(mockAuthService, repo)

	// Define test data
	ctx := context.Background()
	email := "test@example.com"
	otp := "123456"
	profileID := "profile_id_1"
	otpType := constant.OTP_TYPE_EMAIL
	status := constant.USER_STATUS_ACTIVE
	logID := "log_id_1"
	otpStatus := constant.OTP_STATUS_USED
	otpLimit := 4

	// Set up mock behavior for GetUserByEmail
	expectedUser := &repository.User{
		ID:              "user_id_1",
		Phone:           sql.NullString{},
		Email:           email,
		Status:          constant.USER_STATUS_UNVERIFIED,
		OnboardingSteps: buildDefaultOnboardingSteps(),
		CreatedAt:       time.Now(),
	}

	getUserByEmailQueryMock := "SELECT id, phone, email, status, onboarding_steps, created_at, updated_at FROM users WHERE email = \\$1 LIMIT 1"
	rows := sqlmock.NewRows([]string{"id", "phone", "email", "status", "onboarding_steps", "created_at", "updated_at"}).
		AddRow(expectedUser.ID, expectedUser.Phone, expectedUser.Email, expectedUser.Status, expectedUser.OnboardingSteps, expectedUser.CreatedAt, nil)
	mock.ExpectQuery(getUserByEmailQueryMock).WithArgs(email).WillReturnRows(rows)

	getOneTimePasswordLogByUserAndTypeMock := "SELECT id, user_id, onetime_password_type, code, status, created_at, otp_limit FROM one_time_password_logs WHERE user_id = \\$1 AND onetime_password_type = \\$2 ORDER BY created_at DESC LIMIT 1"
	mock.ExpectQuery(getOneTimePasswordLogByUserAndTypeMock).
		WithArgs(expectedUser.ID, otpType).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "onetime_password_type", "code", "status", "created_at", "otp_limit"}).
			AddRow("log_id_1", expectedUser.ID, otpType, "123456", constant.OTP_STATUS_UNUSED, time.Now().UTC(), 5))

	updateStatusUserQueryMock := "UPDATE users SET status = \\$2 WHERE id = \\$1"
	mock.ExpectExec(updateStatusUserQueryMock).WithArgs(expectedUser.ID, status).WillReturnResult(sqlmock.NewResult(0, 1))

	createProfileQueryMock := "INSERT INTO profiles \\(user_id, created_at\\) VALUES \\(\\$1, CURRENT_TIMESTAMP\\) RETURNING id"
	mock.ExpectQuery(createProfileQueryMock).WithArgs(expectedUser.ID).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(profileID))

	updateStatusOtpLogQueryMock := "UPDATE one_time_password_logs SET status = \\$1, updated_at = \\$2, otp_limit = \\$3 WHERE id = \\$4"
	mock.ExpectExec(updateStatusOtpLogQueryMock).
		WithArgs(otpStatus, sqlmock.AnyArg(), otpLimit, logID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Set up mock behavior for GetOneTimePasswordLogByUserAndType
	expectedOtp := &repository.OneTimePasswordLog{
		ID:                  "otp_id_1",
		UserID:              expectedUser.ID,
		OneTimePasswordType: constant.OTP_TYPE_EMAIL,
		Code:                otp,
		Status:              constant.OTP_STATUS_UNUSED,
		CreatedAt:           time.Now().Add(-2 * time.Minute), // OTP created 2 minutes ago
		OTPLimit:            3,
	}

	getOneTimePasswordLogByUserAndTypeQueryMock := "SELECT id, user_id, one_time_password_type, code, status, created_at, otp_limit FROM one_time_password_log WHERE user_id = \\$1 AND one_time_password_type = \\$2 LIMIT 1"
	rows = sqlmock.NewRows([]string{"id", "user_id", "one_time_password_type", "code", "status", "created_at", "otp_limit"}).
		AddRow(expectedOtp.ID, expectedOtp.UserID, expectedOtp.OneTimePasswordType, expectedOtp.Code, expectedOtp.Status, expectedOtp.CreatedAt, expectedOtp.OTPLimit)
	mock.ExpectQuery(getOneTimePasswordLogByUserAndTypeQueryMock).
		WithArgs(expectedUser.ID, constant.OTP_TYPE_EMAIL.String()).
		WillReturnRows(rows)

	//Set up mock behavior for UpdateStatusOneTimePasswordLog
	updateStatusOneTimePasswordLogQueryMock := "UPDATE one_time_password_log SET status = \\$1, otp_limit = \\$2 WHERE id = \\$3"
	mock.ExpectExec(updateStatusOneTimePasswordLogQueryMock).
		WithArgs(constant.OTP_STATUS_USED.String(), expectedOtp.OTPLimit-1, expectedOtp.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Set up mock behavior for CreateToken
	privateKey := "LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlDV3dJQkFBS0JnUUNBb1dyNW9zTGJVS0VGa1pGQlNHNGp0UXZuUm1FOEZJdjc2N01ISHlxcU9SKytqSzcrCnBmdTFJMERZNnlXTU1oUys0TXpiSGlXa1JLczVCbkdnTjEyOXpaYTVNUU1zdGpYWkIxQU9NVXI0TENGdWFrNTQKZW9haG1qdThQRjhkTzFNOElKTXVhcHNCZUhUd2lOVFZPRmdxQW9UTDBUeUhXL2h4WUVlanBzbENKd0lEQVFBQgpBb0dBZmdRbHRrMGpVeE1KdlZmZ0V6SHZYU0lJZUZwMTloTTNGT1hUclgxMkllLzJ6b29yQXFVQUZIUm1HbDA4Cm1yMlJuM0xDbjBSSW9rYjM2OVVKU21vVGRmS2N5eUloL3RFNGMySk5Dc0Y2MklHeUdydmlFRDdtVGR5TWR4T2UKZm5SSS9BMDVuMlR6dlJSWXE0K0ZubVhIdWdjazhiNW1vRkZmcGVJRVArWm9BdmtDUVFES2NQY0RONDhoUkZJdApEM2p4RklqRk1EQlhGTGlnQ0hSWis3NVdiSUJSUVM0TitZYkxnM2lkVGd3SEE0bEp0T2xoTDc4b2ZFZG44R1FDCk1NdEJ1RWV0QWtFQW9xbGRzR2FJRXVjbVMzblUxMDlzVlhLc2x0TXFYbE1NUU5XMHF5QTI3eEtmc1lFb3BGS1kKNTRGZkJKWUN4ZDRvRTg0LzZ1amFWcmNxa1QycFNvSDdvd0pBSVpTQmROZ25kdFkxWjJJVXByRElTeVZHTDN1eApjR0pXb29KK3ZTazhVNzRqSEpCU2lybWhMVDdBQWYzVkxSUEVUcW16NU14UXIrNFJPTWZOUDNhSTlRSkFRbHhZCmRhd08zYTloNXk3b0Q0TStqa20vY2JUcXR4cW9pQmJubzF6OExHTHJ0YTRjMTVVKy9rdkFhUTJPU2cxTlNtODkKa21lM0UrT2NRUzduenhiaWd3SkFDczdEZ2NVaEY3YkVJOFJhYWF2Q1RqSEh1eENpb09kUEEvbWxGQXJEdEZ6ZwpoTG9vam1OTUlLOWRCVkxKczloN0hWRmdTZnNRc3lZTFN5MmUwcE5aenc9PQotLS0tLUVORCBSU0EgUFJJVkFURSBLRVktLS0tLQ=="
	expectedExpireTokenLimit := time.Minute * time.Duration(180)
	expectedMapData := map[ctxsdk.ContextMetadata]string{
		ctxsdk.USER_ID:    expectedUser.ID,
		ctxsdk.EMAIL:      expectedUser.Email,
		ctxsdk.PROFILE_ID: profileID,
	}
	expectedAccessToken, err := auth.CreateToken(expectedExpireTokenLimit, &expectedMapData, privateKey)
	assert.NoError(t, err)

	// Call the function being tested
	req := &domain.AuthRequest{
		Email: email,
		Otp:   otp,
	}
	mockAuthService.Mock.On("LoginOrRegister", ctx, req).Return(&domain.AuthResponse{AccessToken: expectedAccessToken})
	res, err := svc.authService.LoginOrRegister(ctx, req)

	// Assertions
	assert.NoError(t, err, "Expected no error")
	assert.NotNil(t, res, "Expected non-nil response")

	publicKey := "LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUlHZk1BMEdDU3FHU0liM0RRRUJBUVVBQTRHTkFEQ0JpUUtCZ1FDQW9XcjVvc0xiVUtFRmtaRkJTRzRqdFF2bgpSbUU4Rkl2NzY3TUhIeXFxT1IrK2pLNytwZnUxSTBEWTZ5V01NaFMrNE16YkhpV2tSS3M1Qm5HZ04xMjl6WmE1Ck1RTXN0alhaQjFBT01VcjRMQ0Z1YWs1NGVvYWhtanU4UEY4ZE8xTThJSk11YXBzQmVIVHdpTlRWT0ZncUFvVEwKMFR5SFcvaHhZRWVqcHNsQ0p3SURBUUFCCi0tLS0tRU5EIFBVQkxJQyBLRVktLS0tLQ=="
	expectedValidatedToken, _ := auth.ValidateToken(expectedAccessToken, publicKey)
	actualValidatedToken, _ := auth.ValidateToken(res.AccessToken, publicKey)
	assert.Equal(t, expectedValidatedToken, actualValidatedToken, "Access token mismatch")
	assert.Equal(t, expectedAccessToken, res.AccessToken, "Access token mismatch")
}

func TestShowProfile(t *testing.T) {
	// Create a mock DB connection
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbx := sqlx.NewDb(db, "postgres")
	repo := repository.NewUserRepo(dbx)
	mockAuthService := &AuthServiceMock{Mock: mocktest.Mock{}, repo: repo}
	svc := NewAuthService(mockAuthService, repo)

	// Define test data
	ctx := context.Background()
	UserID := "user_id_1"

	// Set up mock behavior for GetProfileByUserID
	expectedProfile := &repository.Profile{
		ID:                  "profile_id_1",
		UserID:              UserID,
		Name:                sql.NullString{String: "John Smith", Valid: true},
		BirthDate:           sql.NullTime{Time: time.Now().AddDate(-25, 0, 0), Valid: true},
		Gender:              sql.NullString{String: "Male", Valid: true},
		Photos:              sql.NullString{String: "photo1.jpg", Valid: true},
		Hobby:               sql.NullString{String: "swimming", Valid: true},
		Interest:            sql.NullString{String: "cooking", Valid: true},
		Location:            sql.NullString{String: "45.1234:-76.5678", Valid: true},
		IsPremium:           true,
		IsPremiumValidUntil: sql.NullTime{Time: time.Now().AddDate(1, 0, 0), Valid: true},
		DailySwapQuota:      10,
		CreatedAt:           time.Now(),
		UpdatedAt:           sql.NullTime{Time: time.Now(), Valid: true},
	}

	getProfileByUserIDQueryMock := "SELECT id, user_id, name, birth_date, gender, photos, hobby, interest, location, is_premium, is_premium_valid_until, daily_swap_quota, created_at, updated_at FROM profiles WHERE user_id = \\$1 LIMIT 1"
	rows := sqlmock.NewRows([]string{"id", "user_id", "name", "birth_date", "gender", "photos", "hobby", "interest", "location", "is_premium", "is_premium_valid_until", "daily_swap_quota", "created_at", "updated_at"}).
		AddRow(expectedProfile.ID, expectedProfile.UserID, expectedProfile.Name, expectedProfile.BirthDate, expectedProfile.Gender, expectedProfile.Photos, expectedProfile.Hobby, expectedProfile.Interest, expectedProfile.Location, expectedProfile.IsPremium, expectedProfile.IsPremiumValidUntil, expectedProfile.DailySwapQuota, expectedProfile.CreatedAt, expectedProfile.UpdatedAt)
	mock.ExpectQuery(getProfileByUserIDQueryMock).WithArgs(UserID).WillReturnRows(rows)

	// Set up mock behavior for GetUserById
	expectedUser := &repository.User{
		ID:     UserID,
		Phone:  sql.NullString{},
		Email:  "test@example.com",
		Status: constant.USER_STATUS_ACTIVE,
	}

	getUserByIdQueryMock := "SELECT id, phone, email, status, onboarding_steps, created_at, updated_at FROM users WHERE id = \\$1 LIMIT 1"
	rows = sqlmock.NewRows([]string{"id", "phone", "email", "status", "onboarding_steps", "created_at", "updated_at"}).
		AddRow(expectedUser.ID, expectedUser.Phone, expectedUser.Email, expectedUser.Status, expectedUser.OnboardingSteps, expectedUser.CreatedAt, nil)
	mock.ExpectQuery(getUserByIdQueryMock).WithArgs(UserID).WillReturnRows(rows)

	count := 5
	getSwipesCountAttributeQueryMock := "SELECT count\\(\\*\\) FROM swipes WHERE swiper_id = \\$1 AND DATE\\(created_at\\) = CURRENT_DATE"
	mock.ExpectQuery(getSwipesCountAttributeQueryMock).WithArgs(expectedProfile.ID).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(count))

	// Set up mock behavior for GetSwipesCount
	dailyCount := 5
	expectedDailyCountQueryMock := "SELECT COUNT\\(id\\) FROM swipes WHERE swiper_id = \\$1"
	rows = sqlmock.NewRows([]string{"count"}).
		AddRow(dailyCount)
	mock.ExpectQuery(expectedDailyCountQueryMock).WithArgs(expectedProfile.ID).WillReturnRows(rows)

	// Call the function being tested
	mockAuthService.Mock.On("ShowProfile", ctx, UserID).Return(ProfileRes(expectedProfile, expectedUser, dailyCount))
	res, err := svc.authService.ShowProfile(ctx, UserID)

	// Assertions
	assert.NoError(t, err, "Expected no error")
	assert.NotNil(t, res, "Expected non-nil response")
	assert.Equal(t, expectedProfile.ID, res.ID, "Profile ID mismatch")
	assert.Equal(t, expectedProfile.UserID, res.UserID, "User ID mismatch")
	assert.Equal(t, expectedProfile.IsPremium, res.IsPremium, "IsPremium mismatch")
}
