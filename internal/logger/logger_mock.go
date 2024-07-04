package logger

import (
	"cotton/internal/result"
	"cotton/internal/stopwatch"
	"cotton/internal/variable"

	"github.com/stretchr/testify/mock"
)

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) PrintTestcaseTitle(title string) error {
	args := m.Called(title)
	return args.Error(0)
}

func (m *MockLogger) PrintExecutableTitle(title string) error {
	args := m.Called(title)
	return args.Error(0)
}

func (m *MockLogger) PrintTestResult(passed bool) error {
	args := m.Called(passed)
	return args.Error(0)
}

func (m *MockLogger) PrintTimeEllapsed(ellapsedTime *stopwatch.EllapsedTime) error {
	args := m.Called(ellapsedTime)
	return args.Error(0)
}

func (m *MockLogger) PrintInlineTimeEllapsed(ellapsedTime *stopwatch.EllapsedTime) error {
	args := m.Called(ellapsedTime)
	return args.Error(0)
}

func (m *MockLogger) PrintSectionedMessage(section, message string) error {
	args := m.Called(section, message)
	return args.Error(0)
}

func (m *MockLogger) PrintInlineTestResult(passed bool) error {
	args := m.Called(passed)
	return args.Error(0)
}

func (m *MockLogger) PrintAssertionResults(assertions []*result.AssertionResult) error {
	args := m.Called(assertions)
	return args.Error(0)
}

func (m *MockLogger) PrintAssertionResult(assertion *result.AssertionResult) error {
	args := m.Called(assertion)
	return args.Error(0)
}

func (m *MockLogger) PrintRequest(req string) error {
	args := m.Called(req)
	return args.Error(0)
}

func (m *MockLogger) PrintResponse(resp string) error {
	args := m.Called(resp)
	return args.Error(0)
}

func (m *MockLogger) PrintError(fileContext string, err error) error {
	args := m.Called(fileContext, err)
	return args.Error(0)
}

func (m *MockLogger) PrintTestcaseSequence(index, total int) error {
	args := m.Called(index, total)
	return args.Error(0)
}

func (m *MockLogger) PrintTestsuiteResult(testsuiteResult *result.TestsuiteResult) error {
	args := m.Called(testsuiteResult)
	return args.Error(0)
}

func (m *MockLogger) PrintSectionTitle(sectionTitle string) error {
	args := m.Called(sectionTitle)
	return args.Error(0)
}

func (m *MockLogger) PrintDebugMessage(message string) error {
	args := m.Called(message)
	return args.Error(0)
}

func (m *MockLogger) PrintDetailedDebugMessage(messages ...string) error {
	margs := make([]interface{}, len(messages))
	for _, message := range messages {
		margs = append(margs, message)
	}
	args := m.Called(margs)
	return args.Error(0)
}

func (m *MockLogger) PrintVariables(variables *variable.Variables) error {
	args := m.Called(variables)
	return args.Error(0)
}
