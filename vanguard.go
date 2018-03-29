package main

import (
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

const baseVGURL = "http://cf-vanguard.com"
const cotdVGPageURL = "http://cf-vanguard.com/todays-card/"
const vanguardName = "vg"

// GetVGCotd returns array of COTD image URLs
func GetVGCotd() []string {
	return getVGCotdURLFromPage(RetrievePage(cotdVGPageURL))
}

func getVGCotdURLFromPage(doc *html.Node) []string {
	cotdUrls := make([]string, 0)

	for _, data := range htmlquery.Find(doc, "//p[contains(@class, 'taC mb08')]/img") {
		cotdUrls = append(cotdUrls, baseVGURL+htmlquery.SelectAttr(data, "src"))
	}

	return cotdUrls
}
