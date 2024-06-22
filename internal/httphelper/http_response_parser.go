package httphelper

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"

	"github.com/chonla/goline"
)

type HTTPResponseParser struct {
}

func (p *HTTPResponseParser) Parse(responseString string) (*HTTPResponse, error) {
	lines := goline.FromMultilineString(responseString)
	if len(lines) == 0 {
		return nil, errors.New("unexpected EOF")
	}

	headers := map[string]string{}
	bodyLines := []string{}
	body := ""
	status := ""
	wrappedBody := ""
	var bodyBuffer *bytes.Buffer

	firstLine := lines[0]
	if captured, ok := firstLine.CaptureAll(`^(.+)/(.+) (.\d+) (.+)$`); ok {
		statusCode, err := strconv.Atoi(captured[3])
		if err != nil {
			return nil, errors.New("unexpected status code")
		}
		statusText := captured[4]
		status = fmt.Sprintf("%d %s", statusCode, statusText)
		protocol := captured[1]
		protocolVersion := captured[2]
		headersCount := 0

		collectingHeader := true
		for _, line := range lines[1:] {
			if collectingHeader {
				if headerCaptured, ok := line.CaptureAll("^([^:]+):(.+)$"); ok {
					headers[goline.Line(headerCaptured[1]).Lower().Trim().Value()] = goline.Line(headerCaptured[2]).Trim().Value()
					headersCount += 1
				} else {
					if line.Value() == "" {
						collectingHeader = false
						continue
					}
				}
			} else {
				bodyLines = append(bodyLines, line.Value())
			}
		}

		r := &http.Response{
			Status:     status,
			StatusCode: statusCode,
			Header:     http.Header{},
		}
		valueMap := ValueMap{
			"Protocol":   protocol,
			"Version":    protocolVersion,
			"StatusCode": statusCode,
			"StatusText": statusText,
			"Status":     status,
		}

		if headersCount > 0 {
			valueMap["Headers"] = map[string]string{}
			for headerKey, headerValue := range headers {
				r.Header.Set(headerKey, headerValue)
				valueMap["Headers"].(map[string]string)[headerKey] = headerValue
			}
		}

		if len(bodyLines) > 0 {
			body = strings.Join(bodyLines, "\n")
			bodyBuffer = bytes.NewBuffer([]byte(body))
			r.Body = io.NopCloser(bodyBuffer)
			wrappedBody = fmt.Sprintf(`{"Body":%s}`, body)
		}

		return &HTTPResponse{
			resp:            r,
			respRaw:         responseString,
			protocol:        protocol,
			protocolVersion: protocolVersion,
			headers:         headers,
			body:            body,
			statusCode:      statusCode,
			status:          status,
			valueMap:        valueMap,
			wrappedBody:     wrappedBody,
		}, nil
	}
	return nil, errors.New("unexpected http response")
}

func (p *HTTPResponseParser) From(response *http.Response) (*HTTPResponse, error) {
	// There is some weird bug when using httputil.DumpResponse(response, true)
	// Response needs separated reading the header and body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	response.Body = io.NopCloser(bytes.NewBuffer(body))

	bodyString := string(body)

	headerBytes, err := httputil.DumpResponse(response, false)
	if err != nil {
		return nil, err
	}
	headerString := string(headerBytes)

	respRaw := fmt.Sprintf("%s\r\n\r\n%s", headerString, bodyString)

	return p.Parse(respRaw)
}
