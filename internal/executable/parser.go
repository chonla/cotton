package executable

import (
	"cotton/internal/capture"
	"cotton/internal/config"
	"cotton/internal/httphelper"
	"cotton/internal/line"
	"cotton/internal/logger"
	"cotton/internal/reader"
	"errors"
	"strings"
)

type ParserOptions struct {
	Configurator  *config.Config
	FileReader    reader.Reader
	RequestParser httphelper.RequestParser
	Logger        logger.Logger
}

type Parser interface {
	FromMarkdownFile(mdFileName string) (*Executable, error)
	FromMarkdownLines(mdLines []line.Line) (*Executable, error)
}

type ExecutableParser struct {
	options *ParserOptions
}

func NewParser(options *ParserOptions) *ExecutableParser {
	return &ExecutableParser{
		options: options,
	}
}

func (p *ExecutableParser) FromMarkdownFile(mdFileName string) (*Executable, error) {
	mdFullPath := p.options.Configurator.ResolvePath(mdFileName)
	lines, err := p.options.FileReader.Read(mdFullPath)
	if err != nil {
		return nil, err
	}
	return p.FromMarkdownLines(lines)
}

func (p *ExecutableParser) FromMarkdownLines(mdLines []line.Line) (*Executable, error) {
	var req []string

	collectingCodeBlockBackTick := false
	collectingCodeBlockTilde := false

	exReqFound := false
	exTitle := "Untitled"
	exReqRaw := ""
	exCaptures := []*capture.Capture{}

	for _, mdLine := range mdLines {
		if mdLine.LookLike("^```http$") && !exReqFound {
			collectingCodeBlockBackTick = true
			continue
		}

		if mdLine.LookLike("^~~~http$") && !exReqFound {
			collectingCodeBlockTilde = true
			continue
		}

		if collectingCodeBlockBackTick {
			if ok := mdLine.LookLike("^```$"); ok {
				collectingCodeBlockBackTick = false

				if len(req) > 0 {
					exReqRaw = line.Line(strings.Join(req, "\n")).Value()
					exReqFound = true
					req = nil
				}
			} else {
				if req == nil {
					req = []string{}
				}
				req = append(req, mdLine.Value())
			}
		} else {
			if collectingCodeBlockTilde {
				if ok := mdLine.LookLike("^~~~$"); ok {
					collectingCodeBlockTilde = false

					if len(req) > 0 {
						exReqRaw = line.Line(strings.Join(req, "\n")).Value()
						exReqFound = true
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
					exCaptures = append(exCaptures, cap)
				}
			}
		}
	}

	if !exReqFound {
		return nil, errors.New("no callable request")
	}

	options := &ExecutableOptions{
		RequestParser: p.options.RequestParser,
		Logger:        p.options.Logger,
	}
	ex := New(exTitle, exReqRaw, options)
	for _, cap := range exCaptures {
		ex.AddCapture(cap)
	}

	return ex, nil
}
