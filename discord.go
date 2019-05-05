package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/adinb/weissbot/rakugaki"
	"github.com/adinb/weissbot/twitter"

	"github.com/bwmarrin/discordgo"
)

const defaultWeissStatus = "with Schwarz | :weiss-help"

var httpClient *http.Client

type discordServiceChannelsStruct struct {
	statusChannel chan string
}

type discordStatusStruct struct {
	Status string
}

func startDiscordBot(channels discordServiceChannelsStruct, client *http.Client) {
	token := os.Getenv("TOKEN")
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

func sendMTGSearchResult(s *discordgo.Session, m *discordgo.MessageCreate) error {
	channel, err := s.State.Channel(m.ChannelID)
	if err != nil {
		return err
	}

	index := strings.Index(m.Content, " ")
	name := []byte(m.Content)[index+1:]
	cards, err := SearchMagicCard(string(name))

	var verb string
	if len(cards) == 1 {
		verb = "card"
	} else {
		verb = "cards"
	}

	s.ChannelMessageSend(channel.ID, fmt.Sprintf("I found **%d** %s with **%s** in its name", len(cards), verb, name))

	for _, card := range cards {
		var imageURL string
		var embed discordgo.MessageEmbed
		var complex discordgo.MessageSend

		if len(card.Faces) == 2 {
			file, err := GenerateCombinedTwoFaceCardImage(card)
			if err != nil {
				return err
			}

			var imgFile discordgo.File
			imgFile.Name = card.ID + ".jpg"
			imgFile.Reader = file
			complex.Files = append(complex.Files, &imgFile)
			imageURL = "attachment://" + card.ID + ".jpg"
		} else {
			imageURL = card.Faces[0].ImageURIs.PNG
		}

		image := new(discordgo.MessageEmbedImage)
		image.URL = imageURL

		embed.Color = 0xea195f
		embed.Title = card.Name
		embed.Image = image
		embed.URL = card.ScryfallURI

		for _, face := range card.Faces {
			if face.Power != "" {
				embed.Description += fmt.Sprintf(
					"%s\n%s\n\n**%s**\n%s\n**%s/%s**\n**Artist:** %s\n*%s*\n",
					face.ManaCost,
					strings.Join(face.Colors, ", "),
					face.TypeLine,
					face.Text,
					face.Power,
					face.Toughness,
					face.Artist,
					face.FlavorText)
			} else {
				embed.Description += fmt.Sprintf(
					"%s\n%s\n\n**%s**\n%s\n**Artist:** %s\n*%s*\n",
					face.ManaCost,
					strings.Join(face.Colors, ", "),
					face.TypeLine,
					face.Text,
					face.Artist,
					face.FlavorText)
			}
		}

		embed.Description += fmt.Sprintf(
			"\n**Format:** %s\n**Rarity:** %s\n**Set:** %s\n**Release date:** %s\n",
			strings.Join(card.Legalities, ", "),
			card.Rarity,
			card.SetName,
			card.ReleaseDate)

		complex.Embed = &embed
		complex.Tts = false

		s.ChannelMessageSendComplex(channel.ID, &complex)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
	}

	return nil
}

func sendDailyRkgk(s *discordgo.Session, m *discordgo.MessageCreate) error {
	channel, err := s.State.Channel(m.ChannelID)
	if err != nil {
		return err
	}

	_, err = s.ChannelMessageSend(channel.ID, ":angry:")
	if err != nil {
		return err
	}

	dailyRkgk, err := rakugaki.GetRakugaki(twitter.SearchTweets)
	if err != nil {
		fmt.Println(err)
		return err
	}

	tweetURL := fmt.Sprintf("https://www.twitter.com/%s/status/%s", dailyRkgk.UserScreenName, dailyRkgk.IDStr)
	_, err = s.ChannelMessageSend(channel.ID, tweetURL)
	sendImageFromURL(dailyRkgk.MediaUrls[0], s, channel)

	return nil
}

func ready(s *discordgo.Session, event *discordgo.Event) {
	s.UpdateStatus(0, defaultWeissStatus)
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	prefix := ""
	if os.Getenv("ENV") != "production" {
		prefix = ":test"
	}

	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, prefix+":cotd vg") {
		sendCotd(vanguardName, s, m)
		return
	}

	if strings.HasPrefix(m.Content, prefix+":cotd ws") {
		sendCotd(wsName, s, m)
		return
	}

	if strings.HasPrefix(m.Content, prefix+":cotd bf") {
		sendCotd(bfName, s, m)
		return
	}

	if strings.HasPrefix(m.Content, prefix+":dailyrkgk") {
		sendDailyRkgk(s, m)
		return
	}

	if strings.HasPrefix(m.Content, prefix+":weiss-help") {
		sendHelpMessage(s, m)
		return
	}

	if strings.HasPrefix(m.Content, prefix+":mtg-search") {
		err := sendMTGSearchResult(s, m)
		if err != nil {
			fmt.Println(err.Error())
		}
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

	MTGSearchField := new(discordgo.MessageEmbedField)
	MTGSearchField.Name = "MTG Card Search"
	MTGSearchField.Value = "You're having MTG chat with your friends and need to do a quick card search? Try typing `:mtg-search <card name>`"

	fields = append(fields, cardOfTheDayField, dailyRkgkField, MTGSearchField)

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
		s.UpdateStatus(0, status)
	}
}
