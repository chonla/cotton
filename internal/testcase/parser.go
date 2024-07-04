package testcase

import (
	"cotton/internal/assertion"
	"cotton/internal/capture"
	"cotton/internal/clock"
	"cotton/internal/config"
	"cotton/internal/executable"
	"cotton/internal/httphelper"
	"cotton/internal/line"
	"cotton/internal/logger"
	"cotton/internal/reader"
	"cotton/internal/variable"
	"errors"
	"strings"
)

type ParserOptions struct {
	Configurator     *config.Config
	FileReader       reader.Reader
	RequestParser    httphelper.RequestParser
	ExecutableParser executable.Parser
	Logger           logger.Logger
	ClockWrapper     clock.ClockWrapper
}

type Parser struct {
	options *ParserOptions
}

func NewParser(options *ParserOptions) *Parser {
	return &Parser{
		options: options,
	}
}

func (p *Parser) FromMarkdownFile(mdFileName string) (*Testcase, error) {
	p.options.Logger.PrintDetailedDebugMessage("Parsing", mdFileName)
	mdFullPath := p.options.Configurator.ResolvePath(mdFileName)
	lines, err := p.options.FileReader.Read(mdFullPath)
	if err != nil {
		return nil, err
	}
	return p.FromMarkdownLines(lines)
}

func (p *Parser) FromMarkdownLines(mdLines []line.Line) (*Testcase, error) {
	title := ""
	description := []string{}
	var req []string

	justTitle := false
	collectingCodeBlockBackTick := false
	collectingCodeBlockTilde := false
	discardingCodeBlockBacktick := false
	discardingCodeBlockTilde := false
	titleCollected := false
	reqRaw := ""
	captures := []*capture.Capture{}
	reqFound := false
	assertions := []*assertion.Assertion{}
	setups := []*executable.Executable{}
	teardowns := []*executable.Executable{}
	defaultVars := variable.New()

	p.options.Logger.PrintDetailedDebugMessage("==========")

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

				if cap, ok := mdLine.Capture(`^ {0,3}#\s+(.*)`, 1); ok && !justTitle && !titleCollected {
					p.options.Logger.PrintDetailedDebugMessage("Test title found")
					title = cap
					justTitle = true
					titleCollected = true
					continue
				}

				if collectingCodeBlockBackTick {
					p.options.Logger.PrintDetailedDebugMessage("Collecting code block backtick")

					if ok := mdLine.LookLike("^```$"); ok {
						p.options.Logger.PrintDetailedDebugMessage("End of code block backtick found")
						collectingCodeBlockBackTick = false

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
							justTitle = false
							captures = append(captures, cap)
						} else {
							if defaultVar, ok := variable.Try(mdLine); ok {
								p.options.Logger.PrintDetailedDebugMessage("Variable found")
								justTitle = false
								defaultVars.Add(defaultVar)
							} else {
								if as, ok := assertion.Try(mdLine); ok {
									p.options.Logger.PrintDetailedDebugMessage("Assertion found")
									justTitle = false
									assertions = append(assertions, as)
								} else {
									if captures, ok := mdLine.CaptureAll(`^\s*[\*\-\+]\s\[([^\]]+)\]\(([^\)]+)\)`); ok {
										// unordered list
										justTitle = false
										if !reqFound {
											p.options.Logger.PrintDetailedDebugMessage("Setup found in unordered list")
											ex, err := p.options.ExecutableParser.FromMarkdownFile(captures[2])
											if err != nil {
												return nil, err
											}
											ex.SetTitle(captures[1])
											setups = append(setups, ex)
										} else {
											p.options.Logger.PrintDetailedDebugMessage("Teardown found in unordered list")
											ex, err := p.options.ExecutableParser.FromMarkdownFile(captures[2])
											if err != nil {
												return nil, err
											}
											ex.SetTitle(captures[1])
											teardowns = append(teardowns, ex)
										}
									} else {
										if captures, ok := mdLine.CaptureAll(`^\s*\d+\.\s\[([^\]]+)\]\(([^\)]+)\)`); ok {
											// ordered list
											justTitle = false
											if !reqFound {
												p.options.Logger.PrintDetailedDebugMessage("Setup found in ordered list")
												ex, err := p.options.ExecutableParser.FromMarkdownFile(captures[2])
												if err != nil {
													return nil, err
												}
												ex.SetTitle(captures[1])
												setups = append(setups, ex)
											} else {
												p.options.Logger.PrintDetailedDebugMessage("Teardown found in ordered list")
												ex, err := p.options.ExecutableParser.FromMarkdownFile(captures[2])
												if err != nil {
													return nil, err
												}
												ex.SetTitle(captures[1])
												teardowns = append(teardowns, ex)
											}
										} else {
											if ok := mdLine.LookLike(`^ {0,3}#{1,6}\s+(.*)`); ok {
												p.options.Logger.PrintDetailedDebugMessage("Other titles found, dropped")
												justTitle = false
												// continue
											} else {
												if mdLine.LookLike("^```http$") && !reqFound {
													p.options.Logger.PrintDetailedDebugMessage("HTTP code block backtick found")
													justTitle = false
													collectingCodeBlockBackTick = true
													continue
												}

												if mdLine.LookLike("^~~~http$") && !reqFound {
													p.options.Logger.PrintDetailedDebugMessage("HTTP code block tilde found")
													justTitle = false
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

												if justTitle {
													p.options.Logger.PrintDetailedDebugMessage("Test description found")
													description = append(description, mdLine.Value())
												}
											}
										}
									}
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

	p.options.Logger.PrintDetailedDebugMessage("==========")

	options := &TestcaseOptions{
		RequestParser: p.options.RequestParser,
		Logger:        p.options.Logger,
		ClockWrapper:  p.options.ClockWrapper,
	}
	tc := NewTestcase(title, line.Line(strings.Join(description, "\n")).Trim().Value(), reqRaw, options)
	for _, cap := range captures {
		tc.AddCapture(cap)
	}
	for _, assertion := range assertions {
		tc.AddAssertion(assertion)
	}
	for _, setup := range setups {
		tc.AddSetup(setup)
	}
	for _, teardown := range teardowns {
		tc.AddTeardown(teardown)
	}
	tc.variables = tc.variables.MergeWith(defaultVars)

	return tc, nil
}
