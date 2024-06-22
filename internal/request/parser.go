package request

import (
	"net/http/httputil"

	"github.com/chonla/httpreqparser"
)

type Parser struct {
	parser httpreqparser.HttpParser
}

func NewParser(parser httpreqparser.HttpParser) *Parser {
	return &Parser{parser}
}

func (p *Parser) Parse(reqString string) (*HTTPRequest, error) {
	req, err := p.parser.Parse(reqString)
	if err != nil {
		return nil, err
	}

	req.Close = true
	reqBytes, err := httputil.DumpRequest(req, true)
	if err != nil {
		return nil, err
	}
	return &HTTPRequest{
		request:      req,
		plainRequest: reqBytes,
	}, nil
}
