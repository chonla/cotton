package config

import "cotton/internal/line"

type Config struct {
	RootDir string
}

func (c *Config) ResolvePath(path string) string {
	pathLine := line.Line(path)
	if pathLine.LookLike("<rootDir>") {
		return pathLine.Replace("<rootDir>", c.RootDir)
	}
	return path
}
