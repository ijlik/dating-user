package domain

import (
	"github.com/ijlik/dating-user/pkg/constant"
	errpkg "github.com/ijlik/dating-user/pkg/error"
	"strings"
)

type PaymentRequest struct {
	Amount        float32                `json:"amount"`
	Identifier    string                 `json:"identifier"`
	Method        string                 `json:"payment_method"`
	PaymentMethod constant.PaymentMethod `json:"-"`
	PaymentData   string                 `json:"payment_data"`
}

func (p *PaymentRequest) Validate() errpkg.ErrorService {
	if p.Amount < 0 {
		return errpkg.DefaultServiceError(errpkg.ErrBadRequest, "missing amount")
	}
	if p.Identifier == "" {
		return errpkg.DefaultServiceError(errpkg.ErrBadRequest, "missing identifier")
	}
	if p.Method == "" {
		return errpkg.DefaultServiceError(errpkg.ErrBadRequest, "missing payment method")
	}
	paymentMethod := constant.GetPaymentMethod(strings.Title(strings.ToLower(p.Method)))
	if paymentMethod == "unknown" {
		return errpkg.DefaultServiceError(errpkg.ErrBadRequest, "payment method not allowed")
	}
	p.PaymentMethod = paymentMethod
	if p.PaymentData == "" {
		return errpkg.DefaultServiceError(errpkg.ErrBadRequest, "missing payment data")
	}

	return nil
}
