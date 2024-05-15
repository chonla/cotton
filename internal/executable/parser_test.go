package executable_test

import (
	"bufio"
	"cotton/internal/capture"
	"cotton/internal/executable"
	"cotton/internal/line"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockFileReader struct {
	mock.Mock
}

func (m *MockFileReader) Read(fileName string) ([]line.Line, error) {
	args := m.Called(fileName)
	return args.Get(0).([]line.Line), args.Error(1)
}

func TestParseMarkdownFileShouldReadFromGivenFile(t *testing.T) {
	lines := []line.Line{""}

	reader := new(MockFileReader)
	reader.On("Read", "mock_file").Return(lines, nil)

	parser := executable.NewParser(reader)
	parser.FromMarkdownFile("mock_file")

	reader.AssertExpectations(t)
}

func TestGetHTTPRequestInHTTPCodeBlock(t *testing.T) {
	lines := []line.Line{
		"```http",
		"POST /some-path HTTP/1.0",
		"Host: http://url",
		"",
		"post",
		"body",
		"```",
	}

	reader := new(MockFileReader)
	reader.On("Read", "mock_file").Return(lines, nil)

	parser := executable.NewParser(reader)
	result, err := parser.FromMarkdownFile("mock_file")

	expectedRequest, _ := http.ReadRequest(bufio.NewReader(strings.NewReader(`POST /some-path HTTP/1.0
Host: http://url

post
body`)))

	assert.NoError(t, err)
	assert.Equal(t, &executable.Executable{
		Request: expectedRequest,
	}, result)

}

func TestDiscardHTTPRequestInNonHTTPCodeBlock(t *testing.T) {
	lines := []line.Line{
		"```",
		"POST /some-path HTTP/1.0",
		"Host: http://url",
		"",
		"post",
		"body",
		"```",
	}

	reader := new(MockFileReader)
	reader.On("Read", "mock_file").Return(lines, nil)

	parser := executable.NewParser(reader)
	result, err := parser.FromMarkdownFile("mock_file")

	assert.NoError(t, err)
	assert.Equal(t, &executable.Executable{
		Request: nil,
	}, result)

}

func TestGetHTTPRequestInOtherHTTPCodeBlock(t *testing.T) {
	lines := []line.Line{
		"```",
		"POST /this-should-be-ignored HTTP/1.0",
		"Host: http://url",
		"",
		"post",
		"body",
		"```",
		"",
		"```http",
		"POST /some-path HTTP/1.0",
		"Host: http://url",
		"",
		"post",
		"body",
		"```",
	}

	reader := new(MockFileReader)
	reader.On("Read", "mock_file").Return(lines, nil)

	parser := executable.NewParser(reader)
	result, err := parser.FromMarkdownFile("mock_file")

	expectedRequest, _ := http.ReadRequest(bufio.NewReader(strings.NewReader(`POST /some-path HTTP/1.0
Host: http://url

post
body`)))

	assert.NoError(t, err)
	assert.Equal(t, &executable.Executable{
		Request: expectedRequest,
	}, result)

}

func TestGetHTTPRequestInOnlyFirstHTTPCodeBlock(t *testing.T) {
	lines := []line.Line{
		"```http",
		"POST /this-should-be-collected HTTP/1.0",
		"Host: http://url",
		"",
		"post",
		"body",
		"```",
		"",
		"```http",
		"POST /some-path HTTP/1.0",
		"Host: http://url",
		"",
		"post",
		"body",
		"```",
	}

	reader := new(MockFileReader)
	reader.On("Read", "mock_file").Return(lines, nil)

	parser := executable.NewParser(reader)
	result, err := parser.FromMarkdownFile("mock_file")

	expectedRequest, _ := http.ReadRequest(bufio.NewReader(strings.NewReader(`POST /this-should-be-collected HTTP/1.0
Host: http://url

post
body`)))

	assert.NoError(t, err)
	assert.Equal(t, &executable.Executable{
		Request: expectedRequest,
	}, result)
}

func TestGetCaptures(t *testing.T) {
	lines := []line.Line{
		"```http",
		"POST /this-should-be-collected HTTP/1.0",
		"Host: http://url",
		"",
		"post",
		"body",
		"```",
		"",
		"* var=sample",
		"* var2=`good.vibe`",
	}

	reader := new(MockFileReader)
	reader.On("Read", "mock_file").Return(lines, nil)

	parser := executable.NewParser(reader)
	result, err := parser.FromMarkdownFile("mock_file")

	expectedRequest, _ := http.ReadRequest(bufio.NewReader(strings.NewReader(`POST /this-should-be-collected HTTP/1.0
Host: http://url

post
body`)))
	expectedCaptures := []*capture.Captured{
		{
			Name:    "var",
			Locator: "sample",
		},
		{
			Name:    "var2",
			Locator: "good.vibe",
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, &executable.Executable{
		Request:  expectedRequest,
		Captures: expectedCaptures,
	}, result)
}
