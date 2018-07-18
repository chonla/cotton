package referrable

import (
	"testing"

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
