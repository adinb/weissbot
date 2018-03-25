package main

import (
	"net/http"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

func retrieveCotdPage() *html.Node {
	resp, err := http.Get("http://cf-vanguard.com/todays-card/")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	doc, err := htmlquery.Parse(resp.Body)
	if err != nil {
		panic(err)
	}

	return doc
}

// GetCotd returns array of COTD image URLs
func GetCotd() []string {
	return getCotdURLFromPage(retrieveCotdPage())
}

func getCotdURLFromPage(doc *html.Node) []string {
	baseurl := "http://cf-vanguard.com"
	cotdUrls := make([]string, 0)

	for _, data := range htmlquery.Find(doc, "//p[contains(@class, 'taC mb08')]/img") {
		cotdUrls = append(cotdUrls, baseurl+htmlquery.SelectAttr(data, "src"))
	}

	return cotdUrls
}
