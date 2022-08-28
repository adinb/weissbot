package discord

import (
	"github.com/adinb/weissbot/pkg/command"
	"github.com/adinb/weissbot/pkg/cotd"
	"log"
	"net/http"
	"time"

	"github.com/adinb/weissbot/pkg/config"
	"github.com/adinb/weissbot/pkg/mtg"
	"github.com/adinb/weissbot/pkg/rakugaki"
	"github.com/adinb/weissbot/pkg/sakuga"
	"github.com/adinb/weissbot/pkg/twitter"

	"github.com/bwmarrin/discordgo"
)

type DiscordController struct {
	logger          *log.Logger
	sakugaService   sakuga.ServiceContract
	rakugakiService rakugaki.ServiceContract
	mtgService      mtg.ServiceContract
	discordSession  *discordgo.Session
	cotdExecuter    command.Executer
}

func New(config *config.Root, logger *log.Logger) (*DiscordController, error) {
	var controller = new(DiscordController)
	discord, err := discordgo.New("Bot " + config.Discord.Token)
	if err != nil {
		return nil, err
	}

	controller.logger = logger
	controller.discordSession = discord

	client := http.Client{Timeout: 10 * time.Second}

	controller.cotdExecuter = cotd.NewExecuter()

	sakugaHttpClient := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			req.URL.Scheme = "https"
			return nil
		},
	}
	sakugaRepository := sakuga.SakugabooruRepository{Client: &sakugaHttpClient, BaseURL: "https://www.sakugabooru.com/post/random"}
	controller.sakugaService = &sakuga.DefaultService{Repo: &sakugaRepository}

	if config.Twitter.Enabled {
		twitterClient := twitter.New(&client, config.Twitter.Token)
		rakugakiTwitterRepository := rakugaki.TwitterRakugakiRepository{Client: twitterClient}
		controller.rakugakiService = &rakugaki.DefaultService{Repo: &rakugakiTwitterRepository}
	}

	mtgRepository := mtg.ScryfallRepository{Client: &client, BaseURL: "https://api.scryfall.com"}
	controller.mtgService = &mtg.DefaultService{Repo: &mtgRepository}

	discord.AddHandler(createReadyHandler(controller.logger))
	discord.AddHandler(controller.createMessageCreateHandler(config))

	return controller, nil
}

func (d *DiscordController) Start() error {
	d.logger.Println("Discord bot started")
	return d.discordSession.Open()
}

func (d *DiscordController) Stop() error {
	return d.discordSession.Close()
}
