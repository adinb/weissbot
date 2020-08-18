package cotd

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWebCOTDRepositoryGet(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadFile("testdata/vanguard_cotd_page.html")
		if err != nil {
			t.Error("Failed to load test file")
		}
		w.Write(data)
	}))
	defer ts.Close()

	var client http.Client
	repo := WebCOTDRepository{
		WebpageURL:   ts.URL,
		ImagePath:    "//p[contains(@class, 'text-center')]/img[contains(@class, 'alignnone')]",
		ImageBaseURL: "",
		Client:       &client,
	}

	cotdList, err := repo.Get()
	if err != nil {
		t.Error("Failed to retrieve cotd", err)
	}
	if len(cotdList) != 1 {
		t.Error("Failed to parse the page")
	}
}

func TestRetrieveServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	}))
	defer ts.Close()

	var client http.Client
	repo := WebCOTDRepository{
		WebpageURL:   ts.URL,
		ImagePath:    "//p[contains(@class, 'text-center')]/img[contains(@class, 'alignnone')]",
		ImageBaseURL: "",
		Client:       &client,
	}
	_, err := repo.Get()
	if err == nil {
		t.Error("No error returned")
	}
}

func TestRetrieveConnectionError(t *testing.T) {
	var ts *httptest.Server
	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ts.CloseClientConnections()
	}))
	defer ts.Close()

	var client http.Client
	repo := WebCOTDRepository{
		WebpageURL:   ts.URL,
		ImagePath:    "//p[contains(@class, 'text-center')]/img[contains(@class, 'alignnone')]",
		ImageBaseURL: "",
		Client:       &client,
	}
	_, err := repo.Get()
	if err == nil {
		t.Error("No error returned")
	}
}
