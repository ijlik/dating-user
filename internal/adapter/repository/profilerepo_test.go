package repository

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCreateProfile(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbx := sqlx.NewDb(db, "postgres")
	repo := NewUserRepo(dbx)

	// Set up the expected query and result for CreateProfile
	UserID := "test_user_id"
	createdProfileId := "created_profile_id"
	createProfileQueryMock := "INSERT INTO profiles \\(user_id, created_at\\) VALUES \\(\\$1, CURRENT_TIMESTAMP\\) RETURNING id"
	rows := sqlmock.NewRows([]string{"id"}).AddRow(createdProfileId)
	mock.ExpectQuery(createProfileQueryMock).WithArgs(UserID).WillReturnRows(rows)

	// Call the CreateProfile function
	ctx := context.Background()
	profile, err := repo.CreateProfile(ctx, UserID)
	assert.NoError(t, err)
	assert.NotNil(t, profile)
	assert.Equal(t, createdProfileId, profile.ID)
	assert.Equal(t, UserID, profile.UserID)
	assert.False(t, profile.IsPremium)
	assert.Equal(t, 10, profile.DailySwapQuota)
}

func TestGetProfileByUserID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbx := sqlx.NewDb(db, "postgres")
	repo := NewUserRepo(dbx)

	// Set up the expected query and result for GetProfileByUserID
	UserID := "test_user_id"
	expectedProfile := &Profile{
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

	// Call the GetProfileByUserID function
	ctx := context.Background()
	profile, err := repo.GetProfileByUserID(ctx, UserID)
	assert.NoError(t, err)
	assert.NotNil(t, profile)
	assert.Equal(t, expectedProfile, profile)
}

func TestUpdateBasicInfoProfile(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbx := sqlx.NewDb(db, "postgres")
	repo := NewUserRepo(dbx)

	// Set up the input data for updating basic info
	profileID := "profile_id_1"
	newName := "John Doe"
	newBirthDate := time.Now().AddDate(-30, 0, 0)
	newGender := "Male"

	// Set up the expected query and result for UpdateBasicInfoProfile
	updateBasicInfoProfileQueryMock := "UPDATE profiles SET name = \\$2, birth_date = \\$3, gender = \\$4 WHERE id = \\$1"
	mock.ExpectExec(updateBasicInfoProfileQueryMock).
		WithArgs(profileID, newName, newBirthDate, newGender).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Call the UpdateBasicInfoProfile function
	ctx := context.Background()
	req := &UpdateProfileInfo{
		ID:        profileID,
		Name:      newName,
		BirthDate: sql.NullTime{Time: newBirthDate, Valid: true},
		Gender:    newGender,
	}
	err = repo.UpdateBasicInfoProfile(ctx, req)
	assert.NoError(t, err)
}

func TestUpdatePhotosProfile(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbx := sqlx.NewDb(db, "postgres")
	repo := NewUserRepo(dbx)

	// Set up the input data for updating photos
	profileID := "profile_id_1"
	newPhotos := "photo1.jpg,photo2.jpg,photo3.jpg"

	// Set up the expected query and result for UpdatePhotosProfile
	updatePhotosProfileQueryMock := "UPDATE profiles SET photos = \\$2 WHERE id = \\$1"
	mock.ExpectExec(updatePhotosProfileQueryMock).
		WithArgs(profileID, newPhotos).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Call the UpdatePhotosProfile function
	ctx := context.Background()
	req := &UpdatePhotos{
		ID:     profileID,
		Photos: newPhotos,
	}
	err = repo.UpdatePhotosProfile(ctx, req)
	assert.NoError(t, err)
}

func TestUpdateHobbyAndInterestProfile(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbx := sqlx.NewDb(db, "postgres")
	repo := NewUserRepo(dbx)

	// Set up the input data for updating hobby and interest
	profileID := "profile_id_1"
	newHobby := "swimming"
	newInterest := "cooking"

	// Set up the expected query and result for UpdateHobbyAndInterestProfile
	updateHobbyAndInterestProfileQueryMock := "UPDATE profiles SET hobby = \\$2, interest = \\$3 WHERE id = \\$1"
	mock.ExpectExec(updateHobbyAndInterestProfileQueryMock).
		WithArgs(profileID, newHobby, newInterest).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Call the UpdateHobbyAndInterestProfile function
	ctx := context.Background()
	req := &UpdateHobbyAndInterest{
		ID:       profileID,
		Hobby:    newHobby,
		Interest: newInterest,
	}
	err = repo.UpdateHobbyAndInterestProfile(ctx, req)
	assert.NoError(t, err)
}

func TestUpdateLocationProfile(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbx := sqlx.NewDb(db, "postgres")
	repo := NewUserRepo(dbx)

	// Set up the input data for updating location
	profileID := "profile_id_1"
	newLocation := "New York"

	// Set up the expected query and result for UpdateLocationProfile
	updateLocationProfileQueryMock := "UPDATE profiles SET location = \\$2 WHERE id = \\$1"
	mock.ExpectExec(updateLocationProfileQueryMock).
		WithArgs(profileID, newLocation).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Call the UpdateLocationProfile function
	ctx := context.Background()
	req := &UpdateLocation{
		ID:       profileID,
		Location: newLocation,
	}
	err = repo.UpdateLocationProfile(ctx, req)
	assert.NoError(t, err)
}

func TestUpdatePremiumStatusProfile(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbx := sqlx.NewDb(db, "postgres")
	repo := NewUserRepo(dbx)

	// Set up the input data for updating premium status
	profileID := "profile_id_1"
	isPremium := true
	premiumValidUntil := time.Now().AddDate(1, 0, 0)
	dailySwapQuota := 20

	// Set up the expected query and result for UpdatePremiumStatusProfile
	updatePremiumStatusProfileQueryMock := "UPDATE profiles SET is_premium = \\$2, is_premium_valid_until = \\$3, daily_swap_quota = \\$4 WHERE id = \\$1"
	mock.ExpectExec(updatePremiumStatusProfileQueryMock).
		WithArgs(profileID, isPremium, premiumValidUntil, dailySwapQuota).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Call the UpdatePremiumStatusProfile function
	ctx := context.Background()
	req := &UpdatePremiumStatus{
		ID:                  profileID,
		IsPremium:           isPremium,
		IsPremiumValidUntil: sql.NullTime{Time: premiumValidUntil, Valid: true},
		DailySwapQuota:      int8(dailySwapQuota),
	}
	err = repo.UpdatePremiumStatusProfile(ctx, req)
	assert.NoError(t, err)
}
