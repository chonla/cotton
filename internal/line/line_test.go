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

func TestReflectValue(t *testing.T) {
	lString, _ := line.Line("\"123\"").ReflectJSValue()
	lInt, _ := line.Line("123").ReflectJSValue()
	lFloat, _ := line.Line("123.2").ReflectJSValue()
	lBoolTrue, _ := line.Line("true").ReflectJSValue()
	lBoolFalse, _ := line.Line("false").ReflectJSValue()
	lNull, _ := line.Line("null").ReflectJSValue()
	_, err := line.Line("udya").ReflectJSValue()

	assert.Equal(t, "123", lString.(string))
	assert.Equal(t, float64(123), lInt.(float64))
	assert.Equal(t, float64(123.2), lFloat.(float64))
	assert.True(t, lBoolTrue.(bool))
	assert.False(t, lBoolFalse.(bool))
	assert.Nil(t, lNull)
	assert.Error(t, err)
}
