package git

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/meinto/glow"
	"github.com/pkg/errors"
)

// NativeGitAdapter implemented with native git
type nativeGitAdapter struct {
	shell string
}

// SetCICDOrigin for pipeline
func (a nativeGitAdapter) SetCICDOrigin(origin string) error {
	cmd := exec.Command(a.shell, "-c", fmt.Sprintf("git config remote.origin.url %s", origin))
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	return errors.Wrap(err, stderr.String())
}

// GitRepoPath returns the path to the root with the .git folder
func (a nativeGitAdapter) GitRepoPath() (string, error) {
	cmd := exec.Command(a.shell, "-c", "git rev-parse --show-toplevel")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	err := cmd.Run()
	return strings.TrimSuffix(stdout.String(), "\n"), err
}

// CurrentBranch returns the current branch name
func (a nativeGitAdapter) CurrentBranch() (glow.Branch, error) {
	cmdBranchList, err := getCMDBranchList(a.shell)
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
	cmdBranchList, err := getCMDBranchList(a.shell)
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
	cmd := exec.Command(a.shell, "-c", "git fetch")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	return errors.Wrap(err, stderr.String())
}

// Add all changes
func (a nativeGitAdapter) AddAll() error {
	cmd := exec.Command(a.shell, "-c", "git add .")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	return errors.Wrap(err, stderr.String())
}

// Commit added changes
func (a nativeGitAdapter) Commit(message string) error {
	cmd := exec.Command(a.shell, "-c", fmt.Sprintf("git commit -m '%s'", message))
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	return errors.Wrap(err, stderr.String())
}

// Push changes
func (a nativeGitAdapter) Push(setUpstream bool) error {
	cmd := exec.Command(a.shell, "-c", "git push")
	if setUpstream {
		currentBranch, err := a.CurrentBranch()
		if err != nil {
			return errors.Wrap(err, "error while getting current branch")
		}
		cmd = exec.Command(a.shell, "-c", fmt.Sprintf("git push -u origin %s", currentBranch.ShortBranchName()))
	}
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	return errors.Wrap(err, stderr.String())
}

// Create a new branch
func (a nativeGitAdapter) Create(b glow.Branch) error {
	cmd := exec.Command(a.shell, "-c", fmt.Sprintf("git branch %s", b.ShortBranchName()))
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	return errors.Wrap(err, stderr.String())
}

// Checkout a branch
func (a nativeGitAdapter) Checkout(b glow.Branch) error {
	cmd := exec.Command(a.shell, "-c", fmt.Sprintf("git checkout %s", b.ShortBranchName()))
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	return errors.Wrap(err, stderr.String())
}

// CleanupBranches removes all unused branches
func (a nativeGitAdapter) CleanupBranches(cleanupGone, cleanupUntracked bool) error {
	if cleanupGone {
		cmd := exec.Command(a.shell, "-c", "git remote prune origin")
		err := cmd.Run()
		if err != nil {
			return errors.Wrap(err, "error pruning branches")
		}

		cmd = exec.Command(a.shell, "-c", "git branch -vv | grep 'origin/.*: gone]' | awk '{print $1}' | xargs --no-run-if-empty git branch -D")
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		err = cmd.Run()
		if err != nil {
			return errors.Wrap(err, stderr.String())
		}
	}

	if cleanupUntracked {
		cmd := exec.Command(a.shell, "-c", "git branch -vv | cut -c 3- | grep -v detached | awk '$3 !~/\\[/ { print $1 }' | xargs --no-run-if-empty git branch -D")
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			return errors.Wrap(err, stderr.String())
		}
	}
	return nil
}

// CleanupTags removes tags from local repo
func (a nativeGitAdapter) CleanupTags(cleanupUntracked bool) error {
	if cleanupUntracked {
		cmd := exec.Command(a.shell, "-c", "git tag -l | xargs --no-run-if-empty git tag -d")
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		err := cmd.Run()

		if err != nil {
			return errors.Wrap(err, stderr.String())
		}

		cmd = exec.Command(a.shell, "-c", "git fetch --tags")
		cmd.Stderr = &stderr
		err = cmd.Run()

		if err != nil {
			return errors.Wrap(err, stderr.String())
		}
	}
	return nil
}

type cmdBranch struct {
	Name            string
	IsCurrentBranch bool
}

func getCMDBranchList(shell string) ([]cmdBranch, error) {
	cmd := exec.Command(shell, "git branch --list")
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
