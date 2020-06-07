package markdown

import (
	"fmt"
	"regexp"
	"strings"
)

// ElementInterface is interface of element
type ElementInterface interface {
	GetType() string
}

// BaseElement is base of all element
type BaseElement struct {
	Type string
}

// SimpleElement contains just a text inside
type SimpleElement struct {
	*BaseElement
	Text string
}

// TableElement contains table data
type TableElement struct {
	*BaseElement
	Header []string
	Values [][]string
	cursor int
}

// RichTextElement contain text with some markdown inside
type RichTextElement struct {
	*BaseElement
	Text string

	// Anchor inside element ordered by appearance
	Anchor []AnchorElement
}

// AnchorElement is anchor [Text](Link)
type AnchorElement struct {
	Text string
	Link string
}

func tryHeader(data []string, level int) (ElementInterface, bool) {
	if len(data) == 0 {
		return nil, false
	}
	if len(data) > 1 {
		return nil, false
	}
	if level < 1 || level > 6 {
		return nil, false
	}
	re := regexp.MustCompile(fmt.Sprintf("^%s (.+)", strings.Repeat("#", level)))
	m := re.FindStringSubmatch(data[0])
	if len(m) > 1 {
		return &SimpleElement{
			BaseElement: &BaseElement{
				Type: fmt.Sprintf("H%d", level),
			},
			Text: m[1],
		}, true
	}
	return nil, false
}

func tryBullet(data []string) (ElementInterface, bool) {
	if len(data) > 1 {
		return nil, false
	}
	re := regexp.MustCompile("^\\* (.+)")
	m := re.FindStringSubmatch(data[0])
	if len(m) > 1 {
		return &RichTextElement{
			BaseElement: &BaseElement{
				Type: "Bullet",
			},
			Text:   m[1],
			Anchor: extractAnchors(m[1]),
		}, true
	}
	return nil, false
}

func tryCodeBlock(data []string) (ElementInterface, bool) {
	body := strings.Join(data, "\n")
	re := regexp.MustCompile("^(?s)```[^\\n]*\\n(.*)\\n```")
	m := re.FindStringSubmatch(body)
	if len(m) == 0 {
		// Alternate code block
		re = regexp.MustCompile("^(?s)~~~[^\\n]*\\n(.*)\\n~~~")
		m = re.FindStringSubmatch(body)
		if len(m) == 0 {
		}
	}

	if len(m) > 1 {
		return &SimpleElement{
			BaseElement: &BaseElement{
				Type: "Code",
			},
			Text: m[1],
		}, true
	}
	return nil, false
}

func tryTable(data []string) (ElementInterface, bool) {
	if len(data) < 3 {
		return nil, false
	}

	var hCol []string
	var bCol []string
	var ok bool

	if hCol, ok = tryColumn(data[0]); !ok {
		return nil, false
	}
	if bCol, ok = tryColumn(data[1]); !ok {
		return nil, false
	}

	colCount := len(hCol)
	if colCount != len(bCol) {
		fmt.Println("column mismatches.")
		return nil, false
	}

	table := &TableElement{
		BaseElement: &BaseElement{
			Type: "Table",
		},
		Header: hCol,
		Values: [][]string{},
		cursor: -1,
	}
	for i, n := 2, len(data); i < n; i++ {
		var vCol []string
		vCol, ok = tryColumn(data[i])
		if colCount != len(vCol) {
			break
		}
		table.Values = append(table.Values, vCol)
	}

	return table, true
}

func tryColumn(data string) ([]string, bool) {
	if strings.Contains(data, " | ") {
		if data[0] == '|' {
			data = data[1:]
		}
		if data[len(data)-1] == '|' {
			data = data[0 : len(data)-1]
		}
		cols := strings.Split(data, " | ")
		colCount := len(cols)

		for i := 0; i < colCount; i++ {
			cols[i] = strings.TrimSpace(cols[i])
		}

		return cols, true
	}
	if data[0] == '|' && data[len(data)-1] == '|' {
		cols := []string{strings.TrimSpace(data[1 : len(data)-1])}
		return cols, true
	}
	return []string{}, false
}

func extractAnchors(data string) []AnchorElement {
	re := regexp.MustCompile(`\[(.+)\]\(([^\)]+)\)`)
	m := re.FindAllStringSubmatch(data, -1)
	a := []AnchorElement{}
	for _, match := range m {
		a = append(a, AnchorElement{
			Text: match[1],
			Link: match[2],
		})
	}
	return a
}

// NewElement creates a new element implementing ElementInterface
func NewElement(data []string) ElementInterface {
	for i := 1; i <= 6; i++ {
		if elm, ok := tryHeader(data, i); ok {
			return elm
		}
	}

	if elm, ok := tryBullet(data); ok {
		return elm
	}

	if elm, ok := tryCodeBlock(data); ok {
		return elm
	}

	if elm, ok := tryTable(data); ok {
		return elm
	}

	return &SimpleElement{
		BaseElement: &BaseElement{
			Type: "Text",
		},
		Text: data[0],
	}
}

// GetType return type of element
func (be *BaseElement) GetType() string {
	return be.Type
}

// Match matches text with given pattern
func (se *SimpleElement) Match(pattern string) bool {
	re := regexp.MustCompile(pattern)
	return re.MatchString(se.Text)
}

// Capture captures data from text with given pattern
func (se *SimpleElement) Capture(pattern string) ([]string, bool) {
	re := regexp.MustCompile(pattern)
	m := re.FindStringSubmatch(se.Text)
	if len(m) > 1 {
		return m[1:], true
	}
	return []string{}, false
}

// ColumnCount return number of columns
func (te *TableElement) ColumnCount() int {
	return len(te.Header)
}

// MatchHeaders matches header text with given corresponding pattern
func (te *TableElement) MatchHeaders(patterns []string) bool {
	if len(patterns) < len(te.Header) {
		return false
	}

	for i, n := 0, len(patterns); i < n; i++ {
		re := regexp.MustCompile(patterns[i])
		if !re.MatchString(te.Header[i]) {
			return false
		}
	}
	return true
}

// RowCount returns number of rows
func (te *TableElement) RowCount() int {
	return len(te.Values)
}

// Reset resets cursor
func (te *TableElement) Reset() {
	te.cursor = -1
}

// Next move cursor to next
func (te *TableElement) Next() bool {
	te.cursor++
	if te.cursor < len(te.Values) {
		return true
	}
	return false
}

// Value returns current value
func (te *TableElement) Value() []string {
	if te.cursor < len(te.Values) {
		return te.Values[te.cursor]
	}
	return []string{}
}
