package twitter

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/adinb/weissbot/pkg/command"
)

type Client interface {
	GetTweetResults(query string) (command.Result, error)
	FindTweets(query string) ([]Tweet, error)
	GetLastTweetsByAuthor(authorHandle string) (command.Result, error)
	BuildTweetSearchURL(q string, expansions []string, fields SearchResultField, maxResults uint) string
	SearchRecentTweets(q string, expansions []string, fields SearchResultField, maxResults uint) (TwitterResponse, error)
}

type client struct {
	httpClient *http.Client
	baseURL    string
	token      string
}

// Tweet contains simplified tweet data
type Tweet struct {
	IDStr          string   `json:"id"`
	Text           string   `json:"text"`
	CreatedAt      string   `json:"created_at"`
	UserScreenName string   `json:",omitempty"`
	RetweetCount   int      `json:",omitempty"`
	FavoriteCount  int      `json:",omitempty"`
	MediaURLs      []string `json:",omitempty"`
	MediaURL       string   `json:",omitempty"`
	URL            string   `json:",omitempty"`
}

type Include struct {
	Users []User
	Media []Media
}

type Media struct {
	URL string `json:"url"`
}

type User struct {
	UserName string `json:"username"`
}

type TwitterResponse struct {
	Statuses []Tweet `json:"data"`
	Include  Include `json:"includes"`
}

type SearchResultField struct {
	Tweet []string
	Media []string `json:"media"`
	User  []string `json:"users"`
}

const tweetSearchBaseURL = "https://api.twitter.com/2/tweets/search/recent?query="

func New(c *http.Client, token string) Client {
	return client{httpClient: c, baseURL: tweetSearchBaseURL, token: token}
}

func (c client) BuildTweetSearchURL(
	q string,
	expansions []string,
	fields SearchResultField, maxResults uint) string {
	query := url.QueryEscape(q)
	query = strings.Replace(query, "+", "%20", -1)
	var b strings.Builder
	b.WriteString(c.baseURL)
	b.WriteString(query)

	if len(expansions) > 0 {
		b.WriteString("&expansions=")
		for i, expansion := range expansions {
			b.WriteString(expansion)
			if i < (len(expansions) - 1) {
				b.WriteString(",")
			}
		}
	}

	if len(fields.Tweet) > 0 {
		b.WriteString("&tweet.fields=")
		for i, tweetField := range fields.Tweet {
			b.WriteString(tweetField)
			if i < (len(fields.Tweet) - 1) {
				b.WriteString(",")
			}
		}
	}

	if len(fields.User) > 0 {
		b.WriteString("&user.fields=")
		for i, userField := range fields.User {
			b.WriteString(userField)
			if i < (len(fields.User) - 1) {
				b.WriteString(",")
			}
		}
	}

	if len(fields.Media) > 0 {
		b.WriteString("&media.fields=")
		for i, mediaField := range fields.Media {
			b.WriteString(mediaField)
			if i < (len(fields.Media) - 1) {
				b.WriteString(",")
			}
		}
	}

	b.WriteString(fmt.Sprintf("&max_results=%d", maxResults))

	return b.String()
}

func (c client) SearchRecentTweets(q string,
	expansions []string, fields SearchResultField, maxResults uint) (TwitterResponse, error) {
	url := c.BuildTweetSearchURL(q, expansions, fields, maxResults)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return TwitterResponse{}, err
	}

	req.Header.Add("Authorization", "Bearer "+c.token)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return TwitterResponse{}, err
	}

	buffer, _ := ioutil.ReadAll(resp.Body)

	var r TwitterResponse
	err = json.Unmarshal(buffer, &r)
	if err != nil {
		return TwitterResponse{}, err
	}

	return r, nil
}

func (c client) GetLastTweetsByAuthor(authorHandle string) (command.Result, error) {
	query := fmt.Sprintf("from: %s -is:retweet -is:reply -is:quote", authorHandle)
	expansions := []string{"author_id"}
	resultField := SearchResultField{
		Tweet: []string{"id", "text", "created_at"},
	}

	t, err := c.SearchRecentTweets(query, expansions, resultField, 10)
	if err != nil {
		return command.Result{}, nil
	}

	if len(t.Statuses) == 0 {
		return command.Result{}, errors.New("no tweet found")
	}

	tweet := t.Statuses[0]
	linkToTweet := fmt.Sprintf("https://www.twitter.com/%s/status/%s",
		t.Include.Users[0].UserName,
		tweet.IDStr,
	)
	message := fmt.Sprintf("@%s tweeted at %s\n%s\n%s",
		t.Include.Users[0].UserName,
		tweet.CreatedAt,
		tweet.Text,
		linkToTweet)
	result := command.Result{Messages: []string{message}}
	return result, nil
}

func (c client) GetTweetResults(query string) (result command.Result, e error) {
	result = command.Result{}
	tweets, e := c.FindTweets(query)
	if e != nil {
		return
	}

	rand.Seed(time.Now().Unix())
	i := rand.Intn(len(tweets))
	imageMessage := command.ImageMessage{
		Message:  tweets[i].Text,
		ImageURL: tweets[i].MediaURL,
	}

	result.ImageMessages = append(result.ImageMessages, imageMessage)

	return
}

func (c client) FindTweets(query string) ([]Tweet, error) {
	searchQuery := c.baseURL + query + "%20-is:retweet%20-is:reply%20-is:quote%20has:images&expansions=author_id,attachments.media_keys&max_results=100&tweet.fields=id,text,created_at&media.fields=url"
	req, err := http.NewRequest("GET", searchQuery, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+c.token)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	buffer, _ := ioutil.ReadAll(resp.Body)

	var r TwitterResponse
	err = json.Unmarshal(buffer, &r)
	if err != nil {
		return nil, err
	}

	for i := range r.Statuses {
		r.Statuses[i].MediaURL = r.Include.Media[i].URL
	}

	return r.Statuses, nil
}

//func (t *Tweet) UnmarshalJSON(b []byte) error {
//	var tweet map[string]interface{}
//	var chosenDataSource map[string]interface{}
//	json.Unmarshal(b, &tweet)
//
//	rts, rtsOk := tweet["retweeted_status"]
//	if rtsOk {
//		chosenDataSource = rts.(map[string]interface{})
//	} else {
//		chosenDataSource = tweet
//	}
//
//	t.IDStr = chosenDataSource["id_str"].(string)
//	t.Text = chosenDataSource["text"].(string)
//	t.RetweetCount = int(chosenDataSource["retweet_count"].(float64))
//	t.FavoriteCount = int(chosenDataSource["favorite_count"].(float64))
//
//	user := chosenDataSource["user"].(map[string](interface{}))
//	t.UserName = user["name"].(string)
//	t.UserScreenName = user["screen_name"].(string)
//	t.URL = fmt.Sprintf("https://www.twitter.com/%s/status/%s", t.UserScreenName, t.IDStr)
//
//	if chosenDataSource["extended_entities"] != nil {
//		extendedEntities := chosenDataSource["extended_entities"].(map[string](interface{}))
//		media := extendedEntities["media"].([]interface{})
//		for _, m := range media {
//			t.MediaURLs = append(t.MediaURLs, m.(map[string](interface{}))["media_url_https"].(string))
//		}
//	}
//
//	return nil
//}
