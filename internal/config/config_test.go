package config_test

import (
	"cotton/internal/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolvePathWithAbsolutePath(t *testing.T) {
	cfg := &config.Config{
		BaseDir: "/some/base",
	}

	result := cfg.ResolvePath("/some/path")

	assert.Equal(t, "/some/path", result)
}

func TestResolvePathWithRelativePath(t *testing.T) {
	cfg := &config.Config{
		BaseDir: "/some/base",
	}

	result := cfg.ResolvePath("some/path")

	assert.Equal(t, "/some/base/some/path", result)
}

func TestResolvePathWithRelativePathWithDoubleDots(t *testing.T) {
	cfg := &config.Config{
		BaseDir: "/some/base",
	}

	result := cfg.ResolvePath("../some/path")

	assert.Equal(t, "/some/base/../some/path", result)
}

func TestResolvePathWithRelativePathAndBaseDirEndsWithSlash(t *testing.T) {
	cfg := &config.Config{
		BaseDir: "/some/base/",
	}

	result := cfg.ResolvePath("some/path")

	assert.Equal(t, "/some/base/some/path", result)
}

func TestResolveAbsolutePathOnWindowsWithDrive(t *testing.T) {
	cfg := &config.Config{
		BaseDir: ".\\some\\base",
	}

	result := cfg.ResolvePath("C:\\some\\path")

	assert.Equal(t, "C:\\some\\path", result)
}

func TestResolveAbsolutePathOnWindowsWithoutDrive(t *testing.T) {
	cfg := &config.Config{
		BaseDir: ".\\some\\base",
	}

	result := cfg.ResolvePath("\\some\\path")

	assert.Equal(t, "\\some\\path", result)
}

func TestResolveRelativePathOnWindows(t *testing.T) {
	cfg := &config.Config{
		BaseDir: ".\\some\\base",
	}

	result := cfg.ResolvePath("some\\path")

	assert.Equal(t, ".\\some\\base\\some\\path", result)
}

func TestResolveRelativePathWithBaseDirHavingMixedPathSeparatorOnWindows(t *testing.T) {
	cfg := &config.Config{
		BaseDir: ".\\some/base",
	}

	result := cfg.ResolvePath("some/path\\inside")

	assert.Equal(t, ".\\some/base\\some/path\\inside", result)
}
