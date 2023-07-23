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
	mailerpkg "github.com/ijlik/dating-user/pkg/mailer"
	"strings"
	"time"
)

func buildDefaultOnboardingSteps() string {
	onboardingSteps := [][]string{
		{"personal-info", "pending"},
		{"photos", "pending"},
		{"hobby-and-interest", "pending"},
		{"location", "pending"},
	}
	var steps string
	for i := 0; i < len(onboardingSteps); i++ {
		steps += fmt.Sprintf("%s:%s,", onboardingSteps[i][0], onboardingSteps[i][1])
	}
	steps = strings.TrimSuffix(steps, ",")
	return steps
}

func (s *service) sendOtpAction(
	ctx context.Context,
	user *repository.User,
) (int, error) {
	var otpNumber = s.math.RandomNumberWithLen(0)
	interval := s.config.GetInt("RESEND_OTP_INTERVAL")
	if interval == 0 {
		interval = 60
	}
	start := s.time.Now().UTC().Add(time.Duration(-1*interval) * time.Second)
	end := s.time.Now().UTC()
	countOtp, err := s.repo.GetCountOneTimePasswordByTime(ctx, user.ID, constant.OTP_TYPE_EMAIL.String(), start, end)
	if err != nil {
		return 0, errpkg.DefaultServiceError(
			errpkg.ErrInternal,
			err.Error(),
		)
	}

	if countOtp > 0 {
		return interval, nil
	}

	err = s.mailer.Send(
		mailerpkg.LOGIN,
		user.Email,
		map[string]interface{}{
			"Code": otpNumber.String(),
		},
	)

	if err != nil {
		return 0, errpkg.DefaultServiceError(
			errpkg.ErrInternal,
			err.Error(),
		)
	}

	otpLimit := s.config.GetInt("OTP_MAX_TRY_LIMIT")

	err = s.repo.CreateOneTimePasswordLog(ctx, &repository.OneTimePasswordLog{
		UserID:              user.ID,
		OneTimePasswordType: constant.OTP_TYPE_EMAIL,
		Code:                otpNumber.String(),
		Status:              constant.OTP_STATUS_UNUSED,
		CreatedAt:           s.time.Now(),
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

func (s *service) ResendOtp(
	ctx context.Context,
	req *domain.ResendOtpRequest,
) (*domain.ResendOtpResponse, errpkg.ErrorService) {
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

func (s *service) LoginOrRegister(
	ctx context.Context,
	req *domain.AuthRequest,
) (*domain.AuthResponse, errpkg.ErrorService) {
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
	expireLoginLimitConfig := time.Minute * time.Duration(s.config.GetInt("OTP_EXPIRY_TIME_IN_MINUTE"))
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
	privateKey := s.config.GetString("TOKEN_SECRET_KEY")
	if privateKey == "" {
		return nil, errpkg.DefaultServiceError(
			errpkg.ErrInternal,
			"internal server error",
		)
	}
	expireTokenLimitConfig := time.Minute * time.Duration(s.config.GetInt("LOGIN_INTERVAL_LIMIT"))
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
	// Set session in redis
	loginInterval := s.config.GetInt("LOGIN_INTERVAL_LIMIT")
	err = s.redis.Set(ctx, accessToken, user.ID, loginInterval)
	if err != nil {
		return nil, errpkg.DefaultServiceError(
			errpkg.ErrInternal,
			"failed to send redis :"+err.Error(),
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

func (s *service) Logout(
	ctx context.Context,
	UserID string,
) errpkg.ErrorService {
	err := s.redis.Del(ctx, UserID)
	if err != nil {
		return errpkg.DefaultServiceError(
			errpkg.ErrInternal,
			"",
		)
	}

	return nil
}

func (s *service) ShowProfile(
	ctx context.Context,
	UserID string,
) (*domain.Profile, errpkg.ErrorService) {
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
