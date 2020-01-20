package referrable

import (
	"fmt"
	"testing"

	"github.com/tidwall/gjson"

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

func TestNewReferrableFromJsonResponseDataShouldBeWrappedUnderDocument(t *testing.T) {
	jsonString := "{ \"data\": \"ok\" }"
	jsonObject := gjson.Parse(fmt.Sprintf("{ \"document\": %s }", jsonString))

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

	result := NewReferrable(response)

	assert.Equal(t, expected, result)
}

func TestNewReferrableFromJsonResponseAsListDataShouldBeWrappedUnderDocument(t *testing.T) {
	jsonString := "[{ \"data\": \"ok\" }]"
	jsonObject := gjson.Parse(fmt.Sprintf("{ \"document\": %s }", jsonString))

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

	result := NewReferrable(response)

	assert.Equal(t, expected, result)
}

func TestNewReferrableFromBrokenJsonResponseShouldContainEmptyData(t *testing.T) {
	jsonString := "{ \"data\": \"ok\""
	jsonObject := gjson.Parse("{}")

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

	result := NewReferrable(response)

	assert.Equal(t, expected, result)
}

func TestNewReferrableFromNonJsonResponseShouldContainEmptyData(t *testing.T) {
	jsonString := "{ \"data\": \"ok\"}"
	jsonObject := gjson.Parse("{}")

	response := &response.Response{
		Proto:      "http",
		Status:     "200 OK",
		StatusCode: 200,
		Header: map[string][]string{
			"content-type": []string{
				"text/plain",
			},
		},
		Body: jsonString,
	}

	expected := &Referrable{
		values: map[string][]string{
			"status":     []string{"200 OK"},
			"statuscode": []string{"200"},
			"header.content-type": []string{
				"text/plain",
			},
		},
		data: jsonObject,
	}

	result := NewReferrable(response)

	assert.Equal(t, expected, result)
}

func TestFindDataStuffInReferrableShouldReturnCorrespondingDataAndTrueWhenFound(t *testing.T) {
	jsonString := "{ \"data\": \"ok\", \"list\": [{ \"name\": \"john\" }, {\"name\": \"jane\" }] }"

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

	ref := NewReferrable(response)
	result, ok := ref.Find("data.list[1].name")

	assert.True(t, ok)
	assert.Equal(t, []string{"jane"}, result)
}

func TestFindDataWithNestedArrayStuffInReferrableShouldReturnCorrespondingDataAndTrueWhenFound(t *testing.T) {
	jsonString := "{ \"data\": \"ok\", \"list\": [{ \"children\": [{ \"name\": \"john\" }] }, { \"children\": [{\"name\": \"jane\" }] }] }"

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

	ref := NewReferrable(response)
	result, ok := ref.Find("data.list[1].children[0].name")

	assert.True(t, ok)
	assert.Equal(t, []string{"jane"}, result)
}

func TestFindDataStuffInDataListReferrableShouldReturnCorrespondingDataAndTrueWhenFound(t *testing.T) {
	jsonString := "[{ \"name\": \"john\" }, {\"name\": \"jane\" }]"

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

	ref := NewReferrable(response)
	result, ok := ref.Find("data[1].name")

	assert.True(t, ok)
	assert.Equal(t, []string{"jane"}, result)
}

func TestFindDataStuffInReferrableShouldReturnCorrespondingDataAndTrueWhenFoundWithCaseInsensitiveHeaderAccess(t *testing.T) {
	jsonString := "{ \"data\": \"ok\", \"list\": [{ \"Name\": \"john\" }, {\"Name\": \"jane\" }] }"

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

	ref := NewReferrable(response)
	result, ok := ref.Find("header.Content-TYPE")

	assert.True(t, ok)
	assert.Equal(t, []string{"application/json; charset=utf-8"}, result)
}

func TestFindDataStuffInReferrableShouldReturnCorrespondingDataAndTrueWhenFoundWithCaseSensitiveJsonAccess(t *testing.T) {
	jsonString := "{ \"data\": \"ok\", \"list\": [{ \"Name\": \"john\" }, {\"Name\": \"jane\" }] }"

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

	ref := NewReferrable(response)
	result, ok := ref.Find("data.list[1].Name")

	assert.True(t, ok)
	assert.Equal(t, []string{"jane"}, result)
}

func TestFindDataStuffInReferrableShouldReturnNilAndTrueWhenNotFound(t *testing.T) {
	jsonString := "{ \"data\": \"ok\" }"

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

	ref := NewReferrable(response)
	result, ok := ref.Find("data.result")

	assert.False(t, ok)
	assert.Nil(t, result)
}

func TestFindHeaderStuffInReferrableShouldReturnCorrespondingDataAndTrueWhenFound(t *testing.T) {
	jsonString := "{ \"data\": \"ok\" }"

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

	ref := NewReferrable(response)
	result, ok := ref.Find("header.content-type")

	assert.True(t, ok)
	assert.Equal(t, []string{"application/json; charset=utf-8"}, result)
}

func TestFindHeaderStuffInReferrableShouldReturnNilAndTrueWhenNotFound(t *testing.T) {
	jsonString := "{ \"data\": \"ok\" }"

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

	ref := NewReferrable(response)
	result, ok := ref.Find("header.content-length")

	assert.False(t, ok)
	assert.Nil(t, result)
}

func TestFindResponseStuffInReferrableShouldReturnCorrespondingDataAndTrueWhenFound(t *testing.T) {
	jsonString := "{ \"data\": \"ok\" }"

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

	ref := NewReferrable(response)
	result, ok := ref.Find("statuscode")

	assert.True(t, ok)
	assert.Equal(t, []string{"200"}, result)
}

func TestFindResponseStuffInReferrableShouldReturnNilAndTrueWhenNotFound(t *testing.T) {
	jsonString := "{ \"data\": \"ok\" }"

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

	ref := NewReferrable(response)
	result, ok := ref.Find("something")

	assert.False(t, ok)
	assert.Nil(t, result)
}
