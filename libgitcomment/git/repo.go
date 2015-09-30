package git

import (
	"errors"
	"path"

	"github.com/kylef/result.go/src/result"
	git "gopkg.in/libgit2/git2go.v23"
)

const (
	invalidHashError = "Invalid commit hash for storage"
	CommentRefBase   = "refs/comments"
	glob             = "*"
)

// @return result.Result<*git.Repository, error>
func OpenRepository(repoPath string) result.Result {
	return result.NewResult(git.OpenRepository(repoPath))
}

func WithRepository(repoPath string, ifSuccess func(repo *git.Repository) result.Result) result.Result {
	return OpenRepository(repoPath).FlatMap(func(value interface{}) result.Result {
		return ifSuccess(value.(*git.Repository))
	})
}

// Find blob for an ID
// @return result.Result<*git.Blob, error>
func LookupBlob(repo *git.Repository, identifier, errorCode string) result.Result {
	return result.NewResult(git.NewOid(identifier)).FlatMap(func(oid interface{}) result.Result {
		return result.NewResult(repo.LookupBlob(oid.(*git.Oid)))
	}).RecoverWith(result.NewFailure(errors.New(errorCode)))
}

// Find commit for an ID
// @return result.Result<*git.Commit, error>
func LookupCommit(repo *git.Repository, identifier string) result.Result {
	return result.NewResult(repo.LookupCommit(git.NewOidFromBytes([]byte(identifier))))
}

// Delete an existing reference
// @return result.Result<bool, error>
func DeleteReference(ref *git.Reference) result.Result {
	return BoolResult(true, ref.Delete())
}

// Reference iterator for all comments on a commit
// @return result.Result<*git.ReferenceIterator, error>
func CommitCommentRefIterator(repo *git.Repository, commitHash string, iteration func(ref *git.Reference)) result.Result {
	return CommitRefDir(commitHash).FlatMap(func(dir interface{}) result.Result {
		return IterateRefs(repo, path.Join(dir.(string), glob), iteration)
	})
}

// Reference iterator for all comments
// @return result.Result<*git.ReferenceIterator, error>
func CommentRefIterator(repo *git.Repository, iteration func(ref *git.Reference)) result.Result {
	return IterateRefs(repo, path.Join(CommentRefBase, glob), iteration)
}

func IterateRefs(repo *git.Repository, refPathGlob string, iteration func(ref *git.Reference)) result.Result {
	iterator := result.NewResult(repo.NewReferenceIteratorGlob(refPathGlob))
	return iterator.FlatMap(func(i interface{}) result.Result {
		iterator := i.(*git.ReferenceIterator)
		ref, err := iterator.Next()
		for {
			if git.IsErrorCode(err, git.ErrIterOver) {
				break
			} else if err != nil {
				return result.NewFailure(err)
			}
			iteration(ref)
			ref, err = iterator.Next()
		}
		return result.NewSuccess(true)
	})
}

// Creates a new blob from string content
// @return result.Result<*git.Oid, error>
func CreateBlob(repo *git.Repository, content string) result.Result {
	return result.NewResult(repo.CreateBlobFromBuffer([]byte(content)))
}

// Base reference path for a commit
// @return result.Result<string, error>
func CommitRefDir(hash string) result.Result {
	if len(hash) > 4 {
		return result.NewSuccess(path.Join(CommentRefBase, hash[:4], hash[4:len(hash)]))
	}
	return result.NewFailure(errors.New(invalidHashError))
}
