package git

import (
	"bufio"
	"bytes"
	"fmt"
	"runtime"
	"strings"

	"github.com/meinto/glow/cmd"

	"github.com/meinto/glow"
	"github.com/pkg/errors"
)

// NativeGitAdapter implemented with native git
type nativeGitAdapter struct {
	exec cmd.CmdExecutor
}

// SetCICDOrigin for pipeline
func (a nativeGitAdapter) SetCICDOrigin(origin string) error {
	cmd := a.exec.Command(fmt.Sprintf("git config remote.origin.url '%s'", origin))
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	return errors.Wrap(err, stderr.String())
}

// GitRepoPath returns the path to the root with the .git folder
func (a nativeGitAdapter) GitRepoPath() (string, error) {
	cmd := a.exec.Command("git rev-parse --show-toplevel")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	err := cmd.Run()
	return strings.TrimSuffix(stdout.String(), "\n"), err
}

// CurrentBranch returns the current branch name
func (a nativeGitAdapter) CurrentBranch() (glow.Branch, error) {
	cmdBranchList, err := getCMDBranchList(a.exec)
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
	cmdBranchList, err := getCMDBranchList(a.exec)
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
	cmd := a.exec.Command("git fetch")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	return errors.Wrap(err, stderr.String())
}

// Add all changes
func (a nativeGitAdapter) AddAll() error {
	cmd := a.exec.Command("git add .")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	return errors.Wrap(err, stderr.String())
}

// Commit added changes
func (a nativeGitAdapter) Commit(message string) error {
	cmd := a.exec.Command(fmt.Sprintf("git commit -m '%s'", message))
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	return errors.Wrap(err, stderr.String())
}

// Push changes
func (a nativeGitAdapter) Push(setUpstream bool) error {
	cmd := a.exec.Command("git push")
	if setUpstream {
		currentBranch, err := a.CurrentBranch()
		if err != nil {
			return errors.Wrap(err, "error while getting current branch")
		}
		cmd = a.exec.Command(fmt.Sprintf("git push -u origin %s", currentBranch.ShortBranchName()))
	}
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	return errors.Wrap(err, stderr.String())
}

// Create a new branch
func (a nativeGitAdapter) Create(b glow.Branch, skipChecks bool) error {
	if !skipChecks {
		sourceBranch, err := a.CurrentBranch()
		if err != nil {
			return err
		}
		if !b.CreationIsAllowedFrom(sourceBranch.BranchName()) {
			return errors.New("creation not allowed from this branch")
		}
	}
	cmd := a.exec.Command(fmt.Sprintf("git branch %s", b.ShortBranchName()))
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	return errors.Wrap(err, stderr.String())
}

// Checkout a branch
func (a nativeGitAdapter) Checkout(b glow.Branch) error {
	cmd := a.exec.Command(fmt.Sprintf("git checkout %s", b.ShortBranchName()))
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	return errors.Wrap(err, stderr.String())
}

// CleanupBranches removes all unused branches
func (a nativeGitAdapter) CleanupBranches(cleanupGone, cleanupUntracked bool) error {
	xargsCmd := "xargs -r git branch -D"
	if runtime.GOOS == "darwin" {
		xargsCmd = "xargs git branch -D"
	}

	if cleanupGone {
		cmd := a.exec.Command("git remote prune origin")
		err := cmd.Run()
		if err != nil {
			return errors.Wrap(err, "error pruning branches")
		}

		cmd = a.exec.Command(fmt.Sprintf("git branch -vv | grep 'origin/.*: gone]' | awk '{print $1}' | %s", xargsCmd))
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		err = cmd.Run()
		if err != nil {
			return errors.Wrap(err, stderr.String())
		}
	}

	if cleanupUntracked {
		cmd := a.exec.Command(fmt.Sprintf("git branch -vv | cut -c 3- | grep -v detached | awk '$3 !~/\\[origin/ { print $1 }' | %s", xargsCmd))
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
	xargsCmd := "xargs -r git tag -d"
	if runtime.GOOS == "darwin" {
		xargsCmd = "xargs git tag -d"
	}

	if cleanupUntracked {
		cmd := a.exec.Command(fmt.Sprintf("git tag -l | %s", xargsCmd))
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		err := cmd.Run()

		if err != nil {
			return errors.Wrap(err, stderr.String())
		}

		cmd = a.exec.Command("git fetch --tags")
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

func getCMDBranchList(exec cmd.CmdExecutor) ([]cmdBranch, error) {
	cmd := exec.Command("git branch --list")
	stdoutReader, err := cmd.StdoutPipe()
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
