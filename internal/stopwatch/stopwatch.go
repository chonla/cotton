package stopwatch

import (
	"cotton/internal/clock"
)

type Stopwatch struct {
	clock            clock.ClockWrapper
	start            int64
	ellapsedMillisec int64
}

func New(clock clock.ClockWrapper) *Stopwatch {
	return &Stopwatch{
		clock: clock,
	}
}

func (s *Stopwatch) Start() {
	nowMs := s.clock.Now().UnixMilli()
	s.start = nowMs
}

func (s *Stopwatch) Stop() *EllapsedTime {
	nowMs := s.clock.Now().UnixMilli()
	s.ellapsedMillisec = nowMs - s.start
	return NewEllapsedTime(s.ellapsedMillisec)
}
