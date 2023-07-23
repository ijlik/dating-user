package disposable

import (
	disposable "github.com/rocketlaunchr/anti-disposable-email"
)

func ParseEmail(email string) (disposable.ParsedEmail, error) {
	item, err := disposable.ParseEmail(email, true)
	if err != nil {
		return disposable.ParsedEmail{}, err
	}

	return item, nil
}

func ValidateIsDisposable(email string) bool {
	item, err := ParseEmail(email)
	if err != nil {
		return true
	}

	return item.Disposable
}
