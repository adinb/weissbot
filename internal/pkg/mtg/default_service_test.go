package mtg

import (
	"os"
	"errors"
	"strings"
	"testing"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

type mockMTGRepository struct {
	isFailed bool
	cards []*MagicCard
}

func (r *mockMTGRepository) Find(query string) ([]*MagicCard, error){
	if r.isFailed {
		return nil, errors.New("Failed to find the card")
	}

	var results []*MagicCard
	for _, card := range(r.cards) {
		if strings.Contains(card.Name, query) {
			results = append(results, card)
		}
	}

	return results, nil
}

func TestDefaultServiceSearchCardByNameFailed(t *testing.T) {
	var repo mockMTGRepository
	repo.isFailed = true
	service := DefaultService{Repo: &repo}
	_, err := service.SearchCardByName("bolas")
	if err == nil {
		t.Error("Service didn't return error")
	}
}

func TestDefaultServiceSearchCardByNameSuccess(t *testing.T) {
	var cards []*MagicCard
	cards = append(cards, &MagicCard{Name: "Nicol Bolas, The Ravager"})
	cards = append(cards, &MagicCard{Name: "Ajani The Greathearted"})
	cards = append(cards, &MagicCard{Name: "Nicol Bolas, Dragon-God"})
	cards = append(cards, &MagicCard{Name: "Teferi, Time Raveler"})
	cards = append(cards, &MagicCard{Name: "Jace, Wielder of Mysteries"})

	var repo mockMTGRepository
	repo.isFailed = false
	repo.cards = cards
	service := DefaultService{Repo: &repo}
	results, err := service.SearchCardByName("Bolas")
	if err != nil {
		t.Error("Failed to search card by name")
	}

	if len(results) != 2 {
		t.Error("Incorrect returned results", results)
	}
}

func TestCombineImage(t *testing.T) {
	imgA, err := os.Open("testdata/a.png")
	if err != nil {
		t.Error("Failed to load test file")
	}

	pngA, err := png.Decode(imgA)
	if err != nil {
		t.Error("Failed to decode image")
	}

	imgB, err := os.Open("testdata/b.png")
	if err != nil {
		t.Error("Failed to load test file")
	}

	pngB, err := png.Decode(imgB)
	if err != nil {
		t.Error("Failed to decode image")
	}

	_, err = jpeg.Decode(CombineImage(pngA, pngB))
	if err != nil {
		t.Error("Failed to decode the resulting image")
	}
}

func TestRetrievePNGSuccess(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadFile("testdata/a.png")
		if err != nil {
			t.Error("Failed to load test file")
		}
		w.Write(data)
	}))
	defer ts.Close()

	_, err := RetrievePNG(ts.URL)
	if err != nil {
		t.Error("Failed to retrieve image", err)
	}
}

func TestRetrievePNGHTTPFail(t *testing.T) {
	var ts *httptest.Server
	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ts.CloseClientConnections()
	}))
	defer ts.Close()

	_, err := RetrievePNG(ts.URL)
	if err == nil {
		t.Error("No error returned")
	}
}

func TestRetrievePNGDecodeFail(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadFile("testdata/dummy_file.png")
		if err != nil {
			t.Error("Failed to load test file")
		}
		w.Write(data)
	}))
	defer ts.Close()

	_, err := RetrievePNG(ts.URL)
	if err == nil {
		t.Error("No error returned")
	}
}