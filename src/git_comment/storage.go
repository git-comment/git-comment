package git_comment

import (
	"errors"
	"fmt"
	"github.com/kylef/result.go/src/result"
	git "github.com/libgit2/git2go"
	"path"
	"time"
)

const (
	DefaultMessageTemplate = "\n# Enter comment content\n# Lines beginning with '#' will be stripped"
	defaultMessageFormat   = "Created a comment ref on [%v] to [%v]"
)

// Create a new comment on a commit, optionally with a file and line
// @return result.Result<*string, error>
func CreateComment(repoPath string, commit *string, fileRef *FileRef, message string) result.Result {
	repo := OpenRepository(repoPath)
	return ConfiguredAuthor(repoPath).FlatMap(func(author interface{}) result.Result {
		return repo.FlatMap(func(value interface{}) result.Result {
			return ResolveSingleCommitHash(value.(*git.Repository), commit)
		}).FlatMap(func(hash interface{}) result.Result {
			return NewComment(message, *(hash).(*string), fileRef, author.(*Person))
		}).FlatMap(func(value interface{}) result.Result {
			comment := value.(*Comment)
			success := writeCommentToDisk(repo.Success.(*git.Repository), comment)
			return success.FlatMap(func(value interface{}) result.Result {
				return result.NewSuccess(comment.ID)
			})
		})
	})
}

// Update an existing comment with a new message
// @return result.Result<*Comment, error>
func UpdateComment(repoPath string, identifier string, message string) result.Result {
	return WithRepository(repoPath, func(repo *git.Repository) result.Result {
		return CommentByID(repo, identifier).FlatMap(func(c interface{}) result.Result {
			comment := c.(*Comment)
			return ConfiguredCommitter(repoPath).FlatMap(func(committer interface{}) result.Result {
				comment.Amend(message, committer.(*Person))
				return writeCommentToDisk(repo, comment)
			})
		})
	})
}

// Remove a comment from a commit
// @return result.Result<*Comment, error>
func DeleteComment(repoPath string, identifier string) result.Result {
	return WithRepository(repoPath, func(repo *git.Repository) result.Result {
		return CommentByID(repo, identifier).FlatMap(func(c interface{}) result.Result {
			comment := c.(*Comment)
			comment.Deleted = true
			return writeCommentToDisk(repo, comment)
		})
	})
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
// @return result.Result<string, error>
func RefPath(comment *Comment, identifier string) result.Result {
	return CommitRefDir(*comment.Commit).FlatMap(func(dir interface{}) result.Result {
		return result.NewSuccess(path.Join(dir.(string), identifier))
	})
}

// Base reference path for a commit
// @return result.Result<string, error>
func CommitRefDir(hash string) result.Result {
	const commentPath = "refs/comments"
	if len(hash) > 4 {
		return result.NewSuccess(path.Join(commentPath, hash[:4], hash[4:len(hash)]))
	}
	return result.NewFailure(errors.New(invalidHashError))
}

// Write git object for a given comment and update the
// comment refs
// @return result.Result<*Comment, error>
func writeCommentToDisk(repo *git.Repository, comment *Comment) result.Result {
	if comment.ID != nil {
		if err := deleteReference(repo, comment, *comment.ID); err != nil {
			return result.NewFailure(err)
		}
	}
	return CreateBlob(repo, comment.Serialize()).FlatMap(func(oid interface{}) result.Result {
		id := fmt.Sprintf("%v", oid)
		return RefPath(comment, id).FlatMap(func(file interface{}) result.Result {
			committer := comment.Amender
			sig := &git.Signature{committer.Name, committer.Email, time.Now()}
			commit := *comment.Commit
			message := fmt.Sprintf(defaultMessageFormat, commit[:7], id[:7])
			return result.NewResult(repo.CreateReference(file.(string), oid.(*git.Oid), false, sig, message))
		}).FlatMap(func(value interface{}) result.Result {
			comment.ID = &id
			return result.NewSuccess(comment)
		})
	})
}

func deleteReference(repo *git.Repository, comment *Comment, identifier string) error {
	_, err := RefPath(comment, identifier).FlatMap(func(refPath interface{}) result.Result {
		return result.NewResult(repo.LookupReference(refPath.(string)))
	}).Analysis(func(ref interface{}) result.Result {
		return DeleteReference(ref.(*git.Reference))
	}, func(err error) result.Result {
		return result.NewFailure(errors.New(commentNotFoundError))
	}).Dematerialize()
	return err
}
