package testenv

import (
	"os"
)

type LocalRepository struct {
	*Repository
}

func NewLocalRepository(folder string) *LocalRepository {
	return &LocalRepository{
		NewRepository(folder + "/local"),
	}
}

func (r *LocalRepository) Init(pathToBare string) {
	os.MkdirAll(r.Folder, os.ModePerm)
	r.Do("git init")
	r.Do("git remote add origin %s", pathToBare)
	r.Do("echo '1.2.3' > VERSION")
	r.Do("git add .")
	r.Do("git commit -m 'initial commit'")
	r.Do("git push -u origin master")
}
