package git_comment

//package git

import (
	"errors"
	"fmt"
	git "gopkg.in/libgit2/git2go.v22"
)

const (
	headCommit            = "HEAD"
	noParentError         = "No parent commit to compare"
	noCommitsMatchedError = "No commits found for '%v'"
)

// Resolve a single commit from a given commitish string
func ResolvedCommit(repoPath string, commitish *string) (*string, error) {
	repo, err := git.OpenRepository(repoPath)
	if err != nil {
		return nil, err
	}
	return ResolveSingleCommitHash(repo, commitish)
}

// Parse the commit hash of a single commit from a reference, converting to HEAD
// where needed
func ResolveSingleCommitHash(repo *git.Repository, commitish *string) (*string, error) {
	var hash string
	if commitish == nil || len(*commitish) == 0 {
		hash = headCommit
	} else {
		hash = *commitish
	}
	object, err := repo.RevparseSingle(hash)
	if err != nil {
		return nil, err
	}
	id := object.Id().String()
	return &id, nil
}

// Parse commits from commitish string. If multiple commits are resolved from
// the string, returns the parent and child respectively.
func ResolveCommits(repo *git.Repository, commitish string) (*git.Commit, *git.Commit, error) {
	spec, err := repo.Revparse(commitish)
	if err != nil {
		return nil, nil, err
	}
	var toCommit *git.Commit
	var fromCommit *git.Commit
	to := spec.To()
	if to != nil && to.Type() == git.ObjectCommit {
		toCommit, err = repo.LookupCommit(to.Id())
		if err != nil {
			return nil, nil, err
		}
	}
	from := spec.From()
	if from != nil && from.Type() == git.ObjectCommit {
		fromCommit, err = repo.LookupCommit(from.Id())
		if err != nil {
			return nil, nil, err
		}
	}
	if toCommit == nil {
		toCommit = fromCommit
		fromCommit = nil
	}
	if toCommit == nil {
		return nil, nil, errors.New(fmt.Sprintf(noCommitsMatchedError, commitish))
	}
	if fromCommit == nil {
		if fromCommit = toCommit.Parent(0); fromCommit == nil {
			return nil, nil, errors.New(noParentError)
		}
	}
	return fromCommit, toCommit, nil
}

// Find all intermediate commits between a parent and child
func CommitsFromRange(fromCommit *git.Commit, toCommit *git.Commit) []*git.Commit {
	commits := make([]*git.Commit, 0)
	if toCommit != nil && fromCommit != nil {
		commit := toCommit
		for {
			if commit == nil {
				break
			}
			commits = append(commits, commit)
			if commit.Id().String() == fromCommit.Id().String() {
				break
			}
			commit = commit.Parent(0)
		}
	} else if fromCommit != nil {
		commits = append(commits, fromCommit)
	} else if toCommit != nil {
		commits = append(commits, toCommit)
	}
	return commits
}
