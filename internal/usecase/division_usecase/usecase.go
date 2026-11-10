package division_usecase

type Usecase struct {
	store divisionStore
}

func New(store divisionStore) *Usecase {
	return &Usecase{store: store}
}