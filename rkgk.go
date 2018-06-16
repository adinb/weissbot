package main

import (
  "log"
  "fmt"
  "os"
  "time"
  "net/http"
  "encoding/json"
  "github.com/joho/godotenv"
)

// GetDailyRkgk will return URL to a rkgk tweet
func GetDailyRkgk(client *http.Client) {
  bodyJSON, err := SearchTweets(client, "rkgk")

  var body map[string]interface{}
  if err != nil {
    log.Panic(err)
  }

  json.Unmarshal(bodyJSON, &body)
  statuses := body["statuses"].([]interface{})
  findTopTweets(statuses)
}

func findTopTweets(statuses []interface{}) {
  const THRESHOLD int = 500
  // var topTweetIDs []string
  for _, status := range statuses {
    statusMap := status.(map[string]interface{})
    // retweetCount := int(statusMap["retweet_count"].(float64))
    retweetedStatus, retweetedStatusOk := statusMap["retweeted_status"].(map[string]interface{})
    if retweetedStatusOk {
      fmt.Println(retweetedStatus["id_str"].(string))
    }

    // if  retweetCount >= THRESHOLD {
    //   statusID := statusMap["id_str"].(string)
    //   fmt.Println("https://www.twitter.com/statuses/" + statusID)
    //   topTweetIDs = append(topTweetIDs, statusID)
    // }
  }
}

func main() {
  if os.Getenv("ENV") != "production" {
    err := godotenv.Load()
    if err != nil {
      log.Fatal(err)
    }
  }

  httpClient := &http.Client {
		Timeout: time.Second * 10,
	}
  GetDailyRkgk(httpClient)
}
