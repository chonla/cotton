package assertable

import (
	"fmt"
	"regexp"
)

// Matcher is matcher
type Matcher struct {
	reg   *regexp.Regexp
	value string
}

// NewMatcher creates a new matcher
func NewMatcher(v string) *Matcher {
	var m *Matcher
	if isRegExp(v) {
		m = &Matcher{
			reg:   regexp.MustCompile(v[1 : len(v)-1]),
			value: v,
		}
	} else {
		m = &Matcher{
			reg:   nil,
			value: v,
		}
	}
	return m
}

func (m *Matcher) String() string {
	if m.reg != nil {
		return fmt.Sprintf("Regex(%s)", m.value)
	}
	return m.value
}

func isRegExp(v string) bool {
	return len(v) > 2 && v[0] == '/' && v[len(v)-1] == '/'
}

// Match to match value
func (m *Matcher) Match(v string) bool {
	if m.reg != nil {
		return m.reg.MatchString(v)
	}
	return m.value == v
}
