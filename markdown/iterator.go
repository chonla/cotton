package markdown

import (
	"fmt"
	"regexp"
	"strings"
)

// Iterator is an iterator for lines to element
type Iterator struct {
	cursor int
	lines  []string
	length int
}

// NewIterator creates an iterator
func NewIterator(lines []string) *Iterator {
	return &Iterator{
		cursor: 0,
		lines:  lines,
		length: len(lines),
	}
}

// Reset cursor to beginning
func (i *Iterator) Reset() {
	i.cursor = 0
}

// Next return if anything in iterator
func (i *Iterator) Next() bool {
	if i.cursor >= i.length {
		return false
	}

	return true
}

// Value return content in current iterator
func (i *Iterator) Value() []string {
	buffer := []string{}

	if i.isHeader(i.lines[i.cursor]) {
		buffer = append(buffer, i.lines[i.cursor])
		i.cursor++
		return buffer
	}

	if i.isCodeBlock(i.lines[i.cursor]) {
		buffer = append(buffer, i.lines[i.cursor])
		i.cursor++
		for i.cursor < i.length && !i.isCodeBlock(i.lines[i.cursor]) {
			buffer = append(buffer, i.lines[i.cursor])
			i.cursor++
		}
		if i.cursor < i.length {
			buffer = append(buffer, i.lines[i.cursor])
		}
		i.cursor++
		return buffer
	}

	if colCount, ok := i.isTable(i.lines[i.cursor]); ok {
		if i.cursor+1 < i.length && i.isTableWithColumnCount(i.lines[i.cursor+1], colCount) {
			buffer = append(buffer, i.lines[i.cursor])
			buffer = append(buffer, i.lines[i.cursor+1])
			i.cursor += 2
			for i.cursor < i.length && i.isTableWithColumnCount(i.lines[i.cursor], colCount) {
				buffer = append(buffer, i.lines[i.cursor])
				i.cursor++
			}
			return buffer
		}
	}

	// default
	buffer = append(buffer, i.lines[i.cursor])
	i.cursor++
	return buffer
}

func (i *Iterator) isHeader(line string) bool {
	re := regexp.MustCompile("^#{1,6} .+")
	return re.MatchString(line)
}

func (i *Iterator) isCodeBlock(line string) bool {
	re := regexp.MustCompile("^```")
	return re.MatchString(line)
}

func (i *Iterator) isTable(line string) (int, bool) {
	re := regexp.MustCompile("^\\|(?:\\s+[^\\|]+\\s+\\|)+")
	if re.MatchString(line) {
		return len(strings.Split(line, "|")) - 2, true
	}
	return 0, false
}

func (i *Iterator) isTableWithColumnCount(line string, count int) bool {
	re := regexp.MustCompile(fmt.Sprintf("^\\|(?:\\s+[^\\|]+\\s+\\|){%d}", count))
	return re.MatchString(line)
}

func (i *Iterator) isBullet(line string) bool {
	re := regexp.MustCompile("^[\\*] .+")
	if re.MatchString(line) {
		return true
	}
	re = regexp.MustCompile("^\\d+\\. .+")
	return re.MatchString(line)
}
