package libgitcomment

import (
	"github.com/stvp/assert"
	"testing"
)

func TestCreateFileRefNoLine(t *testing.T) {
	ref := CreateFileRef("pkg/src/example_item.ft:145", false)
	assert.NotNil(t, ref)
	assert.Equal(t, ref.Path, "pkg/src/example_item.ft")
	assert.Equal(t, ref.Line, 145)
	assert.Equal(t, ref.LineType, RefLineTypeNew)
}

func TestCreateFileRefDeletedLine(t *testing.T) {
	ref := CreateFileRef("pkg/src/example_item.ft:145:old", false)
	assert.NotNil(t, ref)
	assert.Equal(t, ref.Path, "pkg/src/example_item.ft")
	assert.Equal(t, ref.Line, 145)
	assert.Equal(t, ref.LineType, RefLineTypeOld)
}

func TestCreateFileRefDeletedLineOverride(t *testing.T) {
	ref := CreateFileRef("pkg/src/example_item.ft:145", true)
	assert.NotNil(t, ref)
	assert.Equal(t, ref.Path, "pkg/src/example_item.ft")
	assert.Equal(t, ref.Line, 145)
	assert.Equal(t, ref.LineType, RefLineTypeOld)
}

func TestCreateFileMultiColon(t *testing.T) {
	ref := CreateFileRef("pkg/src/item:other.txt:98", false)
	assert.NotNil(t, ref)
	assert.Equal(t, ref.Path, "pkg/src/item:other.txt")
	assert.Equal(t, ref.Line, 98)
	assert.Equal(t, ref.LineType, RefLineTypeNew)
}

func TestCreateFileNoLine(t *testing.T) {
	ref := CreateFileRef("pkg/src/item:other.txt", false)
	assert.NotNil(t, ref)
	assert.Equal(t, ref.Path, "pkg/src/item:other.txt")
	assert.Equal(t, ref.Line, 0)
	assert.Equal(t, ref.LineType, RefLineTypeNew)
}

func TestSerializeRefWithLine(t *testing.T) {
	ref := DeserializeFileRef("pkg/src/example_item.ft:34")
	assert.Equal(t, ref.Serialize(), "pkg/src/example_item.ft:34")
}

func TestSerializeRefWithoutLine(t *testing.T) {
	ref := DeserializeFileRef("pkg/src/example_item:other.txt")
	assert.Equal(t, ref.Serialize(), "pkg/src/example_item:other.txt")
}

func TestSerializeStaleCache(t *testing.T) {
	ref := DeserializeFileRef("src/example.txt")
	ref.Line = 5
	assert.Equal(t, ref.Serialize(), "src/example.txt:5")
}
