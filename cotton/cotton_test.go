package cotton

import (
	"errors"
	"os"
	"testing"

	"github.com/chonla/cotton/testsuite"
	"github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/assert"
)

type MockParser struct {
	mock.Mock
}

func (m *MockParser) Parse(path string) (testsuite.TestSuitesInterface, error) {
	args := m.Called(path)
	return args.Get(0).(testsuite.TestSuitesInterface), args.Error(1)
}

func (m *MockParser) ParseFile(file string) (*testsuite.TestSuite, error) {
	args := m.Called(file)
	return args.Get(0).(*testsuite.TestSuite), args.Error(1)
}

func (m *MockParser) ParseString(content, filePath string) ([]*testsuite.TestCase, error) {
	args := m.Called(content, filePath)
	return args.Get(0).([]*testsuite.TestCase), args.Error(1)
}

type MockTestSuites struct {
	mock.Mock
}

func (m *MockTestSuites) Run() {
	m.Called()
}

func (m *MockTestSuites) Summary() int {
	args := m.Called()
	return args.Int(0)
}

func (m *MockTestSuites) SetVariables(v map[string]string) {
	m.Called(v)
}

func (m *MockTestSuites) Stat() testsuite.TestStat {
	args := m.Called()
	return args.Get(0).(testsuite.TestStat)
}

func (m *MockTestSuites) SetBaseURL(url string) {
	m.Called(url)
}

func (m *MockTestSuites) SetConfig(conf *testsuite.Config) {
	m.Called(conf)
}

func TestNewCottonShouldSuccessWhenTestPathOrFileExists(t *testing.T) {
	// stub
	statFile = func(path string) (os.FileInfo, error) {
		return nil, nil
	}

	c, e := NewCotton("existing/path", Config{
		BaseURL:   "http://www.abc.com",
		Insecure:  true,
		Verbose:   true,
		Variables: []string{},
	})

	assert.Nil(t, e)
	assert.Equal(t, &Cotton{
		path: "existing/path",
		Config: Config{
			BaseURL:   "http://www.abc.com",
			Insecure:  true,
			Verbose:   true,
			Variables: []string{},
		},
	}, c)
}

func TestNewCottonShouldFailWhenTestPathOrFileDoesNotExist(t *testing.T) {
	// stub
	statFile = func(path string) (os.FileInfo, error) {
		return nil, &os.PathError{Op: "stat", Path: "non-existing/path", Err: errors.New("some error")}
	}

	_, e := NewCotton("non-existing/path", Config{
		BaseURL:   "http://www.abc.com",
		Insecure:  true,
		Verbose:   true,
		Variables: []string{},
	})

	assert.NotNil(t, e)
}

func TestRunTestSuitesSuccess(t *testing.T) {
	// stub
	statFile = func(path string) (os.FileInfo, error) {
		return nil, nil
	}

	c, _ := NewCotton("existing/path", Config{
		BaseURL:   "http://www.abc.com",
		Insecure:  true,
		Verbose:   true,
		Variables: []string{},
	})

	parser := new(MockParser)

	suites := new(MockTestSuites)
	suites.On("Run")
	suites.On("Summary").Return(0)
	suites.On("SetVariables", map[string]string{
		"a": "1",
		"b": "2",
	})
	suites.On("SetBaseURL", "http://www.abc.com")
	suites.On("SetConfig", mock.Anything)
	suites.On("Stat").Return(testsuite.TestStat{
		Total:   10,
		Success: 10,
	})

	parser.On("Parse", "existing/path").Return(suites, nil)

	c.SetParser(parser)

	stat, code := c.Run()

	assert.Equal(t, 0, code)
	assert.Equal(t, testsuite.TestStat{
		Total:   10,
		Success: 10,
	}, stat)
}

func TestRunTestSuitesSuccessWhenSomeVarsPassedTo(t *testing.T) {
	// stub
	statFile = func(path string) (os.FileInfo, error) {
		return nil, nil
	}

	c, _ := NewCotton("existing/path", Config{
		BaseURL:   "http://www.abc.com",
		Insecure:  true,
		Verbose:   true,
		Variables: []string{"a=1", "b=2"},
	})

	parser := new(MockParser)

	suites := new(MockTestSuites)
	suites.On("Run")
	suites.On("Summary").Return(100)
	suites.On("SetVariables", map[string]string{
		"a": "1",
		"b": "2",
	})
	suites.On("SetBaseURL", "http://www.abc.com")
	suites.On("SetConfig", mock.Anything)
	suites.On("Stat").Return(testsuite.TestStat{
		Total:   10,
		Success: 10,
	})

	parser.On("Parse", "existing/path").Return(suites, nil)

	c.SetParser(parser)

	c.Run()

	suites.AssertNumberOfCalls(t, "Run", 1)
	suites.AssertNumberOfCalls(t, "Summary", 1)
	suites.AssertNumberOfCalls(t, "Stat", 1)
	suites.AssertCalled(t, "SetVariables", map[string]string{
		"a": "1",
		"b": "2",
	})
	suites.AssertCalled(t, "SetBaseURL", "http://www.abc.com")
}

func TestRunTestSuitesSuccessWhenSomeInvalidVarsPassedTo(t *testing.T) {
	// stub
	statFile = func(path string) (os.FileInfo, error) {
		return nil, nil
	}

	c, _ := NewCotton("existing/path", Config{
		BaseURL:   "http://www.abc.com",
		Insecure:  true,
		Verbose:   true,
		Variables: []string{"a=1", "b=2", "c!!"},
	})

	parser := new(MockParser)

	suites := new(MockTestSuites)
	suites.On("Run")
	suites.On("Summary").Return(100)
	suites.On("SetVariables", map[string]string{
		"a": "1",
		"b": "2",
	})
	suites.On("SetBaseURL", "http://www.abc.com")
	suites.On("SetConfig", mock.Anything)
	suites.On("Stat").Return(testsuite.TestStat{
		Total:   10,
		Success: 10,
	})

	parser.On("Parse", "existing/path").Return(suites, nil)

	c.SetParser(parser)

	c.Run()

	suites.AssertNumberOfCalls(t, "Run", 1)
	suites.AssertNumberOfCalls(t, "Summary", 1)
	suites.AssertNumberOfCalls(t, "Stat", 1)
	suites.AssertCalled(t, "SetVariables", map[string]string{
		"a": "1",
		"b": "2",
	})
	suites.AssertCalled(t, "SetBaseURL", "http://www.abc.com")
}

func TestRunTestSuitesSuccessWhenParseFileError(t *testing.T) {
	// stub
	statFile = func(path string) (os.FileInfo, error) {
		return nil, nil
	}

	c, _ := NewCotton("existing/path", Config{
		BaseURL:   "http://www.abc.com",
		Insecure:  true,
		Verbose:   true,
		Variables: []string{"a=1", "b=2", "c!!"},
	})

	parser := new(MockParser)

	parser.On("Parse", "existing/path").Return(&testsuite.TestSuites{}, errors.New("some error"))

	c.SetParser(parser)

	stat, e := c.Run()

	assert.NotNil(t, e)
	assert.Equal(t, testsuite.TestStat{}, stat)
}
