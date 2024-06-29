package variable_test

import (
	"cotton/internal/line"
	"cotton/internal/variable"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSettingVariableFromLineData(t *testing.T) {
	mdLine := line.Line("* name:\"$.data.firstname\"")

	result, ok := variable.Try(mdLine)

	assert.True(t, ok)
	assert.Equal(t, &variable.Variable{
		Name:  "name",
		Value: "$.data.firstname",
	}, result)
}

func TestSettingVariableFromLineDataWithoutQuoteConsideringAStringIfFailedToParse(t *testing.T) {
	mdLine := line.Line("* name:$.data.firstname")

	result, ok := variable.Try(mdLine)

	assert.True(t, ok)
	assert.Equal(t, &variable.Variable{
		Name:  "name",
		Value: "$.data.firstname",
	}, result)
}

func TestSettingVariableFromLineDataAsNumber(t *testing.T) {
	mdLine := line.Line("* name:123")

	result, ok := variable.Try(mdLine)

	assert.True(t, ok)
	assert.Equal(t, &variable.Variable{
		Name:  "name",
		Value: float64(123),
	}, result)
}
