package httphelper

import (
	"github.com/stretchr/testify/mock"
)

type MockHTTPRequestParser struct {
	mock.Mock
}

func (m *MockHTTPRequestParser) Parse(requestString string) (Request, error) {
	args := m.Called(requestString)
	return args.Get(0).(Request), args.Error(1)
}
