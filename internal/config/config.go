package config

import (
	"cotton/internal/line"
	"fmt"
)

type Config struct {
	BaseDir string
}

func (c *Config) ResolvePath(path string) string {
	pathLine := line.Line(path)
	// windows + linux
	if pathLine.StartsWith("/") || pathLine.LookLike(`^[a-zA-Z]:`) || pathLine.StartsWith("\\") {
		return path
	}
	baseDir := line.Line(c.BaseDir)
	sep := "/"
	if baseDir.Contains("\\") {
		sep = "\\"
	}
	if baseDir.EndsWith("/") || baseDir.EndsWith("\\") || baseDir.Value() == "" {
		sep = ""
	}
	return fmt.Sprintf("%s%s%s", baseDir, sep, path)
}
