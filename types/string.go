package types

import "strings"

type String string

func StringFromBytes(b []byte) String {
	return String(b)
}

func (o String) String() string {
	return string(o)
}

func (o String) LineBreak() []string {
	linebreak := "\n"
	if strings.Contains(o.String(), "\r\n") {
		linebreak = "\r\n"
	}

	lines := strings.Split(o.String(), linebreak)
	output := []string{}
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			output = append(output, line)
		}
	}
	return output
}
