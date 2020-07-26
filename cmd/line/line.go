package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/line/line-bot-sdk-go/linebot"
)

func main() {
	err := godotenv.Load()

	mux := http.NewServeMux()
	line, err := linebot.New(os.Getenv("LINE_CHANNEL_SECRET"), os.Getenv("LINE_CHANNEL_ACCESS_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	mux.HandleFunc("/line", func(w http.ResponseWriter, r *http.Request) {
		events, err := line.ParseRequest(r)
		if err != nil {
			log.Println(err)
		}

		for _, event := range events {
			if event.Type == linebot.EventTypeMessage {
				message := (event.Message).(*linebot.TextMessage)
				log.Println(message.Text)
			}
		}
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Println("Listening at port 8080")
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
