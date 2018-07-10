package assertable

import (
	"errors"
	"fmt"
	"strings"

	"github.com/chonla/yas/response"
	"github.com/fatih/color"
	"github.com/kr/pretty"
)

// Assertable is something assertable
type Assertable struct {
	values map[string][]string
}

// NewAssertable creates an assertable object
func NewAssertable(resp *response.Response) *Assertable {
	values := map[string][]string{}

	values["statuscode"] = []string{fmt.Sprintf("%d", resp.StatusCode)}
	values["status"] = []string{resp.Status}

	for k, v := range resp.Header {
		key := strings.ToLower(fmt.Sprintf("header.%s", k))
		if values[key] == nil {
			values[key] = []string{}
		}
		for _, t := range v {
			values[key] = append(values[key], t)
		}
	}

	return &Assertable{
		values: values,
	}
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

	pretty.Println(ex)

	fmt.Printf("%s\n", magenta("----"))
	for k, v := range ex {
		m := NewMatcher(v)
		fmt.Printf("Assert %s with %s...", blue(k), blue(m))
		k = strings.ToLower(k)
		if val, ok := a.values[k]; ok {
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
