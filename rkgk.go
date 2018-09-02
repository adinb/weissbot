package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
)

type tweetStruct struct {
	id       string
	mediaURL string
}

func getDailyRkgk(client *http.Client) tweetStruct {
	hashtagRkgkBodyJSON, err := searchTweets(client, "%23rkgk")
	if err != nil {
		log.Panic(err)
	}

	rkgkBodyJSON, err := searchTweets(client, "rkgk")
	if err != nil {
		log.Panic(err)
	}

	rkgkJpBodyJSON, err := searchTweets(client, "%e3%82%89%e3%81%8f%e3%81%8c%e3%81%8d")
	if err != nil {
		log.Panic(err)
	}

	var hashtagRkgkBody map[string]interface{}
	var rkgkBody map[string]interface{}
	var rkgkJpBody map[string]interface{}

	var keys []tweetStruct
	json.Unmarshal(hashtagRkgkBodyJSON, &hashtagRkgkBody)
	json.Unmarshal(rkgkBodyJSON, &rkgkBody)
	json.Unmarshal(rkgkJpBodyJSON, &rkgkJpBody)

	statuses := make([]interface{}, 0)
	statuses = append(statuses, hashtagRkgkBody["statuses"].([]interface{})...)
	statuses = append(statuses, rkgkBody["statuses"].([]interface{})...)
	statuses = append(statuses, rkgkJpBody["statuses"].([]interface{})...)

	tweets := findTopTweets(statuses)
	for k := range tweets {
		keys = append(keys, k)
	}

	tweetIndex := rand.Intn(len(keys))
	return keys[tweetIndex]
}

func findTopTweets(statuses []interface{}) map[tweetStruct]bool {
	const THRESHOLD int = 50
	tweets := make(map[tweetStruct]bool, 100)
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
				tweets[tweetStruct{id: tweetID, mediaURL: mediaURL}] = true
			}
		}
	}

	return tweets
}
