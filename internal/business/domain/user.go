package domain

type User struct {
	Email           string            `json:"email"`
	Status          string            `json:"status"`
	OnboardingSteps []OnboardingSteps `json:"onboarding_steps"`
}

type OnboardingSteps struct {
	Step   string `json:"step"`
	Status string `json:"status"`
}
