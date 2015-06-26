package git_comment

import (
	"errors"
	"fmt"
	git "gopkg.in/libgit2/git2go.v22"
	"path"
	"time"
)

const (
	DefaultMessageTemplate = "\n# Enter comment content\n# Lines beginning with '#' will be stripped"
	defaultMessageFormat   = "Created a comment ref on [%v] to [%v]"
)

// Create a new comment on a commit, optionally with a file and line
func CreateComment(repoPath string, commit *string, fileRef *FileRef, message string) (*string, error) {
	repo, err := git.OpenRepository(repoPath)
	if err != nil {
		return nil, err
	}
	author, err := ConfiguredAuthor(repo)
	if err != nil {
		return nil, err
	}
	hash, err := ResolveSingleCommitHash(repo, commit)
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

	return comment.ID, nil
}

// Update an existing comment with a new message
func UpdateComment(repoPath string, ID string, message string) (*string, error) {
	repo, err := git.OpenRepository(repoPath)
	if err != nil {
		return nil, err
	}
	comment, err := CommentByID(repo, ID)
	if err != nil {
		return nil, err
	}
	committer, err := ConfiguredCommitter(repo)
	if err != nil {
		return nil, err
	}
	comment.Amend(message, committer)
	if err := writeCommentToDisk(repo, comment); err != nil {
		return nil, err
	}

	return comment.ID, nil
}

// Remove a comment from a commit
func DeleteComment(repoPath string, ID string) error {
	repo, err := git.OpenRepository(repoPath)
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

// Generate the path within refs for a given comment
//
// Comment refs are nested under refs/comments. The
// format is as follows:
//
// ```
// refs/comments/[<commit prefix>]/[<rest of commit>]/[<comment id>]
// ```
//
func RefPath(comment *Comment, id *string) (*string, error) {
	dir, err := CommitRefDir(comment.Commit)
	if err != nil {
		return nil, err
	}
	hash := path.Join(*dir, *id)
	return &hash, nil
}

// Base reference path for a commit
func CommitRefDir(commit *string) (*string, error) {
	const commentPath = "refs/comments"
	hash := *commit
	if len(hash) > 4 {
		dir := path.Join(commentPath,
			hash[:4],
			hash[4:len(hash)])
		return &dir, nil
	}
	return nil, errors.New(invalidHashError)
}

// Write git object for a given comment and update the
// comment refs
func writeCommentToDisk(repo *git.Repository, comment *Comment) error {
	if comment.ID != nil {
		err := deleteReference(repo, comment, *comment.ID)
		if err != nil {
			return err
		}
	}
	oid, err := repo.CreateBlobFromBuffer([]byte(comment.Serialize()))
	if err != nil {
		return err
	}
	committer := comment.Amender
	sig := &git.Signature{committer.Name, committer.Email, time.Now()}
	id := fmt.Sprintf("%v", oid)
	file, err := RefPath(comment, &id)
	if err != nil {
		return err
	}
	commit := *comment.Commit
	message := fmt.Sprintf(defaultMessageFormat, commit[:7], id[:7])
	_, err = repo.CreateReference(*file, oid, false, sig, message)
	if err != nil {
		return err
	}
	comment.ID = &id
	return nil
}

func deleteReference(repo *git.Repository, comment *Comment, identifier string) error {
	refPath, err := RefPath(comment, &identifier)
	if err != nil {
		return nil
	}
	ref, err := repo.LookupReference(*refPath)
	if err != nil {
		return errors.New(commentNotFoundError)
	}
	err = ref.Delete()
	if err != nil {
		return err
	}
	return nil
}
