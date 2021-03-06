package libgitcomment

import (
	"github.com/stvp/assert"
	"sort"
	"strings"
	"testing"
	"time"
)

func TestSortComments(t *testing.T) {
	var comments = make(CommentSlice, 3)
	comments[0] = &Comment{Author: &Person{Date: time.Now().Add(12 * time.Hour)}}
	comments[1] = &Comment{Author: &Person{Date: time.Now()}}
	comments[2] = &Comment{Author: &Person{Date: time.Now().Add(24 * time.Hour)}}
	sortedComments := []*Comment{comments[1], comments[0], comments[2]}
	sort.Stable(comments)
	for idx, comment := range comments {
		assert.Equal(t, comment, sortedComments[idx])
	}
}

func TestNewCommentAuthor(t *testing.T) {
	author := &Person{"Sam Wafers", "<sam@example.com>", time.Now(), "-0600"}
	c, err := NewComment("Curious decision here.", "123", nil, author).Dematerialize()
	comment := c.(*Comment)
	assert.Nil(t, err)
	assert.NotNil(t, comment)
	assert.Equal(t, comment.Author, author)
	assert.Equal(t, comment.Amender, author)
}

func TestNewCommentAmender(t *testing.T) {
	author := &Person{"Sam Wafers", "<sam@example.com>", time.Now(), "-0600"}
	c, err := NewComment("Doesn't this violate the laws of physics?", "123", nil, author).Dematerialize()
	comment := c.(*Comment)
	assert.Nil(t, err)
	assert.NotNil(t, comment)
	assert.Equal(t, comment.Amender, author)
}

func TestNewCommentTime(t *testing.T) {
	var unix int64 = 1433220431
	c, err := NewComment("ELI5?", "123", nil, &Person{Date: time.Unix(unix, 0)}).Dematerialize()
	comment := c.(*Comment)
	assert.Nil(t, err)
	assert.NotNil(t, comment.Author)
	assert.Equal(t, comment.Author.Date.Unix(), unix)
	assert.NotNil(t, comment.Amender)
	assert.Equal(t, comment.Amender.Date, comment.Author.Date)
}

func TestNewCommentCommit(t *testing.T) {
	c, err := NewComment("Wat?", "abcdefg", nil, nil).Dematerialize()
	comment := c.(*Comment)
	assert.Nil(t, err)
	assert.NotNil(t, comment)
	assert.Equal(t, *comment.Commit, "abcdefg")
}

func TestNewCommentContent(t *testing.T) {
	c, err := NewComment("Season the chex mix", "abcdefg", nil, nil).Dematerialize()
	comment := c.(*Comment)
	assert.Nil(t, err)
	assert.NotNil(t, comment)
	assert.Equal(t, comment.Content, "Season the chex mix")
}

func TestNewCommentID(t *testing.T) {
	c, err := NewComment("This behavior is undocumented", "abcdefg", nil, nil).Dematerialize()
	comment := c.(*Comment)
	assert.Nil(t, err)
	assert.NotNil(t, comment)
	assert.Nil(t, comment.ID)
}

func TestNewCommentDeleted(t *testing.T) {
	c, err := NewComment("What is happening here?", "abcdefg", nil, nil).Dematerialize()
	comment := c.(*Comment)
	assert.Nil(t, err)
	assert.NotNil(t, comment)
	assert.False(t, comment.Deleted)
}

func TestNewCommentFileRef(t *testing.T) {
	ref := &FileRef{"src/example.c", 12, RefLineTypeNew}
	c, err := NewComment("This should be more modular", "abcdefg", ref, nil).Dematerialize()
	comment := c.(*Comment)
	assert.Nil(t, err)
	assert.NotNil(t, comment)
	assert.Equal(t, comment.FileRef, ref)
}

func TestCreateWithoutContent(t *testing.T) {
	_, err := NewComment("", "azerty", new(FileRef), new(Person)).Dematerialize()
	assert.NotNil(t, err)
}

func TestSerializeComment(t *testing.T) {
	ref := &FileRef{"src/example.c", 12, RefLineTypeNew}
	author := &Person{"Selina Kyle", "cat@example.com", time.Unix(1437498360, 0), "+1100"}
	c, _ := NewComment("This line is too long", "acdacdacd", ref, author).Dematerialize()
	comment := c.(*Comment)
	lines := strings.Split(comment.Serialize(), "\n")
	assert.Equal(t, len(lines), 6)
	assert.Equal(t, lines[0], "commit acdacdacd")
	assert.Equal(t, lines[1], "file src/example.c:12")
	assert.Equal(t, lines[2], "author Selina Kyle <cat@example.com> 1437498360 +1100")
	assert.Equal(t, lines[3], "amender Selina Kyle <cat@example.com> 1437498360 +1100")
	assert.Equal(t, lines[4], "")
	assert.Equal(t, lines[5], "This line is too long")
}

func TestSerializeDeletedComment(t *testing.T) {
	author := &Person{"Morpheus", "redpill@example.com", time.Unix(1437498360, 0), "-0600"}
	c, _ := NewComment("Pick one", "afdafdafd", new(FileRef), author).Dematerialize()
	comment := c.(*Comment)
	comment.Deleted = true
	lines := strings.Split(comment.Serialize(), "\n")
	assert.Equal(t, len(lines), 6)
	assert.Equal(t, lines[0], "commit afdafdafd")
	assert.Equal(t, lines[1], "file ")
	assert.Equal(t, lines[2], "author Morpheus <redpill@example.com> 1437498360 -0600")
	assert.Equal(t, lines[3], "amender Morpheus <redpill@example.com> 1437498360 -0600")
	assert.Equal(t, lines[4], "deleted true")
	assert.Equal(t, lines[5], "")
}

func TestDeserializeComment(t *testing.T) {
	author := &Person{"Morpheus", "redpill@example.com", time.Unix(1437498360, 0), "-0600"}
	c, _ := NewComment("Pick one", "afdafdafd", DeserializeFileRef("bin/exec:15"), author).Dematerialize()
	comment := c.(*Comment)
	newC, err := DeserializeComment(comment.Serialize()).Dematerialize()
	newComment := newC.(*Comment)
	assert.Nil(t, err)
	assert.Equal(t, *comment.Commit, *newComment.Commit)
	assert.Equal(t, *comment.FileRef, *newComment.FileRef)
	assert.Equal(t, *comment.Author, *newComment.Author)
	assert.Equal(t, *comment.Amender, *newComment.Amender)
	assert.Equal(t, comment.Content, newComment.Content)
}
