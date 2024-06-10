package executable_test

import (
	"cotton/internal/capture"
	"cotton/internal/config"
	"cotton/internal/executable"
	"cotton/internal/line"
	"cotton/internal/request"
	"net/http"
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
	req, _ := httpreqparser.New().Parse(`POST /some-path HTTP/1.0
Host: url

post
body`)
	expectedRequest, _ := request.New(req)

	reader := new(MockFileReader)
	reader.On("Read", "mock_file").Return(lines, nil)

	reqParser := new(MockRequestParser)
	reqParser.On("Parse", "POST /some-path HTTP/1.0\nHost: url\n\npost\nbody").Return(req, nil)

	parser := executable.NewParser(config, reader, reqParser)
	result, err := parser.FromMarkdownFile("mock_file")

	assert.NoError(t, err)
	assert.Equal(t, &executable.Executable{
		Request: expectedRequest,
	}, result)
	reader.AssertExpectations(t)
	reqParser.AssertExpectations(t)
}

func TestDiscardHTTPRequestInNonHTTPCodeBlockWillCauseANilExecutable(t *testing.T) {
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

	assert.Error(t, err)
	assert.Nil(t, result)
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

	req, _ := httpreqparser.New().Parse(`POST /some-path HTTP/1.0
Host: url

post
body`)
	expectedRequest, _ := request.New(req)

	reader := new(MockFileReader)
	reader.On("Read", "mock_file").Return(lines, nil)

	reqParser := new(MockRequestParser)
	reqParser.On("Parse", "POST /some-path HTTP/1.0\nHost: url\n\npost\nbody").Return(req, nil)

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
	req, _ := httpreqparser.New().Parse(`POST /this-should-be-collected HTTP/1.0
Host: url

post
body`)
	expectedRequest, _ := request.New(req)

	reader := new(MockFileReader)
	reader.On("Read", "mock_file").Return(lines, nil)

	reqParser := new(MockRequestParser)
	reqParser.On("Parse", "POST /this-should-be-collected HTTP/1.0\nHost: url\n\npost\nbody").Return(req, nil)

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
		"* var:sample",
		"* var2:`good.vibe`",
	}

	req, _ := httpreqparser.New().Parse(`POST /this-should-be-collected HTTP/1.0
Host: url

post
body`)
	expectedRequest, _ := request.New(req)
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
	reqParser.On("Parse", "POST /this-should-be-collected HTTP/1.0\nHost: url\n\npost\nbody").Return(req, nil)

	parser := executable.NewParser(config, reader, reqParser)
	result, err := parser.FromMarkdownFile("mock_file")

	assert.NoError(t, err)
	assert.Equal(t, &executable.Executable{
		Request:  expectedRequest,
		Captures: expectedCaptures,
	}, result)
}
