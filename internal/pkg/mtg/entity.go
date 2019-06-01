package mtg

import (
	"encoding/json"
	"strings"
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

type MagicCardRepository interface {
	Find(query string) ([]*MagicCard, error)
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
			m.Faces = append(m.Faces, extractCardFace(face.(map[string]interface{})))
		}
	} else {
		m.Faces = append(m.Faces, extractCardFace(rawCard))
	}

	return nil
}

func extractCardFace(face map[string]interface{}) CardFace {
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

	if face["oracle_text"] != "" {
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
