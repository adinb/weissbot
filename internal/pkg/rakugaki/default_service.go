package rakugaki

import (
	"errors"
	"math/rand"
)

type DefaultService struct {
	Repo RepositoryContract
}

func filterRakugakiByRating(rakugakiList []Rakugaki, rating int) []Rakugaki {
	var filteredRakugakiList []Rakugaki
	for _, r := range rakugakiList {
		if r.Rating >= rating {
			filteredRakugakiList = append(filteredRakugakiList, r)
		}
	}

	return filteredRakugakiList
}

func (s *DefaultService) GetTopRakugaki(threshold int) (Rakugaki, error) {
	const RKGK_QUERY = "rkgk"
	const HASHTAG_RKGK_QUERY = "%23rkgk"
	const JP_RKGK_QUERY = "%e3%82%89%e3%81%8f%e3%81%8c%e3%81%8d"

	rakugakiList, _ := s.Repo.List(RKGK_QUERY)

	hashtagRakugakiList, _ := s.Repo.List(HASHTAG_RKGK_QUERY)

	jpRakugakiList, _ := s.Repo.List(JP_RKGK_QUERY)

	combinedList := append(rakugakiList, hashtagRakugakiList...)
	combinedList = append(combinedList, jpRakugakiList...)

	if len(combinedList) == 0 {
		return Rakugaki{}, errors.New("No rakugaki found")
	}

	filteredList := filterRakugakiByRating(combinedList, threshold)
	if len(filteredList) == 0 {
		return Rakugaki{}, errors.New("No rakugaki with rating over threshold found")
	}
	return filteredList[rand.Intn(len(rakugakiList))], nil
}
