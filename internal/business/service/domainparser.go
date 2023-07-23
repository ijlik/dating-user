package service

import (
	"database/sql"
	"fmt"
	"github.com/ijlik/dating-user/internal/adapter/repository"
	"github.com/ijlik/dating-user/internal/business/domain"
	"strings"
)

func ProfilesFeeds(data []*repository.Profile) []*domain.Profile {
	var result []*domain.Profile
	for _, item := range data {
		result = append(result, ProfileRes(item, &repository.User{}, 0))
	}
	return result
}

func ProfileRes(data *repository.Profile, user *repository.User, dailyCount int) *domain.Profile {
	var photos []string
	if data.Photos.String != "" {
		photos = strings.Split(data.Photos.String, ",")
	}

	var hobby []string
	if data.Hobby.String != "" {
		hobby = strings.Split(data.Hobby.String, ",")
	}

	var interest []string
	if data.Interest.String != "" {
		interest = strings.Split(data.Interest.String, ",")
	}
	var location *domain.Location
	if data.Location.String != "" {
		coordinate := strings.Split(data.Location.String, ":")
		location = &domain.Location{
			Longitude: coordinate[0],
			Latitude:  coordinate[1],
			Url:       fmt.Sprintf("https://www.google.com/maps?q=%s,%s", coordinate[0], coordinate[1]),
		}
	}

	return &domain.Profile{
		ID:                  data.ID,
		UserID:              data.UserID,
		Name:                data.Name.String,
		BirthDate:           data.GetBirthDate(),
		Gender:              data.Gender.String,
		Photos:              photos,
		Hobby:               hobby,
		Interest:            interest,
		Location:            location,
		IsPremium:           data.IsPremium,
		IsPremiumValidUntil: data.GetIsPremiumValidUntil(),
		DailySwapQuota:      data.DailySwapQuota - dailyCount,
		CreatedAt:           data.CreatedAt,
		UpdatedAt:           data.GetUpdatedAt(),
		User: &domain.User{
			Email:           user.Email,
			Status:          user.Status.String(),
			OnboardingSteps: user.GetOnboardingSteps(),
		},
	}
}

func UpdateProfileInfoReq(req *domain.UpdatePersonalInfo, profileId string) *repository.UpdateProfileInfo {
	return &repository.UpdateProfileInfo{
		ID:   profileId,
		Name: req.Name,
		BirthDate: sql.NullTime{
			Time:  req.BirthDate,
			Valid: true,
		},
		Gender: req.Gender,
	}
}
