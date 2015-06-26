package git_comment

import (
	"errors"
	git "gopkg.in/libgit2/git2go.v22"
	"path"
)

func CommentsWithContent(content string) []*Comment {
	return nil
}

// Finds a comment by a given ID
func CommentByID(repo *git.Repository, identifier string) (*Comment, error) {
	oid, err := git.NewOid(identifier)
	if err != nil {
		return nil, errors.New(commentNotFoundError)
	}
	blob, err := repo.LookupBlob(oid)
	if err != nil {
		return nil, errors.New(commentNotFoundError)
	}
	comment, err := DeserializeComment(string(blob.Contents()))
	if err != nil {
		return nil, err
	}
	comment.ID = &identifier
	return comment, nil
}

func CommentsOnCommittish(repoPath string, committish string) ([]*Comment, error) {
	repo, err := git.OpenRepository(repoPath)
	if err != nil {
		return nil, err
	}
	parentCommit, childCommit, err := ResolveCommits(repo, committish)
	if err != nil {
		return nil, err
	}
	if parentCommit != nil && childCommit != nil {
		return CommentsOnCommits(repo, CommitsFromRange(parentCommit, childCommit))
	} else if parentCommit != nil {
		return commentsOnCommit(repo, parentCommit)
	}
	return commentsOnCommit(repo, childCommit)
}

// Finds all comments on a given commit
func CommentsOnCommit(repoPath string, commit *string) ([]*Comment, error) {
	repo, err := git.OpenRepository(repoPath)
	if err != nil {
		return nil, err
	}
	parsedCommit, err := ResolveSingleCommitHash(repo, commit)
	if err != nil {
		return nil, err
	}
	commitObj, err := repo.LookupCommit(git.NewOidFromBytes([]byte(*parsedCommit)))
	if err != nil {
		return nil, err
	}
	return commentsOnCommit(repo, commitObj)
}

func CommentsOnCommits(repo *git.Repository, commits []*git.Commit) ([]*Comment, error) {
	comments := make([]*Comment, 0)
	for _, commit := range commits {
		commitComments, err := commentsOnCommit(repo, commit)
		if err != nil {
			return nil, err
		}
		for _, comment := range commitComments {
			comments = append(comments, comment)
		}
	}
	return comments, nil
}

func commentsOnCommit(repo *git.Repository, commit *git.Commit) ([]*Comment, error) {
	const glob = "*"
	var comments []*Comment
	id := commit.Id().String()
	dir, err := CommitRefDir(&id)
	if err != nil {
		return nil, err
	}
	refIterator, err := repo.NewReferenceIteratorGlob(path.Join(*dir, glob))
	if err != nil {
		return nil, err
	}
	ref, err := refIterator.Next()
	for {
		if err != nil {
			if err.(*git.GitError).Code == git.ErrIterOver {
				break
			} else {
				return nil, err
			}
		}
		comment := commentFromRef(repo, ref.Name())
		if comment != nil {
			comments = append(comments, comment)
		}
		ref, err = refIterator.Next()
	}
	return comments, nil
}

func commentFromRef(repo *git.Repository, refName string) *Comment {
	_, identifier := path.Split(refName)
	comment, err := CommentByID(repo, identifier)
	if err != nil {
		return nil
	}
	return comment
}
