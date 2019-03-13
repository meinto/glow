package git

import (
	"bufio"
	"bytes"
	"os/exec"
	"strings"

	"github.com/meinto/glow"
	"github.com/pkg/errors"
)

// NativeGitAdapter implemented with native git
type nativeGitAdapter struct {
	gitPath string
}

// GitRepoPath returns the path to the root with the .git folder
func (a nativeGitAdapter) GitRepoPath() (string, error) {
	cmd := exec.Command(a.gitPath, "rev-parse", "--show-toplevel")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	err := cmd.Run()
	return strings.TrimSuffix(stdout.String(), "\n"), err
}

// CurrentBranch returns the current branch name
func (a nativeGitAdapter) CurrentBranch() (glow.Branch, error) {
	cmdBranchList, err := getCMDBranchList(a.gitPath)
	if err != nil {
		return nil, err
	}
	for _, b := range cmdBranchList {
		if b.IsCurrentBranch {
			return glow.NewBranch(b.Name)
		}
	}
	return nil, errors.New("cannot detect current branch")
}

// BranchList returns a list of avalilable branches
func (a nativeGitAdapter) BranchList() ([]glow.Branch, error) {
	cmdBranchList, err := getCMDBranchList(a.gitPath)
	if err != nil {
		return nil, err
	}
	branchList := make([]glow.Branch, 0)
	for _, b := range cmdBranchList {
		gb, err := glow.NewBranch(b.Name)
		if err != nil {
			return nil, err
		}
		branchList = append(branchList, gb)
	}
	return branchList, nil
}

// Fetch changes
func (a nativeGitAdapter) Fetch() error {
	cmd := exec.Command(a.gitPath, "fetch")
	err := cmd.Run()
	return errors.Wrap(err, "native Fetch")
}

// Create a new branch
func (a nativeGitAdapter) Create(b glow.Branch) error {
	cmd := exec.Command(a.gitPath, "branch", b.ShortBranchName())
	err := cmd.Run()
	return errors.Wrap(err, "native Create")
}

// Checkout a branch
func (a nativeGitAdapter) Checkout(b glow.Branch) error {
	cmd := exec.Command(a.gitPath, "checkout", b.ShortBranchName())
	err := cmd.Run()
	return errors.Wrap(err, "native Checkout")
}

// CleanupBranches removes all unused branches
func (a nativeGitAdapter) CleanupBranches() error {
	cmd := exec.Command(a.gitPath, "remote", "prune", "origin")
	err := cmd.Run()
	if err != nil {
		return errors.Wrap(err, "error pruning branches")
	}
	// git branch -vv | grep 'origin/.*: gone]' | awk '{print $1}' | xargs git branch -D
	args := []string{
		"-vv",
		"|", "grep", "'origin/.*: gone]'",
		"|", "awk", "'{print $1}'",
		"|", "xargs", "git", "branch", "-D",
	}
	cmd = exec.Command(a.gitPath, args...)
	err = cmd.Run()
	return errors.Wrap(err, "error cleanup branches")
}

type cmdBranch struct {
	Name            string
	IsCurrentBranch bool
}

func getCMDBranchList(gitPath string) ([]cmdBranch, error) {
	cmd := exec.Command(gitPath, "branch", "--list")
	stdoutReader, err := cmd.StdoutPipe()
	if err != nil {
		return []cmdBranch{}, err
	}
	if err != nil {
		return []cmdBranch{}, err
	}

	c := make(chan []cmdBranch)
	go func(c chan []cmdBranch) {
		var branchList []cmdBranch
		scanner := bufio.NewScanner(stdoutReader)
		for scanner.Scan() {
			line := strings.Trim(scanner.Text(), " ")
			parts := strings.Split(line, " ")

			name := parts[0]
			isCurrentBranch := false
			if len(parts) > 1 {
				name = parts[1]
				isCurrentBranch = true
			}

			branchList = append(branchList, cmdBranch{
				name,
				isCurrentBranch,
			})
		}
		c <- branchList
		close(c)
	}(c)
	err = cmd.Run()
	branchList := <-c
	return branchList, err
}
