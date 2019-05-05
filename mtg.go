package main

import (
	"encoding/json"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/fogleman/gg"
)

const scryfallBaseAPIURL = "https://api.scryfall.com"

// Legalities contains list of legal format of a card
type Legalities []string

// ImageURIs hold types of image URIS
type ImageURIs struct {
	Small      string
	Normal     string
	Large      string
	PNG        string
	ArtCrop    string
	BorderCrop string
}

// ScryfallObjectResource represents object resource from scryfall API
type ScryfallObjectResource struct {
	Object string      `json:"object"`
	Data   []MagicCard `json:"data"`
}

// MagicCard represents Magic Card
type MagicCard struct {
	ID          string
	Name        string
	Faces       []CardFace
	Legalities  Legalities
	ScryfallURI string
	Rarity      string
	Set         string
	SetName     string
	ReleaseDate string
}

// CardFace represents face of the card
type CardFace struct {
	Name       string
	ManaCost   string
	TypeLine   string
	Text       string
	Colors     []string
	Power      string
	Toughness  string
	Artist     string
	FlavorText string
	ImageURIs  ImageURIs
}

// UnmarshalJSON will load MagicCard data from JSON
func (m *MagicCard) UnmarshalJSON(b []byte) error {
	var rawCard map[string]interface{}
	err := json.Unmarshal(b, &rawCard)
	if err != nil {
		return err
	}

	m.ID = rawCard["id"].(string)
	m.Name = rawCard["name"].(string)
	m.ScryfallURI = rawCard["scryfall_uri"].(string)
	m.Rarity = rawCard["rarity"].(string)
	m.Set = rawCard["set"].(string)
	m.SetName = rawCard["set_name"].(string)
	m.ReleaseDate = rawCard["released_at"].(string)

	for k, v := range rawCard["legalities"].(map[string]interface{}) {
		if v.(string) == "legal" {
			k = strings.Title(k)
			m.Legalities = append(m.Legalities, k)
		}
	}

	if rawCard["card_faces"] != nil {
		for _, face := range rawCard["card_faces"].([]interface{}) {
			m.Faces = append(m.Faces, createFaceFromMap(face.(map[string]interface{})))
		}
	} else {
		m.Faces = append(m.Faces, createFaceFromMap(rawCard))
	}

	return nil
}

func createFaceFromMap(face map[string]interface{}) CardFace {
	var cardFace CardFace
	cardFace.Name = face["name"].(string)
	cardFace.TypeLine = face["type_line"].(string)
	cardFace.ManaCost = face["mana_cost"].(string)

	if face["power"] != nil {
		cardFace.Power = face["power"].(string)
	}

	if face["toughness"] != nil {
		cardFace.Toughness = face["toughness"].(string)
	}

	cardFace.Artist = face["artist"].(string)

	if face["oracle_text"] != nil {
		cardFace.Text = face["oracle_text"].(string)
	} else {
		cardFace.Text = "-"
	}

	if face["flavor_text"] != nil {
		cardFace.FlavorText = face["flavor_text"].(string)
	} else {
		cardFace.FlavorText = "-"
	}

	if face["colors"] != nil {
		for _, color := range face["colors"].([]interface{}) {
			cardFace.Colors = append(cardFace.Colors, color.(string))
		}
	}

	if face["image_uris"] != nil {
		cardFace.ImageURIs.Small = face["image_uris"].(map[string]interface{})["small"].(string)
		cardFace.ImageURIs.Normal = face["image_uris"].(map[string]interface{})["normal"].(string)
		cardFace.ImageURIs.Large = face["image_uris"].(map[string]interface{})["large"].(string)
		cardFace.ImageURIs.PNG = face["image_uris"].(map[string]interface{})["png"].(string)
		cardFace.ImageURIs.ArtCrop = face["image_uris"].(map[string]interface{})["art_crop"].(string)
		cardFace.ImageURIs.BorderCrop = face["image_uris"].(map[string]interface{})["border_crop"].(string)
	}

	return cardFace
}

// LoadScryfallObjectResourceFromJSON loads the resource from JSON
func LoadScryfallObjectResourceFromJSON(b []byte) (ScryfallObjectResource, error) {
	var s ScryfallObjectResource
	err := json.Unmarshal(b, &s)
	if err != nil {
		return s, err
	}
	return s, nil
}

// SearchMagicCard retrieve resources from Scryfall API, returning internal representation of Magic cards
func SearchMagicCard(name string) ([]MagicCard, error) {
	resp, err := http.Get(scryfallBaseAPIURL + "/cards/search?unique=cards&q=" + url.QueryEscape(name))
	if err != nil {
		return nil, err
	}

	buffer, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	obj, err := LoadScryfallObjectResourceFromJSON(buffer)
	if err != nil {
		return nil, err
	}

	cards := obj.Data
	return cards, nil
}

// GenerateCombinedTwoFaceCardImage returns a combined image of a two-faced card
func GenerateCombinedTwoFaceCardImage(card MagicCard) (io.Reader, error) {
	resp, err := http.Get(card.Faces[0].ImageURIs.PNG)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	resp2, err := http.Get(card.Faces[1].ImageURIs.PNG)
	if err != nil {
		return nil, err
	}
	defer resp2.Body.Close()

	imgA, err := png.Decode(resp.Body)
	if err != nil {
		return nil, err
	}

	imgB, err := png.Decode(resp2.Body)
	if err != nil {
		return nil, err
	}

	w := imgA.Bounds().Size().X * 2
	h := imgA.Bounds().Size().Y
	dc := gg.NewContext(w, h)

	dc.DrawImage(imgA, 0, 0)
	dc.DrawImage(imgB, imgA.Bounds().Size().X, 0)

	reader, writer := io.Pipe()
	go func() {
		jpeg.Encode(writer, dc.Image(), &jpeg.Options{Quality: 80})
		writer.Close()
	}()

	return reader, nil
}
