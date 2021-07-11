package discord

import (
	"log"
	"net/http"
	"time"

	"github.com/adinb/weissbot/pkg/cotd"
	"github.com/adinb/weissbot/pkg/mtg"
	"github.com/adinb/weissbot/pkg/rakugaki"
	"github.com/adinb/weissbot/pkg/sakuga"
	"github.com/adinb/weissbot/pkg/twitter"

	"github.com/bwmarrin/discordgo"
)

type DiscordController struct {
	logger                  *log.Logger
	vanguardCOTDService     cotd.ServiceContract
	buddyfightCOTDService   cotd.ServiceContract
	weissschwarzCOTDService cotd.ServiceContract
	sakugaService           sakuga.ServiceContract
	rakugakiService         rakugaki.ServiceContract
	mtgService              mtg.ServiceContract
	discordSession          *discordgo.Session
}

func NewDiscordService(env string, twitterToken string, discordToken string, logger *log.Logger) (*DiscordController, error) {
	var controller = new(DiscordController)
	discord, err := discordgo.New("Bot " + discordToken)
	if err != nil {
		return nil, err
	}

	controller.logger = logger
	controller.discordSession = discord

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

	discord.AddHandler(createReadyHandler(controller.logger))
	discord.AddHandler(controller.createMessageCreateHandler(env))

	return controller, nil
}

func (d *DiscordController) Start() error {
	return d.discordSession.Open()
}

func (d *DiscordController) Stop() error {
	return d.discordSession.Close()
}
