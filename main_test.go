package main

// import (
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// )

// type MockTestSuites struct {
// 	mock.Mock

// 	Variables map[string]string
// }

// func (m *MockTestSuites) Run() {
// 	m.Called()
// }

// func (m *MockTestSuites) Summary() int {
// 	args := m.Called()
// 	return args.Int(0)
// }

// func (m *MockTestSuites) SetVariables(v map[string]string) {
// 	m.Called(v)
// }

// func TestDispatchWhenNoVarsPassedTo(t *testing.T) {
// 	suites := new(MockTestSuites)
// 	suites.On("Run")
// 	suites.On("Summary").Return(100)
// 	suites.On("SetVariables", map[string]string{})

// 	result := dispatchTests(suites, []string{}, false)

// 	suites.AssertNumberOfCalls(t, "Run", 1)
// 	suites.AssertNumberOfCalls(t, "Summary", 1)
// 	suites.AssertNumberOfCalls(t, "SetVariables", 0)
// 	assert.Equal(t, result, 100)
// }

// func TestDispatchWhenSomeVarsPassedTo(t *testing.T) {
// 	suites := new(MockTestSuites)
// 	suites.On("Run")
// 	suites.On("Summary").Return(100)
// 	suites.On("SetVariables", map[string]string{
// 		"a": "1",
// 		"b": "2",
// 	})

// 	result := dispatchTests(suites, []string{"a=1", "b=2"}, false)

// 	suites.AssertNumberOfCalls(t, "Run", 1)
// 	suites.AssertNumberOfCalls(t, "Summary", 1)
// 	suites.AssertCalled(t, "SetVariables", map[string]string{
// 		"a": "1",
// 		"b": "2",
// 	})
// 	assert.Equal(t, result, 100)
// }

// func TestDispatchWhenSomeInvalidVarsPassedTo(t *testing.T) {
// 	suites := new(MockTestSuites)
// 	suites.On("Run")
// 	suites.On("Summary").Return(100)
// 	suites.On("SetVariables", map[string]string{
// 		"a": "1",
// 		"b": "2",
// 	})

// 	result := dispatchTests(suites, []string{"a=1", "b=2", "ccc:d"}, false)

// 	suites.AssertNumberOfCalls(t, "Run", 1)
// 	suites.AssertNumberOfCalls(t, "Summary", 1)
// 	suites.AssertCalled(t, "SetVariables", map[string]string{
// 		"a": "1",
// 		"b": "2",
// 	})
// 	assert.Equal(t, result, 100)
// }
