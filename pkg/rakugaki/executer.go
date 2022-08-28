package rakugaki

import (
	"errors"
	"math/rand"
	"time"

	"github.com/adinb/weissbot/pkg/command"
	"github.com/adinb/weissbot/pkg/twitter"
)

const rkgkCommand = "rkgk"

type Executer struct {
	command string
	client  twitter.Client
}

func (e Executer) GetCommand() string {
	return e.command
}

func NewExecuter(client twitter.Client) Executer {
	return Executer{
		command: rkgkCommand,
		client:  client,
	}
}
func (e Executer) Execute(args []string) (command.Result, error) {
	var result command.Result
	// TODO: Add threshold heuristic to filter tweets based on like / retweet counts
	result, err := e.getRakugaki()
	if err != nil {
		return result, errors.New("failed to retrieve the rkgk")
	}

	return result, nil
}

func filterRakugakiTweets(t twitter.TwitterResponse) twitter.TwitterResponse {
	filterResult := t

	return filterResult
}

func (e Executer) getRakugaki() (command.Result, error) {
	result := command.Result{}
	rand.Seed(time.Now().Unix())

	baseQuery := " -is:retweet -is:reply -is:quote has:images"
	rkgkQuery := "rkgk" + baseQuery
	hashtagRkgkQuery := "%23rkgk" + baseQuery
	jpRkgkQuery := "らくがき" + baseQuery
	expansions := []string{"author_id", "attachments.media_keys"}
	resultField := twitter.SearchResultField{
		Tweet: []string{"id", "text", "created_at"},
		Media: []string{"url"},
	}

	rkgk, _ := e.client.SearchRecentTweets(rkgkQuery, expansions, resultField, 10)
	hashtagRkgk, _ := e.client.SearchRecentTweets(hashtagRkgkQuery, expansions, resultField, 10)
	jpRkgk, _ := e.client.SearchRecentTweets(jpRkgkQuery, expansions, resultField, 10)

	var selectedRkgk twitter.TwitterResponse

	poolId := rand.Intn(3)
	switch poolId {
	case 0:
		selectedRkgk = rkgk
	case 1:
		selectedRkgk = hashtagRkgk
	case 2:
		selectedRkgk = jpRkgk
	}

	filteredRakugakiTweets := filterRakugakiTweets(selectedRkgk)
	if len(filteredRakugakiTweets.Statuses) > 0 {
		selectedRkgk = filteredRakugakiTweets
	}

	id := rand.Intn(len(selectedRkgk.Statuses))
	result.ImageMessages = append(result.ImageMessages, command.ImageMessage{
		Message:  selectedRkgk.Statuses[id].Text,
		ImageURL: selectedRkgk.Include.Media[id].URL,
	})

	if (len(result.ImageMessages)) == 0 {
		return command.Result{}, errors.New("no tweet found")
	}

	return result, nil
}
