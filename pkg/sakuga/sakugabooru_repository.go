package sakuga

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/antchfx/htmlquery"
)

type SakugabooruRepository struct {
	Client  *http.Client
	BaseURL string
}

func (s SakugabooruRepository) Get() (Sakuga, error) {
	const randomLimit int = 70000

	resp, err := s.Client.Get(s.BaseURL)
	if err != nil {
		return Sakuga{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return Sakuga{}, errors.New("Couldn't download HTML page: " + strconv.Itoa(resp.StatusCode))
	}

	defer resp.Body.Close()

	// The error is intentionally ignored because the HTML parser is very permissive.
	doc, _ := htmlquery.Parse(resp.Body)
	data := htmlquery.FindOne(doc, "//*[@id=\"highres\"]")

	return Sakuga{URL: htmlquery.SelectAttr(data, "href")}, nil
}
