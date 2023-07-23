package repository

import (
	"context"
	"github.com/ijlik/dating-user/pkg/constant"
	"time"
)

type UserRepository interface {
	UserRepo
	ProfileRepo
	OneTimePasswordRepo
	SwipesRepo
	PaymentRepo
}

type UserRepo interface {
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserById(ctx context.Context, UserID string) (*User, error)
	CreateUser(ctx context.Context, req *CreateUser) (*User, error)
	UpdateOnboardingStepsUser(ctx context.Context, req *UpdateOnboardingSteps) error
	UpdateStatusUser(ctx context.Context, req *UpdateStatus) error
}

type ProfileRepo interface {
	CreateProfile(ctx context.Context, UserID string) (*Profile, error)
	GetProfileByUserID(ctx context.Context, UserID string) (*Profile, error)
	UpdateBasicInfoProfile(ctx context.Context, req *UpdateProfileInfo) error
	UpdatePhotosProfile(ctx context.Context, req *UpdatePhotos) error
	UpdateHobbyAndInterestProfile(ctx context.Context, req *UpdateHobbyAndInterest) error
	UpdateLocationProfile(ctx context.Context, req *UpdateLocation) error
	UpdatePremiumStatusProfile(ctx context.Context, req *UpdatePremiumStatus) error
}

type OneTimePasswordRepo interface {
	CreateOneTimePasswordLog(ctx context.Context, req *OneTimePasswordLog) error
	GetOneTimePasswordLogByUserAndType(ctx context.Context, UserID, otpType string) (*OneTimePasswordLog, error)
	UpdateStatusOneTimePasswordLog(ctx context.Context, status constant.OneTimeLogStatus, id string, otpLimit int) error
	GetCountOneTimePasswordByTime(ctx context.Context, UserID, otpType string, startDate, endDate time.Time) (int, error)
}

type SwipesRepo interface {
	GetProfileBySwiperId(ctx context.Context, swiperId string) ([]*Profile, error)
	GetProfileBySwiperIdWithProfileId(ctx context.Context, swiperId, profileId string) ([]*Profile, error)
	CreateSwipes(ctx context.Context, req *Swipe) error
	GetSwipesCount(ctx context.Context, swiperId string) (int, error)
}

type PaymentRepo interface {
	CreatePayment(ctx context.Context, payment *Payment) error
}
