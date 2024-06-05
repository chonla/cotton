package executable_test

import (
	"bufio"
	"cotton/internal/capture"
	"cotton/internal/config"
	"cotton/internal/executable"
	"cotton/internal/line"
	"net/http"
	"strings"
	"testing"

	"github.com/chonla/httpreqparser"
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

type MockRequestParser struct {
	mock.Mock
}

func (m *MockRequestParser) Parse(req string) (*http.Request, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Request), args.Error(1)
}

func TestParseMarkdownFileShouldReadFromGivenFile(t *testing.T) {
	config := &config.Config{
		RootDir: "",
	}

	lines := []line.Line{""}

	reader := new(MockFileReader)
	reader.On("Read", "mock_file").Return(lines, nil)

	reqParser := new(MockRequestParser)

	parser := executable.NewParser(config, reader, reqParser)

	parser.FromMarkdownFile("mock_file")

	reader.AssertExpectations(t)
	reqParser.AssertNotCalled(t, "Parse")
}

func TestGetHTTPRequestInHTTPCodeBlock(t *testing.T) {
	config := &config.Config{
		RootDir: "",
	}

	lines := []line.Line{
		"```http",
		"POST /some-path HTTP/1.0",
		"Host: url",
		"",
		"post",
		"body",
		"```",
	}
	expectedRequest, _ := http.ReadRequest(bufio.NewReader(strings.NewReader(`POST /some-path HTTP/1.0
Host: url

post
body`)))

	reader := new(MockFileReader)
	reader.On("Read", "mock_file").Return(lines, nil)

	reqParser := new(MockRequestParser)
	reqParser.On("Parse", "POST /some-path HTTP/1.0\nHost: url\n\npost\nbody").Return(expectedRequest, nil)

	parser := executable.NewParser(config, reader, reqParser)
	result, err := parser.FromMarkdownFile("mock_file")

	assert.NoError(t, err)
	assert.Equal(t, &executable.Executable{
		Request: expectedRequest,
	}, result)
	reader.AssertExpectations(t)
	reqParser.AssertExpectations(t)
}

func TestDiscardHTTPRequestInNonHTTPCodeBlock(t *testing.T) {
	config := &config.Config{
		RootDir: "",
	}

	lines := []line.Line{
		"```",
		"POST /some-path HTTP/1.0",
		"Host: url",
		"",
		"post",
		"body",
		"```",
	}

	reader := new(MockFileReader)
	reader.On("Read", "mock_file").Return(lines, nil)

	reqParser := new(MockRequestParser)

	parser := executable.NewParser(config, reader, reqParser)
	result, err := parser.FromMarkdownFile("mock_file")

	assert.NoError(t, err)
	assert.Equal(t, &executable.Executable{
		Request: nil,
	}, result)
	reader.AssertExpectations(t)
	reqParser.AssertNotCalled(t, "Parse")
}

func TestGetHTTPRequestInOtherHTTPCodeBlock(t *testing.T) {
	config := &config.Config{
		RootDir: "",
	}

	lines := []line.Line{
		"```",
		"POST /this-should-be-ignored HTTP/1.0",
		"Host: url",
		"",
		"post",
		"body",
		"```",
		"",
		"```http",
		"POST /some-path HTTP/1.0",
		"Host: url",
		"",
		"post",
		"body",
		"```",
	}

	expectedRequest, _ := httpreqparser.New().Parse(`POST /some-path HTTP/1.0
Host: url

post
body`)

	reader := new(MockFileReader)
	reader.On("Read", "mock_file").Return(lines, nil)

	reqParser := new(MockRequestParser)
	reqParser.On("Parse", "POST /some-path HTTP/1.0\nHost: url\n\npost\nbody").Return(expectedRequest, nil)

	parser := executable.NewParser(config, reader, reqParser)
	result, err := parser.FromMarkdownFile("mock_file")

	assert.NoError(t, err)
	assert.Equal(t, &executable.Executable{
		Request: expectedRequest,
	}, result)
	reader.AssertExpectations(t)
	reqParser.AssertExpectations(t)
}

func TestGetHTTPRequestInOnlyFirstHTTPCodeBlock(t *testing.T) {
	config := &config.Config{
		RootDir: "",
	}

	lines := []line.Line{
		"```http",
		"POST /this-should-be-collected HTTP/1.0",
		"Host: url",
		"",
		"post",
		"body",
		"```",
		"",
		"```http",
		"POST /some-path HTTP/1.0",
		"Host: url",
		"",
		"post",
		"body",
		"```",
	}
	expectedRequest, _ := httpreqparser.New().Parse(`POST /this-should-be-collected HTTP/1.0
Host: url

post
body`)

	reader := new(MockFileReader)
	reader.On("Read", "mock_file").Return(lines, nil)

	reqParser := new(MockRequestParser)
	reqParser.On("Parse", "POST /this-should-be-collected HTTP/1.0\nHost: url\n\npost\nbody").Return(expectedRequest, nil)

	parser := executable.NewParser(config, reader, reqParser)
	result, err := parser.FromMarkdownFile("mock_file")

	assert.NoError(t, err)
	assert.Equal(t, &executable.Executable{
		Request: expectedRequest,
	}, result)
	reader.AssertExpectations(t)
	reqParser.AssertExpectations(t)
}

func TestGetCaptures(t *testing.T) {
	config := &config.Config{
		RootDir: "",
	}

	lines := []line.Line{
		"```http",
		"POST /this-should-be-collected HTTP/1.0",
		"Host: url",
		"",
		"post",
		"body",
		"```",
		"",
		"* var=sample",
		"* var2=`good.vibe`",
	}

	expectedRequest, _ := httpreqparser.New().Parse(`POST /this-should-be-collected HTTP/1.0
Host: url

post
body`)
	expectedCaptures := []*capture.Capture{
		{
			Name:    "var",
			Locator: "sample",
		},
		{
			Name:    "var2",
			Locator: "good.vibe",
		},
	}

	reader := new(MockFileReader)
	reader.On("Read", "mock_file").Return(lines, nil)

	reqParser := new(MockRequestParser)
	reqParser.On("Parse", "POST /this-should-be-collected HTTP/1.0\nHost: url\n\npost\nbody").Return(expectedRequest, nil)

	parser := executable.NewParser(config, reader, reqParser)
	result, err := parser.FromMarkdownFile("mock_file")

	assert.NoError(t, err)
	assert.Equal(t, &executable.Executable{
		Request:  expectedRequest,
		Captures: expectedCaptures,
	}, result)
}
