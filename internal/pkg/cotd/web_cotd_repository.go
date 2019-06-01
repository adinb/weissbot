package cotd

import (
	"errors"
	"net/http"
	"github.com/antchfx/htmlquery"
)

type WebCOTDRepository struct {
	WebpageURL   string
	ImageBaseURL string
	ImagePath	string
	Client       *http.Client
}

func (r *WebCOTDRepository) Get() ([]COTD, error) {
	resp, err := r.Client.Get(r.WebpageURL)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Couldn't download HTML page")
	}
	defer resp.Body.Close()

	// The error is intentionally ignored because the HTML parser is very permissive.
	doc, _ := htmlquery.Parse(resp.Body)
	cotdList := make([]COTD, 0)
	for _, data := range htmlquery.Find(doc, r.ImagePath) {
		cotd := COTD{ImageURL: r.ImageBaseURL + htmlquery.SelectAttr(data, "src")}
		cotdList = append(cotdList, cotd)
	}

	return cotdList, nil
}