package testsuite

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/chonla/cotton/referrable"
	"github.com/chonla/cotton/response"

	"github.com/chonla/cotton/request"
	"github.com/fatih/color"
)

// Task is task for setup and teardown
type Task struct {
	Name        string
	Method      string
	BaseURL     string
	Config      *Config
	Path        string
	ContentType string
	RequestBody string
	Headers     map[string]string
	Captures    map[string]string
	Captured    map[string]string
	Variables   map[string]string
	Cookies     []*http.Cookie
	UploadList  request.UploadFiles
}

// NewTask to create a new task
func NewTask(t *TestCase) *Task {
	return &Task{
		Name:    t.Name,
		Method:  t.Method,
		BaseURL: t.BaseURL,
		Config: &Config{
			Insecure: false,
			Detail:   false,
		},
		Path:        t.Path,
		ContentType: t.ContentType,
		RequestBody: t.RequestBody,
		Headers:     t.Headers,
		Captures:    t.Captures,
		Variables:   map[string]string{},
		Captured:    map[string]string{},
		Cookies:     []*http.Cookie{},
		UploadList:  t.UploadList,
	}
}

// Run executes test case
func (t *Task) Run() error {
	red := color.New(color.FgRed).SprintFunc()
	url := fmt.Sprintf("%s%s", t.BaseURL, t.Path)

	req, e := request.NewRequester(t.Method, t.Config.Insecure, t.Config.Detail)
	if e != nil {
		fmt.Printf("%s: %s\n", red("Error"), e)
		return e
	}
	req.SetHeaders(t.applyVarsToMap(t.Headers))
	req.SetCookies(t.Cookies)

	reqBody := t.applyVars(t.RequestBody)
	if len(t.UploadList) > 0 {
		uploadRequest, e := t.UploadList.ToRequestBody()
		if e != nil {
			return e
		}

		reqBody = uploadRequest.RequestBody
		req.SetHeaders(map[string]string{
			"Content-Type":   uploadRequest.ContentType,
			"Content-Length": fmt.Sprintf("%d", len(reqBody)),
		})
	}

	resp, e := req.Request(t.applyVars(url), reqBody)
	if e != nil {
		fmt.Printf("%s: %s\n", red("Error"), e)
		return e
	}

	r := response.NewResponse(resp, t.Config.Detail)
	if t.Config.Detail {
		r.LogResponse()
	}

	ref := referrable.NewReferrable(r)

	t.Cookies = []*http.Cookie{}

	for _, v := range r.Cookies {
		t.Cookies = append(t.Cookies, v)
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

func (t *Task) SetCookies(c []*http.Cookie) {
	for _, v := range c {
		t.Cookies = append(t.Cookies, v)
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
