package twitter

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

// Tweet contains simplified tweet data
type Tweet struct {
	IDStr          string
	Text           string
	UserName       string
	UserScreenName string
	RetweetCount   int
	FavoriteCount  int
	MediaUrls      []string
}

// UnmarshalJSON unmarshall tweet resource to a Tweet struct. It always read the original tweet data, even if the resource is a retweet tweet
func (t *Tweet) UnmarshalJSON(b []byte) error {
	var source map[string]interface{}
	var rawTweetData map[string]interface{}
	err := json.Unmarshal(b, &rawTweetData)
	if err != nil {
		return err
	}

	rts, rtsOk := rawTweetData["retweeted_status"]
	if rtsOk {
		source = rts.(map[string]interface{})
	} else {
		source = rawTweetData
	}

	t.IDStr = source["id_str"].(string)
	t.Text = source["text"].(string)
	t.RetweetCount = int(source["retweet_count"].(float64))
	t.FavoriteCount = int(source["favorite_count"].(float64))

	user := source["user"].(map[string](interface{}))
	t.UserName = user["name"].(string)
	t.UserScreenName = user["screen_name"].(string)

	extendedEntities := source["extended_entities"].(map[string](interface{}))
	media := extendedEntities["media"].([]interface{})

	for _, m := range media {
		t.MediaUrls = append(t.MediaUrls, m.(map[string](interface{}))["media_url_https"].(string))
	}
	return nil
}

func makeTweetFromMap(data map[string](interface{})) Tweet {
	var source map[string]interface{}
	var t Tweet

	_, extEntitiesOk := data["extended_entities"]
	if !extEntitiesOk {
		return Tweet{}
	}

	rts, rtsOk := data["retweeted_status"]
	if rtsOk {
		source = rts.(map[string]interface{})
	} else {
		source = data
	}

	t.IDStr = source["id_str"].(string)
	t.Text = source["text"].(string)
	t.RetweetCount = int(source["retweet_count"].(float64))
	t.FavoriteCount = int(source["favorite_count"].(float64))

	user := source["user"].(map[string](interface{}))
	t.UserName = user["name"].(string)
	t.UserScreenName = user["screen_name"].(string)

	extendedEntities := source["extended_entities"].(map[string](interface{}))
	media := extendedEntities["media"].([]interface{})

	for _, m := range media {
		t.MediaUrls = append(t.MediaUrls, m.(map[string](interface{}))["media_url_https"].(string))
	}

	return t
}

// SearchTweets returns raw array of tweets JSON payload
func SearchTweets(query string) ([]Tweet, error) {
	rawData, err := searchRawTweets(query)
	if err != nil {
		return nil, err
	}

	var tweetsData map[string](interface{})
	err = json.Unmarshal(rawData, &tweetsData)
	if err != nil {
		return nil, err
	}

	var tweets []Tweet
	statuses := tweetsData["statuses"].([]interface{})
	for _, status := range statuses {
		tweet := makeTweetFromMap(status.(map[string]interface{}))
		tweets = append(tweets, tweet)
	}

	return tweets, nil
}

func searchRawTweets(query string) ([]byte, error) {
	bearerToken := os.Getenv("TWITTER_TOKEN")
	searchQuery := "https://api.twitter.com/1.1/search/tweets.json?q=" + query + "&count=100"

	req, err := http.NewRequest("GET", searchQuery, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+bearerToken)

	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	buffer, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return buffer, nil
}

// FindTopTweets returns Tweets that has favorite count more than threshold
func FindTopTweets(threshold int, tweets []Tweet) []Tweet {
	var topTweets []Tweet
	for _, tweet := range tweets {
		if tweet.FavoriteCount > threshold {
			topTweets = append(topTweets, tweet)
		}
	}

	return topTweets
}
