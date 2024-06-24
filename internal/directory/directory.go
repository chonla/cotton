package directory

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Directory struct{}

func New() *Directory {
	return &Directory{}
}

func (d *Directory) ListFiles(path string, ext string) ([]string, error) {
	var files []string
	extension := fmt.Sprintf(".%s", ext)
	err := filepath.Walk(path, d.scan(&files, extension))
	return files, err

}

func (d *Directory) scan(files *[]string, ext string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
		}
		if !info.IsDir() && info.Name()[0] == '_' {
			return nil
		}
		if !info.IsDir() && strings.ToLower(filepath.Ext(path)) == ext {
			*files = append(*files, path)
		}
		return nil
	}
}
