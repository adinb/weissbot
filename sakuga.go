package main

import (
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"math/rand"
	"fmt"
)

const baseSakugaURI = "https://www.sakugabooru.com/post/show/"
const randomLimit = 70000

func GetSakuga() string {
	videoUrl := ""
	for videoUrl == "" {
		id := rand.Intn(randomLimit)
		uri := fmt.Sprintf("%s%d", baseSakugaURI, id)
		html := RetrievePage(uri)
		videoUrl = getVideoURL(html)
	}
	
	return videoUrl
}

func getVideoURL(doc *html.Node) string {
	data := htmlquery.FindOne(doc, "//*[@id=\"highres\"]")
	if data == nil {
		return ""
	} else {
		return htmlquery.SelectAttr(data, "href")	
	}
}