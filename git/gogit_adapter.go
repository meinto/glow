package git

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/meinto/glow"
	"github.com/pkg/errors"
	"gopkg.in/src-d/go-billy.v4/helper/chroot"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/storage/filesystem"
)

// GoGitAdapter implemented with go-git
type goGitAdapter struct{}

// GitRepoPath returns the path to the root with the .git folder
func (a goGitAdapter) GitRepoPath() (string, error) {
	r, err := git.PlainOpenWithOptions(".", &git.PlainOpenOptions{
		DetectDotGit: true,
	})
	if err != nil {
		return "", errors.Wrap(err, "error opening repository")
	}

	// Try to grab the repository Storer
	s, ok := r.Storer.(*filesystem.Storage)
	if !ok {
		return "", errors.New("Repository storage is not filesystem.Storage")
	}

	// Try to get the underlying billy.Filesystem
	fs, ok := s.Filesystem().(*chroot.ChrootHelper)
	if !ok {
		return "", errors.New("Filesystem is not chroot.ChrootHelper")
	}

	p := fs.Root()

	// Clean up the path
	if !filepath.IsAbs(p) {
		pwd, err := os.Getwd()
		if err != nil {
			return "", err
		}
		p = filepath.Join(pwd, p)
	}
	return strings.TrimSuffix(p, "/.git"), nil
}

// CurrentBranch returns the current branch name
func (a goGitAdapter) CurrentBranch() (glow.Branch, error) {
	r, err := git.PlainOpenWithOptions(".", &git.PlainOpenOptions{
		DetectDotGit: true,
	})
	if err != nil {
		return nil, errors.Wrap(err, "error opening repository")
	}

	headRef, err := r.Head()
	if err != nil {
		return nil, errors.Wrap(err, "error getting current branch")
	}

	refName := headRef.Name().String()
	return glow.NewBranch(refName)
}

// BranchList returns a list of avalilable branches
func (s goGitAdapter) BranchList() ([]glow.Branch, error) {
	r, err := git.PlainOpenWithOptions(".", &git.PlainOpenOptions{
		DetectDotGit: true,
	})
	if err != nil {
		return nil, errors.Wrap(err, "error opening repository")
	}

	refList, err := r.References()
	if err != nil {
		return nil, errors.Wrap(err, "error getting ref list")
	}

	branches := make([]glow.Branch, 0)
	refPrefix := "refs/heads/"
	refList.ForEach(func(ref *plumbing.Reference) error {
		refName := ref.Name().String()
		if strings.HasPrefix(refName, refPrefix) {
			b, _ := glow.NewBranch(refName)
			branches = append(branches, b)
		}
		return nil
	})

	return branches, nil
}

// Fetch changes
func (a goGitAdapter) Fetch() error {
	r, err := git.PlainOpenWithOptions(".", &git.PlainOpenOptions{
		DetectDotGit: true,
	})
	if err != nil {
		return errors.Wrap(err, "error while fetching")
	}
	return r.Fetch(&git.FetchOptions{})
}

// Create a new branch
func (a goGitAdapter) Create(b glow.Branch) error {
	r, err := git.PlainOpenWithOptions(".", &git.PlainOpenOptions{
		DetectDotGit: true,
	})
	if err != nil {
		return errors.Wrap(err, "error opening repository")
	}

	headRef, err := r.Head()
	if err != nil {
		return errors.Wrap(err, "error getting current branch")
	}

	refName := string(headRef.Name())
	if !b.CreationIsAllowedFrom(refName) {
		return errors.New("You are not on the develop branch.\nPlease switch branch...\n")
	}

	ref := plumbing.NewHashReference(plumbing.ReferenceName(b.BranchName()), headRef.Hash())

	err = r.Storer.SetReference(ref)
	return errors.Wrap(err, "error while creating branch")
}

// Checkout a branch
func (a goGitAdapter) Checkout(b glow.Branch) error {
	r, err := git.PlainOpenWithOptions(".", &git.PlainOpenOptions{
		DetectDotGit: true,
	})
	if err != nil {
		return errors.Wrap(err, "error opening repository")
	}

	w, err := r.Worktree()
	if err != nil {
		return errors.Wrap(err, "error getting worktree")
	}

	err = w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName(b.BranchName()),
	})
	return errors.Wrapf(err, "error while checkout branch %s", b.BranchName())
}

// CleanupBranches removes all unused branches
func (a goGitAdapter) CleanupBranches(cleanupGone, cleanupUntracked bool) error {
	return errors.New("not implemented yet")
}
