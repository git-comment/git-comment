package git_comment

import (
  "github.com/wayn3h0/go-uuid/random"
  "github.com/wayn3h0/go-uuid"
  "time"
)

type FileRef struct {
  Path string
  Line int
}

type Comment struct {
  Author     *Person
  CreateTime time.Time
  Content    string
  Amender    *Person
  AmendTime  time.Time
  Commit     string
  ID         string
  Deleted    bool
  FileRef    *FileRef
}

// Creates a new comment using provided content and author
func NewComment(message string, commit string, fileRef *FileRef, author *Person) (*Comment, error) {
  id, idErr := random.New()
  if idErr != nil {
    return nil, idErr
  }
  createTime := time.Now()
  return &Comment{
    author,
    createTime,
    message,
    author,
    createTime,
    commit,
    id.Format(uuid.StyleWithoutDash),
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
//   1243040974 -0900
//   amender Delisa Mason <name@example.com>
//   1243040974 -0900
//
//   Too many levels of indentation here.
// ```
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
