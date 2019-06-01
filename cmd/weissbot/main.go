package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/adinb/weissbot/internal/pkg/controller/discord"
	"github.com/adinb/weissbot/internal/pkg/controller/http"
	"github.com/adinb/weissbot/internal/pkg/meta"
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
	discordToken := os.Getenv("TOKEN")
	twitterToken := os.Getenv("TWITTER_TOKEN")
	env := os.Getenv("ENV")

	errc := make(chan error, 1)

	discordMetac := make(chan meta.Meta, 1)
	discord := discord.NewDiscordController(env, twitterToken, discordToken, discordMetac, errc)
	discord.Start()

	httpMetac := make(chan meta.Meta, 1)
	srv := http.CreateAndStartHTTPServer(port, httpMetac, errc)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	for {
		select {
		case meta := <-httpMetac:
			discordMetac <- meta
		case sig := <-sc:
			fmt.Printf("Got %s signal. Shutting down...\n", sig)

			close(errc)
			close(discordMetac)
			close(httpMetac)

			srv.Shutdown(context.Background())
			err := discord.Stop()
			if err != nil {
				errc <- err
			} else {
				return
			}
		case err := <-errc:
			log.Fatalf("Weissbot has encountered an error: %s", err.Error())
		}
	}
}
