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
	"mime/multipart"
	"testing"
	"time"
)

func TestUpdatePersonalInfo(t *testing.T) {
	// Create a mock DB connection
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbx := sqlx.NewDb(db, "postgres")
	repo := repository.NewUserRepo(dbx)
	mockOnboardingService := &OnboardingServiceMock{Mock: mocktest.Mock{}, repo: repo}
	svc := NewOnboardingService(mockOnboardingService, repo)

	// Define test data
	ctx := context.Background()
	UserID := "test_user_id"
	times, _ := time.Parse("2006-01-02 15:04:05", "1990-01-01 00:00:00")
	req := &domain.UpdatePersonalInfo{
		Name:      "John Doe",
		BirthDate: times,
		Gender:    "Female",
	}

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
		Location:            sql.NullString{String: "San Francisco", Valid: true},
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

	// Set up mock behavior for UpdateBasicInfoProfile
	profileID := "profile_id_1"
	newName := "John Doe"
	newBirthDate := times
	newGender := "Female"

	// Set up the expected query and result for UpdateBasicInfoProfile
	updateBasicInfoProfileQueryMock := "UPDATE profiles SET name = \\$2, birth_date = \\$3, gender = \\$4 WHERE id = \\$1"
	mock.ExpectExec(updateBasicInfoProfileQueryMock).
		WithArgs(profileID, newName, newBirthDate, newGender).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Set up mock behavior for GetUserById
	UserID = "test_user_id"
	expectedData := &repository.User{
		ID:              UserID,
		Phone:           sql.NullString{},
		Email:           "test@example.com",
		Status:          constant.USER_STATUS_UNVERIFIED,
		OnboardingSteps: "pending",
		CreatedAt:       time.Now(),
	}

	getUserByIdQueryMock := "SELECT id, phone, email, status, onboarding_steps, created_at, updated_at FROM users WHERE id = \\$1 LIMIT 1"
	mock.ExpectQuery(getUserByIdQueryMock).WithArgs(UserID).WillReturnRows(sqlmock.NewRows([]string{"id", "phone", "email", "status", "onboarding_steps", "created_at", "updated_at"}).
		AddRow(expectedData.ID, expectedData.Phone, expectedData.Email, expectedData.Status, expectedData.OnboardingSteps, expectedData.CreatedAt, nil))

	// Set up mock behavior for UpdateOnboardingStepsUser
	UserID = "test_user_id"
	onboardingSteps := "completed"

	updateOnboardingStepsUserQueryMock := "UPDATE users SET onboarding_steps = \\$2 WHERE id = \\$1"
	mock.ExpectExec(updateOnboardingStepsUserQueryMock).WithArgs(UserID, onboardingSteps).WillReturnResult(sqlmock.NewResult(0, 1))

	mockOnboardingService.Mock.On("UpdatePersonalInfo", ctx, req, UserID).Return(nil)
	// Call the function being tested
	err = svc.onboardingService.UpdatePersonalInfo(ctx, req, UserID)

	// Assertions
	assert.Nil(t, err, "Expected no error")
}

func TestUpdatePhotos(t *testing.T) {
	// Create a mock DB connection
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbx := sqlx.NewDb(db, "postgres")
	repo := repository.NewUserRepo(dbx)
	mockOnboardingService := &OnboardingServiceMock{Mock: mocktest.Mock{}, repo: repo}
	svc := NewOnboardingService(mockOnboardingService, repo)

	// Define test data
	ctx := context.Background()
	UserID := "test_user_id"
	req := &domain.UpdatePhotos{
		Photos: []*multipart.FileHeader{
			{
				Filename: "file1.png",
				Size:     2048,
				Header:   make(map[string][]string),
			},
			{
				Filename: "file2.png",
				Size:     2048,
				Header:   make(map[string][]string),
			},
		},
	}

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
		Location:            sql.NullString{String: "San Francisco", Valid: true},
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

	// Set up mock behavior for UpdatePhotosProfile
	profileID := "profile_id_1"
	photo1FileName := "storage/photos/profile_id_1/file1.png"
	photo2FileName := "storage/photos/profile_id_1/file2.png"

	// Set up the expected query and result for UpdatePhotosProfile
	updatePhotosProfileQueryMock := "UPDATE profiles SET photos = \\$2 WHERE id = \\$1"
	mock.ExpectExec(updatePhotosProfileQueryMock).
		WithArgs(profileID, photo1FileName+","+photo2FileName).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Set up mock behavior for GetUserById
	UserID = "test_user_id"
	expectedData := &repository.User{
		ID:              UserID,
		Phone:           sql.NullString{},
		Email:           "test@example.com",
		Status:          constant.USER_STATUS_UNVERIFIED,
		OnboardingSteps: "pending",
		CreatedAt:       time.Now(),
	}

	getUserByIdQueryMock := "SELECT id, phone, email, status, onboarding_steps, created_at, updated_at FROM users WHERE id = \\$1 LIMIT 1"
	mock.ExpectQuery(getUserByIdQueryMock).WithArgs(UserID).WillReturnRows(sqlmock.NewRows([]string{"id", "phone", "email", "status", "onboarding_steps", "created_at", "updated_at"}).
		AddRow(expectedData.ID, expectedData.Phone, expectedData.Email, expectedData.Status, expectedData.OnboardingSteps, expectedData.CreatedAt, nil))

	// Set up mock behavior for UpdateOnboardingStepsUser
	UserID = "test_user_id"
	onboardingSteps := "completed"

	updateOnboardingStepsUserQueryMock := "UPDATE users SET onboarding_steps = \\$2 WHERE id = \\$1"
	mock.ExpectExec(updateOnboardingStepsUserQueryMock).WithArgs(UserID, onboardingSteps).WillReturnResult(sqlmock.NewResult(0, 1))

	mockOnboardingService.Mock.On("UpdatePhotos", ctx, req, UserID).Return(nil)
	// Call the function being tested
	err = svc.onboardingService.UpdatePhotos(ctx, req, UserID)

	// Assertions
	assert.Nil(t, err, "Expected no error")
}

func TestUpdateHobbyAndInterest(t *testing.T) {
	// Create a mock DB connection
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbx := sqlx.NewDb(db, "postgres")
	repo := repository.NewUserRepo(dbx)
	mockOnboardingService := &OnboardingServiceMock{Mock: mocktest.Mock{}, repo: repo}
	svc := NewOnboardingService(mockOnboardingService, repo)

	// Define test data
	ctx := context.Background()
	UserID := "test_user_id"
	req := &domain.UpdateHobbyAndInterest{
		Hobby:    []string{"Swimming", "Cycling"},
		Interest: []string{"Cooking", "Reading"},
	}

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
		Location:            sql.NullString{String: "San Francisco", Valid: true},
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

	// Set up mock behavior for UpdateHobbyAndInterestProfile
	profileID := "profile_id_1"
	newHobby := "Swimming,Cycling"
	newInterest := "Cooking,Reading"

	// Set up the expected query and result for UpdateHobbyAndInterestProfile
	updateHobbyAndInterestProfileQueryMock := "UPDATE profiles SET hobby = \\$2, interest = \\$3 WHERE id = \\$1"
	mock.ExpectExec(updateHobbyAndInterestProfileQueryMock).
		WithArgs(profileID, newHobby, newInterest).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Set up mock behavior for GetUserById
	UserID = "test_user_id"
	expectedData := &repository.User{
		ID:              UserID,
		Phone:           sql.NullString{},
		Email:           "test@example.com",
		Status:          constant.USER_STATUS_UNVERIFIED,
		OnboardingSteps: "pending",
		CreatedAt:       time.Now(),
	}

	getUserByIdQueryMock := "SELECT id, phone, email, status, onboarding_steps, created_at, updated_at FROM users WHERE id = \\$1 LIMIT 1"
	mock.ExpectQuery(getUserByIdQueryMock).WithArgs(UserID).WillReturnRows(sqlmock.NewRows([]string{"id", "phone", "email", "status", "onboarding_steps", "created_at", "updated_at"}).
		AddRow(expectedData.ID, expectedData.Phone, expectedData.Email, expectedData.Status, expectedData.OnboardingSteps, expectedData.CreatedAt, nil))

	// Set up mock behavior for UpdateOnboardingStepsUser
	UserID = "test_user_id"
	onboardingSteps := "completed"

	updateOnboardingStepsUserQueryMock := "UPDATE users SET onboarding_steps = \\$2 WHERE id = \\$1"
	mock.ExpectExec(updateOnboardingStepsUserQueryMock).WithArgs(UserID, onboardingSteps).WillReturnResult(sqlmock.NewResult(0, 1))

	mockOnboardingService.Mock.On("UpdateHobbyAndInterest", ctx, req, UserID).Return(nil)
	// Call the function being tested
	err = svc.onboardingService.UpdateHobbyAndInterest(ctx, req, UserID)

	// Assertions
	assert.Nil(t, err, "Expected no error")
}

func TestUpdateLocation(t *testing.T) {
	// Create a mock DB connection
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbx := sqlx.NewDb(db, "postgres")
	repo := repository.NewUserRepo(dbx)
	mockOnboardingService := &OnboardingServiceMock{Mock: mocktest.Mock{}, repo: repo}
	svc := NewOnboardingService(mockOnboardingService, repo)

	// Define test data
	ctx := context.Background()
	UserID := "test_user_id"
	req := &domain.Location{
		Longitude: "45.1234",
		Latitude:  "-76.5678",
	}

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
		Location:            sql.NullString{String: "San Francisco", Valid: true},
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

	// Set up mock behavior for UpdateLocationProfile
	profileID := "profile_id_1"
	newLocation := "45.1234:-76.5678"

	// Set up the expected query and result for UpdateLocationProfile
	updateLocationProfileQueryMock := "UPDATE profiles SET location = \\$2 WHERE id = \\$1"
	mock.ExpectExec(updateLocationProfileQueryMock).
		WithArgs(profileID, newLocation).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Set up mock behavior for GetUserById
	UserID = "test_user_id"
	expectedData := &repository.User{
		ID:              UserID,
		Phone:           sql.NullString{},
		Email:           "test@example.com",
		Status:          constant.USER_STATUS_UNVERIFIED,
		OnboardingSteps: "pending",
		CreatedAt:       time.Now(),
	}

	getUserByIdQueryMock := "SELECT id, phone, email, status, onboarding_steps, created_at, updated_at FROM users WHERE id = \\$1 LIMIT 1"
	mock.ExpectQuery(getUserByIdQueryMock).WithArgs(UserID).WillReturnRows(sqlmock.NewRows([]string{"id", "phone", "email", "status", "onboarding_steps", "created_at", "updated_at"}).
		AddRow(expectedData.ID, expectedData.Phone, expectedData.Email, expectedData.Status, expectedData.OnboardingSteps, expectedData.CreatedAt, nil))

	// Set up mock behavior for UpdateOnboardingStepsUser
	UserID = "test_user_id"
	onboardingSteps := "completed"

	updateOnboardingStepsUserQueryMock := "UPDATE users SET onboarding_steps = \\$2 WHERE id = \\$1"
	mock.ExpectExec(updateOnboardingStepsUserQueryMock).WithArgs(UserID, onboardingSteps).WillReturnResult(sqlmock.NewResult(0, 1))

	mockOnboardingService.Mock.On("UpdateLocation", ctx, req, UserID).Return(nil)
	// Call the function being tested
	err = svc.onboardingService.UpdateLocation(ctx, req, UserID)

	// Assertions
	assert.Nil(t, err, "Expected no error")
}
