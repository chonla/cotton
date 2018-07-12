package testsuite

import (
	"fmt"
	"strings"

	"github.com/chonla/yas/referrable"
	"github.com/chonla/yas/response"

	"github.com/chonla/yas/request"
	"github.com/fatih/color"
)

// Task is task for setup and teardown
type Task struct {
	Name        string
	Method      string
	BaseURL     string
	Path        string
	ContentType string
	RequestBody string
	Headers     map[string]string
	Captures    map[string]string
	Captured    map[string]string
	Variables   map[string]string
}

// NewTask to create a new task
func NewTask(t *TestCase) *Task {
	return &Task{
		Name:        t.Name,
		Method:      t.Method,
		BaseURL:     t.BaseURL,
		Path:        t.Path,
		ContentType: t.ContentType,
		RequestBody: t.RequestBody,
		Headers:     t.Headers,
		Captures:    t.Captures,
		Variables:   map[string]string{},
		Captured:    map[string]string{},
	}
}

// Run executes test case
func (t *Task) Run() error {
	white := color.New(color.FgHiWhite, color.Bold).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	url := fmt.Sprintf("%s%s", t.BaseURL, t.Path)

	fmt.Printf("%s\n", white(".........."))
	fmt.Printf("Task: %s\n", white(t.Name))
	fmt.Printf("%s\n", white(".........."))

	req, e := request.NewRequester(t.Method)
	if e != nil {
		fmt.Printf("%s: %s\n", red("Error"), e)
		return e
	}
	req.SetHeaders(t.applyVarsToMap(t.Headers))
	resp, e := req.Request(t.applyVars(url), t.applyVars(t.RequestBody))
	if e != nil {
		fmt.Printf("%s: %s\n", red("Error"), e)
		return e
	}

	ref, e := referrable.NewReferrable(response.NewResponse(resp))
	if e != nil {
		fmt.Printf("%s: %s\n", red("Error"), e)
		return e
	}

	for k, v := range t.Captures {
		r, ok := ref.Find(v)
		if ok {
			t.Captured[k] = r[0]
		} else {
			e = fmt.Errorf("unable to capture data from response: %s", k)
			return e
		}
	}

	return e
}

// Value return a value of corresponding captured key
func (t *Task) Value(k string) (string, bool) {
	if v, ok := t.Captured[k]; ok {
		return v, true
	}
	return "", false
}

// MergeVariables merge a given set of vars to local one
func (t *Task) MergeVariables(vars map[string]string) {
	for k, v := range vars {
		t.Variables[k] = v
	}
}

func (t *Task) applyVarsToMap(data map[string]string) map[string]string {
	out := map[string]string{}
	for k, v := range data {
		out[k] = t.applyVars(v)
	}
	return out
}

func (t *Task) applyVars(data string) string {
	for k, v := range t.Variables {
		data = strings.Replace(data, fmt.Sprintf("{%s}", k), v, -1)
	}
	return data
}
