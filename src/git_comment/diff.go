package git_comment

import (
	"fmt"
	git "gopkg.in/libgit2/git2go.v22"
)

type Diff struct {
	FileRef     *FileRef
	LinesBefore string
	LinesAfter  string
}

// Create a representation of the changed, added, and/or removed
// lines at a given file ref on a commit
func DiffLines(repoPath string, commit string, fileRef *FileRef) *Diff {
	if fileRef == nil || fileRef.Line == 0 {
		return nil
	}
	repo, err := git.OpenRepository(repoPath)
	if err != nil {
		return nil
	}
	// find commit object
	obj, err := repo.RevparseSingle(commit)
	if err != nil {
		return nil
	}
	commitObj, err := repo.LookupCommit(obj.Id())
	if err != nil {
		return nil
	}
	// find commit parent
	parent := commitObj.Parent(0)
	if parent == nil {
		return nil
	}
	// create diff from tree
	commitTree, err := commitObj.Tree()
	if err != nil {
		return nil
	}
	parentTree, err := parent.Tree()
	if err != nil {
		return nil
	}
	opts, err := git.DefaultDiffOptions()
	if err != nil {
		return nil
	}
	var linesBefore = []byte{}
	var linesAfter = []byte{}
	opts.Pathspec = []string{fileRef.Path}
	diff, err := repo.DiffTreeToTree(parentTree, commitTree, &opts)
	if err != nil {
		return nil
	}
	cbLine := func(line git.DiffLine) error {
		formatted := []byte(formattedLine(line))
		if line.OldLineno > fileRef.Line || line.NewLineno > fileRef.Line {
			linesAfter = append(linesAfter, formatted...)
		} else {
			linesBefore = append(linesBefore, formatted...)
		}
		return nil
	}
	cbHunk := func(hunk git.DiffHunk) (git.DiffForEachLineCallback, error) {
		return cbLine, nil
	}
	cbFile := func(delta git.DiffDelta, progress float64) (git.DiffForEachHunkCallback, error) {
		return cbHunk, nil
	}
	diff.ForEach(cbFile, git.DiffDetailLines)
	// iterate over the diff deltas to find the correct file
	// grab lines at the fileref lines, or all if none provided
	// return in structure
	return &Diff{fileRef, string(linesBefore), string(linesAfter)}
}

func formattedLine(line git.DiffLine) string {
	switch line.Origin {
	case git.DiffLineContext:
		return line.Content
	case git.DiffLineAddition:
		return fmt.Sprintf("+ %v", line.Content)
	case git.DiffLineDeletion:
		return fmt.Sprintf("- %v", line.Content)
	default:
		return line.Content
	}
}
