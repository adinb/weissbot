package line

import (
	"github.com/adinb/weissbot/pkg/command"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"log"
	"net/http"
)

const weissbotLineCommandPrefix = "?w"

type handler struct {
	service *service
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	events, err := h.service.linebotClient.ParseRequest(r)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
			return
		}

		w.WriteHeader(500)
		return
	}

	h.service.processLineEvents(events)
	w.WriteHeader(200)
}

func (s service) processLineEvents(events []*linebot.Event) {
	prefix := weissbotLineCommandPrefix
	if s.environment != "production" {
		prefix = "test" + prefix
	}

	for _, event := range events {
		id := getSourceIdFromEvent(event)
		if m, ok := getTextMessage(event); ok {
			if botParameter, isValidCommand := command.ParseBotCommand(m.Text, prefix); isValidCommand {
				for _, executer := range s.executers {
					if botParameter.Command == executer.GetCommand() {
						result, err := executer.Execute(botParameter.Args)
						if err != nil {
							log.Println(err)
							return
						}

						s.sendMessage(id, result)
					}
				}
			}
		}
	}
}
