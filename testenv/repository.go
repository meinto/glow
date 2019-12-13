package testenv

import (
	"bytes"
	"fmt"
)

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

func (r *Repository) Do(str string, args ...interface{}) (stdout, stderr bytes.Buffer, err error) {
	moveToDir := fmt.Sprintf("cd %s", r.Folder)
	return r.exec.Do(moveToDir+" && "+str, args...)
}
