package discord

import (
	"fmt"
	"strings"

	"github.com/adinb/weissbot/internal/pkg/cotd"
	"github.com/adinb/weissbot/internal/pkg/mtg"
	"github.com/adinb/weissbot/internal/pkg/rakugaki"
	"github.com/adinb/weissbot/internal/pkg/sakuga"
	"github.com/bwmarrin/discordgo"
)

const defaultWeissStatus = "with Schwarz | :weiss-help"

func ready(s *discordgo.Session, event *discordgo.Event) {
	s.UpdateStatus(0, defaultWeissStatus)
}

func (d *DiscordController) createMessageCreateHandler(env string) func(s *discordgo.Session, m *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		prefix := ""
		if env != "production" {
			prefix = ":test"
		}

		if m.Author.ID == s.State.User.ID {
			return
		}

		if strings.HasPrefix(m.Content, prefix+":weiss-help") {
			sendHelpMessage(s, m)
			return
		}

		if strings.HasPrefix(m.Content, prefix+":cotd vg") {
			cotd, err := d.vanguardCOTDService.GetCOTD()
			if err != nil {
				d.errorChannel <- err
				return
			}

			sendCOTD(cotd, s, m)
			return
		}

		if strings.HasPrefix(m.Content, prefix+":cotd bf") {
			cotd, err := d.buddyfightCOTDService.GetCOTD()
			if err != nil {
				d.errorChannel <- err
				return
			}

			sendCOTD(cotd, s, m)
			return
		}

		if strings.HasPrefix(m.Content, prefix+":cotd ws") {
			cotd, err := d.weissschwarzCOTDService.GetCOTD()
			if err != nil {
				d.errorChannel <- err
				return
			}

			sendCOTD(cotd, s, m)
			return
		}

		if strings.HasPrefix(m.Content, prefix+":dailysakuga") {
			sakuga, err := d.sakugaService.GetSakuga()
			if err != nil {
				d.errorChannel <- err
				return
			}

			sendDailySakuga(sakuga, s, m)
			return
		}

		if strings.HasPrefix(m.Content, prefix+":dailyrkgk") {
			rkgk, err := d.rakugakiService.GetTopRakugaki(100)
			if err != nil {
				d.errorChannel <- err
				return
			}

			sendDailyRkgk(rkgk, s, m)
			return
		}

		if strings.HasPrefix(m.Content, prefix+":mtg-search") {
			index := strings.Index(m.Content, " ")
			name := []byte(m.Content)[index+1:]
			cards, err := d.mtgService.SearchCardByName(string(name))
			err = sendMTGSearchResult(cards, s, m)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}
}

func sendDailyRkgk(rkgk rakugaki.Rakugaki, s *discordgo.Session, m *discordgo.MessageCreate) error {
	channel, err := s.State.Channel(m.ChannelID)
	if err != nil {
		return err
	}

	_, err = s.ChannelMessageSend(channel.ID, ":angry:")
	if err != nil {
		return err
	}

	_, err = s.ChannelMessageSend(channel.ID, rkgk.SourceURL)

	embedImage := discordgo.MessageEmbedImage{URL: rkgk.ImageURL}
	embed := discordgo.MessageEmbed{Image: &embedImage}
	s.ChannelMessageSendEmbed(channel.ID, &embed)

	return nil
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

	DailySakugaField := new(discordgo.MessageEmbedField)
	DailySakugaField.Name = "Daily Sakuga"
	DailySakugaField.Value = "Everyone loves sakuga. Enjoy random sakuga by typing `:dailysakuga`"

	fields = append(fields, cardOfTheDayField, dailyRkgkField, MTGSearchField, DailySakugaField)

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

func sendCOTD(cards []cotd.COTD, s *discordgo.Session, m *discordgo.MessageCreate) error {
	channel, err := s.State.Channel(m.ChannelID)
	if err != nil {
		return err
	}

	_, err = s.ChannelMessageSend(channel.ID, "Here are CoTD you asked :video_game:")
	if err != nil {
		return err
	}

	for _, card := range cards {
		embedImage := discordgo.MessageEmbedImage{URL: card.ImageURL}
		embed := discordgo.MessageEmbed{Image: &embedImage}

		s.ChannelMessageSendEmbed(channel.ID, &embed)
	}

	return nil
}

func sendDailySakuga(sakuga sakuga.Sakuga, s *discordgo.Session, m *discordgo.MessageCreate) error {
	channel, err := s.State.Channel(m.ChannelID)
	if err != nil {
		return err
	}

	_, err = s.ChannelMessageSend(channel.ID, ":open_mouth:")
	if err != nil {
		return err
	}
	_, err = s.ChannelMessageSend(channel.ID, sakuga.URL)
	if err != nil {
		return err
	}

	return nil
}

func sendMTGSearchResult(cards []*mtg.MagicCard, s *discordgo.Session, m *discordgo.MessageCreate) error {
	channel, err := s.State.Channel(m.ChannelID)
	if err != nil {
		return err
	}

	index := strings.Index(m.Content, " ")
	name := []byte(m.Content)[index+1:]

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
			imgA, err := mtg.RetrievePNG(card.Faces[0].ImageURIs.PNG)
			if err != nil {
				return err
			}

			imgB, err := mtg.RetrievePNG(card.Faces[1].ImageURIs.PNG)
			if err != nil {
				return err
			}

			tiledImage := mtg.TileImagesHorizontally(imgA, imgB)
			tiledImageReader := mtg.CreateImageReader(tiledImage)

			var imgFile discordgo.File
			imgFile.Name = card.ID + ".jpg"
			imgFile.Reader = tiledImageReader
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
