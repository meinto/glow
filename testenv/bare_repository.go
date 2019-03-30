package testenv

import "os"

type BareRepository struct {
	*Repository
}

func NewBareRepository(folder string) *BareRepository {
	return &BareRepository{
		NewRepository(folder + "/bare.git"),
	}
}

func (r *BareRepository) Init() {
	os.MkdirAll(r.Folder, os.ModePerm)
	r.Do("git --bare init")
}
