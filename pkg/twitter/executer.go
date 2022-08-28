package twitter

import (
	"errors"
	"github.com/adinb/weissbot/pkg/command"
)

const tweetCommand = "tweet"

type Executer struct {
	command string
	client  Client
}

func NewExecuter(client Client) Executer {
	return Executer{
		command: tweetCommand,
		client:  client,
	}
}
func (e Executer) Execute(args []string) (command.Result, error) {
	var result command.Result
	if len(args) == 0 {
		return result, errors.New("no arguments supplied")
	}

	result, err := e.client.GetLastTweetsByAuthor(args[0])
	if err != nil {
		return result, errors.New("failed to retrieve the tweet")
	}

	return result, nil
}

func (e Executer) GetCommand() string {
	return e.command
}
