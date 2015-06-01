package git_comment

import (
	"github.com/stvp/assert"
	"testing"
)

func TestCreateFileRefFull(t *testing.T) {
	ref := CreateFileRef("pkg/src/example_item.ft:145")
	assert.NotNil(t, ref)
	assert.Equal(t, ref.Path, "pkg/src/example_item.ft")
	assert.Equal(t, ref.Line, 145)
}

func TestCreateFileMultiColon(t *testing.T) {
	ref := CreateFileRef("pkg/src/item:other.txt:98")
	assert.NotNil(t, ref)
	assert.Equal(t, ref.Path, "pkg/src/item:other.txt")
	assert.Equal(t, ref.Line, 98)
}

func TestCreateFileNoLine(t *testing.T) {
	ref := CreateFileRef("pkg/src/item:other.txt")
	assert.NotNil(t, ref)
	assert.Equal(t, ref.Path, "pkg/src/item:other.txt")
	assert.Equal(t, ref.Line, 0)
}

func TestSerializeRefWithLine(t *testing.T) {
	data := "pkg/src/example_item.ft:34"
	ref := CreateFileRef(data)
	assert.Equal(t, ref.Serialize(), data)
}

func TestSerializeRefWithoutLine(t *testing.T) {
	data := "pkg/src/example_item:other.txt"
	ref := CreateFileRef(data)
	assert.Equal(t, ref.Serialize(), data)
}

func TestSerializeStaleCache(t *testing.T) {
	ref := CreateFileRef("src/example.txt")
	ref.Line = 5
	assert.Equal(t, ref.Serialize(), "src/example.txt:5")
}
