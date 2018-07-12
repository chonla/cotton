package referrable

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/chonla/yas/response"
	"github.com/fatih/color"
	"github.com/stretchr/objx"
)

// Referrable is referrable items
type Referrable struct {
	values map[string][]string
	data   objx.Map
}

// NewReferrable creates an referrable object
func NewReferrable(resp *response.Response) (*Referrable, error) {
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

	var jsonObj objx.Map
	var e error
	if isJSONContent(values["header.content-type"]) {
		jsonObj, e = objx.FromJSON(resp.Body)
		if e != nil {
			fmt.Printf("%s: %s\n", red("Error"), e)
			jsonObj, _ = objx.FromJSON("{}")
		}
	} else {
		jsonObj, _ = objx.FromJSON("{}")
	}

	return &Referrable{
		values: values,
		data:   jsonObj,
	}, nil
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

// Find to find a value of given key
func (a *Referrable) Find(k string) ([]string, bool) {
	re := regexp.MustCompile("(?i)^data\\.(.+)")
	match := re.FindStringSubmatch(strings.ToLower(k))
	if len(match) > 1 {
		if a.data.Has(match[1]) {
			return []string{a.data.Get(match[1]).String()}, true
		}
		return []string{}, false
	} else {
		val, ok := a.values[k]
		return val, ok
	}
}
