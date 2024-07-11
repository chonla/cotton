package clock

import "time"

type ClockWrapper interface {
	Now() time.Time
	Epoch() int64
}

type Clock struct{}

func New() ClockWrapper {
	return &Clock{}
}

func (c *Clock) Now() time.Time {
	return time.Now()
}

func (c *Clock) Epoch() int64 {
	return time.Now().Unix()
}
