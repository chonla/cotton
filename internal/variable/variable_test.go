package variable_test

import (
	"cotton/internal/variable"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateEmptyVariables(t *testing.T) {
	result := variable.New()

	assert.Equal(t, []string{}, result.Names())
}

func TestSetIntVariable(t *testing.T) {
	vars := variable.New()

	vars.Set("k", 1)

	result, err := vars.ValueOf("k")

	assert.NoError(t, err)
	assert.Equal(t, 1, result.(int))
}

func TestSetStringVariable(t *testing.T) {
	vars := variable.New()

	vars.Set("k", "1")

	result, err := vars.ValueOf("k")

	assert.NoError(t, err)
	assert.Equal(t, "1", result.(string))
}

func TestGetNames(t *testing.T) {
	vars := variable.New()

	vars.Set("k", "1")
	vars.Set("v", "2")
	vars.Set("x", "3")

	result := vars.Names()

	assert.Equal(t, []string{"k", "v", "x"}, result)
}

func TestGetMap(t *testing.T) {
	vars := variable.New()

	vars.Set("k", 1)
	vars.Set("v", errors.New("yeah!"))
	vars.Set("x", 3)

	result := vars.ToStringMap()

	assert.Equal(t, map[string]interface{}{
		"k": "1",
		"v": "yeah!",
		"x": "3",
	}, result)
}

func TestMergeVars(t *testing.T) {
	vars1 := variable.New()
	vars1.Set("a", 8)
	vars1.Set("b", 9)
	vars1.Set("c", 10)

	vars2 := variable.New()
	vars2.Set("c", "1")
	vars2.Set("d", "2")
	vars2.Set("e", "3")

	mergedVars := vars1.MergeWith(vars2)

	result := mergedVars.ToStringMap()

	assert.Equal(t, map[string]interface{}{
		"a": "8",
		"b": "9",
		"c": "1",
		"d": "2",
		"e": "3",
	}, result)
}

func TestResetVars(t *testing.T) {
	vars := variable.New()
	vars.Set("a", 8)
	vars.Set("b", 9)
	vars.Set("c", 10)
	vars.Set("b", 11)

	result := vars.ToStringMap()

	assert.Equal(t, map[string]interface{}{
		"a": "8",
		"b": "11",
		"c": "10",
	}, result)
}

func TestCloneShouldReturnACopyOfVars(t *testing.T) {
	vars1 := variable.New()
	vars1.Set("a", 8)
	vars1.Set("b", 9)
	vars1.Set("c", 10)

	vars2 := vars1.Clone()

	assert.Equal(t, vars1, vars2)
}

func TestCloneShouldNotMutateTheOrigin(t *testing.T) {
	vars1 := variable.New()
	vars1.Set("a", 8)
	vars1.Set("b", 9)
	vars1.Set("c", 10)

	vars2 := vars1.Clone()
	vars2.Set("a", 11)

	result1 := vars1.ToStringMap()
	result2 := vars2.ToStringMap()

	assert.Equal(t, map[string]interface{}{
		"a": "8",
		"b": "9",
		"c": "10",
	}, result1)
	assert.Equal(t, map[string]interface{}{
		"a": "11",
		"b": "9",
		"c": "10",
	}, result2)
}
