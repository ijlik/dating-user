package error

import "net/http"

type ErrCode int

func GetHttpStatus(code ErrCode) int {
	val, ok := mapHttpStatus[code]
	if !ok {
		return http.StatusOK
	}

	return val
}

func GetCode(code ErrCode) string {
	val, ok := mapCode[code]
	if !ok {
		return "00"
	}

	return val
}

func GetMessage(code ErrCode) string {
	val, ok := mapText[code]
	if !ok {
		return "Success"
	}

	return val
}

const (
	ErrBadRequest ErrCode = iota + 1
	ErrInternal
	ErrInvalidToken
	ErrNotFound
	ErrAlreadyRegistered
	ErrEmptyPassword
	ErrFailedToSendDeeplink
	ErrTemporaryBlocked
	ErrUnauthorize
	ErrInvalidPassword
	ErrReloginNeeded
	ErrTokenAlreadyUsed
	ErrMaxUserReached
	ErrAccessLimited
)

var mapCode = map[ErrCode]string{
	ErrBadRequest:           "01",
	ErrInternal:             "99",
	ErrInvalidToken:         "02",
	ErrNotFound:             "03",
	ErrAlreadyRegistered:    "04",
	ErrEmptyPassword:        "05",
	ErrFailedToSendDeeplink: "06",
	ErrTemporaryBlocked:     "07",
	ErrUnauthorize:          "08",
	ErrInvalidPassword:      "09",
	ErrReloginNeeded:        "10",
	ErrTokenAlreadyUsed:     "11",
	ErrMaxUserReached:       "12",
	ErrAccessLimited:        "13",
}

var mapHttpStatus = map[ErrCode]int{
	ErrBadRequest:           http.StatusBadRequest,
	ErrInternal:             http.StatusInternalServerError,
	ErrInvalidToken:         http.StatusUnauthorized,
	ErrNotFound:             http.StatusNotFound,
	ErrAlreadyRegistered:    http.StatusConflict,
	ErrEmptyPassword:        http.StatusUnprocessableEntity,
	ErrFailedToSendDeeplink: http.StatusInternalServerError,
	ErrTemporaryBlocked:     http.StatusUnprocessableEntity,
	ErrUnauthorize:          http.StatusUnauthorized,
	ErrInvalidPassword:      http.StatusBadRequest,
	ErrReloginNeeded:        http.StatusUnauthorized,
	ErrTokenAlreadyUsed:     http.StatusUnprocessableEntity,
	ErrMaxUserReached:       http.StatusUnprocessableEntity,
	ErrAccessLimited:        http.StatusForbidden,
}

var mapText = map[ErrCode]string{
	ErrBadRequest:           "Bad Request",
	ErrInternal:             "Internal Error",
	ErrInvalidToken:         "Token Invalid",
	ErrNotFound:             "Not Found",
	ErrAlreadyRegistered:    "User Already Registered, Please Login",
	ErrEmptyPassword:        "User Already Registered, Please Set Password",
	ErrFailedToSendDeeplink: "Failed To Send DeepLink",
	ErrTemporaryBlocked:     "User Temporary Blocked",
	ErrUnauthorize:          "Unauthorize",
	ErrInvalidPassword:      "Invalid password",
	ErrReloginNeeded:        "User Already Logout Please Re-Login",
	ErrTokenAlreadyUsed:     "Token Already Use",
	ErrMaxUserReached:       "Maximum 5 Users",
	ErrAccessLimited:        "Access limited",
}
