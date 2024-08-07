package reader_test

import (
	"cotton/internal/line"
	"cotton/internal/reader"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadFileWithSingleLine(t *testing.T) {
	readerFunction := func(fileName string) ([]byte, error) {
		return []byte(``), nil
	}

	fileReader := reader.New(readerFunction)

	lines, err := fileReader.Read("somefile")

	assert.Nil(t, err)
	assert.Equal(t, []line.Line{""}, lines)
}

func TestReadFileWithMultipleLines(t *testing.T) {
	readerFunction := func(fileName string) ([]byte, error) {
		return []byte(`a file
	with
some lines`), nil
	}

	fileReader := reader.New(readerFunction)

	lines, err := fileReader.Read("somefile")

	assert.Nil(t, err)
	assert.Equal(t, []line.Line{"a file", "	with", "some lines"}, lines)
}

func TestReadFileOnWindows(t *testing.T) {
	readerFunction := func(fileName string) ([]byte, error) {
		return []byte("a file\r\n\twith\r\nsome lines"), nil
	}

	fileReader := reader.New(readerFunction)

	lines, err := fileReader.Read("somefile")

	assert.Nil(t, err)
	assert.Equal(t, []line.Line{"a file", "	with", "some lines"}, lines)
}

func TestReadFileOnNonWindows(t *testing.T) {
	readerFunction := func(fileName string) ([]byte, error) {
		return []byte("a file\n\twith\nsome lines"), nil
	}

	fileReader := reader.New(readerFunction)

	lines, err := fileReader.Read("somefile")

	assert.Nil(t, err)
	assert.Equal(t, []line.Line{"a file", "	with", "some lines"}, lines)
}

func TestReadFileErrorShouldReturnError(t *testing.T) {
	expectedError := errors.New("unable to read file")

	readerFunction := func(fileName string) ([]byte, error) {
		return nil, expectedError
	}

	fileReader := reader.New(readerFunction)

	lines, err := fileReader.Read("somefile")

	assert.Equal(t, expectedError, err)
	assert.Nil(t, lines)
}
