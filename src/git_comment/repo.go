package git_comment

import (
	"errors"
	git "gopkg.in/libgit2/git2go.v22"
)

const (
	authorNotFoundError  = "No name or email found in git config for commenting"
	commitNotFoundError  = "Commit not found"
	commentNotFoundError = "Comment not found"
	headCommit           = "HEAD"
)

// Create a new comment on a commit, optionally with a file and line
func CreateComment(repoPath string, commit *string, fileRef *FileRef, message string) (*string, error) {
	repo, err := repo(repoPath)
	if err != nil {
		return nil, err
	}
	author, err := author(repo)
	if err != nil {
		return nil, err
	}
	hash, err := parseCommit(repo, commit)
	if err != nil {
		return nil, err
	}
	comment, err := NewComment(message, *hash, fileRef, author)
	if err != nil {
		return nil, err
	}
	if err := writeCommentToDisk(repoPath, comment); err != nil {
		return nil, err
	}
	return &comment.ID, nil
}

func UpdateComment(repoPath string, ID string, message string) error {
	repo, err := repo(repoPath)
	if err != nil {
		return err
	}
	comment, err := CommentByID(repoPath, ID)
	if err != nil {
		return err
	}
	author, err := author(repo)
	if err != nil {
		return err
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

func author(repo *git.Repository) (*Person, error) {
	sig, err := repo.DefaultSignature()
	if err != nil {
		return nil, err
	}
	return &Person{sig.Name, sig.Email}, nil
}

// parse a commit hash, converting to the HEAD commit where needed
func parseCommit(repo *git.Repository, commit *string) (*string, error) {
	var hash string
	var id string
	if commit == nil {
		hash = headCommit
	} else {
		hash = *commit
	}
	ref, err := repo.LookupReference(hash)
	if err != nil {
		oid, err := git.NewOid(hash)
		if err != nil {
			return nil, errors.New(commitNotFoundError)
		}
		obj, err := repo.Lookup(oid)
		if err != nil {
			return nil, errors.New(commitNotFoundError)
		}
		id = obj.Id().String()
		return nil, errors.New(commitNotFoundError)
	}
	res, err := ref.Resolve()
	if err != nil {
		return nil, err
	}
	id = res.Target().String()
	return &id, nil
}

func head(repo *git.Repository) (*string, error) {
	head, hErr := repo.Head()
	if hErr != nil {
		return nil, hErr
	}
	hash := head.Name()
	return &hash, nil
}
