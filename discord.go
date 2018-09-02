package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"net/http"

	"github.com/bwmarrin/discordgo"
)

var weissStatus string
var httpClient *http.Client

func sendImageFromURL(url string, s *discordgo.Session, c *discordgo.Channel) {
	var embed discordgo.MessageEmbed
	var embedImage discordgo.MessageEmbedImage

	embedImage.URL = url
	embed.Image = &embedImage
	fmt.Println(embed.URL)
	s.ChannelMessageSendEmbed(c.ID, &embed)
}

func speedCheck(s *discordgo.Session, m *discordgo.MessageCreate) {
	channel, err := s.State.Channel(m.ChannelID)
	if err != nil {
		return
	}

	_, err = s.ChannelMessageSend(channel.ID, "Speedcheck!!")
	if err != nil {
		return
	}

	sendImageFromURL(getRandomDorenoCardURL(), s, channel)
}

func sendCotd(game string, s *discordgo.Session, m *discordgo.MessageCreate) {

	var cotdUrls []string
	switch game {
	case vanguardName:
		cotdUrls = getVGCotd()
	case wsName:
		cotdUrls = getWSCotd()
	case bfName:
		cotdUrls = getBFCotd()
	}

	channel, err := s.State.Channel(m.ChannelID)
	if err != nil {
		return
	}

	_, err = s.ChannelMessageSend(channel.ID, "Here are CoTD you asked :video_game:")
	if err != nil {
		return
	}

	for _, cotdURL := range cotdUrls {
		sendImageFromURL(cotdURL, s, channel)
	}
}

func sendDailyRkgk(s *discordgo.Session, m *discordgo.MessageCreate) {
	channel, err := s.State.Channel(m.ChannelID)
	if err != nil {
		return
	}

	_, err = s.ChannelMessageSend(channel.ID, ":angry:")
	dailyRkgk := getDailyRkgk(httpClient)
	_, err = s.ChannelMessageSend(channel.ID, dailyRkgk.id)
	sendImageFromURL(dailyRkgk.mediaURL, s, channel)
}

func ready(s *discordgo.Session, event *discordgo.Event) {
	s.UpdateStatus(0, weissStatus)
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, ":cotd vg") {
		sendCotd(vanguardName, s, m)
		return
	}

	if strings.HasPrefix(m.Content, ":cotd ws") {
		sendCotd(wsName, s, m)
		return
	}

	if strings.HasPrefix(m.Content, ":cotd bf") {
		sendCotd(bfName, s, m)
		return
	}

	if strings.HasPrefix(m.Content, ":speedcheck") {
		speedCheck(s, m)
		return
	}

	if strings.HasPrefix(m.Content, ":dailyrkgk") {
		sendDailyRkgk(s, m)
		return
	}
}

func statusPoller(statusChannel <-chan string, s *discordgo.Session) {
	for status := range statusChannel {
		if weissStatus != status {
			weissStatus = status
			s.UpdateStatus(0, weissStatus)
		}
	}
}

// StartDiscordBot will start the discord bot
func StartDiscordBot(statusChannel <-chan string, client *http.Client) {
	token := os.Getenv("TOKEN")
	weissStatus = "with Schwarz"
	httpClient = client
	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		panic(err)
	}

	discord.AddHandler(ready)
	discord.AddHandler(messageCreate)

	go statusPoller(statusChannel, discord)

	err = discord.Open()
	if err != nil {
		panic(err)
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	discord.Close()
}
