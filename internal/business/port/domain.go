package port

import (
	"context"
	"github.com/ijlik/dating-user/internal/business/domain"
	errpkg "github.com/ijlik/dating-user/pkg/error"
)

type UserDomainService interface {
	ResendOtp(ctx context.Context, req *domain.ResendOtpRequest) (*domain.ResendOtpResponse, errpkg.ErrorService)
	LoginOrRegister(ctx context.Context, req *domain.AuthRequest) (*domain.AuthResponse, errpkg.ErrorService)
	Logout(ctx context.Context, UserID string) errpkg.ErrorService
	ShowProfile(ctx context.Context, UserID string) (*domain.Profile, errpkg.ErrorService)

	UpdatePersonalInfo(ctx context.Context, req *domain.UpdatePersonalInfo, UserID string) errpkg.ErrorService
	UpdatePhotos(ctx context.Context, req *domain.UpdatePhotos, UserID string) errpkg.ErrorService
	UpdateHobbyAndInterest(ctx context.Context, req *domain.UpdateHobbyAndInterest, UserID string) errpkg.ErrorService
	UpdateLocation(ctx context.Context, req *domain.Location, UserID string) errpkg.ErrorService

	ShowFeeds(ctx context.Context, UserID, profileId string) ([]*domain.Profile, errpkg.ErrorService)
	Swipes(ctx context.Context, req *domain.SwipeRequest, UserID string) errpkg.ErrorService

	CreatePayment(ctx context.Context, req *domain.PaymentRequest, UserID string) errpkg.ErrorService
}
