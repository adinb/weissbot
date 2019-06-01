package twitter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Client interface {
	FindTweets(query string) ([]Tweet, error)
}

type client struct {
	httpClient *http.Client
	baseURL    string
	token      string
}

// Tweet contains simplified tweet data
type Tweet struct {
	IDStr          string
	Text           string
	UserName       string
	UserScreenName string
	RetweetCount   int
	FavoriteCount  int
	MediaURLs      []string
	URL            string
}

type twitterResponse struct {
	Statuses []Tweet `json:statuses`
}

func New(c *http.Client, baseURL string, token string) client {
	return client{httpClient: c, baseURL: baseURL, token: token}
}

func (c *client) FindTweets(query string) ([]Tweet, error) {
	searchQuery := c.baseURL + query + "&count=100"
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

	var r twitterResponse
	err = json.Unmarshal(buffer, &r)
	if err != nil {
		return nil, err
	}
	return r.Statuses, nil
}

func (t *Tweet) UnmarshalJSON(b []byte) error {
	var tweet map[string]interface{}
	var chosenDataSource map[string]interface{}
	json.Unmarshal(b, &tweet)

	rts, rtsOk := tweet["retweeted_status"]
	if rtsOk {
		chosenDataSource = rts.(map[string]interface{})
	} else {
		chosenDataSource = tweet
	}

	t.IDStr = chosenDataSource["id_str"].(string)
	t.Text = chosenDataSource["text"].(string)
	t.RetweetCount = int(chosenDataSource["retweet_count"].(float64))
	t.FavoriteCount = int(chosenDataSource["favorite_count"].(float64))

	user := chosenDataSource["user"].(map[string](interface{}))
	t.UserName = user["name"].(string)
	t.UserScreenName = user["screen_name"].(string)
	t.URL = fmt.Sprintf("https://www.twitter.com/%s/status/%s", t.UserScreenName, t.IDStr)

	if chosenDataSource["extended_entities"] != nil {
		extendedEntities := chosenDataSource["extended_entities"].(map[string](interface{}))
		media := extendedEntities["media"].([]interface{})
		for _, m := range media {
			t.MediaURLs = append(t.MediaURLs, m.(map[string](interface{}))["media_url_https"].(string))
		}
	}

	return nil
}
