package testcase_test

import (
	"bufio"
	"cotton/internal/capture"
	"cotton/internal/executable"
	"cotton/internal/line"
	"cotton/internal/testcase"
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

	parser := testcase.NewParser(reader)
	parser.FromMarkdownFile("mock_file")

	reader.AssertExpectations(t)
}

// First H1 is considered a test title
func TestH1AtVeryFirstLine(t *testing.T) {
	lines := []line.Line{"# Title"}

	reader := new(MockFileReader)
	reader.On("Read", "mock_file").Return(lines, nil)

	parser := testcase.NewParser(reader)
	result, err := parser.FromMarkdownFile("mock_file")

	assert.NoError(t, err)
	assert.Equal(t, &testcase.TestCase{
		Title:   "Title",
		Request: nil,
	}, result)
}

func TestH1AtSomeOtherLines(t *testing.T) {
	lines := []line.Line{
		"## Ok",
		"",
		"# Title",
	}

	reader := new(MockFileReader)
	reader.On("Read", "mock_file").Return(lines, nil)

	parser := testcase.NewParser(reader)
	result, err := parser.FromMarkdownFile("mock_file")

	assert.NoError(t, err)
	assert.Equal(t, &testcase.TestCase{
		Title:   "Title",
		Request: nil,
	}, result)
}

func TestMultipleH1s(t *testing.T) {
	lines := []line.Line{
		"# Title",
		"# Other Title",
	}

	reader := new(MockFileReader)
	reader.On("Read", "mock_file").Return(lines, nil)

	parser := testcase.NewParser(reader)
	result, err := parser.FromMarkdownFile("mock_file")

	assert.NoError(t, err)
	assert.Equal(t, &testcase.TestCase{
		Title:   "Title",
		Request: nil,
	}, result)
}

func TestDescription(t *testing.T) {
	lines := []line.Line{
		"# Title",
		"",
		"Wonderful",
		"world",
	}

	reader := new(MockFileReader)
	reader.On("Read", "mock_file").Return(lines, nil)

	parser := testcase.NewParser(reader)
	result, err := parser.FromMarkdownFile("mock_file")

	assert.NoError(t, err)
	assert.Equal(t, &testcase.TestCase{
		Title:       "Title",
		Description: "Wonderful\nworld",
		Request:     nil,
	}, result)
}

func TestDescriptionAtOtherH1(t *testing.T) {
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

	parser := testcase.NewParser(reader)
	result, err := parser.FromMarkdownFile("mock_file")

	assert.NoError(t, err)
	assert.Equal(t, &testcase.TestCase{
		Title:       "Title",
		Description: "Wonderful",
		Request:     nil,
	}, result)
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

	parser := testcase.NewParser(reader)
	result, err := parser.FromMarkdownFile("mock_file")

	expectedRequest, _ := http.ReadRequest(bufio.NewReader(strings.NewReader(`POST /some-path HTTP/1.0
Host: http://url

post
body`)))

	assert.NoError(t, err)
	assert.Equal(t, &testcase.TestCase{
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

	parser := testcase.NewParser(reader)
	result, err := parser.FromMarkdownFile("mock_file")

	assert.NoError(t, err)
	assert.Equal(t, &testcase.TestCase{
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

	parser := testcase.NewParser(reader)
	result, err := parser.FromMarkdownFile("mock_file")

	expectedRequest, _ := http.ReadRequest(bufio.NewReader(strings.NewReader(`POST /some-path HTTP/1.0
Host: http://url

post
body`)))

	assert.NoError(t, err)
	assert.Equal(t, &testcase.TestCase{
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

	parser := testcase.NewParser(reader)
	result, err := parser.FromMarkdownFile("mock_file")

	expectedRequest, _ := http.ReadRequest(bufio.NewReader(strings.NewReader(`POST /this-should-be-collected HTTP/1.0
Host: http://url

post
body`)))

	assert.NoError(t, err)
	assert.Equal(t, &testcase.TestCase{
		Request: expectedRequest,
	}, result)
}

func TestGetExecutablesBeforeTest(t *testing.T) {
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
		"Host: http://url",
		"",
		"post",
		"body",
		"```",
	}

	beforeLink1Lines := []line.Line{
		"This is before first step",
		"```http",
		"GET /healthcheck HTTP/1.0",
		"Host: http://localhost",
		"```",
	}

	beforeLink2Lines := []line.Line{
		"This is before second step",
		"```http",
		"GET /readiness HTTP/1.0",
		"Host: http://localhost",
		"```",
		"# Capture part",
		"* readiness=`$.readiness.status`",
	}

	reader := new(MockFileReader)
	reader.On("Read", "mock_file").Return(lines, nil)
	reader.On("Read", "link1").Return(beforeLink1Lines, nil)
	reader.On("Read", "link2").Return(beforeLink2Lines, nil)

	parser := testcase.NewParser(reader)
	result, err := parser.FromMarkdownFile("mock_file")

	expectedBefore, _ := http.ReadRequest(bufio.NewReader(strings.NewReader(`POST /this-should-be-collected-in-before HTTP/1.0
	Host: http://url
	
	post
	body`)))
	expectedRequest, _ := http.ReadRequest(bufio.NewReader(strings.NewReader(`POST /some-path HTTP/1.0
Host: http://url

post
body`)))

	assert.NoError(t, err)
	assert.Equal(t, &testcase.TestCase{
		Request: expectedRequest,
		Setups: []*executable.Executable{
			{
				Title:   "link to executable",
				Request: expectedBefore,
			},
			{
				Title:   "link to another executable",
				Request: expectedBefore,
				Captures: []*capture.Captured{
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
	lines := []line.Line{
		"```http",
		"POST /this-should-be-collected HTTP/1.0",
		"Host: http://url",
		"",
		"post",
		"body",
		"```",
		"",
		"* varname=`$.result`",
	}

	reader := new(MockFileReader)
	reader.On("Read", "mock_file").Return(lines, nil)

	parser := testcase.NewParser(reader)
	result, err := parser.FromMarkdownFile("mock_file")

	expectedRequest, _ := http.ReadRequest(bufio.NewReader(strings.NewReader(`POST /this-should-be-collected HTTP/1.0
Host: http://url

post
body`)))

	assert.NoError(t, err)
	assert.Equal(t, &testcase.TestCase{
		Request: expectedRequest,
		Captures: []*capture.Captured{
			{
				Name:    "varname",
				Locator: "$.result",
			},
		},
	}, result)

}
