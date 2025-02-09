package commands

import (
	"fmt"
	"strings"
)

type SingleCommand struct {
	command string
	params  []Param
}

func (c *SingleCommand) String() string {
	if len(c.params) == 0 {
		return fmt.Sprintf("%s\n", c.command)
	}
	stringSlice := make([]string, len(c.params))
	for i, param := range c.params {
		stringSlice[i] = param.AsString()
	}
	return fmt.Sprintf("%s %s\n", c.command, strings.Join(stringSlice, " "))
}

func NewSingleCommand(command CommandType) *SingleCommand {
	return &SingleCommand{
		command: command.String(),
	}
}

func (c *SingleCommand) AddParams(params ...any) (*SingleCommand, error) {
	for _, param := range params {
		var p Param
		switch v := param.(type) {
		case string:
			p = StringParam(v)
		case int:
			p = IntParam(v)
		case bool:
			p = BoolParam(v)
		default:
			return nil, fmt.Errorf("unknown param type: %T", v)
		}
		c.params = append(c.params, p)
	}
	return c, nil
}
