package testcase_test

import (
	"cotton/internal/assertion"
	"cotton/internal/capture"
	"cotton/internal/config"
	"cotton/internal/executable"
	"cotton/internal/httphelper"
	"cotton/internal/line"
	"cotton/internal/logger"
	"cotton/internal/reader"
	"cotton/internal/testcase"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// First H1 is considered a test title
func TestH1AtVeryFirstLine(t *testing.T) {
	config := &config.Config{
		BaseDir: "",
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

	mockFileReader := new(reader.MockFileReader)
	mockFileReader.On("Read", "mock_file").Return(lines, nil)

	mockHTTPRequestParser := new(httphelper.MockHTTPRequestParser)

	mockLogger := new(logger.MockLogger)
	mockLogger.On("PrintDetailedDebugMessage", mock.Anything).Return(nil)

	mockExecutableParser := new(executable.MockExecutableParser)

	options := &testcase.ParserOptions{
		Configurator:     config,
		FileReader:       mockFileReader,
		RequestParser:    mockHTTPRequestParser,
		Logger:           mockLogger,
		ExecutableParser: mockExecutableParser,
	}

	testcaseOptions := &testcase.TestcaseOptions{
		Logger:        mockLogger,
		RequestParser: mockHTTPRequestParser,
	}
	expectedRawRequest := `POST /some-path HTTP/1.0
Host: url

post
body`
	expectedTestcase := testcase.NewTestcase("Title", "", expectedRawRequest, testcaseOptions)

	parser := testcase.NewParser(options)
	result, err := parser.FromMarkdownFile("mock_file")

	assert.Nil(t, err)
	assert.Equal(t, expectedTestcase, result)
	mockFileReader.AssertExpectations(t)
}

func TestH1AtSomeOtherLines(t *testing.T) {
	config := &config.Config{
		BaseDir: "",
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

	mockFileReader := new(reader.MockFileReader)
	mockFileReader.On("Read", "mock_file").Return(lines, nil)

	mockHTTPRequestParser := new(httphelper.MockHTTPRequestParser)

	mockLogger := new(logger.MockLogger)
	mockLogger.On("PrintDetailedDebugMessage", mock.Anything).Return(nil)

	mockExecutableParser := new(executable.MockExecutableParser)

	options := &testcase.ParserOptions{
		Configurator:     config,
		FileReader:       mockFileReader,
		RequestParser:    mockHTTPRequestParser,
		Logger:           mockLogger,
		ExecutableParser: mockExecutableParser,
	}

	testcaseOptions := &testcase.TestcaseOptions{
		Logger:        mockLogger,
		RequestParser: mockHTTPRequestParser,
	}
	expectedRawRequest := `POST /some-path HTTP/1.0
Host: url

post
body`
	expectedTestcase := testcase.NewTestcase("Title", "", expectedRawRequest, testcaseOptions)

	parser := testcase.NewParser(options)
	result, err := parser.FromMarkdownFile("mock_file")

	assert.NoError(t, err)
	assert.Equal(t, expectedTestcase, result)
	mockFileReader.AssertExpectations(t)
	mockHTTPRequestParser.AssertExpectations(t)
}

func TestMultipleH1sWillGrabTheFirstH1AsTitle(t *testing.T) {
	config := &config.Config{
		BaseDir: "",
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

	mockFileReader := new(reader.MockFileReader)
	mockFileReader.On("Read", "mock_file").Return(lines, nil)

	mockHTTPRequestParser := new(httphelper.MockHTTPRequestParser)

	mockLogger := new(logger.MockLogger)
	mockLogger.On("PrintDetailedDebugMessage", mock.Anything).Return(nil)

	mockExecutableParser := new(executable.MockExecutableParser)

	options := &testcase.ParserOptions{
		Configurator:     config,
		FileReader:       mockFileReader,
		RequestParser:    mockHTTPRequestParser,
		Logger:           mockLogger,
		ExecutableParser: mockExecutableParser,
	}

	testcaseOptions := &testcase.TestcaseOptions{
		Logger:        mockLogger,
		RequestParser: mockHTTPRequestParser,
	}
	expectedRawRequest := `POST /some-path HTTP/1.0
Host: url

post
body`
	expectedTestcase := testcase.NewTestcase("Title", "", expectedRawRequest, testcaseOptions)

	parser := testcase.NewParser(options)
	result, err := parser.FromMarkdownFile("mock_file")

	assert.NoError(t, err)
	assert.Equal(t, expectedTestcase, result)
	mockFileReader.AssertExpectations(t)
	mockHTTPRequestParser.AssertExpectations(t)
}

func TestDescription(t *testing.T) {
	config := &config.Config{
		BaseDir: "",
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

	mockFileReader := new(reader.MockFileReader)
	mockFileReader.On("Read", "mock_file").Return(lines, nil)

	mockHTTPRequestParser := new(httphelper.MockHTTPRequestParser)

	mockLogger := new(logger.MockLogger)
	mockLogger.On("PrintDetailedDebugMessage", mock.Anything).Return(nil)

	mockExecutableParser := new(executable.MockExecutableParser)

	options := &testcase.ParserOptions{
		Configurator:     config,
		FileReader:       mockFileReader,
		RequestParser:    mockHTTPRequestParser,
		Logger:           mockLogger,
		ExecutableParser: mockExecutableParser,
	}

	testcaseOptions := &testcase.TestcaseOptions{
		Logger:        mockLogger,
		RequestParser: mockHTTPRequestParser,
	}
	expectedRawRequest := `POST /some-path HTTP/1.0
Host: url

post
body`
	expectedTestcase := testcase.NewTestcase("Title", "Wonderful\nworld", expectedRawRequest, testcaseOptions)

	parser := testcase.NewParser(options)
	result, err := parser.FromMarkdownFile("mock_file")

	assert.NoError(t, err)
	assert.Equal(t, expectedTestcase, result)
	mockFileReader.AssertExpectations(t)
	mockHTTPRequestParser.AssertExpectations(t)
}

func TestDescriptionAtOtherH1(t *testing.T) {
	config := &config.Config{
		BaseDir: "",
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

	mockFileReader := new(reader.MockFileReader)
	mockFileReader.On("Read", "mock_file").Return(lines, nil)

	mockHTTPRequestParser := new(httphelper.MockHTTPRequestParser)

	mockLogger := new(logger.MockLogger)
	mockLogger.On("PrintDetailedDebugMessage", mock.Anything).Return(nil)

	mockExecutableParser := new(executable.MockExecutableParser)

	options := &testcase.ParserOptions{
		Configurator:     config,
		FileReader:       mockFileReader,
		RequestParser:    mockHTTPRequestParser,
		Logger:           mockLogger,
		ExecutableParser: mockExecutableParser,
	}

	testcaseOptions := &testcase.TestcaseOptions{
		Logger:        mockLogger,
		RequestParser: mockHTTPRequestParser,
	}
	expectedRawRequest := `POST /some-path HTTP/1.0
Host: url

post
body`
	expectedTestcase := testcase.NewTestcase("Title", "Wonderful", expectedRawRequest, testcaseOptions)

	parser := testcase.NewParser(options)
	result, err := parser.FromMarkdownFile("mock_file")

	assert.NoError(t, err)
	assert.Equal(t, expectedTestcase, result)
	mockFileReader.AssertExpectations(t)
	mockHTTPRequestParser.AssertExpectations(t)
}

func TestGetHTTPRequestWithoutTitle(t *testing.T) {
	config := &config.Config{
		BaseDir: "",
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

	mockFileReader := new(reader.MockFileReader)
	mockFileReader.On("Read", "mock_file").Return(lines, nil)

	mockHTTPRequestParser := new(httphelper.MockHTTPRequestParser)

	mockLogger := new(logger.MockLogger)
	mockLogger.On("PrintDetailedDebugMessage", mock.Anything).Return(nil)

	mockExecutableParser := new(executable.MockExecutableParser)

	options := &testcase.ParserOptions{
		Configurator:     config,
		FileReader:       mockFileReader,
		RequestParser:    mockHTTPRequestParser,
		Logger:           mockLogger,
		ExecutableParser: mockExecutableParser,
	}

	testcaseOptions := &testcase.TestcaseOptions{
		Logger:        mockLogger,
		RequestParser: mockHTTPRequestParser,
	}
	expectedRawRequest := `POST /some-path HTTP/1.0
Host: url

post
body`
	expectedTestcase := testcase.NewTestcase("", "", expectedRawRequest, testcaseOptions)

	parser := testcase.NewParser(options)
	result, err := parser.FromMarkdownFile("mock_file")

	assert.NoError(t, err)
	assert.Equal(t, expectedTestcase, result)
	mockFileReader.AssertExpectations(t)
	mockHTTPRequestParser.AssertExpectations(t)
}

func TestNotGetHTTPRequestWithinOtherCodeblock(t *testing.T) {
	config := &config.Config{
		BaseDir: "",
	}

	lines := []line.Line{
		"~~~markdown",
		"```http",
		"POST /some-other-path HTTP/1.0",
		"Host: url",
		"",
		"post",
		"body",
		"```",
		"~~~",
		"",
		"```http",
		"POST /some-path HTTP/1.0",
		"Host: url",
		"",
		"post",
		"body",
		"```",
	}

	mockFileReader := new(reader.MockFileReader)
	mockFileReader.On("Read", "mock_file").Return(lines, nil)

	mockHTTPRequestParser := new(httphelper.MockHTTPRequestParser)

	mockLogger := new(logger.MockLogger)
	mockLogger.On("PrintDetailedDebugMessage", mock.Anything).Return(nil)

	mockExecutableParser := new(executable.MockExecutableParser)

	options := &testcase.ParserOptions{
		Configurator:     config,
		FileReader:       mockFileReader,
		RequestParser:    mockHTTPRequestParser,
		Logger:           mockLogger,
		ExecutableParser: mockExecutableParser,
	}

	testcaseOptions := &testcase.TestcaseOptions{
		Logger:        mockLogger,
		RequestParser: mockHTTPRequestParser,
	}
	expectedRawRequest := `POST /some-path HTTP/1.0
Host: url

post
body`
	expectedTestcase := testcase.NewTestcase("", "", expectedRawRequest, testcaseOptions)

	parser := testcase.NewParser(options)
	result, err := parser.FromMarkdownFile("mock_file")

	assert.NoError(t, err)
	assert.Equal(t, expectedTestcase, result)
	mockFileReader.AssertExpectations(t)
	mockHTTPRequestParser.AssertExpectations(t)
}

func TestNotGetHTTPRequestWithinOtherCodeblockFlip(t *testing.T) {
	config := &config.Config{
		BaseDir: "",
	}

	lines := []line.Line{
		"```markdown",
		"~~~http",
		"POST /some-other-path HTTP/1.0",
		"Host: url",
		"",
		"post",
		"body",
		"~~~",
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

	mockFileReader := new(reader.MockFileReader)
	mockFileReader.On("Read", "mock_file").Return(lines, nil)

	mockHTTPRequestParser := new(httphelper.MockHTTPRequestParser)

	mockLogger := new(logger.MockLogger)
	mockLogger.On("PrintDetailedDebugMessage", mock.Anything).Return(nil)

	mockExecutableParser := new(executable.MockExecutableParser)

	options := &testcase.ParserOptions{
		Configurator:     config,
		FileReader:       mockFileReader,
		RequestParser:    mockHTTPRequestParser,
		Logger:           mockLogger,
		ExecutableParser: mockExecutableParser,
	}

	testcaseOptions := &testcase.TestcaseOptions{
		Logger:        mockLogger,
		RequestParser: mockHTTPRequestParser,
	}
	expectedRawRequest := `POST /some-path HTTP/1.0
Host: url

post
body`
	expectedTestcase := testcase.NewTestcase("", "", expectedRawRequest, testcaseOptions)

	parser := testcase.NewParser(options)
	result, err := parser.FromMarkdownFile("mock_file")

	assert.NoError(t, err)
	assert.Equal(t, expectedTestcase, result)
	mockFileReader.AssertExpectations(t)
	mockHTTPRequestParser.AssertExpectations(t)
}

func TestGetHTTPRequestFromThreeTildedCodeBlock(t *testing.T) {
	config := &config.Config{
		BaseDir: "",
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

	mockFileReader := new(reader.MockFileReader)
	mockFileReader.On("Read", "mock_file").Return(lines, nil)

	mockHTTPRequestParser := new(httphelper.MockHTTPRequestParser)

	mockLogger := new(logger.MockLogger)
	mockLogger.On("PrintDetailedDebugMessage", mock.Anything).Return(nil)

	mockExecutableParser := new(executable.MockExecutableParser)

	options := &testcase.ParserOptions{
		Configurator:     config,
		FileReader:       mockFileReader,
		RequestParser:    mockHTTPRequestParser,
		Logger:           mockLogger,
		ExecutableParser: mockExecutableParser,
	}

	testcaseOptions := &testcase.TestcaseOptions{
		Logger:        mockLogger,
		RequestParser: mockHTTPRequestParser,
	}
	expectedRawRequest := `POST /some-path HTTP/1.0
Host: url

post
body`
	expectedTestcase := testcase.NewTestcase("", "", expectedRawRequest, testcaseOptions)

	parser := testcase.NewParser(options)
	result, err := parser.FromMarkdownFile("mock_file")

	assert.NoError(t, err)
	assert.Equal(t, expectedTestcase, result)
	mockFileReader.AssertExpectations(t)
	mockHTTPRequestParser.AssertExpectations(t)
}

func TestDiscardHTTPRequestInNonHTTPCodeBlock(t *testing.T) {
	config := &config.Config{
		BaseDir: "",
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

	mockFileReader := new(reader.MockFileReader)
	mockFileReader.On("Read", "mock_file").Return(lines, nil)

	mockHTTPRequestParser := new(httphelper.MockHTTPRequestParser)

	mockLogger := new(logger.MockLogger)
	mockLogger.On("PrintDetailedDebugMessage", mock.Anything).Return(nil)

	mockExecutableParser := new(executable.MockExecutableParser)

	options := &testcase.ParserOptions{
		Configurator:     config,
		FileReader:       mockFileReader,
		RequestParser:    mockHTTPRequestParser,
		Logger:           mockLogger,
		ExecutableParser: mockExecutableParser,
	}

	parser := testcase.NewParser(options)
	result, err := parser.FromMarkdownFile("mock_file")

	assert.Error(t, err)
	assert.Nil(t, result)
	mockFileReader.AssertExpectations(t)
	mockHTTPRequestParser.AssertNotCalled(t, "Parse")
}

func TestGetHTTPRequestInOtherHTTPCodeBlock(t *testing.T) {
	config := &config.Config{
		BaseDir: "",
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

	mockFileReader := new(reader.MockFileReader)
	mockFileReader.On("Read", "mock_file").Return(lines, nil)

	mockHTTPRequestParser := new(httphelper.MockHTTPRequestParser)

	mockLogger := new(logger.MockLogger)
	mockLogger.On("PrintDetailedDebugMessage", mock.Anything).Return(nil)

	mockExecutableParser := new(executable.MockExecutableParser)

	options := &testcase.ParserOptions{
		Configurator:     config,
		FileReader:       mockFileReader,
		RequestParser:    mockHTTPRequestParser,
		Logger:           mockLogger,
		ExecutableParser: mockExecutableParser,
	}

	testcaseOptions := &testcase.TestcaseOptions{
		Logger:        mockLogger,
		RequestParser: mockHTTPRequestParser,
	}
	expectedRawRequest := `POST /some-path HTTP/1.0
Host: url

post
body`
	expectedTestcase := testcase.NewTestcase("", "", expectedRawRequest, testcaseOptions)

	parser := testcase.NewParser(options)
	result, err := parser.FromMarkdownFile("mock_file")

	assert.NoError(t, err)
	assert.Equal(t, expectedTestcase, result)
	mockFileReader.AssertExpectations(t)
	mockHTTPRequestParser.AssertExpectations(t)
}

func TestGetHTTPRequestInMixedHTTPCodeBlock(t *testing.T) {
	config := &config.Config{
		BaseDir: "",
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

	mockFileReader := new(reader.MockFileReader)
	mockFileReader.On("Read", "mock_file").Return(lines, nil)

	mockHTTPRequestParser := new(httphelper.MockHTTPRequestParser)

	mockLogger := new(logger.MockLogger)
	mockLogger.On("PrintDetailedDebugMessage", mock.Anything).Return(nil)

	mockExecutableParser := new(executable.MockExecutableParser)

	options := &testcase.ParserOptions{
		Configurator:     config,
		FileReader:       mockFileReader,
		RequestParser:    mockHTTPRequestParser,
		Logger:           mockLogger,
		ExecutableParser: mockExecutableParser,
	}

	testcaseOptions := &testcase.TestcaseOptions{
		Logger:        mockLogger,
		RequestParser: mockHTTPRequestParser,
	}
	expectedRawRequest := `POST /some-path HTTP/1.0
Host: url

post
body`
	expectedTestcase := testcase.NewTestcase("", "", expectedRawRequest, testcaseOptions)

	parser := testcase.NewParser(options)
	result, err := parser.FromMarkdownFile("mock_file")

	assert.NoError(t, err)
	assert.Equal(t, expectedTestcase, result)
	mockFileReader.AssertExpectations(t)
	mockHTTPRequestParser.AssertExpectations(t)
}

func TestGetHTTPRequestInOnlyFirstHTTPCodeBlock(t *testing.T) {
	config := &config.Config{
		BaseDir: "",
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
	mockFileReader := new(reader.MockFileReader)
	mockFileReader.On("Read", "mock_file").Return(lines, nil)

	mockHTTPRequestParser := new(httphelper.MockHTTPRequestParser)

	mockLogger := new(logger.MockLogger)
	mockLogger.On("PrintDetailedDebugMessage", mock.Anything).Return(nil)

	mockExecutableParser := new(executable.MockExecutableParser)

	options := &testcase.ParserOptions{
		Configurator:     config,
		FileReader:       mockFileReader,
		RequestParser:    mockHTTPRequestParser,
		Logger:           mockLogger,
		ExecutableParser: mockExecutableParser,
	}

	testcaseOptions := &testcase.TestcaseOptions{
		Logger:        mockLogger,
		RequestParser: mockHTTPRequestParser,
	}
	expectedRawRequest := `POST /this-should-be-collected HTTP/1.0
Host: url

post
body`
	expectedTestcase := testcase.NewTestcase("", "", expectedRawRequest, testcaseOptions)

	parser := testcase.NewParser(options)
	result, err := parser.FromMarkdownFile("mock_file")

	assert.NoError(t, err)
	assert.Equal(t, expectedTestcase, result)
	mockFileReader.AssertExpectations(t)
	mockHTTPRequestParser.AssertExpectations(t)
}

func TestGetExecutablesBeforeTest(t *testing.T) {
	config := &config.Config{
		BaseDir: "",
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

	mockHTTPRequestParser := new(httphelper.MockHTTPRequestParser)

	mockLogger := new(logger.MockLogger)
	mockLogger.On("PrintDetailedDebugMessage", mock.Anything).Return(nil)

	executableOption := &executable.ExecutableOptions{
		Logger:        mockLogger,
		RequestParser: mockHTTPRequestParser,
	}
	mockSetup1 := executable.New("link to executable", "GET /healthcheck HTTP/1.0\nHost: localhost", executableOption)
	mockSetup2 := executable.New("link to another executable", "GET /readiness HTTP/1.0\nHost: localhost", executableOption)
	mockSetup2.AddCapture(capture.New("readiness", "$.readiness.status"))

	mockFileReader := new(reader.MockFileReader)
	mockFileReader.On("Read", "mock_file").Return(lines, nil)

	mockExecutableParser := new(executable.MockExecutableParser)
	mockExecutableParser.On("FromMarkdownFile", "link1").Return(mockSetup1, nil)
	mockExecutableParser.On("FromMarkdownFile", "link2").Return(mockSetup2, nil)

	options := &testcase.ParserOptions{
		Configurator:     config,
		FileReader:       mockFileReader,
		RequestParser:    mockHTTPRequestParser,
		Logger:           mockLogger,
		ExecutableParser: mockExecutableParser,
	}

	testcaseOptions := &testcase.TestcaseOptions{
		Logger:        mockLogger,
		RequestParser: mockHTTPRequestParser,
	}
	expectedRawRequest := `POST /some-path HTTP/1.0
Host: url

post
body`
	expectedTestcase := testcase.NewTestcase("", "", expectedRawRequest, testcaseOptions)
	expectedTestcase.AddSetup(mockSetup1)
	expectedTestcase.AddSetup(mockSetup2)

	parser := testcase.NewParser(options)
	result, err := parser.FromMarkdownFile("mock_file")

	assert.NoError(t, err)
	assert.Equal(t, expectedTestcase, result)
}

func TestGetExecutablesAfterTest(t *testing.T) {
	config := &config.Config{
		BaseDir: "",
	}

	lines := []line.Line{
		"## Method under test",
		"",
		"```http",
		"POST /some-path HTTP/1.0",
		"Host: url",
		"",
		"post",
		"body",
		"```",
		"",
		"## After",
		"",
		"* [link to executable](link1)",
		"* [link to another executable](link2)",
	}

	mockHTTPRequestParser := new(httphelper.MockHTTPRequestParser)

	mockLogger := new(logger.MockLogger)
	mockLogger.On("PrintDetailedDebugMessage", mock.Anything).Return(nil)

	executableOption := &executable.ExecutableOptions{
		Logger:        mockLogger,
		RequestParser: mockHTTPRequestParser,
	}
	mockTeardown1 := executable.New("link to executable", "GET /healthcheck HTTP/1.0\nHost: localhost", executableOption)
	mockTeardown2 := executable.New("link to another executable", "GET /readiness HTTP/1.0\nHost: localhost", executableOption)
	mockTeardown2.AddCapture(capture.New("readiness", "$.readiness.status"))

	mockFileReader := new(reader.MockFileReader)
	mockFileReader.On("Read", "mock_file").Return(lines, nil)

	mockExecutableParser := new(executable.MockExecutableParser)
	mockExecutableParser.On("FromMarkdownFile", "link1").Return(mockTeardown1, nil)
	mockExecutableParser.On("FromMarkdownFile", "link2").Return(mockTeardown2, nil)

	options := &testcase.ParserOptions{
		Configurator:     config,
		FileReader:       mockFileReader,
		RequestParser:    mockHTTPRequestParser,
		Logger:           mockLogger,
		ExecutableParser: mockExecutableParser,
	}

	testcaseOptions := &testcase.TestcaseOptions{
		Logger:        mockLogger,
		RequestParser: mockHTTPRequestParser,
	}
	expectedRawRequest := `POST /some-path HTTP/1.0
Host: url

post
body`
	expectedTestcase := testcase.NewTestcase("", "", expectedRawRequest, testcaseOptions)
	expectedTestcase.AddTeardown(mockTeardown1)
	expectedTestcase.AddTeardown(mockTeardown2)

	parser := testcase.NewParser(options)
	result, err := parser.FromMarkdownFile("mock_file")

	assert.NoError(t, err)
	assert.Equal(t, expectedTestcase, result)
}

func TestGetCapturesInTestcase(t *testing.T) {
	config := &config.Config{
		BaseDir: "",
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
		"* varname:`$.result`",
	}

	mockFileReader := new(reader.MockFileReader)
	mockFileReader.On("Read", "mock_file").Return(lines, nil)

	mockHTTPRequestParser := new(httphelper.MockHTTPRequestParser)

	mockLogger := new(logger.MockLogger)
	mockLogger.On("PrintDetailedDebugMessage", mock.Anything).Return(nil)

	mockExecutableParser := new(executable.MockExecutableParser)

	options := &testcase.ParserOptions{
		Configurator:     config,
		FileReader:       mockFileReader,
		RequestParser:    mockHTTPRequestParser,
		Logger:           mockLogger,
		ExecutableParser: mockExecutableParser,
	}

	testcaseOptions := &testcase.TestcaseOptions{
		Logger:        mockLogger,
		RequestParser: mockHTTPRequestParser,
	}
	expectedRawRequest := `POST /this-should-be-collected HTTP/1.0
Host: url

post
body`
	expectedTestcase := testcase.NewTestcase("", "", expectedRawRequest, testcaseOptions)
	expectedTestcase.AddCapture(capture.New("varname", "$.result"))

	parser := testcase.NewParser(options)
	result, err := parser.FromMarkdownFile("mock_file")

	assert.NoError(t, err)
	assert.Equal(t, expectedTestcase, result)

}

func TestGetAssertion(t *testing.T) {
	config := &config.Config{
		BaseDir: "",
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

	mockFileReader := new(reader.MockFileReader)
	mockFileReader.On("Read", "mock_file").Return(lines, nil)

	mockHTTPRequestParser := new(httphelper.MockHTTPRequestParser)

	mockLogger := new(logger.MockLogger)
	mockLogger.On("PrintDetailedDebugMessage", mock.Anything).Return(nil)

	mockExecutableParser := new(executable.MockExecutableParser)

	options := &testcase.ParserOptions{
		Configurator:     config,
		FileReader:       mockFileReader,
		RequestParser:    mockHTTPRequestParser,
		Logger:           mockLogger,
		ExecutableParser: mockExecutableParser,
	}

	testcaseOptions := &testcase.TestcaseOptions{
		Logger:        mockLogger,
		RequestParser: mockHTTPRequestParser,
	}
	expectedRawRequest := `POST /this-should-be-collected HTTP/1.0
Host: url

post
body`
	eqOp, _ := assertion.NewOp("==")
	expectedTestcase := testcase.NewTestcase("", "", expectedRawRequest, testcaseOptions)
	expectedTestcase.AddAssertion(assertion.New("$.var", eqOp, float64(3)))
	expectedTestcase.AddAssertion(assertion.New("$.var2", eqOp, "good.vibe"))

	parser := testcase.NewParser(options)
	result, err := parser.FromMarkdownFile("mock_file")

	assert.NoError(t, err)
	assert.Equal(t, expectedTestcase, result)
}
