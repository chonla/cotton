package stopwatch_test

import (
	"cotton/internal/stopwatch"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRenderEllapsedTimeLowerBoundMillisec(t *testing.T) {
	ellapsed := stopwatch.NewEllapsedTime(0)

	result := ellapsed.String()

	assert.Equal(t, "0ms", result)
}

func TestRenderEllapsedTimeUpperBoundMillisec(t *testing.T) {
	ellapsed := stopwatch.NewEllapsedTime(999)

	result := ellapsed.String()

	assert.Equal(t, "999ms", result)
}

func TestRenderEllapsedTimeLowerBoundSec(t *testing.T) {
	ellapsed := stopwatch.NewEllapsedTime(1000)

	result := ellapsed.String()

	assert.Equal(t, "1s", result)
}

func TestRenderEllapsedTimeUpperBoundSec(t *testing.T) {
	ellapsed := stopwatch.NewEllapsedTime(59000)

	result := ellapsed.String()

	assert.Equal(t, "59s", result)
}

func TestRenderEllapsedTimeLowerBoundMin(t *testing.T) {
	ellapsed := stopwatch.NewEllapsedTime(60000)

	result := ellapsed.String()

	assert.Equal(t, "1m", result)
}

func TestRenderEllapsedTimeMinWithSec(t *testing.T) {
	ellapsed := stopwatch.NewEllapsedTime(62000)

	result := ellapsed.String()

	assert.Equal(t, "1m 2s", result)
}

func TestRenderEllapsedTimeMinWithMillisec(t *testing.T) {
	ellapsed := stopwatch.NewEllapsedTime(60020)

	result := ellapsed.String()

	assert.Equal(t, "1m 20ms", result)
}

func TestRenderEllapsedTimeSecWithMillisec(t *testing.T) {
	ellapsed := stopwatch.NewEllapsedTime(59020)

	result := ellapsed.String()

	assert.Equal(t, "59s 20ms", result)
}

func TestRenderEllapsedTimeMinWithSecWithMillisec(t *testing.T) {
	ellapsed := stopwatch.NewEllapsedTime(61020)

	result := ellapsed.String()

	assert.Equal(t, "1m 1s 20ms", result)
}
