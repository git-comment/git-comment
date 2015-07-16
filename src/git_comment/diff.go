package git_comment

import (
	"fmt"
	"github.com/kylef/result.go/src/result"
	git "github.com/libgit2/git2go"
)

type DiffLineType int

const (
	DiffAdd DiffLineType = iota
	DiffAddNewline
	DiffRemove
	DiffRemoveNewline
	DiffContext
	DiffOther
)

type Diff struct {
	Files      []*DiffFile
	FromCommit string
	ToCommit   string
}

type DiffFile struct {
	OldPath string
	NewPath string
	Lines   []*DiffLine
}

type DiffLine struct {
	Type          DiffLineType
	Content       string
	OldLineNumber int
	NewLineNumber int
	Comments      []*Comment
}

// Find diffs on given commits
//
// If commitish resolves to a single commit, the diff is performed
// between the commit and its parent.
// @return result.Result<*Diff, error>
func DiffCommits(repoPath string, commitish string) result.Result {
	return WithRepository(repoPath, func(repo *git.Repository) result.Result {
		return ResolveCommits(repo, commitish).FlatMap(func(commitRange interface{}) result.Result {
			return diffCommits(repo, commitRange.(*CommitRange))
		})
	})
}

// @return result.Result<*Diff, error>
func diffCommits(repo *git.Repository, commitRange *CommitRange) result.Result {
	comments := CommentsOnCommits(repo, commitRange.Commits())
	diff := diffRange(repo, commitRange)
	return result.Combine(func(values ...interface{}) result.Result {
		parentID := commitRange.Parent.Id().String()
		childID := commitRange.Child.Id().String()
		files := parseDiffForLines(values[0].(*git.Diff), values[1].([]interface{}))
		return result.NewSuccess(&Diff{files, parentID, childID})
	}, diff, comments)
}

func commitTree(commit *git.Commit) result.Result {
	return result.NewResult(commit.Tree())
}

func diffRange(repo *git.Repository, commitRange *CommitRange) result.Result {
	return result.Combine(func(values ...interface{}) result.Result {
		opts := values[2].(git.DiffOptions)
		return result.NewResult(repo.DiffTreeToTree(
			values[0].(*git.Tree),
			values[1].(*git.Tree),
			&opts))
	}, commitTree(commitRange.Parent), commitTree(commitRange.Child), defaultDiffOptions())
}

func parseDiffForLines(diff *git.Diff, comments []interface{}) []*DiffFile {
	commentMapping := commentsByFileRef(comments)
	files := make([]*DiffFile, 0)
	var file *DiffFile
	var delta git.DiffDelta
	cbLine := func(line git.DiffLine) error {
		var comments []*Comment = nil
		commentKey := fileRefMappingKey(delta.NewFile.Path, line.NewLineno)
		if list, ok := commentMapping[commentKey]; ok {
			comments = list
		}
		file.Lines = append(file.Lines, &DiffLine{
			diffTypeFromLine(line),
			line.Content,
			line.OldLineno,
			line.NewLineno,
			comments,
		})
		return nil
	}
	cbHunk := func(hunk git.DiffHunk) (git.DiffForEachLineCallback, error) {
		return cbLine, nil
	}
	cbFile := func(diffDelta git.DiffDelta, progress float64) (git.DiffForEachHunkCallback, error) {
		delta = diffDelta
		lines := make([]*DiffLine, 0)
		file = &DiffFile{delta.OldFile.Path, delta.NewFile.Path, lines}
		files = append(files, file)
		return cbHunk, nil
	}
	diff.ForEach(cbFile, git.DiffDetailLines)
	return files
}

func commentsByFileRef(comments []interface{}) map[string][]*Comment {
	mapping := make(map[string][]*Comment)
	for _, c := range comments {
		comment := c.(*Comment)
		ref := comment.FileRef
		if ref != nil && len(ref.Path) > 0 && ref.Line > 0 {
			key := fileRefMappingKey(ref.Path, ref.Line)
			if list, ok := mapping[key]; ok {
				mapping[key] = append(list, comment)
			} else {
				list := make([]*Comment, 0)
				mapping[key] = append(list, comment)
			}
		}
	}
	return mapping
}

func fileRefMappingKey(path string, line int) string {
	return fmt.Sprintf("%v:%d", path, line)
}

func diffTypeFromLine(line git.DiffLine) DiffLineType {
	switch line.Origin {
	case git.DiffLineContext, git.DiffLineContextEOFNL:
		return DiffContext
	case git.DiffLineAddition:
		return DiffAdd
	case git.DiffLineDeletion:
		return DiffRemove
	case git.DiffLineAddEOFNL:
		return DiffAddNewline
	case git.DiffLineDelEOFNL:
		return DiffRemoveNewline
	default:
		return DiffOther
	}
}

// @return result.Result<git.DiffOptions, error>
func defaultDiffOptions() result.Result {
	return result.NewResult(git.DefaultDiffOptions())
}
