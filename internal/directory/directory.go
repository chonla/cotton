package directory

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// Directory represents a file system directory utility.
type Directory struct{}

// New creates a new Directory instance.
func New() *Directory {
	return &Directory{}
}

// ListFiles lists files with the specified extension in the given path, excluding those starting with an underscore.
func (d *Directory) ListFiles(rootPath, ext string) ([]string, error) {
	var files []string
	ext = strings.ToLower("." + ext) // Ensure lowercase extension for case-insensitive matching

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err // Propagate errors instead of logging and exiting
		}

		if info.IsDir() {
			if info.Name()[0] == '_' {
				return filepath.SkipDir // Skip hidden directories efficiently
			}
			return nil
		}

		if strings.ToLower(filepath.Ext(path)) == ext {
			files = append(files, path)
		}
		return nil
	})

	return files, err
}

// DirectoryOf gets a directory for given path.
// If path is referencing to a file, DirectoryOf returns path containing that file.
// If path is referencing to a directory, DirectoryOf returns itself.
func (d *Directory) DirectoryOf(path string) (string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return "", err
	}

	if os.IsNotExist(err) {
		return "", errors.New("path does not exist")
	} else if err != nil {
		return "", err
	}

	if info.IsDir() {
		return path, nil
	}
	parentPath := filepath.Dir(path)
	return parentPath, nil
}
