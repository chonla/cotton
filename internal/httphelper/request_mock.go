package httphelper

import "github.com/stretchr/testify/mock"

type MockHTTPRequest struct {
	mock.Mock
}

func (m *MockHTTPRequest) Do() (*HTTPResponse, error) {
	args := m.Called()
	return args.Get(0).(*HTTPResponse), args.Error(1)
}

func (m *MockHTTPRequest) String() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockHTTPRequest) Specs() map[string]interface{} {
	args := m.Called()
	return args.Get(0).(map[string]interface{})
}
