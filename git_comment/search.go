package git_comment

import (
	gitg "git_comment/git"
	"github.com/kylef/result.go/src/result"
	git "github.com/libgit2/git2go"
	"path"
	"sort"
)

// Find all comments matching text
// @return result.Result<[]*Comment, error>
func CommentsWithContent(content string) result.Result {
	return result.NewSuccess(nil)
}

// Finds a comment by ID
// @return result.Result<*Comment, error>
func CommentByID(repo *git.Repository, identifier string) result.Result {
	return gitg.LookupBlob(repo, identifier, commentNotFoundError).FlatMap(func(blob interface{}) result.Result {
		return DeserializeComment(string(blob.(*git.Blob).Contents()))
	}).FlatMap(func(c interface{}) result.Result {
		comment := c.(*Comment)
		comment.ID = &identifier
		return result.NewSuccess(comment)
	})
}

// Find comments in a commit range or on a single commit
// @return result.Result<[]*Comment, error>
func CommentsOnCommittish(repoPath string, committish string) result.Result {
	return gitg.WithRepository(repoPath, func(repo *git.Repository) result.Result {
		resolution := gitg.ResolveCommits(repo, committish)
		return resolution.FlatMap(func(commitRange interface{}) result.Result {
			return CommentsOnCommits(repo, commitRange.(*gitg.CommitRange).Commits())
		})
	})
}

// Count comments on commit
// @return result.Result<uint16, error>
func CommentCountOnCommit(repo *git.Repository, commit string) result.Result {
	return gitg.CommentRefIterator(repo, commit).FlatMap(func(iterator interface{}) result.Result {
		refIterator := iterator.(*git.ReferenceIterator)
		_, err := refIterator.Next()
		var count uint16 = 0
		for {
			if err != nil && err.(*git.GitError).Code == git.ErrIterOver {
				break
			} else if err != nil {
				return result.NewFailure(err)
			}
			count += 1
			_, err = refIterator.Next()
		}
		return result.NewSuccess(count)
	})
}

// Finds all comments on a given commit
// @return result.Result<[]*Comment, error>
func CommentsOnCommit(repoPath string, commitHash *string) result.Result {
	return gitg.WithRepository(repoPath, func(repo *git.Repository) result.Result {
		hash := gitg.ResolveSingleCommitHash(repo, commitHash)
		return hash.FlatMap(func(commit interface{}) result.Result {
			return gitg.LookupCommit(repo, *(commit.(*string)))
		}).FlatMap(func(commit interface{}) result.Result {
			return commentsOnCommit(repo, commit.(*git.Commit))
		})
	})
}

// Finds all comments on an array of commits
// @return result.Result<[]*Comment, error>
func CommentsOnCommits(repo *git.Repository, commits []*git.Commit) result.Result {
	results := make([]result.Result, len(commits))
	for index, commit := range commits {
		results[index] = commentsOnCommit(repo, commit)
	}
	return result.Combine(func(values ...interface{}) result.Result {
		comments := make(CommentSlice, 0)
		for _, list := range values {
			for _, comment := range list.([]interface{}) {
				comments = append(comments, comment.(*Comment))
			}
		}
		sort.Stable(comments)
		return result.NewSuccess(comments)
	}, results...)
}

// Finds all comments on a commit
// @return result.Result<[]*Comment, error>
func commentsOnCommit(repo *git.Repository, commit *git.Commit) result.Result {
	return gitg.CommentRefIterator(repo, commit.Id().String()).FlatMap(func(iterator interface{}) result.Result {
		refIterator := iterator.(*git.ReferenceIterator)
		var comments []interface{}
		ref, err := refIterator.Next()
		for {
			if err != nil && err.(*git.GitError).Code == git.ErrIterOver {
				break
			} else if err != nil {
				return result.NewFailure(err)
			}
			commentFromRef(repo, ref.Name()).FlatMap(func(comment interface{}) result.Result {
				comments = append(comments, comment)
				return result.Result{}
			})
			ref, err = refIterator.Next()
		}
		return result.NewSuccess(comments)
	})
}

func commentFromRef(repo *git.Repository, refName string) result.Result {
	_, identifier := path.Split(refName)
	return CommentByID(repo, identifier)
}
