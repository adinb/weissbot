package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var weissStatus string
var httpClient *http.Client

type discordServiceChannelsStruct struct {
	statusChannel chan string
}

type discordStatusStruct struct {
	Status string
}

func createDiscordServiceChannel() discordServiceChannelsStruct {
	var chanStruct discordServiceChannelsStruct
	chanStruct.statusChannel = make(chan string)
	return chanStruct
}

func sendImageFromURL(url string, s *discordgo.Session, c *discordgo.Channel) {
	var embed discordgo.MessageEmbed
	var embedImage discordgo.MessageEmbedImage

	embedImage.URL = url
	embed.Image = &embedImage
	fmt.Println(embed.URL)
	s.ChannelMessageSendEmbed(c.ID, &embed)
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

	if strings.HasPrefix(m.Content, ":dailyrkgk") {
		sendDailyRkgk(s, m)
		return
	}

	if strings.HasPrefix(m.Content, ":weiss-help") {
		sendHelpMessage(s, m)
		return
	}
}

func sendHelpMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	channel, err := s.State.Channel(m.ChannelID)
	if err != nil {
		return
	}

	fields := make([]*discordgo.MessageEmbedField, 0)

	cardOfTheDayField := new(discordgo.MessageEmbedField)
	cardOfTheDayField.Name = "Card of The Day"
	cardOfTheDayField.Value = "Weiss can help you get **Vanguard** | **Buddyfight** | **Weiss Schwarz** CoTD by using `:cotd vg` | `:cotd bf` | `:cotd ws` respectively"

	dailyRkgkField := new(discordgo.MessageEmbedField)
	dailyRkgkField.Name = "Daily Rakugaki"
	dailyRkgkField.Value = "Want to get triggered by `#rkgk`? Weiss can help you with that. Type `:dailyrgk` and prepare your :angry: react"

	fields = append(fields, cardOfTheDayField, dailyRkgkField)

	var footer discordgo.MessageEmbedFooter
	footer.Text = "Weiss will learn more tricks in the future, stay tuned!"

	var embed discordgo.MessageEmbed
	embed.Color = 0xea195f
	embed.Title = "Need help?"
	embed.Description = "Here's what Weiss can help you with:"
	embed.Fields = fields
	embed.Footer = &footer

	s.ChannelMessageSendEmbed(channel.ID, &embed)
	if err != nil {
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

func startDiscordBot(channels discordServiceChannelsStruct, client *http.Client) {
	token := os.Getenv("TOKEN")
	weissStatus = "with Schwarz | :weiss-help"
	httpClient = client
	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		panic(err)
	}

	discord.AddHandler(ready)
	discord.AddHandler(messageCreate)

	go statusPoller(channels.statusChannel, discord)

	err = discord.Open()
	if err != nil {
		panic(err)
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	discord.Close()
}
