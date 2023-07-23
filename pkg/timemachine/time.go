package timemachine

import (
	"time"

	"github.com/jinzhu/now"
)

type TimeMachine interface {
	Now() time.Time
	GetStartAndEndDayTime() (time.Time, time.Time)
	GetStartAndEndMonthTime() (time.Time, time.Time)
	GetTimeAfterNow(t time.Duration) time.Time
}

type machine struct {
	config now.Config
}

func NewTimeMachine() TimeMachine {

	location, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return nil
	}

	timeConfig := now.Config{
		WeekStartDay: time.Monday,
		TimeLocation: location,
		TimeFormats:  []string{"2006-01-02 15:04:05"},
	}

	return &machine{
		config: timeConfig,
	}
}

func (m *machine) Now() time.Time {
	return time.Now().UTC()
}

func (m *machine) GetStartAndEndDayTime() (time.Time, time.Time) {
	timeNow := m.config.With(m.Now())
	return timeNow.BeginningOfDay(), timeNow.EndOfDay()
}

func (m *machine) GetStartAndEndMonthTime() (time.Time, time.Time) {
	timeNow := m.config.With(m.Now())
	return timeNow.BeginningOfMonth(), timeNow.EndOfMonth()
}

func (m *machine) GetTimeAfterNow(t time.Duration) time.Time {
	timeNow := m.config.With(m.Now())
	return timeNow.Add(t)
}
