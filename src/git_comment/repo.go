package git_comment

import (
	"errors"
	"fmt"
	git "gopkg.in/libgit2/git2go.v22"
)

const (
	authorNotFoundError  = "No name or email found in git config for commenting"
	commitNotFoundError  = "Commit not found"
	commentNotFoundError = "Comment not found"
	userNameKey          = "user.name"
	userEmailKey         = "user.email"
	headCommit           = "HEAD"
)

// Create a new comment on a commit, optionally with a file and line
func CreateComment(repoPath string, commit *string, fileRef *FileRef, message string) (*string, error) {
	var hash = commit
	author, cErr := author(repoPath)
	if cErr != nil {
		return nil, cErr
	}
	if *hash == "HEAD" || hash == nil {
		head, hErr := head(repoPath)
		if hErr != nil {
			return nil, hErr
		}
		hash = head
	} else {
		// check if commit hash exists
	}
	comment, err := NewComment(message, *hash, fileRef, author)
	if err != nil {
		return nil, err
	}
	fmt.Println(comment.Serialize())
	writeErr := writeCommentToDisk(repoPath, comment)
	if writeErr != nil {
		return nil, writeErr
	}
	return &comment.ID, nil
}

func UpdateComment(repoPath string, ID string, message string) error {
	comment, err := CommentByID(repoPath, ID)
	if err != nil {
		return err
	}
	author, cErr := author(repoPath)
	if cErr != nil {
		return cErr
	}
	comment.Amend(message, author)
	return writeCommentToDisk(repoPath, comment)
}

// Remove a comment from a commit
func DeleteComment(repoPath string, ID string) error {
	comment, err := CommentByID(repoPath, ID)
	if err != nil {
		return err
	}
	comment.Deleted = true
	return writeCommentToDisk(repoPath, comment)
}

// Finds a comment by a given ID
func CommentByID(repoPath string, identifier string) (*Comment, error) {
	return &Comment{}, errors.New(commentNotFoundError)
}

// Finds all comments on a given commit
func CommentsOnCommit(repoPath string, commit string) []*Comment {
	return []*Comment{}
}

// Write git object for a given comment and update the
// comment refs
func writeCommentToDisk(repoPath string, comment *Comment) error {
	return nil
}

func repo(repoPath string) (*git.Repository, error) {
	return git.OpenRepository(repoPath)
}

func author(repoPath string) (*Person, error) {
	repo, err := repo(repoPath)
	if err != nil {
		return nil, err
	}
	config, cErr := repo.Config()
	if cErr != nil {
		return nil, cErr
	}
	name, nErr := config.LookupString(userNameKey)
	email, eErr := config.LookupString(userEmailKey)
	if nErr != nil && eErr != nil {
		return nil, errors.New(authorNotFoundError)
	}
	return &Person{name, email}, nil
}

func head(repoPath string) (*string, error) {
	repo, err := repo(repoPath)
	if err != nil {
		return nil, err
	}
	head, hErr := repo.Head()
	if hErr != nil {
		return nil, hErr
	}
	hash := head.Name()
	return &hash, nil
}
