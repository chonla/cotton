package executable

import (
	"bufio"
	"cotton/internal/capture"
	"cotton/internal/line"
	"cotton/internal/reader"
	"net/http"
	"strings"
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
	var req []string
	var exReq *http.Request

	collectingCodeBlockBackTick := false

	ex := &Executable{}
	for _, mdLine := range mdLines {
		if mdLine.LookLike("^```http$") && exReq == nil {
			collectingCodeBlockBackTick = true
			continue
		}

		if collectingCodeBlockBackTick {
			if ok := mdLine.LookLike("^```$"); ok {
				collectingCodeBlockBackTick = false

				if len(req) > 0 {
					reqReader := bufio.NewReader(strings.NewReader(line.Line(strings.Join(req, "\n")).Trim().Value()))
					httpRequest, err := http.ReadRequest(reqReader)
					if err == nil {
						exReq = httpRequest
					}
					req = nil
				}
			} else {
				if req == nil {
					req = []string{}
				}
				req = append(req, mdLine.Value())
			}
		} else {
			if captured, ok := capture.Try(mdLine); ok {
				if ex.Captures == nil {
					ex.Captures = []*capture.Captured{}
				}
				ex.Captures = append(ex.Captures, captured)
			}
		}
	}

	ex.Request = exReq

	return ex, nil
}
