package assertable

import (
	"errors"
	"testing"

	"github.com/chonla/cotton/referrable"
	"github.com/stretchr/testify/assert"

	"github.com/chonla/cotton/response"
)

func TestNewAssertableShouldReturnAssertableIfReferrableObjectCanBeCreated(t *testing.T) {
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

	assertable := NewAssertable(response)
	ref := referrable.NewReferrable(response)
	expectedAssertable := &Assertable{
		Referrable: ref,
	}

	assert.Equal(t, expectedAssertable, assertable)
}

func TestAssertWithNoAssertion(t *testing.T) {
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

	assertable := NewAssertable(response)

	result := assertable.Assert(map[string]string{})

	assert.Equal(t, errors.New("no assertion given"), result)
}

func TestAssertionAllPass(t *testing.T) {
	jsonString := "{ \"data\": \"ok\", \"list\": [0, 1, 2] }"

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

	assertable := NewAssertable(response)

	result := assertable.Assert(map[string]string{
		"header.content-type": "application/json; charset=utf-8",
		"data.list[1]":        "1",
	})

	assert.Nil(t, result)
}

func TestAssertionSomeFailShouldGiveResultFail(t *testing.T) {
	jsonString := "{ \"data\": \"ok\", \"list\": [0, 1, 2] }"

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

	assertable := NewAssertable(response)

	result := assertable.Assert(map[string]string{
		"header.content-type": "application/json; charset=utf-8",
		"data.list[1]":        "0",
	})

	assert.NotNil(t, result)
}

func TestAssertionNonExistingShouldGiveResultFail(t *testing.T) {
	jsonString := "{ \"data\": \"ok\", \"list\": [0, 1, 2] }"

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

	assertable := NewAssertable(response)

	result := assertable.Assert(map[string]string{
		"header.content-type": "application/json; charset=utf-8",
		"data.not-here":       "0",
	})

	assert.NotNil(t, result)
}
