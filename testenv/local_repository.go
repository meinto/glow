package testenv

import (
	"fmt"
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
	r.Do("mkdir subfolder")
	r.Do("touch subfolder/somefile")
	r.Do("git add .")
	r.Do("git commit -m 'initial commit'")
	r.Do("git push -u origin master")
}

func (r *LocalRepository) CreateBranch(branch string) {
	r.Do(fmt.Sprintf("git branch %s", branch))
}

func (r *LocalRepository) Checkout(branch string) {
	r.Do(fmt.Sprintf("git checkout %s", branch))
}
