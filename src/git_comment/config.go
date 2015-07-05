package git_comment

import (
	"errors"
	"fmt"
	git "gopkg.in/libgit2/git2go.v22"
	"os"
	"os/exec"
	"strings"
)

const (
	defaultEditor = "vi"
	defaultPager  = "less"
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

// The editor to use for editing comments interactively, as
// configured through git-var(1)
func ConfiguredEditor(repoPath string) *string {
	return gitVariable(repoPath, "GIT_EDITOR", defaultEditor)
}

// The text viewer to use for viewing text interactively, as
// configured through git-var(1)
func ConfiguredPager(repoPath string) *string {
	return gitVariable(repoPath, "GIT_PAGER", defaultPager)
}

// The author of a piece of code as configured through git-var(1)
func ConfiguredAuthor(repoPath string) (*Person, error) {
	author := gitVariable(repoPath, "GIT_AUTHOR_IDENT", "")
	if len(*author) == 0 {
		return nil, errors.New(authorNotFoundError)
	}
	return CreatePerson(*author), nil
}

// The committer of a piece of code as configured through git-var(1)
func ConfiguredCommitter(repoPath string) (*Person, error) {
	author := gitVariable(repoPath, "GIT_COMMITTER_IDENT", "")
	if len(*author) == 0 {
		return nil, errors.New(committerNotFoundError)
	}
	return CreatePerson(*author), nil
}

func ConfiguredString(repoPath, name string) (*string, error) {
	config, err := repoConfig(repoPath)
	if err != nil {
		return nil, err
	}
	option, err := config.LookupString(name)
	if err != nil {
		return nil, err
	}
	return &option, nil
}

func ConfiguredInt32(repoPath, name string, fallback int32) int32 {
	config, err := repoConfig(repoPath)
	if err != nil {
		return fallback
	}
	option, err := config.LookupInt32(name)
	if err != nil {
		return fallback
	}
	return option
}

func ConfiguredBool(repoPath, name string, fallback bool) bool {
	config, err := repoConfig(repoPath)
	if err != nil {
		return fallback
	}
	option, err := config.LookupBool(name)
	if err != nil {
		return fallback
	}
	return option
}

func repoConfig(repoPath string) (*git.Config, error) {
	repo, err := git.OpenRepository(repoPath)
	if err != nil {
		return nil, err
	}
	return repo.Config()
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func gitVariable(path, name, fallback string) *string {
	err := os.Chdir(path)
	if err != nil {
		return &fallback
	}
	cmd := exec.Command("git", "var", name)
	cmd.Env = os.Environ()
	output, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		return &fallback
	}
	variable := strings.TrimSpace(string(output))
	return &variable
}
