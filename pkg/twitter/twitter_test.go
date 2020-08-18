package twitter

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFindTweets(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadFile("testdata/response.json")
		if err != nil {
			t.Error("Failed to load test file")
		}
		w.Write(data)
	}))
	defer ts.Close()
	var httpClient http.Client
	client := New(&httpClient, ts.URL+"?q=", "token")
	tweets, err := client.FindTweets("rkgk")

	if err != nil {
		t.Error("Error while trying to get Tweets", err)
		return
	}

	if len(tweets) == 0 {
		t.Error("Failed to parse tweets")
	}
}

func TestFindTweetsInvalidURL(t *testing.T) {
	var httpClient http.Client
	client := New(&httpClient, ":invalid/url", "token")
	_, err := client.FindTweets("")

	if err == nil {
		t.Error("No error returned")
		return
	}
}

func TestFindTweetsConnectionError(t *testing.T) {
	var ts *httptest.Server
	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ts.CloseClientConnections()
	}))
	defer ts.Close()
	var httpClient http.Client
	client := New(&httpClient, ts.URL+"?q=", "token")
	_, err := client.FindTweets("rkgk")

	if err == nil {
		t.Error("No error returned", err)
		return
	}
}

func TestFindTweetsInvalidBody(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadFile("testdata/broken_json.json")
		if err != nil {
			t.Error("Failed to load test file")
		}
		w.Write(data)
	}))
	defer ts.Close()
	var httpClient http.Client
	client := New(&httpClient, ts.URL+"?q=", "token")
	_, err := client.FindTweets("rkgk")

	if err == nil {
		t.Error("No error returned", err)
		return
	}
}

