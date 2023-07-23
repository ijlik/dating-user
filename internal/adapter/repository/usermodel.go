package repository

import (
	"database/sql"
	"github.com/ijlik/dating-user/internal/business/domain"
	"github.com/ijlik/dating-user/pkg/constant"
	"strings"
	"time"
)

type User struct {
	ID              string              `db:"id"`
	Phone           sql.NullString      `db:"phone"`
	Email           string              `db:"email"`
	Status          constant.UserStatus `db:"status"`
	OnboardingSteps string              `db:"onboarding_steps"`
	CreatedAt       time.Time           `db:"created_at"`
	UpdatedAt       sql.NullTime        `db:"updated_at"`
}

func (u *User) GetOnboardingSteps() []domain.OnboardingSteps {
	parts := strings.Split(u.OnboardingSteps, ",")
	var onboardingSteps []domain.OnboardingSteps

	for _, part := range parts {
		subParts := strings.Split(part, ":")
		if len(subParts) == 2 {
			step := strings.TrimSpace(subParts[0])
			status := strings.TrimSpace(subParts[1])

			onboardingSteps = append(onboardingSteps, domain.OnboardingSteps{
				Step:   step,
				Status: status,
			})
		}
	}
	return onboardingSteps
}

type CreateUser struct {
	Email           string `db:"email"`
	OnboardingSteps string `db:"onboarding_steps"`
}

func (u *CreateUser) RowData() []interface{} {
	var data = []interface{}{
		u.Email,
		u.OnboardingSteps,
	}
	return data
}

type UpdateOnboardingSteps struct {
	ID              string `db:"id"`
	OnboardingSteps string `db:"onboarding_steps"`
}

func (u *UpdateOnboardingSteps) RowData() []interface{} {
	var data = []interface{}{
		u.ID,
		u.OnboardingSteps,
	}
	return data
}

type UpdateStatus struct {
	ID     string              `db:"id"`
	Status constant.UserStatus `db:"status"`
}

func (u *UpdateStatus) RowData() []interface{} {
	var data = []interface{}{
		u.ID,
		u.Status,
	}
	return data
}
