package rakugaki

import (
	"testing"
	"weissbot/twitter"
)

func mockTweetSearch(query string) ([]twitter.Tweet, error) {
	var tweets []twitter.Tweet

	if query == "%23rkgk" {
		tweets = append(tweets, twitter.Tweet{FavoriteCount: 10})
		tweets = append(tweets, twitter.Tweet{FavoriteCount: 100})
		tweets = append(tweets, twitter.Tweet{FavoriteCount: 40})
	}

	if query == "rkgk" {
		tweets = append(tweets, twitter.Tweet{FavoriteCount: 150})
		tweets = append(tweets, twitter.Tweet{FavoriteCount: 90})
	}

	if query == "%e3%82%89%e3%81%8f%e3%81%8c%e3%81%8d" {
		tweets = append(tweets, twitter.Tweet{FavoriteCount: 10})
		tweets = append(tweets, twitter.Tweet{FavoriteCount: 120})
		tweets = append(tweets, twitter.Tweet{FavoriteCount: 400})
	}

	return tweets, nil
}

func TestGetRakugaki(t *testing.T) {
	const favoriteThreshold = 100
	rkgkTweet, _ := GetRakugaki(mockTweetSearch)

	if rkgkTweet.FavoriteCount < favoriteThreshold {
		t.Error("Retrieved tweet has less favorite count than threshold:", rkgkTweet.FavoriteCount)
	}
}
