package capture_test

import (
	"cotton/internal/capture"
	"cotton/internal/line"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCaptureFromInlineAssignment(t *testing.T) {
	mdLine := line.Line("* name:`$.data.firstname`")

	result, ok := capture.Try(mdLine)

	assert.True(t, ok)
	assert.Equal(t, &capture.Capture{
		Name:     "name",
		Selector: "$.data.firstname",
	}, result)
}

func TestCaptureFromInlineAssignmentPlus(t *testing.T) {
	mdLine := line.Line("+ name:`$.data.firstname`")

	result, ok := capture.Try(mdLine)

	assert.True(t, ok)
	assert.Equal(t, &capture.Capture{
		Name:     "name",
		Selector: "$.data.firstname",
	}, result)
}

func TestCaptureFromInlineAssignmentMinus(t *testing.T) {
	mdLine := line.Line("- name:`$.data.firstname`")

	result, ok := capture.Try(mdLine)

	assert.True(t, ok)
	assert.Equal(t, &capture.Capture{
		Name:     "name",
		Selector: "$.data.firstname",
	}, result)
}

func TestCaptureFromInlineAssignmentOrderedList(t *testing.T) {
	mdLine := line.Line("33. name:`$.data.firstname`")

	result, ok := capture.Try(mdLine)

	assert.True(t, ok)
	assert.Equal(t, &capture.Capture{
		Name:     "name",
		Selector: "$.data.firstname",
	}, result)
}

func TestCaptureFromMoreIndentsPlainAssignment(t *testing.T) {
	mdLine := line.Line("*   name:`$.data.firstname`")

	result, ok := capture.Try(mdLine)

	assert.True(t, ok)
	assert.Equal(t, &capture.Capture{
		Name:     "name",
		Selector: "$.data.firstname",
	}, result)
}

func TestCaptureFromPlainAssignmentWithWhiteSpaces(t *testing.T) {
	mdLine := line.Line("*   name :  `$.data.firstname`")

	result, ok := capture.Try(mdLine)

	assert.True(t, ok)
	assert.Equal(t, &capture.Capture{
		Name:     "name",
		Selector: "$.data.firstname",
	}, result)
}

func TestCaptureFromNonCapture(t *testing.T) {
	mdLine := line.Line("=   name :  `$.data.firstname`")

	result, ok := capture.Try(mdLine)

	assert.False(t, ok)
	assert.Nil(t, result)
}

func TestCapturesAreSimilarIfTheyHaveSameNameAndSelector(t *testing.T) {
	cap1 := capture.New("cool.name", "cool.value")
	cap2 := capture.New("cool.name", "cool.value")

	result1 := cap1.SimilarTo(cap2)
	result2 := cap2.SimilarTo(cap1)

	assert.True(t, result1)
	assert.True(t, result2)
}

func TestCapturesAreSimilarIfTheyHaveDifferentNames(t *testing.T) {
	cap1 := capture.New("cool.Name", "cool.value")
	cap2 := capture.New("cool.name", "cool.value")

	result1 := cap1.SimilarTo(cap2)
	result2 := cap2.SimilarTo(cap1)

	assert.False(t, result1)
	assert.False(t, result2)
}

func TestCapturesAreSimilarIfTheyHaveDifferentValues(t *testing.T) {
	cap1 := capture.New("cool.name", "cool.Value")
	cap2 := capture.New("cool.name", "cool.value")

	result1 := cap1.SimilarTo(cap2)
	result2 := cap2.SimilarTo(cap1)

	assert.False(t, result1)
	assert.False(t, result2)
}

func TestCloningMustReturnSameValue(t *testing.T) {
	cap1 := capture.New("cool.name", "cool.Value")
	cap2 := cap1.Clone()

	result1 := cap1.SimilarTo(cap2)
	result2 := cap2.SimilarTo(cap1)

	assert.True(t, result1)
	assert.True(t, result2)
}

func TestCloningMustReturnACopyOfTheSource(t *testing.T) {
	cap1 := capture.New("cool.name", "cool.Value")
	cap2 := cap1.Clone()
	cap2.Name = "cooler.name"

	result1 := cap1.SimilarTo(cap2)
	result2 := cap2.SimilarTo(cap1)

	assert.False(t, result1)
	assert.False(t, result2)
}
