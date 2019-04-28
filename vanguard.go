package main

import (
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

const baseVGURL = "https://cf-vanguard.com"
const cotdVGPageURL = "https://cf-vanguard.com/todays-card/"
const vanguardName = "vg"

// GetVGCotd returns array of COTD image URLs
func getVGCotd() []string {
	return getVGCotdURLFromPage(RetrievePage(cotdVGPageURL))
}

func getVGCotdURLFromPage(doc *html.Node) []string {
	cotdUrls := make([]string, 0)

	for _, data := range htmlquery.Find(doc, "//p[contains(@class, 'text-center')]/img[contains(@class, 'alignnone')]") {
		cotdUrls = append(cotdUrls, htmlquery.SelectAttr(data, "src"))
	}

	return cotdUrls
}
