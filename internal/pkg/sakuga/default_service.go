package sakuga

type DefaultService struct {
	Repo RepositoryContract
}

func (s *DefaultService) GetSakuga() (Sakuga, error) {
	sakuga, err := s.Repo.Get()
	if err != nil {
		return Sakuga{}, err
	}

	return sakuga, nil
}
