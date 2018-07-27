package assertable

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/fatih/color"
)

// Matcher is matcher
type Matcher struct {
	reg     *regexp.Regexp
	builtIn bool
	key     string
	value   string
}

// NewMatcher creates a new matcher
func NewMatcher(k, v string) *Matcher {
	var m *Matcher
	if isRegExp(v) {
		m = &Matcher{
			reg:     regexp.MustCompile(v[1 : len(v)-1]),
			builtIn: false,
			key:     k,
			value:   v,
		}
	} else {
		if isBuiltIn(v) {
			m = &Matcher{
				reg:     nil,
				builtIn: true,
				key:     k,
				value:   v,
			}
		} else {
			m = &Matcher{
				reg:     nil,
				builtIn: false,
				key:     k,
				value:   v,
			}
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

func isBuiltIn(v string) bool {
	return len(v) > 2 && v[0] == '*' && v[len(v)-1] == '*'
}

func isRegExp(v string) bool {
	return len(v) > 2 && v[0] == '/' && v[len(v)-1] == '/'
}

func (m *Matcher) match(v string) bool {
	if m.reg != nil {
		return m.reg.MatchString(v)
	}
	return m.value == v
}

// Match to match value
func (m *Matcher) Match(a *Assertable) (bool, error) {
	red := color.New(color.FgRed).SprintFunc()

	k := strings.ToLower(m.key)
	if val, ok := a.Find(k); ok {
		match := false
		for _, t := range val {
			if m.match(t) {
				match = true
				break
			}
		}
		if match {
			return true, nil
		}
		return false, fmt.Errorf("expect %s in %s, but not", red(m.value), red(m.key))
	}
	return false, fmt.Errorf("response does not contain %s", k)
}
