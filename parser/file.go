package parser

import "strings"

func (p *Parser) readTestSuiteFile(file string) ([]string, error) {
	b, err := readFileFn(file)
	if err != nil {
		return []string{}, err
	}
	contents := string(b)

	linebreak := "\n"
	if strings.Contains(contents, "\r\n") {
		linebreak = "\r\n"
	}

	lines := strings.Split(contents, linebreak)
	output := []string{}
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			output = append(output, line)
		}
	}

	return output, nil
}
