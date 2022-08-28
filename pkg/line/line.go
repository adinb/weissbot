package line

import (
	"github.com/adinb/weissbot/pkg/command"
	"github.com/adinb/weissbot/pkg/config"
	"github.com/adinb/weissbot/pkg/cotd"
	"github.com/adinb/weissbot/pkg/rakugaki"
	"github.com/adinb/weissbot/pkg/twitter"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"log"
	"net/http"
	"strconv"
	"time"
)

const httpServerReadHeaderTimeout = 3
const httpServerWriteTimeout = 10

type service struct {
	logger        *log.Logger
	environment   string
	listenAddress string
	linebotClient *linebot.Client
	executers     []command.Executer
	twitterClient twitter.Client
}

func New(config *config.Root) (*http.Server, error) {
	bot, err := linebot.New(config.Line.ChannelSecret, config.Line.ChannelAccessToken)
	if err != nil {
		return nil, err
	}

	controller := new(service)
	controller.environment = config.Weissbot.Environment
	controller.linebotClient = bot
	controller.listenAddress = config.Line.Addr + ":" + strconv.Itoa(int(config.Line.Port))

	// Executers
	controller.executers = append(controller.executers, cotd.NewExecuter())
	controller.executers = append(controller.executers, command.NewHelpExecuter())

	if config.Twitter.Enabled {
		httpClient := http.Client{Timeout: 10 * time.Second}
		twitterClient := twitter.New(&httpClient, config.Twitter.Token)

		controller.executers = append(controller.executers, twitter.NewExecuter(twitterClient))
		controller.executers = append(controller.executers, rakugaki.NewExecuter(twitterClient))
	}

	handler := handler{controller}
	mux := http.NewServeMux()
	mux.Handle("/line", handler)

	httpServer := http.Server{
		Addr:              controller.listenAddress,
		ReadHeaderTimeout: httpServerReadHeaderTimeout * time.Second,
		WriteTimeout:      httpServerWriteTimeout * time.Second,
		Handler:           mux,
	}

	return &httpServer, nil
}

func Start(server *http.Server, logger *log.Logger) error {
	logger.Printf("LINE webhook server listening at %s\n", server.Addr)
	return server.ListenAndServe()
}
