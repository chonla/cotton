package testcase_test

import (
	"cotton/internal/assertion"
	"cotton/internal/capture"
	"cotton/internal/config"
	"cotton/internal/executable"
	"cotton/internal/line"
	"cotton/internal/request"
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

// func TestParseMarkdownFileShouldReadFromGivenFile(t *testing.T) {
// 	config := &config.Config{
// 		RootDir: "",
// 	}

// 	lines := []line.Line{""}

// 	reader := new(MockFileReader)
// 	reader.On("Read", "mock_file").Return(lines, nil)

// 	reqParser := new(MockRequestParser)
// 	reqParser.On("Parse", "POST /some-path HTTP/1.0\nHost: url\n\npost\nbody").Return(expectedRequest, nil)

// 	parser := testcase.NewParser(config, reader, reqParser)
// 	parser.FromMarkdownFile("mock_file")

// 	reader.AssertExpectations(t)
// 	reqParser.AssertExpectations(t)
// }

// First H1 is considered a test title
func TestH1AtVeryFirstLine(t *testing.T) {
	config := &config.Config{
		RootDir: "",
	}

	lines := []line.Line{
		"# Title",
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

	parser := testcase.NewParser(config, reader, reqParser)
	result, err := parser.FromMarkdownFile("mock_file")

	assert.Nil(t, err)
	assert.Equal(t, &testcase.TestCase{
		Title:   "Title",
		Request: expectedRequest,
	}, result)
	reader.AssertExpectations(t)
}

func TestH1AtSomeOtherLines(t *testing.T) {
	config := &config.Config{
		RootDir: "",
	}

	lines := []line.Line{
		"## Ok",
		"",
		"# Title",
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

	parser := testcase.NewParser(config, reader, reqParser)
	result, err := parser.FromMarkdownFile("mock_file")

	assert.NoError(t, err)
	assert.Equal(t, &testcase.TestCase{
		Title:   "Title",
		Request: expectedRequest,
	}, result)
	reader.AssertExpectations(t)
	reqParser.AssertExpectations(t)
}

func TestMultipleH1sWillGrabTheFirstH1AsTitle(t *testing.T) {
	config := &config.Config{
		RootDir: "",
	}

	lines := []line.Line{
		"# Title",
		"",
		"```http",
		"POST /some-path HTTP/1.0",
		"Host: url",
		"",
		"post",
		"body",
		"```",
		"",
		"# Other Title",
	}

	req, _ := httpreqparser.New().Parse(`POST /some-path HTTP/1.0
Host: url

post
body`)
	expectedReq, _ := request.New(req)

	reader := new(MockFileReader)
	reader.On("Read", "mock_file").Return(lines, nil)

	reqParser := new(MockRequestParser)
	reqParser.On("Parse", "POST /some-path HTTP/1.0\nHost: url\n\npost\nbody").Return(req, nil)

	parser := testcase.NewParser(config, reader, reqParser)
	result, err := parser.FromMarkdownFile("mock_file")

	assert.NoError(t, err)
	assert.Equal(t, &testcase.TestCase{
		Title:   "Title",
		Request: expectedReq,
	}, result)
	reader.AssertExpectations(t)
	reqParser.AssertExpectations(t)
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
	expectedReq, _ := request.New(req)

	reader := new(MockFileReader)
	reader.On("Read", "mock_file").Return(lines, nil)

	reqParser := new(MockRequestParser)
	reqParser.On("Parse", "POST /some-path HTTP/1.0\nHost: url\n\npost\nbody").Return(req, nil)

	parser := testcase.NewParser(config, reader, reqParser)
	result, err := parser.FromMarkdownFile("mock_file")

	assert.NoError(t, err)
	assert.Equal(t, &testcase.TestCase{
		Title:       "Title",
		Description: "Wonderful\nworld",
		Request:     expectedReq,
	}, result)
	reader.AssertExpectations(t)
	reqParser.AssertExpectations(t)
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
	expectedReq, _ := request.New(req)

	reader := new(MockFileReader)
	reader.On("Read", "mock_file").Return(lines, nil)

	reqParser := new(MockRequestParser)
	reqParser.On("Parse", "POST /some-path HTTP/1.0\nHost: url\n\npost\nbody").Return(req, nil)

	parser := testcase.NewParser(config, reader, reqParser)
	result, err := parser.FromMarkdownFile("mock_file")

	assert.NoError(t, err)
	assert.Equal(t, &testcase.TestCase{
		Title:       "Title",
		Description: "Wonderful",
		Request:     expectedReq,
	}, result)
	reader.AssertExpectations(t)
	reqParser.AssertExpectations(t)
}

func TestGetHTTPRequestWithoutTitle(t *testing.T) {
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

	parser := testcase.NewParser(config, reader, reqParser)
	result, err := parser.FromMarkdownFile("mock_file")

	assert.NoError(t, err)
	assert.Equal(t, &testcase.TestCase{
		Request: expectedRequest,
	}, result)
	reader.AssertExpectations(t)
	reqParser.AssertExpectations(t)
}

func TestGetHTTPRequestFromThreeTildedCodeBlock(t *testing.T) {
	config := &config.Config{
		RootDir: "",
	}

	lines := []line.Line{
		"~~~http",
		"POST /some-path HTTP/1.0",
		"Host: url",
		"",
		"post",
		"body",
		"~~~",
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

	parser := testcase.NewParser(config, reader, reqParser)
	result, err := parser.FromMarkdownFile("mock_file")

	assert.NoError(t, err)
	assert.Equal(t, &testcase.TestCase{
		Request: expectedRequest,
	}, result)
	reader.AssertExpectations(t)
	reqParser.AssertExpectations(t)
}

func TestGetHTTPRequestInMixedHTTPCodeBlock(t *testing.T) {
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
		"~~~http",
		"POST /some-path HTTP/1.0",
		"Host: url",
		"",
		"post",
		"body",
		"~~~",
		"",
		"```http",
		"POST /this-should-be-ignored-too HTTP/1.0",
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
	req, _ := httpreqparser.New().Parse(`POST /this-should-be-collected HTTP/1.0
Host: url

post
body`)
	expectedRequest, _ := request.New(req)

	reader := new(MockFileReader)
	reader.On("Read", "mock_file").Return(lines, nil)

	reqParser := new(MockRequestParser)
	reqParser.On("Parse", "POST /this-should-be-collected HTTP/1.0\nHost: url\n\npost\nbody").Return(req, nil)

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
		"* readiness:`$.readiness.status`",
	}

	reader := new(MockFileReader)
	reader.On("Read", "mock_file").Return(lines, nil)
	reader.On("Read", "link1").Return(beforeLink1Lines, nil)
	reader.On("Read", "link2").Return(beforeLink2Lines, nil)

	beforeReq1, _ := httpreqparser.New().Parse(`GET /healthcheck HTTP/1.0
Host: localhost`)
	expectedBefore1, _ := request.New(beforeReq1)
	beforeReq2, _ := httpreqparser.New().Parse(`GET /readiness HTTP/1.0
Host: localhost`)
	expectedBefore2, _ := request.New(beforeReq2)
	req, _ := httpreqparser.New().Parse(`POST /some-path HTTP/1.0
Host: url

post
body`)
	expectedRequest, _ := request.New(req)

	reqParser := new(MockRequestParser)
	reqParser.On("Parse", "POST /some-path HTTP/1.0\nHost: url\n\npost\nbody").Return(req, nil)
	reqParser.On("Parse", "GET /healthcheck HTTP/1.0\nHost: localhost").Return(beforeReq1, nil)
	reqParser.On("Parse", "GET /readiness HTTP/1.0\nHost: localhost").Return(beforeReq2, nil)

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
		"* varname:`$.result`",
	}

	req, _ := httpreqparser.New().Parse(`POST /this-should-be-collected HTTP/1.0
Host: localhost

post
body`)
	expectedRequest, _ := request.New(req)

	reader := new(MockFileReader)
	reader.On("Read", "mock_file").Return(lines, nil)

	reqParser := new(MockRequestParser)
	reqParser.On("Parse", "POST /this-should-be-collected HTTP/1.0\nHost: localhost\n\npost\nbody").Return(req, nil)

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

func TestGetAssertion(t *testing.T) {
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
		"* `$.var`==`3`",
		"* `$.var2`==`\"good.vibe\"`",
	}

	req, _ := httpreqparser.New().Parse(`POST /this-should-be-collected HTTP/1.0
Host: url

post
body`)
	expectedRequest, _ := request.New(req)

	expectedAssertions := []*assertion.Assertion{
		{
			Selector: "$.var",
			Value:    float64(3),
			Operator: &assertion.EqAssertion{},
		},
		{
			Selector: "$.var2",
			Value:    "good.vibe",
			Operator: &assertion.EqAssertion{},
		},
	}

	reader := new(MockFileReader)
	reader.On("Read", "mock_file").Return(lines, nil)

	reqParser := new(MockRequestParser)
	reqParser.On("Parse", "POST /this-should-be-collected HTTP/1.0\nHost: url\n\npost\nbody").Return(req, nil)

	parser := testcase.NewParser(config, reader, reqParser)
	result, err := parser.FromMarkdownFile("mock_file")

	assert.NoError(t, err)
	assert.Equal(t, &testcase.TestCase{
		Request:    expectedRequest,
		Assertions: expectedAssertions,
	}, result)
}
