package git_comment

import (
	"errors"
	"github.com/kylef/result.go/src/result"
	"os"
	"os/exec"
	"strings"
)

const (
	defaultEditor          = "vi"
	defaultPager           = "less"
	authorNotFoundError    = "No name or email found in git config for commenting"
	committerNotFoundError = "No name or email found in git config for creating a comment"
)

// The editor to use for editing comments interactively, as
// configured through git-var(1)
func ConfiguredEditor(repoPath string) string {
	return gitVariable(repoPath, "GIT_EDITOR", defaultEditor)
}

// The text viewer to use for viewing text interactively, as
// configured through git-var(1)
func ConfiguredPager(repoPath string) string {
	return gitVariable(repoPath, "GIT_PAGER", defaultPager)
}

// The author of a piece of code as configured through git-var(1)
func ConfiguredAuthor(repoPath string) result.Result {
	defaultError := result.NewFailure(errors.New(authorNotFoundError))
	return CreatePerson(gitVariable(repoPath, "GIT_AUTHOR_IDENT", "")).RecoverWith(defaultError)
}

// The committer of a piece of code as configured through git-var(1)
func ConfiguredCommitter(repoPath string) result.Result {
	defaultError := result.NewFailure(errors.New(committerNotFoundError))
	return CreatePerson(gitVariable(repoPath, "GIT_COMMITTER_IDENT", "")).RecoverWith(defaultError)
}

func gitVariable(repoPath, name, fallback string) string {
	if err := os.Chdir(repoPath); err != nil {
		return fallback
	}
	cmd := exec.Command("git", "var", name)
	cmd.Env = os.Environ()
	output, err := cmd.Output()
	if err != nil {
		return fallback
	}
	return strings.TrimSpace(string(output))
}
