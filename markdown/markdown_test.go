package markdown

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseWellFormedMarkdownString(t *testing.T) {
	markdown := `
		# Title
	`

	md := NewMD()

	result := md.ParseString(markdown)

	assert.Nil(t, result)
	assert.Equal(t, 1, md.Len())
}

func TestParseWellFormedMarkdownFile(t *testing.T) {
	readFileFn = func(filename string) ([]byte, error) {
		markdown := `
			# Title

			## POST /login
		`

		return []byte(markdown), nil
	}

	md := NewMD()

	result := md.Parse("file.md")

	assert.Nil(t, result)
	assert.Equal(t, 2, md.Len())
}

func TestParseWindowsMarkdownFile(t *testing.T) {
	readFileFn = func(filename string) ([]byte, error) {
		markdown := "# Title\r\n\r\n## POST /login"

		return []byte(markdown), nil
	}

	md := NewMD()

	result := md.Parse("file.md")

	assert.Nil(t, result)
	assert.Equal(t, 2, md.Len())
}

func TestParseErrorMarkdownFile(t *testing.T) {
	readFileFn = func(filename string) ([]byte, error) {
		markdown := `
			# Title

			## POST /login
		`

		return []byte(markdown), errors.New("oh")
	}

	md := NewMD()

	result := md.Parse("file.md")

	assert.NotNil(t, result)
}

func TestCursorShouldReturnTrueIfItemAvailable(t *testing.T) {
	markdown := `
		# Title

		## PUT /login
	`

	md := NewMD()

	md.ParseString(markdown)

	assert.True(t, md.Next())
}

func TestCursorShouldReturnFalseIfItemNotAvailable(t *testing.T) {
	markdown := `
		# Title

		## PUT /login
	`

	md := NewMD()

	md.ParseString(markdown)

	md.Next()
	md.Next()

	assert.False(t, md.Next())
}

func TestResetCursorShouldSetCursorToOrigin(t *testing.T) {
	markdown := `
		# Title

		## PUT /login
	`

	md := NewMD()

	md.ParseString(markdown)

	md.Next()
	md.Next()
	md.Next()
	md.Reset()

	assert.True(t, md.Next())
	assert.True(t, md.Next())
	assert.False(t, md.Next())
}

func TestGetValueAtCursor(t *testing.T) {
	markdown := `
		# Title

		## PUT /login
	`

	md := NewMD()

	md.ParseString(markdown)

	md.Next()
	md.Next()

	assert.Equal(t, md.Value(), &SimpleElement{
		BaseElement: &BaseElement{
			Type: "H2",
		},
		Text: "PUT /login",
	})
}

func TestGetValueAtTheEnd(t *testing.T) {
	markdown := `
		# Title

		## PUT /login
	`

	md := NewMD()

	md.ParseString(markdown)

	md.Next()
	md.Next()
	md.Next()

	assert.Nil(t, md.Value())
}
