package repository

import "database/sql"

type Swipe struct {
	SwiperId string       `db:"swiper_id"`
	SwipedId string       `db:"swiped_id"`
	IsLike   sql.NullBool `db:"is_like"`
}

func (s *Swipe) RowData() []interface{} {
	var data = []interface{}{
		s.SwiperId,
		s.SwipedId,
		s.IsLike,
	}
	return data
}
