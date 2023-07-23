package service

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ijlik/dating-user/internal/adapter/repository"
	"github.com/ijlik/dating-user/internal/business/domain"
	"github.com/ijlik/dating-user/pkg/constant"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	mocktest "github.com/stretchr/testify/mock"
	"strings"
	"testing"
	"time"
)

func TestShowFeeds(t *testing.T) {
	// Create a mock DB connection
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbx := sqlx.NewDb(db, "postgres")
	repo := repository.NewUserRepo(dbx)
	mockFeedService := &FeedServiceMock{Mock: mocktest.Mock{}, repo: repo}
	svc := NewFeedService(mockFeedService, repo)

	// Define test data
	ctx := context.Background()
	// Set up mock behavior for GetProfileBySwiperId
	swiperID := "test_swiper_id"
	profileID := ""
	randomProfile1 := &repository.Profile{
		ID:                  "profile_id_1",
		UserID:              "user_id_1",
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

	getRandomProfileQueryMock := "SELECT id, user_id, name, birth_date, gender, photos, hobby, interest, location, is_premium, is_premium_valid_until, daily_swap_quota, created_at, updated_at FROM profiles WHERE name <> '' AND birth_date < CURRENT_TIMESTAMP AND gender <> '' AND photos <> '' AND hobby <> '' AND interest <> '' AND location <> '' AND id <> \\$1 AND id NOT IN \\(SELECT swiped_id FROM swipes WHERE swiper_id = \\$1 AND DATE\\(created_at\\) = CURRENT_DATE\\) ORDER BY RANDOM\\(\\) LIMIT 1"

	rows := sqlmock.NewRows([]string{"id", "user_id", "name", "birth_date", "gender", "photos", "hobby", "interest", "location", "is_premium", "is_premium_valid_until", "daily_swap_quota", "created_at", "updated_at"}).
		AddRow(randomProfile1.ID, "user_id_1", randomProfile1.Name, randomProfile1.BirthDate, randomProfile1.Gender, randomProfile1.Photos, randomProfile1.Hobby, randomProfile1.Interest, randomProfile1.Location, randomProfile1.IsPremium, randomProfile1.IsPremiumValidUntil, randomProfile1.DailySwapQuota, randomProfile1.CreatedAt, randomProfile1.UpdatedAt)
	mock.ExpectQuery(getRandomProfileQueryMock).WithArgs(swiperID).WillReturnRows(rows)

	// Set up the expected query and result for getProfileWithoutId
	getProfileWithoutIdQueryMock := "SELECT id, user_id, name, birth_date, gender, photos, hobby, interest, location, is_premium, is_premium_valid_until, daily_swap_quota, created_at, updated_at FROM profiles WHERE name <> '' AND birth_date < CURRENT_TIMESTAMP AND gender <> '' AND photos <> '' AND hobby <> '' AND interest <> '' AND location <> '' AND id <> \\$1 AND id <> \\$2 AND id NOT IN \\(SELECT swiped_id FROM swipes WHERE swiper_id = \\$1 AND DATE\\(created_at\\) = CURRENT_DATE\\) LIMIT 1"
	mock.ExpectQuery(getProfileWithoutIdQueryMock).WithArgs(swiperID, randomProfile1.ID).WillReturnRows(rows)

	deleteSwipesShowOnlyQueryMock := "DELETE FROM swipes WHERE swiper_id = \\$1 AND swiped_id = \\$2 AND DATE\\(created_at\\) = CURRENT_DATE"
	mock.ExpectExec(deleteSwipesShowOnlyQueryMock).WithArgs(swiperID, randomProfile1.ID).WillReturnResult(sqlmock.NewResult(1, 1))

	createSwipesQueryMock := "INSERT INTO swipes \\(swiper_id, swiped_id, is_like, created_at\\) VALUES \\(\\$1, \\$2, \\$3, CURRENT_TIMESTAMP\\)"
	mock.ExpectExec(createSwipesQueryMock).WithArgs(swiperID, randomProfile1.ID, nil).WillReturnResult(sqlmock.NewResult(1, 1))

	mockFeedService.Mock.On("ShowFeeds", ctx, swiperID, profileID).Return([]*domain.Profile{
		ProfileRes(randomProfile1, &repository.User{
			Email:           "test@email.com",
			Status:          constant.USER_STATUS_ACTIVE,
			OnboardingSteps: "personal-info:pending,photos:pending,hobby-and-interest:pending,location:pending",
		}, 0),
	})
	// Call the function being tested
	profiles, err := svc.feedService.ShowFeeds(ctx, swiperID, profileID)

	// Assertions
	assert.Nil(t, err, "Expected no error")
	assert.NotNil(t, profiles, "Expected non-nil profiles")
	assert.Len(t, profiles, 1, "Expected 1 profile in the result")

	// Check profile data
	profile := profiles[0]
	assert.Equal(t, randomProfile1.ID, profile.ID, "Profile ID mismatch")
	assert.Equal(t, randomProfile1.UserID, profile.UserID, "User ID mismatch")
	assert.Equal(t, randomProfile1.Name.String, profile.Name, "Name mismatch")
	assert.Equal(t, randomProfile1.BirthDate.Time.String(), profile.BirthDate, "BirthDate mismatch")
	assert.Equal(t, randomProfile1.Gender.String, profile.Gender, "Gender mismatch")
	assert.Equal(t, strings.Split(randomProfile1.Photos.String, ","), profile.Photos, "Photos mismatch")
	assert.Equal(t, strings.Split(randomProfile1.Hobby.String, ","), profile.Hobby, "Hobby mismatch")
	assert.Equal(t, strings.Split(randomProfile1.Interest.String, ","), profile.Interest, "Interest mismatch")
	coordinate := strings.Split(randomProfile1.Location.String, ":")
	location := &domain.Location{
		Longitude: coordinate[0],
		Latitude:  coordinate[1],
		Url:       fmt.Sprintf("https://www.google.com/maps?q=%s,%s", coordinate[0], coordinate[1]),
	}
	assert.Equal(t, location, profile.Location, "Location mismatch")
	assert.Equal(t, randomProfile1.IsPremium, profile.IsPremium, "IsPremium mismatch")
	assert.Equal(t, randomProfile1.IsPremiumValidUntil.Time, profile.IsPremiumValidUntil, "IsPremiumValidUntil mismatch")
	assert.Equal(t, randomProfile1.DailySwapQuota, profile.DailySwapQuota, "DailySwapQuota mismatch")
	assert.Equal(t, randomProfile1.CreatedAt, profile.CreatedAt, "CreatedAt mismatch")
	assert.Equal(t, randomProfile1.UpdatedAt.Time, profile.UpdatedAt, "UpdatedAt mismatch")
}

func TestSwipes(t *testing.T) {
	// Create a mock DB connection
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbx := sqlx.NewDb(db, "postgres")
	repo := repository.NewUserRepo(dbx)
	mockFeedService := &FeedServiceMock{Mock: mocktest.Mock{}, repo: repo}
	svc := NewFeedService(mockFeedService, repo)

	// Define test data
	ctx := context.Background()
	UserID := "test_user_id"
	swiperID := "swiper_id_1"
	swipedID := "swiped_id_1"
	isLike := true

	req := &domain.SwipeRequest{
		SwiperId: swiperID,
		SwipedId: swipedID,
		IsLike:   isLike,
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

	deleteSwipesShowOnlyQueryMock := "DELETE FROM swipes WHERE swiper_id = \\$1 AND swiped_id = \\$2 AND DATE\\(created_at\\) = CURRENT_DATE"
	mock.ExpectExec(deleteSwipesShowOnlyQueryMock).WithArgs(swiperID, swipedID).WillReturnResult(sqlmock.NewResult(1, 1))

	createSwipesQueryMock := "INSERT INTO swipes \\(swiper_id, swiped_id, is_like, created_at\\) VALUES \\(\\$1, \\$2, \\$3, CURRENT_TIMESTAMP\\)"
	mock.ExpectExec(createSwipesQueryMock).WithArgs(swiperID, swipedID, isLike).WillReturnResult(sqlmock.NewResult(1, 1))

	// Set up mock behavior for GetSwipesCount
	dailyQuota := 5
	getSwipesCountQueryMock := "SELECT COUNT\\(\\*\\) FROM swipes WHERE swiper_id = \\$1 AND created_at >= current_date\\(\\)"
	mock.ExpectQuery(getSwipesCountQueryMock).WithArgs(swiperID).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(dailyQuota))

	// Set up mock behavior for UpdatePremiumStatusProfile
	updatePremiumStatusProfileQueryMock := "UPDATE profiles SET is_premium = \\$2, is_premium_valid_until = \\$3, daily_swap_quota = \\$4 WHERE id = \\$1"
	mock.ExpectExec(updatePremiumStatusProfileQueryMock).
		WithArgs(expectedProfile.ID, false, nil, expectedProfile.DailySwapQuota).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mockFeedService.Mock.On("Swipes", ctx, req, UserID).Return(nil)
	// Call the function being tested
	err = svc.feedService.Swipes(ctx, req, UserID)

	// Assertions
	assert.Nil(t, err, "Expected no error")
}
