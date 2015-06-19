package git_comment

import (
	"errors"
	"fmt"
	git "gopkg.in/libgit2/git2go.v22"
	"path"
	"time"
)

const (
	DefaultMessageTemplate = "\n# Enter comment content\n# Lines beginning with '#' will be stripped"
	authorNotFoundError    = "No name or email found in git config for commenting"
	invalidHashError       = "Invalid commit hash for storage"
	commitNotFoundError    = "Commit not found"
	commentNotFoundError   = "Comment not found"
	headCommit             = "HEAD"
	defaultMessageFormat   = "Created a comment ref on [%v] to [%v]"
)

// Create a new comment on a commit, optionally with a file and line
func CreateComment(repoPath string, commit *string, fileRef *FileRef, message string) (*string, error) {
	repo, err := git.OpenRepository(repoPath)
	if err != nil {
		return nil, err
	}
	author, err := ConfiguredAuthor(repo)
	if err != nil {
		return nil, err
	}
	hash, err := parseCommit(repo, commit)
	if err != nil {
		return nil, err
	}
	comment, err := NewComment(message, *hash, fileRef, author)
	if err != nil {
		return nil, err
	}
	if err := writeCommentToDisk(repo, comment); err != nil {
		return nil, err
	}

	return comment.ID, nil
}

// Update an existing comment with a new message
func UpdateComment(repoPath string, ID string, message string) (*string, error) {
	repo, err := git.OpenRepository(repoPath)
	if err != nil {
		return nil, err
	}
	comment, err := CommentByID(repo, ID)
	if err != nil {
		return nil, err
	}
	committer, err := ConfiguredCommitter(repo)
	if err != nil {
		return nil, err
	}
	comment.Amend(message, committer)
	if err := writeCommentToDisk(repo, comment); err != nil {
		return nil, err
	}

	return comment.ID, nil
}

// Remove a comment from a commit
func DeleteComment(repoPath string, ID string) error {
	repo, err := git.OpenRepository(repoPath)
	if err != nil {
		return err
	}
	comment, err := CommentByID(repo, ID)
	if err != nil {
		return err
	}
	comment.Deleted = true
	return writeCommentToDisk(repo, comment)
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

// Finds all comments on a given commit
func CommentsOnCommit(repoPath string, commit *string) ([]*Comment, error) {
	const glob = "*"
	var comments []*Comment
	repo, err := git.OpenRepository(repoPath)
	if err != nil {
		return nil, err
	}
	parsedCommit, err := parseCommit(repo, commit)
	if err != nil {
		return nil, err
	}
	dir, err := commitRefDir(parsedCommit)
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

// Write git object for a given comment and update the
// comment refs
func writeCommentToDisk(repo *git.Repository, comment *Comment) error {
	if comment.ID != nil {
		err := deleteReference(repo, comment, *comment.ID)
		if err != nil {
			return err
		}
	}
	oid, err := repo.CreateBlobFromBuffer([]byte(comment.Serialize()))
	if err != nil {
		return err
	}
	committer := comment.Amender
	sig := &git.Signature{committer.Name, committer.Email, time.Now()}
	id := fmt.Sprintf("%v", oid)
	file, err := refPath(comment, &id)
	if err != nil {
		return err
	}
	commit := *comment.Commit
	message := fmt.Sprintf(defaultMessageFormat, commit[:7], id[:7])
	_, err = repo.CreateReference(*file, oid, false, sig, message)
	if err != nil {
		return err
	}
	comment.ID = &id
	return nil
}

func deleteReference(repo *git.Repository, comment *Comment, identifier string) error {
	refPath, err := refPath(comment, &identifier)
	if err != nil {
		return nil
	}
	ref, err := repo.LookupReference(*refPath)
	if err != nil {
		return errors.New(commentNotFoundError)
	}
	err = ref.Delete()
	if err != nil {
		return err
	}
	return nil
}

// Generate the path within refs for a given comment
//
// Comment refs are nested under refs/comments. The
// format is as follows:
//
// ```
// refs/comments/[<commit prefix>]/[<rest of commit>]/[<comment id>]
// ```
//
func refPath(comment *Comment, id *string) (*string, error) {
	dir, err := commitRefDir(comment.Commit)
	if err != nil {
		return nil, err
	}
	hash := path.Join(*dir, *id)
	return &hash, nil
}

func commitRefDir(commit *string) (*string, error) {
	const commentPath = "refs/comments"
	hash := *commit
	if len(hash) > 4 {
		dir := path.Join(commentPath,
			hash[:4],
			hash[4:len(hash)])
		return &dir, nil
	}
	return nil, errors.New(invalidHashError)
}

// parse a commit hash, converting to the HEAD commit where needed
func parseCommit(repo *git.Repository, commit *string) (*string, error) {
	var hash string
	var id string
	if commit == nil || len(*commit) == 0 {
		hash = headCommit
	} else {
		hash = *commit
	}
	ref, err := repo.LookupReference(hash)
	if err != nil {
		oid, err := git.NewOid(hash)
		if err != nil {
			return nil, errors.New(commitNotFoundError)
		}
		obj, err := repo.Lookup(oid)
		if err != nil {
			return nil, errors.New(commitNotFoundError)
		}
		id = obj.Id().String()
		return nil, errors.New(commitNotFoundError)
	}
	res, err := ref.Resolve()
	if err != nil {
		return nil, err
	}
	id = res.Target().String()
	return &id, nil
}

func head(repo *git.Repository) (*string, error) {
	head, hErr := repo.Head()
	if hErr != nil {
		return nil, hErr
	}
	hash := head.Name()
	return &hash, nil
}
