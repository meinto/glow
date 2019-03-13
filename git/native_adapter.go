package git

import (
	"bufio"
	"bytes"
	"os"
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
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	return errors.Wrap(err, stderr.String())
}

// Create a new branch
func (a nativeGitAdapter) Create(b glow.Branch) error {
	cmd := exec.Command(a.gitPath, "branch", b.ShortBranchName())
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	return errors.Wrap(err, stderr.String())
}

// Checkout a branch
func (a nativeGitAdapter) Checkout(b glow.Branch) error {
	cmd := exec.Command(a.gitPath, "checkout", b.ShortBranchName())
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	return errors.Wrap(err, stderr.String())
}

// CleanupBranches removes all unused branches
func (a nativeGitAdapter) CleanupBranches(cleanupGone, cleanupUntracked bool) error {
	if cleanupGone {
		cmd := exec.Command(a.gitPath, "remote", "prune", "origin")
		err := cmd.Run()
		if err != nil {
			return errors.Wrap(err, "error pruning branches")
		}
		// /usr/bin/git branch -vv | /usr/bin/grep 'origin/.*: gone]' | /usr/bin/awk '{print $1}' | /usr/bin/xargs /usr/bin/git branch -D
		c1 := exec.Command(a.gitPath, "branch", "-vv")
		c2 := exec.Command("/usr/bin/grep", "origin/.*: gone]")
		c3 := exec.Command("/usr/bin/awk", "{print $1}")
		c4 := exec.Command("/usr/bin/xargs", a.gitPath, "branch", "-D")

		r1, w1, err := os.Pipe()
		c1.Stdout = w1
		c2.Stdin = r1
		if err != nil {
			return err
		}

		r2, w2, err := os.Pipe()
		c2.Stdout = w2
		c3.Stdin = r2
		if err != nil {
			return err
		}

		r3, w3, err := os.Pipe()
		c3.Stdout = w3
		c4.Stdin = r3
		if err != nil {
			return err
		}

		var b1, b2, b3, b4 bytes.Buffer
		c1.Stderr = &b1
		c2.Stderr = &b2
		c3.Stderr = &b3
		c4.Stderr = &b4

		c1.Start()

		c2.Start()
		c1.Wait()
		w1.Close()

		c3.Start()
		c2.Wait()
		w2.Close()

		c4.Start()
		c3.Wait()
		w3.Close()

		c4.Wait()

		errorString := b1.String() + b2.String() + b3.String() + b4.String()
		if errorString != "" {
			return errors.New(errorString)
		}
	}

	if cleanupUntracked {
		c1 := exec.Command(a.gitPath, "branch", "-vv")
		c2 := exec.Command("/usr/bin/cut", "-c", "3-")
		c3 := exec.Command("/usr/bin/grep", "-v", "detached")
		c4 := exec.Command("/usr/bin/awk", "$3 !~/\\[/ { print $1 }")
		c5 := exec.Command("/usr/bin/xargs", a.gitPath, "branch", "-D")

		r1, w1, err := os.Pipe()
		c1.Stdout = w1
		c2.Stdin = r1
		if err != nil {
			return err
		}

		r2, w2, err := os.Pipe()
		c2.Stdout = w2
		c3.Stdin = r2
		if err != nil {
			return err
		}

		r3, w3, err := os.Pipe()
		c3.Stdout = w3
		c4.Stdin = r3
		if err != nil {
			return err
		}

		r4, w4, err := os.Pipe()
		c4.Stdout = w4
		c5.Stdin = r4
		if err != nil {
			return err
		}

		var b1, b2, b3, b4, b5 bytes.Buffer
		c1.Stderr = &b1
		c2.Stderr = &b2
		c3.Stderr = &b3
		c4.Stderr = &b4
		c5.Stderr = &b5

		c1.Start()

		c2.Start()
		c1.Wait()
		w1.Close()

		c3.Start()
		c2.Wait()
		w2.Close()

		c4.Start()
		c3.Wait()
		w3.Close()

		c5.Start()
		c4.Wait()
		w4.Close()

		c5.Wait()

		errorString := b1.String() + b2.String() + b3.String() + b4.String() + b5.String()
		if errorString != "" {
			return errors.New(errorString)
		}
	}
	return nil
}

// CleanupTags removes tags from local repo
func (a nativeGitAdapter) CleanupTags(cleanupUntracked bool) error {
	if cleanupUntracked {
		c1 := exec.Command(a.gitPath, "tag", "-l")
		c2 := exec.Command("/usr/bin/xargs", a.gitPath, "tag", "-d")
		c3 := exec.Command(a.gitPath, "fetch", "--tags")

		r1, w1, err := os.Pipe()
		c1.Stdout = w1
		c2.Stdin = r1
		if err != nil {
			return err
		}

		var b1, b2, b3 bytes.Buffer
		c1.Stderr = &b1
		c2.Stderr = &b2
		c3.Stderr = &b3

		c1.Start()

		c2.Start()
		c1.Wait()
		w1.Close()
		c2.Wait()

		c3.Start()
		c3.Wait()

		errorString := b1.String() + b2.String() + b3.String()
		if errorString != "" {
			return errors.New(errorString)
		}
	}
	return nil
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
