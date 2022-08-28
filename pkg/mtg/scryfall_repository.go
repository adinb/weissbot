package mtg

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

type ScryfallRepository struct {
	Client  *http.Client
	BaseURL string
}

// ScryfallObjectResource represents object resource from scryfall API
type ScryfallObjectResource struct {
	Object string       `json:"object"`
	Data   []*MagicCard `json:"data"`
}

func (r *ScryfallRepository) Find(query string) ([]*MagicCard, error) {
	resp, err := http.Get(r.BaseURL + "/cards/search?unique=cards&q=" + url.QueryEscape(query))
	if err != nil {
		return nil, err
	}

	buffer, _ := ioutil.ReadAll(resp.Body)

	var s ScryfallObjectResource
	err = json.Unmarshal(buffer, &s)
	if err != nil {
		return nil, err
	}

	cards := s.Data
	return cards, nil
}
