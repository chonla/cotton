package executable_test

import (
	"cotton/internal/capture"
	"cotton/internal/config"
	"cotton/internal/executable"
	"cotton/internal/httphelper"
	"cotton/internal/line"
	"cotton/internal/logger"
	"cotton/internal/reader"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestParseMarkdownFileButFailed(t *testing.T) {
	config := &config.Config{
		BaseDir: "",
	}

	mockFileReader := new(reader.MockFileReader)
	mockFileReader.On("Read", "mock_file").Return([]line.Line(nil), errors.New("file not found"))

	mockHTTPRequestParser := new(httphelper.MockHTTPRequestParser)

	mockLogger := new(logger.MockLogger)
	mockLogger.On("PrintDetailedDebugMessage", mock.Anything).Return(nil)

	options := &executable.ParserOptions{
		Configurator:  config,
		FileReader:    mockFileReader,
		RequestParser: mockHTTPRequestParser,
		Logger:        mockLogger,
	}

	parser := executable.NewParser(options)
	result, err := parser.FromMarkdownFile("mock_file")

	assert.Nil(t, result)
	assert.Equal(t, errors.New("file not found"), err)
	mockFileReader.AssertExpectations(t)
	mockHTTPRequestParser.AssertNotCalled(t, "Parse")
}

func TestParseMarkdownFileShouldReadFromGivenFile(t *testing.T) {
	config := &config.Config{
		BaseDir: "",
	}

	lines := []line.Line{""}

	mockFileReader := new(reader.MockFileReader)
	mockFileReader.On("Read", "mock_file").Return(lines, nil)

	mockHTTPRequestParser := new(httphelper.MockHTTPRequestParser)

	mockLogger := new(logger.MockLogger)
	mockLogger.On("PrintDetailedDebugMessage", mock.Anything).Return(nil)

	options := &executable.ParserOptions{
		Configurator:  config,
		FileReader:    mockFileReader,
		RequestParser: mockHTTPRequestParser,
		Logger:        mockLogger,
	}

	parser := executable.NewParser(options)
	parser.FromMarkdownFile("mock_file")

	mockFileReader.AssertExpectations(t)
	mockHTTPRequestParser.AssertNotCalled(t, "Parse")
}

func TestGetHTTPRequestInHTTPCodeBlock(t *testing.T) {
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

	expectedExecutableOptions := &executable.ExecutableOptions{
		Logger:        mockLogger,
		RequestParser: mockHTTPRequestParser,
	}
	expectedRawRequest := `POST /some-path HTTP/1.0
Host: url

post
body`
	expectedExecutable := executable.New("Untitled", expectedRawRequest, expectedExecutableOptions)

	options := &executable.ParserOptions{
		Configurator:  config,
		FileReader:    mockFileReader,
		RequestParser: mockHTTPRequestParser,
		Logger:        mockLogger,
	}

	parser := executable.NewParser(options)
	result, err := parser.FromMarkdownFile("mock_file")

	assert.NoError(t, err)
	assert.Equal(t, expectedExecutable, result)
	mockFileReader.AssertExpectations(t)
}

func TestGetHTTPRequestInThreeTildedHTTPCodeBlock(t *testing.T) {
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

	expectedExecutableOptions := &executable.ExecutableOptions{
		Logger:        mockLogger,
		RequestParser: mockHTTPRequestParser,
	}
	expectedRawRequest := `POST /some-path HTTP/1.0
Host: url

post
body`
	expectedExecutable := executable.New("Untitled", expectedRawRequest, expectedExecutableOptions)

	options := &executable.ParserOptions{
		Configurator:  config,
		FileReader:    mockFileReader,
		RequestParser: mockHTTPRequestParser,
		Logger:        mockLogger,
	}

	parser := executable.NewParser(options)
	result, err := parser.FromMarkdownFile("mock_file")

	assert.NoError(t, err)
	assert.Equal(t, expectedExecutable, result)
	mockFileReader.AssertExpectations(t)
}

func TestNotGetHTTPRequestInHTTPCodeBlockInOtherCodeBlock(t *testing.T) {
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

	expectedExecutableOptions := &executable.ExecutableOptions{
		Logger:        mockLogger,
		RequestParser: mockHTTPRequestParser,
	}
	expectedRawRequest := `POST /some-path HTTP/1.0
Host: url

post
body`
	expectedExecutable := executable.New("Untitled", expectedRawRequest, expectedExecutableOptions)

	options := &executable.ParserOptions{
		Configurator:  config,
		FileReader:    mockFileReader,
		RequestParser: mockHTTPRequestParser,
		Logger:        mockLogger,
	}

	parser := executable.NewParser(options)
	result, err := parser.FromMarkdownFile("mock_file")

	assert.NoError(t, err)
	assert.Equal(t, expectedExecutable, result)
	mockFileReader.AssertExpectations(t)
}

func TestNotGetHTTPRequestInHTTPCodeBlockInOtherCodeBlockFlip(t *testing.T) {
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

	expectedExecutableOptions := &executable.ExecutableOptions{
		Logger:        mockLogger,
		RequestParser: mockHTTPRequestParser,
	}
	expectedRawRequest := `POST /some-path HTTP/1.0
Host: url

post
body`
	expectedExecutable := executable.New("Untitled", expectedRawRequest, expectedExecutableOptions)

	options := &executable.ParserOptions{
		Configurator:  config,
		FileReader:    mockFileReader,
		RequestParser: mockHTTPRequestParser,
		Logger:        mockLogger,
	}

	parser := executable.NewParser(options)
	result, err := parser.FromMarkdownFile("mock_file")

	assert.NoError(t, err)
	assert.Equal(t, expectedExecutable, result)
	mockFileReader.AssertExpectations(t)
}

func TestDiscardHTTPRequestInNonHTTPCodeBlockWillCauseANilExecutable(t *testing.T) {
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

	options := &executable.ParserOptions{
		Configurator:  config,
		FileReader:    mockFileReader,
		RequestParser: mockHTTPRequestParser,
		Logger:        mockLogger,
	}

	parser := executable.NewParser(options)
	result, err := parser.FromMarkdownFile("mock_file")

	assert.Equal(t, errors.New("no callable request"), err)
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

	expectedExecutableOptions := &executable.ExecutableOptions{
		Logger:        mockLogger,
		RequestParser: mockHTTPRequestParser,
	}
	expectedRawRequest := `POST /some-path HTTP/1.0
Host: url

post
body`
	expectedExecutable := executable.New("Untitled", expectedRawRequest, expectedExecutableOptions)

	options := &executable.ParserOptions{
		Configurator:  config,
		FileReader:    mockFileReader,
		RequestParser: mockHTTPRequestParser,
		Logger:        mockLogger,
	}

	parser := executable.NewParser(options)
	result, err := parser.FromMarkdownFile("mock_file")

	assert.NoError(t, err)
	assert.Equal(t, expectedExecutable, result)
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

	expectedExecutableOptions := &executable.ExecutableOptions{
		Logger:        mockLogger,
		RequestParser: mockHTTPRequestParser,
	}
	expectedRawRequest := `POST /some-path HTTP/1.0
Host: url

post
body`
	expectedExecutable := executable.New("Untitled", expectedRawRequest, expectedExecutableOptions)

	options := &executable.ParserOptions{
		Configurator:  config,
		FileReader:    mockFileReader,
		RequestParser: mockHTTPRequestParser,
		Logger:        mockLogger,
	}

	parser := executable.NewParser(options)
	result, err := parser.FromMarkdownFile("mock_file")

	assert.NoError(t, err)
	assert.Equal(t, expectedExecutable, result)
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

	expectedExecutableOptions := &executable.ExecutableOptions{
		Logger:        mockLogger,
		RequestParser: mockHTTPRequestParser,
	}
	expectedRawRequest := `POST /this-should-be-collected HTTP/1.0
Host: url

post
body`
	expectedExecutable := executable.New("Untitled", expectedRawRequest, expectedExecutableOptions)

	options := &executable.ParserOptions{
		Configurator:  config,
		FileReader:    mockFileReader,
		RequestParser: mockHTTPRequestParser,
		Logger:        mockLogger,
	}

	parser := executable.NewParser(options)
	result, err := parser.FromMarkdownFile("mock_file")

	assert.NoError(t, err)
	assert.Equal(t, expectedExecutable, result)
	mockFileReader.AssertExpectations(t)
	mockHTTPRequestParser.AssertExpectations(t)
}

func TestGetCaptures(t *testing.T) {
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
		"* var:`sample`",
		"* var2:`good.vibe`",
	}

	expectedRawRequest := `POST /this-should-be-collected HTTP/1.0
Host: url

post
body`
	mockFileReader := new(reader.MockFileReader)
	mockFileReader.On("Read", "mock_file").Return(lines, nil)

	mockHTTPRequestParser := new(httphelper.MockHTTPRequestParser)

	mockLogger := new(logger.MockLogger)
	mockLogger.On("PrintDetailedDebugMessage", mock.Anything).Return(nil)

	expectedExecutableOptions := &executable.ExecutableOptions{
		Logger:        mockLogger,
		RequestParser: mockHTTPRequestParser,
	}
	expectedExecutable := executable.New("Untitled", expectedRawRequest, expectedExecutableOptions)
	expectedExecutable.AddCapture(capture.New("var", "sample"))
	expectedExecutable.AddCapture(capture.New("var2", "good.vibe"))

	options := &executable.ParserOptions{
		Configurator:  config,
		FileReader:    mockFileReader,
		RequestParser: mockHTTPRequestParser,
		Logger:        mockLogger,
	}

	parser := executable.NewParser(options)
	result, err := parser.FromMarkdownFile("mock_file")

	assert.NoError(t, err)
	assert.Equal(t, expectedExecutable, result)
	mockFileReader.AssertExpectations(t)
	mockHTTPRequestParser.AssertExpectations(t)
}
