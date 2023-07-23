package repository

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ijlik/dating-user/pkg/constant"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGetUserByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbx := sqlx.NewDb(db, "postgres")
	repo := NewUserRepo(dbx)

	// Set up the expected query and result
	email := "test@example.com"
	expectedData := &User{
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

	ctx := context.Background()
	data, err := repo.GetUserByEmail(ctx, email)
	assert.NoError(t, err)
	assert.Equal(t, expectedData, data)
}

func TestGetUserById(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbx := sqlx.NewDb(db, "postgres")
	repo := NewUserRepo(dbx)

	// Set up the expected query and result
	UserID := "test_user_id"
	expectedData := &User{
		ID:              UserID,
		Phone:           sql.NullString{},
		Email:           "test@example.com",
		Status:          constant.USER_STATUS_UNVERIFIED,
		OnboardingSteps: "pending",
		CreatedAt:       time.Now(),
	}

	getUserByIdQueryMock := "SELECT id, phone, email, status, onboarding_steps, created_at, updated_at FROM users WHERE id = \\$1 LIMIT 1"

	rows := sqlmock.NewRows([]string{"id", "phone", "email", "status", "onboarding_steps", "created_at", "updated_at"}).
		AddRow(expectedData.ID, expectedData.Phone, expectedData.Email, expectedData.Status, expectedData.OnboardingSteps, expectedData.CreatedAt, nil)

	mock.ExpectQuery(getUserByIdQueryMock).WithArgs(UserID).WillReturnRows(rows)

	ctx := context.Background()
	data, err := repo.GetUserById(ctx, UserID)
	assert.NoError(t, err)
	assert.Equal(t, expectedData, data)
}

func TestCreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbx := sqlx.NewDb(db, "postgres")
	repo := NewUserRepo(dbx)

	// Set up the expected query and result
	email := "test@example.com"
	onboardingSteps := "pending"
	expectedData := &User{
		ID:              "test_user_id",
		Phone:           sql.NullString{},
		Email:           email,
		Status:          constant.USER_STATUS_UNVERIFIED,
		OnboardingSteps: onboardingSteps,
		CreatedAt:       time.Now(),
	}

	createUserQueryMock := "INSERT INTO users \\(email, onboarding_steps, created_at\\) VALUES \\(\\$1, \\$2, CURRENT_TIMESTAMP\\) RETURNING id"

	rows := sqlmock.NewRows([]string{"id"}).AddRow(expectedData.ID)

	mock.ExpectQuery(createUserQueryMock).WithArgs(email, onboardingSteps).WillReturnRows(rows)

	ctx := context.Background()
	req := &CreateUser{
		Email:           email,
		OnboardingSteps: onboardingSteps,
	}
	data, err := repo.CreateUser(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, expectedData.ID, data.ID)
}

func TestUpdateOnboardingStepsUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbx := sqlx.NewDb(db, "postgres")
	repo := NewUserRepo(dbx)

	// Set up the expected query and parameters
	UserID := "test_user_id"
	onboardingSteps := "completed"

	updateOnboardingStepsUserQueryMock := "UPDATE users SET onboarding_steps = \\$2 WHERE id = \\$1"

	mock.ExpectExec(updateOnboardingStepsUserQueryMock).WithArgs(UserID, onboardingSteps).WillReturnResult(sqlmock.NewResult(0, 1))

	ctx := context.Background()
	req := &UpdateOnboardingSteps{
		ID:              UserID,
		OnboardingSteps: onboardingSteps,
	}
	err = repo.UpdateOnboardingStepsUser(ctx, req)
	assert.NoError(t, err)
}

func TestUpdateStatusUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbx := sqlx.NewDb(db, "postgres")
	repo := NewUserRepo(dbx)

	// Set up the expected query and parameters
	UserID := "test_user_id"
	status := constant.USER_STATUS_ACTIVE

	updateStatusUserQueryMock := "UPDATE users SET status = \\$2 WHERE id = \\$1"
	mock.ExpectExec(updateStatusUserQueryMock).WithArgs(UserID, status).WillReturnResult(sqlmock.NewResult(0, 1))

	ctx := context.Background()
	req := &UpdateStatus{
		ID:     UserID,
		Status: status,
	}
	err = repo.UpdateStatusUser(ctx, req)
	assert.NoError(t, err)
}
