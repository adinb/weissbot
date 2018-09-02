package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func searchTweets(client *http.Client, query string) ([]byte, error) {
	bearerToken := os.Getenv("TWITTER_TOKEN")
	searchQuery := "https://api.twitter.com/1.1/search/tweets.json?q=" + query + "&count=100"
	req, err := http.NewRequest("GET", searchQuery, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+bearerToken)
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
