package executable

import (
	"cotton/internal/capture"
	"cotton/internal/clock"
	"cotton/internal/config"
	"cotton/internal/httphelper"
	"cotton/internal/line"
	"cotton/internal/logger"
	"cotton/internal/reader"
	"cotton/internal/variable"
	"errors"
	"strings"
)

type ParserOptions struct {
	Configurator  *config.Config
	FileReader    reader.Reader
	RequestParser httphelper.RequestParser
	Logger        logger.Logger
	ClockWrapper  clock.ClockWrapper
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
	discardingCodeBlockBacktick := false
	discardingCodeBlockTilde := false

	reqFound := false
	title := "Untitled"
	reqRaw := ""
	captures := []*capture.Capture{}
	defaultVars := variable.New()

	for _, mdLine := range mdLines {
		if discardingCodeBlockBacktick {
			// discard everything after opening unsupport ```
			if ok := mdLine.LookLike("^```$"); ok {
				discardingCodeBlockBacktick = false
			}
		} else {
			if discardingCodeBlockTilde {
				// discard everything after opening unsupport ~~~
				if ok := mdLine.LookLike("^~~~$"); ok {
					discardingCodeBlockTilde = false
				}
			} else {
				if collectingCodeBlockBackTick {
					if ok := mdLine.LookLike("^```$"); ok {
						collectingCodeBlockBackTick = false
						discardingCodeBlockBacktick = false

						if len(req) > 0 {
							reqRaw = line.Line(strings.Join(req, "\n")).Value()
							reqFound = true
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
								reqRaw = line.Line(strings.Join(req, "\n")).Value()
								reqFound = true
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
							captures = append(captures, cap)
						} else {
							if defaultVar, ok := variable.Try(mdLine); ok {
								defaultVars.Add(defaultVar)
							} else {
								if mdLine.LookLike("^```http$") && !reqFound {
									collectingCodeBlockBackTick = true
									continue
								}

								if mdLine.LookLike("^~~~http$") && !reqFound {
									collectingCodeBlockTilde = true
									continue
								}

								if ok := mdLine.LookLike("^```"); ok {
									discardingCodeBlockBacktick = true
									continue
								}

								if ok := mdLine.LookLike("^~~~"); ok {
									discardingCodeBlockTilde = true
									continue
								}
							}
						}
					}
				}
			}
		}
	}

	if !reqFound {
		return nil, errors.New("no callable request")
	}

	options := &ExecutableOptions{
		RequestParser: p.options.RequestParser,
		Logger:        p.options.Logger,
	}
	ex := New(title, reqRaw, options)
	for _, cap := range captures {
		ex.AddCapture(cap)
	}
	ex.variables = ex.variables.MergeWith(defaultVars)

	return ex, nil
}
