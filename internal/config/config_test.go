package config_test

import (
	"cotton/internal/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPathWithoutRootDir(t *testing.T) {
	cfg := &config.Config{
		RootDir: "/hello",
	}

	result := cfg.ResolvePath("/some/path")

	assert.Equal(t, "/some/path", result)
}

func TestGetPathWithRootDir(t *testing.T) {
	cfg := &config.Config{
		RootDir: "/hello",
	}

	result := cfg.ResolvePath("<rootDir>/some/path")

	assert.Equal(t, "/hello/some/path", result)
}
