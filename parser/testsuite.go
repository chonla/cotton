package parser

import (
	"errors"
	"io/ioutil"
	"regexp"
	"strings"

	ts "github.com/chonla/yas/testsuite"
)

var readFileFn = ioutil.ReadFile

// Parser is test parser
type Parser struct {
	state    string
	substate string
}

// NewParser create a new parser
func NewParser() *Parser {
	return &Parser{}
}

// Parse a test file into *TestSuite
func (p *Parser) Parse(file string) (*ts.TestSuite, error) {
	testcases, e := p.parseTestSuiteFile(file)
	if e != nil {
		return nil, e
	}

	suite := &ts.TestSuite{
		Name:      p.parseTestSuiteFileName(file),
		TestCases: testcases,
	}
	return suite, nil
}

func (p *Parser) parseTestSuiteFileName(file string) string {
	re := regexp.MustCompile("(?i).+\\.md")
	if re.MatchString(file) {
		file = file[0 : len(file)-3]
	}
	return p.parseTestSuiteName(file)
}

func (p *Parser) parseTestSuiteFile(file string) ([]*ts.TestCase, error) {
	testcases := []*ts.TestCase{}
	var tc *ts.TestCase
	blockData := []string{}

	lines, e := p.readTestSuiteFile(file)
	if e != nil {
		return []*ts.TestCase{}, e
	}

	for _, line := range lines {
		if p.state == "request" {
			if p.substate == "openblock" {
				if _, ok := isCodeBlock(line); ok {
					p.substate = "closeblock"
					tc.RequestBody = strings.Join(blockData, "\n")
					blockData = []string{}
				} else {
					blockData = append(blockData, line)
				}
			} else {
				if title, ok := isTitle(line); ok {
					if tc != nil {
						testcases = append(testcases, tc)
					}
					tc = &ts.TestCase{
						Name: title,
					}
					p.substate = "title"
				} else {
					if method, path, ok := isMethod(line); ok {
						tc.Method = strings.ToUpper(method)
						tc.Path = path
						p.substate = "method"
					} else {
						if blocktype, ok := isCodeBlock(line); ok {
							p.substate = "openblock"
							tc.SetContentType(blocktype)
						} else {
							if isExpectation(line) {
								p.state = "expectation"
							} else {
								// nothing to do now
								if p.substate == "tableheader" {
									if isRequestHeaderEnd(line) {
										p.substate = "tableheaderend"
									} else {
										return nil, errors.New("unexpected line found. expect expectation table ending here: " + line)
									}
								} else {
									if p.substate == "tableheaderend" {
										if item, value, ok := isRequestHeaderContent(line); ok {
											tc.Headers[item] = value
										}
									} else {
										p.substate = ""

										if isRequestHeaderTableHeader(line) {
											p.substate = "tableheader"
										}
									}
								}
							}
						}
					}
				}
			}
		} else {
			if p.state == "expectation" {
				if p.substate == "tableheader" {
					if isExpectationTableHeaderEnd(line) {
						p.substate = "tableheaderend"
					} else {
						return nil, errors.New("unexpected line found. expect expectation table ending here: " + line)
					}
				} else {
					if p.substate == "tableheaderend" {
						if item, value, ok := isExpectationTableContent(line); ok {
							tc.Expectations[item] = value
						} else {
							if title, ok := isTitle(line); ok {
								if tc != nil {
									testcases = append(testcases, tc)
								}
								tc = ts.NewTestCase(title)
								p.state = "request"
								p.substate = "title"
							}
						}
					} else {
						// table ended
						p.substate = ""

						if isExpectationTableHeader(line) {
							p.substate = "tableheader"
						}
					}
				}
			} else {
				if p.state == "" {
					if title, ok := isTitle(line); ok {
						tc = ts.NewTestCase(title)
						p.state = "request"
						p.substate = "title"
					} else {
						return nil, errors.New("unexpected line found. expect testcase name here: " + line)
					}
				}
			}
		}
	}
	if tc != nil {
		testcases = append(testcases, tc)
	}

	return testcases, nil
}

func isCodeBlock(line string) (string, bool) {
	re := regexp.MustCompile("^```(.*)$")
	matches := re.FindStringSubmatch(line)
	if len(matches) > 1 {
		return matches[1], true
	}
	return "", false
}

func isTitle(line string) (string, bool) {
	re := regexp.MustCompile("^# (.+)$")
	matches := re.FindStringSubmatch(line)
	if len(matches) > 1 {
		return matches[1], true
	}
	return "", false
}

func isMethod(line string) (string, string, bool) {
	re := regexp.MustCompile("(?i)^## (GET|POST|DELETE|PUT|PATCH|OPTIONS) (.+)$")
	matches := re.FindStringSubmatch(line)
	if len(matches) > 2 {
		return matches[1], matches[2], true
	}
	return "", "", false
}

func isExpectation(line string) bool {
	re := regexp.MustCompile("(?i)^## Expectation$")
	return re.MatchString(line)
}

func isExpectationTableHeader(line string) bool {
	re := regexp.MustCompile("(?i)^\\|\\s+Assert\\s+\\|\\s+Expected\\s+\\|$")
	return re.MatchString(line)
}

func isExpectationTableHeaderEnd(line string) bool {
	re := regexp.MustCompile("(?i)^\\|\\s+\\-+\\s+\\|\\s+\\-+\\s+\\|$")
	return re.MatchString(line)
}

func isExpectationTableContent(line string) (string, string, bool) {
	re := regexp.MustCompile("(?i)^\\|\\s+([^\\|]+)\\s+\\|\\s+([^\\|]+)\\s+\\|$")
	matches := re.FindStringSubmatch(line)
	if len(matches) > 2 {
		return matches[1], matches[2], true
	}
	return "", "", false
}

func isRequestHeaderEnd(line string) bool {
	re := regexp.MustCompile("(?i)^\\|\\s+\\-+\\s+\\|\\s+\\-+\\s+\\|$")
	return re.MatchString(line)
}

func isRequestHeaderTableHeader(line string) bool {
	re := regexp.MustCompile("(?i)^\\|\\s+Header\\s+\\|\\s+Value\\s+\\|$")
	return re.MatchString(line)
}

func isRequestHeaderContent(line string) (string, string, bool) {
	re := regexp.MustCompile("(?i)^\\|\\s+([^\\|]+)\\s+\\|\\s+([^\\|]+)\\s+\\|$")
	matches := re.FindStringSubmatch(line)
	if len(matches) > 2 {
		return matches[1], matches[2], true
	}
	return "", "", false
}
