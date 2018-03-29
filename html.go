package main

import (
	"net/http"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

// RetrievePage get parsed html page
func RetrievePage(url string) *html.Node {
	resp, err := http.Get(url)
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
