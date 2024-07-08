package httphelper

import (
	"cotton/internal/line"
	"cotton/internal/value"
	"errors"
	"net/http"
	"reflect"

	"github.com/tidwall/gjson"
)

type HTTPResponse struct {
	resp            *http.Response
	respRaw         string
	headers         map[string]string
	body            string
	protocol        string
	protocolVersion string
	statusCode      float64
	status          string
	valueMap        ValueMap
	wrappedBody     string
}

func (r *HTTPResponse) Specs() map[string]interface{} {
	return map[string]interface{}{
		"statusCode":      r.statusCode,
		"status":          r.status,
		"headers":         r.headers,
		"body":            r.body,
		"protocol":        r.protocol,
		"protocolVersion": r.protocolVersion,
	}
}

func (r *HTTPResponse) SimilarTo(r2 *HTTPResponse) bool {
	if r2 == nil {
		return false
	}

	return r.status == r2.status &&
		r.body == r2.body &&
		r.protocol == r2.protocol &&
		r.protocolVersion == r2.protocolVersion &&
		reflect.DeepEqual(r.headers, r2.headers)
}

func (r *HTTPResponse) String() string {
	return r.respRaw
}

func (r *HTTPResponse) ValueOf(key string) (*value.Value, error) {
	if r.respRaw == "" {
		return value.NewUndefined(), errors.New("response is empty")
	}
	k := line.Line(key).Trim()
	if k.LookLike(`^Body\W`) {
		val := gjson.Get(r.wrappedBody, k.Value())
		if val.Exists() {
			return value.New(val.Value()), nil
		}
		return value.NewUndefined(), nil
	}
	if k.StartsWith("Headers.") {
		if val, ok := r.valueMap["Headers"].(map[string]string)[k.Lower().Value()[8:]]; ok {
			return value.New(val), nil
		}
		return value.NewUndefined(), nil
	}
	if k.Value() == "Protocol" ||
		k.Value() == "Version" ||
		k.Value() == "StatusCode" ||
		k.Value() == "StatusText" ||
		k.Value() == "Status" {
		return value.New(r.valueMap[k.Value()]), nil
	}
	return value.NewUndefined(), nil
}
