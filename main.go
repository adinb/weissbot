package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
)

var discordStatus chan string

type discordStatusStruct struct {
	Status string
}

func handleMainPage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hi, I'm Weiss! What can I do for you?"))
}

func updateDiscordStatus(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	var message discordStatusStruct
	err = json.Unmarshal(body, &message)
	if err != nil {
		panic(err)
	}

	discordStatus <- message.Status
	w.Write([]byte(""))
}

func main() {
	discordStatus = make(chan string)

	go StartDiscordBot(discordStatus)

	port := os.Getenv("PORT")
	http.HandleFunc("/", handleMainPage)
	http.HandleFunc("/updateDiscordStatus", updateDiscordStatus)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}

}
