package directory_test

import (
	"cotton/internal/directory"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDirectoryOfDirectory(t *testing.T) {
	tmpDir, _ := os.MkdirTemp("", "tmpdir-for-directoryof")
	defer os.RemoveAll(tmpDir)

	dir := directory.New()

	result, err := dir.DirectoryOf(tmpDir)

	assert.NoError(t, err)
	assert.Equal(t, tmpDir, result)
}

func TestDirectoryOfFile(t *testing.T) {
	tmpDir, _ := os.MkdirTemp("", "tmpdir-for-directoryof")
	tmpFile, _ := os.CreateTemp(tmpDir, "tmpfile-for-directoryof")
	defer func() {
		os.RemoveAll(tmpFile.Name())
		os.RemoveAll(tmpDir)
	}()

	dir := directory.New()

	result, err := dir.DirectoryOf(tmpFile.Name())

	assert.NoError(t, err)
	assert.Equal(t, tmpDir, result)
}
