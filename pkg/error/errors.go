package error

type ErrorService interface {
	Error() string
	GetCode() ErrCode
}

type errorService struct {
	Code ErrCode
	Msg  string
}

func (e *errorService) Error() string {
	return e.Msg
}

func (e *errorService) GetCode() ErrCode {
	return e.Code
}

func DefaultServiceError(code ErrCode, msg string) ErrorService {
	return &errorService{
		Code: code,
		Msg:  msg,
	}
}
