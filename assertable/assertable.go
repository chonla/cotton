package assertable

import (
	"errors"
	"fmt"
	"strings"

	"github.com/chonla/yas/referrable"
	"github.com/chonla/yas/response"
	"github.com/fatih/color"
)

// Assertable is something assertable
type Assertable struct {
	*referrable.Referrable
}

// NewAssertable creates an assertable object
func NewAssertable(resp *response.Response) (*Assertable, error) {
	ref, e := referrable.NewReferrable(resp)
	if e != nil {
		return nil, e
	}
	return &Assertable{
		Referrable: ref,
	}, nil
}

// Assert to assert with expectations
func (a *Assertable) Assert(ex map[string]string) error {
	magenta := color.New(color.FgMagenta, color.Bold).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	blue := color.New(color.FgBlue).SprintFunc()

	if len(ex) == 0 {
		return errors.New("no assertion given")
	}

	fmt.Printf("%s\n", magenta("----"))
	for k, v := range ex {
		m := NewMatcher(v)
		fmt.Printf("Assert %s with %s...", blue(k), blue(m))
		k = strings.ToLower(k)
		if val, ok := a.Find(k); ok {
			match := false
			for _, t := range val {
				if m.Match(t) {
					match = true
					break
				}
			}
			if match {
				fmt.Printf("%s\n", green("PASS"))
			} else {
				fmt.Printf("%s\n", red("FAILED"))
				return fmt.Errorf("expect %s in %s, but not", red(v), red(k))
			}
		} else {
			fmt.Printf("%s\n", red("FAILED"))
			return fmt.Errorf("response does not contain %s", k)
		}
	}

	return nil
}
