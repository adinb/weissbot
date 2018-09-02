package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func startHTTPServer(port string) {
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}
}

func initializeHTTPServer(discordServiceChannels discordServiceChannelsStruct) {
	discordStatusHandler := createDiscordStatusHandler(discordServiceChannels.statusChannel)
	http.HandleFunc("/", handleMainPage)
	http.HandleFunc("/discord_status", discordStatusHandler)
}

func createDiscordStatusHandler(channel chan<- string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}

		var message discordStatusStruct
		err = json.Unmarshal(body, &message)
		if err != nil {
			panic(err)
		}

		channel <- message.Status
		w.Write([]byte(message.Status))
	}
}

func handleMainPage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hi, I'm Weiss! What can I do for you?"))
}
