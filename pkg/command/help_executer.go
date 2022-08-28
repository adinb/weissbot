package command

const helpCommand = "help"

type HelpExecuter struct {
	message string
}

func NewHelpExecuter() Executer {
	return HelpExecuter{
		message: "Hi, I'm here to help!",
	}
}

func (e HelpExecuter) GetCommand() string {
	return helpCommand
}

func (e HelpExecuter) Execute(args []string) (Result, error) {
	return Result{Messages: []string{e.message}}, nil
}
