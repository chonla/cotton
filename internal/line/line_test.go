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

	cap, ok := l.Capture("# (.*)", 1)

	assert.True(t, ok)
	assert.Equal(t, "Title", cap)
}

func TestOOBCapture(t *testing.T) {
	l := line.Line("# Title")

	cap, ok := l.Capture("# (.*)", 2)

	assert.False(t, ok)
	assert.Equal(t, "", cap)
}

func TestUnmatchedCapture(t *testing.T) {
	l := line.Line("# Title")

	cap, ok := l.Capture("## (.*)", 1)

	assert.False(t, ok)
	assert.Equal(t, "", cap)
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

func TestReplaceSingleString(t *testing.T) {
	l := line.Line("<rootDir>/some/path")

	newPath := l.Replace("<rootDir>", "/tmp")

	assert.Equal(t, "/tmp/some/path", newPath)
}

func TestReplaceMultipleStrings(t *testing.T) {
	l := line.Line("<rootDir>/some/path<rootDir>")

	newPath := l.Replace("<rootDir>", "/tmp")

	assert.Equal(t, "/tmp/some/path/tmp", newPath)
}

func TestReplaceUnmatchString(t *testing.T) {
	l := line.Line("/some/path")

	newPath := l.Replace("<rootDir>", "/tmp")

	assert.Equal(t, "/some/path", newPath)
}
