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
