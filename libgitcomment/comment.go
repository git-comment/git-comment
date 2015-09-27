package libgitcomment

import (
	"errors"
	"github.com/kylef/result.go/src/result"
	"strings"
	"time"
)

type Comment struct {
	Author  *Person
	Content string
	Amender *Person
	Commit  *string
	ID      *string
	Deleted bool
	FileRef *FileRef
}

const timeFormat string = time.RFC822Z

type CommentSlice []*Comment

const (
	authorKey  = "author"
	commitKey  = "commit"
	amenderKey = "amender"
	fileRefKey = "file"
	deletedKey = "deleted"
)

func (cs CommentSlice) Len() int {
	return len(cs)
}

func (cs CommentSlice) Less(i, j int) bool {
	return cs[i].Author.Date.Before(cs[j].Author.Date)
}

func (cs CommentSlice) Swap(i, j int) {
	cs[i], cs[j] = cs[j], cs[i]
}

// Creates a new comment using provided content and author
func NewComment(message string, commit string, fileRef *FileRef, author *Person) result.Result {
	const missingContentMessage = "No message content provided"
	const missingCommitMessage = "No commit provided"
	if len(message) == 0 {
		return result.NewFailure(errors.New(missingContentMessage))
	} else if len(commit) == 0 {
		return result.NewFailure(errors.New(missingCommitMessage))
	}
	return result.NewSuccess(&Comment{
		author,
		message,
		author,
		&commit,
		nil,
		false,
		fileRef,
	})
}

func DeserializeComment(content string) result.Result {
	const serializationErrorMessage = "Could not deserialize object into comment"
	blob := CreatePropertyBlob(content)
	comment := &Comment{}
	comment.Content = blob.Message
	comment.Commit = blob.Get(commitKey)
	comment.Author = blob.GetPerson(authorKey)
	comment.Amender = blob.GetPerson(amenderKey)
	comment.FileRef = blob.GetFileRef(fileRefKey)
	return result.NewSuccess(comment)
}

// First line of the comment content
func (c *Comment) Title() string {
	return strings.Split(c.Content, "\n")[0]
}

// Update the message content of the comment
func (c *Comment) Amend(message string, amender *Person) {
	c.Content = message
	c.Amender = amender
}

// Generate content of git object for comment
// Comment ref file format:
//
// ```
//   commit 0155eb4229851634a0f03eb265b69f5a2d56f341
//   file src/example.txt:12
//   author Delisa Mason <name@example.com>
//   created 1243040974 -0900
//   amender Delisa Mason <name@example.com>
//   amended 1243040974 -0900
//
//   Too many levels of indentation here.
// ```
//
func (c *Comment) Serialize() string {
	blob := NewPropertyBlob()
	blob.Set(commitKey, *c.Commit)
	blob.Set(fileRefKey, c.FileRef.Serialize())
	blob.Set(authorKey, c.Author.Serialize())
	blob.Set(amenderKey, c.Amender.Serialize())
	if c.Deleted {
		blob.Set(deletedKey, "true")
	} else {
		blob.Message = c.Content
	}
	return blob.Serialize()
}
