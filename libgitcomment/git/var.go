package git

import (
	"os"
	"os/exec"
	"strings"
)

type variable string

const (
	gitEditor    variable = "GIT_EDITOR"
	gitPager     variable = "GIT_PAGER"
	gitAuthor    variable = "GIT_AUTHOR_IDENT"
	gitCommitter variable = "GIT_COMMITTER_IDENT"
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
	return gitVariable(repoPath, gitEditor, defaultEditor)
}

// The text viewer to use for viewing text interactively, as
// configured through git-var(1)
func ConfiguredPager(repoPath string) string {
	return gitVariable(repoPath, gitPager, defaultPager)
}

// The author of a piece of code as configured through git-var(1)
func ConfiguredAuthor(repoPath string) string {
	return gitVariable(repoPath, gitAuthor, "")
}

// The committer of a piece of code as configured through git-var(1)
func ConfiguredCommitter(repoPath string) string {
	return gitVariable(repoPath, gitCommitter, "")
}

func gitVariable(repoPath string, name variable, fallback string) string {
	if env := os.Getenv(string(name)); len(env) > 0 {
		return env
	}
	if err := os.Chdir(repoPath); err != nil {
		return fallback
	}
	cmd := exec.Command("git", "var", string(name))
	cmd.Env = os.Environ()
	output, err := cmd.Output()
	if err != nil {
		return fallback
	}
	return strings.TrimSpace(string(output))
}
