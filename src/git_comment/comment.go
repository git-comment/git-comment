package git_comment

import (
	"errors"
	"github.com/kylef/result.go/src/result"
	"time"
)

type Comment struct {
	Author     *Person
	CreateTime time.Time
	Content    string
	Amender    *Person
	AmendTime  time.Time
	Commit     *string
	ID         *string
	Deleted    bool
	FileRef    *FileRef
}

const timeFormat string = time.RFC822Z

type CommentSlice []*Comment

const (
	authorKey  = "author"
	commitKey  = "commit"
	createdKey = "created"
	amenderKey = "amender"
	amendedKey = "amended"
	fileRefKey = "file"
	deletedKey = "deleted"
)

func (cs CommentSlice) Len() int {
	return len(cs)
}

func (cs CommentSlice) Less(i, j int) bool {
	return cs[i].CreateTime.Before(cs[j].CreateTime)
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
	createTime := time.Now()
	return result.NewSuccess(&Comment{
		author,
		createTime,
		message,
		author,
		createTime,
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
	cTime := blob.GetTime(createdKey)
	if cTime == nil {
		return result.NewFailure(errors.New(serializationErrorMessage))
	}
	comment.CreateTime = *cTime
	aTime := blob.GetTime(amendedKey)
	if aTime == nil {
		return result.NewFailure(errors.New(serializationErrorMessage))
	}
	comment.AmendTime = *aTime
	return result.NewSuccess(comment)
}

// Update the message content of the comment
func (c *Comment) Amend(message string, amender *Person) {
	c.Content = message
	c.Amender = amender
	c.AmendTime = time.Now()
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
	blob.Set(createdKey, c.CreateTime.Format(timeFormat))
	blob.Set(amenderKey, c.Amender.Serialize())
	blob.Set(amendedKey, c.CreateTime.Format(timeFormat))
	if c.Deleted {
		blob.Set(deletedKey, "true")
	} else {
		blob.Message = c.Content
	}
	return blob.Serialize()
}
