package discord

import (
	"fmt"
	"net/http"
	"time"

	"github.com/adinb/weissbot/pkg/cotd"
	"github.com/adinb/weissbot/pkg/meta"
	"github.com/adinb/weissbot/pkg/mtg"
	"github.com/adinb/weissbot/pkg/rakugaki"
	"github.com/adinb/weissbot/pkg/sakuga"
	"github.com/adinb/weissbot/pkg/twitter"

	"github.com/bwmarrin/discordgo"
)

type DiscordController struct {
	vanguardCOTDService     cotd.ServiceContract
	buddyfightCOTDService   cotd.ServiceContract
	weissschwarzCOTDService cotd.ServiceContract
	sakugaService           sakuga.ServiceContract
	rakugakiService         rakugaki.ServiceContract
	mtgService              mtg.ServiceContract
	discordSession          *discordgo.Session
	metaChannel             <-chan meta.Meta
	errorChannel            chan<- error
}

func NewDiscordController(env string, twitterToken string, discordToken string, metac <-chan meta.Meta, errc chan<- error) DiscordController {
	var controller DiscordController
	discord, err := discordgo.New("Bot " + discordToken)
	if err != nil {
		errc <- err
	}

	controller.discordSession = discord
	controller.metaChannel = metac
	controller.errorChannel = errc

	client := http.Client{Timeout: time.Duration(10 * time.Second)}

	vanguardCOTDRepository := cotd.WebCOTDRepository{
		WebpageURL:   "https://cf-vanguard.com/todays-card/",
		ImagePath:    "//p[contains(@class, 'text-center')]/img[contains(@class, 'alignnone')]",
		ImageBaseURL: "",
		Client:       &client,
	}
	controller.vanguardCOTDService = &cotd.DefaultService{Repo: &vanguardCOTDRepository}

	buddyfightCOTDRepository := cotd.WebCOTDRepository{
		WebpageURL:   "https://fc-buddyfight.com/todays-card/",
		ImagePath:    "//div[contains(@class, 'lp_bg')]/div/div/img",
		ImageBaseURL: "",
		Client:       &client,
	}
	controller.buddyfightCOTDService = &cotd.DefaultService{Repo: &buddyfightCOTDRepository}

	weissschwarzCOTDRepository := cotd.WebCOTDRepository{
		WebpageURL:   "https://ws-tcg.com/todays-card/",
		ImagePath:    "//div[contains(@class, 'entry-content')]/p/img",
		ImageBaseURL: "https://ws-tcg.com",
		Client:       &client,
	}
	controller.weissschwarzCOTDService = &cotd.DefaultService{Repo: &weissschwarzCOTDRepository}

	sakugaHttpClient := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			req.URL.Scheme = "https"
			return nil
		},
	}
	sakugaRepository := sakuga.SakugabooruRepository{Client: &sakugaHttpClient, BaseURL: "https://www.sakugabooru.com/post/random"}
	controller.sakugaService = &sakuga.DefaultService{Repo: &sakugaRepository}

	twitterClient := twitter.New(&client, "https://api.twitter.com/1.1/search/tweets.json?q=", twitterToken)
	rakugakiTwitterRepository := rakugaki.TwitterRakugakiRepository{Client: &twitterClient}
	controller.rakugakiService = &rakugaki.DefaultService{Repo: &rakugakiTwitterRepository}

	mtgRepository := mtg.ScryfallRepository{Client: &client, BaseURL: "https://api.scryfall.com"}
	controller.mtgService = &mtg.DefaultService{Repo: &mtgRepository}

	discord.AddHandler(ready)
	discord.AddHandler(controller.createMessageCreateHandler(env))

	return controller
}

func (d *DiscordController) metaPoller() {
	for meta := range d.metaChannel {
		d.discordSession.UpdateStatus(0, meta.Status)
	}
}

func (d *DiscordController) Start() {

	go d.metaPoller()
	go func() {
		fmt.Println("Opening Discord websocket connection")
		err := d.discordSession.Open()
		if err != nil {
			d.errorChannel <- err
		}
	}()
}

func (d *DiscordController) Stop() error {
	err := d.discordSession.Close()
	if err != nil {
		return err
	}

	return nil
}
