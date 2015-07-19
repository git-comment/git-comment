package git_comment

import (
	"github.com/kylef/result.go/src/result"
	git "github.com/libgit2/git2go"
)

// Lookup a remote by name, performing a block if found
// @return result.Result<*git.Remote, error>
func WithRemote(repoPath, remoteName string, ifSuccess func(*git.Remote) result.Result) result.Result {
	return WithRepository(repoPath, func(repo *git.Repository) result.Result {
		return result.NewResult(repo.LookupRemote(remoteName))
	}).FlatMap(func(remote interface{}) result.Result {
		return ifSuccess(remote.(*git.Remote))
	})
}

// Add a push refspec to a remote. Return true if added.
// @return result.Result<bool, error>
func addPush(remote *git.Remote, pushRef string) result.Result {
	p := result.NewResult(remote.PushRefspecs())
	return p.FlatMap(func(pushes interface{}) result.Result {
		if !contains(pushes.([]string), pushRef) {
			return BoolResult(true, remote.AddPush(pushRef))
		}
		return result.NewSuccess(false)
	})
}

// Add a fetch refspec to a remote. Return true if added.
// @return result.Result<bool, error>
func addFetch(remote *git.Remote, fetchRef string) result.Result {
	f := result.NewResult(remote.FetchRefspecs())
	return f.FlatMap(func(fetches interface{}) result.Result {
		if !contains(fetches.([]string), fetchRef) {
			return BoolResult(true, remote.AddFetch(fetchRef))
		}
		return result.NewSuccess(false)
	})
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
