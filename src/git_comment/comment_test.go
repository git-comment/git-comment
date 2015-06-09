package git_comment

import (
	"github.com/stvp/assert"
	"strings"
	"testing"
	"time"
)

func TestNewCommentAuthor(t *testing.T) {
	author := CreatePerson("Sam Wafers <sam@example.com>")
	comment, err := NewComment("Curious decision here.", "123", nil, author)
	assert.Nil(t, err)
	assert.NotNil(t, comment)
	assert.Equal(t, comment.Author, author)
	assert.Equal(t, comment.Amender, author)
}

func TestNewCommentAmender(t *testing.T) {
	author := CreatePerson("Sam Wafers <sam@example.com>")
	comment, err := NewComment("Doesn't this violate the laws of physics?", "123", nil, author)
	assert.Nil(t, err)
	assert.NotNil(t, comment)
	assert.Equal(t, comment.Amender, author)
}

func TestNewCommentTime(t *testing.T) {
	comment, err := NewComment("ELI5?", "123", nil, nil)
	assert.Nil(t, err)
	assert.NotNil(t, comment)
	assert.NotNil(t, comment.CreateTime)
	assert.Equal(t, comment.AmendTime, comment.CreateTime)
}

func TestNewCommentCommit(t *testing.T) {
	comment, err := NewComment("Wat?", "abcdefg", nil, nil)
	assert.Nil(t, err)
	assert.NotNil(t, comment)
	assert.Equal(t, *comment.Commit, "abcdefg")
}

func TestNewCommentContent(t *testing.T) {
	comment, err := NewComment("Season the chex mix", "abcdefg", nil, nil)
	assert.Nil(t, err)
	assert.NotNil(t, comment)
	assert.Equal(t, comment.Content, "Season the chex mix")
}

func TestNewCommentID(t *testing.T) {
	comment, err := NewComment("This behavior is undocumented", "abcdefg", nil, nil)
	assert.Nil(t, err)
	assert.NotNil(t, comment)
	assert.Nil(t, comment.ID)
}

func TestNewCommentDeleted(t *testing.T) {
	comment, err := NewComment("What is happening here?", "abcdefg", nil, nil)
	assert.Nil(t, err)
	assert.NotNil(t, comment)
	assert.False(t, comment.Deleted)
}

func TestNewCommentFileRef(t *testing.T) {
	ref := &FileRef{"src/example.c", 12}
	comment, err := NewComment("This should be more modular", "abcdefg", ref, nil)
	assert.Nil(t, err)
	assert.NotNil(t, comment)
	assert.Equal(t, comment.FileRef, ref)
}

func TestCreateWithoutContent(t *testing.T) {
	_, err := NewComment("", "azerty", new(FileRef), new(Person))
	assert.NotNil(t, err)
}

func TestSerializeComment(t *testing.T) {
	ref := &FileRef{"src/example.c", 12}
	author := &Person{"Selina Kyle", "cat@example.com"}
	comment, _ := NewComment("This line is too long", "acdacdacd", ref, author)
	lines := strings.Split(comment.Serialize(), "\n")
	assert.Equal(t, len(lines), 8)
	assert.Equal(t, lines[0], "commit acdacdacd")
	assert.Equal(t, lines[1], "file src/example.c:12")
	assert.Equal(t, lines[2], "author Selina Kyle <cat@example.com>")
	assert.Equal(t, lines[4], "amender Selina Kyle <cat@example.com>")
	assert.Equal(t, lines[6], "")
	assert.Equal(t, lines[7], "This line is too long")
}

func TestSerializeDeletedComment(t *testing.T) {
	author := &Person{"Morpheus", "redpill@example.com"}
	comment, _ := NewComment("Pick one", "afdafdafd", new(FileRef), author)
	comment.Deleted = true
	lines := strings.Split(comment.Serialize(), "\n")
	assert.Equal(t, len(lines), 8)
	assert.Equal(t, lines[0], "commit afdafdafd")
	assert.Equal(t, lines[1], "file ")
	assert.Equal(t, lines[2], "author Morpheus <redpill@example.com>")
	assert.Equal(t, lines[4], "amender Morpheus <redpill@example.com>")
	assert.Equal(t, lines[6], "deleted true")
	assert.Equal(t, lines[7], "")
}

func TestDeserializeComment(t *testing.T) {
	author := &Person{"Morpheus", "redpill@example.com"}
	comment, _ := NewComment("Pick one", "afdafdafd", CreateFileRef("bin/exec:15"), author)
	newComment, err := DeserializeComment(comment.Serialize())
	assert.Nil(t, err)
	assert.Equal(t, *comment.Commit, *newComment.Commit)
	assert.Equal(t, *comment.FileRef, *newComment.FileRef)
	assert.Equal(t, *comment.Author, *newComment.Author)
	assert.Equal(t, *comment.Amender, *newComment.Amender)
	assert.Equal(t, comment.Content, newComment.Content)
	assert.Equal(t, comment.CreateTime.Format(time.RFC822Z), newComment.CreateTime.Format(time.RFC822Z))
	assert.Equal(t, comment.AmendTime.Format(time.RFC822Z), newComment.AmendTime.Format(time.RFC822Z))
}
