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

func TestGetProfileBySwiperId(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbx := sqlx.NewDb(db, "postgres")
	repo := NewUserRepo(dbx)

	// Set up the expected query and result for getRandomProfile
	swiperId := "test_swiper_id"
	randomProfile1 := &Profile{
		ID:                  "profile_id_1",
		UserID:              "user_id_1",
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
	randomProfile2 := &Profile{
		ID:                  "profile_id_2",
		UserID:              "user_id_2",
		Name:                sql.NullString{String: "John Doe", Valid: true},
		BirthDate:           sql.NullTime{Time: time.Now().AddDate(-25, 0, 0), Valid: true},
		Gender:              sql.NullString{String: "Male", Valid: true},
		Photos:              sql.NullString{String: "photo2.jpg", Valid: true},
		Hobby:               sql.NullString{String: "slot", Valid: true},
		Interest:            sql.NullString{String: "money", Valid: true},
		Location:            sql.NullString{String: "Los Angles", Valid: true},
		IsPremium:           true,
		IsPremiumValidUntil: sql.NullTime{Time: time.Now().AddDate(1, 0, 0), Valid: true},
		DailySwapQuota:      10,
		CreatedAt:           time.Now(),
		UpdatedAt:           sql.NullTime{Time: time.Now(), Valid: true},
	}

	getRandomProfileQueryMock := "SELECT id, user_id, name, birth_date, gender, photos, hobby, interest, location, is_premium, is_premium_valid_until, daily_swap_quota, created_at, updated_at FROM profiles WHERE name <> '' AND birth_date < CURRENT_TIMESTAMP AND gender <> '' AND photos <> '' AND hobby <> '' AND interest <> '' AND location <> '' AND id <> \\$1 AND id NOT IN \\(SELECT swiped_id FROM swipes WHERE swiper_id = \\$1 AND DATE\\(created_at\\) = CURRENT_DATE\\) ORDER BY RANDOM\\(\\) LIMIT 1"

	rows := sqlmock.NewRows([]string{"id", "user_id", "name", "birth_date", "gender", "photos", "hobby", "interest", "location", "is_premium", "is_premium_valid_until", "daily_swap_quota", "created_at", "updated_at"}).
		AddRow(randomProfile1.ID, "user_id_1", randomProfile1.Name, randomProfile1.BirthDate, randomProfile1.Gender, randomProfile1.Photos, randomProfile1.Hobby, randomProfile1.Interest, randomProfile1.Location, randomProfile1.IsPremium, randomProfile1.IsPremiumValidUntil, randomProfile1.DailySwapQuota, randomProfile1.CreatedAt, randomProfile1.UpdatedAt).
		AddRow(randomProfile2.ID, "user_id_2", randomProfile2.Name, randomProfile2.BirthDate, randomProfile2.Gender, randomProfile2.Photos, randomProfile2.Hobby, randomProfile2.Interest, randomProfile2.Location, randomProfile2.IsPremium, randomProfile2.IsPremiumValidUntil, randomProfile2.DailySwapQuota, randomProfile2.CreatedAt, randomProfile2.UpdatedAt)
	mock.ExpectQuery(getRandomProfileQueryMock).WithArgs(swiperId).WillReturnRows(rows)

	// Set up the expected query and result for getProfileWithoutId
	getProfileWithoutIdQueryMock := "SELECT id, user_id, name, birth_date, gender, photos, hobby, interest, location, is_premium, is_premium_valid_until, daily_swap_quota, created_at, updated_at FROM profiles WHERE name <> '' AND birth_date < CURRENT_TIMESTAMP AND gender <> '' AND photos <> '' AND hobby <> '' AND interest <> '' AND location <> '' AND id <> \\$1 AND id <> \\$2 AND id NOT IN \\(SELECT swiped_id FROM swipes WHERE swiper_id = \\$1 AND DATE\\(created_at\\) = CURRENT_DATE\\) LIMIT 1"
	mock.ExpectQuery(getProfileWithoutIdQueryMock).WithArgs(swiperId, randomProfile1.ID).WillReturnRows(rows)

	deleteSwipesShowOnlyQueryMock := "DELETE FROM swipes WHERE swiper_id = \\$1 AND swiped_id = \\$2 AND DATE\\(created_at\\) = CURRENT_DATE"
	mock.ExpectExec(deleteSwipesShowOnlyQueryMock).WithArgs(swiperId, randomProfile1.ID).WillReturnResult(sqlmock.NewResult(1, 1))

	createSwipesQueryMock := "INSERT INTO swipes \\(swiper_id, swiped_id, is_like, created_at\\) VALUES \\(\\$1, \\$2, \\$3, CURRENT_TIMESTAMP\\)"
	mock.ExpectExec(createSwipesQueryMock).WithArgs(swiperId, randomProfile1.ID, nil).WillReturnResult(sqlmock.NewResult(1, 1))

	ctx := context.Background()
	profiles, err := repo.GetProfileBySwiperId(ctx, swiperId)
	assert.NoError(t, err)
	assert.Len(t, profiles, 2)
	assert.Equal(t, randomProfile1, profiles[0])
	assert.Equal(t, randomProfile2, profiles[1])
}

func TestGetProfileBySwiperIdWithProfileId(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbx := sqlx.NewDb(db, "postgres")
	repo := NewUserRepo(dbx)

	// Set up the expected query and result for getProfileById
	swiperId := "test_swiper_id"
	profileId := "test_profile_id"
	currentProfile := &Profile{
		ID:                  profileId,
		UserID:              "user_id_1",
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

	getProfileByIdQueryMock := "SELECT id, user_id, name, birth_date, gender, photos, hobby, interest, location, is_premium, is_premium_valid_until, daily_swap_quota, created_at, updated_at FROM profiles WHERE name <> '' AND birth_date < CURRENT_TIMESTAMP AND gender <> '' AND photos <> '' AND hobby <> '' AND interest <> '' AND location <> '' AND id <> \\$1 AND id = \\$2 AND id NOT IN \\(SELECT swiped_id FROM swipes WHERE swiper_id = \\$1 AND DATE\\(created_at\\) = CURRENT_DATE\\) LIMIT 1"
	rows := sqlmock.NewRows([]string{"id", "user_id", "name", "birth_date", "gender", "photos", "hobby", "interest", "location", "is_premium", "is_premium_valid_until", "daily_swap_quota", "created_at", "updated_at"}).
		AddRow(currentProfile.ID, "user_id_1", currentProfile.Name, currentProfile.BirthDate, currentProfile.Gender, currentProfile.Photos, currentProfile.Hobby, currentProfile.Interest, currentProfile.Location, currentProfile.IsPremium, currentProfile.IsPremiumValidUntil, currentProfile.DailySwapQuota, currentProfile.CreatedAt, currentProfile.UpdatedAt)
	mock.ExpectQuery(getProfileByIdQueryMock).WithArgs(swiperId, profileId).WillReturnRows(rows)

	// Set up the expected query and result for getProfileWithoutId
	randomProfile := &Profile{
		ID:                  "profile_id_2",
		UserID:              "user_id_2",
		Name:                sql.NullString{String: "John Doe", Valid: true},
		BirthDate:           sql.NullTime{Time: time.Now().AddDate(-25, 0, 0), Valid: true},
		Gender:              sql.NullString{String: "Male", Valid: true},
		Photos:              sql.NullString{String: "photo2.jpg", Valid: true},
		Hobby:               sql.NullString{String: "slot", Valid: true},
		Interest:            sql.NullString{String: "money", Valid: true},
		Location:            sql.NullString{String: "Los Angles", Valid: true},
		IsPremium:           true,
		IsPremiumValidUntil: sql.NullTime{Time: time.Now().AddDate(1, 0, 0), Valid: true},
		DailySwapQuota:      10,
		CreatedAt:           time.Now(),
		UpdatedAt:           sql.NullTime{Time: time.Now(), Valid: true},
	}

	getProfileWithoutIdQueryMock := "SELECT id, user_id, name, birth_date, gender, photos, hobby, interest, location, is_premium, is_premium_valid_until, daily_swap_quota, created_at, updated_at FROM profiles WHERE name <> '' AND birth_date < CURRENT_TIMESTAMP AND gender <> '' AND photos <> '' AND hobby <> '' AND interest <> '' AND location <> '' AND id <> \\$1 AND id <> \\$2 AND id NOT IN \\(SELECT swiped_id FROM swipes WHERE swiper_id = \\$1 AND DATE\\(created_at\\) = CURRENT_DATE\\) LIMIT 1"
	rows = sqlmock.NewRows([]string{"id", "user_id", "name", "birth_date", "gender", "photos", "hobby", "interest", "location", "is_premium", "is_premium_valid_until", "daily_swap_quota", "created_at", "updated_at"}).
		AddRow(randomProfile.ID, "user_id_2", randomProfile.Name, randomProfile.BirthDate, randomProfile.Gender, randomProfile.Photos, randomProfile.Hobby, randomProfile.Interest, randomProfile.Location, randomProfile.IsPremium, randomProfile.IsPremiumValidUntil, randomProfile.DailySwapQuota, randomProfile.CreatedAt, randomProfile.UpdatedAt)
	mock.ExpectQuery(getProfileWithoutIdQueryMock).WithArgs(swiperId, currentProfile.ID).WillReturnRows(rows)

	deleteSwipesShowOnlyQueryMock := "DELETE FROM swipes WHERE swiper_id = \\$1 AND swiped_id = \\$2 AND DATE\\(created_at\\) = CURRENT_DATE"
	mock.ExpectExec(deleteSwipesShowOnlyQueryMock).WithArgs(swiperId, currentProfile.ID).WillReturnResult(sqlmock.NewResult(1, 1))

	createSwipesQueryMock := "INSERT INTO swipes \\(swiper_id, swiped_id, is_like, created_at\\) VALUES \\(\\$1, \\$2, \\$3, CURRENT_TIMESTAMP\\)"
	mock.ExpectExec(createSwipesQueryMock).WithArgs(swiperId, currentProfile.ID, nil).WillReturnResult(sqlmock.NewResult(1, 1))

	ctx := context.Background()
	profiles, err := repo.GetProfileBySwiperIdWithProfileId(ctx, swiperId, profileId)
	assert.NoError(t, err)
	assert.Len(t, profiles, 2)
	assert.Equal(t, currentProfile, profiles[0])
	assert.Equal(t, randomProfile, profiles[1])
}

func TestCreateSwipes(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbx := sqlx.NewDb(db, "postgres")
	repo := NewUserRepo(dbx)

	// Set up the expected query and result for checkIfHasSwipe
	swiperId := "test_swiper_id"
	swipedId := "test_swiped_id"
	checkIfHasSwipeQueryMock := "SELECT count\\(\\*\\) FROM swipes WHERE swiper_id = \\$1 AND swiped_id = \\$2 AND is_like IN \\(true, false\\) AND DATE\\(created_at\\) = CURRENT_DATE"
	rows := sqlmock.NewRows([]string{"count"}).AddRow(0)
	mock.ExpectQuery(checkIfHasSwipeQueryMock).WithArgs(swiperId, swipedId).WillReturnRows(rows)

	// Set up the expected query and result for deleteSwipesShowOnly
	deleteSwipesShowOnlyQueryMock := "DELETE FROM swipes WHERE swiper_id = \\$1 AND swiped_id = \\$2 AND DATE\\(created_at\\) = CURRENT_DATE"
	mock.ExpectExec(deleteSwipesShowOnlyQueryMock).WithArgs(swiperId, swipedId).WillReturnResult(sqlmock.NewResult(1, 1))

	// Set up the expected query and result for createSwipes
	createSwipesQueryMock := "INSERT INTO swipes \\(swiper_id, swiped_id, is_like, created_at\\) VALUES \\(\\$1, \\$2, \\$3, CURRENT_TIMESTAMP\\)"
	mock.ExpectExec(createSwipesQueryMock).WithArgs(swiperId, swipedId, nil).WillReturnResult(sqlmock.NewResult(1, 1))

	// Call the CreateSwipes function
	ctx := context.Background()
	err = repo.CreateSwipes(ctx, &Swipe{
		SwiperId: swiperId,
		SwipedId: swipedId,
		IsLike:   sql.NullBool{},
	})
	assert.NoError(t, err)
}

func TestGetSwipesCount(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbx := sqlx.NewDb(db, "postgres")
	repo := NewUserRepo(dbx)

	// Set up the expected query and result for getSwipesCount
	swiperId := "test_swiper_id"
	count := 5
	getSwipesCountAttributeQueryMock := "SELECT count\\(\\*\\) FROM swipes WHERE swiper_id = \\$1 AND DATE\\(created_at\\) = CURRENT_DATE"
	rows := sqlmock.NewRows([]string{"count"}).AddRow(count)
	mock.ExpectQuery(getSwipesCountAttributeQueryMock).WithArgs(swiperId).WillReturnRows(rows)

	// Call the GetSwipesCount function
	ctx := context.Background()
	result, err := repo.GetSwipesCount(ctx, swiperId)
	assert.NoError(t, err)
	assert.Equal(t, count, result)
}
