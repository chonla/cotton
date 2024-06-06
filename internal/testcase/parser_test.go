package testcase_test

import (
	"cotton/internal/capture"
	"cotton/internal/config"
	"cotton/internal/executable"
	"cotton/internal/line"
	"cotton/internal/testcase"
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

	parser := testcase.NewParser(config, reader, reqParser)
	parser.FromMarkdownFile("mock_file")

	reader.AssertExpectations(t)
	reqParser.AssertNotCalled(t, "Parse")
}

// First H1 is considered a test title
func TestH1AtVeryFirstLine(t *testing.T) {
	config := &config.Config{
		RootDir: "",
	}

	lines := []line.Line{"# Title"}

	reader := new(MockFileReader)
	reader.On("Read", "mock_file").Return(lines, nil)

	reqParser := new(MockRequestParser)

	parser := testcase.NewParser(config, reader, reqParser)
	result, err := parser.FromMarkdownFile("mock_file")

	assert.NoError(t, err)
	assert.Equal(t, &testcase.TestCase{
		Title:   "Title",
		Request: nil,
	}, result)
	reader.AssertExpectations(t)
	reqParser.AssertNotCalled(t, "Parse")
}

func TestH1AtSomeOtherLines(t *testing.T) {
	config := &config.Config{
		RootDir: "",
	}

	lines := []line.Line{
		"## Ok",
		"",
		"# Title",
	}

	reader := new(MockFileReader)
	reader.On("Read", "mock_file").Return(lines, nil)

	reqParser := new(MockRequestParser)

	parser := testcase.NewParser(config, reader, reqParser)
	result, err := parser.FromMarkdownFile("mock_file")

	assert.NoError(t, err)
	assert.Equal(t, &testcase.TestCase{
		Title:   "Title",
		Request: nil,
	}, result)
	reader.AssertExpectations(t)
	reqParser.AssertNotCalled(t, "Parse")
}

func TestMultipleH1s(t *testing.T) {
	config := &config.Config{
		RootDir: "",
	}

	lines := []line.Line{
		"# Title",
		"# Other Title",
	}

	reader := new(MockFileReader)
	reader.On("Read", "mock_file").Return(lines, nil)

	reqParser := new(MockRequestParser)

	parser := testcase.NewParser(config, reader, reqParser)
	result, err := parser.FromMarkdownFile("mock_file")

	assert.NoError(t, err)
	assert.Equal(t, &testcase.TestCase{
		Title:   "Title",
		Request: nil,
	}, result)
	reader.AssertExpectations(t)
	reqParser.AssertNotCalled(t, "Parse")
}

func TestDescription(t *testing.T) {
	config := &config.Config{
		RootDir: "",
	}

	lines := []line.Line{
		"# Title",
		"",
		"Wonderful",
		"world",
	}

	reader := new(MockFileReader)
	reader.On("Read", "mock_file").Return(lines, nil)

	reqParser := new(MockRequestParser)

	parser := testcase.NewParser(config, reader, reqParser)
	result, err := parser.FromMarkdownFile("mock_file")

	assert.NoError(t, err)
	assert.Equal(t, &testcase.TestCase{
		Title:       "Title",
		Description: "Wonderful\nworld",
		Request:     nil,
	}, result)
	reader.AssertExpectations(t)
	reqParser.AssertNotCalled(t, "Parse")
}

func TestDescriptionAtOtherH1(t *testing.T) {
	config := &config.Config{
		RootDir: "",
	}

	lines := []line.Line{
		"# Title",
		"",
		"Wonderful",
		"",
		"# Title 2",
		"",
		"Wonderful world",
		"Lalala",
	}

	reader := new(MockFileReader)
	reader.On("Read", "mock_file").Return(lines, nil)

	reqParser := new(MockRequestParser)

	parser := testcase.NewParser(config, reader, reqParser)
	result, err := parser.FromMarkdownFile("mock_file")

	assert.NoError(t, err)
	assert.Equal(t, &testcase.TestCase{
		Title:       "Title",
		Description: "Wonderful",
		Request:     nil,
	}, result)
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

	expectedRequest, _ := httpreqparser.New().Parse(`POST /some-path HTTP/1.0
Host: url

post
body`)

	reader := new(MockFileReader)
	reader.On("Read", "mock_file").Return(lines, nil)

	reqParser := new(MockRequestParser)
	reqParser.On("Parse", "POST /some-path HTTP/1.0\nHost: url\n\npost\nbody").Return(expectedRequest, nil)

	parser := testcase.NewParser(config, reader, reqParser)
	result, err := parser.FromMarkdownFile("mock_file")

	assert.NoError(t, err)
	assert.Equal(t, &testcase.TestCase{
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

	parser := testcase.NewParser(config, reader, reqParser)
	result, err := parser.FromMarkdownFile("mock_file")

	assert.NoError(t, err)
	assert.Equal(t, &testcase.TestCase{
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

	parser := testcase.NewParser(config, reader, reqParser)
	result, err := parser.FromMarkdownFile("mock_file")

	assert.NoError(t, err)
	assert.Equal(t, &testcase.TestCase{
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

	parser := testcase.NewParser(config, reader, reqParser)
	result, err := parser.FromMarkdownFile("mock_file")

	assert.NoError(t, err)
	assert.Equal(t, &testcase.TestCase{
		Request: expectedRequest,
	}, result)
	reader.AssertExpectations(t)
	reqParser.AssertExpectations(t)
}

func TestGetExecutablesBeforeTest(t *testing.T) {
	config := &config.Config{
		RootDir: "",
	}

	lines := []line.Line{
		"## Before",
		"",
		"* [link to executable](link1)",
		"* [link to another executable](link2)",
		"",
		"## Method under test",
		"",
		"```http",
		"POST /some-path HTTP/1.0",
		"Host: url",
		"",
		"post",
		"body",
		"```",
	}

	beforeLink1Lines := []line.Line{
		"This is before first step",
		"```http",
		"GET /healthcheck HTTP/1.0",
		"Host: localhost",
		"```",
	}

	beforeLink2Lines := []line.Line{
		"This is before second step",
		"```http",
		"GET /readiness HTTP/1.0",
		"Host: localhost",
		"```",
		"# Capture part",
		"* readiness=`$.readiness.status`",
	}

	reader := new(MockFileReader)
	reader.On("Read", "mock_file").Return(lines, nil)
	reader.On("Read", "link1").Return(beforeLink1Lines, nil)
	reader.On("Read", "link2").Return(beforeLink2Lines, nil)

	expectedBefore1, _ := httpreqparser.New().Parse(`GET /healthcheck HTTP/1.0
Host: localhost`)
	expectedBefore2, _ := httpreqparser.New().Parse(`GET /readiness HTTP/1.0
Host: localhost`)
	expectedRequest, _ := httpreqparser.New().Parse(`POST /some-path HTTP/1.0
Host: url

post
body`)

	reqParser := new(MockRequestParser)
	reqParser.On("Parse", "POST /some-path HTTP/1.0\nHost: url\n\npost\nbody").Return(expectedRequest, nil)
	reqParser.On("Parse", "GET /healthcheck HTTP/1.0\nHost: localhost").Return(expectedBefore1, nil)
	reqParser.On("Parse", "GET /readiness HTTP/1.0\nHost: localhost").Return(expectedBefore2, nil)

	parser := testcase.NewParser(config, reader, reqParser)
	result, err := parser.FromMarkdownFile("mock_file")

	assert.NoError(t, err)
	assert.Equal(t, &testcase.TestCase{
		Request: expectedRequest,
		Setups: []*executable.Executable{
			{
				Title:   "link to executable",
				Request: expectedBefore1,
			},
			{
				Title:   "link to another executable",
				Request: expectedBefore2,
				Captures: []*capture.Capture{
					{
						Name:    "readiness",
						Locator: "$.readiness.status",
					},
				},
			},
		},
	}, result)
}

func TestGetCapturesInTestcase(t *testing.T) {
	config := &config.Config{
		RootDir: "",
	}

	lines := []line.Line{
		"```http",
		"POST /this-should-be-collected HTTP/1.0",
		"Host: localhost",
		"",
		"post",
		"body",
		"```",
		"",
		"* varname=`$.result`",
	}

	expectedRequest, _ := httpreqparser.New().Parse(`POST /this-should-be-collected HTTP/1.0
Host: localhost

post
body`)

	reader := new(MockFileReader)
	reader.On("Read", "mock_file").Return(lines, nil)

	reqParser := new(MockRequestParser)
	reqParser.On("Parse", "POST /this-should-be-collected HTTP/1.0\nHost: localhost\n\npost\nbody").Return(expectedRequest, nil)

	parser := testcase.NewParser(config, reader, reqParser)
	result, err := parser.FromMarkdownFile("mock_file")

	assert.NoError(t, err)
	assert.Equal(t, &testcase.TestCase{
		Request: expectedRequest,
		Captures: []*capture.Capture{
			{
				Name:    "varname",
				Locator: "$.result",
			},
		},
	}, result)

}
