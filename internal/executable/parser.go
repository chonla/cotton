package executable

import (
	"cotton/internal/capture"
	"cotton/internal/config"
	"cotton/internal/line"
	"cotton/internal/reader"
	"cotton/internal/request"
	"net/http"
	"strings"

	"github.com/chonla/httpreqparser"
)

type Parser struct {
	config        *config.Config
	fileReader    reader.Reader
	requestParser httpreqparser.Parser
}

func NewParser(config *config.Config, fileReader reader.Reader, requestParser httpreqparser.Parser) *Parser {
	return &Parser{
		config:        config,
		fileReader:    fileReader,
		requestParser: requestParser,
	}
}

func (p *Parser) FromMarkdownFile(mdFileName string) (*Executable, error) {
	mdFullPath := p.config.ResolvePath(mdFileName)
	lines, err := p.fileReader.Read(mdFullPath)
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
					requestString := line.Line(strings.Join(req, "\n")).Value()
					httpRequest, err := p.requestParser.Parse(requestString)
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
			if cap, ok := capture.Try(mdLine); ok {
				if ex.Captures == nil {
					ex.Captures = []*capture.Capture{}
				}
				ex.Captures = append(ex.Captures, cap)
			}
		}
	}

	wrappedReq, err := request.New(exReq)
	if err != nil {
		return nil, err
	}
	ex.Request = wrappedReq

	return ex, nil
}
