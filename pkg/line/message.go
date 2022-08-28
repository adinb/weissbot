package line

import (
	"github.com/adinb/weissbot/pkg/command"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"log"
)

func getSourceIdFromEvent(event *linebot.Event) string {
	if event.Source.Type == linebot.EventSourceTypeGroup {
		return event.Source.GroupID
	}

	if event.Source.Type == linebot.EventSourceTypeUser {
		return event.Source.UserID
	}

	return event.Source.RoomID
}

func (s service) sendMessage(id string, result command.Result) {
	// Send text messages
	for _, message := range result.Messages {
		sendingMessage := linebot.NewTextMessage(message)
		if _, err := s.linebotClient.PushMessage(id, sendingMessage).Do(); err != nil {
			log.Println(err)
			return
		}
	}

	// Send image messages
	for _, imageMessage := range result.ImageMessages {
		var sendingMessage []linebot.SendingMessage

		// Add the image
		sendingImageMessage := linebot.NewImageMessage(
			imageMessage.ImageURL,
			imageMessage.ImageURL)
		sendingMessage = append(sendingMessage, sendingImageMessage)

		// Add a caption text if any
		if imageMessage.Message != "" {
			sendingMessage = append(sendingMessage, linebot.NewTextMessage(imageMessage.Message))
		}
		if _, err := s.linebotClient.PushMessage(id, sendingMessage...).Do(); err != nil {
			log.Println(err)
			return
		}
	}
}
func getTextMessage(event *linebot.Event) (message *linebot.TextMessage, ok bool) {
	if event.Type == linebot.EventTypeMessage {
		switch m := event.Message.(type) {
		case *linebot.TextMessage:
			message = m
			ok = true
			return
		}
	}

	message = nil
	ok = false

	return
}
