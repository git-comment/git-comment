package git_comment

import (
	"fmt"
	"github.com/kylef/result.go/src/result"
	git "github.com/libgit2/git2go"
)

const (
	commentDefaultFetch = "+refs/comments/*:refs/remotes/%v/comments/*"
	commentDefaultPush  = "refs/comments/*"
)

// Configure a remote to fetch and push comment changes by default
// @return result.Result<bool, error>
func ConfigureRemoteForComments(repoPath, remoteName string) result.Result {
	return WithRemote(repoPath, remoteName, func(remote *git.Remote) result.Result {
		pushRef := commentDefaultPush
		fetchRef := fmt.Sprintf(commentDefaultFetch, remoteName)
		success := func(values ...interface{}) result.Result {
			if err := remote.Save(); err != nil {
				return result.NewFailure(err)
			}
			return result.NewSuccess(true)
		}
		return result.Combine(success, addPush(remote, pushRef), addFetch(remote, fetchRef))
	})
}

func ConfiguredString(repoPath, name string) result.Result {
	return WithConfig(repoPath, func(config *git.Config) result.Result {
		return result.NewResult(config.LookupString(name))
	})
}

func ConfiguredInt32(repoPath, name string, fallback int32) int32 {
	return WithConfig(repoPath, func(config *git.Config) result.Result {
		return result.NewResult(config.LookupInt32(name))
	}).Recover(fallback).(int32)
}

func ConfiguredBool(repoPath, name string, fallback bool) bool {
	return WithConfig(repoPath, func(config *git.Config) result.Result {
		return result.NewResult(config.LookupBool(name))
	}).Recover(fallback).(bool)
}

func WithConfig(repoPath string, ifSuccess func(config *git.Config) result.Result) result.Result {
	return WithRepository(repoPath, func(repo *git.Repository) result.Result {
		return result.NewResult(repo.Config())
	}).FlatMap(func(config interface{}) result.Result {
		return ifSuccess(config.(*git.Config))
	})
}
