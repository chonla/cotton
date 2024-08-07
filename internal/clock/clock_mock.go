package clock

import (
	"time"

	"github.com/stretchr/testify/mock"
)

type MockClock struct {
	mock.Mock
}

func (m *MockClock) Now() time.Time {
	args := m.Called()
	return args.Get(0).(time.Time)
}

func (m *MockClock) Epoch() int64 {
	args := m.Called()
	return args.Get(0).(int64)
}
