package cotd

import (
	"errors"
	"github.com/adinb/weissbot/pkg/command"
)

type Executer struct {
	command          string
	vanguardCOTD     webCOTD
	weissSchwarzCOTD webCOTD
}

func NewExecuter() Executer {
	return Executer{
		command: cotdCommand,
		vanguardCOTD: webCOTD{
			WebpageURL:   "https://cf-vanguard.com/todays-card/",
			ImagePath:    "//p[contains(@class, 'text-center')]/img[contains(@class, 'alignnone')]",
			ImageBaseURL: "",
		},
		weissSchwarzCOTD: webCOTD{
			WebpageURL:   "https://ws-tcg.com/todays-card/",
			ImagePath:    "//div[contains(@class, 'entry-content')]/p/img",
			ImageBaseURL: "https://ws-tcg.com",
		},
	}
}
func (e Executer) Execute(args []string) (command.Result, error) {
	var result command.Result
	if len(args) == 0 {
		return result, errors.New("no arguments supplied")
	}

	if command.FindSubcommand(args[0], vanguardSubcommands) {
		result, err := e.vanguardCOTD.get()
		if err != nil {
			return result, err
		}
		return result, nil
	}

	if command.FindSubcommand(args[0], weissSchwarzSubcommands) {
		result, err := e.weissSchwarzCOTD.get()
		if err != nil {
			return result, err
		}
		return result, nil
	}

	return result, errors.New("unsupported subcommand")
}

func (e Executer) GetCommand() string {
	return e.command
}
