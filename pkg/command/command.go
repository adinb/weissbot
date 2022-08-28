package command

import (
	"strings"
)

const (
	TwitterCommand = "tweet"
)

type Parameter struct {
	Command string
	Args    []string
}

type Executer interface {
	Execute([]string) (Result, error)
	GetCommand() string
}

type ImageMessage struct {
	Message  string
	ImageURL string
}

type Result struct {
	Messages      []string
	ImageMessages []ImageMessage
}

func ParseBotCommand(input string, commandPrefix string) (command Parameter, ok bool) {
	tokens := strings.Split(input, " ")

	if tokens[0] != commandPrefix {
		command = Parameter{}
		ok = false
		return
	}

	if len(tokens) == 1 {
		return Parameter{
			Command: "help",
		}, true
	}

	commandTokens := tokens[1:]
	command = Parameter{
		Command: commandTokens[0],
		Args:    commandTokens[1:],
	}

	ok = true
	return
}

func FindSubcommand(arg string, subcommands []string) bool {
	for _, subcmd := range subcommands {
		if arg == subcmd {
			return true
		}
	}

	return false
}
