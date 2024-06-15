package response

import (
	"bytes"
	"cotton/internal/line"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"strconv"

	"github.com/tidwall/gjson"
)

type Response struct {
	response     *http.Response
	headerValues ValueMap
	wrappedBody  string
	body         string
	fullResponse string
}

type ValueMap map[string]interface{}

func New(resp *http.Response) (*Response, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body = io.NopCloser(bytes.NewBuffer(body))

	bodyString := string(body)

	headerBytes, err := httputil.DumpResponse(resp, false)
	if err != nil {
		return nil, err
	}

	values, err := parseHeaders(string(headerBytes))
	if err != nil {
		return nil, err
	}

	return &Response{
		response:     resp,
		headerValues: values,
		wrappedBody:  fmt.Sprintf(`{"Body":%s}`, bodyString),
		body:         bodyString,
		fullResponse: fmt.Sprintf("%s\r\n\r\n%s", string(headerBytes), bodyString),
	}, nil
}

func (r *Response) String() string {
	return string(r.fullResponse)
}

func (r *Response) ValueOf(key string) (interface{}, error) {
	k := line.Line(key).Trim()
	if k.StartsWith("Body.") {
		value := gjson.Get(r.wrappedBody, k.Value())
		if value.Exists() {
			return value.Value(), nil
		}
		return nil, errors.New("value not found in response body")
	}
	if value, ok := r.headerValues[key]; ok {
		return value, nil
	}
	return nil, errors.New("value not found")
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
