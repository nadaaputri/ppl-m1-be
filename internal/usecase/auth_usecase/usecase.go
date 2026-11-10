package auth_usecase

import (
	"ppl-m1-be/internal/usecase/audit_usecase"
	"ppl-m1-be/internal/usecase/utils"
)

type Usecase struct {
	store     userStore
	tokens    *utils.TokenMaker
	audit     audit_usecase.Logger
	dummyHash string
}

func New(store userStore, tokens *utils.TokenMaker, audit audit_usecase.Logger) *Usecase {
	return &Usecase{
		store:     store,
		tokens:    tokens,
		audit:     audit,
		dummyHash: newDummyHash(),
	}
}

func newDummyHash() string {
	hash, err := utils.HashPassword("dummy-password-for-timing")
	if err != nil {
		panic(err)
	}
	return hash
}