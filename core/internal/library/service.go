package library

type Service struct {
	store Store
}

func NewApp(store Store) *Service {

	return &Service{}
}
