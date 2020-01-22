package glow

import (
	l "github.com/meinto/glow/logging"
)

type branchLoggingService struct {
	next Branch
}

func NewBranchLoggingService(b Branch) Branch {
	defer func() {
		l.Log().Info(l.Fields{"branch": b})
	}()
	return &branchLoggingService{b}
}

func (s *branchLoggingService) CreationIsAllowedFrom(sourceBranch Branch) (creationAllowed bool) {
	defer func() {
		l.Log().Info(l.Fields{
			"sourceBranch":    sourceBranch,
			"creationAllowed": creationAllowed,
		})
	}()
	return s.next.CreationIsAllowedFrom(sourceBranch)
}

func (s *branchLoggingService) CanBeClosed() (canBeClosed bool) {
	defer func() {
		l.Log().Info(l.Fields{"canBeClosed": canBeClosed})
	}()
	return s.next.CanBeClosed()
}

func (s *branchLoggingService) CanBePublished() (canBePublished bool) {
	defer func() {
		l.Log().Info(l.Fields{"canBePublished": canBePublished})
	}()
	return s.next.CanBePublished()
}

func (s *branchLoggingService) CloseBranches(availableBranches []Branch) (closeBranches []Branch) {
	defer func() {
		l.Log().Info(l.Fields{
			"availableBranches": availableBranches,
			"closeBranches":     closeBranches,
		})
	}()
	return s.next.CloseBranches(availableBranches)
}

func (s *branchLoggingService) PublishBranch() (publishBranch Branch) {
	defer func() {
		l.Log().Info(l.Fields{"publishBranch": publishBranch})
	}()
	return s.next.PublishBranch()
}

func (s *branchLoggingService) BranchName() (branchName string) {
	defer func() {
		l.Log().Info(l.Fields{"branchName": branchName})
	}()
	return s.next.BranchName()
}

func (s *branchLoggingService) ShortBranchName() (shortBranchName string) {
	defer func() {
		l.Log().Info(l.Fields{"shortBranchName": shortBranchName})
	}()
	return s.next.ShortBranchName()
}
