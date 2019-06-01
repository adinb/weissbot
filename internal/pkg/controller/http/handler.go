package http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/adinb/weissbot/internal/pkg/meta"
)

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
