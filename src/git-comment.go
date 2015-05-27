package main

import (
	"git_comment"
	"time"
	"gopkg.in/libgit2/git2go.v22"
)

// Finds a comment by a given ID
func CommentByID(repoPath string, ID Identifier) *Comment {
  return &Comment{}
}

// Finds all comments on a given commit
func CommentsOnCommit(repoPath string, commit Identifier) []*Comment {
  return []*Comment{}
}

func CreateComment(repoPath string, commit Identifier, fileRef FileRef, message string, author *Person) error {
	comment, err := git_comment.NewComment(message, commit, fileRef, author)
	if comment {
		return WriteCommentToDisk(repoPath, comment)
	}
	return err
}

func DeleteComment(repoPath string, ID Identifier, amender *Person) {

}

// Write git object for a given comment and update the
// comment refs
func WriteCommentToDisk(repoPath string, comment *Comment) error {

}

func repoWithPath(repoPath string) (*git2go.Repository, error) {
	return git2go.OpenRepository(repoPath);
}

func main() {

}
