package service

import (
	"context"
	"database/sql"
	"github.com/ijlik/dating-user/internal/adapter/repository"
	"github.com/ijlik/dating-user/internal/business/domain"
	errpkg "github.com/ijlik/dating-user/pkg/error"
	"github.com/stretchr/testify/mock"
	"time"
)

type FeedServiceMock struct {
	Mock mock.Mock
	repo repository.UserRepository
}

type FeedService interface {
	ShowFeeds(ctx context.Context, swiperId, profileId string) ([]*domain.Profile, errpkg.ErrorService)
	Swipes(ctx context.Context, req *domain.SwipeRequest, UserID string) errpkg.ErrorService
}

func NewFeedService(feedService FeedService, repo repository.UserRepository) *MockFeedService {
	return &MockFeedService{
		feedService: feedService,
		repo:        repo,
	}
}

type MockFeedService struct {
	feedService FeedService
	repo        repository.UserRepository
}

func (s *FeedServiceMock) ShowFeeds(
	ctx context.Context,
	swiperId,
	profileId string,
) ([]*domain.Profile, errpkg.ErrorService) {
	arguments := s.Mock.Called(ctx, swiperId, profileId)
	if arguments.Get(0) == nil {
		return nil, errpkg.DefaultServiceError(
			errpkg.ErrInternal,
			"empty req",
		)
	}
	//response := arguments.Get(0).([]*domain.Profile)
	var profiles []*domain.Profile
	data, err := s.repo.GetProfileBySwiperId(ctx, swiperId)
	if err != nil {
		return nil, errpkg.DefaultServiceError(
			errpkg.ErrInternal,
			err.Error(),
		)
	}
	profiles = ProfilesFeeds(data)

	return profiles, nil
}

func (s *FeedServiceMock) Swipes(
	ctx context.Context,
	req *domain.SwipeRequest,
	UserID string,
) errpkg.ErrorService {
	_ = s.Mock.Called(ctx, req, UserID)
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
	if !profile.IsPremium {
		// Check daily quota
		dailyQuota, err := s.repo.GetSwipesCount(ctx, req.SwiperId)
		if err != nil {
			return errpkg.DefaultServiceError(
				errpkg.ErrInternal,
				err.Error(),
			)
		}
		if dailyQuota >= profile.DailySwapQuota {
			return errpkg.DefaultServiceError(
				errpkg.ErrUnauthorize,
				"daily swipes quota exceed",
			)
		}
	} else {
		if profile.GetIsPremiumValidUntil().UTC().Unix() < time.Now().Unix() {
			// Update Is Premium Status to false
			err = s.repo.UpdatePremiumStatusProfile(ctx, &repository.UpdatePremiumStatus{
				ID:                  profile.ID,
				IsPremium:           false,
				IsPremiumValidUntil: sql.NullTime{},
				DailySwapQuota:      10,
			})
			if err != nil {
				return errpkg.DefaultServiceError(
					errpkg.ErrInternal,
					err.Error(),
				)
			}

			return errpkg.DefaultServiceError(
				errpkg.ErrUnauthorize,
				"premium membership expired",
			)
		}
	}

	// Create Swipes
	err = s.repo.CreateSwipes(ctx, &repository.Swipe{
		SwiperId: req.SwiperId,
		SwipedId: req.SwipedId,
		IsLike: sql.NullBool{
			Bool:  req.IsLike,
			Valid: true,
		},
	})

	if err != nil {
		return errpkg.DefaultServiceError(
			errpkg.ErrInternal,
			err.Error(),
		)
	}

	return nil
}
