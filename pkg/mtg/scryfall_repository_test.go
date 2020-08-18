package mtg

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestScryfallRepositoryFindSuccess(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadFile("testdata/bolas.json")
		if err != nil {
			t.Error("Failed to load test file")
		}
		w.Write(data)
	}))
	defer ts.Close()

	var client http.Client	
	repo := ScryfallRepository{Client: &client, BaseURL: ts.URL}
	_, err := repo.Find("Bolas")
	if err != nil {
		t.Error("Failed to find card from scryfall", err)
	}
}

func TestScryfallRepositoryHTTPFail(t *testing.T) {
	var ts *httptest.Server
	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ts.CloseClientConnections()
	}))
	defer ts.Close()

	var client http.Client	
	repo := ScryfallRepository{Client: &client, BaseURL: ts.URL}
	_, err := repo.Find("Bolas")
	if err == nil {
		t.Error("No error returned")
	}
}

func TestScryfallRepositoryResourceFail(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadFile("testdata/broken_card.json")
		if err != nil {
			t.Error("Failed to load test file", err)
		}
		w.Write(data)
	}))
	defer ts.Close()

	var client http.Client	
	repo := ScryfallRepository{Client: &client, BaseURL: ts.URL}
	_, err := repo.Find("Bolas")
	if err == nil {
		t.Error("No error returned")
	}
}

func TestMagicCardUnmarshallJSONNoFlavorTextCard( t *testing.T) {
	data, err := ioutil.ReadFile("testdata/gideon_blackblade.json")
	if err != nil {
		t.Error("Failed to load test file")
	}

	var card MagicCard
	err = json.Unmarshal(data, &card)
	if err != nil {
		t.Error("Failed to decode JSON data", err)
	}
}

func TestMagicCardUnmarshallJSONNoTextCard( t *testing.T) {
	data, err := ioutil.ReadFile("testdata/ironclad_krovod.json")
	if err != nil {
		t.Error("Failed to load test file")
	}

	var card MagicCard
	err = json.Unmarshal(data, &card)
	if err != nil {
		t.Error("Failed to decode JSON data", err)
	}
}