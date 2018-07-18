package parser

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/chonla/cotton/markdown"
	ts "github.com/chonla/cotton/testsuite"
)

var readFileFn = ioutil.ReadFile

// Parser is test parser
type Parser struct{}

// NewParser create a new parser
func NewParser() *Parser {
	return &Parser{}
}

// Parse a test path into []*TestSuite
func (p *Parser) Parse(path string) (*ts.TestSuites, error) {
	suites := &ts.TestSuites{}
	suite := []*ts.TestSuite{}

	files, e := p.listFiles(path)
	if e != nil {
		return suites, e
	}
	for _, f := range files {
		ts, e := p.ParseFile(f)
		if e != nil {
			fmt.Printf("Unable to parse file %s: %s", f, e)
		} else {
			suite = append(suite, ts)
		}
	}
	suites.Suites = suite
	return suites, nil
}

func (p *Parser) listFiles(path string) ([]string, error) {
	var files []string
	e := filepath.Walk(path, p.scan(&files))
	return files, e
}

func (p *Parser) scan(files *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
		}
		if !info.IsDir() && info.Name()[0] == '_' {
			return nil
		}
		if !info.IsDir() && strings.ToLower(filepath.Ext(path)) == ".md" {
			*files = append(*files, path)
		}
		return nil
	}
}

// ParseFile a test file into *TestSuite
func (p *Parser) ParseFile(file string) (*ts.TestSuite, error) {
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
	section := ""
	filePath := filepath.Dir(file)

	md := markdown.NewMD()
	e := md.Parse(file)
	if e != nil {
		return []*ts.TestCase{}, e
	}

	for md.Next() {
		elm := md.Value()

		switch section {
		case "":
			switch elm.GetType() {
			case "H1":
				tc = ts.NewTestCase(elm.(*markdown.SimpleElement).Text)
				section = "suite"
			}
		case "suite":
			switch elm.GetType() {
			case "H1":
				testcases = append(testcases, tc)

				tc = ts.NewTestCase(elm.(*markdown.SimpleElement).Text)
				section = "suite"
			case "H2":
				se := elm.(*markdown.SimpleElement)
				if m, ok := se.Capture("(?i)^(GET|POST|DELETE|PUT|PATCH|OPTIONS) (.+)$"); ok {
					tc.Method = m[0]
					tc.Path = m[1]
					section = "request"
				} else {
					if se.Match("(?i)^expectations?$") {
						section = "expectations"
					} else {
						if se.Match("(?i)^preconditions?$") {
							section = "preconditions"
						} else {
							if se.Match("(?i)^captures?$") {
								section = "captures"
							}
						}
					}
				}
			}
		case "request":
			switch elm.GetType() {
			case "H1":
				testcases = append(testcases, tc)

				tc = ts.NewTestCase(elm.(*markdown.SimpleElement).Text)
				section = "suite"
			case "H2":
				se := elm.(*markdown.SimpleElement)
				if se.Match("(?i)^expectations?$") {
					section = "expectations"
				} else {
					if se.Match("(?i)^preconditions?$") {
						section = "preconditions"
					} else {
						if se.Match("(?i)^captures?$") {
							section = "captures"
						}
					}
				}
			case "Code":
				se := elm.(*markdown.SimpleElement)
				tc.RequestBody = se.Text
			case "Table":
				te := elm.(*markdown.TableElement)
				if te.ColumnCount() == 2 && te.MatchHeaders([]string{"(?i)^header$", "(?i)^value$"}) {
					for te.Next() {
						row := te.Value()
						tc.Headers[row[0]] = row[1]
					}
				}
			}
		case "expectations":
			switch elm.GetType() {
			case "Table":
				te := elm.(*markdown.TableElement)
				if te.ColumnCount() == 2 && te.MatchHeaders([]string{"(?i)^assert$", "(?i)^expected$"}) {
					for te.Next() {
						row := te.Value()
						tc.Expectations[row[0]] = row[1]
					}
				}
			case "H1":
				testcases = append(testcases, tc)

				tc = ts.NewTestCase(elm.(*markdown.SimpleElement).Text)
				section = "suite"
			case "H2":
				se := elm.(*markdown.SimpleElement)
				if m, ok := se.Capture("(?i)^(GET|POST|DELETE|PUT|PATCH|OPTIONS) (.+)$"); ok {
					tc.Method = m[0]
					tc.Path = m[1]
					section = "request"
				} else {
					if se.Match("(?i)^preconditions?$") {
						section = "preconditions"
					} else {
						if se.Match("(?i)^captures?$") {
							section = "captures"
						}
					}
				}
			}
		case "captures":
			switch elm.GetType() {
			case "H1":
				testcases = append(testcases, tc)

				tc = ts.NewTestCase(elm.(*markdown.SimpleElement).Text)
				section = "suite"
			case "Table":
				te := elm.(*markdown.TableElement)
				if te.ColumnCount() == 2 && te.MatchHeaders([]string{"(?i)^name$", "(?i)^value$"}) {
					for te.Next() {
						row := te.Value()
						tc.Captures[row[0]] = row[1]
					}
				}
			case "H2":
				se := elm.(*markdown.SimpleElement)
				if m, ok := se.Capture("(?i)^(GET|POST|DELETE|PUT|PATCH|OPTIONS) (.+)$"); ok {
					tc.Method = m[0]
					tc.Path = m[1]
					section = "request"
				} else {
					if se.Match("(?i)^preconditions?$") {
						section = "preconditions"
					} else {
						if se.Match("(?i)^expectations?$") {
							section = "expectations"
						}
					}
				}
			}
		case "preconditions":
			switch elm.GetType() {
			case "H1":
				testcases = append(testcases, tc)

				tc = ts.NewTestCase(elm.(*markdown.SimpleElement).Text)
				section = "suite"
			case "Bullet":
				se := elm.(*markdown.RichTextElement)
				if len(se.Anchor) > 0 {
					for _, anc := range se.Anchor {
						setupParser := NewParser()
						setup, e := setupParser.Parse(filepath.Clean(fmt.Sprintf("%s/%s", filePath, anc.Link)))
						if e != nil {
							return []*ts.TestCase{}, e
						}
						if len(setup.Suites) > 0 {
							tc.Setups = append(tc.Setups, ts.NewTask(setup.Suites[0].TestCases[0]))
						}
					}
				}
			case "H2":
				se := elm.(*markdown.SimpleElement)
				if m, ok := se.Capture("(?i)^(GET|POST|DELETE|PUT|PATCH|OPTIONS) (.+)$"); ok {
					tc.Method = m[0]
					tc.Path = m[1]
					section = "request"
				} else {
					if se.Match("(?i)^captures?$") {
						section = "captures"
					} else {
						if se.Match("(?i)^expectations?$") {
							section = "expectations"
						}
					}
				}
			}
		}
	}

	if tc != nil {
		testcases = append(testcases, tc)
	}

	// pretty.Println(testcases)

	return testcases, nil
}
