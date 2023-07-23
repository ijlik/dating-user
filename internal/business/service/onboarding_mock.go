package service

import (
	"context"
	"fmt"
	"github.com/ijlik/dating-user/internal/adapter/repository"
	"github.com/ijlik/dating-user/internal/business/domain"
	errpkg "github.com/ijlik/dating-user/pkg/error"
	"github.com/stretchr/testify/mock"
	"path/filepath"
	"strings"
)

type OnboardingServiceMock struct {
	Mock mock.Mock
	repo repository.UserRepository
}

type OnboardingService interface {
	UpdatePersonalInfo(ctx context.Context, req *domain.UpdatePersonalInfo, UserID string) errpkg.ErrorService
	UpdatePhotos(ctx context.Context, req *domain.UpdatePhotos, UserID string) errpkg.ErrorService
	UpdateHobbyAndInterest(ctx context.Context, req *domain.UpdateHobbyAndInterest, UserID string) errpkg.ErrorService
	UpdateLocation(ctx context.Context, req *domain.Location, UserID string) errpkg.ErrorService
}

func NewOnboardingService(onboardingService OnboardingService, repo repository.UserRepository) *MockOnboardingService {
	return &MockOnboardingService{
		onboardingService: onboardingService,
		repo:              repo,
	}
}

type MockOnboardingService struct {
	onboardingService OnboardingService
	repo              repository.UserRepository
}

func (s *OnboardingServiceMock) UpdatePersonalInfo(
	ctx context.Context,
	req *domain.UpdatePersonalInfo,
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
	err = s.repo.UpdateBasicInfoProfile(ctx, UpdateProfileInfoReq(req, profile.ID))
	if err != nil {
		return errpkg.DefaultServiceError(
			errpkg.ErrInternal,
			err.Error(),
		)
	}

	user, err := s.repo.GetUserById(ctx, UserID)
	if err != nil {
		return errpkg.DefaultServiceError(
			errpkg.ErrInternal,
			err.Error(),
		)
	}
	onboardingSteps := user.GetOnboardingSteps()
	currentStepStatus := ""
	for _, step := range onboardingSteps {
		if step.Step == "personal-info" {
			currentStepStatus = step.Status
			break
		}
	}

	if currentStepStatus == "pending" {
		// Update Onboarding Step
		var steps string
		for _, step := range onboardingSteps {
			status := step.Status
			if step.Step == "personal-info" {
				status = "done"
			}
			steps += fmt.Sprintf("%s:%s,", step.Step, status)
		}
		steps = strings.TrimSuffix(steps, ",")
		err = s.repo.UpdateOnboardingStepsUser(ctx, &repository.UpdateOnboardingSteps{
			ID:              UserID,
			OnboardingSteps: steps,
		})
		if err != nil {
			return errpkg.DefaultServiceError(
				errpkg.ErrInternal,
				err.Error(),
			)
		}
	}

	return nil
}

func (s *OnboardingServiceMock) UpdatePhotos(
	ctx context.Context,
	req *domain.UpdatePhotos,
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

	uploadDir := fmt.Sprintf("storage/photos/%s", profile.ID)
	var allPhotos string
	for _, photo := range req.Photos {
		// Generate a unique filename for each photo
		fileName := filepath.Join(uploadDir, photo.Filename)
		allPhotos += fmt.Sprintf("%s,", fileName)
	}

	err = s.repo.UpdatePhotosProfile(ctx, &repository.UpdatePhotos{
		ID:     profile.ID,
		Photos: strings.TrimSuffix(allPhotos, ","),
	})

	if err != nil {
		return errpkg.DefaultServiceError(
			errpkg.ErrInternal,
			err.Error(),
		)
	}

	user, err := s.repo.GetUserById(ctx, UserID)
	if err != nil {
		return errpkg.DefaultServiceError(
			errpkg.ErrInternal,
			err.Error(),
		)
	}
	onboardingSteps := user.GetOnboardingSteps()
	currentStepStatus := ""
	for _, step := range onboardingSteps {
		if step.Step == "photos" {
			currentStepStatus = step.Status
			break
		}
	}

	if currentStepStatus == "pending" {
		// Update Onboarding Step
		var steps string
		for _, step := range onboardingSteps {
			status := step.Status
			if step.Step == "photos" {
				status = "done"
			}
			steps += fmt.Sprintf("%s:%s,", step.Step, status)
		}
		steps = strings.TrimSuffix(steps, ",")
		err = s.repo.UpdateOnboardingStepsUser(ctx, &repository.UpdateOnboardingSteps{
			ID:              UserID,
			OnboardingSteps: steps,
		})
		if err != nil {
			return errpkg.DefaultServiceError(
				errpkg.ErrInternal,
				err.Error(),
			)
		}
	}

	return nil

	return nil
}

func (s *OnboardingServiceMock) UpdateHobbyAndInterest(
	ctx context.Context,
	req *domain.UpdateHobbyAndInterest,
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
	err = s.repo.UpdateHobbyAndInterestProfile(ctx, &repository.UpdateHobbyAndInterest{
		ID:       profile.ID,
		Hobby:    strings.Join(req.Hobby, ","),
		Interest: strings.Join(req.Interest, ","),
	})
	if err != nil {
		return errpkg.DefaultServiceError(
			errpkg.ErrInternal,
			err.Error(),
		)
	}

	user, err := s.repo.GetUserById(ctx, UserID)
	if err != nil {
		return errpkg.DefaultServiceError(
			errpkg.ErrInternal,
			err.Error(),
		)
	}
	onboardingSteps := user.GetOnboardingSteps()
	currentStepStatus := ""
	for _, step := range onboardingSteps {
		if step.Step == "hobby-and-interest" {
			currentStepStatus = step.Status
			break
		}
	}

	if currentStepStatus == "pending" {
		// Update Onboarding Step
		var steps string
		for _, step := range onboardingSteps {
			status := step.Status
			if step.Step == "hobby-and-interest" {
				status = "done"
			}
			steps += fmt.Sprintf("%s:%s,", step.Step, status)
		}
		steps = strings.TrimSuffix(steps, ",")
		err = s.repo.UpdateOnboardingStepsUser(ctx, &repository.UpdateOnboardingSteps{
			ID:              UserID,
			OnboardingSteps: steps,
		})
		if err != nil {
			return errpkg.DefaultServiceError(
				errpkg.ErrInternal,
				err.Error(),
			)
		}
	}

	return nil
}

func (s *OnboardingServiceMock) UpdateLocation(
	ctx context.Context,
	req *domain.Location,
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
	err = s.repo.UpdateLocationProfile(ctx, &repository.UpdateLocation{
		ID:       profile.ID,
		Location: fmt.Sprintf("%s:%s", req.Longitude, req.Latitude),
	})
	if err != nil {
		return errpkg.DefaultServiceError(
			errpkg.ErrInternal,
			err.Error(),
		)
	}

	user, err := s.repo.GetUserById(ctx, UserID)
	if err != nil {
		return errpkg.DefaultServiceError(
			errpkg.ErrInternal,
			err.Error(),
		)
	}
	onboardingSteps := user.GetOnboardingSteps()
	currentStepStatus := ""
	for _, step := range onboardingSteps {
		if step.Step == "location" {
			currentStepStatus = step.Status
			break
		}
	}

	if currentStepStatus == "pending" {
		// Update Onboarding Step
		var steps string
		for _, step := range onboardingSteps {
			status := step.Status
			if step.Step == "location" {
				status = "done"
			}
			steps += fmt.Sprintf("%s:%s,", step.Step, status)
		}
		steps = strings.TrimSuffix(steps, ",")
		err = s.repo.UpdateOnboardingStepsUser(ctx, &repository.UpdateOnboardingSteps{
			ID:              UserID,
			OnboardingSteps: steps,
		})
		if err != nil {
			return errpkg.DefaultServiceError(
				errpkg.ErrInternal,
				err.Error(),
			)
		}
	}
	return nil

	return nil
}
