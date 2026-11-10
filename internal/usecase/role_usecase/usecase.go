package role_usecase

type Usecase struct {
	store roleStore
}

func New(store roleStore) *Usecase {
	return &Usecase{store: store}
}
