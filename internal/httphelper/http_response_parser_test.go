package httphelper_test

import (
	"cotton/internal/httphelper"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseHTTPResponse(t *testing.T) {
	respString := `HTTP/1.1 200 OK
Content-Length: 15
Content-Type: application/json

{"key":"value"}`
	responseParser := &httphelper.HTTPResponseParser{}
	resp, err := responseParser.Parse(respString)

	expectedSpecs := map[string]interface{}{
		"statusCode": 200,
		"status":     "200 OK",
		"headers": map[string]string{
			"content-type":   "application/json",
			"content-length": "15",
		},
		"body":            "{\"key\":\"value\"}",
		"protocol":        "HTTP",
		"protocolVersion": "1.1",
	}

	assert.NoError(t, err)
	assert.Equal(t, expectedSpecs, resp.Specs())
}
