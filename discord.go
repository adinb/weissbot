package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func sendVanguardCotd(s *discordgo.Session, m *discordgo.MessageCreate) {
	cotdUrls := GetCotd()

	channel, err := s.State.Channel(m.ChannelID)
	if err != nil {
		return
	}

	_, err = s.ChannelMessageSend(channel.ID, "Here are CoTD you asked :video_game:")

	if err != nil {
		return
	}

	for _, cotdURL := range cotdUrls {
		var embed discordgo.MessageEmbed
		var embedImage discordgo.MessageEmbedImage

		embedImage.URL = cotdURL
		embed.Image = &embedImage
		fmt.Println(embed.URL)
		s.ChannelMessageSendEmbed(channel.ID, &embed)
	}
}

func ready(s *discordgo.Session, event *discordgo.Event) {
	s.UpdateStatus(0, "with Schwarz")
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, ":cotd vg") {
		sendVanguardCotd(s, m)
	}
}

// StartDiscordBot will start the discord bot
func StartDiscordBot() {
	token := os.Getenv("TOKEN")
	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		panic(err)
		return
	}

	discord.AddHandler(ready)
	discord.AddHandler(messageCreate)

	err = discord.Open()
	if err != nil {
		panic(err)
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	discord.Close()
}
