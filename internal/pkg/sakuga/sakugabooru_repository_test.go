package sakuga

import (
	"testing"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

func TestSakugabooruRepositoryGetSuccess(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadFile("testdata/success_page.html")
		if err != nil {
			t.Error("Failed to load test file")
		}
		w.Write(data)
	}))
	defer ts.Close()
	repo := SakugabooruRepository{Client: new(http.Client), BaseURL: ts.URL}
	sakuga, _ := repo.Get()

	if sakuga.URL != "https://sakugabooru.com/data/ed900fb57e3031f19f95077e0ebdfead.mp4" {
		t.Error("Failed to get the video URL")
	}
}

func TestSakugabooruRepositoryConnectionError(t *testing.T) {
	var ts *httptest.Server
	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ts.CloseClientConnections()
	}))
	defer ts.Close()
	repo := SakugabooruRepository{Client: new(http.Client), BaseURL: ts.URL}
	_, err := repo.Get()
	if err == nil {
		t.Error("Failed to return error")
	}
}

func TestSakugabooruRepositoryHTTPError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	}))
	defer ts.Close()
	repo := SakugabooruRepository{Client: new(http.Client), BaseURL: ts.URL}
	_, err := repo.Get()
	if err == nil {
		t.Error("Failed to return error")
	}
}