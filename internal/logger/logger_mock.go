package logger

import (
	"cotton/internal/result"

	"github.com/stretchr/testify/mock"
)

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) PrintTestCaseTitle(title string) error {
	args := m.Called(title)
	return args.Error(0)
}

func (m *MockLogger) PrintExecutableTitle(title string) error {
	args := m.Called(title)
	return args.Error(0)
}

func (m *MockLogger) PrintBlockTitle(title string) error {
	args := m.Called(title)
	return args.Error(0)
}

func (m *MockLogger) PrintTestResult(passed bool) error {
	args := m.Called(passed)
	return args.Error(0)
}

func (m *MockLogger) PrintInlineTestResult(passed bool) error {
	args := m.Called(passed)
	return args.Error(0)
}

func (m *MockLogger) PrintAssertionResults(assertions []result.AssertionResult) error {
	args := m.Called(assertions)
	return args.Error(0)
}

func (m *MockLogger) PrintAssertionResult(assertion result.AssertionResult) error {
	args := m.Called(assertion)
	return args.Error(0)
}

func (m *MockLogger) PrintRequest(req string) error {
	args := m.Called(req)
	return args.Error(0)
}
