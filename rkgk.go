package main

import (
  "log"
  "net/http"
  "math/rand"
  "encoding/json"
)

type TweetStruct struct {
  id string
  mediaURL string
}

// GetDailyRkgk will return URL to a rkgk tweet
func GetDailyRkgk(client *http.Client) TweetStruct {
  bodyJSON, err := SearchTweets(client, "%23rkgk")

  var body map[string]interface{}
  if err != nil {
    log.Panic(err)
  }

  var keys []TweetStruct
  json.Unmarshal(bodyJSON, &body)
  statuses := body["statuses"].([]interface{})
  tweets := findTopTweets(statuses)
  for k := range tweets {
    keys = append(keys, k)
  }

  tweetIndex := rand.Intn(len(keys))
  return keys[tweetIndex]
}

func findTopTweets(statuses []interface{}) map[TweetStruct]bool {
  const THRESHOLD int = 100
  tweets := make(map[TweetStruct]bool, 100)
  for _, status := range statuses {
    statusMap := status.(map[string]interface{})
    retweetCount := int(statusMap["retweet_count"].(float64))
    retweetedStatus, retweetedStatusOk := statusMap["retweeted_status"].(map[string]interface{})
    if retweetedStatusOk && retweetCount >= THRESHOLD {
      tweetID := "https://www.twitter.com/statuses/" + retweetedStatus["id_str"].(string)
      entities := retweetedStatus["entities"].(map[string]interface{})
      medias, mediasOK := entities["media"].([]interface{})
      if mediasOK {
        media := medias[0].(map[string]interface{})
        mediaURL := media["media_url_https"].(string)
        tweets[TweetStruct{id: tweetID, mediaURL: mediaURL}] = true
      }
    }
  }

  return tweets
}
