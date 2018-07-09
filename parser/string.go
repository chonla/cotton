package parser

import (
	"bytes"
	"strings"
)

func (p *Parser) parseTestSuiteName(name string) string {
	tokens := p.tokenizeTestSuiteName(name)
	name = strings.Title(strings.ToLower(strings.Join(tokens, " ")))
	return name
}

func (p *Parser) tokenizeTestSuiteName(name string) []string {
	tokens := []string{}
	var buffer bytes.Buffer
	found := "others"
	lastFound := ""
	for _, c := range name {
		found = p.charType(c)

		if p.isFirstToken(lastFound, found) {
			b := buffer.String()
			if b != "" {
				tokens = append(tokens, b)
				buffer.Reset()
			}
		}
		if found != "others" {
			buffer.WriteRune(c)
		}

		lastFound = found
	}

	tokens = append(tokens, buffer.String())
	buffer.Reset()

	return tokens
}

func (p *Parser) isFirstToken(lastFound, currentFound string) bool {
	if lastFound == "upper" && currentFound == "lower" {
		return false
	}
	if lastFound != currentFound {
		return true
	}
	return false
}

func (p *Parser) charType(c rune) string {
	t := ""
	switch {
	case 'a' <= c && c <= 'z':
		t = "lower"
	case 'A' <= c && c <= 'Z':
		t = "upper"
	case '0' <= c && c <= '9':
		t = "num"
	default:
		t = "others"
	}
	return t
}
