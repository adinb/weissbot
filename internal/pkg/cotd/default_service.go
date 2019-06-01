package cotd

import (
	"errors"
)
type DefaultService struct {
	Repo RepositoryContract
}

func (s *DefaultService) GetCOTD() ([]COTD, error) {
	cotdList, err := s.Repo.Get()
	if err != nil {
		return nil, err
	}

	if len(cotdList) == 0 {
		return nil, errors.New("No card found")
	}

	return cotdList, nil
}
