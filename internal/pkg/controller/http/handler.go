package http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/adinb/weissbot/internal/pkg/meta"
)

type LineEventSource struct {
	UserId  string `json:userId`
	GroupId string `json:groupId`
	Type    string `json:type`
}

type LineEventMessage struct {
	Type string `json:type`
	Id   string `json:id`
	Text string `json:text`
}

type LineEvent struct {
	Type       string           `type`
	ReplyToken string           `json:replyToken`
	Source     LineEventSource  `json:source`
	Timestamp  int64            `json:timestamp`
	Mode       string           `json:mode`
	Message    LineEventMessage `json:message`
}

type LineWebhookEvents struct {
	Events       []LineEvent `json:events`
	Destaination string      `json:destination`
}

func createWeissMetaHandler(channel chan<- meta.Meta) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}

		var message meta.Meta
		err = json.Unmarshal(body, &message)
		if err != nil {
			panic(err)
		}

		channel <- message
		w.Write([]byte(message.Status))
	}
}

func handleMainPage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hi, I'm Weiss! What can I do for you?"))
}

func handleLineEvent(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	var events LineWebhookEvents
	err = json.Unmarshal(body, &events)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(200)
}
