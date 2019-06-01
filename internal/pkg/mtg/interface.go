package mtg

type RepositoryContract interface {
	Find(query string) ([]*MagicCard, error)
}

type ServiceContract interface {
	SearchCardByName(name string) ([]*MagicCard, error)
}