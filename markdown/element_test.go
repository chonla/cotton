package markdown

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTryHeader(t *testing.T) {
	data := [][]string{
		[]string{"# Title 1"},
		[]string{"## Title 2"},
		[]string{"### Title 3"},
		[]string{"#### Title 4"},
		[]string{"##### Title 5"},
		[]string{"###### Title 6"},
	}
	expected := []*SimpleElement{
		&SimpleElement{
			BaseElement: &BaseElement{
				Type: "H1",
			},
			Text: "Title 1",
		},
		&SimpleElement{
			BaseElement: &BaseElement{
				Type: "H2",
			},
			Text: "Title 2",
		},
		&SimpleElement{
			BaseElement: &BaseElement{
				Type: "H3",
			},
			Text: "Title 3",
		},
		&SimpleElement{
			BaseElement: &BaseElement{
				Type: "H4",
			},
			Text: "Title 4",
		},
		&SimpleElement{
			BaseElement: &BaseElement{
				Type: "H5",
			},
			Text: "Title 5",
		},
		&SimpleElement{
			BaseElement: &BaseElement{
				Type: "H6",
			},
			Text: "Title 6",
		},
	}

	for i := 0; i < 6; i++ {
		r, ok := tryHeader(data[i], i+1)

		assert.True(t, ok)
		assert.Equal(t, expected[i], r)
	}
}

func TestTryHeaderWithMultiline(t *testing.T) {
	data := []string{"# Title", "ok"}

	r, ok := tryHeader(data, 1)

	assert.False(t, ok)
	assert.Nil(t, r)
}

func TestTryHeaderWithNoLine(t *testing.T) {
	data := []string{}

	r, ok := tryHeader(data, 1)

	assert.False(t, ok)
	assert.Nil(t, r)
}

func TestTryHeaderWithOutOfRange(t *testing.T) {
	data := []string{"####### Not a title"}

	r, ok := tryHeader(data, 7)

	assert.False(t, ok)
	assert.Nil(t, r)
}
