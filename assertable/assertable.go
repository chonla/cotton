package assertable

import (
	"errors"
	"fmt"

	"github.com/chonla/cotton/referrable"
	"github.com/chonla/cotton/response"
	"github.com/fatih/color"
)

// Assertable is something assertable
type Assertable struct {
	*referrable.Referrable
}

// NewAssertable creates an assertable object
func NewAssertable(resp *response.Response) *Assertable {
	ref := referrable.NewReferrable(resp)
	return &Assertable{
		Referrable: ref,
	}
}

// Assert to assert with expectations
func (a *Assertable) Assert(ex map[string]string) error {
	// magenta := color.New(color.FgMagenta, color.Bold).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	blue := color.New(color.FgBlue).SprintFunc()

	if len(ex) == 0 {
		return errors.New("no assertion given")
	}

	for k, v := range ex {
		m := NewMatcher(k, v)
		fmt.Printf("* Assert %s with %s...", blue(k), blue(m))
		r, e := m.Match(a)
		if r {
			fmt.Printf("%s\n", green("PASSED"))
		} else {
			fmt.Printf("%s\n", red("FAILED"))
			return e
		}
	}

	return nil
}
