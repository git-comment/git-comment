package git_comment

// package git

import (
	"fmt"
	git "gopkg.in/libgit2/git2go.v22"
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
func DiffCommits(repoPath string, commitish string) (*Diff, error) {
	repo, err := git.OpenRepository(repoPath)
	if err != nil {
		return nil, err
	}
	parent, child, err := ResolveCommits(repo, commitish)
	if err != nil {
		return nil, err
	}
	if child == nil {
		child = parent
		parent = child.Parent(0)
	}
	return diffCommits(repo, parent, child)
}

func diffCommits(repo *git.Repository, parent *git.Commit, child *git.Commit) (*Diff, error) {
	commitTree, err := child.Tree()
	if err != nil {
		return nil, err
	}
	parentTree, err := parent.Tree()
	if err != nil {
		return nil, err
	}
	opts, err := defaultDiffOptions()
	if err != nil {
		return nil, err
	}
	diff, err := repo.DiffTreeToTree(parentTree, commitTree, opts)
	if err != nil {
		return nil, err
	}
	commits := CommitsFromRange(parent, child)
	comments, err := CommentsOnCommits(repo, commits)
	if err != nil {
		return nil, err
	}
	files := parseDiffForLines(diff, comments)
	return &Diff{files, parent.Id().String(), child.Id().String()}, nil
}

func parseDiffForLines(diff *git.Diff, comments []*Comment) []*DiffFile {
	commentMapping := commentsByFileRef(comments)
	files := make([]*DiffFile, 0)
	cbFile := func(delta git.DiffDelta, progress float64) (git.DiffForEachHunkCallback, error) {
		lines := make([]*DiffLine, 0)
		file := &DiffFile{delta.OldFile.Path, delta.NewFile.Path, lines}
		files = append(files, file)
		return func(hunk git.DiffHunk) (git.DiffForEachLineCallback, error) {
			return func(line git.DiffLine) error {
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
			}, nil
		}, nil
	}
	diff.ForEach(cbFile, git.DiffDetailLines)
	return files
}

func commentsByFileRef(comments []*Comment) map[string][]*Comment {
	mapping := make(map[string][]*Comment)
	for _, comment := range comments {
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

func defaultDiffOptions() (*git.DiffOptions, error) {
	opts, err := git.DefaultDiffOptions()
	if err != nil {
		return nil, err
	}
	// TODO: set context lines from config
	return &opts, nil
}
