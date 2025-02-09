package commands

import (
	"fmt"
	"strings"
)

type Param interface {
	AsString() string
}

type StringParam string

func (p StringParam) AsString() string {
	escapeValue := strings.ReplaceAll(string(p), "\\", "\\\\")
	escapeValue = strings.ReplaceAll(escapeValue, "\"", "\\\"")
	return fmt.Sprintf("\"%s\"", escapeValue)
}

type IntParam int

func (p IntParam) AsString() string {
	return fmt.Sprintf("%d", p)
}

type BoolParam bool

func (p BoolParam) AsString() string {
	if p {
		return "\"1\""
	}
	return "\"0\""
}
