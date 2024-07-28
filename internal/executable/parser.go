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
	Configurator    *config.Config
	FileReader      reader.Reader
	RequestParser   httphelper.RequestParser
	Logger          logger.Logger
	ClockWrapper    clock.ClockWrapper
	InsecureRequest bool
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
	p.options.Logger.PrintDetailedDebugMessage("Parsing", mdFileName)
	lines, err := p.options.FileReader.Read(mdFileName)
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

	p.options.Logger.PrintDetailedDebugMessage("----------")
	for _, mdLine := range mdLines {
		p.options.Logger.PrintDetailedDebugMessage("Line data", mdLine.Value())
		if discardingCodeBlockBacktick {
			p.options.Logger.PrintDetailedDebugMessage("Discarding code block backtick")
			// discard everything after opening unsupport ```
			if ok := mdLine.LookLike("^```$"); ok {
				p.options.Logger.PrintDetailedDebugMessage("End of code block backtick found")
				discardingCodeBlockBacktick = false
			}
		} else {
			if discardingCodeBlockTilde {
				p.options.Logger.PrintDetailedDebugMessage("Discarding code block tilde")
				// discard everything after opening unsupport ~~~
				if ok := mdLine.LookLike("^~~~$"); ok {
					p.options.Logger.PrintDetailedDebugMessage("End of code block tilde found")
					discardingCodeBlockTilde = false
				}
			} else {
				if collectingCodeBlockBackTick {
					p.options.Logger.PrintDetailedDebugMessage("Collecting code block backtick")
					if ok := mdLine.LookLike("^```$"); ok {
						p.options.Logger.PrintDetailedDebugMessage("End of code block backtick found")
						collectingCodeBlockBackTick = false
						discardingCodeBlockBacktick = false

						if len(req) > 0 {
							reqRaw = line.Line(strings.Join(req, "\n")).Value()
							reqFound = true
							req = nil
						}
					} else {
						p.options.Logger.PrintDetailedDebugMessage("Collecting request")
						if req == nil {
							req = []string{}
						}
						req = append(req, mdLine.Value())
					}
				} else {
					if collectingCodeBlockTilde {
						p.options.Logger.PrintDetailedDebugMessage("Collecting code block tilde")

						if ok := mdLine.LookLike("^~~~$"); ok {
							p.options.Logger.PrintDetailedDebugMessage("End of code block tilde found")
							collectingCodeBlockTilde = false

							if len(req) > 0 {
								p.options.Logger.PrintDetailedDebugMessage("Request available, store request")
								reqRaw = line.Line(strings.Join(req, "\n")).Value()
								reqFound = true
								req = nil
							}
						} else {
							p.options.Logger.PrintDetailedDebugMessage("Collecting request")
							if req == nil {
								req = []string{}
							}
							req = append(req, mdLine.Value())
						}
					} else {
						if cap, ok := capture.Try(mdLine); ok {
							p.options.Logger.PrintDetailedDebugMessage("Capture found")
							captures = append(captures, cap)
						} else {
							if defaultVar, ok := variable.Try(mdLine); ok {
								p.options.Logger.PrintDetailedDebugMessage("Variable found")
								defaultVars.Add(defaultVar)
							} else {
								if mdLine.LookLike("^```http$") && !reqFound {
									p.options.Logger.PrintDetailedDebugMessage("HTTP code block backtick found")
									collectingCodeBlockBackTick = true
									continue
								}

								if mdLine.LookLike("^~~~http$") && !reqFound {
									p.options.Logger.PrintDetailedDebugMessage("HTTP code block tilde found")
									collectingCodeBlockTilde = true
									continue
								}

								if ok := mdLine.LookLike("^```"); ok {
									p.options.Logger.PrintDetailedDebugMessage("Unsupport code block backtick found")
									discardingCodeBlockBacktick = true
									continue
								}

								if ok := mdLine.LookLike("^~~~"); ok {
									p.options.Logger.PrintDetailedDebugMessage("Unsupport code block tilde found")
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
		RequestParser:   p.options.RequestParser,
		Logger:          p.options.Logger,
		InsecureRequest: p.options.InsecureRequest,
	}
	ex := New(title, reqRaw, options)
	for _, cap := range captures {
		ex.AddCapture(cap)
	}
	ex.variables = ex.variables.MergeWith(defaultVars)

	p.options.Logger.PrintDetailedDebugMessage("----------")

	return ex, nil
}
