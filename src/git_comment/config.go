package git_comment

import (
	"errors"
	"fmt"
	git "gopkg.in/libgit2/git2go.v22"
	"os"
)

// Configure a remote to fetch and push comment changes by default
func ConfigureRemoteForComments(repoPath string, remoteName string) error {
	const (
		commentDefaultFetch = "+refs/comments/*:refs/remotes/%v/comments/*"
		commentDefaultPush  = "refs/comments/*"
	)
	repo, err := git.OpenRepository(repoPath)
	if err != nil {
		return err
	}
	remote, err := repo.LookupRemote(remoteName)
	if err != nil {
		return err
	}
	fetch := fmt.Sprintf(commentDefaultFetch, remoteName)
	fetches, err := remote.FetchRefspecs()
	if err != nil {
		return err
	}
	if !contains(fetches, fetch) {
		err = remote.AddFetch(fetch)
		if err != nil {
			return err
		}
	}
	pushes, err := remote.PushRefspecs()
	if err != nil {
		return err
	}
	if !contains(pushes, commentDefaultPush) {
		err = remote.AddPush(commentDefaultPush)
		if err != nil {
			return err
		}
	}
	if err = remote.Save(); err != nil {
		return err
	}
	return nil
}

// The editor to use for editing comments interactively.
// Emulates the behavior of `git-var(1)` to determine which
// editor to use from this list of options:
//
// * `$GIT_EDITOR` environment variable
// * `core.editor` configuration
// * `$VISUAL`
// * `$EDITOR`
// * vi
func ConfiguredEditor(repoPath string) *string {
	const defaultEditor = "vi"
	repo, err := git.OpenRepository(repoPath)
	if err != nil {
		return nil
	}

	if gitEditor := os.Getenv("GIT_EDITOR"); len(gitEditor) > 0 {
		return &gitEditor
	}
	config, err := repo.Config()
	if err == nil {
		confEditor, err := config.LookupString("core.editor")
		if err == nil {
			if len(confEditor) > 0 {
				return &confEditor
			}
		}
	}

	if visual := os.Getenv("VISUAL"); len(visual) > 0 {
		return &visual
	} else if envEditor := os.Getenv("EDITOR"); len(envEditor) > 0 {
		return &envEditor
	}
	editor := defaultEditor
	return &editor
}

// The text viewer to use for viewing text interactively.
// Emulates the behavior of `git-var(1)` by checking the
// options in this list of options:
//
// * `$GIT_PAGER` environment variable
// * `core.pager` configuration
// * `$PAGER`
// * less
func ConfiguredPager(repoPath string) *string {
	const defaultPager = "less"
	repo, err := git.OpenRepository(repoPath)
	if err != nil {
		return nil
	}

	if pager := os.Getenv("GIT_PAGER"); len(pager) > 0 {
		return &pager
	}
	config, err := repo.Config()
	if err == nil {
		pager, err := config.LookupString("core.pager")
		if err == nil {
			if len(pager) > 0 {
				return &pager
			}
		}
	}

	if pager := os.Getenv("PAGER"); len(pager) > 0 {
		return &pager
	}
	pager := defaultPager
	return &pager
}

// The author of a piece of code, fetched from:
//
// * `$GIT_AUTHOR_NAME` and `$GIT_AUTHOR_EMAIL`
// * configured default from `user.name` and `user.email`
func ConfiguredAuthor(repo *git.Repository) (*Person, error) {
	// TODO: update impl
	sig, err := repo.DefaultSignature()
	if err != nil {
		return nil, errors.New(authorNotFoundError)
	}
	return &Person{sig.Name, sig.Email}, nil
}

// The committer of a piece of code
//
// * `$GIT_COMMITTER_NAME` and `$GIT_COMMITTER_EMAIL`
// * configured default from `user.name` and `user.email`
func ConfiguredCommitter(repo *git.Repository) (*Person, error) {
	return ConfiguredAuthor(repo)
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
