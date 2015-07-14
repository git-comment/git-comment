package git_comment

import (
	"errors"
	"github.com/kylef/result.go/src/result"
	git "gopkg.in/libgit2/git2go.v22"
	"path"
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
func LookupBlob(repo *git.Repository, identifier string) result.Result {
	return result.NewResult(git.NewOid(identifier)).FlatMap(func(oid interface{}) result.Result {
		return result.NewResult(repo.LookupBlob(oid.(*git.Oid)))
	}).RecoverWith(result.NewFailure(errors.New(commentNotFoundError)))
}

// Find commit for an ID
// @return result.Result<*git.Commit, error>
func LookupCommit(repo *git.Repository, identifier string) result.Result {
	return result.NewResult(repo.LookupCommit(git.NewOidFromBytes([]byte(identifier))))
}

// Delete an existing reference
// @return result.Result<bool, error>
func DeleteReference(ref *git.Reference) result.Result {
	if err := ref.Delete(); err != nil {
		return result.NewFailure(err)
	}
	return result.NewSuccess(true)
}

// Reference iterator for all comments on a commit
// @return result.Result<*git.ReferenceIterator, error>
func CommentRefIterator(repo *git.Repository, commitHash string) result.Result {
	const glob = "*"
	return CommitRefDir(commitHash).FlatMap(func(dir interface{}) result.Result {
		return result.NewResult(repo.NewReferenceIteratorGlob(path.Join(dir.(string), glob)))
	})
}

// Creates a new blob from string content
// @return result.Result<*git.Oid, error>
func CreateBlob(repo *git.Repository, content string) result.Result {
	return result.NewResult(repo.CreateBlobFromBuffer([]byte(content)))
}
