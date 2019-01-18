package testsuite

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/chonla/cotton/assertable"
	"github.com/chonla/cotton/referrable"
	"github.com/chonla/cotton/request"
	"github.com/chonla/cotton/response"
	"github.com/fatih/color"
)

// TestCase holds a test case
type TestCase struct {
	Name         string
	Method       string
	BaseURL      string
	Config       *Config
	Path         string
	ContentType  string
	RequestBody  string
	Headers      map[string]string
	Expectations []assertable.Row
	Captures     map[string]string
	Setups       []*Task
	Teardowns    []*Task
	Variables    map[string]string
	Captured     map[string]string
	Cookies      []*http.Cookie
}

// NewTestCase creates a new testcase
func NewTestCase(name string) *TestCase {
	return &TestCase{
		Name:         name,
		Headers:      map[string]string{},
		Expectations: []assertable.Row{},
		Config: &Config{
			Insecure: false,
			Detail:   false,
		},
		Captures:  map[string]string{},
		Setups:    []*Task{},
		Teardowns: []*Task{},
		Variables: map[string]string{},
		Captured:  map[string]string{},
		Cookies:   []*http.Cookie{},
	}
}

// SetContentType set a corresponding content type
func (tc *TestCase) SetContentType(ct string) {
	switch strings.ToLower(ct) {
	case "json":
		ct = "application/json"
	default:
		ct = "application/json"
	}
	tc.ContentType = ct
}

// Run executes test case
func (tc *TestCase) Run() error {
	// Skip if no assertion
	if len(tc.Expectations) == 0 {
		return nil
	}

	white := color.New(color.FgHiWhite, color.Bold).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	url := fmt.Sprintf("%s%s", tc.BaseURL, tc.Path)
	blue := color.New(color.FgBlue).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	grey := color.New(color.FgWhite, color.Faint).SprintfFunc()

	fmt.Printf("%s\n", white("================================================================================"))
	fmt.Printf("Testcase: %s\n", white(tc.Name))
	fmt.Printf("%s\n", white("================================================================================"))

	tc.Cookies = []*http.Cookie{}

	if len(tc.Setups) > 0 {
		fmt.Printf("Setup:\n")
		for _, s := range tc.Setups {

			fmt.Printf("* %s...", blue(s.Name))

			s.BaseURL = tc.BaseURL
			s.Config = tc.Config
			s.MergeVariables(tc.Variables)

			if len(tc.Cookies) > 0 {
				s.SetCookies(tc.Cookies)
			}

			if s.Config.Detail {
				fmt.Println()
			}

			e := s.Run()

			if s.Config.Detail {
				fmt.Printf("...%s...", grey(fmt.Sprintf("(%s)", s.Name)))
			}

			if e != nil {
				fmt.Printf("%s: %s\n", red("FAILED"), e)
				return e
			}
			fmt.Printf("%s\n", green("PASSED"))

			for k, v := range s.Captured {
				tc.Variables[k] = v
			}

			if len(s.Cookies) > 0 {
				for _, v := range s.Cookies {
					tc.Cookies = append(tc.Cookies, v)
				}
			}
		}

		fmt.Println()
	}

	req, e := request.NewRequester(tc.Method, tc.Config.Insecure, tc.Config.Detail)
	if e != nil {
		return e
	}

	req.SetHeaders(tc.applyVarsToMap(tc.Headers))
	req.SetCookies(tc.Cookies)

	targetURL := tc.applyVars(url)
	fmt.Printf("Action: %s %s\n", white(tc.Method), yellow(targetURL))
	resp, e := req.Request(targetURL, tc.applyVars(tc.RequestBody))
	if e != nil {
		fmt.Printf("%s: %s\n", red("Error"), e)
		return e
	}

	r := response.NewResponse(resp, tc.Config.Detail)

	if tc.Config.Detail {
		r.LogResponse()
	}

	if len(tc.Captures) > 0 {
		ref := referrable.NewReferrable(r)

		for k, v := range tc.Captures {
			r, ok := ref.Find(v)
			if ok {
				tc.Captured[k] = r[0]
			} else {
				e = fmt.Errorf("unable to capture data from response: %s", k)
				return e
			}
		}

		for k, v := range tc.Captured {
			tc.Variables[k] = v
		}
	}

	if len(r.Cookies) > 0 {
		for _, v := range r.Cookies {
			tc.Cookies = append(tc.Cookies, v)
		}
	}

	as := assertable.NewAssertable(r)

	assertResult := as.Assert(tc.Expectations)

	if len(tc.Teardowns) > 0 {
		fmt.Printf("\nTeardown:\n")
		for _, s := range tc.Teardowns {

			fmt.Printf("* %s...", blue(s.Name))

			s.BaseURL = tc.BaseURL
			s.Config = tc.Config
			s.MergeVariables(tc.Variables)

			if len(tc.Cookies) > 0 {
				s.SetCookies(tc.Cookies)
			}

			if s.Config.Detail {
				fmt.Println()
			}

			e := s.Run()

			if s.Config.Detail {
				fmt.Printf("...%s...", grey(fmt.Sprintf("(%s)", s.Name)))
			}

			if e != nil {
				fmt.Printf("%s: %s\n", red("FAILED"), e)
				return e
			}
			fmt.Printf("%s\n", green("PASSED"))

			for k, v := range s.Captured {
				tc.Variables[k] = v
			}

			if len(s.Cookies) > 0 {
				for _, v := range s.Cookies {
					tc.Cookies = append(tc.Cookies, v)
				}
			}

		}
	}

	return assertResult
}

func (tc *TestCase) applyVarsToMap(data map[string]string) map[string]string {
	out := map[string]string{}
	for k, v := range data {
		out[k] = tc.applyVars(v)
	}
	return out
}

func (tc *TestCase) applyVars(data string) string {
	for k, v := range tc.Variables {
		data = strings.Replace(data, fmt.Sprintf("{%s}", k), v, -1)
	}
	return data
}
