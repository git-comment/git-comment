package git_comment

import (
	"github.com/stvp/assert"
	"github.com/wayn3h0/go-uuid"
	"strings"
	"testing"
)

func TestNewCommentAuthor(t *testing.T) {
	author, _ := CreatePerson("Sam Wafers <sam@example.com>")
	comment, err := NewComment("", "123", nil, author)
	assert.Nil(t, err)
	assert.NotNil(t, comment)
	assert.Equal(t, comment.Author, author)
	assert.Equal(t, comment.Amender, author)
}

func TestNewCommentAmender(t *testing.T) {
	author, _ := CreatePerson("Sam Wafers <sam@example.com>")
	comment, err := NewComment("", "123", nil, author)
	assert.Nil(t, err)
	assert.NotNil(t, comment)
	assert.Equal(t, comment.Amender, author)
}

func TestNewCommentTime(t *testing.T) {
	comment, err := NewComment("", "123", nil, nil)
	assert.Nil(t, err)
	assert.NotNil(t, comment)
	assert.NotNil(t, comment.CreateTime)
	assert.Equal(t, comment.AmendTime, comment.CreateTime)
}

func TestNewCommentCommit(t *testing.T) {
	comment, err := NewComment("", "abcdefg", nil, nil)
	assert.Nil(t, err)
	assert.NotNil(t, comment)
	assert.Equal(t, comment.Commit, "abcdefg")
}

func TestNewCommentContent(t *testing.T) {
	comment, err := NewComment("Season the chex mix", "abcdefg", nil, nil)
	assert.Nil(t, err)
	assert.NotNil(t, comment)
	assert.Equal(t, comment.Content, "Season the chex mix")
}

func TestNewCommentID(t *testing.T) {
	comment, err := NewComment("", "abcdefg", nil, nil)
	assert.Nil(t, err)
	assert.NotNil(t, comment)
	assert.NotNil(t, comment.ID)
	identifier, uErr := uuid.Parse(comment.ID)
	assert.Nil(t, uErr)
	assert.NotNil(t, identifier)
}

func TestNewCommentDeleted(t *testing.T) {
	comment, err := NewComment("", "abcdefg", nil, nil)
	assert.Nil(t, err)
	assert.NotNil(t, comment)
	assert.False(t, comment.Deleted)
}

func TestNewCommentFileRef(t *testing.T) {
	ref := &FileRef{"src/example.c", 12}
	comment, err := NewComment("", "abcdefg", ref, nil)
	assert.Nil(t, err)
	assert.NotNil(t, comment)
	assert.Equal(t, comment.FileRef, ref)
}

func TestObjectContentFull(t *testing.T) {
	ref := &FileRef{"src/example.c", 12}
	author := &Person{"Selina Kyle", "cat@example.com"}
	comment, _ := NewComment("This line is too long", "acdacdacd", ref, author)
	assert.NotNil(t, comment)
	lines := strings.Split(comment.ObjectContent(), "\n")
	assert.Equal(t, len(lines), 8)
	assert.Equal(t, lines[0], "commit acdacdacd")
	assert.Equal(t, lines[2], "Selina Kyle <cat@example.com>")
	assert.Equal(t, lines[4], "Selina Kyle <cat@example.com>")
}
