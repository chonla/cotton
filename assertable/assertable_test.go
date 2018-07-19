package assertable

import (
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

	assertable, e := NewAssertable(response)
	ref, _ := referrable.NewReferrable(response)
	expectedAssertable := &Assertable{
		Referrable: ref,
	}

	assert.Nil(t, e)
	assert.Equal(t, expectedAssertable, assertable)
}
