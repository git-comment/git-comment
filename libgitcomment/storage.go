package libgitcomment

import (
	"errors"
	"fmt"
	gg "git"
	"github.com/kylef/result.go/src/result"
	git "gopkg.in/libgit2/git2go.v23"
	"path"
)

const (
	CommentStorageDir    = ".git/comments"
	maxCommentsOnCommit  = 4096
	defaultMessageFormat = "Created a comment ref on [%v] to [%v]"
	maxCommentError      = "Maximum comments on [%v] reached."
)

// Create a new comment on a commit, optionally with a file and line
// @return result.Result<*string, error>
func CreateComment(repoPath, commit, author, message string, fileRef *FileRef) result.Result {
	return gg.WithRepository(repoPath, func(repo *git.Repository) result.Result {
		return validatedCommitForComment(repo, commit).FlatMap(func(hash interface{}) result.Result {
			return commentAuthor(repoPath, author).FlatMap(func(author interface{}) result.Result {
				return NewComment(message, *(hash).(*string), fileRef, author.(*Person))
			}).FlatMap(func(value interface{}) result.Result {
				comment := value.(*Comment)
				success := writeCommentToDisk(repo, comment)
				return success.FlatMap(func(value interface{}) result.Result {
					return result.NewSuccess(comment.ID)
				})
			})
		})
	})
}

// Update an existing comment with a new message
// @return result.Result<*Comment, error>
func UpdateComment(repoPath, identifier, committer, message string) result.Result {
	return gg.WithRepository(repoPath, func(repo *git.Repository) result.Result {
		return CommentByID(repo, identifier).FlatMap(func(c interface{}) result.Result {
			comment := c.(*Comment)
			return commentCommitter(repoPath, committer).FlatMap(func(committer interface{}) result.Result {
				comment.Amend(message, committer.(*Person))
				return writeCommentToDisk(repo, comment)
			})
		})
	})
}

// Remove a comment from a commit
// @return result.Result<*Comment, error>
func DeleteComment(repoPath string, identifier string) result.Result {
	return gg.WithRepository(repoPath, func(repo *git.Repository) result.Result {
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
	return gg.CommitRefDir(*comment.Commit).FlatMap(func(dir interface{}) result.Result {
		return result.NewSuccess(path.Join(dir.(string), identifier))
	})
}

// Determine author for comment preferring the author string if
// available.
// @return result.Result<*Person, error>
func commentAuthor(repoPath, author string) result.Result {
	if len(author) > 0 {
		return CreatePerson(author)
	}
	return CreatePerson(gg.ConfiguredAuthor(repoPath))
}

// Determine committer for comment preferring the committer string if
// available.
// @return result.Result<*Person, error>
func commentCommitter(repoPath, committer string) result.Result {
	if len(committer) > 0 {
		return CreatePerson(committer)
	}
	return CreatePerson(gg.ConfiguredCommitter(repoPath))
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
	return gg.CreateBlob(repo, comment.Serialize()).FlatMap(func(oid interface{}) result.Result {
		id := fmt.Sprintf("%v", oid)
		return RefPath(comment, id).FlatMap(func(file interface{}) result.Result {
			commit := *comment.Commit
			message := fmt.Sprintf(defaultMessageFormat, commit[:7], id[:7])
			return result.NewResult(repo.References.Create(file.(string), oid.(*git.Oid), false, message))
		}).FlatMap(func(value interface{}) result.Result {
			comment.ID = &id
			return result.NewSuccess(comment)
		})
	})
}

func validatedCommitForComment(repo *git.Repository, commit string) result.Result {
	return gg.ResolveSingleCommitHash(repo, commit).FlatMap(func(hash interface{}) result.Result {
		return CommentCountOnCommit(repo, *(hash.(*string))).FlatMap(func(count interface{}) result.Result {
			if count.(uint16) >= maxCommentsOnCommit {
				return result.NewFailure(errors.New(maxCommentError))
			}
			return result.NewSuccess(hash)
		})
	})
}

func deleteReference(repo *git.Repository, comment *Comment, identifier string) error {
	_, err := RefPath(comment, identifier).FlatMap(func(refPath interface{}) result.Result {
		return result.NewResult(repo.References.Lookup(refPath.(string)))
	}).Analysis(func(ref interface{}) result.Result {
		return gg.DeleteReference(ref.(*git.Reference))
	}, func(err error) result.Result {
		return result.NewFailure(errors.New(commentNotFoundError))
	}).Dematerialize()
	return err
}
