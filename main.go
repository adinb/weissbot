package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	if os.Getenv("ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal(err)
		}
	}

	port := os.Getenv("PORT")
	discordServiceChannels := createDiscordServiceChannel()
	initializeHTTPServer(discordServiceChannels)

	go startDiscordBot(discordServiceChannels, &http.Client{Timeout: time.Second * 10})
	go startHTTPServer(port)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
