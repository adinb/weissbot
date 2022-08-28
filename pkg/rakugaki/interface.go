package rakugaki

type ServiceContract interface {
	GetTopRakugaki(threshold int) (Rakugaki, error)
}

type RepositoryContract interface {
	List(query string) ([]Rakugaki, error)
}
