package git_comment

import (
	"errors"
	"fmt"
	"time"
)

type Comment struct {
	Author     *Person
	CreateTime time.Time
	Content    string
	Amender    *Person
	AmendTime  time.Time
	Commit     string
	ID         *string
	Deleted    bool
	FileRef    *FileRef
}

// Creates a new comment using provided content and author
func NewComment(message string, commit string, fileRef *FileRef, author *Person) (*Comment, error) {
	const missingContentMessage = "No message content provided"
	const missingCommitMessage = "No commit provided"
	if len(message) == 0 {
		return nil, errors.New(missingContentMessage)
	} else if len(commit) == 0 {
		return nil, errors.New(missingCommitMessage)
	}
	createTime := time.Now()
	return &Comment{
		author,
		createTime,
		message,
		author,
		createTime,
		commit,
		nil,
		false,
		fileRef,
	}, nil
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
	blob.Properties.Set("commit", c.Commit)
	blob.Properties.Set("file", c.FileRef.Serialize())
	blob.Properties.Set("author", c.Author.Serialize())
	blob.Properties.Set("created", fmt.Sprintf("%d %v", c.CreateTime.Unix(), c.CreateTime.Format("-0700")))
	blob.Properties.Set("amender", c.Amender.Serialize())
	blob.Properties.Set("amended", fmt.Sprintf("%d %v", c.AmendTime.Unix(), c.CreateTime.Format("-0700")))
	if c.Deleted {
		blob.Properties.Set("deleted", "true")
	} else {
		blob.Message = c.Content
	}
	return blob.Serialize()
}
