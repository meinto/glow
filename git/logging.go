package git

import (
	"bytes"

	"github.com/meinto/glow"
	l "github.com/meinto/glow/logging"
	"github.com/sirupsen/logrus"
)

type loggingService struct {
	next Service
}

// NewLoggingService returns a new instance of a logging Service.
func NewLoggingService(s Service) Service {
	return &loggingService{s}
}

// SetCICDOrigin for pipeline
func (s loggingService) SetCICDOrigin(origin string) (err error) {
	defer l.Log().WithFields(logrus.Fields{
		"origin": origin,
		"error":  err,
	}).Info()
	return s.next.SetCICDOrigin(origin)
}

// GitRepoPath returns the path to the root with the .git folder
func (s loggingService) GitRepoPath() (_ string, err error) {
	defer l.Log().WithFields(logrus.Fields{
		"error": err,
	}).Info()
	return s.next.GitRepoPath()
}

// CurrentBranch returns the current branch name
func (s loggingService) CurrentBranch() (_ glow.Branch, err error) {
	defer l.Log().WithFields(logrus.Fields{
		"error": err,
	}).Info()
	return s.next.CurrentBranch()
}

// BranchList returns a list of avalilable branches
func (s loggingService) BranchList() (_ []glow.Branch, err error) {
	defer l.Log().WithFields(logrus.Fields{
		"error": err,
	}).Info()
	return s.next.BranchList()
}

// Fetch changes
func (s loggingService) Fetch() (err error) {
	defer l.Log().WithFields(logrus.Fields{
		"error": err,
	}).Info()
	return s.next.Fetch()
}

// Add all changes
func (s loggingService) AddAll() (err error) {
	defer l.Log().WithFields(logrus.Fields{
		"error": err,
	}).Info()
	return s.next.AddAll()
}

// Stash all changes
func (s loggingService) Stash() (err error) {
	defer l.Log().WithFields(logrus.Fields{
		"error": err,
	}).Info()
	return s.next.Stash()
}

// Pop all stashed changes
func (s loggingService) StashPop() (stdout, stderr bytes.Buffer, err error) {
	defer l.Log().WithFields(logrus.Fields{
		"error": err,
	}).Info()
	return s.next.StashPop()
}

// Commit added changes
func (s loggingService) Commit(message string) (err error) {
	defer l.Log().WithFields(logrus.Fields{
		"error": err,
	}).Info()
	return s.next.Commit(message)
}

// Push changes
func (s loggingService) Push(setUpstream bool) (stdout, stderr bytes.Buffer, err error) {
	defer l.Log().WithFields(logrus.Fields{
		"stdout": stdout.String(),
		"stderr": stderr.String(),
		"error":  err,
	}).Info()
	return s.next.Push(setUpstream)
}

// Create a new branch
func (s loggingService) Create(b glow.Branch, skipChecks bool) (err error) {
	defer l.Log().WithFields(logrus.Fields{
		"error": err,
	}).Info()
	return s.next.Create(b, skipChecks)
}

// Checkout a branch
func (s loggingService) Checkout(b glow.Branch) (err error) {
	defer l.Log().WithFields(logrus.Fields{
		"error": err,
	}).Info()
	return s.next.Checkout(b)
}

// CleanupBranches removes all unused branches
func (s loggingService) CleanupBranches(cleanupGone, cleanupUntracked bool) (err error) {
	defer l.Log().WithFields(logrus.Fields{
		"error": err,
	}).Info()
	return s.next.CleanupBranches(cleanupGone, cleanupUntracked)
}

// CleanupTags removes tags from local repo
func (s loggingService) CleanupTags(cleanupUntracked bool) (err error) {
	defer l.Log().WithFields(logrus.Fields{
		"error": err,
	}).Info()
	return s.next.CleanupTags(cleanupUntracked)
}

func (s loggingService) RemoteBranchExists(branchName string) (err error) {
	defer l.Log().WithFields(logrus.Fields{
		"error": err,
	}).Info()
	return s.next.RemoteBranchExists(branchName)
}
