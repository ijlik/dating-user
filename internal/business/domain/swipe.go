package domain

import errpkg "github.com/ijlik/dating-user/pkg/error"

type SwipeRequest struct {
	SwiperId string `json:"swiper_id"`
	SwipedId string `json:"swiped_id"`
	IsLike   bool   `json:"is_like"`
}

func (s *SwipeRequest) Validate() errpkg.ErrorService {
	if s.SwiperId == "" {
		return errpkg.DefaultServiceError(errpkg.ErrBadRequest, "missing swiper id")
	}
	if s.SwipedId == "" {
		return errpkg.DefaultServiceError(errpkg.ErrBadRequest, "missing swiped id")
	}

	return nil
}
