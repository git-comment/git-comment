package git_comment

import (
	"github.com/stvp/assert"
	"testing"
)

func TestContains(t *testing.T) {
	assert.True(t, contains([]string{"a", "b", "c"}, "b"))
	assert.False(t, contains([]string{"a", "b", "c"}, "d"))
	assert.False(t, contains([]string{"a", "b", "c"}, "A"))
	assert.False(t, contains([]string{"a", "b", "c"}, "ab"))
	assert.False(t, contains([]string{}, "d"))
}
