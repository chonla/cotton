package template_test

import (
	"cotton/internal/template"
	"cotton/internal/variable"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApplyTemplateWithVariables(t *testing.T) {
	vars := variable.New()

	vars.Set("Name", "Sam")
	vars.Set("Age", 18)

	text := "My name is {{Name}}. I am {{Age}} years old."

	result := template.New(text).Apply(vars)

	expected := "My name is Sam. I am 18 years old."

	assert.Equal(t, expected, result)
}

func TestApplyTemplateWithRepeatedVariables(t *testing.T) {
	vars := variable.New()

	vars.Set("Name", "Sam")
	vars.Set("Age", 18)

	text := "My name is {{Name}}. {{Name}} is {{Age}} years old."

	result := template.New(text).Apply(vars)

	expected := "My name is Sam. Sam is 18 years old."

	assert.Equal(t, expected, result)
}

func TestVariableTagWithWhitespacesInsideWillNotMatch(t *testing.T) {
	vars := variable.New()

	vars.Set("Name", "Sam")
	vars.Set("Age", 18)

	text := "My name is {{ Name}}. {{Name  }} is {{ Age }} years old. {{Name}} is so confused."

	result := template.New(text).Apply(vars)

	expected := "My name is {{ Name}}. {{Name  }} is {{ Age }} years old. Sam is so confused."

	assert.Equal(t, expected, result)
}

func TestNilVariables(t *testing.T) {
	var vars *variable.Variables

	text := "My name is {{ Name}}. {{Name  }} is {{ Age }} years old. {{Name}} is so confused."

	result := template.New(text).Apply(vars)

	expected := "My name is {{ Name}}. {{Name  }} is {{ Age }} years old. {{Name}} is so confused."

	assert.Equal(t, expected, result)
}
