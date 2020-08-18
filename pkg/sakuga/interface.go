package sakuga

type RepositoryContract interface {
	Get() (Sakuga, error)
}


type ServiceContract interface {
	GetSakuga() (Sakuga, error)
}