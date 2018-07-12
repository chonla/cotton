package testsuite

import (
	"fmt"

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
	req.SetHeaders(t.Headers)
	resp, e := req.Request(url, t.RequestBody)
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
