package parser

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/chonla/cotton/assertable"
	"github.com/chonla/cotton/markdown"
	"github.com/chonla/cotton/request"
	ts "github.com/chonla/cotton/testsuite"
)

var readFileFn = ioutil.ReadFile

// Parser is test parser
type Parser struct{}

// Interface is interface of parser
type Interface interface {
	Parse(string) (ts.TestSuitesInterface, error)
	ParseFile(string) (*ts.TestSuite, error)
	ParseString(string, string) ([]*ts.TestCase, error)
}

// NewParser create a new parser
func NewParser() *Parser {
	return &Parser{}
}

// Parse a test path into []*TestSuite
func (p *Parser) Parse(path string) (ts.TestSuitesInterface, error) {
	suites := &ts.TestSuites{
		Variables: map[string]string{},
	}
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
	filePath := filepath.Dir(file)
	content, e := readFileFn(file)
	if e != nil {
		return []*ts.TestCase{}, e
	}

	return p.ParseString(string(content), filePath)
}

// ParseString to parse string content on working path
func (p *Parser) ParseString(content, filePath string) ([]*ts.TestCase, error) {
	testcases := []*ts.TestCase{}
	var tc *ts.TestCase
	section := ""

	md := markdown.NewMD()
	e := md.ParseString(content)
	if e != nil {
		return []*ts.TestCase{}, e
	}

	for md.Next() {
		elm := md.Value()

		switch elm.GetType() {
		case "H1":
			if tc != nil {
				testcases = append(testcases, tc)
			}
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
						} else {
							if se.Match("(?i)^finally$") {
								section = "finally"
							} else {
								section = ""
							}
						}
					}
				}
			}
		case "Code":
			if section == "request" {
				se := elm.(*markdown.SimpleElement)
				tc.RequestBody = se.Text
			}
		case "Table":
			switch section {
			case "request":
				te := elm.(*markdown.TableElement)
				if te.ColumnCount() == 2 && te.MatchHeaders([]string{"(?i)^header$", "(?i)^value$"}) {
					for te.Next() {
						row := te.Value()
						tc.Headers[row[0]] = row[1]
					}
				}
			case "expectations":
				te := elm.(*markdown.TableElement)
				if te.ColumnCount() == 2 && te.MatchHeaders([]string{"(?i)^assert$", "(?i)^expected$"}) {
					for te.Next() {
						row := te.Value()
						tc.Expectations = append(tc.Expectations, assertable.Row{
							Field:       row[0],
							Expectation: row[1],
						})
					}
				}
			case "captures":
				te := elm.(*markdown.TableElement)
				if te.ColumnCount() == 2 && te.MatchHeaders([]string{"(?i)^name$", "(?i)^value$"}) {
					for te.Next() {
						row := te.Value()
						tc.Captures[row[0]] = row[1]
					}
				}
			}
		case "Bullet":
			switch section {
			case "request":
				se := elm.(*markdown.RichTextElement)
				if len(se.Anchor) > 0 {
					for _, anc := range se.Anchor {
						tc.UploadList = append(tc.UploadList, &request.UploadFile{
							FieldName: anc.Text,
							FileName:  anc.Link,
						})
					}
				}
			case "preconditions":
				se := elm.(*markdown.RichTextElement)
				if len(se.Anchor) > 0 {
					for _, anc := range se.Anchor {
						setupParser := NewParser()
						setup, e := setupParser.ParseFile(filepath.Clean(fmt.Sprintf("%s/%s", filePath, anc.Link)))
						if e != nil {
							return []*ts.TestCase{}, e
						}
						if len(setup.TestCases) > 0 {
							tc.Setups = append(tc.Setups, ts.NewTask(setup.TestCases[0]))
						}
					}
				}
			case "finally":
				se := elm.(*markdown.RichTextElement)
				if len(se.Anchor) > 0 {
					for _, anc := range se.Anchor {
						teardownParser := NewParser()
						teardown, e := teardownParser.ParseFile(filepath.Clean(fmt.Sprintf("%s/%s", filePath, anc.Link)))
						if e != nil {
							return []*ts.TestCase{}, e
						}
						if len(teardown.TestCases) > 0 {
							tc.Teardowns = append(tc.Teardowns, ts.NewTask(teardown.TestCases[0]))
						}
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
