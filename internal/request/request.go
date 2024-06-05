package request

import (
	"bytes"
	"net/http"
	"net/http/httputil"
)

func Similar(req1, req2 *http.Request) bool {
	req1Bytes, err := httputil.DumpRequest(req1, true)
	if err != nil {
		return false
	}
	req2Bytes, err := httputil.DumpRequest(req2, true)
	if err != nil {
		return false
	}
	return bytes.Equal(req1Bytes, req2Bytes)
}
