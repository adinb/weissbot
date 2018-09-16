package twitter

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestFindTopTweets(t *testing.T) {
	var tweets []Tweet
	expectedTweets := []Tweet{{FavoriteCount: 120}, {FavoriteCount: 400}}
	const favoriteThreshold = 100

	tweets = append(tweets, Tweet{FavoriteCount: 10})
	tweets = append(tweets, Tweet{FavoriteCount: 120})
	tweets = append(tweets, Tweet{FavoriteCount: 400})

	topTweets := FindTopTweets(favoriteThreshold, tweets)
	if !cmp.Equal(expectedTweets, topTweets) {
		t.Error("Top tweets list doesn't match expectation", cmp.Diff(expectedTweets, topTweets))
	}
}

func TestMakeTweetFromMap(t *testing.T) {
	const tweetJSONString = `
	{
		"id_str": "1041045876280504320",
		"text": "RT @H_cho3: 사요히나 낙서 rkgk\nさよひな https://t.co/1H8uRMU0l4",
		"extended_entities": {
			"media": [{
				"id": 1040913570265853952,
				"id_str": "1040913570265853952",
				"media_url_https": "https://pbs.twimg.com/media/DnIRYI2UcAADxts.jpg",
				"url": "https://t.co/1H8uRMU0l4",
				"display_url": "pic.twitter.com/1H8uRMU0l4",
				"expanded_url": "https://twitter.com/H_cho3/status/1040914206751522816/photo/1",
				"type": "photo"
			}]
		},
		"user": {
			"id_str": "825482847280455680",
			"name": "토퍼",
			"screen_name": "Thanks664"
		},
		"retweeted_status": {
			"created_at": "Sat Sep 15 10:44:20 +0000 2018",
			"id_str": "1040914206751522816",
			"text": "사요히나 낙서 rkgk\nさよひな https://t.co/1H8uRMU0l4",
			"extended_entities": {
				"media": [{
					"id": 1040913570265853952,
					"id_str": "1040913570265853952",
					"media_url_https": "https://pbs.twimg.com/media/DnIRYI2UcAADxts.jpg",
					"url": "https://t.co/1H8uRMU0l4",
					"display_url": "pic.twitter.com/1H8uRMU0l4",
					"expanded_url": "https://twitter.com/H_cho3/status/1040914206751522816/photo/1",
					"type": "photo"
				}]
			},
			"user": {
				"id": 2920809560,
				"id_str": "2920809560",
				"name": "初3",
				"screen_name": "H_cho3"
			},
			"retweet_count": 265,
			"favorite_count": 765
		},
		"retweet_count": 265,
		"favorite_count": 0
	}
	`

	var tweetDataMap map[string](interface{})
	expectedRetweetedTweet := Tweet{
		IDStr:          "1040914206751522816",
		Text:           "사요히나 낙서 rkgk\nさよひな https://t.co/1H8uRMU0l4",
		RetweetCount:   265,
		FavoriteCount:  765,
		UserName:       "初3",
		UserScreenName: "H_cho3",
		MediaUrls:      []string{"https://pbs.twimg.com/media/DnIRYI2UcAADxts.jpg"}}

	err := json.Unmarshal([]byte(tweetJSONString), &tweetDataMap)
	if err != nil {
		t.Error("Something wrong happened", err)
	}

	tweetData := makeTweetFromMap(tweetDataMap)

	if !cmp.Equal(tweetData, expectedRetweetedTweet) {
		t.Error("Tweet struct isn't properly created\n", cmp.Diff(tweetData, expectedRetweetedTweet))
	}
}
