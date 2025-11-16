package entities

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type ErrDuplicateUserIDs struct {
	IDs []uuid.UUID
}

func (e *ErrDuplicateUserIDs) Error() string {
	return fmt.Sprintf("duplicate user IDs in request: %v", e.IDs)
}

type ErrTeamAlreadyExists struct {
	Name string
}

func (e *ErrTeamAlreadyExists) Error() string {
	return fmt.Sprintf("team already exists: %s", e.Name)
}

type ErrTeamNotFound struct {
	Name string
}

type ErrUserNameValidation struct {
	Reason string
}

func (e *ErrUserNameValidation) Error() string {
	return fmt.Sprintf("user name is invalid: %s", e.Reason)
}

func (e *ErrTeamNotFound) Error() string {
	return fmt.Sprintf("team not found: %s", e.Name)
}

type ErrUserNotFound struct {
	UserID *uuid.UUID
	Name   *string
}

func (e *ErrUserNotFound) Error() string {
	if e.Name != nil {
		return fmt.Sprintf("user not found: %s", *e.Name)
	} else if e.UserID != nil {
		return fmt.Sprintf("user not found: %s", *e.UserID)
	}
	return "user not found"
}

type ErrTeamNameValidation struct {
	Reason string
}

func (e *ErrTeamNameValidation) Error() string {
	return fmt.Sprintf("team name is invalid: %s", e.Reason)
}

type ErrPullRequestAlreadyExists struct {
	ID uuid.UUID
}

func (e *ErrPullRequestAlreadyExists) Error() string {
	return fmt.Sprintf("pull request already exists: %v", e.ID)
}

type ErrPullRequestNotFound struct {
	ID uuid.UUID
}

func (e *ErrPullRequestNotFound) Error() string {
	return fmt.Sprintf("pull request not found: %v", e.ID)
}

type ErrPullRequestAlreadyMerged struct {
	ID uuid.UUID
}

func (e *ErrPullRequestAlreadyMerged) Error() string {
	return fmt.Sprintf("pull request alreaddy merged: %v", e.ID)
}

type ErrReviewerNotAssigned struct {
	ID uuid.UUID
}

func (e *ErrReviewerNotAssigned) Error() string {
	return fmt.Sprintf("reviewer not assigned: %v", e.ID)
}

var ErrNoReplacementCandidate = errors.New("no replacement candidate")

type ErrPullRequestNameValidation struct {
	Reason string
}

func (e *ErrPullRequestNameValidation) Error() string {
	return fmt.Sprintf("pull request name is invalid: %s", e.Reason)
}
