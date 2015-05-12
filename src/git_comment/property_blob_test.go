package git_comment

import (
	"github.com/stvp/assert"
	"strings"
	"testing"
)

func TestPropertiesFromBlob(t *testing.T) {
	content := "fruit Green Apple\ntree_type Redwood\n\nThe Redwood is the tallest tree in North America"
	blob := CreatePropertyBlob(content)
	assert.Equal(t, "Green Apple", blob.Properties["fruit"])
	assert.Equal(t, "Redwood", blob.Properties["tree_type"])
}

func TestMessageFromBlob(t *testing.T) {
	content := "fruit Green Apple\ntree_type Redwood\n\nThe Redwood is the tallest tree in North America"
	blob := CreatePropertyBlob(content)
	assert.Equal(t, "The Redwood is the tallest tree in North America", blob.Message)
}

func TestNewlinesInMessage(t *testing.T) {
	props := "author Elira <elira@example.com>\nfile src/example.txt\n\n"
	message := "I have a few questions.\nHow do we plan on handling the latter case?\n\nWhere can I get some chili dogs?"
	content := strings.Join([]string{props, message}, "")
	blob := CreatePropertyBlob(content)
	assert.Equal(t, message, blob.Message)
}

func TestEmptyMessage(t *testing.T) {
	content := "item1 some stuff\nitem2 some other stuff"
	blob := CreatePropertyBlob(content)
	assert.Equal(t, "", blob.Message)
}

func TestRawContent(t *testing.T) {
	message := "I have a few questions.\nHow do we plan on handling the latter case?\n\nWhere can I get some chili dogs?"
	blob := &PropertyBlob{
		map[string]string{
			"author": "Elira <elira@example.com>",
			"file":   "src/example.txt",
		},
		message,
	}
	props := "author Elira <elira@example.com>\nfile src/example.txt\n\n"
	content := strings.Join([]string{props, message}, "")
	assert.Equal(t, content, blob.RawContent())
}
