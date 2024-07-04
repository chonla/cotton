package line

import (
	"errors"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/samber/lo"
)

const (
	windowsLineSeparator    = "\r\n"
	nonWindowsLineSeparator = "\n"
)

type Line string

func DetectLineSeparator() string {
	if runtime.GOOS == "windows" {
		return windowsLineSeparator
	}
	return nonWindowsLineSeparator
}

func FromMultilineString(content string) []Line {
	normalized := normalizeLineSeparator(content)
	lines := strings.Split(normalized, nonWindowsLineSeparator)
	wrappedLines := lo.Map(lines, func(l string, index int) Line {
		return Line(l)
	})
	return wrappedLines
}

func isWindowContent(content string) bool {
	return strings.Contains(content, windowsLineSeparator)
}

func normalizeLineSeparator(content string) string {
	if isWindowContent(content) {
		content = strings.ReplaceAll(content, windowsLineSeparator, nonWindowsLineSeparator)
	}
	return content
}

func (l Line) Trim() Line {
	return Line(strings.TrimSpace(string(l)))
}

func (l Line) Capture(pattern string, index int) (string, bool) {
	captures, ok := l.CaptureAll(pattern)
	if ok {
		if len(captures)-1 < index {
			return "", false
		}
		return captures[index], true
	}
	return "", false
}

func (l Line) CaptureAll(pattern string) ([]string, bool) {
	pat, err := regexp.Compile(pattern)
	if err != nil {
		return nil, false
	}

	matches := pat.FindStringSubmatch(string(l))
	if len(matches) <= 1 {
		return nil, false
	}
	return matches, true
}

func (l Line) LookLike(pattern string) bool {
	pat, err := regexp.Compile(pattern)
	if err != nil {
		return false
	}

	return pat.MatchString(string(l))
}

func (l Line) Value() string {
	return string(l)
}

func (l Line) Replace(needle, with string) string {
	return strings.ReplaceAll(string(l), needle, with)
}

func (l Line) Lower() Line {
	return Line(strings.ToLower(string(l)))
}

func (l Line) StartsWith(partial string) bool {
	return strings.HasPrefix(string(l), partial)
}

func (l Line) ReflectJSValue() (interface{}, error) {
	if cap, ok := l.Capture(`"(.+)"`, 1); ok {
		return cap, nil
	}
	if l.LookLike(`^\d+$`) {
		// ALL numbers in JSON considered a floating point.
		v, err := strconv.ParseFloat(l.Value(), 64)
		if err != nil {
			return nil, err
		}
		return v, nil
	}
	if l.LookLike(`^\d+\.\d+$`) {
		v, err := strconv.ParseFloat(l.Value(), 64)
		if err != nil {
			return nil, err
		}
		return v, nil
	}
	if l.LookLike("true") {
		return true, nil
	}
	if l.LookLike("false") {
		return false, nil
	}
	if l.LookLike("null") {
		return nil, nil
	}
	return nil, errors.New("unexpected value")
}

func (l Line) ReflectRegexValue() (interface{}, error) {
	return regexp.Compile(l.Value())
}

func (l Line) String() string {
	return string(l)
}
