package commands

import (
	"fmt"
	"strings"
)

type Param interface {
	fmt.Stringer
}

type StringParam string

func (p StringParam) String() string {
	escapeValue := strings.ReplaceAll(string(p), "\\", "\\\\")
	escapeValue = strings.ReplaceAll(escapeValue, "\"", "\\\"")
	return fmt.Sprintf("\"%s\"", escapeValue)
}

type IntParam int

func (p IntParam) String() string {
	return fmt.Sprintf("%d", p)
}

type BoolParam bool

func (p BoolParam) String() string {
	if p {
		return "\"1\""
	}
	return "\"0\""
}
