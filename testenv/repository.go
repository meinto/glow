package testenv

import "fmt"

type Repository struct {
	Folder string
	exec   *Command
}

func NewRepository(folder string) *Repository {
	return &Repository{
		folder,
		NewCommand(),
	}
}

func (r *Repository) Do(str string, args ...interface{}) {
	moveToDir := fmt.Sprintf("cd %s", r.Folder)
	r.exec.Do(moveToDir+" && "+str, args...)
}
