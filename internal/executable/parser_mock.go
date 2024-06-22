package executable

import (
	"cotton/internal/line"

	"github.com/stretchr/testify/mock"
)

type MockExecutableParser struct {
	mock.Mock
}

func (m *MockExecutableParser) FromMarkdownLines(mdLines []line.Line) (*Executable, error) {
	args := m.Called(mdLines)
	return args.Get(0).(*Executable), args.Error(1)
}

func (m *MockExecutableParser) FromMarkdownFile(mdFileName string) (*Executable, error) {
	args := m.Called(mdFileName)
	return args.Get(0).(*Executable), args.Error(1)
}
