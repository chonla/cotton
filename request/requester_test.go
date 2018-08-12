package request

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateNewUnsupportedMethodRequest(t *testing.T) {
	_, e := NewRequester("HEAD", false, false)

	assert.NotNil(t, e)
}

func TestCreateNewGetRequest(t *testing.T) {
	r, e := NewRequester("GET", false, false)

	xreq := &Getter{
		Requester: &Requester{
			headers:     map[string]string{},
			client:      &http.Client{},
			insecure:    false,
			printDetail: false,
		},
	}

	assert.Nil(t, e)
	assert.Equal(t, xreq, r)
}

func TestCreateNewPostRequest(t *testing.T) {
	r, e := NewRequester("POST", false, false)

	xreq := &Poster{
		Requester: &Requester{
			headers:     map[string]string{},
			client:      &http.Client{},
			insecure:    false,
			printDetail: false,
		},
	}

	assert.Nil(t, e)
	assert.Equal(t, xreq, r)
}

func TestCreateNewPutRequest(t *testing.T) {
	r, e := NewRequester("PUT", false, false)

	xreq := &Putter{
		Requester: &Requester{
			headers:     map[string]string{},
			client:      &http.Client{},
			insecure:    false,
			printDetail: false,
		},
	}

	assert.Nil(t, e)
	assert.Equal(t, xreq, r)
}

func TestCreateNewPatchRequest(t *testing.T) {
	r, e := NewRequester("PATCH", false, false)

	xreq := &Patcher{
		Requester: &Requester{
			headers:     map[string]string{},
			client:      &http.Client{},
			insecure:    false,
			printDetail: false,
		},
	}

	assert.Nil(t, e)
	assert.Equal(t, xreq, r)
}

func TestCreateNewDeleteRequest(t *testing.T) {
	r, e := NewRequester("DELETE", false, false)

	xreq := &Deleter{
		Requester: &Requester{
			headers:     map[string]string{},
			client:      &http.Client{},
			insecure:    false,
			printDetail: false,
		},
	}

	assert.Nil(t, e)
	assert.Equal(t, xreq, r)
}

func TestEscapeURL(t *testing.T) {
	r := &Requester{}
	u := "http://www.google.com/q?param=data"
	expected := u

	result := r.EscapeURL(u)
	assert.Equal(t, expected, result)
}

func TestEscapeURLWithEscapableQueryString(t *testing.T) {
	r := &Requester{}
	u := "http://www.google.com/q?param=data test"
	expected := "http://www.google.com/q?param=data+test"

	result := r.EscapeURL(u)
	assert.Equal(t, expected, result)
}

func TestEscapeURLWithThaiLanguage(t *testing.T) {
	r := &Requester{}
	u := "http://www.google.com/q?param=ทดสอบ"
	expected := "http://www.google.com/q?param=%E0%B8%97%E0%B8%94%E0%B8%AA%E0%B8%AD%E0%B8%9A"

	result := r.EscapeURL(u)
	assert.Equal(t, expected, result)
}
