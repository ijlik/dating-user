package service

import (
	"context"
	"database/sql"
	"github.com/ijlik/dating-user/internal/adapter/repository"
	"github.com/ijlik/dating-user/internal/business/domain"
	"github.com/ijlik/dating-user/pkg/constant"
	errpkg "github.com/ijlik/dating-user/pkg/error"
	"github.com/stretchr/testify/mock"
	"time"
)

type PaymentServiceMock struct {
	Mock mock.Mock
	repo repository.UserRepository
}

type PaymentService interface {
	CreatePayment(ctx context.Context, req *domain.PaymentRequest, UserID string) errpkg.ErrorService
}

func NewPaymentService(paymentService PaymentService, repo repository.UserRepository) *MockPaymentService {
	return &MockPaymentService{
		paymentService: paymentService,
		repo:           repo,
	}
}

type MockPaymentService struct {
	paymentService PaymentService
	repo           repository.UserRepository
}

func (s *PaymentServiceMock) CreatePayment(
	ctx context.Context,
	req *domain.PaymentRequest,
	UserID string,
) errpkg.ErrorService {
	_ = s.Mock.Called(req, UserID)
	if req == nil {
		return errpkg.DefaultServiceError(
			errpkg.ErrInternal,
			"empty req",
		)
	}
	if UserID == "" {
		return errpkg.DefaultServiceError(
			errpkg.ErrInternal,
			"empty req",
		)
	}
	profile, err := s.repo.GetProfileByUserID(ctx, UserID)
	if err != nil {
		return errpkg.DefaultServiceError(
			errpkg.ErrInternal,
			err.Error(),
		)
	}
	// To do payment with 3rd party
	// Code here

	// Save payment result
	err = s.repo.CreatePayment(ctx, &repository.Payment{
		UserID:        UserID,
		Amount:        req.Amount,
		Identifier:    req.Identifier,
		PaymentMethod: req.PaymentMethod,
		PaymentData:   req.PaymentData,
		Status:        constant.PAYMENT_STATUS_SUCCESS,
	})
	if err != nil {
		return errpkg.DefaultServiceError(
			errpkg.ErrInternal,
			err.Error(),
		)
	}

	times, err := time.Parse("2006-01-02 15:04:05", "2023-07-21 14:30:00")
	if err != nil {
		return errpkg.DefaultServiceError(
			errpkg.ErrInternal,
			err.Error(),
		)
	}
	err = s.repo.UpdatePremiumStatusProfile(ctx, &repository.UpdatePremiumStatus{
		ID:        profile.ID,
		IsPremium: true,
		IsPremiumValidUntil: sql.NullTime{
			Time:  times,
			Valid: true,
		},
		DailySwapQuota: -1,
	})
	if err != nil {
		return errpkg.DefaultServiceError(
			errpkg.ErrInternal,
			err.Error(),
		)
	}
	return nil
}
