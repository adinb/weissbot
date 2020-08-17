package line

import (
	"fmt"
	"log"
	"net/http"

	"github.com/line/line-bot-sdk-go/linebot"
)

const commandPrefix = ":weiss"

// WebhookOptions contains webhook server configurable options
type WebhookOptions struct {
	Path    string
	Port    uint
	Address string
}

// WebhookServer handles incoming Line events
type WebhookServer struct {
	Line   *linebot.Client
	Server *http.Server
}

// NewWebhookServer creates a new Line Webhook Server
func NewWebhookServer(line *linebot.Client, options WebhookOptions) *WebhookServer {
	mux := createWebhookMux(options.Path, line)

	address := fmt.Sprintf("%s:%d", options.Address, options.Port)
	server := &http.Server{
		Addr:    address,
		Handler: mux,
	}

	webhookServer := WebhookServer{
		Line:   line,
		Server: server,
	}

	return &webhookServer
}

// NewLineClient initializes and returns a new Line client
func NewLineClient(lineChannelSecret string, lineChannelAccessToken string) *linebot.Client {
	line, err := linebot.New(lineChannelSecret, lineChannelAccessToken)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	return line
}

// Start starts the webhook server
func (webhook *WebhookServer) Start() {
	log.Printf("Listening at %s\n", webhook.Server.Addr)

	err := webhook.Server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func createWebhookMux(path string, line *linebot.Client) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		events, err := line.ParseRequest(r)
		if err != nil {
			log.Println(err)
		}

		for _, event := range events {
			if event.Type == linebot.EventTypeMessage {
				message := (event.Message).(*linebot.TextMessage)
				if message.Text[:1] == commandPrefix {
					sendingMessage := linebot.NewTextMessage("Thanks for sending me a message!")
					_, err := line.ReplyMessage(event.ReplyToken, sendingMessage).Do()

					if err != nil {
						log.Println(err)
					}
				}
			}
		}
	})

	return mux
}
