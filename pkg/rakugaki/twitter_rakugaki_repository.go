package rakugaki

import (
	"github.com/adinb/weissbot/pkg/twitter"
)

type TwitterRakugakiRepository struct {
	Client twitter.Client
}

func (t *TwitterRakugakiRepository) List(query string) ([]Rakugaki, error) {
	var rakugakiList []Rakugaki
	tweets, err := t.Client.FindTweets(query)
	if err != nil {
		return rakugakiList, err
	}

	for _, tweet := range tweets {
		if len(tweet.MediaURLs) > 0 {
			r := Rakugaki{ImageURL: tweet.MediaURLs[0], SourceURL: tweet.URL, Rating: tweet.FavoriteCount}
			rakugakiList = append(rakugakiList, r)
		}
	}
	return rakugakiList, nil
}
