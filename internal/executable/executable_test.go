package executable_test

import (
	"cotton/internal/capture"
	"cotton/internal/executable"
	"cotton/internal/logger"
	"cotton/internal/request"
	"cotton/internal/response"
	"cotton/internal/variable"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockHTTPRequest struct {
	mock.Mock
}

func (m *MockHTTPRequest) Similar(anotherRequest request.Request) bool {
	margs := m.Called(anotherRequest)
	return margs.Bool(0)
}

func (m *MockHTTPRequest) Data() []byte {
	margs := m.Called()
	return margs.Get(0).([]byte)
}

func (m *MockHTTPRequest) Do() (*http.Response, error) {
	margs := m.Called()
	return margs.Get(0).(*http.Response), margs.Error(1)
}

func (m *MockHTTPRequest) String() string {
	margs := m.Called()
	return margs.String(0)
}

func TestParsingCaptureFromResponse(t *testing.T) {
	recorder := httptest.NewRecorder()
	recorder.Body.WriteString(`{"form":{"key1":"val1","key2":"val2"}}`)
	recorder.Code = 200
	recorder.Header().Add("Content-Length", "15")
	recorder.Header().Add("Content-Type", "application/json")
	recorder.Flush()

	resp := recorder.Result()

	reqMock := new(MockHTTPRequest)
	reqMock.On("Do").Return(resp, nil)

	ex := &executable.Executable{
		Title:   "test",
		Request: reqMock,
		Captures: []*capture.Capture{
			{
				Name:     "var1",
				Selector: "Body.form.key1",
			},
			{
				Name:     "var2",
				Selector: "Body.form.key2",
			},
		},
	}

	expectedVariables := variable.New()
	expectedVariables.Set("var1", &response.DataValue{
		Value:       "val1",
		TypeName:    "string",
		IsUndefined: false,
	})
	expectedVariables.Set("var2", &response.DataValue{
		Value:       "val2",
		TypeName:    "string",
		IsUndefined: false,
	})

	log := logger.NewNilLogger(true)

	result, err := ex.Execute(variable.New(), log)

	assert.NoError(t, err)
	assert.Equal(t, expectedVariables, result.Variables)
}
