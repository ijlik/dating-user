package service

import (
	"context"
	"database/sql"
	"github.com/ijlik/dating-user/internal/adapter/repository"
	"github.com/ijlik/dating-user/internal/business/domain"
	"github.com/ijlik/dating-user/pkg/constant"
	errpkg "github.com/ijlik/dating-user/pkg/error"
	"time"
)

func (s *service) CreatePayment(
	ctx context.Context,
	req *domain.PaymentRequest,
	UserID string,
) errpkg.ErrorService {
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

	err = s.repo.UpdatePremiumStatusProfile(ctx, &repository.UpdatePremiumStatus{
		ID:        profile.ID,
		IsPremium: true,
		IsPremiumValidUntil: sql.NullTime{
			Time:  time.Now().AddDate(0, 1, 0),
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
