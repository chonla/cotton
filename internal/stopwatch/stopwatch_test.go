package stopwatch_test

import (
	"cotton/internal/clock"
	"cotton/internal/stopwatch"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStartThenStopImmidiately(t *testing.T) {
	mockStartTime, _ := time.Parse(time.RFC3339, "2024-06-26T15:27:05+07:00")
	mockStopTime, _ := time.Parse(time.RFC3339, "2024-06-26T15:27:05+07:00")

	mockClock := new(clock.MockClock)
	mockClock.On("Now").Return(mockStartTime).Once()
	mockClock.On("Now").Return(mockStopTime).Once()

	watch := stopwatch.New(mockClock)
	watch.Start()
	result := watch.Stop()

	assert.Equal(t, "0ms", result.String())
	mockClock.AssertExpectations(t)
}

func TestStartThenStop50secLater(t *testing.T) {
	mockStartTime, _ := time.Parse(time.RFC3339, "2024-06-26T15:27:05+07:00")
	mockStopTime, _ := time.Parse(time.RFC3339, "2024-06-26T15:27:55+07:00")

	mockClock := new(clock.MockClock)
	mockClock.On("Now").Return(mockStartTime).Once()
	mockClock.On("Now").Return(mockStopTime).Once()

	watch := stopwatch.New(mockClock)
	watch.Start()
	result := watch.Stop()

	assert.Equal(t, "50s", result.String())
	mockClock.AssertExpectations(t)
}
