package service

import (
	"context"
	"fmt"
	"github.com/ijlik/dating-user/internal/adapter/repository"
	"github.com/ijlik/dating-user/internal/business/domain"
	errpkg "github.com/ijlik/dating-user/pkg/error"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func (s *service) UpdatePersonalInfo(
	ctx context.Context,
	req *domain.UpdatePersonalInfo,
	UserID string,
) errpkg.ErrorService {
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

func (s *service) UpdatePhotos(
	ctx context.Context,
	req *domain.UpdatePhotos,
	UserID string,
) errpkg.ErrorService {
	profile, err := s.repo.GetProfileByUserID(ctx, UserID)
	if err != nil {
		return errpkg.DefaultServiceError(
			errpkg.ErrInternal,
			err.Error(),
		)
	}

	uploadDir := fmt.Sprintf("storage/photos/%s", profile.ID)
	var allPhotos string
	if err = os.RemoveAll(uploadDir); err != nil {
		return errpkg.DefaultServiceError(
			errpkg.ErrInternal,
			err.Error(),
		)
	}
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return errpkg.DefaultServiceError(
			errpkg.ErrInternal,
			err.Error(),
		)
	}
	for _, photo := range req.Photos {
		// Generate a unique filename for each photo
		fileName := filepath.Join(uploadDir, photo.Filename)

		// Create the destination file
		dst, err := os.Create(fileName)
		if err != nil {
			return errpkg.DefaultServiceError(
				errpkg.ErrInternal,
				err.Error(),
			)
		}
		defer dst.Close()

		// Open the uploaded file
		src, err := photo.Open()
		if err != nil {
			return errpkg.DefaultServiceError(
				errpkg.ErrInternal,
				err.Error(),
			)
		}
		defer src.Close()

		// Copy the uploaded file to the destination
		if _, err = io.Copy(dst, src); err != nil {
			return errpkg.DefaultServiceError(
				errpkg.ErrInternal,
				err.Error(),
			)
		}
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
}

func (s *service) UpdateHobbyAndInterest(
	ctx context.Context,
	req *domain.UpdateHobbyAndInterest,
	UserID string,
) errpkg.ErrorService {
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

func (s *service) UpdateLocation(
	ctx context.Context,
	req *domain.Location,
	UserID string,
) errpkg.ErrorService {
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
}
