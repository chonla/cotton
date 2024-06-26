package clock

import "time"

type ClockWrapper interface {
	Now() time.Time
}

type Clock struct{}

func New() ClockWrapper {
	return &Clock{}
}

func (c *Clock) Now() time.Time {
	return time.Now()
}
