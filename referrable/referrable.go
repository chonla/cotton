package referrable

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/chonla/cotton/response"
	"github.com/fatih/color"
	"github.com/tidwall/gjson"
)

// Referrable is referrable items
type Referrable struct {
	values map[string][]string
	data   gjson.Result
}

// Any type
type Any interface{}

// NewReferrable creates an referrable object
func NewReferrable(resp *response.Response) *Referrable {
	red := color.New(color.FgRed).SprintFunc()
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

	var jsonObj gjson.Result
	var e error
	if isJSONContent(values["header.content-type"]) {
		jsonObj, e = tryParse(resp.Body)
		if e != nil {
			fmt.Printf("%s: %s\n", red("Error"), e)
			jsonObj = gjson.Parse("{}")
		}
	} else {
		jsonObj = gjson.Parse("{}")
	}

	return &Referrable{
		values: values,
		data:   jsonObj,
	}
}

func tryParse(jsonString string) (gjson.Result, error) {
	jsonString = fmt.Sprintf("{ \"document\": %s }", jsonString)
	if gjson.Valid(jsonString) {
		return gjson.Parse(jsonString), nil
	}
	return gjson.Result{}, errors.New("not a well-formed json")
}

func isJSONContent(contenttype []string) bool {
	for _, r := range contenttype {
		token := strings.Split(strings.ToLower(r), ";")
		if token[0] == "application/json" {
			return true
		}
	}
	return false
}

// find is internal finder
func (a *Referrable) find(k string) (gjson.Result, error) {
	re := regexp.MustCompile("(?i)^data\\.(.+)")
	match := re.FindStringSubmatch(k)
	if len(match) > 1 {
		key := fmt.Sprintf("document.%s", match[1])
		key = a.convertToGJsonPath(key)
		if a.data.Get(key).Exists() {
			return a.data.Get(key), nil
		}
		return gjson.Result{}, errors.New("not found")
	}
	re = regexp.MustCompile("(?i)^data\\[(\\d+)\\](.*)")
	match = re.FindStringSubmatch(k)
	if len(match) > 1 {
		key := fmt.Sprintf("document.%s%s", match[1], match[2])
		key = a.convertToGJsonPath(key)
		if a.data.Get(key).Exists() {
			return a.data.Get(key), nil
		}
		return gjson.Result{}, errors.New("not found")
	}
	return gjson.Result{}, errors.New("not found")
}

func (a *Referrable) convertToGJsonPath(k string) string {
	re := regexp.MustCompile("(.*)\\[(\\d+)\\](.*)")
	k = re.ReplaceAllString(k, "$1.$2$3")
	return k
}

// Find to find a value of given key
func (a *Referrable) Find(k string) ([]string, bool) {
	val, err := a.find(k)
	if err == nil {
		return []string{val.String()}, true
	}

	k = strings.ToLower(k)
	if val, ok := a.values[k]; ok {
		return val, true
	}
	return nil, false
}

// FindBoolean to find a value of given key and treat it as boolean.
// If the value is non-boolean, error will be raised.
// All header stuffs will be treated as non-boolean.
func (a *Referrable) FindBoolean(k string) (bool, bool) {
	val, err := a.find(k)
	if err == nil {
		if v, ok := val.Value().(bool); ok {
			return v, true
		}
	}

	return false, false
}

// FindNull to find a value of given key and treat it as null.
// If the value is not-null, error will be raised.
func (a *Referrable) FindNull(k string) (bool, bool) {
	val, err := a.find(k)
	if err == nil {
		if val.Type == gjson.Null {
			return true, true
		}
		return false, true
	}

	return false, false
}
