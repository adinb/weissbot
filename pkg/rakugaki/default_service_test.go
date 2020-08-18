package rakugaki

import (
	"errors"
	"testing"
)

type mockRakugakiRepository struct {
	isFailed bool
	rakugakiList []Rakugaki
}

func (m *mockRakugakiRepository) List(query string) ([]Rakugaki, error) {
	if m.isFailed {
		return m.rakugakiList, errors.New("Failed to get rakugaki")
	}
	
	return m.rakugakiList, nil
}

func TestDefaultServiceGetTopRakugakiSuccess(t *testing.T) {
	var rakugakiList []Rakugaki
	rakugakiList = append(rakugakiList, Rakugaki{ImageURL: "url", SourceURL: "url2", Rating: 50})
	rakugakiList = append(rakugakiList, Rakugaki{ImageURL: "url", SourceURL: "url2", Rating: 100})
	rakugakiList = append(rakugakiList, Rakugaki{ImageURL: "url", SourceURL: "url3", Rating: 25})
	repo := mockRakugakiRepository{isFailed: false, rakugakiList: rakugakiList}
	service := DefaultService{Repo: &repo}
	rakugaki, err := service.GetTopRakugaki(100)
	if err != nil {
		t.Error("Failed to get top rakugaki", err)
		return
	}

	if rakugaki.Rating < 100 {
		t.Error("Filter is not working properly")
	}
}

func TestDefaultServiceGetTopRakugakiNotFound(t *testing.T) {
	var rakugakiList []Rakugaki
	rakugakiList = append(rakugakiList, Rakugaki{ImageURL: "url", SourceURL: "url2", Rating: 50})
	repo := mockRakugakiRepository{isFailed: false, rakugakiList: rakugakiList}
	service := DefaultService{Repo: &repo}
	_, err := service.GetTopRakugaki(100)
	if err == nil {
		t.Error("Failed to throw error", err)
		return
	}
}

func TestDefaultServiceGetTopRakugakiFail(t *testing.T) {
	var rakugakiList []Rakugaki
	repo := mockRakugakiRepository{isFailed: false, rakugakiList: rakugakiList}
	service := DefaultService{Repo: &repo}
	_, err := service.GetTopRakugaki(100)
	if err == nil {
		t.Error("Failed to return error")
	}
}