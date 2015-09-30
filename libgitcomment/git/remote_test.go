package git

import (
	"testing"

	"github.com/stvp/assert"
)

func TestContains(t *testing.T) {
	assert.True(t, contains([]string{"a", "b", "c"}, "b"))
	assert.False(t, contains([]string{"a", "b", "c"}, "d"))
	assert.False(t, contains([]string{"a", "b", "c"}, "A"))
	assert.False(t, contains([]string{"a", "b", "c"}, "ab"))
	assert.False(t, contains([]string{}, "d"))
}
