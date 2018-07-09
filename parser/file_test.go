package parser

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadLinesFromFileContainingWindowsNewLine(t *testing.T) {
	readFileFn = func(string) ([]byte, error) {
		return []byte("test\r\ntest2\r\ntest3"), nil
	}

	p := NewParser()
	result, _ := p.readTestSuiteFile("test")

	assert.Equal(t, []string{"test", "test2", "test3"}, result)
}

func TestReadLinesFromFileContainingLinuxNewLine(t *testing.T) {
	readFileFn = func(string) ([]byte, error) {
		return []byte("test\ntest2\ntest3"), nil
	}

	p := NewParser()
	result, _ := p.readTestSuiteFile("test")

	assert.Equal(t, []string{"test", "test2", "test3"}, result)
}

func TestReadLinesFromFileWithError(t *testing.T) {
	readFileFn = func(string) ([]byte, error) {
		return []byte{}, errors.New("some error")
	}

	p := NewParser()
	result, _ := p.readTestSuiteFile("test")

	assert.Equal(t, []string{}, result)
}
