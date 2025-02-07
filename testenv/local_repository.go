package testenv

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type LocalRepository struct {
	*Repository
}

func NewLocalRepository(folder string) *LocalRepository {
	return &LocalRepository{
		NewRepository(folder + "/local"),
	}
}

func Clone(pathToBare, repoName string) *LocalRepository {
	exec := NewCommand()
	execDir := strings.TrimSuffix(pathToBare, "/bare.git")
	exec.Do(fmt.Sprintf("cd %s && git clone %s %s", execDir, pathToBare, repoName))
	return &LocalRepository{
		NewRepository(execDir + "/" + repoName),
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

	mainBranch := viper.GetString("mainBranch")
	r.Do(fmt.Sprintf("git push -u origin %s", mainBranch))
}

func (r *LocalRepository) CreateBranch(branch string) {
	r.Do(fmt.Sprintf("git branch %s", branch))
}

func (r *LocalRepository) Checkout(branch string) {
	r.Do(fmt.Sprintf("git checkout %s", branch))
}

func (r *LocalRepository) Push(branch string) {
	r.Do(fmt.Sprintf("git push -u origin %s", branch))
}

func (r *LocalRepository) Exists(branch string) (bool, string) {
	stdout, _, _ := r.Do(fmt.Sprintf("git rev-parse --abbrev-ref %s", branch))
	return strings.TrimSpace(stdout.String()) == branch, strings.TrimSpace(stdout.String())
}
