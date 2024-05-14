package reader

import (
	"cotton/internal/line"
)

type FileReadingFunction func(fileName string) ([]byte, error)

type Reader interface {
	Read(fileName string) ([]line.Line, error)
}

type FileReader struct {
	readerFunc FileReadingFunction
}

func New(readerFunc FileReadingFunction) *FileReader {
	return &FileReader{
		readerFunc: readerFunc,
	}
}

func (fr *FileReader) Read(fileName string) ([]line.Line, error) {
	b, err := fr.readerFunc(fileName)
	if err != nil {
		return nil, err
	}

	s := line.FromMultilineString(string(b))
	return s, nil
}
