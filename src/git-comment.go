package main

import (
	"git_comment"
	"github.com/wayn3h0/go-uuid"
	"time"
	// "gopkg.in/libgit2/git2go.v22"
)

type Identifier uuid.UUID

type Person struct {
	Name  string
	Email string
}

type FileRef struct {
	Path string
	Line int
}

type Comment struct {
	Author     Person
	CreateTime time.Time
	Content    string
	Amender    Person
	AmendTime  time.Time
	Commit     Identifier
	ID         Identifier
}

// Finds a comment by a given ID
func CommentByID(ID Identifier) *Comment {
	return &Comment{}
}

// Finds all comments on a given commit
func CommentsOnCommit(commit Identifier) []*Comment {
	return []*Comment{}
}

// Creates a new comment using provided content and author
func NewComment(message string, author *Person) *Comment {
	return &Comment{}
}

// Write git object for a given comment and update the
// comment refs
func WriteCommentToDisk(comment *Comment) {

}

// Generate content of git object for comment
// Comment ref file format:
//
//   commit 0155eb4229851634a0f03eb265b69f5a2d56f341
//   file src/example.txt:12
//   author Delisa Mason <name@example.com>
//   1243040974 -0900
//   amender Delisa Mason <name@example.com>
//   1243040974 -0900
//
//   Too many levels of indentation here.
//
func (comment *Comment) ObjectContent() string {
	return ""
}

// Generate the path within refs for a given comment
//
// Comment refs are nested under refs/comments. The
// directory name is the first four characters of the
// commit identifier, and the file name are the
// remaining characters. The contents of the file are
// the identifiers of all comments on the commit
func (comment *Comment) RefPath() string {
	return ""
}

func main() {

}
