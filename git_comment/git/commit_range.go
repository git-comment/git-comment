package git

import (
	git "github.com/libgit2/git2go"
)

type CommitRange struct {
	Parent *git.Commit
	Child  *git.Commit
}

// Find all intermediate commits between a parent and child
func (c *CommitRange) Commits() []*git.Commit {
	commits := make([]*git.Commit, 0)
	if c.Child != nil && c.Parent != nil {
		commit := c.Child
		for {
			if commit == nil {
				break
			}
			commits = append(commits, commit)
			if commit.Id().String() == c.Parent.Id().String() {
				break
			}
			commit = commit.Parent(0)
		}
	} else if c.Parent != nil {
		commits = append(commits, c.Parent)
	} else if c.Child != nil {
		commits = append(commits, c.Child)
	}
	return commits
}
