package pagination

import (
	"math"
)

type Pagination struct {
	Limit      int         `json:"limit"`
	Page       int         `json:"page"`
	NextPage   int         `json:"nextPage"`
	Offset     int         `json:"-"`
	TotalData  int         `json:"totalData"`
	TotalPages int         `json:"totalPages"`
	Data       interface{} `json:"data"`
}

func NewPaginate(limit int, page int) *Pagination {
	return &Pagination{
		Limit:  limit,
		Page:   page,
		Offset: (page - 1) * limit,
	}
}

func (p *Pagination) SetData(data interface{}, count int64) {
	p.Data = data
	p.TotalData = int(count)

	p.TotalPages = int(math.Ceil(float64(count) / float64(p.Limit)))
	if math.Ceil(float64(count)/float64(p.Limit)) > float64(p.Page) {
		p.NextPage = p.Page + 1
	}
}
