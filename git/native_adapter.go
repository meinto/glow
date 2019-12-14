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
func (a nativeGitAdapter) SetCICDOrigin(origin string) (stdout, stderr bytes.Buffer, err error) {
	cmd := a.exec.Command(fmt.Sprintf("git remote set-url origin %s", origin))
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	return stdout, stderr, err
}

// GitRepoPath returns the path to the root with the .git folder
func (a nativeGitAdapter) GitRepoPath() (repoPath string, stdout, stderr bytes.Buffer, err error) {
	cmd := a.exec.Command("git rev-parse --show-toplevel")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	return strings.TrimSuffix(stdout.String(), "\n"), stdout, stderr, err
}

// CurrentBranch returns the current branch name
func (a nativeGitAdapter) CurrentBranch() (branch glow.Branch, stdout, stderr bytes.Buffer, err error) {
	cmdBranchList, stdout, stderr, err := getCMDBranchList(a.exec)
	if err != nil {
		return nil, stdout, stderr, err
	}
	for _, b := range cmdBranchList {
		if b.IsCurrentBranch {
			branch, err := glow.NewBranch(b.Name)
			return branch, stdout, stderr, err
		}
	}
	return nil, stdout, stderr, errors.New("cannot detect current branch")
}

// BranchList returns a list of avalilable branches
func (a nativeGitAdapter) BranchList() (branchList []glow.Branch, stdout, stderr bytes.Buffer, err error) {
	cmdBranchList, stdout, stderr, err := getCMDBranchList(a.exec)
	if err != nil {
		return nil, stdout, stderr, err
	}
	branchList = make([]glow.Branch, 0)
	for _, b := range cmdBranchList {
		gb, err := glow.NewBranch(b.Name)
		if err != nil {
			return nil, stdout, stderr, err
		}
		branchList = append(branchList, gb)
	}
	return branchList, stdout, stderr, nil
}

// Fetch changes
func (a nativeGitAdapter) Fetch() (stdout, stderr bytes.Buffer, err error) {
	cmd := a.exec.Command("git fetch")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	return stdout, stderr, err
}

// Add all changes
func (a nativeGitAdapter) AddAll() (stdout, stderr bytes.Buffer, err error) {
	cmd := a.exec.Command("git add -A")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	return stdout, stderr, err
}

// Stash all changes
func (a nativeGitAdapter) Stash() (stdout, stderr bytes.Buffer, err error) {
	cmd := a.exec.Command("git stash")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	return bytes.Buffer{}, bytes.Buffer{}, err
}

// Pop all stashed changes
func (a nativeGitAdapter) StashPop() (stdout, stderr bytes.Buffer, err error) {
	cmd := a.exec.Command("git stash pop")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	return stdout, stderr, err
}

// Commit added changes
func (a nativeGitAdapter) Commit(message string) (stdout, stderr bytes.Buffer, err error) {
	cmd := a.exec.Command(fmt.Sprintf("git commit -m '%s'", message))
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	return stdout, stderr, err
}

// Push changes
func (a nativeGitAdapter) Push(setUpstream bool) (stdout, stderr bytes.Buffer, err error) {
	cmd := a.exec.Command("git push")
	if setUpstream {
		currentBranch, stdout, stderr, err := a.CurrentBranch()
		if err != nil {
			return stdout, stderr, err
		}
		cmd = a.exec.Command(fmt.Sprintf("git push -u origin %s", currentBranch.ShortBranchName()))
	}
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	return stdout, stderr, err
}

// Create a new branch
func (a nativeGitAdapter) Create(b glow.Branch, skipChecks bool) (stdout, stderr bytes.Buffer, err error) {
	if !skipChecks {
		sourceBranch, stdout, stderr, err := a.CurrentBranch()
		if err != nil {
			return stdout, stderr, err
		}
		if !b.CreationIsAllowedFrom(sourceBranch.BranchName()) {
			return stdout, stderr, errors.New("creation not allowed from this branch")
		}
	}
	cmd := a.exec.Command(fmt.Sprintf("git branch %s", b.ShortBranchName()))
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	return stdout, stderr, err
}

// Checkout a branch
func (a nativeGitAdapter) Checkout(b glow.Branch) (stdout, stderr bytes.Buffer, err error) {
	cmd := a.exec.Command(fmt.Sprintf("git checkout %s", b.ShortBranchName()))
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	return stdout, stderr, err
}

// CleanupBranches removes all unused branches
func (a nativeGitAdapter) CleanupBranches(cleanupGone, cleanupUntracked bool) (stdout, stderr bytes.Buffer, err error) {
	xargsCmd := "xargs -r git branch -D"
	if runtime.GOOS == "darwin" {
		xargsCmd = "xargs git branch -D"
	}

	if cleanupGone {
		cmd := a.exec.Command("git remote prune origin")
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			return stdout, stderr, err
		}

		cmd = a.exec.Command(fmt.Sprintf("git branch -vv | grep 'origin/.*: gone]' | awk '{print $1}' | %s", xargsCmd))
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		err = cmd.Run()
		if err != nil {
			return stdout, stderr, err
		}
	}

	if cleanupUntracked {
		cmd := a.exec.Command(fmt.Sprintf("git branch -vv | cut -c 3- | grep -v detached | awk '$3 !~/\\[origin/ { print $1 }' | %s", xargsCmd))
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			return stdout, stderr, err
		}
	}
	return stdout, stderr, err
}

// CleanupTags removes tags from local repo
func (a nativeGitAdapter) CleanupTags(cleanupUntracked bool) (stdout, stderr bytes.Buffer, err error) {
	xargsCmd := "xargs -r git tag -d"
	if runtime.GOOS == "darwin" {
		xargsCmd = "xargs git tag -d"
	}

	if cleanupUntracked {
		cmd := a.exec.Command(fmt.Sprintf("git tag -l | %s", xargsCmd))
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		err := cmd.Run()

		if err != nil {
			return stdout, stderr, err
		}

		cmd = a.exec.Command("git fetch --tags")
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		err = cmd.Run()

		if err != nil {
			return stdout, stderr, err
		}
	}
	return stdout, stderr, err
}

func (a nativeGitAdapter) RemoteBranchExists(branchName string) (stdout, stderr bytes.Buffer, err error) {
	cmd := a.exec.Command(fmt.Sprintf("git ls-remote --heads $(git remote get-url origin) %s | wc -l", branchName))
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()

	if err != nil {
		return stdout, stderr, err
	}

	branchCount := strings.TrimSpace(stdout.String())
	if branchCount == "1" {
		return stdout, stderr, err
	}
	return stdout, stderr, err
}

type cmdBranch struct {
	Name            string
	IsCurrentBranch bool
}

func getCMDBranchList(exec cmd.CmdExecutor) (branch []cmdBranch, stdout, stderr bytes.Buffer, err error) {
	cmd := exec.Command("git branch --list")
	cmd.Stderr = &stderr
	stdoutReader, err := cmd.StdoutPipe()
	if err != nil {
		return []cmdBranch{}, bytes.Buffer{}, stderr, err
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
	return branchList, bytes.Buffer{}, bytes.Buffer{}, err
}
