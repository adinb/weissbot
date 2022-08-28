package cotd

import (
	"errors"
	"github.com/adinb/weissbot/pkg/command"
	"github.com/antchfx/htmlquery"
	"net/http"
	"time"
)

const cotdCommand = "cotd"
const WebCOTDHttpClientTimeoutDurationSeconds = 10

var vanguardSubcommands = []string{"vg", "vanguard"}
var weissSchwarzSubcommands = []string{"ws", "weissschwarz"}

type webCOTD struct {
	WebpageURL   string
	ImageBaseURL string
	ImagePath    string
}

func (r webCOTD) get() (command.Result, error) {
	result := command.Result{}
	httpClient := http.Client{Timeout: WebCOTDHttpClientTimeoutDurationSeconds * time.Second}

	resp, err := httpClient.Get(r.WebpageURL)
	if err != nil {
		return result, err
	}

	if resp.StatusCode != http.StatusOK {
		return result, errors.New("failed to download HTML page")
	}
	defer resp.Body.Close()

	doc, err := htmlquery.Parse(resp.Body)
	if err != nil {
		return result, err
	}

	for _, data := range htmlquery.Find(doc, r.ImagePath) {
		imageUrl := r.ImageBaseURL + htmlquery.SelectAttr(data, "src")
		imageMessage := command.ImageMessage{
			Message:  "",
			ImageURL: imageUrl,
		}
		result.ImageMessages = append(result.ImageMessages, imageMessage)
	}

	return result, nil
}
