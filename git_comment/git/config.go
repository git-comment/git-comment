package git

import (
	"github.com/kylef/result.go/src/result"
	git "gopkg.in/libgit2/git2go.v23"
)

func ConfiguredString(repoPath, name, fallback string) string {
	return WithConfig(repoPath, func(config *git.Config) result.Result {
		return result.NewResult(config.LookupString(name))
	}).Recover(fallback).(string)
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
