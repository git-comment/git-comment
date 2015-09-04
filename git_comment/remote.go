package git_comment

import (
	"errors"
	"fmt"
	gitg "git_comment/git"
	"github.com/kylef/result.go/src/result"
	git "gopkg.in/libgit2/git2go.v23"
)

const (
	commentDefaultFetch = "+refs/comments/*:refs/remotes/%v/comments/*"
	commentDefaultPush  = "refs/comments/*"
)

// Configure a remote to fetch and push comment changes by default
// @return result.Result<bool, error>
func ConfigureRemoteForComments(repoPath, remoteName string) result.Result {
	return gitg.WithRemote(repoPath, remoteName, func(remote *git.Remote) result.Result {
		success := func(values ...interface{}) result.Result {
			return result.NewSuccess(true)
		}
		return gitg.WithRepository(repoPath, func(repo *git.Repository) result.Result {
			pushRef := commentDefaultPush
			fetchRef := fmt.Sprintf(commentDefaultFetch, remoteName)
			return result.Combine(success, gitg.AddPush(repo, remote, pushRef), gitg.AddFetch(repo, remote, fetchRef))
		})
	})
}

func DeleteRemoteComment(repoPath, remoteName, commentID string) result.Result {
	return gitg.WithRemote(repoPath, remoteName, func(remote *git.Remote) result.Result {
		return CreatePerson(gitg.ConfiguredCommitter(repoPath)).Analysis(func(val interface{}) result.Result {
			sig := val.(*Person).Signature()
			refspec := fmt.Sprintf(":%v", commentID)
			return gitg.Push(repoPath, remoteName, []string{refspec}, sig)
		}, func(err error) result.Result {
			return result.NewFailure(errors.New(noCommitterError))
		})
	})
}
