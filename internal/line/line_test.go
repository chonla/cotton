package line_test

import (
	"cotton/internal/line"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrim(t *testing.T) {
	l := line.Line("  hello      ")

	assert.Equal(t, line.Line("hello"), l.Trim())
}

func TestCapture(t *testing.T) {
	l := line.Line("# Title")

	captured, ok := l.Capture("# (.*)", 1)

	assert.True(t, ok)
	assert.Equal(t, "Title", captured)
}

func TestOOBCapture(t *testing.T) {
	l := line.Line("# Title")

	captured, ok := l.Capture("# (.*)", 2)

	assert.False(t, ok)
	assert.Equal(t, "", captured)
}

func TestUnmatchedCapture(t *testing.T) {
	l := line.Line("# Title")

	captured, ok := l.Capture("## (.*)", 1)

	assert.False(t, ok)
	assert.Equal(t, "", captured)
}

func TestLookLike(t *testing.T) {
	l := line.Line("# Title")

	ok := l.LookLike("# (.*)")

	assert.True(t, ok)
}

func TestUnmatchedLookLike(t *testing.T) {
	l := line.Line("# Title")

	ok := l.LookLike("^## (.*)")

	assert.False(t, ok)
}
