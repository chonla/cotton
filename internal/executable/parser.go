package executable

import (
	"cotton/internal/capture"
	"cotton/internal/config"
	"cotton/internal/line"
	"cotton/internal/reader"
	"cotton/internal/request"
	"fmt"
	"net/http"
	"strings"

	"github.com/chonla/httpreqparser"
	"github.com/kr/pretty"
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
			fmt.Println("START COLLECTING CODE BLOCK")
			collectingCodeBlockBackTick = true
			continue
		}

		if collectingCodeBlockBackTick {
			if ok := mdLine.LookLike("^```$"); ok {
				fmt.Println("END COLLECTING CODE BLOCK")
				collectingCodeBlockBackTick = false

				if len(req) > 0 {
					requestString := line.Line(strings.Join(req, "\n")).Value()
					fmt.Println("=============")
					pretty.Println(requestString)
					fmt.Println("=============")
					httpRequest, err := p.requestParser.Parse(requestString)
					if err == nil {
						exReq = httpRequest
						// pretty.Println(exReq)
					}
					req = nil
				}
			} else {
				if req == nil {
					req = []string{}
				}
				req = append(req, mdLine.Value())
				fmt.Println(mdLine.Value())
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
