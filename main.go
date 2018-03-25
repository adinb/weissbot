package main

import (
	"net/http"
	"os"
)

func handleMainPage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hi, I'm Weiss! What can I do for you?"))
}

func main() {
	StartDiscordBot()
	port := os.Getenv("PORT")

	http.HandleFunc("/", handleMainPage)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}

}
