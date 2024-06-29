package executable_test

import (
	"cotton/internal/capture"
	"cotton/internal/executable"
	"cotton/internal/httphelper"
	"cotton/internal/logger"
	"cotton/internal/value"
	"cotton/internal/variable"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewReturnEmptyExecutable(t *testing.T) {
	mockReqParser := new(httphelper.MockHTTPRequestParser)
	mockLogger := new(logger.MockLogger)

	options := &executable.ExecutableOptions{
		RequestParser: mockReqParser,
		Logger:        mockLogger,
	}

	result := executable.New("title", "req", options)

	assert.Equal(t, "title", result.Title())
	assert.Equal(t, "req", result.RawRequest())
	assert.Equal(t, []*capture.Capture{}, result.Captures())
}

func TestAddCaptures(t *testing.T) {
	mockReqParser := new(httphelper.MockHTTPRequestParser)
	mockLogger := new(logger.MockLogger)

	options := &executable.ExecutableOptions{
		RequestParser: mockReqParser,
		Logger:        mockLogger,
	}

	ex := executable.New("title", "req", options)
	ex.AddCapture(capture.New("key", "value"))

	result := ex.Captures()

	assert.Equal(t, 1, len(result))
	assert.Equal(t, "key", result[0].Name)
	assert.Equal(t, "value", result[0].Selector)
}

func TestEditCapturesShouldNotMutateTheOriginalCaptures(t *testing.T) {
	mockReqParser := new(httphelper.MockHTTPRequestParser)
	mockLogger := new(logger.MockLogger)

	options := &executable.ExecutableOptions{
		RequestParser: mockReqParser,
		Logger:        mockLogger,
	}

	ex := executable.New("title", "req", options)
	ex.AddCapture(capture.New("key", "value"))

	captures := ex.Captures()
	captures[0].Name = "key1"
	captures[0].Selector = "value1"

	result := ex.Captures()

	assert.NotEqual(t, result[0].Name, captures[0].Name)
	assert.NotEqual(t, result[0].Selector, captures[0].Selector)
}

func TestParsingCaptureFromResponse(t *testing.T) {
	reqRaw := "GET http://localhost/path HTTP/1.1"
	respRaw := `HTTP/1.1 200 OK
Content-Length: 38
Content-Type: application/json

{"form":{"key1":"val1","key2":"val2"}}`

	mockHTTPResponse, _ := (&httphelper.HTTPResponseParser{}).Parse(respRaw)

	mockHTTPRequest := new(httphelper.MockHTTPRequest)
	mockHTTPRequest.On("Do", mock.AnythingOfType("bool")).Return(mockHTTPResponse, nil)

	mockReqParser := new(httphelper.MockHTTPRequestParser)
	mockReqParser.On("Parse", reqRaw).Return(mockHTTPRequest, nil)

	expectedCapturedVariables := variable.New()
	expectedCapturedVariables.Set("var1", value.New("val1"))
	expectedCapturedVariables.Set("var2", value.New("val2"))

	mockLogger := new(logger.MockLogger)
	mockLogger.On("PrintExecutableTitle", mock.Anything).Return(nil)
	mockLogger.On("PrintRequest", mock.Anything).Return(nil)
	mockLogger.On("PrintVariables", mock.Anything).Return(nil)

	options := &executable.ExecutableOptions{
		RequestParser: mockReqParser,
		Logger:        mockLogger,
		// ClockWrapper:  mockClock,
	}

	ex := executable.New("test", reqRaw, options)
	ex.AddCapture(capture.New("var1", "Body.form.key1"))
	ex.AddCapture(capture.New("var2", "Body.form.key2"))

	initialVars := variable.New()
	result, err := ex.Execute(initialVars)

	assert.NoError(t, err)
	assert.Equal(t, expectedCapturedVariables, result.Variables)
}
