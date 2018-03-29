package main

import (
	"fmt"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

const baseWSURL = "https://ws-tcg.com"
const cotdWSPageURL = "https://ws-tcg.com/todays-card/"
const wsName = "ws"

// GetWSCotd returns array of COTD image URLs
func GetWSCotd() []string {
	return getWSCotdURLFromPage(RetrievePage(cotdWSPageURL))
}

func getWSCotdURLFromPage(doc *html.Node) []string {
	cotdUrls := make([]string, 0)

	for _, data := range htmlquery.Find(doc, "//div[contains(@class, 'entry-content')]/p/img") {
		cotdUrls = append(cotdUrls, baseWSURL+htmlquery.SelectAttr(data, "src"))
		fmt.Println(cotdUrls)
	}

	return cotdUrls
}
