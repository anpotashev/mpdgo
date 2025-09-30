package commands

import (
	"fmt"
	"strings"
)

type BatchCommand struct {
	commands []SingleCommand
}

func NewBatchCommands(cmds []SingleCommand, maxCommandsCount int) []BatchCommand {
	var result = make([]BatchCommand, 0)
	for len(cmds) > 0 {
		batchSize := min(len(cmds), maxCommandsCount)
		multiCommand := BatchCommand{commands: cmds[:batchSize]}
		result = append(result, multiCommand)
		cmds = cmds[batchSize:]
	}
	return result
}

func (m BatchCommand) String() string {
	stringSlice := make([]string, len(m.commands))
	for i, command := range m.commands {
		stringSlice[i] = command.String()
	}
	return fmt.Sprintf("command_list_begin\n%scommand_list_end\n", strings.Join(stringSlice, ""))
}
