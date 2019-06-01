package meta

type Meta struct {
	Status string `json:status`
}

func PublishMetaChanges(status Meta, sc chan<- Meta) {
	sc <- status
}
