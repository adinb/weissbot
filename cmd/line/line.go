package main

import (
	"log"
	"os"
	"strconv"

	"github.com/adinb/weissbot/pkg/line"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	lineChannelSecret := os.Getenv("LINE_CHANNEL_SECRET")
	lineChannelAccessToken := os.Getenv("LINE_CHANNEL_ACCESS_TOKEN")

	port, err := strconv.Atoi(os.Getenv("LINE_WEBHOOK_PORT"))
	if err != nil {
		log.Fatal(err)
	}

	lineClient := line.NewLineClient(lineChannelSecret, lineChannelAccessToken)
	lineWebhook := line.NewWebhookServer(
		lineClient, 
		line.WebhookOptions{Path: "/line", Port: uint(port)},
	)

	lineWebhook.Start()
}
