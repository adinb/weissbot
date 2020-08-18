package cotd

type RepositoryContract interface {
	Get() ([]COTD, error)
}

type ServiceContract interface {
	GetCOTD() ([]COTD, error)
}
