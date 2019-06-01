package rakugaki

import (
	"fmt"
	"errors"
	"testing"

	"github.com/adinb/weissbot/internal/pkg/twitter"
)

type mockTwitterClient struct {
	isFailed bool
	tweets []twitter.Tweet
}

func (t *mockTwitterClient) FindTweets(query string) ([]twitter.Tweet, error) {
	if t.isFailed {
		return nil, errors.New("Failed to find tweets")
	}

	return t.tweets, nil
}

func TestTwitterRakugakiRepositoryListSuccess(t *testing.T) {
	var tweets []twitter.Tweet
	var mediaURLs []string
	mediaURLs = append(mediaURLs, "url1")
	mediaURLs = append(mediaURLs, "url2")
	tweets = append(tweets, twitter.Tweet{FavoriteCount: 100, URL: "url", MediaURLs: mediaURLs})

	m := mockTwitterClient{isFailed: false, tweets: tweets}
	r := TwitterRakugakiRepository{Client: &m}
	rakugakiList, err := r.List("query")
	if err != nil {
		t.Error("Failed to get rakugaki list", err)
	}
	if len(rakugakiList) != len(tweets) {
		t.Error("Record mismatch")
		return
	}

	for i, r := range(rakugakiList) {
		if r.ImageURL != tweets[i].MediaURLs[0] || r.SourceURL != tweets[i].URL || r.Rating != tweets[i].FavoriteCount {
			t.Error("Record mismatch")
			fmt.Println(r)
			fmt.Println(tweets[i])
		}
	}
}

func TestTwitterRakugakiRepositoryListFail(t *testing.T) {
	var tweets []twitter.Tweet
	tweets = append(tweets, twitter.Tweet{})
	m := mockTwitterClient{isFailed: true, tweets: tweets}
	r := TwitterRakugakiRepository{Client: &m}
	_, err := r.List("query")
	if err == nil {
		t.Error("Failed to throw error")
	}
}