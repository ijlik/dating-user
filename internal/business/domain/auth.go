package domain

import (
	"errors"
	"regexp"

	pkgdisposable "github.com/ijlik/dating-user/pkg/disposable"
)

var (
	rgxEmail     = regexp.MustCompile(`^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`)
	rgxLatitude  = regexp.MustCompile(`^(-?\d+(\.\d+)?)$`)
	rgxLongitude = regexp.MustCompile(`^(-?\d+(\.\d+)?)$`)
)

type ResendOtpRequest struct {
	Email string `json:"email"`
}

func (req *ResendOtpRequest) Validate(allowedDisposableEmail bool) error {
	if req.Email == "" {
		return errors.New("email is required")
	}

	if !rgxEmail.Match([]byte(req.Email)) {
		return errors.New("invalid email")
	}

	if !allowedDisposableEmail && pkgdisposable.ValidateIsDisposable(req.Email) {
		return errors.New("free and disposable email isnt allowed")
	}

	return nil
}

type ResendOtpResponse struct {
	Message           string `json:"message"`
	ResendOTPInterval int    `json:"resendOTPInterval"`
}

type AuthRequest struct {
	Email string `json:"email"`
	Otp   string `json:"otp"`
}

func (req *AuthRequest) Validate() error {
	if req.Email == "" {
		return errors.New("email is required")
	}

	if !rgxEmail.Match([]byte(req.Email)) {
		return errors.New("invalid email")
	}

	if req.Otp == "" {
		return errors.New("otp is required")
	}

	return nil
}

type AuthResponse struct {
	AccessToken string `json:"access_token"`
}
