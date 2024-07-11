package stopwatch

import "fmt"

type EllapsedTime struct {
	duration    int64
	millisecond int64
	second      int64
	minute      int64
}

func NewEllapsedTime(millisec int64) *EllapsedTime {
	ms := millisec % 1000
	s := (millisec - ms) / 1000
	m := s / 60
	s %= 60

	return &EllapsedTime{
		duration:    millisec,
		millisecond: ms,
		second:      s,
		minute:      m,
	}
}

func (e *EllapsedTime) Duration() int64 {
	return e.duration
}

func (e *EllapsedTime) String() string {
	output := ""
	space := ""

	if e.minute > 0 {
		output = fmt.Sprintf("%dm", e.minute)
		space = " "
	}
	if e.second > 0 {
		output = fmt.Sprintf("%s%s%ds", output, space, e.second)
		space = " "
	}
	if e.millisecond > 0 || (e.millisecond == 0 && output == "") {
		output = fmt.Sprintf("%s%s%dms", output, space, e.millisecond)
	}

	return output
}
