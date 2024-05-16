package testcase

import (
	"bufio"
	"cotton/internal/executable"
	"cotton/internal/line"
	"cotton/internal/reader"
	"net/http"
	"strings"
)

type Parser struct {
	fileReader       reader.Reader
	executableParser *executable.Parser
}

func NewParser(fileReader reader.Reader) *Parser {
	return &Parser{
		fileReader:       fileReader,
		executableParser: executable.NewParser(fileReader),
	}
}

func (p *Parser) FromMarkdownFile(mdFileName string) (*TestCase, error) {
	lines, err := p.fileReader.Read(mdFileName)
	if err != nil {
		return nil, err
	}
	return p.FromMarkdownLines(lines)
}

func (p *Parser) FromMarkdownLines(mdLines []line.Line) (*TestCase, error) {
	title := ""
	description := []string{}
	var sutReq *http.Request
	var req []string

	justTitle := false
	collectingCodeBlockBackTick := false

	tc := &TestCase{}
	for _, mdLine := range mdLines {
		if captured, ok := mdLine.Capture(`^ {0,3}#\s+(.*)`, 1); ok && !justTitle {
			title = captured
			justTitle = true
			continue
		}

		if mdLine.LookLike("^```http$") && sutReq == nil {
			justTitle = false
			collectingCodeBlockBackTick = true
			continue
		}

		if justTitle {
			if ok := mdLine.LookLike(`^ {0,3}#{1,6}\s+(.*)`); ok {
				justTitle = false
				continue
			}

			description = append(description, mdLine.Value())
			continue
		}

		if collectingCodeBlockBackTick {
			if ok := mdLine.LookLike("^```$"); ok {
				collectingCodeBlockBackTick = false

				if len(req) > 0 {
					reqReader := bufio.NewReader(strings.NewReader(line.Line(strings.Join(req, "\n")).Trim().Value()))
					httpRequest, err := http.ReadRequest(reqReader)
					if err == nil {
						sutReq = httpRequest
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
			if captured, ok := mdLine.CaptureAll(`^\s*\*\s\[([^\]]+)\]\(([^\)]+)\)`); ok {
				if sutReq == nil {
					ex, err := p.executableParser.FromMarkdownFile(captured[2])
					if err != nil {
						return nil, err
					}
					ex.Title = captured[1]
					if tc.Setups == nil {
						tc.Setups = []*executable.Executable{}
					}
					tc.Setups = append(tc.Setups, ex)
				} else {
					ex, err := p.executableParser.FromMarkdownFile(captured[2])
					if err != nil {
						return nil, err
					}
					ex.Title = captured[1]
					if tc.Teardowns == nil {
						tc.Teardowns = []*executable.Executable{}
					}
					tc.Teardowns = append(tc.Teardowns, ex)
				}
			}
		}
	}

	tc.Title = title
	tc.Description = line.Line(strings.Join(description, "\n")).Trim().Value()
	tc.Request = sutReq

	return tc, nil
}
