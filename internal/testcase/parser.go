package testcase

import (
	"cotton/internal/capture"
	"cotton/internal/config"
	"cotton/internal/executable"
	"cotton/internal/line"
	"cotton/internal/reader"
	"net/http"
	"strings"

	"github.com/chonla/httpreqparser"
)

type Parser struct {
	config           *config.Config
	fileReader       reader.Reader
	executableParser *executable.Parser
	requestParser    httpreqparser.Parser
}

func NewParser(config *config.Config, fileReader reader.Reader, requestParser httpreqparser.Parser) *Parser {
	return &Parser{
		config:           config,
		fileReader:       fileReader,
		executableParser: executable.NewParser(config, fileReader, requestParser),
		requestParser:    requestParser,
	}
}

func (p *Parser) FromMarkdownFile(mdFileName string) (*TestCase, error) {
	mdFullPath := p.config.ResolvePath(mdFileName)
	lines, err := p.fileReader.Read(mdFullPath)
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
		if cap, ok := mdLine.Capture(`^ {0,3}#\s+(.*)`, 1); ok && !justTitle {
			title = cap
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
					requestString := line.Line(strings.Join(req, "\n")).Value()
					httpRequest, err := p.requestParser.Parse(requestString)
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
			if cap, ok := capture.Try(mdLine); ok {
				if tc.Captures == nil {
					tc.Captures = []*capture.Capture{}
				}
				tc.Captures = append(tc.Captures, cap)
			} else {
				if captures, ok := mdLine.CaptureAll(`^\s*\*\s\[([^\]]+)\]\(([^\)]+)\)`); ok {
					if sutReq == nil {
						ex, err := p.executableParser.FromMarkdownFile(captures[2])
						if err != nil {
							return nil, err
						}
						ex.Title = captures[1]
						if tc.Setups == nil {
							tc.Setups = []*executable.Executable{}
						}
						tc.Setups = append(tc.Setups, ex)
					} else {
						ex, err := p.executableParser.FromMarkdownFile(captures[2])
						if err != nil {
							return nil, err
						}
						ex.Title = captures[1]
						if tc.Teardowns == nil {
							tc.Teardowns = []*executable.Executable{}
						}
						tc.Teardowns = append(tc.Teardowns, ex)
					}
				}
			}
		}
	}

	tc.Title = title
	tc.Description = line.Line(strings.Join(description, "\n")).Trim().Value()
	tc.Request = sutReq

	return tc, nil
}
