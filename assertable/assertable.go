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

// Row is assertion entry
type Row struct {
	Field       string
	Expectation string
}

// NewAssertable creates an assertable object
func NewAssertable(resp *response.Response) *Assertable {
	ref := referrable.NewReferrable(resp)

	return &Assertable{
		Referrable: ref,
	}
}

// Assert to assert with expectations
func (a *Assertable) Assert(ex []Row) error {
	// magenta := color.New(color.FgMagenta, color.Bold).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	blue := color.New(color.FgBlue).SprintFunc()

	if len(ex) == 0 {
		return errors.New("no assertion given")
	}

	for _, row := range ex {
		m := NewMatcher(row.Field, row.Expectation)
		fmt.Printf("* Assert %s %s...", blue(row.Field), m)
		r, e := m.Match(a)
		if r {
			fmt.Printf("%s\n", green("PASSED"))
		} else {
			fmt.Printf("%s\n", red("FAILED"))
			fmt.Printf("    %s\n", e)
			return e
		}
	}

	return nil
}
