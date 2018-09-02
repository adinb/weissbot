package main

import (
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

const baseBFURL = "https://fc-buddyfight.com/"
const cotdBFURL = "https://fc-buddyfight.com/todays-card/"

const bfName = "bf"

func getBFCotd() []string {
	return getBFCotdURLFromPage(RetrievePage(cotdBFURL))
}

func getBFCotdURLFromPage(doc *html.Node) []string {
	cotdUrls := make([]string, 0)

	for _, data := range htmlquery.Find(doc, "//div[contains(@class, 'lp_bg')]/div/div/img") {
		cotdUrls = append(cotdUrls, htmlquery.SelectAttr(data, "src"))
	}

	return cotdUrls
}
