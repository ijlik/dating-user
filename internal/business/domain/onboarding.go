package domain

import (
	errpkg "github.com/ijlik/dating-user/pkg/error"
	"mime/multipart"
	"strings"
	"time"
)

type UpdatePersonalInfo struct {
	Name       string    `json:"name"`
	BirthDates string    `json:"birth_date"`
	BirthDate  time.Time `json:"-"`
	Gender     string    `json:"gender"`
}

func (p *UpdatePersonalInfo) Validate() errpkg.ErrorService {
	if p.Name == "" {
		return errpkg.DefaultServiceError(errpkg.ErrBadRequest, "missing name")
	}
	if strings.ToLower(p.Gender) != "male" && strings.ToLower(p.Gender) != "female" {
		return errpkg.DefaultServiceError(errpkg.ErrBadRequest, "allowed gender Male or Female")
	}
	p.Gender = strings.Title(strings.ToLower(p.Gender))
	date, err := time.Parse(time.RFC3339, p.BirthDates)
	if err != nil {
		return errpkg.DefaultServiceError(errpkg.ErrBadRequest, "invalid format birth date")
	}

	p.BirthDate = date.UTC()

	return nil
}

type UpdatePhotos struct {
	Photos []*multipart.FileHeader `form:"photos"`
}

func (p *UpdatePhotos) Validate() errpkg.ErrorService {
	if len(p.Photos) == 0 {
		return errpkg.DefaultServiceError(errpkg.ErrBadRequest, "missing photos: at least one image file is required")
	}
	maxSize := int64(5 * 1024 * 1024)

	for _, file := range p.Photos {
		if file == nil {
			return errpkg.DefaultServiceError(errpkg.ErrBadRequest, "invalid file: one of the images is missing")
		}

		contentType := file.Header.Get("Content-Type")
		if contentType != "image/jpeg" && contentType != "image/png" {
			return errpkg.DefaultServiceError(errpkg.ErrBadRequest, "invalid file type: only JPEG and PNG images are allowed")
		}

		if file.Size > maxSize {
			return errpkg.DefaultServiceError(errpkg.ErrBadRequest, "file size too large: maximum allowed size is 5MB")
		}
	}

	return nil
}

type UpdateHobbyAndInterest struct {
	Hobby    []string `json:"hobby"`
	Interest []string `json:"interest"`
}

func (p *UpdateHobbyAndInterest) Validate() errpkg.ErrorService {
	for _, hobby := range p.Hobby {
		if hobby == "" {
			return errpkg.DefaultServiceError(errpkg.ErrBadRequest, "missing hobby")
		}
	}

	for _, interest := range p.Interest {
		if interest == "" {
			return errpkg.DefaultServiceError(errpkg.ErrBadRequest, "missing interest")
		}
	}

	return nil
}

type Location struct {
	Longitude string `json:"longitude"`
	Latitude  string `json:"latitude"`
	Url       string `json:"url"`
}

func (u *Location) Validate() errpkg.ErrorService {
	if u.Longitude == "" {
		return errpkg.DefaultServiceError(errpkg.ErrBadRequest, "missing longitude")
	}
	if !rgxLongitude.Match([]byte(u.Longitude)) {
		return errpkg.DefaultServiceError(errpkg.ErrBadRequest, "invalid format longitude")
	}
	if u.Latitude == "" {
		return errpkg.DefaultServiceError(errpkg.ErrBadRequest, "missing latitude")
	}
	if !rgxLatitude.Match([]byte(u.Latitude)) {
		return errpkg.DefaultServiceError(errpkg.ErrBadRequest, "invalid format latitude")
	}

	return nil
}
