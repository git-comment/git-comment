package git

//package git

import (
	"errors"
	"github.com/kylef/result.go/src/result"
	git "github.com/libgit2/git2go"
)

const (
	headCommit            = "HEAD"
	noParentError         = "No parent commit to compare"
	noCommitsMatchedError = "No commits found for '%v'"
	noCommitError         = "No commmit found"
)

// Resolve a single commit from a given commitish string
//
// return result.Result<*string, error>
func ResolvedCommit(repoPath, commitish string) result.Result {
	return WithRepository(repoPath, func(repo *git.Repository) result.Result {
		return ResolveSingleCommitHash(repo, commitish)
	})
}

// Parse the commit hash of a single commit from a reference, converting to HEAD
// where needed
//
// return result.Result<*string, error>
func ResolveSingleCommitHash(repo *git.Repository, commitish string) result.Result {
	return result.NewResult(repo.RevparseSingle(ExpandCommitish(commitish))).FlatMap(getObjectId)
}

// Parse commits from commitish string, populating a CommitRange. If a
// single commit is matched, it is paired with its first parent commit
// or resolves to an error
//
// return result.Result<*git.Commit, error>
func ResolveCommits(repo *git.Repository, commitish string) result.Result {
	return result.NewResult(repo.Revparse(commitish)).FlatMap(func(value interface{}) result.Result {
		spec := value.(*git.Revspec)
		return resolveCommit(repo, spec.From()).FlatMap(func(f interface{}) result.Result {
			fromCommit := f.(*git.Commit)
			return resolveCommit(repo, spec.To()).Analysis(func(t interface{}) result.Result {
				toCommit := t.(*git.Commit)
				return result.NewSuccess(&CommitRange{fromCommit, toCommit})
			}, func(err error) result.Result {
				return resolveCommitParent(fromCommit).FlatMap(func(parent interface{}) result.Result {
					return result.NewSuccess(&CommitRange{parent.(*git.Commit), fromCommit})
				})
			})
		})
	})
}

// Resolve empty or nil strings to HEAD
func ExpandCommitish(commitish string) string {
	if len(commitish) == 0 {
		return headCommit
	}
	return commitish
}

func getObjectId(value interface{}) result.Result {
	id := value.(git.Object).Id().String()
	return result.NewSuccess(&id)
}

func resolveCommit(repo *git.Repository, object git.Object) result.Result {
	if object != nil && object.Type() == git.ObjectCommit {
		return result.NewResult(repo.LookupCommit(object.Id()))
	}
	return result.NewFailure(errors.New(noCommitError))
}

func resolveCommitParent(commit *git.Commit) result.Result {
	parent := commit.Parent(0)
	if parent == nil {
		return result.NewFailure(errors.New(noParentError))
	}
	return result.NewSuccess(parent)
}
