package libgitcomment

import (
	"fmt"
	"github.com/kylef/result.go/src/result"
	git "gopkg.in/libgit2/git2go.v23"
	gg "libgitcomment/git"
)

type DiffLineType int

const AdditionalCommentsFile = "comments:"

const (
	DiffAdd DiffLineType = iota
	DiffAddNewline
	DiffRemove
	DiffRemoveNewline
	DiffContext
	DiffOther
	DiffUnassignedComments
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
func DiffCommits(repoPath, commitish string, contextLines uint32) result.Result {
	return gg.WithRepository(repoPath, func(repo *git.Repository) result.Result {
		commits := gg.ResolveCommits(repo, gg.ExpandCommitish(commitish))
		return commits.FlatMap(func(commitRange interface{}) result.Result {
			return diffCommits(repo, commitRange.(*gg.CommitRange), contextLines)
		})
	})
}

// @return result.Result<*Diff, error>
func diffCommits(repo *git.Repository, commitRange *gg.CommitRange, contextLines uint32) result.Result {
	comments := CommentsOnCommits(repo, commitRange.Commits())
	diff := diffRange(repo, commitRange, contextLines)
	return result.Combine(func(values ...interface{}) result.Result {
		parentID := commitRange.Parent.Id().String()
		childID := commitRange.Child.Id().String()
		files := parseDiffForLines(values[0].(*git.Diff), values[1].(CommentSlice))
		return result.NewSuccess(&Diff{files, parentID, childID})
	}, diff, comments)
}

func commitTree(commit *git.Commit) result.Result {
	return result.NewResult(commit.Tree())
}

func diffRange(repo *git.Repository, commitRange *gg.CommitRange, contextLines uint32) result.Result {
	return result.Combine(func(values ...interface{}) result.Result {
		opts := values[2].(git.DiffOptions)
		opts.ContextLines = contextLines
		return result.NewResult(repo.DiffTreeToTree(
			values[0].(*git.Tree),
			values[1].(*git.Tree),
			&opts))
	}, commitTree(commitRange.Parent), commitTree(commitRange.Child), diffOptions())
}

func parseDiffForLines(diff *git.Diff, comments CommentSlice) []*DiffFile {
	commentMapping := commentsByFileRef(comments)
	files := make([]*DiffFile, 0)
	var file *DiffFile
	var delta git.DiffDelta
	cbLine := func(line git.DiffLine) error {
		file.Lines = append(file.Lines, &DiffLine{
			diffTypeFromLine(line),
			line.Content,
			line.OldLineno,
			line.NewLineno,
			commentsForLine(delta, line, commentMapping),
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
	file = fileForAdditionalComments(commentMapping)
	if file != nil {
		files = append(files, file)
	}
	return files
}

func fileForAdditionalComments(mapping map[string][]*Comment) *DiffFile {
	var comments []*Comment
	if list, ok := mapping[AdditionalCommentsFile]; ok {
		comments = list
	} else {
		return nil
	}

	return &DiffFile{AdditionalCommentsFile, "",
		[]*DiffLine{&DiffLine{
			DiffUnassignedComments,
			"",
			-1,
			-1,
			comments,
		}}}
}

func commentsForLine(delta git.DiffDelta, line git.DiffLine, mapping map[string][]*Comment) []*Comment {
	if line.Origin == git.DiffLineDeletion {
		return removedCommentForLine(delta, line, mapping)
	}
	return addedCommentsForLine(delta, line, mapping)
}

func removedCommentForLine(delta git.DiffDelta, line git.DiffLine, mapping map[string][]*Comment) []*Comment {
	var comments []*Comment = nil
	oldCommentKey := fileRefMappingKey(delta.OldFile.Path, line.OldLineno)
	if list, ok := mapping[oldCommentKey]; ok {
		for _, comment := range list {
			if comment.FileRef.LineType == RefLineTypeOld {
				comments = append(comments, comment)
			}
		}
	}
	return comments
}

func addedCommentsForLine(delta git.DiffDelta, line git.DiffLine, mapping map[string][]*Comment) []*Comment {
	var comments []*Comment = nil
	newCommentKey := fileRefMappingKey(delta.NewFile.Path, line.NewLineno)
	if list, ok := mapping[newCommentKey]; ok {
		for _, comment := range list {
			if comment.FileRef.LineType == RefLineTypeNew {
				comments = append(comments, comment)
			}
		}
	}
	return comments
}

func commentsByFileRef(comments CommentSlice) map[string][]*Comment {
	mapping := make(map[string][]*Comment)
	for _, comment := range comments {
		ref := comment.FileRef
		var key string
		if ref != nil && len(ref.Path) > 0 && ref.Line > 0 {
			key = fileRefMappingKey(ref.Path, ref.Line)
		} else {
			key = AdditionalCommentsFile
		}
		if list, ok := mapping[key]; ok {
			mapping[key] = append(list, comment)
		} else {
			list := make([]*Comment, 0)
			mapping[key] = append(list, comment)
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
func diffOptions() result.Result {
	return result.NewResult(git.DefaultDiffOptions())
}
