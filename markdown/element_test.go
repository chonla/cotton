package markdown

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewElementWithKnownElements(t *testing.T) {
	data := [][]string{
		[]string{"# Title 1"},
		[]string{"## Title 2"},
		[]string{"### Title 3"},
		[]string{"#### Title 4"},
		[]string{"##### Title 5"},
		[]string{"###### Title 6"},
		[]string{"* Bullet 1"},
		[]string{"* Bullet 1 with [Anchor Text](Anchor Link)"},
		[]string{"* [Anchor Text](Anchor Link)"},
		[]string{"```", "text in code block", "another text", "```"},
		[]string{"~~~", "text in code block", "another text", "~~~"},
		[]string{"```json", "text in code block", "another text", "```"},
		[]string{"~~~json", "text in code block", "another text", "~~~"},
	}
	expected := []ElementInterface{
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
		&RichTextElement{
			BaseElement: &BaseElement{
				Type: "Bullet",
			},
			Text:   "Bullet 1",
			Anchor: []AnchorElement{},
		},
		&RichTextElement{
			BaseElement: &BaseElement{
				Type: "Bullet",
			},
			Text: "Bullet 1 with [Anchor Text](Anchor Link)",
			Anchor: []AnchorElement{
				AnchorElement{
					Text: "Anchor Text",
					Link: "Anchor Link",
				},
			},
		},
		&RichTextElement{
			BaseElement: &BaseElement{
				Type: "Bullet",
			},
			Text: "[Anchor Text](Anchor Link)",
			Anchor: []AnchorElement{
				AnchorElement{
					Text: "Anchor Text",
					Link: "Anchor Link",
				},
			},
		},
		&SimpleElement{
			BaseElement: &BaseElement{
				Type: "Code",
			},
			Text: "text in code block\nanother text",
		},
		&SimpleElement{
			BaseElement: &BaseElement{
				Type: "Code",
			},
			Text: "text in code block\nanother text",
		},
		&SimpleElement{
			BaseElement: &BaseElement{
				Type: "Code",
			},
			Text: "text in code block\nanother text",
		},
		&SimpleElement{
			BaseElement: &BaseElement{
				Type: "Code",
			},
			Text: "text in code block\nanother text",
		},
	}

	for i := 0; i < len(data); i++ {
		r := NewElement(data[i])

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

func TestTryHeaderWithUnmatchedHeader(t *testing.T) {
	data := []string{"## Title"}

	r, ok := tryHeader(data, 1)

	assert.False(t, ok)
	assert.Nil(t, r)
}

func TestTryHeaderWithNotAHeaderElement(t *testing.T) {
	data := []string{"Not a header"}

	r, ok := tryHeader(data, 1)

	assert.False(t, ok)
	assert.Nil(t, r)
}
