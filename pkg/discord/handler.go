package discord

import (
	"fmt"
	"log"
	"strings"

	"github.com/adinb/weissbot/pkg/cotd"
	"github.com/adinb/weissbot/pkg/mtg"
	"github.com/adinb/weissbot/pkg/rakugaki"
	"github.com/adinb/weissbot/pkg/sakuga"
	"github.com/bwmarrin/discordgo"
)

const defaultWeissStatus = "with Schwarz | :weiss-help"

func createReadyHandler(logger *log.Logger) func(*discordgo.Session, *discordgo.Event) {
	return func(s *discordgo.Session, _ *discordgo.Event) {
		if err := s.UpdateStatus(0, defaultWeissStatus); err != nil {
			logger.Println(err)
		}
	}
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
			if err := sendHelpMessage(s, m); err != nil {
				d.logger.Println(err)
				return
			}
			return
		}

		if strings.HasPrefix(m.Content, prefix+":bushiroadcotd vg") {
			cotd, err := d.vanguardCOTDService.GetCOTD()
			if err != nil {
				d.logger.Println(err)
				return
			}

			if err = sendCOTD(cotd, s, m); err != nil {
				d.logger.Println(err)
			}
			return
		}

		if strings.HasPrefix(m.Content, prefix+":bushiroadcotd bf") {
			cotd, err := d.buddyfightCOTDService.GetCOTD()
			if err != nil {
				d.logger.Println(err)
				return
			}

			if err = sendCOTD(cotd, s, m); err != nil {
				d.logger.Println(err)
			}
			return
		}

		if strings.HasPrefix(m.Content, prefix+":bushiroadcotd ws") {
			cotd, err := d.weissschwarzCOTDService.GetCOTD()
			if err != nil {
				d.logger.Println(err)
				return
			}

			if err = sendCOTD(cotd, s, m); err != nil {
				d.logger.Println(err)
			}
			return
		}

		if strings.HasPrefix(m.Content, prefix+":dailysakuga") {
			sakugas, err := d.sakugaService.GetSakuga()
			if err != nil {
				d.logger.Println(err)
				return
			}

			if err = sendDailySakuga(sakugas, s, m); err != nil {
				d.logger.Println(err)
			}
			return
		}

		if strings.HasPrefix(m.Content, prefix+":dailyrkgk") {
			rkgk, err := d.rakugakiService.GetTopRakugaki(100)
			if err != nil {
				d.logger.Println(err)
				return
			}

			if err = sendDailyRkgk(rkgk, s, m); err != nil {
				d.logger.Println(err)
			}
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
	if _, err = s.ChannelMessageSendEmbed(channel.ID, &embed); err != nil {
		return err
	}

	return nil
}

func sendHelpMessage(s *discordgo.Session, m *discordgo.MessageCreate) error {
	channel, err := s.State.Channel(m.ChannelID)
	if err != nil {
		return err
	}

	fields := make([]*discordgo.MessageEmbedField, 0)

	cardOfTheDayField := new(discordgo.MessageEmbedField)
	cardOfTheDayField.Name = "Card of The Day"
	cardOfTheDayField.Value = "Weiss can help you get **Vanguard** | **Buddyfight** | **Weiss Schwarz** CoTD by using `:bushiroadcotd vg` | `:bushiroadcotd bf` | `:bushiroadcotd ws` respectively"

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

	_, err = s.ChannelMessageSendEmbed(channel.ID, &embed)
	if err != nil {
		return err
	}

	return nil
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

		if _, err = s.ChannelMessageSendEmbed(channel.ID, &embed); err != nil {
			return err
		}
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

	if _, err = s.ChannelMessageSend(channel.ID, fmt.Sprintf("I found **%d** %s with **%s** in its name", len(cards), verb, name)); err != nil {
		return err
	}

	for _, card := range cards {
		var imageURL string
		var embed discordgo.MessageEmbed
		var m discordgo.MessageSend

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
			m.Files = append(m.Files, &imgFile)
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

		m.Embed = &embed
		m.Tts = false

		_, err = s.ChannelMessageSendComplex(channel.ID, &m)
		if err != nil {
			return err
		}
	}

	return nil
}
