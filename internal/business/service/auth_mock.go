package service

import (
	"context"
	"fmt"
	"github.com/ijlik/dating-user/internal/adapter/repository"
	"github.com/ijlik/dating-user/internal/business/domain"
	"github.com/ijlik/dating-user/pkg/auth"
	"github.com/ijlik/dating-user/pkg/constant"
	ctxsdk "github.com/ijlik/dating-user/pkg/context"
	errpkg "github.com/ijlik/dating-user/pkg/error"
	commonmath "github.com/ijlik/dating-user/pkg/math"
	"github.com/stretchr/testify/mock"
	"time"
)

type AuthServiceMock struct {
	Mock mock.Mock
	repo repository.UserRepository
}

type AuthService interface {
	ResendOtp(ctx context.Context, req *domain.ResendOtpRequest) (*domain.ResendOtpResponse, errpkg.ErrorService)
	LoginOrRegister(ctx context.Context, req *domain.AuthRequest) (*domain.AuthResponse, errpkg.ErrorService)
	ShowProfile(ctx context.Context, UserID string) (*domain.Profile, errpkg.ErrorService)
}

func NewAuthService(authService AuthService, repo repository.UserRepository) *MockAuthService {
	return &MockAuthService{
		authService: authService,
		repo:        repo,
	}
}

type MockAuthService struct {
	authService AuthService
	repo        repository.UserRepository
}

func (s *AuthServiceMock) ShowProfile(
	ctx context.Context,
	UserID string,
) (*domain.Profile, errpkg.ErrorService) {
	arguments := s.Mock.Called(ctx, UserID)
	if arguments.Get(0) == nil {
		return nil, errpkg.DefaultServiceError(
			errpkg.ErrInternal,
			"empty req",
		)
	}

	data, err := s.repo.GetProfileByUserID(ctx, UserID)
	if err != nil {
		return nil, errpkg.DefaultServiceError(
			errpkg.ErrInternal,
			err.Error(),
		)
	}
	user, err := s.repo.GetUserById(ctx, UserID)
	if err != nil {
		return nil, errpkg.DefaultServiceError(
			errpkg.ErrInternal,
			err.Error(),
		)
	}

	dailyCount, err := s.repo.GetSwipesCount(ctx, data.ID)
	if err != nil {
		return nil, errpkg.DefaultServiceError(
			errpkg.ErrInternal,
			err.Error(),
		)
	}

	return ProfileRes(data, user, dailyCount), nil
}

func (s *AuthServiceMock) ResendOtp(
	ctx context.Context,
	req *domain.ResendOtpRequest,
) (*domain.ResendOtpResponse, errpkg.ErrorService) {
	arguments := s.Mock.Called(ctx, req)
	if arguments.Get(0) == nil {
		return nil, errpkg.DefaultServiceError(
			errpkg.ErrInternal,
			"empty req",
		)
	}

	var interval int
	var user *repository.User
	var err error
	user, err = s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, errpkg.DefaultServiceError(
			errpkg.ErrInternal,
			err.Error(),
		)
	}

	if user == nil {
		user, err = s.repo.CreateUser(ctx, &repository.CreateUser{
			Email:           req.Email,
			OnboardingSteps: buildDefaultOnboardingSteps(),
		})
		if err != nil {
			return nil, errpkg.DefaultServiceError(
				errpkg.ErrInternal,
				err.Error(),
			)
		}

	}

	interval, err = s.sendOtpAction(ctx, user)
	if err != nil {
		return nil, errpkg.DefaultServiceError(
			errpkg.ErrInternal,
			err.Error(),
		)
	}

	return &domain.ResendOtpResponse{
		Message:           fmt.Sprintf("OTP already sent to your email: %s.", req.Email),
		ResendOTPInterval: interval,
	}, nil
}

func (s *AuthServiceMock) LoginOrRegister(
	ctx context.Context,
	req *domain.AuthRequest,
) (*domain.AuthResponse, errpkg.ErrorService) {
	arguments := s.Mock.Called(ctx, req)
	if arguments.Get(0) == nil {
		return nil, errpkg.DefaultServiceError(
			errpkg.ErrInternal,
			"empty req",
		)
	}

	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, errpkg.DefaultServiceError(
			errpkg.ErrInternal,
			err.Error(),
		)
	}

	otp, err := s.repo.GetOneTimePasswordLogByUserAndType(ctx, user.ID, constant.OTP_TYPE_EMAIL.String())
	if err != nil {
		return nil, errpkg.DefaultServiceError(
			errpkg.ErrInternal,
			err.Error(),
		)
	}
	// Validate expired OTP by max attempt try
	if otp.OTPLimit <= 0 {
		err := s.repo.UpdateStatusOneTimePasswordLog(ctx, constant.OTP_STATUS_EXPIRED, otp.ID, 0)
		if err != nil {
			return nil, errpkg.DefaultServiceError(
				errpkg.ErrInternal,
				err.Error(),
			)
		}
		return nil, errpkg.DefaultServiceError(
			errpkg.ErrUnauthorize,
			"To many attempt. Otp Expired",
		)
	}

	// validate same otp code
	if req.Otp != otp.Code {
		err := s.repo.UpdateStatusOneTimePasswordLog(ctx, constant.OTP_STATUS_UNUSED, otp.ID, otp.OTPLimit-1)
		if err != nil {
			return nil, errpkg.DefaultServiceError(
				errpkg.ErrInternal,
				err.Error(),
			)
		}
		return nil, errpkg.DefaultServiceError(
			errpkg.ErrUnauthorize,
			"Invalid OTP",
		)
	}

	// Validate otp status
	if otp.Status != constant.OTP_STATUS_UNUSED {
		return nil, errpkg.DefaultServiceError(
			errpkg.ErrUnauthorize,
			"OTP Expired",
		)
	}

	// Validate expired OTP bu time
	expireLoginLimitConfig := time.Minute * time.Duration(5)
	var expireLoginLimit = otp.CreatedAt.Add(expireLoginLimitConfig)

	if time.Now().UTC().Unix() > expireLoginLimit.Unix() {
		err = s.repo.UpdateStatusOneTimePasswordLog(ctx, constant.OTP_STATUS_EXPIRED, otp.ID, 0)
		if err != nil {
			return nil, errpkg.DefaultServiceError(
				errpkg.ErrInternal,
				err.Error(),
			)
		}
		return nil, errpkg.DefaultServiceError(
			errpkg.ErrBadRequest,
			"otp is expired",
		)
	}

	// OTP VALID
	// Find or Create Profile
	var profile *repository.Profile
	if user.Status == constant.USER_STATUS_UNVERIFIED {
		// Update User Status
		err = s.repo.UpdateStatusUser(ctx, &repository.UpdateStatus{
			ID:     user.ID,
			Status: constant.USER_STATUS_ACTIVE,
		})
		if err != nil {
			return nil, errpkg.DefaultServiceError(
				errpkg.ErrInternal,
				err.Error(),
			)
		}
		// Create Profile
		profile, err = s.repo.CreateProfile(ctx, user.ID)
		if err != nil {
			return nil, errpkg.DefaultServiceError(
				errpkg.ErrInternal,
				err.Error(),
			)
		}
	} else {
		profile, err = s.repo.GetProfileByUserID(ctx, user.ID)
		if err != nil {
			return nil, errpkg.DefaultServiceError(
				errpkg.ErrInternal,
				err.Error(),
			)
		}
	}
	// Create token
	privateKey := "LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlDV3dJQkFBS0JnUUNBb1dyNW9zTGJVS0VGa1pGQlNHNGp0UXZuUm1FOEZJdjc2N01ISHlxcU9SKytqSzcrCnBmdTFJMERZNnlXTU1oUys0TXpiSGlXa1JLczVCbkdnTjEyOXpaYTVNUU1zdGpYWkIxQU9NVXI0TENGdWFrNTQKZW9haG1qdThQRjhkTzFNOElKTXVhcHNCZUhUd2lOVFZPRmdxQW9UTDBUeUhXL2h4WUVlanBzbENKd0lEQVFBQgpBb0dBZmdRbHRrMGpVeE1KdlZmZ0V6SHZYU0lJZUZwMTloTTNGT1hUclgxMkllLzJ6b29yQXFVQUZIUm1HbDA4Cm1yMlJuM0xDbjBSSW9rYjM2OVVKU21vVGRmS2N5eUloL3RFNGMySk5Dc0Y2MklHeUdydmlFRDdtVGR5TWR4T2UKZm5SSS9BMDVuMlR6dlJSWXE0K0ZubVhIdWdjazhiNW1vRkZmcGVJRVArWm9BdmtDUVFES2NQY0RONDhoUkZJdApEM2p4RklqRk1EQlhGTGlnQ0hSWis3NVdiSUJSUVM0TitZYkxnM2lkVGd3SEE0bEp0T2xoTDc4b2ZFZG44R1FDCk1NdEJ1RWV0QWtFQW9xbGRzR2FJRXVjbVMzblUxMDlzVlhLc2x0TXFYbE1NUU5XMHF5QTI3eEtmc1lFb3BGS1kKNTRGZkJKWUN4ZDRvRTg0LzZ1amFWcmNxa1QycFNvSDdvd0pBSVpTQmROZ25kdFkxWjJJVXByRElTeVZHTDN1eApjR0pXb29KK3ZTazhVNzRqSEpCU2lybWhMVDdBQWYzVkxSUEVUcW16NU14UXIrNFJPTWZOUDNhSTlRSkFRbHhZCmRhd08zYTloNXk3b0Q0TStqa20vY2JUcXR4cW9pQmJubzF6OExHTHJ0YTRjMTVVKy9rdkFhUTJPU2cxTlNtODkKa21lM0UrT2NRUzduenhiaWd3SkFDczdEZ2NVaEY3YkVJOFJhYWF2Q1RqSEh1eENpb09kUEEvbWxGQXJEdEZ6ZwpoTG9vam1OTUlLOWRCVkxKczloN0hWRmdTZnNRc3lZTFN5MmUwcE5aenc9PQotLS0tLUVORCBSU0EgUFJJVkFURSBLRVktLS0tLQ=="

	expireTokenLimitConfig := time.Minute * time.Duration(180)
	var mapData = map[ctxsdk.ContextMetadata]string{
		ctxsdk.USER_ID:    user.ID,
		ctxsdk.EMAIL:      user.Email,
		ctxsdk.PROFILE_ID: profile.ID,
	}
	accessToken, err := auth.CreateToken(expireTokenLimitConfig, &mapData, privateKey)
	if err != nil {
		return nil, errpkg.DefaultServiceError(
			errpkg.ErrInternal,
			"internal server error create token",
		)
	}

	// Update Otp Status
	err = s.repo.UpdateStatusOneTimePasswordLog(ctx, constant.OTP_STATUS_USED, otp.ID, otp.OTPLimit-1)
	if err != nil {
		return nil, errpkg.DefaultServiceError(
			errpkg.ErrInternal,
			err.Error(),
		)
	}

	return &domain.AuthResponse{
		AccessToken: accessToken,
	}, nil
}

func (s *AuthServiceMock) sendOtpAction(
	ctx context.Context,
	user *repository.User,
) (int, error) {
	math := commonmath.NewMath()
	var otpNumber = math.RandomNumberWithLen(0)

	var interval = 60

	startDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2023, 1, 10, 0, 0, 0, 0, time.UTC)
	countOtp, err := s.repo.GetCountOneTimePasswordByTime(ctx, user.ID, constant.OTP_TYPE_EMAIL.String(), startDate, endDate)
	if err != nil {
		return 0, errpkg.DefaultServiceError(
			errpkg.ErrInternal,
			err.Error(),
		)
	}

	if countOtp > 0 {
		return interval, nil
	}

	otpLimit := 3

	err = s.repo.CreateOneTimePasswordLog(ctx, &repository.OneTimePasswordLog{
		UserID:              user.ID,
		OneTimePasswordType: constant.OTP_TYPE_EMAIL,
		Code:                otpNumber.String(),
		Status:              constant.OTP_STATUS_UNUSED,
		CreatedAt:           time.Now(),
		OTPLimit:            otpLimit,
	})

	if err != nil {
		return 0, errpkg.DefaultServiceError(
			errpkg.ErrInternal,
			err.Error(),
		)
	}

	return interval, nil
}
