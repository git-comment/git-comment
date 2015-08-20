package git_comment

import (
	"errors"
	"fmt"
	gitg "git_comment/git"
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
	return gitg.WithRemote(repoPath, remoteName, func(remote *git.Remote) result.Result {
		pushRef := commentDefaultPush
		fetchRef := fmt.Sprintf(commentDefaultFetch, remoteName)
		success := func(values ...interface{}) result.Result {
			if err := remote.Save(); err != nil {
				return result.NewFailure(err)
			}
			return result.NewSuccess(true)
		}
		return result.Combine(success, gitg.AddPush(remote, pushRef), gitg.AddFetch(remote, fetchRef))
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
