package httphelper_test

import (
	"cotton/internal/httphelper"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseHTTPRequest(t *testing.T) {
	reqString := `POST http://localhost/path HTTP/1.1
Content-Length: 15
Content-Type: application/json

{"key":"value"}`
	requestParser := &httphelper.HTTPRequestParser{}
	req, err := requestParser.Parse(reqString)

	expectedSpecs := map[string]interface{}{
		"method": "POST",
		"path":   "http://localhost/path",
		"headers": map[string]string{
			"content-type":   "application/json",
			"content-length": "15",
		},
		"body":            "{\"key\":\"value\"}",
		"protocol":        "HTTP",
		"protocolVersion": "1.1",
	}

	assert.NoError(t, err)
	assert.Equal(t, expectedSpecs, req.Specs())
}
