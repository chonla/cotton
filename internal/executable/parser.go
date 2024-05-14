package executable

import (
	"cotton/internal/line"
	"cotton/internal/reader"
)

type Parser struct {
	fileReader reader.Reader
}

func NewParser(fileReader reader.Reader) *Parser {
	return &Parser{
		fileReader: fileReader,
	}
}

func (p *Parser) FromMarkdownFile(mdFileName string) (*Executable, error) {
	lines, err := p.fileReader.Read(mdFileName)
	if err != nil {
		return nil, err
	}
	return p.FromMarkdownLines(lines)
}

func (p *Parser) FromMarkdownLines(mdLines []line.Line) (*Executable, error) {
	return nil, nil
}
