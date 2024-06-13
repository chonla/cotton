package response

import (
	"cotton/internal/line"
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
)

type Response struct {
	response      *http.Response
	plainResponse []byte
	headerValues  ValueMap
	wrappedBody   string
	body          string
}

type ValueMap map[string]interface{}

func New(resp *http.Response) (*Response, error) {
	respBytes, err := httputil.DumpResponse(resp, true)
	if err != nil {
		return nil, err
	}

	values, err := parseHeaders(string(respBytes))
	if err != nil {
		return nil, err
	}

	body := extractBody(string(respBytes))

	return &Response{
		response:      resp,
		plainResponse: respBytes,
		headerValues:  values,
		wrappedBody:   fmt.Sprintf(`{"Body":%s}`, body),
		body:          body,
	}, nil
}

func (r *Response) String() string {
	return string(r.plainResponse)
}

func (r *Response) ValueOf(key string) (interface{}, error) {
	k := line.Line(key).Trim()
	if k.StartsWith("Body.") {
		value := gjson.Get(r.wrappedBody, k.Value())
		if value.Exists() {
			return value.Value(), nil
		}
		return nil, errors.New("value not found")
	}
	if value, ok := r.headerValues[key]; ok {
		return value, nil
	}
	return nil, errors.New("value not found")
}

func extractBody(resp string) string {
	idx := strings.Index(resp, "\r\n\r\n")
	if idx == -1 {
		return ""
	}
	return resp[idx+4:]
}

func parseHeaders(resp string) (ValueMap, error) {
	values := ValueMap{}

	data := line.FromMultilineString(resp)
	if len(data) < 1 {
		return nil, errors.New("invalid http response")
	}

	// Status Line
	if captures, ok := data[0].CaptureAll(`^([^\s]+) (\d+) (.+)$`); ok {
		statusCode, _ := strconv.Atoi(captures[2])
		values["Version"] = captures[1]
		values["StatusCode"] = statusCode
		values["StatusText"] = captures[3]
	}

	if len(data) > 1 {
		// Headers
		headers := map[string]string{}
		for _, headerLine := range data[1:] {
			if captures, ok := headerLine.CaptureAll(`^([^:]+):(.+)$`); ok {
				headers[line.Line(captures[1]).Trim().Lower().Value()] = line.Line(captures[2]).Trim().Value()
			} else {
				break
			}
		}
		values["Headers"] = headers
	}

	return values, nil
}
