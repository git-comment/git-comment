package git_comment

import (
	"errors"
	"fmt"
	git "gopkg.in/libgit2/git2go.v22"
	"path"
	"time"
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
	if err := writeCommentToDisk(repo, comment); err != nil {
		return nil, err
	}

	id := comment.ID
	return id, nil
}

// Update an existing comment with a new message
func UpdateComment(repoPath string, ID string, message string) error {
	repo, err := repo(repoPath)
	if err != nil {
		return err
	}
	comment, err := CommentByID(repo, ID)
	if err != nil {
		return err
	}
	author, err := author(repo)
	if err != nil {
		return err
	}
	comment.Amend(message, author)
	return writeCommentToDisk(repo, comment)
}

// Remove a comment from a commit
func DeleteComment(repoPath string, ID string) error {
	repo, err := repo(repoPath)
	if err != nil {
		return err
	}
	comment, err := CommentByID(repo, ID)
	if err != nil {
		return err
	}
	comment.Deleted = true
	return writeCommentToDisk(repo, comment)
}

// Finds a comment by a given ID
func CommentByID(repo *git.Repository, identifier string) (*Comment, error) {
	return &Comment{}, errors.New(commentNotFoundError)
}

// Finds all comments on a given commit
func CommentsOnCommit(repoPath string, commit string) []*Comment {
	return []*Comment{}
}

// Write git object for a given comment and update the
// comment refs
func writeCommentToDisk(repo *git.Repository, comment *Comment) error {
	oid, err := repo.CreateBlobFromBuffer([]byte(comment.Serialize()))
	if err != nil {
		return err
	}
	committer := comment.Amender
	sig := &git.Signature{committer.Name, committer.Email, time.Now()}
	id := fmt.Sprintf("%v", oid)
	file, err := refPath(comment, &id)
	if err != nil {
		return err
	}
	_, err = repo.CreateReference(*file, oid, false, sig, "some message")
	if err != nil {
		return err
	}
	comment.ID = &id
	return nil
}

// Generate the path within refs for a given comment
//
// Comment refs are nested under refs/comments. The
// format is as follows:
//
// ```
// refs/comments/[<commit prefix>]/[<rest of commit>]/[<comment id>]
// ```
//
func refPath(comment *Comment, id *string) (*string, error) {
	dir, err := commitRefDir(&comment.Commit)
	if err != nil {
		return nil, err
	}
	hash := path.Join(*dir, *id)
	return &hash, nil
}

func commitRefDir(commit *string) (*string, error) {
	const invalidHash = "Invalid commit hash for storage"
	const commentPath = "refs/comments"
	hash := *commit
	if len(hash) > 4 {
		dir := path.Join(commentPath,
			hash[0:4],
			hash[4:len(hash)])
		return &dir, nil
	}
	return nil, errors.New(invalidHash)
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
