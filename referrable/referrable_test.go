package referrable

import (
	"testing"

	"github.com/stretchr/objx"

	"github.com/chonla/cotton/response"
	"github.com/stretchr/testify/assert"
)

func TestIsJsonObjectShouldReturnTrueIfArrayContainsContentTypeOfApplicationJson(t *testing.T) {
	contentType := []string{
		"some-content-type",
		"application/json",
	}
	result := isJSONContent(contentType)
	assert.True(t, result)
}

func TestIsJsonObjectShouldReturnTrueIfArrayContainsContentTypeOfApplicationJsonWithCharset(t *testing.T) {
	contentType := []string{
		"some-content-type",
		"application/json; charset=utf-8",
	}
	result := isJSONContent(contentType)
	assert.True(t, result)
}

func TestIsJsonObjectShouldReturnTrueIfArrayNotContainsContentTypeOfApplicationJson(t *testing.T) {
	contentType := []string{
		"some-content-type",
		"application/pdf",
	}
	result := isJSONContent(contentType)
	assert.False(t, result)
}

func TestNewReferrableFromNonJsonResponse(t *testing.T) {
	jsonString := "{ \"data\": \"ok\" }"
	jsonObject, _ := objx.FromJSON(jsonString)

	response := &response.Response{
		Proto:      "http",
		Status:     "200 OK",
		StatusCode: 200,
		Header: map[string][]string{
			"content-type": []string{
				"application/json; charset=utf-8",
			},
		},
		Body: jsonString,
	}

	expected := &Referrable{
		values: map[string][]string{
			"status":     []string{"200 OK"},
			"statuscode": []string{"200"},
			"header.content-type": []string{
				"application/json; charset=utf-8",
			},
		},
		data: jsonObject,
	}

	result, e := NewReferrable(response)

	assert.Nil(t, e)
	assert.Equal(t, expected, result)
}
