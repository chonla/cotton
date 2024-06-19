package capture_test

import (
	"cotton/internal/capture"
	"cotton/internal/line"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCaptureFromPlainAssignment(t *testing.T) {
	mdLine := line.Line("* name:$.data.firstname")

	result, ok := capture.Try(mdLine)

	assert.True(t, ok)
	assert.Equal(t, &capture.Capture{
		Name:     "name",
		Selector: "$.data.firstname",
	}, result)
}

func TestCaptureFromInlineAssignment(t *testing.T) {
	mdLine := line.Line("* name:`$.data.firstname`")

	result, ok := capture.Try(mdLine)

	assert.True(t, ok)
	assert.Equal(t, &capture.Capture{
		Name:     "name",
		Selector: "$.data.firstname",
	}, result)
}

func TestCaptureFromMoreIndentsPlainAssignment(t *testing.T) {
	mdLine := line.Line("*   name:$.data.firstname")

	result, ok := capture.Try(mdLine)

	assert.True(t, ok)
	assert.Equal(t, &capture.Capture{
		Name:     "name",
		Selector: "$.data.firstname",
	}, result)
}

func TestCaptureFromPlainAssignmentWithWhiteSpaces(t *testing.T) {
	mdLine := line.Line("*   name :  $.data.firstname")

	result, ok := capture.Try(mdLine)

	assert.True(t, ok)
	assert.Equal(t, &capture.Capture{
		Name:     "name",
		Selector: "$.data.firstname",
	}, result)
}
