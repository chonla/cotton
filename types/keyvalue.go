package types

import (
	"errors"
	"strings"
)

type KeyValue struct {
	Key   string
	Value string
}

func ParseKeyValue(kv string) (*KeyValue, error) {
	s := strings.SplitN(kv, "=", 2)
	if len(s) == 2 {
		return &KeyValue{s[0], s[1]}, nil
	}
	return nil, errors.New("unable to parse key-value variable")
}
