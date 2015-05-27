package git_comment

import (
	"github.com/stvp/assert"
	"strings"
	"testing"
	"github.com/cevaris/ordered_map"
)

func TestPropertiesFromBlob(t *testing.T) {
	content := "fruit Green Apple\ntree_type Redwood\n\nThe Redwood is the tallest tree in North America"
	blob := CreatePropertyBlob(content)
	fruit, _ := blob.Properties.Get("fruit")
	tree, _ := blob.Properties.Get("tree_type")
	assert.Equal(t, "Green Apple", fruit)
	assert.Equal(t, "Redwood", tree)
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
	propertyMap := ordered_map.NewOrderedMap()
	propertyMap.Set("author", "Elira <elira@example.com>")
	propertyMap.Set("file", "src/example.txt")
	blob := &PropertyBlob{
		propertyMap,
		message,
	}
	props := "author Elira <elira@example.com>\nfile src/example.txt\n\n"
	content := strings.Join([]string{props, message}, "")
	assert.Equal(t, content, blob.RawContent())
}
