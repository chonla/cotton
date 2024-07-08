package httphelper_test

import (
	"cotton/internal/httphelper"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetResponseHeaderValues(t *testing.T) {
	respString := `HTTP/1.1 200 OK
Content-Length: 15
Content-Type: application/json

{"key":"value"}`
	responseParser := &httphelper.HTTPResponseParser{}
	resp, _ := responseParser.Parse(respString)

	protocolValue, _ := resp.ValueOf("Protocol")
	assert.Equal(t, "HTTP", protocolValue.Value())

	protocolVersionValue, _ := resp.ValueOf("Version")
	assert.Equal(t, "1.1", protocolVersionValue.Value())

	statusCodeValue, _ := resp.ValueOf("StatusCode")
	assert.Equal(t, float64(200), statusCodeValue.Value())

	statusTextValue, _ := resp.ValueOf("StatusText")
	assert.Equal(t, "OK", statusTextValue.Value())

	statusValue, _ := resp.ValueOf("Status")
	assert.Equal(t, "200 OK", statusValue.Value())

	contentLengthValue, _ := resp.ValueOf("Headers.Content-Length")
	assert.Equal(t, "15", contentLengthValue.Value())

	contentTypeValue, _ := resp.ValueOf("Headers.Content-Type")
	assert.Equal(t, "application/json", contentTypeValue.Value())

	bodyValue, _ := resp.ValueOf("Body.key")
	assert.Equal(t, "value", bodyValue.Value())
}
