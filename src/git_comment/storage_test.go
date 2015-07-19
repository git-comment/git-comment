package git_comment

import (
	"github.com/stvp/assert"
	"path"
	"testing"
)

func TestCreateComment(t *testing.T) {

}

func TestCreateInvalidComment(t *testing.T) {

}

func TestDeleteComment(t *testing.T) {

}

func TestDeleteNonexistentComment(t *testing.T) {

}

func TestRefPath(t *testing.T) {
	commit := "0155eb4229851634a0f03eb265b69f5a2d56f341"
	id := "23caf9710a71e3736597415c57bdcf5eebae6bcb"
	comment, _ := NewComment("Unsure of the intent here.",
		commit, new(FileRef), new(Person)).Dematerialize()
	p, err := RefPath(comment.(*Comment), id).Dematerialize()
	assert.Nil(t, err)
	expected := path.Join("refs/comments", "0155",
		"eb4229851634a0f03eb265b69f5a2d56f341", id)
	assert.Equal(t, p, expected)
}
